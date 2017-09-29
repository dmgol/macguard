package mb

import (
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
