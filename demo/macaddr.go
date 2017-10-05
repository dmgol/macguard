package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/dmgol/macguard/agent/snmp"
	"github.com/dmgol/macguard/mb"
	"github.com/dmgol/macguard/models"
	"github.com/dmgol/macguard/utils"
	"github.com/markbates/pop"
	uuid "github.com/satori/go.uuid"
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

func main() {
	db, err := pop.Connect("development")
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	mb, err := mb.Connect("amqp://guest:guest@10.44.32.99:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer mb.Close()

	q, err := mb.DeclareQueue("mac_address_table")
	utils.FailOnError(err, "Failed to declare a queue")

	callback, err := mb.DeclareCallbackQueue()
	utils.FailOnError(err, "Failed to declare a callback queue")

	msgs, err := callback.Consume()
	utils.FailOnError(err, "Failed to consume a callback queue")

	corrId := uuid.NewV4().String()

	req := snmp.SnmpParams{Target: "10.102.7.61", Community: "test"}
	body, err := json.Marshal(req)
	utils.FailOnError(err, "Failed to marshal json")
	err = q.Call(body, corrId, callback.Name())
	utils.FailOnError(err, "Failed to declare a queue")

	for d := range msgs {
		if d.CorrelationId == corrId {
			var data models.MacAddrTableEntries
			err = json.Unmarshal(d.Body, &data)
			utils.FailOnError(err, "Failed to unmarshall a reply")
			fmt.Println(data)
			writeToDb(data, db)
			break
		}
	}
}
