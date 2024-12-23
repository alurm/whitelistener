package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

func log(args ...any) {
	fmt.Fprint(os.Stderr, args...)
}

func usage() {
	log("Whitelistener: a TCP/IPv6 to TCP/IP reverse proxy with an IPv6 whitelist.\n")
	log("Usage: whitelistener from <source> to <destination> < <whitelist>\n")
	log("The standard input must be a list of allowed IPv6 addresses, one per line.\n")
	log("Lines starting with a hash character are treated as comments.\n")
	log("Example usage: echo ::1 | whitelistener from [::1]:1024 to [::1]:8080\n")
	os.Exit(1)
}

type configuration struct {
	source      string
	destination string
	whitelist   map[string]bool
}

func parse(r io.Reader, args ...string) (result configuration) {
	// Parse args.

	var haveSource, haveDestination bool

	for i := 0; i+1 < len(args); i += 2 {
		key := args[i]
		value := args[i+1]

		switch key {
		case "from":
			haveSource = true
			result.source = value
		case "to":
			haveDestination = true
			result.destination = value
		default:
			usage()
		}
	}

	if !haveSource || !haveDestination {
		usage()
	}

	// Read the whitelist.

	result.whitelist = map[string]bool{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 || line[0] != '#' {
			result.whitelist[line] = true
		}
	}

	if err := scanner.Err(); err != nil {
		log("An error occurred during reading the whitelist:", err)
		os.Exit(1)
	}

	return
}

func handle(client net.Conn, destination string) {
	d, err := net.Dial("tcp", destination)
	if err != nil {
		log("An error occurred while dialing the destination:", err)
		os.Exit(1)
	}

	go func() {
		io.Copy(client, d)
		client.Close()
	}()

	go func() {
		io.Copy(d, client)
		d.Close()
	}()
}
