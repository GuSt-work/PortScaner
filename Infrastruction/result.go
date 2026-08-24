package infrastruction

import (
	"fmt"
	"net"
)

type ScanResult struct {
	IP                 string
	Port               uint16
	State              PortState
	PortType           string
	Error              error
	PacketsPortOpen    uint
	PacketsPortClosed  uint
	PacketsPortTimeout uint
}

type PortState int

const (
	PortOpen PortState = iota
	PortClosed
	PortFiltered
)

func (s PortState) String() string {

	switch s {

	case PortOpen:
		return "open"

	case PortClosed:
		return "closed"

	case PortFiltered:
		return "open|filtered"

	default:
		return "unknown"
	}
}

func PrintHeader(ip string) {

	fmt.Printf(
		"\nScan report for %s\n\n",
		ip,
	)

	fmt.Printf(
		"%-12s%s\n",
		"PORT",
		"STATE",
	)
}

func PrintBody(result ScanResult, pcktCount uint) {

	fmt.Printf(
		"%-12s%-12s Open: %d/%d  Open|filtered: %d/%d  Closed: %d/%d\n",
		fmt.Sprintf("%d/%s", result.Port, result.PortType),
		result.State,
		result.PacketsPortOpen, pcktCount,
		result.PacketsPortTimeout, pcktCount,
		result.PacketsPortClosed, pcktCount,
	)
}

func PrintResult(results []ScanResult, ip net.IP, pcktCount uint) {

	PrintHeader(ip.String())

	for _, result := range results {
		//if result.State == PortOpen {
		PrintBody(result, pcktCount)
		//}
	}
}
