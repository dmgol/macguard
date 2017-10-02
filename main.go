package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/dmgol/macguard/agent/snmp"
	"github.com/dmgol/macguard/mb"
	"github.com/dmgol/macguard/utils"
)

func main() {
	//const targetIP = "10.102.7.2"
	// macAddressTable, _ := snmp.GetMACAddressTable(net.ParseIP(targetIP), "test@5")
	// fmt.Printf("Mac address table for %s (%d total)\n", targetIP, len(macAddressTable))
	// for _, row := range macAddressTable {
	// 	fmt.Println(row)
	// }
	// vlans := snmp.TestGetCiscoVlanList(net.ParseIP(targetIP), "test")
	// for _, v := range vlans {
	// 	fmt.Println(v)
	// }

	mb, err := mb.Connect("amqp://guest:guest@10.44.32.99:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer mb.Close()

	q, err := mb.DeclareQueue("mac_address_table")
	utils.FailOnError(err, "Failed to declare a queue")

	msgs, err := q.Consume()
	utils.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var r snmp.MacAddressTableRequest
			json.Unmarshal(d.Body, &r)
			log.Printf("Received a message: %v", r)
			macAddressTable, _ := snmp.GetMacAddressTable(net.ParseIP(r.Target), r.Community)
			log.Printf("Mac address table for %s (%d total)\n", r.Target, len(macAddressTable))
			reply, _ := json.Marshal(macAddressTable)
			q.Reply(reply, d.CorrelationId, d.ReplyTo)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
