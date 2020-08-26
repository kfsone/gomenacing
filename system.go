package main

import (
	"strings"
	"time"

	"github.com/tidwall/gjson"
)

type System struct {
	// System describes a star system within Elite Dangerous.
	DbEntity
	Position   Coordinate  `json:"pos"`
	Permit     bool        `json:"permit"`
	Facilities []*Facility `json:"-"`
	Updated    time.Time   `json:"updated"`
}

func NewSystem(entity DbEntity, position Coordinate, permit bool) (*System, error) {
	return &System{DbEntity: entity, Position: position, Permit: permit}, nil
}

func NewSystemFromJson(json []gjson.Result) (system *System, err error) {
	if entity, err := NewDbEntityFromJSON(json); err == nil {
		position := Coordinate{json[2].Float(), json[3].Float(), json[4].Float()}
		system, err = NewSystem(entity, position, json[5].Bool())
	}
	return
}

func (s System) Distance(to *System) Square {
	return s.Position.Distance(to.Position)
}

func (s System) Name() string {
	return s.DbName
}

func (s System) String() string {
	return s.DbName
}

func (s System) GetFacility(name string) *Facility {
	for _, facility := range s.Facilities {
		if strings.EqualFold(name, facility.DbName) {
			return facility
		}
	}
	return nil
}
