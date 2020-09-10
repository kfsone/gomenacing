package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	flag "github.com/spf13/pflag"
)

func failOnError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	flag.Parse()
	failOnError(SetupEnv())
	doImports := *eddbPath != ""

	fmt.Println("GoMenacing v0.01 (C) Oliver 'kfsone' Smith, 2020")

	// Populate the parameters.
	/// TODO: Parse arguments

	var db *Database
	db, err := OpenDatabase(*DefaultPath, *DefaultDbName)
	failOnError(err)
	defer db.Close()

	sdb := NewSystemDatabase(db)
	if doImports {
		fmt.Printf("import not implemented")
		//importEddbData(db)
	}
	failOnError(db.LoadDatabase(sdb))

	reader := bufio.NewReader(os.Stdin)

	repl, err := NewRepl(db, sdb, bufio.NewScanner(reader), os.Stdout)
	failOnError(err)

	if db != nil {
		err = repl.Run("GoM> ")
		failOnError(err)
	}
}
