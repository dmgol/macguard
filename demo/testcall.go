package main

import (
	"encoding/json"
	"fmt"

	"github.com/dmgol/macguard/mb"
	"github.com/dmgol/macguard/models"
	"github.com/dmgol/macguard/snmp"
	"github.com/dmgol/macguard/utils"
	"github.com/markbates/pop"
)

func writeToDb(macAddrTable models.MacAddrTableEntries, db *pop.Connection) error {
	db.Transaction(func(tx *pop.Connection) error {
		var table models.MacAddrTable
		if err := tx.Save(&table); err != nil {
			return err
		}
		for _, entry := range macAddrTable {
			entry.TableID = table.ID
			if err := tx.Save(&entry); err != nil {
				return err
			}
		}
		return nil
	})
	return nil
}

func getMacAddrTable(bus *mb.MessageBus, db *pop.Connection) {
	bus.CallSmtpQueue("snmp_mac_addr_table", snmp.Params{Target: "10.102.7.2", Community: "test@5"}, func(data []byte) error {
		var table models.MacAddrTableEntries
		err := json.Unmarshal(data, &table)
		if err != nil {
			return err
		}
		fmt.Println(table)
		writeToDb(table, db)
		return nil
	})
}

func getVlanList(bus *mb.MessageBus, db *pop.Connection) {
	bus.CallSmtpQueue("snmp_vlan_list", snmp.Params{Target: "10.102.7.2", Community: "test", Vendor: "Cisco"}, func(data []byte) error {
		var vlans models.Vlans
		err := json.Unmarshal(data, &vlans)
		if err != nil {
			return err
		}
		fmt.Println(vlans)
		return nil
	})
}

func getPortList(bus *mb.MessageBus, db *pop.Connection) {
	bus.CallSmtpQueue("snmp_port_list", snmp.Params{Target: "10.102.7.2", Community: "test"}, func(data []byte) error {
		var ports models.Ports
		err := json.Unmarshal(data, &ports)
		if err != nil {
			return err
		}
		fmt.Println(ports)
		return nil
	})
}

func main() {
	db, err := pop.Connect("development")
	utils.FailOnError(err, "Failed to connect to database")
	defer db.Close()

	bus, err := mb.Connect("amqp://guest:guest@10.44.32.99:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer bus.Close()

	//getMacAddrTable(bus, db)
	getVlanList(bus, db)
	//getPortList(bus, db)
}
