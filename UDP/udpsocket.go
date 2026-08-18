package udp

import (
	//infrastruct "PortScaner/Infrastruction"
	"fmt"
	"syscall"
	"time"
	"unsafe"
)

var (
	ws2_32 = syscall.NewLazyDLL("ws2_32.dll")

	procSocket      = ws2_32.NewProc("socket")
	procSetsockopt  = ws2_32.NewProc("setsockopt")
	procIoctlsocket = ws2_32.NewProc("ioctlsocket")
	procClosesocket = ws2_32.NewProc("closesocket")
	procGetsockopt  = ws2_32.NewProc("getsockopt")
	procBind        = ws2_32.NewProc("bind")
	procGetsockname = ws2_32.NewProc("getsockname")
	procSendto      = ws2_32.NewProc("sendto")

	WSAIoctlProc = ws2_32.NewProc("WSAIoctl")
)

type Socket syscall.Handle

const (
	IPPROTO_IP   = 0
	IPPROTO_ICMP = 1
	AF_INET      = 2
	SOCK_RAW     = 3
	SOCK_DGRAM   = 2

	IP_TTL      = 4
	IPPROTO_TCP = 6
	IPPROTO_UDP = 17

	FIONBIO = 0x8004667E

	INVALID_SOCKET    = ^uintptr(0)
	IP_HDRINCL        = 2
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

// func CreateUDPSocket() (Socket, error) {
// 	sock, _, err := procSocket.Call(
// 		uintptr(AF_INET),
// 		uintptr(SOCK_DGRAM),
// 		uintptr(IPPROTO_UDP),
// 	)

// 	if sock == INVALID_SOCKET {

// 		fmt.Println(
// 			"[ERROR] UDP socket() failed:",
// 			err,
// 		)

// 		return 0, err
// 	}

// 	err = BindUDP(
// 		Socket(sock),
// 		[4]byte{0, 0, 0, 0},
// 		0,
// 	)

// 	if err != nil {
// 		CloseSocket(Socket(sock))
// 		return 0, err
// 	}
// 	err = SetNonBlocking(Socket(sock))

// 	if err != nil {

// 		CloseSocket(Socket(sock))

// 		return 0, err
// 	}
// 	return Socket(sock), nil
// }

// func BindUDP(sock Socket, ip [4]byte, port uint16) error {

// 	addr := infrastruct.SockaddrInet4{
// 		Family: AF_INET,
// 		Port:   Htons(port),
// 		Addr:   ip,
// 	}

// 	ret, _, err := procBind.Call(
// 		uintptr(sock),
// 		uintptr(unsafe.Pointer(&addr)),
// 		uintptr(unsafe.Sizeof(addr)),
// 	)

// 	if ret != 0 {

// 		return fmt.Errorf(
// 			"bind failed: %w",
// 			err,
// 		)
// 	}

// 	return nil
// }

func Htons(port uint16) uint16 {
	return (port<<8)&0xff00 | port>>8
}
func ntohs(port uint16) uint16 {
	return (port<<8)&0xff00 | port>>8
}

// func SetNonBlocking(sock Socket) error {

// 	var mode uint32 = 1

// 	ret, _, err := procIoctlsocket.Call(
// 		uintptr(sock),
// 		uintptr(FIONBIO),
// 		uintptr(unsafe.Pointer(&mode)),
// 	)

// 	if ret != 0 {

// 		return err
// 	}

// 	return nil
// }

// func CreateICMPSocket() (Socket, error) {
// 	fd, err := syscall.Socket(
// 		syscall.AF_INET,
// 		syscall.SOCK_RAW,
// 		IPPROTO_ICMP,
// 	)

// 	if err != nil {
// 		return 0, err
// 	}

// 	return Socket(fd), nil
// }

func CloseSocket(sock Socket) {

	if sock != Socket(INVALID_SOCKET) {

		procClosesocket.Call(
			uintptr(sock),
		)
	}
}

// func GetSocketName(sock Socket) (infrastruct.SockaddrInet4, error) {

// 	var addr infrastruct.SockaddrInet4

// 	size := uint32(unsafe.Sizeof(addr))

// 	ret, _, err := procGetsockname.Call(
// 		uintptr(sock),
// 		uintptr(unsafe.Pointer(&addr)),
// 		uintptr(unsafe.Pointer(&size)),
// 	)

// 	if ret != 0 {

// 		return addr, fmt.Errorf(
// 			"getsockname failed: %v",
// 			err,
// 		)
// 	}

// 	addr.Port = ntohs(addr.Port)

// 	return addr, nil
// }

// func SendUDP(
// 	sock Socket,
// 	ip [4]byte,
// 	port uint16,
// 	data []byte,
// ) error {

// 	addr := infrastruct.SockaddrInet4{
// 		Family: AF_INET,
// 		Port:   Htons(port),
// 		Addr:   ip,
// 	}

// 	var dataPtr uintptr

// 	if len(data) > 0 {
// 		dataPtr = uintptr(unsafe.Pointer(&data[0]))
// 	}

// 	ret, _, err := procSendto.Call(
// 		uintptr(sock),
// 		dataPtr,
// 		uintptr(len(data)),
// 		0,
// 		uintptr(unsafe.Pointer(&addr)),
// 		uintptr(unsafe.Sizeof(addr)),
// 	)

// 	if ret == ^uintptr(0) {
// 		return fmt.Errorf("sendto failed: %v", err)
// 	}

// 	if int(ret) != len(data) {
// 		return fmt.Errorf(
// 			"sendto sent %d bytes, expected %d",
// 			ret,
// 			len(data),
// 		)
// 	}

// 	return nil
// }
