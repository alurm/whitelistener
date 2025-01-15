package main

import (
	_ "embed"
	"fmt"
	"net"
	"os"
)

//go:embed README.md
var readme string

func main() {
	if err := whitelisten(); err != nil {
		fmt.Fprintf(os.Stderr, "An error: %v.\n\n%s", err, readme)
		os.Exit(1)
	}
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

		if !conf.whitelist[host] {
			client.Close()
		} else {
			handle(client, conf.destination)
		}
	}

	return nil
}
