package main

import (
	"fmt"
	"os"
	"os/exec"
	"net/netip"
)

type configuration struct {
	source      string
	destination string
	receiver    []string
	// If nil, all connections are allowed.
	list map[netip.Addr]bool
	// To-do: add netip.Prefixes.
}

func child(argv []string) {
	cmd := exec.Command(argv[0], argv[1:]...)

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
}

func listen() error {
	conf, err := parse(os.Args[1:]...)
	if err != nil {
		return err
	}

	if len(conf.receiver) != 0 {
		go child(conf.receiver)
	}

	return tcp(conf.source, conf.destination, conf.list)
}
