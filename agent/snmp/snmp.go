package snmp

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dmgol/macguard/models"
	"github.com/soniah/gosnmp"
)

type Params struct {
	Target    string
	Community string
}

func snmpConnect(target string, community string) (*gosnmp.GoSNMP, error) {
	result := &gosnmp.GoSNMP{
		Target:    target,
		Port:      161,
		Community: community,
		Version:   gosnmp.Version2c,
		Timeout:   time.Duration(10 * time.Second),
		Retries:   3,
		MaxOids:   gosnmp.MaxOids,
	}
	return result, result.Connect()
}

func getLastNumberFromPduName(name string) int {
	if last := strings.LastIndex(name, "."); last != -1 {
		result, _ := strconv.Atoi(name[last+1:])
		return result
	}
	return -1
}

const macAddressOctetCount = 6

func extractMacAddrFromPdu(pdu *gosnmp.SnmpPDU) string {
	split := strings.Split(pdu.Name, ".")
	macAddrStart := len(split) - macAddressOctetCount
	hexStrings := make([]string, macAddressOctetCount)
	for i := 0; i < macAddressOctetCount; i++ {
		v, _ := strconv.Atoi(split[i+macAddrStart])
		hexStrings[i] = strconv.FormatInt(int64(v), 16)
	}
	return strings.Join(hexStrings, ":")
}

func pduToMacAddrTableEntry(pdu *gosnmp.SnmpPDU, portNumMap map[int]int) models.MacAddrTableEntry {
	result := models.MacAddrTableEntry{}
	result.MacAddr = extractMacAddrFromPdu(pdu)
	if pdu.Type == gosnmp.Integer {
		n := pdu.Value.(int)
		result.PortNumber = portNumMap[n]
	}
	return result
}

const portIfIndexOid = ".1.3.6.1.2.1.17.1.4.1.2"

func getPortNumberMap(session *gosnmp.GoSNMP) (map[int]int, error) {
	pduList, err := session.BulkWalkAll(portIfIndexOid)
	if err != nil {
		return nil, fmt.Errorf("bulk walk err: %v", err)
	}
	result := make(map[int]int)
	for _, pdu := range pduList {
		result[getLastNumberFromPduName(pdu.Name)] = pdu.Value.(int)
	}
	return result, nil
}

const macAddressTableOid = ".1.3.6.1.2.1.17.4.3.1.2"

func GetMacAddrTable(params Params) (models.MacAddrTableEntries, error) {

	session, err := snmpConnect(params.Target, params.Community)
	if err != nil {
		return nil, fmt.Errorf("connect err: %v", err)
	}
	defer session.Conn.Close()

	portNumMap, err := getPortNumberMap(session)
	if err != nil {
		return nil, fmt.Errorf("getPortNumberMap err: %v", err)
	}

	pduList, err := session.BulkWalkAll(macAddressTableOid)
	if err != nil {
		return nil, fmt.Errorf("bulk walk err: %v", err)
	}
	result := make(models.MacAddrTableEntries, len(pduList))
	for i, pdu := range pduList {
		result[i] = pduToMacAddrTableEntry(&pdu, portNumMap)
	}
	return result, nil
}

const vlanListOid = ".1.3.6.1.4.1.9.9.46.1.3.1.1.2"

func GetVlanList(params Params) (models.Vlans, error) {
	session, err := snmpConnect(params.Target, params.Community)
	if err != nil {
		return nil, err
	}
	defer session.Conn.Close()
	pduList, err := session.BulkWalkAll(vlanListOid)
	if err != nil {
		return nil, fmt.Errorf("bulk walk err: %v", err)
	}
	result := make(models.Vlans, len(pduList))
	for i, pdu := range pduList {
		result[i].Number = getLastNumberFromPduName(pdu.Name)
	}
	return result, nil
}

const portListOid = ".1.3.6.1.2.1.31.1.1.1.1"

func GetPortList(params Params) (models.Ports, error) {
	session, err := snmpConnect(params.Target, params.Community)
	if err != nil {
		return nil, err
	}
	defer session.Conn.Close()
	pduList, err := session.BulkWalkAll(portListOid)
	if err != nil {
		return nil, fmt.Errorf("bulk walk err: %v", err)
	}
	result := make(models.Ports, len(pduList))
	for i, pdu := range pduList {
		result[i].Number = getLastNumberFromPduName(pdu.Name)
		result[i].Name = string(pdu.Value.([]byte))
	}
	return result, nil
}
