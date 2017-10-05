package main

import (
	"encoding/json"
	"log"

	"github.com/dmgol/macguard/agent/snmp"
	"github.com/dmgol/macguard/mb"
	"github.com/dmgol/macguard/utils"
)

type ReplyFunc func(reply []byte, correlationId, replyTo string)
type ConsumeFunc func(msgs mb.DeliveryChan, reply ReplyFunc)

func declareAndConsumeQueue(mb *mb.MessageBus, queueName string, consumeFunc ConsumeFunc) {
	q, err := mb.DeclareQueue(queueName)
	utils.FailOnError(err, "Failed to declare a queue")

	msgs, err := q.Consume()
	go consumeFunc(msgs, func(reply []byte, correlationId, replyTo string) {
		q.Reply(reply, correlationId, replyTo)
	})
	utils.FailOnError(err, "Failed to register a consumer")
}

func main() {

	bus, err := mb.Connect("amqp://guest:guest@10.44.32.99:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer bus.Close()

	forever := make(chan bool)

	declareAndConsumeQueue(bus, "mac_address_table", func(msgs mb.DeliveryChan, reply ReplyFunc) {
		for d := range msgs {
			var params snmp.SnmpParams
			json.Unmarshal(d.Body, &params)
			log.Printf("Received a message[mac_address_table]: %v", params)
			macAddrTable, _ := snmp.GetMacAddrTable(params)
			log.Printf("Mac address table for %s (%d total)\n", params.Target, len(macAddrTable))
			result, _ := json.Marshal(macAddrTable)
			reply(result, d.CorrelationId, d.ReplyTo)
		}
	})
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
