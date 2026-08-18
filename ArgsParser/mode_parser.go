package argsparser

import (
	"fmt"
	"strings"
)

type ScanMode int

const (
	ScanTCP ScanMode = iota
	ScanUDP
	ScanBoth
)

func (m ScanMode) String() string {

	switch m {

	case ScanTCP:
		return "tcp"

	case ScanUDP:
		return "udp"

	case ScanBoth:
		return "both"

	default:
		return "unknown"
	}
}

func ParseScanMode(value string) (ScanMode, error) {

	switch strings.ToLower(strings.TrimSpace(value)) {

	case "tcp":
		return ScanTCP, nil

	case "udp":
		return ScanUDP, nil

	case "both":
		return ScanBoth, nil

	default:
		return 0, fmt.Errorf(
			"invalid scan mode: %s (expected tcp, udp or tcp+udp)",
			value,
		)
	}
}
