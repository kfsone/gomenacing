package main

import (
	"fmt"
	"strings"
)

// DbEntity normalizes the id/name identification of any database storable.
type DbEntity struct {
	ID     EntityID
	DbName string
}

func NewDbEntity(id int64, name string) (DbEntity, error) {
	if dbName, err := validateEntity(id, name); err == nil {
		return DbEntity{ID: EntityID(id), DbName: dbName}, nil
	} else {
		return DbEntity{}, err
	}
}

func (e DbEntity) GetId() uint32 {
	return uint32(e.ID)
}

func (e DbEntity) GetName() string {
	return e.DbName
}

func validateEntity(id int64, name string) (string, error) {
	if id <= 0 || id >= 1<<32 {
		return "", fmt.Errorf("invalid id: %d", id)
	}

	dbName := strings.TrimSpace(name)
	if len(dbName) <= 1 {
		return "", fmt.Errorf("invalid/empty name: \"%s\"", name)
	}
	return dbName, nil
}
