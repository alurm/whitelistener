Whitelistener is a TCP/IP to TCP/IP reverse proxy with an optional whitelist.

Syntax: whitelistener ( unsafe | list <list> ) from <source> to <destination> [ receiver: args... ]

List must be a path to a file with a list of allowed IPv6 addresses, one per line.
Empty lines and lines starting with a pound sign are ignored.

If "unsafe" is specified, all connections are allowed.

If "receiver:" is specified, a child will be launched specified by args.
Whitelistener will terminate after the child dies.

Example: whitelistener list <(echo ::1) from [::1]:1024 to [::1]:8000 receiver: nc -l ::1 8000

The license for this project is MIT.
