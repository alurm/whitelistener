package main

import (
	"net"
	"os"
)

func main() {
	conf := parse(os.Stdin, os.Args[1:]...)

	server, err := net.Listen("tcp6", conf.source)
	if err != nil {
		log("An error: listen:", err)
		os.Exit(1)
	}

	for {
		client, err := server.Accept()
		if err != nil {
			log("An error occurred while accepting a connection:", err)
			os.Exit(1)
		}

		host, _, err := net.SplitHostPort(client.RemoteAddr().String())
		if err != nil {
			log("Bug: a remote address didn't split into host and port parts: ", err)
			os.Exit(1)
		}

		if !conf.whitelist[host] {
			client.Close()
		} else {
			handle(client, conf.destination)
		}
	}
}
