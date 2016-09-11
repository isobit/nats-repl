package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"time"
	"io"
	"github.com/urfave/cli"
	"github.com/chzyer/readline"
	"github.com/nats-io/nats"
)

func colorize(code interface{}, s string) string {
	return fmt.Sprintf("\x1b[%vm%s\x1b[0m", code, s)
}
func logInfo(s string) {
	fmt.Printf("[INFO] %s\n", s)
}
func logWarn(s string) {
	fmt.Printf("[%s] %s\n", colorize(33, "WARNING"), s)
}
func logError(s string) {
	fmt.Printf("[%s] %s\n", colorize(31, "ERROR"), s)
}

func main() {
	app := cli.NewApp()
	app.Name = "nats-repl"
	app.Usage = "REPL for NATS"
	app.Flags = []cli.Flag {
		cli.StringFlag{
			Name: "server, s",
			EnvVar: "NATS_DEFAULT_URL",
			Usage: "NATS `URL` to connect to",
		},
	}
	app.Action = func(ctx *cli.Context) error {
		nc, err := nats.Connect(ctx.String("server"))
		if err != nil {
			logError(fmt.Sprintf("%v", err))
			os.Exit(1)
		}

		logInfo(fmt.Sprintf("connected to %s", ctx.String("server")))

		rl, err := readline.NewEx(&readline.Config{
			Prompt:          colorize("1;37", "nats> "),
			HistoryFile:     "/tmp/nats-repl-history.tmp",
			InterruptPrompt: "^C",
			EOFPrompt:       "exit",
		})
		if err != nil {
			panic(err)
		}
		defer rl.Close()

		repl:
		for {
			line, err := rl.Readline()
			switch err {
			case readline.ErrInterrupt:
				if len(line) == 0 {
					break repl
				} else {
					continue repl
				}
			case io.EOF:
				break repl
			}

			args := strings.Fields(line)
			switch {
			case len(args) == 0:
				continue repl
			case args[0] == "pub":
				var subject string
				if len(args) >= 2 {
					subject = args[1]
				} else {
					logError("subject is required")
					continue repl
				}
				var data string
				if len(args) >= 3 {
					data = strings.Join(args[2:], " ")
				} else {
					data = ""
				}
				nc.Publish(subject, []byte(data))
			case args[0] == "sub":
				sigch := make(chan os.Signal, 1)
				signal.Notify(sigch, os.Interrupt)

				var subject string
				if len(args) >= 2 {
					subject = args[1]
				} else {
					logError("subject is required")
					continue repl
				}
				subch := make(chan *nats.Msg, 64)
				sub, _ := nc.ChanSubscribe(subject, subch)

				sub:
				for {
					select {
					case msg := <-subch:
						fmt.Printf("[%s] %s\n", msg.Subject, string(msg.Data))
					case <-sigch:
						fmt.Println()
						break sub
					}
				}
				close(sigch)
				close(subch)
				signal.Reset(os.Interrupt)
				sub.Unsubscribe()
			case args[0] == "req":
				var subject string
				if len(args) >= 2 {
					subject = args[1]
				} else {
					logError("subject is required")
					continue repl
				}
				var data string
				if len(args) >= 3 {
					data = strings.Join(args[2:], " ")
				} else {
					data = ""
				}
				msg, err := nc.Request(subject, []byte(data), 5000*time.Millisecond)
				if err != nil {
					logError(fmt.Sprintf("%v", err))
					break
				}
				fmt.Println(string(msg.Data))
			case args[0] == "help":
				logInfo("COMMANDS:")
				logInfo("pub <subject> [data]")
				logInfo("sub <subject>")
				logInfo("req <subject> [data]")
			case args[0] == "exit":
				break repl
			default:
				logError(fmt.Sprintf("unknown command: %s", args[0]))
			}
		}
		return nil
	}
	app.Run(os.Args)
}
