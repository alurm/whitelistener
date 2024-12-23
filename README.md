# Whitelistener: a TCP/IPv6 to TCP/IP reverse proxy with an IPv6 whitelist

Usage: `whitelistener from <source> to <destination> < <whitelist>`

The standard input must be a list of allowed IPv6 addresses, one per line.

Lines starting with a hash character are treated as comments.

Example usage: `echo ::1 | whitelistener from [::1]:1024 to [::1]:8000`

Building: `go build`

Testing: `go test`

License: MIT

