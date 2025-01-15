# Whitelistener: a TCP/IP to TCP/IP reverse proxy with an optional whitelist

## Usage

```sh
whitelistener from <source> to <destination> [ option unsafe ]
```

The standard input must be a list of allowed IPv6 addresses, one per line.
Empty lines and lines starting with a pound sign are ignored.

If `option unsafe` is given, the whitelist is not read and all connections are proxied.

Example: `echo ::1 | whitelistener from [::1]:1024 to [::1]:8000`

## Development

To build run `go build`, to test run `go test`.

License: MIT.
