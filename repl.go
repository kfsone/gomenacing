package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kfsone/gomenacing/pkg/gomschema"
	"github.com/mattn/go-shellwords"
	"google.golang.org/protobuf/proto"
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

func importCommodity(r *Repl, item *gomschema.Commodity) error {
	return nil
}

func importSystem(r *Repl, item *gomschema.System) error {
	return nil
}

func importFacility(r *Repl, item *gomschema.Facility) error {
	return nil
}

func importListings(r *Repl, item *gomschema.FacilityListing) error {
	return nil
}

func fnImportFile(r *Repl, pathname string, required bool) bool {
	file, err := os.Open(pathname)
	if err != nil {
		if required {
			fmt.Fprintln(r, file, ": error opening file: ", err)
		}
		return false
	}
	defer func() { Must(file.Close()) }()

	gomFile, err := gomschema.OpenGOMFile(file)
	if err != nil {
		fmt.Fprintln(r, file, ": ", err)
		return false
	}
	defer gomFile.Close()

	count := 0
	err = gomFile.Read(func(message proto.Message, index uint) (err error) {
		switch m := message.(type) {
		case *gomschema.Commodity:
			err = importCommodity(r, m)

		case *gomschema.System:
			err = importSystem(r, m)

		case *gomschema.Facility:
			err = importFacility(r, m)

		case *gomschema.FacilityListing:
			err = importListings(r, m)

		default:
			return fmt.Errorf("unrecognized message type: %t", message)
		}
		if err == nil {
			count++
		}
		return err
	})

	fmt.Fprintf(r, "%s: read %d items.\n", pathname, count)

	return true
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
		for _, filename := range []string{"commodities.gom", "systems.gom", "stations.gom", "listings.gom"} {
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

var commands = CommandParser{
	commands: map[string]CommandParser{
		"exit":   {help: "Exit the application.", action: func(r *Repl, _ []string, _ *CommandParser) { r.terminated = true }},
		"quit":   {help: "", action: func(r *Repl, _ []string, _ *CommandParser) { r.terminated = true }},
		"import": {help: "Import data from a file or directory.", action: cmdImport},
		"system": {commands: map[string]CommandParser{
			"find": {help: "Lookup a system by name.", action: cmdSystemFind}},
			help: "System-related commands."},
	},
	action: func(r *Repl, _ []string, cp *CommandParser) {
		cp.Info(r, "")
	},
}
