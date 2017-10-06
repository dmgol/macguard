package main

import (
	"encoding/json"
	"log"

	"github.com/dmgol/macguard/mb"
	"github.com/dmgol/macguard/snmp"
	"github.com/dmgol/macguard/utils"
)

func main() {

	bus, err := mb.Connect("amqp://guest:guest@10.44.32.99:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer bus.Close()

	forever := make(chan bool)

	err = bus.DeclareAndConsumeSnmpQueue("snmp_mac_addr_table", func(params snmp.Params) []byte {
		macAddrTable, _ := snmp.GetMacAddrTable(params)
		log.Printf("Mac address table for %s (%d total)\n", params.Target, len(macAddrTable))
		result, _ := json.Marshal(macAddrTable)
		return result
	})
	utils.FailOnError(err, "Failed to declare the queue 'snmp_mac_addr_table'")

	err = bus.DeclareAndConsumeSnmpQueue("snmp_vlan_list", func(params snmp.Params) []byte {
		vlans, _ := snmp.GetVlanList(params)
		log.Printf("Vlan list for %s (%d total)\n", params.Target, len(vlans))
		result, _ := json.Marshal(vlans)
		return result
	})
	utils.FailOnError(err, "Failed to declare the queue 'snmp_port_list'")

	err = bus.DeclareAndConsumeSnmpQueue("snmp_port_list", func(params snmp.Params) []byte {
		ports, _ := snmp.GetPortList(params)
		log.Printf("Port list for %s (%d total)\n", params.Target, len(ports))
		result, _ := json.Marshal(ports)
		return result
	})
	utils.FailOnError(err, "Failed to declare the queue 'snmp_port_list'")

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
