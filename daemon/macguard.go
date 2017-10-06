package main

import (
	"encoding/json"
	"fmt"

	"github.com/dmgol/macguard/mb"
	"github.com/dmgol/macguard/models"
	"github.com/dmgol/macguard/snmp"
	"github.com/dmgol/macguard/utils"
	"github.com/markbates/pop"
	uuid "github.com/satori/go.uuid"
)

func createSwitches(db *pop.Connection) error {
	err := db.Transaction(func(tx *pop.Connection) error {
		var community = models.Community{Value: "test"}
		var hp2510Model = models.SwitchModel{Name: "HP Procurve 2510", Vendor: "HP"}
		var hp2510List = models.Switches{
			models.Switch{Name: "hp1", Location: "location1", IpAddr: "10.102.7.60"},
			models.Switch{Name: "hp2", Location: "location2", IpAddr: "10.102.7.62"},
		}
		var cisco4948Model = models.SwitchModel{Name: "Cisco Catalyst 4948", Vendor: "Cisco"}
		var cisco4948List = models.Switches{
			models.Switch{Name: "cisco1", Location: "location1", IpAddr: "10.102.7.1"},
			models.Switch{Name: "cisco2", Location: "location2", IpAddr: "10.102.7.2"},
		}

		err := tx.Save(&community)
		if err != nil {
			return err
		}

		err = tx.Save(&hp2510Model)
		if err != nil {
			return err
		}

		err = tx.Save(&cisco4948Model)
		if err != nil {
			return err
		}

		for _, sw := range cisco4948List {
			sw.ModelID = cisco4948Model.ID
			sw.CommunityID = community.ID
			err = tx.Save(&sw)
			if err != nil {
				return err
			}
		}

		for _, sw := range hp2510List {
			sw.ModelID = hp2510Model.ID
			sw.CommunityID = community.ID
			err = tx.Save(&sw)
			if err != nil {
				return err
			}
		}
		return nil

	})
	return err
}

func getVlanList(tx *pop.Connection, bus *mb.MessageBus, params snmp.Params, switchID uuid.UUID) error {
	err := bus.CallSmtpQueue("snmp_vlan_list", params, func(data []byte) error {
		var vlans models.Vlans
		if err := json.Unmarshal(data, &vlans); err != nil {
			return err
		}

		deleteQuery := tx.RawQuery("DELETE FROM vlans WHERE switch_id = ?", switchID)
		if err := deleteQuery.Exec(); err != nil {
			return err
		}

		for _, v := range vlans {
			v.SwitchID = switchID
			if err := tx.Save(&v); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func getPortList(tx *pop.Connection, bus *mb.MessageBus, params snmp.Params, switchID uuid.UUID) error {
	err := bus.CallSmtpQueue("snmp_port_list", params, func(data []byte) error {
		var ports models.Ports
		if err := json.Unmarshal(data, &ports); err != nil {
			return err
		}

		deleteQuery := tx.RawQuery("DELETE FROM ports WHERE switch_id = ?", switchID)
		if err := deleteQuery.Exec(); err != nil {
			return err
		}

		for _, p := range ports {
			p.SwitchID = switchID
			if err := tx.Save(&p); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func loadVlansAndPorts(db *pop.Connection, bus *mb.MessageBus) error {
	err := db.Transaction(func(tx *pop.Connection) error {
		var allSwithes models.Switches
		tx.All(&allSwithes)

		for _, sw := range allSwithes {
			var model models.SwitchModel
			if err := tx.Find(&model, sw.ModelID); err != nil {
				return err
			}
			var community models.Community
			if err := tx.Find(&community, sw.CommunityID); err != nil {
				return err
			}
			fmt.Println(sw)
			if err := getVlanList(tx,
				bus,
				snmp.Params{Target: sw.IpAddr, Community: community.Value, Vendor: model.Vendor},
				sw.ID); err != nil {
				return err
			}
			if err := getPortList(tx,
				bus,
				snmp.Params{Target: sw.IpAddr, Community: community.Value, Vendor: model.Vendor},
				sw.ID); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func getMacAddrTableCisco(tx *pop.Connection, bus *mb.MessageBus, params snmp.Params, switchID uuid.UUID) error {
	var vlans models.Vlans

	var table = models.MacAddrTable{SwitchID: switchID}
	if err := tx.Create(&table); err != nil {
		return err
	}

	if err := tx.Where("switch_id = ?", switchID).All(&vlans); err != nil {
		return err
	}

	community := params.Community

	for _, v := range vlans {
		params.Community = fmt.Sprintf("%s@%d", community, v.Number)
		err := bus.CallSmtpQueue("snmp_mac_addr_table", params, func(data []byte) error {
			var entries models.MacAddrTableEntries
			if err := json.Unmarshal(data, &entries); err != nil {
				return err
			}
			for _, e := range entries {
				e.TableID = table.ID
				e.PortID = findPortIDByNumber(tx, switchID, e.PortNumber)
				e.VlanID = v.ID
				if err := tx.Save(&e); err != nil {
					return err
				}
			}
			return nil
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func findPortIDByNumber(tx *pop.Connection, switchID uuid.UUID, number int) uuid.UUID {
	var port models.Port

	tx.Where("number = ?", number).First(&port)
	return port.ID
}

func getMacAddrTable(tx *pop.Connection, bus *mb.MessageBus, params snmp.Params, switchID uuid.UUID) error {
	err := bus.CallSmtpQueue("snmp_mac_addr_table", params, func(data []byte) error {
		var entries models.MacAddrTableEntries
		if err := json.Unmarshal(data, &entries); err != nil {
			return err
		}

		var table = models.MacAddrTable{SwitchID: switchID}
		if err := tx.Create(&table); err != nil {
			return err
		}

		for _, e := range entries {
			e.TableID = table.ID
			e.PortID = findPortIDByNumber(tx, switchID, e.PortNumber)
			if err := tx.Save(&e); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

type getMacAddrTableFunc func(tx *pop.Connection, bus *mb.MessageBus, params snmp.Params, switchID uuid.UUID) error

func loadMacAddrTables(db *pop.Connection, bus *mb.MessageBus) error {
	err := db.Transaction(func(tx *pop.Connection) error {
		var allSwithes models.Switches
		tx.All(&allSwithes)

		for _, sw := range allSwithes {
			var model models.SwitchModel
			if err := tx.Find(&model, sw.ModelID); err != nil {
				return err
			}
			var community models.Community
			if err := tx.Find(&community, sw.CommunityID); err != nil {
				return err
			}
			fmt.Println(sw)

			var getFunc getMacAddrTableFunc

			switch model.Vendor {
			case "Cisco":
				getFunc = getMacAddrTableCisco

			default:
				getFunc = getMacAddrTable
			}

			if err := getFunc(tx,
				bus,
				snmp.Params{Target: sw.IpAddr, Community: community.Value, Vendor: model.Vendor},
				sw.ID); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func main() {
	db, err := pop.Connect("development")
	utils.FailOnError(err, "Failed to connect to database")
	defer db.Close()

	bus, err := mb.Connect("amqp://guest:guest@10.44.32.99:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer bus.Close()

	// utils.FailOnError(createSwitches(db), "Failed to create swithes")

	//utils.FailOnError(loadVlansAndPorts(db, bus), "Failed to load vlans and ports")

	utils.FailOnError(loadMacAddrTables(db, bus), "Failed to load macAddrTables")
}
