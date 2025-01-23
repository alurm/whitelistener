package main

import (
	"fmt"
	"bufio"
	"strings"
	"os"
	"io"
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
