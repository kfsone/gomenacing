package main

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-shellwords"
	"io"
	"log"
)

type Database struct{}

type Repl struct {
	db         *Database
	src        *bufio.Scanner
	out        io.StringWriter
	terminated bool
}

func NewRepl(db *Database, src *bufio.Scanner, out io.StringWriter) (*Repl, error) {
	repl := Repl{db: db, src: src, out: out}
	return &repl, nil
}

func (r *Repl) evaluate(args []string) {
	switch args[0] {
	case "exit":
		r.terminated = true
		break

	case "quit":
		r.terminated = true
		break

	case "help":
		r.out.WriteString("Commands: exit, quit.\n")
		break

	default:
		r.out.WriteString(fmt.Sprintf("Unrecognized command '%s', try 'help' if you're stuck.\n", args[0]))
	}
}

func (r *Repl) Run(prompt string) error {
	parser := shellwords.NewParser()
	parser.ParseEnv = true
	parser.ParseBacktick = false

	for !r.terminated {
		if len(prompt) != 0 {
			if _, err := r.out.WriteString(prompt); err != nil {
				return err
			}
		}

		if !r.src.Scan() {
			break
		}

		args, err := parser.Parse(r.src.Text())
		if err != nil {
			log.Printf(err.Error())
			continue
		}
		if len(args) >= 0 {
			r.evaluate(args)
		}
	}

	return r.src.Err()
}
