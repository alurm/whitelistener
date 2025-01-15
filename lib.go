package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

type configuration struct {
	source      string
	destination string
	whitelist   map[string]bool
}

func parse(r io.Reader, args ...string) (configuration, error) {
	var res configuration

	// Parse args.

	var haveSource, haveDestination bool

	for i := 0; i+1 < len(args); i += 2 {
		key := args[i]
		value := args[i+1]

		switch key {
		case "from":
			haveSource = true
			res.source = value
		case "to":
			haveDestination = true
			res.destination = value
		default:
			return res, fmt.Errorf("an unknown flag: %s", key)
		}
	}

	if !haveSource {
		return res, fmt.Errorf("the source must be provided")
	}

	if !haveDestination {
		return res, fmt.Errorf("the destination must be provided")
	}

	// Read the whitelist.

	res.whitelist = map[string]bool{}

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) != 0 && line[0] != '#' {
			res.whitelist[line] = true
		}
	}

	if err := scanner.Err(); err != nil {
		return res, fmt.Errorf("during reading the whitelist: %v", err)
	}

	return res, nil
}

func handle(client net.Conn, destination string) error {
	d, err := net.Dial("tcp", destination)
	if err != nil {
		return fmt.Errorf("while dialing the destination: %v", err)
	}

	go func() {
		io.Copy(client, d)
		client.Close()
	}()

	go func() {
		io.Copy(d, client)
		d.Close()
	}()

	return nil
}
