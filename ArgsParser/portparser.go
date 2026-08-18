package argsparser

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	MinPort = 1
	MaxPort = 65535
)

func ParsePorts(input string) ([]uint16, error) {

	// Все порты
	if input == "-" {
		return AllPorts(), nil
	}

	var ports []uint16

	// Разделяем список
	items := strings.Split(input, ",")

	for _, item := range items {

		item = strings.TrimSpace(item)

		// Диапазон
		if strings.Contains(item, "-") {

			rangePorts, err := ParsePortRange(item)

			if err != nil {
				return nil, err
			}

			ports = append(ports, rangePorts...)
			continue
		}

		// Одиночный порт
		port, err := strconv.Atoi(item)

		if err != nil {
			return nil, fmt.Errorf("invalid port: %s", item)
		}

		if port < MinPort || port > MaxPort {
			return nil, fmt.Errorf("port out of range: %d", port)
		}

		ports = append(ports, uint16(port))
	}

	return RemoveDuplicatePorts(ports), nil
}

func ParsePortRange(value string) ([]uint16, error) {

	parts := strings.Split(value, "-")

	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid port range: %s", value)
	}

	start, err := strconv.Atoi(parts[0])

	if err != nil {
		return nil, err
	}

	end, err := strconv.Atoi(parts[1])

	if err != nil {
		return nil, err
	}

	if start < MinPort || end > MaxPort || start > end {
		return nil, fmt.Errorf("invalid port range: %s", value)
	}

	var result []uint16

	for i := start; i <= end; i++ {
		result = append(result, uint16(i))
	}

	return result, nil
}

func AllPorts() []uint16 {

	ports := make([]uint16, 0, MaxPort)

	for i := MinPort; i <= MaxPort; i++ {
		ports = append(ports, uint16(i))
	}

	return ports
}

func RemoveDuplicatePorts(ports []uint16) []uint16 {

	result := make([]uint16, 0)

	seen := make(map[uint16]bool)

	for _, port := range ports {

		if !seen[port] {

			seen[port] = true
			result = append(result, port)

		}
	}

	return result
}
