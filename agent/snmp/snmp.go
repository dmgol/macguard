package snmp

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/soniah/gosnmp"
)

type GetMacAddressTableRequest struct {
	Target    string
	Community string
}

type MacAddressTableEntry struct {
	Address string
	Port    int
}

type MacAddressTable []MacAddressTableEntry

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

func extractMACFromName(pduName string) string {
	split := strings.Split(pduName, ".")
	macAddrStart := len(split) - macAddressOctetCount
	hexStrings := make([]string, macAddressOctetCount)
	for i := 0; i < macAddressOctetCount; i++ {
		v, _ := strconv.Atoi(split[i+macAddrStart])
		hexStrings[i] = strconv.FormatInt(int64(v), 16)
	}
	return strings.Join(hexStrings, ":")
}

func pduToMacAddressTableEntry(pdu gosnmp.SnmpPDU) MacAddressTableEntry {
	result := MacAddressTableEntry{}
	result.Address = extractMACFromName(pdu.Name)
	if pdu.Type == gosnmp.Integer {
		result.Port = pdu.Value.(int)
	}
	return result
}

func GetMacAddressTable(target net.IP, community string) (MacAddressTable, error) {

	session, err := snmpConnect(target.String(), community)
	if err != nil {
		return nil, fmt.Errorf("connect err: %v", err)
	}
	defer session.Conn.Close()

	pduList, err := session.BulkWalkAll(macAddressTableOid)
	if err != nil {
		return nil, fmt.Errorf("bulk walk err: %v", err)
	}
	result := make(MacAddressTable, len(pduList))
	for i, pdu := range pduList {
		result[i] = pduToMacAddressTableEntry(pdu)
	}
	return result, nil
}

const vlanListOid = ".1.3.6.1.4.1.9.9.46.1.3.1.1.2"

func getLastNumberFromName(name string) int {
	if last := strings.LastIndex(name, "."); last != -1 {
		result, _ := strconv.Atoi(name[last+1:])
		return result
	}
	return -1
}

func pduToVlan(pdu gosnmp.SnmpPDU) int {
	return getLastNumberFromName(pdu.Name)
}

func getCiscoVlanList(session *gosnmp.GoSNMP) ([]int, error) {
	pduList, err := session.BulkWalkAll(vlanListOid)
	if err != nil {
		return nil, fmt.Errorf("bulk walk err: %v", err)
	}
	result := make([]int, len(pduList))
	for i, pdu := range pduList {
		result[i] = pduToVlan(pdu)
	}
	return result, nil
}

func TestGetCiscoVlanList(target net.IP, community string) []int {
	session, err := snmpConnect(target.String(), community)
	if err != nil {
		return nil
	}
	defer session.Conn.Close()
	vlans, _ := getCiscoVlanList(session)
	return vlans
}
