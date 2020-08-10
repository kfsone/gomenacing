package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func fatalize(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}

func main() {
	fmt.Println("GoMenacing v0.01 (C) Oliver 'kfsone' Smith, 2020")

	// Create a destination for parameters.
	env, err := NewEnv("", "")
	env.SilenceWarnings = true
	fatalize(err)

	// Populate the parameters.
	/// TODO: Parse arguments

	var db *Database
	//db, err := OpenDatabase(env)
	//fatalize(err)
	//defer db.Close()

	// Parse the systems file.

	sdb := NewSystemDatabase()
	fatalize(sdb.loadPopulatedSystems(env))
	fatalize(sdb.loadFacilities(env))

	reader := bufio.NewReader(os.Stdin)

	repl, err := NewRepl(env, db, bufio.NewScanner(reader), os.Stdout)
	fatalize(err)

	err = repl.Run("GoM> ")
	fatalize(err)
}
