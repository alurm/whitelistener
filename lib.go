package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
)

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

func listen() error {
	conf, err := parse(os.Args[1:]...)
	if err != nil {
		return err
	}

	if len(conf.receiver) != 0 {
		name := conf.receiver[0]
		conf.receiver = conf.receiver[1:]
		go func() {
			cmd := exec.Command(name, conf.receiver...)
			// Fix me: inherit ALL file descriptors.
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()
			if err == nil {
				os.Exit(0)
			}
			switch err := err.(type) {
			case *exec.ExitError:
				os.Exit(err.ExitCode())
			default:
				fmt.Fprintf(os.Stderr, "An error: receiver: %v.\n\n\n", err)
				os.Exit(1)
			}
		}()
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

		if conf.list == nil || conf.list[host] {
			handle(client, conf.destination)
		} else {
			client.Close()
		}
	}
}
