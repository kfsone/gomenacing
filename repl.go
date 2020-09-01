package main

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-shellwords"
	"io"
	"log"
	"strings"
)

type Repl struct {
	db         *Database
	sdb        *SystemDatabase
	src        *bufio.Scanner
	out        io.StringWriter
	terminated bool
}

func NewRepl(db *Database, sdb *SystemDatabase, src *bufio.Scanner, out io.StringWriter) (*Repl, error) {
	repl := Repl{db: db, sdb: sdb, src: src, out: out}
	return &repl, nil
}

func cmdSystemFind(r *Repl, args []string, _ *CommandParser) {
	name := strings.Join(args, " ")
	system := r.sdb.GetSystem(name)
	if system == nil {
		r.out.WriteString("Not found.\n")
	} else {
		r.out.WriteString(fmt.Sprintf("%+v\n", *system))
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
			commands.Parse(r, args)
		}
	}

	return r.src.Err()
}

type CommandParser struct {
	commands map[string]CommandParser
	help     string
	action   func(*Repl, []string, *CommandParser)
}

func (c CommandParser) Info(r *Repl, command string) {
	outputs := make([]string, 0, 16)

	if command != "" {
		outputs = append(outputs, command, "sub-commands:")
	} else {
		outputs = append(outputs, "Commands:")
	}

	for cmd, info := range c.commands {
		if info.help != "" {
			outputs = append(outputs, cmd)
		}
	}

	r.out.WriteString(strings.Join(outputs, " ") + "\n")
	if c.help != "" {
		r.out.WriteString(c.help + "\n")
	}
}

func (c CommandParser) Parse(r *Repl, args []string) {
	parse := &c
	commandSeq := make([]string, 0, 16)
	for {
		if len(args) == 0 || parse.commands == nil {
			if parse.action != nil {
				parse.action(r, args, parse)
				return
			}
			if parse.commands != nil {
				r.out.WriteString("Error: Missing sub-command.\n")
				parse.Info(r, strings.Join(commandSeq, " "))
				return
			}
			panic("Invalid command: " + strings.Join(commandSeq, " "))
		}
		command := strings.ToLower(args[0])
		commandSeq = append(commandSeq, command)
		args = args[1:]
		subParse, exists := parse.commands[command]
		if exists {
			parse = &subParse
			continue
		}
		if command == "help" {
			if parse.help != "" {
				r.out.WriteString(strings.Join(commandSeq, " ") + ": " + parse.help + "\n")
				if parse.commands != nil {
					r.out.WriteString("Sub-commands:\n")
					for cmd, info := range parse.commands {
						if info.help != "" {
							r.out.WriteString(fmt.Sprintf("  %16s  %s\n", cmd, info.help))
						}
					}
				}
			}
			break
		}
		if parse.action != nil {
			parse.action(r, args, parse)
			break
		}
		r.out.WriteString("Sorry, unrecognized command: " + strings.Join(commandSeq, " ") + "\n")
	}
}

var commands = CommandParser{
	commands: map[string]CommandParser{
		"exit": {help: "Exit the application.", action: func(r *Repl, _ []string, _ *CommandParser) { r.terminated = true }},
		"quit": {help: "", action: func(r *Repl, _ []string, _ *CommandParser) { r.terminated = true }},
		"system": {commands: map[string]CommandParser{
			"find": {help: "Lookup a system by name.", action: cmdSystemFind}},
			help: "System-related commands."},
	},
	action: func(r *Repl, _ []string, cp *CommandParser) {
		cp.Info(r, "")
	},
}
