package main

import (
	"log"

	"github.com/dmgol/macguard/models"
	"github.com/markbates/pop"
)

func main() {
	db, err := pop.Connect("development")
	defer db.Close()
	if err != nil {
		log.Panic(err)
	}
	// db.Transaction(func(tx *pop.Connection) error {
	// 	d := &models.MacAddrTableEntry{MacAddr: "002255-ad8f0d", PortNumber: 1}
	// 	tx.Save(d)
	// 	log.Println(d)
	// 	return nil
	// })
	var d models.MacAddrTableEntry
	err = db.Find(&d, "9ce17b7a-0fc3-42fa-a52f-6922403017fb")
	if err != nil {
		log.Panic(err)
	}
	log.Println(d)

	// var users models.Users
	// err = db.Where("name  LIKE '%est%'").All(&users)
	// if err != nil {
	// 	log.Panic(err)
	// }
	// log.Println(users)

}
