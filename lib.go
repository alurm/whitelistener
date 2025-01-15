package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
)

type configuration struct {
	source      string
	destination string
	whitelist   map[string]bool
	unsafe bool
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
		case "option":
			switch value {
			case "unsafe":
				res.unsafe = true
			default:
				return res, fmt.Errorf("an unknown option: %s", value)
			}
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

	// If we are allowing all connections (res.unsafe), the whitelist is not even read.
	if !res.unsafe {
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

func whitelisten() error {
	conf, err := parse(os.Stdin, os.Args[1:]...)
	if err != nil {
		return err
	}

	server, err := net.Listen("tcp", conf.source)
	if err != nil {
		return fmt.Errorf("during preparing to listen: %v", err)
	}

	for {
		client, err := server.Accept()
		if err != nil {
			return fmt.Errorf("during accepting a connection: %v", err)
		}

		address := client.RemoteAddr().String()
		host, _, err := net.SplitHostPort(address)
		if err != nil {
			return fmt.Errorf("a bug: a remote address %v didn't split into host and port parts: %v", address, err)
		}

		if conf.unsafe || conf.whitelist[host] {
			handle(client, conf.destination)
		} else {
			client.Close()
		}
	}

	return nil
}
