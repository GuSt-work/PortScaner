package main

import (
	argsparser "PortScaner/ArgsParser"
	infrastruct "PortScaner/Infrastruction"
	udp "PortScaner/UDP"
	"net"
)

func Scan(config argsparser.Config) {

	var (
		udpSock udp.Socket
		err     error
	)

	if config.Mode == argsparser.ScanUDP || config.Mode == argsparser.ScanBoth {
		udpSock, err = udp.CreateUDPSocket(config.Timeout)
		if err != nil {
			return
		}
		defer udp.CloseSocket(udpSock)
	}

	for _, ip := range config.Targets {

		results := ScanHost(
			udpSock,
			ip,
			config,
		)

		infrastruct.PrintResult(results, ip, config.PacketsCount)
	}
}

func ScanHost(
	udpSock udp.Socket,
	ip net.IP,
	config argsparser.Config,
) []infrastruct.ScanResult {

	results := make([]infrastruct.ScanResult, 0)

	for _, port := range config.Ports {

		if config.Mode == argsparser.ScanTCP || config.Mode == argsparser.ScanBoth {
			result := ScanTcp(
				ip.String(),
				port,
				config,
			)
			results = append(
				results,
				result,
			)
		}
		if config.Mode == argsparser.ScanUDP || config.Mode == argsparser.ScanBoth {
			result := udp.ScanUDPPort(
				udpSock,
				ip,
				port,
				config,
			)
			results = append(
				results,
				result,
			)
		}
	}

	return results
}
