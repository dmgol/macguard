package snmp

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/dmgol/macguard/models"
	"github.com/soniah/gosnmp"
)

type SnmpParams struct {
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

const macAddressTableOid = ".1.3.6.1.2.1.17.4.3.1.2"
const macAddressOctetCount = 6

func getLastNumberFromPduName(name string) int {
	if last := strings.LastIndex(name, "."); last != -1 {
		result, _ := strconv.Atoi(name[last+1:])
		return result
	}
	return -1
}

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

func pduToMacAddrTableEntry(pdu *gosnmp.SnmpPDU) models.MacAddrTableEntry {
	result := models.MacAddrTableEntry{}
	result.MacAddr = extractMacAddrFromPdu(pdu)
	if pdu.Type == gosnmp.Integer {
		result.PortNumber = pdu.Value.(int)
	}
	return result
}

func GetMacAddrTable(params SnmpParams) (models.MacAddrTableEntries, error) {

	session, err := snmpConnect(params.Target, params.Community)
	if err != nil {
		return nil, fmt.Errorf("connect err: %v", err)
	}
	defer session.Conn.Close()

	pduList, err := session.BulkWalkAll(macAddressTableOid)
	if err != nil {
		return nil, fmt.Errorf("bulk walk err: %v", err)
	}
	result := make(models.MacAddrTableEntries, len(pduList))
	for i, pdu := range pduList {
		result[i] = pduToMacAddrTableEntry(&pdu)
	}
	return result, nil
}

const vlanListOid = ".1.3.6.1.4.1.9.9.46.1.3.1.1.2"

func GetVlanList(params SnmpParams) (models.Vlans, error) {
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
