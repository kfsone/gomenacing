package main

import (
	"errors"
	"fmt"
	"strings"
)

type System struct {
	// System describes a star system within Elite Dangerous.
	DbEntity
	Position   Coordinate  `json:"pos"`
	Permit     bool        `json:"permit"`
	Facilities []*Facility `json:"-"`
}

func NewSystem(id int64, dbName string, position Coordinate, permit bool) (*System, error) {
	if id <= 0 {
		return nil, errors.New(fmt.Sprintf("invalid system id: %d", id))
	}
	if id >= (1 << 32) {
		return nil, errors.New(fmt.Sprintf("invalid system id (too large): %d", id))
	}
	name := strings.TrimSpace(dbName)
	if len(name) == 0 {
		return nil, errors.New("empty system name")
	}
	return &System{DbEntity: DbEntity{EntityID(id), strings.ToUpper(name)}, Position: position, Permit: permit}, nil
}

func (s *System) NewFacility(id int64, dbName string, features FacilityFeatureMask) (*Facility, error) {
	if id <= 0 || id >= 1<<32 {
		return nil, errors.New(fmt.Sprintf("invalid facility id: %d", id))
	}
	dbName = strings.TrimSpace(dbName)
	if len(dbName) == 0 {
		return nil, errors.New("invalid (empty) facility name")
	}

	facility := &Facility{
		DbEntity: DbEntity{
			Id:     EntityID(id),
			DbName: strings.ToUpper(dbName),
		},
		System:   s,
		Features: features,
	}

	return facility, nil
}

func (s System) Distance(to *System) Square {
	return s.Position.Distance(to.Position)
}

func (s System) Name(_ int) string {
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
