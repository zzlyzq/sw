package sw

import (
	"fmt"
	"github.com/alouca/gosnmp"
	"log"
	"strconv"
	"strings"
	"time"
)

type IfStats struct {
	IfName           string
	IfIndex          int
	IfHCInOctets     int64
	IfHCOutOctets    int64
	IfHCInUcastPkts  int64
	IfHCOutUcastPkts int64
	TS               int64
}

func (this *IfStats) String() string {
	return fmt.Sprintf("<IfName:%s, IfIndex:%d, IfHCInOctets:%d, IfHCOutOctets:%d>", this.IfName, this.IfIndex, this.IfHCInOctets, this.IfHCOutOctets)
}

func ListIfStats(ip, community string, timeout int, onlyPrefix []string) ([]IfStats, error) {
	var ifStatsList []IfStats

	chIfInList := make(chan []gosnmp.SnmpPDU)
	chIfOutList := make(chan []gosnmp.SnmpPDU)

	chIfInPktList := make(chan []gosnmp.SnmpPDU)
	chIfOutPktList := make(chan []gosnmp.SnmpPDU)

	chIfNameList := make(chan []gosnmp.SnmpPDU)

	go ListIfHCInOctets(ip, community, timeout, chIfInList)
	go ListIfHCOutOctets(ip, community, timeout, chIfOutList)

	go ListIfHCInUcastPkts(ip, community, timeout, chIfInPktList)
	go ListIfHCOutUcastPkts(ip, community, timeout, chIfOutPktList)

	go ListIfName(ip, community, timeout, chIfNameList)

	ifInList := <-chIfInList
	ifOutList := <-chIfOutList

	ifInPktList := <-chIfInPktList
	ifOutPktList := <-chIfOutPktList

	ifNameList := <-chIfNameList

	if len(ifNameList) > 0 && len(ifInList) > 0 && len(ifOutList) > 0 && len(ifInPktList) > 0 && len(ifOutPktList) > 0 {
		for _, ifNamePDU := range ifNameList {

			ifName := ifNamePDU.Value.(string)

			var found bool
			if len(onlyPrefix) > 0 {
				found = false
				for _, prefix := range onlyPrefix {
					if strings.Contains(ifName, prefix) {
						found = true
						break
					}
				}
			} else {
				found = true
			}

			if ifName == "Nu0" || strings.Contains(ifName, "Stack") {
				found = false
			}

			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered in ListIfStats", r)
				}
			}()

			if found {
				var ifStats IfStats

				ifIndexStr := strings.Replace(ifNamePDU.Name, ".1.3.6.1.2.1.31.1.1.1.1.", "", 1)

				ifStats.IfIndex, _ = strconv.Atoi(ifIndexStr)

				for ti, ifHCInOctetsPDU := range ifInList {
					if strings.Replace(ifHCInOctetsPDU.Name, ".1.3.6.1.2.1.31.1.1.1.6.", "", 1) == ifIndexStr {

						ifStats.IfHCInOctets = ifInList[ti].Value.(int64)
						ifStats.IfHCOutOctets = ifOutList[ti].Value.(int64)

						ifStats.IfHCInUcastPkts = ifInPktList[ti].Value.(int64)
						ifStats.IfHCOutUcastPkts = ifOutPktList[ti].Value.(int64)

						ifStats.TS = time.Now().Unix()
						ifStats.IfName = ifName
					}
				}

				ifStatsList = append(ifStatsList, ifStats)

			}
		}
	}

	return ifStatsList, nil
}

func ListIfHCInOctets(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU) {
	oid := "1.3.6.1.2.1.31.1.1.1.6"
	method := "walk"

	snmpPDUs, _ := RunSnmp(ip, community, oid, method, timeout)

	ch <- snmpPDUs
	return
}

func ListIfHCOutOctets(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU) {
	oid := "1.3.6.1.2.1.31.1.1.1.10"
	method := "walk"

	snmpPDUs, _ := RunSnmp(ip, community, oid, method, timeout)

	ch <- snmpPDUs
	return
}

func ListIfHCInUcastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU) {
	oid := "1.3.6.1.2.1.31.1.1.1.7"
	method := "walk"

	snmpPDUs, _ := RunSnmp(ip, community, oid, method, timeout)

	ch <- snmpPDUs
	return
}

func ListIfHCOutUcastPkts(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU) {
	oid := "1.3.6.1.2.1.31.1.1.1.11"
	method := "walk"

	snmpPDUs, _ := RunSnmp(ip, community, oid, method, timeout)

	ch <- snmpPDUs
	return
}

func ListIfName(ip, community string, timeout int, ch chan []gosnmp.SnmpPDU) {
	oid := "1.3.6.1.2.1.31.1.1.1.1"
	method := "walk"

	snmpPDUs, _ := RunSnmp(ip, community, oid, method, timeout)

	ch <- snmpPDUs
	return
}