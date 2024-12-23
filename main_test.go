package main

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	r := strings.NewReader(
		`# Localhost.
::1
# A device.
200:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx
`,
	)
	p := parse(r, "from", "[::1]:1024", "to", "[::1]:9999")
	if p.source != "[::1]:1024" || p.destination != "[::1]:9999" {
		t.Fail()
	}
	whitelist := []string{"::1", "200:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx"}
	for _, v := range whitelist {
		if !p.whitelist[v] {
			t.Fail()
		}
	}
}
