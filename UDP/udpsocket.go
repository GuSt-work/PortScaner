package udp

import (
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

var (
	ws2_32          = syscall.NewLazyDLL("ws2_32.dll")
	procSocket      = ws2_32.NewProc("socket")
	procClosesocket = ws2_32.NewProc("closesocket")
	WSAIoctlProc    = ws2_32.NewProc("WSAIoctl")
)

type Socket syscall.Handle

const (
	IPPROTO_IP = 0
	AF_INET    = 2
	SOCK_DGRAM = 2

	IPPROTO_UDP = 17

	INVALID_SOCKET    = ^uintptr(0)
	SIO_UDP_CONNRESET = 0x9800000C
)

func CreateUDPSocket(timeout time.Duration) (Socket, error) {

	// Создаём UDP socket
	sock, _, err := socketProc.Call(
		uintptr(AF_INET),
		uintptr(SOCK_DGRAM),
		uintptr(IPPROTO_UDP),
	)

	if sock == uintptr(^uint(0)) {
		//fmt.Println("socket error: %d", err)
		return 0, err
	}

	// Разрешаем получать ICMP Port Unreachable
	// как WSAECONNRESET через UDP socket.
	enable := uint32(1)
	var bytesReturned uint32

	ret, _, err := WSAIoctlProc.Call(
		sock,
		uintptr(SIO_UDP_CONNRESET),
		uintptr(unsafe.Pointer(&enable)),
		uintptr(unsafe.Sizeof(enable)),
		0,
		0,
		uintptr(unsafe.Pointer(&bytesReturned)),
		0,
		0,
	)

	if ret != 0 {
		fmt.Printf("WSAIoctl(SIO_UDP_CONNRESET) failed: %v\n", err)
		closesocketProc.Call(sock)
		return 0, err
	}

	// SO_RCVTIMEO
	ms := uint32(timeout / time.Millisecond)

	ret, _, err = setsockoptProc.Call(
		sock,
		uintptr(SOL_SOCKET),
		uintptr(SO_RCVTIMEO),
		uintptr(unsafe.Pointer(&ms)),
		uintptr(unsafe.Sizeof(ms)),
	)

	if ret != 0 {
		fmt.Printf("setsockopt error: %v", err)
		return 0, err
	}

	return Socket(sock), nil
}

func Htons(port uint16) uint16 {
	return (port<<8)&0xff00 | port>>8
}
func ntohs(port uint16) uint16 {
	return (port<<8)&0xff00 | port>>8
}

func CloseSocket(sock Socket) {

	if sock != Socket(INVALID_SOCKET) {

		procClosesocket.Call(
			uintptr(sock),
		)
	}
}
