package main

import (
	"errors"
	"log"
	"path/filepath"
)

var (
	// DefaultDbPath is the directory data files will be stored in unless overridden.
	DefaultDbPath = "data"
	// DefaultDbFile is the default filename for the database.
	DefaultDbFile = "dangerous.db"
)

// Env is a holder of runtime environment options/configuration values.
type Env struct {
	dataPath         string // path to the data directory.
	dataFile         string // full path and filename of the data file.
	ErrorOnDuplicate bool   // reject duplicate systems or facilities.
	ErrorOnUnknown   bool   // reject references to unknown systems etc.
	SilenceWarnings  bool   // don't log warnings.
}

// NewEnv creates an environment object for holding options/flags, after
// ensuring the directory for the environment exists.
// Leave 'path' and/or 'filename' blank for default values.
func NewEnv(path string, filename string) (*Env, error) {
	if path == "" {
		path = DefaultDbPath
	}
	if filename == "" {
		filename = DefaultDbFile
	}

	// Create the directory if we need to
	if _, err := ensureDirectory(path); err != nil {
		return nil, err
	}

	return &Env{
		dataPath: path,
		dataFile: filepath.Join(path, filename),
	}, nil
}

// DataPath is the runtime location for data files.
func (env Env) DataPath() string {
	return env.dataPath
}

// DataFile is the runtime path and filename of the database file.
func (env Env) DataFile() string {
	return env.dataFile
}

// FilterError checks if the input error should be demoted to a warning, based on
// the environment configuration. Demoted errors are then logged unless warnings
// are silenced. A demoted error is indicated by a nil return.
func (env Env) FilterError(err error) error {
	if errors.Is(err, ErrDuplicateEntity) {
		if env.ErrorOnDuplicate {
			return err
		}
	} else if errors.Is(err, ErrUnknownEntity) {
		if env.ErrorOnUnknown {
			return err
		}
	} else {
		return err
	}

	if !env.SilenceWarnings {
		log.Printf("NOTE: %s", err.Error())
	}
	return nil
}
