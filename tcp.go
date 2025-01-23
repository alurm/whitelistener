package main

import (
	"fmt"
	"net"
	"io"
	"net/netip"
)

func forwardTCP(client net.Conn, destination string) error {
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

func tcp(source, destination string, list map[netip.Addr]bool) error {
	listener, err := net.Listen("tcp", source)
	if err != nil {
		return fmt.Errorf("during preparing to listen: %v", err)
	}

	for {
		client, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("during accepting a connection: %v", err)
		}

		hostPort := client.RemoteAddr().String()

		host, _, err := net.SplitHostPort(hostPort)
		if err != nil {
			return fmt.Errorf("a bug: a remote address %v didn't split into host and port parts: %v", hostPort, err)
		}

		address, err := netip.ParseAddr(host)
		if err != nil {
			return fmt.Errorf("failed to parse the address %v: %v", address, err)
		}

		if list == nil || list[address] {
			forwardTCP(client, destination)
		} else {
			client.Close()
		}
	}
}
