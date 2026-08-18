package main

import (
	argsparser "PortScaner/ArgsParser"
	//udp "PortScaner/UDP"
	"fmt"
	"os"
)

// func TestSendUDP(sock udp.Socket) {
// 	// Цель
// 	targetIP := [4]byte{10, 33, 21, 1}
// 	targetPort := uint16(5000)

// 	// UDP probe
// 	data := []byte("UDP probe")

// 	fmt.Printf(
// 		"Sending UDP packet to %d.%d.%d.%d:%d\n",
// 		targetIP[0],
// 		targetIP[1],
// 		targetIP[2],
// 		targetIP[3],
// 		targetPort,
// 	)

// 	err := udp.SendUDP(
// 		sock,
// 		targetIP,
// 		targetPort,
// 		data,
// 	)

// 	if err != nil {
// 		fmt.Println("SendUDP error:", err)
// 		return
// 	}

// 	fmt.Println("UDP packet sent")
// }

// go run . -m tcp 10.33.21.1
// go run . -m udp 10.33.21.1
// go run . -m both 10.33.21.1
func main() {

	cfg, err := argsparser.ValidateArgs(os.Args[1:])

	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	Scan(*cfg)
	//UdpScan1()

}

// func UdpScan() {
// 	sock, err := udp.CreateUDPSocket()
// 	if err != nil {
// 		fmt.Println("Create UDP socket error:", err)
// 		return
// 	}
// 	defer udp.CloseSocket(sock)

// 	fmt.Println("UDP socket created:", sock)
// 	addr, err := udp.GetSocketName(sock)
// 	if err != nil {

// 		fmt.Println(
// 			"getsockname error:",
// 			err,
// 		)

// 		return
// 	}
// 	fmt.Println(
// 		"Local port:",
// 		addr.Port,
// 	)

// 	icmpSock, err := CreateICMPSocket()

// 	if err != nil {
// 		fmt.Println("ICMP socket error:", err)
// 		return
// 	}

// 	defer CloseSocket(icmpSock)

// 	err = BindICMP(icmpSock)

// 	if err != nil {
// 		fmt.Println("ICMP bind error:", err)
// 		return
// 	}

// 	fmt.Println("ICMP socket bound")

// 	packetCount := 0
// 	recvBuf := make([]byte, 4096)
// 	msTimeout := (5000) * time.Millisecond
// 	startTimer := time.Now()

// 	fmt.Println("FDSet size:", unsafe.Sizeof(FDSet{}))
// 	fmt.Println("FDSet Array offset:", unsafe.Offsetof(FDSet{}.Array))
// 	fmt.Println("Socket size:", unsafe.Sizeof(Socket(0)))
// 	fmt.Println("Timeval size:", unsafe.Sizeof(Timeval{}))

// 	for { //
// 		if time.Since(startTimer) >= msTimeout {
// 			fmt.Println("MyTimeout!")
// 			break
// 		}
// 		if packetCount < 1 { //len(cfg.Targets)
// 			TestSendUDP(sock)

// 			// time.Sleep(1 * time.Second)

// 			// fmt.Println("Trying recvfrom directly...")

// 			// err = ReceiveICMP(
// 			// 	icmpSock,
// 			// 	recvBuf,
// 			// )

// 			// fmt.Println("ReceiveICMP:", err)

// 			packetCount++
// 		}
// 		var fds FDSet

// 		FDZero(&fds)
// 		FDSetSocket(icmpSock, &fds)

// 		// fmt.Println("ICMP socket:", icmpSock)
// 		// fmt.Println("FD count:", fds.Count)
// 		// fmt.Println("FD[0]:", fds.Array[0])

// 		tv := Timeval{
// 			Sec:  0,
// 			Usec: 0,
// 		}

// 		ret, err := Select(&fds, &tv)

// 		// fmt.Println(
// 		// 	"SELECT RESULT:",
// 		// 	ret,
// 		// 	"ERROR:",
// 		// 	err,
// 		// )

// 		if ret > 0 {
// 			fmt.Println("ret > 0")

// 			if FDIsSet(icmpSock, &fds) {

// 				for {
// 					// Приём ответа
// 					isRecieve := ReceiveICMP(
// 						icmpSock,
// 						recvBuf,
// 					)
// 					if isRecieve == syscall.Errno(10035) {
// 						break
// 					} else if isRecieve != nil {
// 						break
// 					}
// 				}

// 			}

// 		} else if ret != 0 {

// 			fmt.Println("select error:", err)
// 		} else {
// 			//fmt.Println("select = 0:")
// 		}

// 	}

// 	udp.CloseSocket(sock)
// 	fmt.Println("UDP socket closed")
// }
