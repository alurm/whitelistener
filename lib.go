package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
)

type arguments struct {
	source      string
	destination string
	listPath    string
	unsafe      bool
	receiver    []string
}

type configuration struct {
	source      string
	destination string
	receiver    []string
	// If nil, all connections are allowed.
	list map[string]bool
}

func parseArgs(args ...string) (arguments, error) {
	var res arguments

For:
	for len(args) != 0 {
		key := args[0]
		args = args[1:]

		switch key {
		case "from":
			if len(args) == 0 {
				return res, fmt.Errorf("expected the source address")
			}
			res.source = args[0]
			args = args[1:]
		case "to":
			if len(args) == 0 {
				return res, fmt.Errorf("expected the destination address")
			}
			res.destination = args[0]
			args = args[1:]
		case "list":
			if len(args) == 0 {
				return res, fmt.Errorf("expected the whitelist")
			}
			res.listPath = args[0]
			args = args[1:]
		case "unsafe":
			res.unsafe = true
		case "receiver:":
			if len(args) == 0 {
				return res, fmt.Errorf("expected args")
			}
			res.receiver = args
			break For
		default:
			return res, fmt.Errorf("an unknown argument: %s", key)
		}
	}

	if res.source == "" {
		return res, fmt.Errorf("the source must be specified")
	}

	if res.destination == "" {
		return res, fmt.Errorf("the destination must be specified")
	}

	if res.unsafe && res.listPath != "" {
		return res, fmt.Errorf("can't specify both \"unsafe\" and \"list\"")
	}

	if !res.unsafe && res.listPath == "" {
		return res, fmt.Errorf("either \"unsafe\" or \"list\" have to be specified")
	}

	return res, nil
}

func parseList(s string) map[string]bool {
	list := map[string]bool{}

	scanner := bufio.NewScanner(strings.NewReader(string(s)))
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) != 0 && line[0] != '#' {
			list[line] = true
		}
	}

	return list
}

func parse(argsSlice ...string) (configuration, error) {
	args, err := parseArgs(argsSlice...)

	if err != nil {
		return configuration{}, err
	}

	res := configuration{
		source:      args.source,
		destination: args.destination,
		receiver:    args.receiver,
	}

	if !args.unsafe {
		file, err := os.Open(args.listPath)

		if err != nil {
			return res, fmt.Errorf("during opening the list: %v", err)
		}

		content, err := io.ReadAll(file)
		if err != nil {
			return res, fmt.Errorf("during reading the list: %v", err)
		}

		res.list = parseList(string(content))
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
