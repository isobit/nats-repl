# nats-repl
A simple REPL console for interacting with a NATS (http://nats.io) server.

# Installation
`go get github.com/joshglendenning/nats-repl`

# Usage
```
USAGE:
   nats-repl [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --server URL, -s URL  NATS URL to connect to [$NATS_DEFAULT_URL]
   --help, -h            show help
   --version, -v         print the version
```

# Commands
- `pub <subject> [data]`
- `sub <subject>`
- `req <subject> [data]`
