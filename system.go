package main

import (
	"fmt"
	gom "github.com/kfsone/gomenacing/pkg/gomschema"
	"strings"
	"time"
)

type System struct {
	// System describes a star system within Elite Dangerous.
	DbEntity
	TimestampUtc  time.Time
	position      Coordinate
	Populated     bool
	NeedsPermit   bool
	SecurityLevel gom.SecurityLevel
	Government    gom.GovernmentType
	Allegiance    gom.AllegianceType

	facilities []*Facility
}

func NewSystem(dbEntity DbEntity, position Coordinate) *System {
	return &System{
		DbEntity:      dbEntity,
		position:      position,
	}
}

func (s *System) GetDbId() string {
	return fmt.Sprintf("%08x", s.DbEntity.ID)
}

func (s *System) GetFacility(name string) *Facility {
	for _, facility := range s.facilities {
		if strings.EqualFold(name, facility.DbName) {
			return facility
		}
	}
	return nil
}
func (s *System) Name() string {
	return s.DbName
}

func (s *System) Position() *Coordinate {
	return &s.position
}
func (s *System) String() string {
	return s.DbName
}
