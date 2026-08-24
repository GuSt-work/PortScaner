package main

import (
	argsparser "PortScaner/ArgsParser"
	infrastruct "PortScaner/Infrastruction"
	"errors"
	"fmt"
	"net"
	"os"
	"syscall"
)

const (
	WSAECONNREFUSED = 10061
	WSAETIMEDOUT    = 10060
)

func InitResult(
	ip string,
	port uint16,
) infrastruct.ScanResult {
	return infrastruct.ScanResult{
		IP:                 ip,
		Port:               port,
		PortType:           "tcp",
		PacketsPortOpen:    0,
		PacketsPortTimeout: 0,
		PacketsPortClosed:  0,
	}
}

func GetTCPPortState(err error) infrastruct.PortState {

	if err == nil {
		return infrastruct.PortOpen
	}

	var syscallErr *os.SyscallError

	if errors.As(err, &syscallErr) {

		errno, ok := syscallErr.Err.(syscall.Errno)

		if ok {

			switch errno {

			case WSAECONNREFUSED:
				return infrastruct.PortClosed

			case WSAETIMEDOUT:
				return infrastruct.PortFiltered
			}
		}
	}

	return infrastruct.PortFiltered
}

func SetResult(
	currentResult *infrastruct.ScanResult,
	err error,
) {
	result := GetTCPPortState(err)

	switch result {
	case infrastruct.PortOpen:
		currentResult.PacketsPortOpen++
	case infrastruct.PortFiltered:
		currentResult.PacketsPortTimeout++
	case infrastruct.PortClosed:
		currentResult.PacketsPortClosed++
	}
}

func ScanTcp(
	ip string,
	port uint16,
	config argsparser.Config,
) infrastruct.ScanResult {

	result := InitResult(ip, port)

	for count := uint(0); count < config.PacketsCount; count++ {

		address := fmt.Sprintf(
			"%s:%d",
			ip,
			port,
		)
		conn, err := net.DialTimeout(
			"tcp",
			address,
			config.Timeout,
		)

		// Соединение установлено
		if err == nil {
			conn.Close()
		}
		SetResult(&result, err)

	}

	GetPortState(&result)
	return result
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
