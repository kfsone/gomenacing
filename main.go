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

func main() {
	flag.Parse()
	SetupEnv()

	fmt.Println("GoMenacing v0.01 (C) Oliver 'kfsone' Smith, 2020")

	// Populate the parameters.
	/// TODO: Parse arguments

	var db *Database
	//db, err := OpenDatabase()
	//fatalize(err)
	//defer db.Close()

	// Parse the systems file.

	sdb := NewSystemDatabase()
	fatalize(sdb.importSystems())
	fatalize(sdb.importFacilities())

	reader := bufio.NewReader(os.Stdin)

	repl, err := NewRepl(db, bufio.NewScanner(reader), os.Stdout)
	fatalize(err)

	if db != nil {
		err = repl.Run("GoM> ")
		fatalize(err)
	}
}
