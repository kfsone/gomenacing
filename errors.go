package main

import (
	"errors"
)

// ErrDuplicateEntity represents detection that an entity ID is being reused.
var ErrDuplicateEntity = errors.New("duplicate")

// ErrUnknownEntity represents detection that an ID references an unknown entity.
var ErrUnknownEntity = errors.New("unknown")
