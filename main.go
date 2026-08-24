package main

import (
	argsparser "PortScaner/ArgsParser"
	"fmt"
	"os"
)

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

}
