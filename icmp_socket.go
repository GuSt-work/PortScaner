package main

import (
	infrastruct "PortScaner/Infrastruction"
	"fmt"
	"syscall"
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
	procRecvfrom    = ws2_32.NewProc("recvfrom")
	procSelect      = ws2_32.NewProc("select")
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

	INVALID_SOCKET = ^uintptr(0)
	IP_HDRINCL     = 2
)

type FDSet struct {
	Count uint32
	Array [64]Socket //64 — это значение FD_SETSIZE в WinSock по умолчанию
}

func CreateICMPSocket() (Socket, error) {

	sock, _, err := procSocket.Call(
		uintptr(AF_INET),
		uintptr(SOCK_RAW),
		uintptr(IPPROTO_ICMP),
	)

	if sock == INVALID_SOCKET {

		return 0, fmt.Errorf(
			"ICMP socket failed: %v",
			err,
		)
	}

	err = SetNonBlocking(Socket(sock))

	if err != nil {

		CloseSocket(Socket(sock))

		return 0, err
	}

	return Socket(sock), nil
}

func SetNonBlocking(sock Socket) error {

	var mode uint32 = 1

	ret, _, err := procIoctlsocket.Call(
		uintptr(sock),
		uintptr(FIONBIO),
		uintptr(unsafe.Pointer(&mode)),
	)

	if ret != 0 {

		return err
	}

	return nil
}

func ReceiveICMP(sock Socket, buffer []byte) error {

	//var addr infrastruct.SockaddrInet4
	//addrLen := uint32(unsafe.Sizeof(addr))

	ret, _, err := procRecvfrom.Call(
		uintptr(sock),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(len(buffer)),
		0,
		0,
		0,
	)

	if int32(ret) == -1 {

		//WSAEWOULDBLOCK
		if err == syscall.Errno(10035) {
			fmt.Println("ICMP WSAEWOULDBLOCK")
			return err
		}

		fmt.Println("recvfrom error:", err)

		return err
	}

	bytesReceived := int(ret)

	if bytesReceived == 0 {
		fmt.Println("ICMP bytesReceived == 0")
		return err
	}

	fmt.Println(
		"ICMP recvfrom:", ret, "bytes",
	)

	return nil
}

func CloseSocket(sock Socket) {

	if sock != Socket(INVALID_SOCKET) {

		procClosesocket.Call(
			uintptr(sock),
		)

	}
}

func FDZero(set *FDSet) {
	set.Count = 0
}

func FDSetSocket(sock Socket, set *FDSet) {

	if set.Count >= uint32(len(set.Array)) {
		return
	}

	set.Array[set.Count] = sock
	set.Count++
}

func FDIsSet(sock Socket, set *FDSet) bool {

	for i := uint32(0); i < set.Count; i++ {

		if set.Array[i] == sock {
			return true
		}
	}

	return false
}

type Timeval struct {
	Sec  int32
	Usec int32
}

func Select(readSet *FDSet, timeout *Timeval) (int, error) {

	ret, _, err := procSelect.Call(
		0,
		uintptr(unsafe.Pointer(readSet)),
		0,
		0,
		uintptr(unsafe.Pointer(timeout)),
	)

	if int32(ret) == -1 {

		if err != syscall.Errno(0) {
			return 0, err
		}

		return 0, syscall.EINVAL
	}

	return int(ret), nil
}

// func Select(readSet *FDSet, timeout *Timeval) (int, error) {

// 	// fmt.Println("Calling select...")
// 	// fmt.Println("  readSet:", readSet)
// 	// fmt.Println("  count:", readSet.Count)
// 	// fmt.Println("  fd[0]:", readSet.Array[0])
// 	// fmt.Println("  timeout:", timeout.Sec, timeout.Usec)

// 	ret, _, err := procSelect.Call(
// 		0,
// 		uintptr(unsafe.Pointer(readSet)),
// 		0,
// 		0,
// 		uintptr(unsafe.Pointer(timeout)),
// 	)

// 	// fmt.Println("select raw return:", ret)
// 	// fmt.Println("select raw error:", err)

// 	if ret == INVALID_SOCKET {
// 		return -1, err
// 	}

// 	return int(ret), nil
// }

func BindICMP(sock Socket) error {

	addr := infrastruct.SockaddrInet4{
		Family: AF_INET,
		Port:   0,
		Addr:   [4]byte{0, 0, 0, 0},
	}

	ret, _, err := procBind.Call(
		uintptr(sock),
		uintptr(unsafe.Pointer(&addr)),
		uintptr(unsafe.Sizeof(addr)),
	)

	if int32(ret) == -1 {
		return fmt.Errorf("ICMP bind failed: %v", err)
	}

	return nil
}
