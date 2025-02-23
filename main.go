package main

import "flag"

func main() {
	mode := flag.String("mode", "initiator", "Mode of operation: initiator or target")
	flag.Parse()

	switch *mode {
	case "initiator":
		initiator()
	case "target":
		target()
	default:
		panic("Invalid mode. Please specify either 'initiator' or 'target'")
	}
}
