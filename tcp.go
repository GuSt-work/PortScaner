package main

import (
	//argsparser "PortScaner/ArgsParser"
	argsparser "PortScaner/ArgsParser"
	infrastruct "PortScaner/Infrastruction"
	"errors"
	"fmt"
	"net"
	"syscall"
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
func SetResult(
	currentResult *infrastruct.ScanResult,
	err error,
) {
	if err == nil {
		currentResult.PacketsPortOpen++
		return
	}

	if errors.Is(err, syscall.ECONNREFUSED) {
		currentResult.PacketsPortClosed++
		return
	}

	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		currentResult.PacketsPortTimeout++
		return
	}

	// Все остальные ошибки пока считаем закрытым портом.
	currentResult.PacketsPortClosed++
}

func ScanTcp(
	ip string,
	port uint16,
	config argsparser.Config,
) infrastruct.ScanResult {

	result := InitResult(ip, port)

	for count := uint(0); count < config.PacketsCount; count++ {

		fmt.Println("Packets numb: ", count)

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
