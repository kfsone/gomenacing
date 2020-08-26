package main

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
)

// DbEntity normalizes the id/name identification of any database storable.
type DbEntity struct {
	ID     EntityID `json:"id"`
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

func NewDbEntityFromJSON(json []gjson.Result) (entity DbEntity, err error) {
	return NewDbEntity(json[0].Int(), json[1].String())
}
