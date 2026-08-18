package udp

import (
	argsparser "PortScaner/ArgsParser"
	infrastruct "PortScaner/Infrastruction"

	//"errors"
	"fmt"
	"net"
	"syscall"
	"unsafe"
)

const (
	SOL_SOCKET    = 0xffff
	SO_RCVTIMEO   = 0x1006
	WSAETIMEDOUT  = 10060
	WSAECONNRESET = 10054
)

var (
	socketProc      = ws2_32.NewProc("socket")
	setsockoptProc  = ws2_32.NewProc("setsockopt")
	sendtoProc      = ws2_32.NewProc("sendto")
	recvfromProc    = ws2_32.NewProc("recvfrom")
	closesocketProc = ws2_32.NewProc("closesocket")
)

type sockaddrIn struct {
	Family uint16
	Port   uint16
	Addr   uint32
	Zero   [8]byte
}

func htons(v uint16) uint16 {
	return (v << 8) | (v >> 8)
}

func inetAddr(ip net.IP) uint32 {
	ip4 := ip.To4()
	if ip4 == nil {
		return 0
	}

	return uint32(ip4[0]) |
		uint32(ip4[1])<<8 |
		uint32(ip4[2])<<16 |
		uint32(ip4[3])<<24
}

func InitResult(
	ip string,
	port uint16,
) infrastruct.ScanResult {
	return infrastruct.ScanResult{
		IP:                 ip,
		Port:               port,
		PortType:           "udp",
		PacketsPortOpen:    0,
		PacketsPortTimeout: 0,
		PacketsPortClosed:  0,
	}
}
func SetResult(
	currentResult *infrastruct.ScanResult,
	err error,
) {
	errCode := syscall.Errno(err.(syscall.Errno))
	fmt.Println("errCode: ", errCode)
	switch errCode {
	case WSAECONNRESET:
		currentResult.PacketsPortClosed++

	case WSAETIMEDOUT: // "Open|Filtered"
		currentResult.PacketsPortTimeout++

	case 0: // "Open
		currentResult.PacketsPortOpen++

	default:
		currentResult.PacketsPortClosed++
	}
}

func ResultForming(
	ip string,
	port uint16,
	state infrastruct.PortState,
	err error) infrastruct.ScanResult {
	return infrastruct.ScanResult{
		IP:       ip,
		Port:     port,
		State:    state,
		PortType: "udp",
		Error:    err,
	}
}

func GetPortState(currentResult *infrastruct.ScanResult) {

	if currentResult.PacketsPortOpen > 0 {
		currentResult.State = infrastruct.PortOpen
	} else if currentResult.PacketsPortTimeout > 0 {
		currentResult.State = infrastruct.PortFiltered
	} else {
		currentResult.State = infrastruct.PortClosed
	}
}

func ScanUDPPort(sock Socket, ip net.IP, port uint16, config argsparser.Config) infrastruct.ScanResult {

	// Адрес назначения
	dest := sockaddrIn{
		Family: AF_INET,
		Port:   htons(port),
		Addr:   inetAddr(ip),
	}

	// Наш UDP probe
	payload := []byte("X")
	result := InitResult(ip.String(), port)

	for count := uint(0); count < config.PacketsCount; count++ {

		fmt.Println("count: ", count)
		fmt.Println("port: ", port)

		// sendto()
		ret, _, err := sendtoProc.Call(
			uintptr(sock),
			uintptr(unsafe.Pointer(&payload[0])),
			uintptr(len(payload)),
			0,
			uintptr(unsafe.Pointer(&dest)),
			uintptr(unsafe.Sizeof(dest)),
		)

		if ret == ^uintptr(0) {
			fmt.Printf("sendto error: %v", err)
			return ResultForming(ip.String(), port, infrastruct.PortClosed, err)
		}

		buffer := make([]byte, 512)

		var from sockaddrIn
		fromLen := int32(unsafe.Sizeof(from))

		ret, _, err = recvfromProc.Call(
			uintptr(sock),
			uintptr(unsafe.Pointer(&buffer[0])),
			uintptr(len(buffer)),
			0,
			uintptr(unsafe.Pointer(&from)),
			uintptr(unsafe.Pointer(&fromLen)),
		)
		SetResult(&result, err)

	}
	GetPortState(&result)
	return result
}
