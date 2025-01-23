package main

import (
	"net/netip"
	"testing"
)

func TestParseArgs(t *testing.T) {
	args, err := parseArgs(
		"list", "l",
		"from", "[::1]:1024",
		"to", "[::1]:9999",
		"receiver:", "a", "b",
	)

	if err != nil ||
		args.source != "[::1]:1024" ||
		args.destination != "[::1]:9999" ||
		args.listPath != "l" ||
		args.unsafe != false ||
		len(args.receiver) != 2 ||
		args.receiver[0] != "a" ||
		args.receiver[1] != "b" {
		t.Fail()
	}
}

func TestParseList(t *testing.T) {
	l := `
# Localhost.
::1
# A device.
200::
`

	list, err := parseList(l)
	if err != nil {
		t.Fail()
	}

	strings := []string{"::1", "200:0::"}
	addresses := []netip.Addr{}

	for _, v := range strings {
		address, err := netip.ParseAddr(v)
		if err != nil {
			t.Fail()
		}
		addresses = append(addresses, address)
	}

	for _, v := range addresses {
		if !list[v] {
			t.Fail()
		}
	}
}
