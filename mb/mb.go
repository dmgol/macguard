package mb

import (
	"encoding/json"
	"log"

	"github.com/dmgol/macguard/snmp"
	"github.com/dmgol/macguard/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/streadway/amqp"
)

type MessageBus struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

type Queue struct {
	q  *amqp.Queue
	ch *amqp.Channel
}

type DeliveryChan <-chan amqp.Delivery

func (me *Queue) Consume() (DeliveryChan, error) {
	return me.ch.Consume(
		me.q.Name, // queue
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

}

func (me *Queue) Name() string {
	return me.q.Name
}

func (me *Queue) Publish(body []byte) error {
	return me.ch.Publish(
		"",        // exchange
		me.q.Name, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			Body: body,
		})
}

func (me *Queue) Call(body []byte, correlationId, replyTo string) error {
	return me.ch.Publish(
		"",        // exchange
		me.q.Name, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			Body:          body,
			CorrelationId: correlationId,
			ReplyTo:       replyTo,
		})
}

func (me *Queue) Reply(body []byte, correlationId, replyTo string) error {
	return me.ch.Publish(
		"",      // exchange
		replyTo, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			Body:          body,
			CorrelationId: correlationId,
		})
}

func (me *MessageBus) DeclareQueue(name string) (*Queue, error) {
	q, err := me.ch.QueueDeclare(
		name,  // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}
	return &Queue{q: &q, ch: me.ch}, nil
}

func (me *MessageBus) DeclareCallbackQueue() (*Queue, error) {
	q, err := me.ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}
	return &Queue{q: &q, ch: me.ch}, nil
}

func Connect(uri string) (*MessageBus, error) {
	var (
		err error
		me  MessageBus
	)
	me.conn, err = amqp.Dial(uri)
	if err != nil {
		return nil, err
	}
	me.ch, err = me.conn.Channel()
	if err != nil {
		me.Close()
		return nil, err
	}
	err = me.ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		me.Close()
		return nil, err
	}
	return &me, nil
}

func (me *MessageBus) Close() {
	if me.ch != nil {
		me.ch.Close()
	}
	if me.conn != nil {
		me.conn.Close()
	}
}

type ReplyFunc func(params snmp.Params) []byte

func (me *MessageBus) DeclareAndConsumeSnmpQueue(queueName string, replyFunc ReplyFunc) error {
	q, err := me.DeclareQueue(queueName)
	utils.FailOnError(err, "Failed to declare a queue")

	msgs, err := q.Consume()
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			var params snmp.Params
			json.Unmarshal(d.Body, &params)
			log.Printf("Received a message[%s]: %v", queueName, params)
			result := replyFunc(params)
			q.Reply(result, d.CorrelationId, d.ReplyTo)
		}
	}()

	return nil
}

type CallbackFunc func(data []byte) error

func (me *MessageBus) CallSmtpQueue(queueName string, params snmp.Params, callbackFunc CallbackFunc) error {
	q, err := me.DeclareQueue(queueName)
	if err != nil {
		return err
	}

	cbQueue, err := me.DeclareCallbackQueue()
	if err != nil {
		return err
	}

	msgs, err := cbQueue.Consume()
	if err != nil {
		return err
	}

	corrId := uuid.NewV4().String()

	body, err := json.Marshal(params)
	if err != nil {
		return err
	}

	err = q.Call(body, corrId, cbQueue.Name())
	if err != nil {
		return err
	}

	for d := range msgs {
		if d.CorrelationId == corrId {
			err := callbackFunc(d.Body)
			if err != nil {
				return err
			}
			break
		}
	}
	return nil
}
