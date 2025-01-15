package main

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	r := strings.NewReader(
		`
# Localhost.
::1
# A device.
200:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx
`,
	)
	conf, err := parse(r, "from", "[::1]:1024", "to", "[::1]:9999")
	if err != nil || conf.source != "[::1]:1024" || conf.destination != "[::1]:9999" {
		t.Fail()
	}
	whitelist := []string{"::1", "200:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx:xxxx"}
	for _, v := range whitelist {
		if !conf.whitelist[v] {
			t.Fail()
		}
	}
}
