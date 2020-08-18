package main

import (
	"fmt"
	"strings"
)

// Anything that can go into the database is a DbEntity.
type DbEntity struct {
	Id     EntityID `json:"id"`
	DbName string   `json:"name"`
}

func NewDbEntity(id int64, name string) (entity DbEntity, err error) {
	if id > 0 && id < 1<<32 {
		dbName := strings.TrimSpace(name)
		if len(dbName) > 0 {
			entity = DbEntity{EntityID(id), dbName}
		} else {
			err = fmt.Errorf("invalid/empty name: \"%s\"", name)
		}
	} else {
		err = fmt.Errorf("invalid id: %d", id)
	}
	return
}
