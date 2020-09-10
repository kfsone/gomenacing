package main

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-shellwords"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

type Repl struct {
	db         *Database
	sdb        *SystemDatabase
	src        *bufio.Scanner
	out        io.Writer
	terminated bool
}

func (r *Repl) Write(p []byte) (int, error) {
	return r.out.Write(p)
}

func NewRepl(db *Database, sdb *SystemDatabase, src *bufio.Scanner, out io.Writer) (*Repl, error) {
	repl := Repl{db: db, sdb: sdb, src: src, out: out}
	return &repl, nil
}

func cmdSystemFind(r *Repl, args []string, _ *CommandParser) {
	name := strings.Join(args, " ")
	system := r.sdb.GetSystem(name)
	if system == nil {
		fmt.Fprintln(r, "Not found.")
	} else {
		fmt.Fprintf(r, "%+v\n", *system)
	}
}

func cmdImport(r *Repl, args []string, _ *CommandParser) {
	pathname := strings.Join(args, " ")

	// If they named a specific .gom file, go ahead and import just that.
	if strings.HasSuffix(pathname, ".gom") {
		fnImportFile(r, pathname, true)
	} else {
		stat, err := os.Stat(pathname)
		if err != nil && !os.IsExist(err) {
			fmt.Fprintf(r, "import %s: %s\n", pathname, err)
			return
		}
		if !stat.IsDir() {
			fmt.Fprintf(r, "import %s: not a directory.\n", pathname)
			return
		}

		imports := 0
		for _, filename := range GetImportFilenames() {
			if fnImportFile(r, filepath.Join(pathname, filename), false) {
				imports++
			}
		}

		if imports == 0 {
			fmt.Fprintln(r, "Nothing to import.")
		}
	}
}

func (r *Repl) Run(prompt string) error {
	parser := shellwords.NewParser()
	parser.ParseEnv = true
	parser.ParseBacktick = false

	for !r.terminated {
		if len(prompt) != 0 {
			if _, err := fmt.Fprint(r, prompt); err != nil {
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

	fmt.Fprintln(r, strings.Join(outputs, " "))
	if c.help != "" {
		fmt.Fprintln(r, c.help)
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
				fmt.Fprintln(r, "Error: Missing sub-command.")
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
				fmt.Fprintln(r, strings.Join(commandSeq, " ")+": "+parse.help)
				if parse.commands != nil {
					fmt.Fprintln(r, "Sub-commands:")
					for cmd, info := range parse.commands {
						if info.help != "" {
							fmt.Fprintf(r, "  %16s  %s\n", cmd, info.help)
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
		fmt.Fprintln(r, "Sorry, unrecognized command: ", strings.Join(commandSeq, " "))
	}
}

func toggleBool(boolPtr *bool) string {
	*boolPtr = !*boolPtr
	boolState := "off"
	if *boolPtr {
		boolState = "on"
	}
	return boolState
}

func cmdToggleWarnings(r *Repl, args []string, _ *CommandParser) {
	fmt.Fprintln(r, "Warnings are now:", toggleBool(ShowWarnings))
}

func cmdToggleOnDuplicate(r *Repl, args []string, _ *CommandParser) {
	fmt.Fprintln(r, "OnDuplicate is now:", toggleBool(ErrorOnUnknown))
}

func cmdToggleOnUnknown(r *Repl, args []string, _ *CommandParser) {
	fmt.Fprintln(r, "OnUnknown is now:", toggleBool(ErrorOnDuplicate))
}

func cmdProbe(r *Repl, args []string, _ *CommandParser) {
	// Assume 'args' is a system name
	hard := false
	if len(args) >= 1 && args[0] == "--hard" {
		hard = true
		args = args[1:]
	}
	if len(args) < 2 {
		fmt.Fprintln(r, "Please specify <distance in ly> and <system name>, e.g: probe 16.2 sol")
		return
	}
	distance, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		fmt.Fprintf(r, "Invalid distance value: %s\n", err)
		return
	}

	systemName := strings.Join(args[1:], " ")
	system := r.sdb.GetSystem(systemName)
	if system == nil {
		fmt.Fprintf(r, "Unrecognized system: %s\n", err)
		return
	}

	fmt.Fprintf(r, "Searching for systems within %.fly of %s (%d)\n", distance, system.Name(), system.GetId())

	start := time.Now()
	matches := false
	if hard {
		rangeSq := NewSquareFloat(distance)
		for _, target := range r.sdb.systemsByID {
			if distSq := Distance(system, target); distSq <= rangeSq {
				fmt.Fprintf(r, "- %8.2fly %s\n", distSq.Root(), target.Name())
				matches = true
			}
		}
	} else {
		matches, err = r.sdb.getSystemsWithinRange(system, distance, func(neighbor *System, distSq SquareFloat) bool {
			fmt.Fprintf(r, "- %8.2fly %s\n", distSq.Root(), neighbor.Name())
			return true
		})
		if err != nil {
			fmt.Fprintf(r, "Error: %s\n", err)
			return
		}
	}
	if !matches {
		fmt.Fprintln(r, "Nothing found, which is odd.")
	} else {
		fmt.Fprintf(r, "Took: %s\n", time.Since(start))
	}
}

var commands = CommandParser{
	commands: map[string]CommandParser{
		"exit":   {help: "Exit the application.", action: func(r *Repl, _ []string, _ *CommandParser) { r.terminated = true }},
		"quit":   {help: "", action: func(r *Repl, _ []string, _ *CommandParser) { r.terminated = true }},
		"import": {help: "Import data from a file or directory.", action: cmdImport},
		"stats": {help: "Show stats on current database.", action: func(r *Repl, _ []string, _ *CommandParser) {
			r.sdb.Stats(r.out)
		}},
		"set": {commands: map[string]CommandParser{
			"warnings":    {help: "Toggle warnings on/off.", action: cmdToggleWarnings},
			"onduplicate": {help: "Toggle on-duplicate on/off.", action: cmdToggleOnDuplicate},
			"onunknown":   {help: "Toggle on-unknown on/off.", action: cmdToggleOnUnknown},
		},
			help: "Change environment settings.",
		},
		"system": {commands: map[string]CommandParser{
			"find": {help: "Lookup a system by name.", action: cmdSystemFind},
			"scan": {help: "Find other systems within a given distance of a system.", action: cmdProbe},
		},
			help: "System-related commands."},
	},
	action: func(r *Repl, _ []string, cp *CommandParser) {
		cp.Info(r, "")
	},
}
