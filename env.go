package main

import (
	"errors"
	"log"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

// DefaultPath is the directory data files will be stored in unless overridden.
var DefaultPath = flag.StringP("path", "p", ".", "Top-level directory to operate from.")

// DefaultDbName is the default filename for the database.
var DefaultDbName = flag.StringP("db", "d", "menace.db", "Name of the gomenace database.")

// ShowWarnings enables noisy warnings about json parsing and loading during imports.
var ShowWarnings = flag.BoolP("warnings", "W", false, "Show warnings about data imports (can be noisy).")

// Should it be an error when a duplicate system is encountered?
var ErrorOnDuplicate = flag.Bool("erronduplicate", false, "Turn duplicates in json files into errors.")

// Should it be an error when something references an unknown parent.
var ErrorOnUnknown = flag.Bool("erronunknown", false, "Make unknown system/station references in json files into errors.")

// SetupEnv prepares the environment options/flags, after
// ensuring the directory for the environment exists.
// Leave 'path' and/or 'filename' blank for default values.
func SetupEnv() error {
	// Make sure the directory actually exists.
	if _, err := ensureDirectory(*DefaultPath); err != nil {
		return err
	}
	return nil
}

// DataPath is the runtime location for data files.
func DataPath() string {
	return *DefaultPath
}

// FilterError checks if the input error should be demoted to a warning, based on
// the environment configuration. Demoted errors are then logged when warnings
// are enabled. A demoted error is indicated by a nil return.
func FilterError(err error) error {
	if errors.Is(err, ErrDuplicateEntity) {
		if *ErrorOnDuplicate {
			return err
		}
	} else if errors.Is(err, ErrUnknownEntity) {
		if *ErrorOnUnknown {
			return err
		}
	} else {
		return err
	}

	if *ShowWarnings {
		log.Printf("NOTE: %s", err.Error())
	}
	return nil
}

/// Returns the path to a file under the data directory.
func DataFilePath(pathToFile ...string) string {
	return filepath.Join(*DefaultPath, filepath.Join(pathToFile...))
}
