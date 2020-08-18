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

func importEddbData(db *Database) {
	errorsCh := make(chan error, 2)
	defer close(errorsCh)
	go func() {
		errorsCh <- importFacilities(db)
	}()
	go func() {
		errorsCh <- importSystems(db)
	}()
	go func() {
		errorsCh <- importCommodities(db)
	}()
	errorCount := 0
	for i := 0; i < 3; i++ {
		if err := <-errorsCh; err != nil {
			log.Print(err)
			errorCount += 1
		}
	}
	if errorCount > 0 {
		log.Fatal("Stopping because of errors.")
	}
}

func main() {
	flag.Parse()
	failOnError(SetupEnv())

	fmt.Println("GoMenacing v0.01 (C) Oliver 'kfsone' Smith, 2020")

	// Populate the parameters.
	/// TODO: Parse arguments

	var db *Database
	db, err := GetDatabase(*DefaultPath)
	failOnError(err)

	if *EddbPath != "" {
		importEddbData(db)
	}

	sdb := NewSystemDatabase()
	failOnError(db.LoadData(sdb))

	reader := bufio.NewReader(os.Stdin)

	repl, err := NewRepl(db, bufio.NewScanner(reader), os.Stdout)
	failOnError(err)

	if db != nil {
		err = repl.Run("GoM> ")
		failOnError(err)
	}
}
