package argsparser

import (
	"flag"
	"fmt"
	"time"
)

func ParseTimeout(ms int) (time.Duration, error) {
	if ms <= 0 {
		return 0, fmt.Errorf("timeout must be greater than 0 ms")
	}

	return time.Duration(ms) * time.Millisecond, nil
}

func ValidateArgs(args []string) (*Config, error) {

	fs := flag.NewFlagSet("PortScanner", flag.ContinueOnError)

	portArg := fs.String(
		"p",
		"80",
		"ports to scan",
	)

	timeoutArg := fs.Int(
		"t",
		1000,
		"connection timeout in milliseconds",
	)

	countArg := fs.Uint(
		"n",
		1,
		"number of packets to send",
	)

	modeArg := fs.String(
		"m",
		"tcp",
		"scan mode: tcp, udp or both",
	)

	err := fs.Parse(args)
	if err != nil {
		return nil, err
	}

	if *countArg < 1 {
		return nil, fmt.Errorf(
			"packet count must be greater than 0",
		)
	}

	timeout, err := ParseTimeout(*timeoutArg)
	if err != nil {
		return nil, err
	}

	mode, err := ParseScanMode(*modeArg)
	if err != nil {
		return nil, err
	}

	targetArgs := fs.Args()

	if len(targetArgs) == 0 {
		return nil, fmt.Errorf("no targets specified")
	}

	targets, err := ParseTargets(targetArgs)
	if err != nil {
		return nil, err
	}

	ports, err := ParsePorts(*portArg)
	if err != nil {
		return nil, err
	}

	return &Config{
		Targets:      targets,
		Ports:        ports,
		PacketsCount: *countArg,
		Timeout:      timeout,
		Mode:         mode,
	}, nil
}
