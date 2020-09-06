package main

import (
	"fmt"
	gom "github.com/kfsone/gomenacing/pkg/gomschema"
	"strings"
	"time"
)

type GomDbEntity interface {
	GetId() uint32
	GetName() string
}

func validateEntityForSerialization(kind string, entity GomDbEntity) (string, error) {
	if entity.GetId() == 0 {
		return "", fmt.Errorf("%w %s: without an id", ErrUnknownEntity, kind)
	}
	name := strings.TrimSpace(entity.GetName())
	if name == "" {
		return "", fmt.Errorf("%w %s: #%d: missing name", ErrUnknownEntity, kind, entity.GetId())
	}
	return name, nil
}

//////////////////////////////////////////////////////////////////////////////////////////
// Commodity

// SerializeCommodity converts from a local Commodity into a schema Commodity
func SerializeCommodity(into *gom.Commodity, from *Commodity) {
	into.Id = uint32(from.DbEntity.ID)
	into.Name = from.DbEntity.DbName
	into.CategoryId = from.CategoryID
	if from.IsRare {
		into.IsRare = true
	}
	if from.IsNonMarketable {
		into.IsNonMarketable = true
	}
	into.AverageCr = from.AverageCr
}

// DeserializeCommodity converts from a schema Commodity into a local Commodity
func DeserializeCommodity(into *Commodity, from *gom.Commodity, _ *SystemDatabase) error {
	if name, err := validateEntityForSerialization("commodity", from); err == nil {
		if from.GetCategoryId() != 0 {
			into.DbEntity.ID = EntityID(from.GetId())
			into.DbEntity.DbName = name
			into.CategoryID = from.GetCategoryId()
			into.IsRare = from.GetIsRare()
			into.IsNonMarketable = from.GetIsNonMarketable()
			into.AverageCr = from.GetAverageCr()
			return nil
		}
		return fmt.Errorf("%w: commodity %d: missing category", ErrUnknownEntity, from.GetId())
	} else {
		return err
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
// System

// SerializeSystem converts from a local System into a schema System
func SerializeSystem(into *gom.System, from *System) {
	into.Id = uint32(from.DbEntity.ID)
	into.Name = from.DbEntity.DbName
	into.TimestampUtc = uint64(from.TimestampUtc.Unix())
	into.Position.X = from.Position().X
	into.Position.Y = from.Position().Y
	into.Position.Z = from.Position().Z
	if from.Populated {
		into.Populated = true
	}
	if from.NeedsPermit {
		into.NeedsPermit = true
	}
	into.SecurityLevel = from.SecurityLevel
	into.Government = from.Government
	into.Allegiance = from.Allegiance
}

// DeserializeSystem converts from a schema System into a local System
func DeserializeSystem(into *System, from *gom.System, _ *SystemDatabase) error {
	if name, err := validateEntityForSerialization("system", from); err == nil {
		into.DbEntity.ID = EntityID(from.GetId())
		into.DbEntity.DbName = name
		into.TimestampUtc = time.Unix(int64(from.GetTimestampUtc()), 0)
		into.position = Coordinate{from.Position.X, from.Position.Y, from.Position.Z}
		into.Populated = from.Populated
		into.NeedsPermit = from.NeedsPermit
		into.SecurityLevel = from.SecurityLevel
		into.Government = from.Government
		into.Allegiance = from.Allegiance
		return nil
	} else {
		return err
	}
}

//////////////////////////////////////////////////////////////////////////////////////////
// Facility

// SerializeFacility converts from a local Facility into a schema Facility
func SerializeFacility(into *gom.Facility, from *Facility) error {
	into.Id = uint32(from.DbEntity.ID)
	into.Name = from.DbEntity.DbName
	into.TimestampUtc = uint64(from.TimestampUtc.Unix())
	into.FacilityType = from.FacilityType
	FeatureMaskToServices(from.Features, into.GetServices())
	into.PadSize = FeatureMaskToPadSize(from.Features)
	into.Government = from.Government
	into.Allegiance = from.Allegiance
	return nil
}

// Deserialize converts from a schema System into a local System
func Deserialize(into *Facility, from *gom.Facility, db *SystemDatabase) (err error) {
	var name string
	if name, err = validateEntityForSerialization("facility", from); err != nil {
		return err
	}
	systemId := from.GetSystemId()
	if systemId == 0 {
		return fmt.Errorf("%w system: facility %d (%s) with missing system id", ErrUnknownEntity, from.GetId(), name)
	}
	system := db.GetSystemByID(EntityID(systemId))
	if system == nil {
		return fmt.Errorf("%w system: facility %d (%s) references unknown system: %d", ErrUnknownEntity, from.GetId(), name, systemId)
	}

	into.DbEntity.ID = EntityID(from.GetId())
	into.DbEntity.DbName = name
	into.System = system
	into.TimestampUtc = time.Unix(int64(from.GetTimestampUtc()), 0)
	into.FacilityType = from.FacilityType
	into.Features = ServicesToFeatures(from.GetServices(), from.GetPadSize())
	into.LsFromStar = from.LsFromStar
	into.Government = from.Government
	into.Allegiance = from.Allegiance

	return nil
}
