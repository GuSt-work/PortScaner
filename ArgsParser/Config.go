package argsparser

import (
	"net"
	"time"
)

type Config struct {
	Targets      []net.IP
	Ports        []uint16
	PacketsCount uint
	Timeout      time.Duration
	Mode         ScanMode
}
