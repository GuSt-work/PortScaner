package argsparser

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

func ParseTargets(args []string) ([]net.IP, error) {
	var result []net.IP

	for _, arg := range args {

		// Проверяем обычный IP
		if ip := net.ParseIP(arg); ip != nil {
			result = append(result, ip)
			continue
		}

		// Проверяем диапазон
		ips, err := ParseIPRange(arg)
		if err != nil {
			return nil, err
		}

		result = append(result, ips...)
	}

	return result, nil
}

func ParseIPRange(s string) ([]net.IP, error) {

	parts := strings.Split(s, "-")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid target: %s", s)
	}

	startIP := net.ParseIP(parts[0])
	if startIP == nil {
		return nil, fmt.Errorf("invalid start IP: %s", parts[0])
	}

	// Берём только последний октет
	octets := strings.Split(parts[0], ".")

	if len(octets) != 4 {
		return nil, fmt.Errorf("only IPv4 ranges supported")
	}

	start, err := strconv.Atoi(octets[3])
	if err != nil {
		return nil, err
	}

	end, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, err
	}

	if start > end || end > 255 {
		return nil, fmt.Errorf("invalid range %s", s)
	}

	var result []net.IP

	for i := start; i <= end; i++ {

		ipstr := fmt.Sprintf(
			"%s.%s.%s.%d",
			octets[0],
			octets[1],
			octets[2],
			i,
		)
		if ip := net.ParseIP(ipstr); ip != nil {
			result = append(result, ip)
			continue
		}
	}

	return result, nil
}
