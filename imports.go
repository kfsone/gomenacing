package main

import (
	"fmt"
	"github.com/kfsone/gomenacing/pkg/gomschema"
	"google.golang.org/protobuf/proto"
	"os"
)

func GetImportFilenames() []string {
	return []string{"commodities.gom", "systems.gom", "stations.gom", "listings.gom"}
}

/*
 * Import Implementation:
 * The data needs to be written into both the datastore and the in-memory listings,
 * Do we write to the DB and register with the in-memory?
 * Or do we write to the DB and mark the in-memory as dirty (needs reload)?
 * For now it's going to do both.
 */

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

	schema, err := getSchemaForMessage(r.db, *gomFile.Item())
	if err != nil {
		fmt.Println(r, "Error:", err)
		return false
	}
	defer schema.Close()

	count := 0
	err = gomFile.Read(func(message proto.Message, index uint) (err error) {
		err = r.sdb.registerFromMessage(message, schema)
		if err == nil {
			count++
		}
		return FilterError(err)
	})

	fmt.Fprintf(r, "%s: read %d items.\n", pathname, count)

	return true
}
