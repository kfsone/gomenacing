package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	flag "github.com/spf13/pflag"
)

func fatalize(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func importEddbData(db *Database) {
	errorsCh := make(chan error, 2)
	defer close(errorsCh)
	go func () {
		errorsCh <- importFacilities(db)
	} ()
	go func () {
		errorsCh <- importSystems(db)
	} ()
	firstErr, secondErr := <-errorsCh, <-errorsCh
	fatalize(firstErr)
	fatalize(secondErr)
}


func main() {
	flag.Parse()
	SetupEnv()

	fmt.Println("GoMenacing v0.01 (C) Oliver 'kfsone' Smith, 2020")

	// Populate the parameters.
	/// TODO: Parse arguments

	var db *Database
	db, err := GetDatabase(*DefaultPath)
	fatalize(err)

	if *EddbPath != "" {
		importEddbData(db)
	}

	sdb := NewSystemDatabase()
	fatalize(db.LoadData(sdb))

	reader := bufio.NewReader(os.Stdin)

	repl, err := NewRepl(db, bufio.NewScanner(reader), os.Stdout)
	fatalize(err)

	if db != nil {
		err = repl.Run("GoM> ")
		fatalize(err)
	}
}
