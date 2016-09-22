# nats-repl
A simple REPL console for interacting with a NATS (http://nats.io) server.

# Installation
`go get github.com/joshglendenning/nats-repl`

# Usage
## CLI
```
USAGE:
   nats-repl [global options] command [command options] [arguments...]

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --server URL, -s URL  NATS URL to connect to [$NATS_DEFAULT_URL]
   --help, -h            show help
   --version, -v         print the version
```

## Commands
- `pub <subject> [data]`
- `sub <subject>`
- `req <subject> [data]`

# Development
This project uses [gopm](https://github.com/gpmgo/gopm) for dependency management and building.

```sh
go get -u github.com/gpmgo/gopm
gopm get
GOOS=<os> GOARCH=<arch> gopm build
```
