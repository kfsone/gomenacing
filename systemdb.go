package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/kfsone/gomenacing/pkg/gomschema"
	"google.golang.org/protobuf/proto"
	"log"
	"strings"

	flag "github.com/spf13/pflag"
)

var eddbPath = flag.StringP("eddbdir", "e", "", "Path to EDDB json files to import.")

const (
	EddbSystems     string = "systems_populated.jsonl"
	EddbFacilities  string = "stations.jsonl"
	EddbCommodities string = "commodities.json"
)

type SystemDatabase struct {
	db *Database
	// Index of Systems by their database ids.
	systemsByID map[EntityID]*System
	// Look-up a system's EntityID by it's name.
	systemIDs map[string]EntityID
	// Index of Facilities by their database ids.
	facilitiesByID map[EntityID]*Facility
	// Index of Commodities by their database ids.
	commoditiesByID map[EntityID]*Commodity
	// Look-up a commodity's EntityID by it's name.
	commodityIDs map[string]EntityID
	// Localized index of systems based on their sector keys.
	sectors map[SectorKey][]*System
}

func NewSystemDatabase(db *Database) *SystemDatabase {
	return &SystemDatabase{
		db:              db,
		systemsByID:     make(map[EntityID]*System, 4096),
		systemIDs:       make(map[string]EntityID, 4096),
		facilitiesByID:  make(map[EntityID]*Facility, 8192),
		commoditiesByID: make(map[EntityID]*Commodity, 500),
		commodityIDs:    make(map[string]EntityID, 500),
		sectors:         make(map[SectorKey][]*System, 1024),
	}
}

func registerIDLookup(entity *DbEntity, ids map[string]EntityID) bool {
	name := strings.ToLower(entity.DbName)
	if _, present := ids[name]; present != false {
		return false
	}
	ids[name] = entity.ID
	return true
}

func (sdb *SystemDatabase) registerCommodity(commodity *Commodity) (err error) {
	if _, present := sdb.commoditiesByID[commodity.ID]; present == false {
		if registerIDLookup(&commodity.DbEntity, sdb.commodityIDs) {
			sdb.commoditiesByID[commodity.ID] = commodity
			return nil
		}
		err = fmt.Errorf("%w: item name", ErrDuplicateEntity)
	} else {
		err = fmt.Errorf("%w: item id", ErrDuplicateEntity)
	}
	return fmt.Errorf("%s (#%d): %w", commodity.DbName, commodity.ID, err)
}

func (sdb *SystemDatabase) registerSystem(system *System) (err error) {
	if _, present := sdb.systemsByID[system.ID]; present == false {
		if registerIDLookup(&system.DbEntity, sdb.systemIDs) {
			sdb.systemsByID[system.ID] = system
			return nil
		}
		err = fmt.Errorf("%w: system name", ErrDuplicateEntity)
	} else {
		err = fmt.Errorf("%w: system id", ErrDuplicateEntity)
	}

	return fmt.Errorf("%s (#%d): %w", system.DbName, system.ID, err)
}

func (sdb *SystemDatabase) registerSystemToSector(system *System) {
	key := system.Position().SectorKey()
	sector, exists := sdb.sectors[key]
	if !exists {
		sector = make([]*System, 0, 8)
	}
	// Find an existing match.
	for idx, existing := range sector {
		if existing.GetId() == system.GetId() {
			sector[idx] = system
			return
		}
	}
	// Just add to the end.
	sdb.sectors[key] = append(sector, system)
}

func (sdb *SystemDatabase) registerFacility(facility *Facility) error {
	var exists bool
	system := facility.System
	if system == nil {
		return fmt.Errorf("%s (#%d): attempted to register facility without a system", facility.DbName, facility.ID)
	}
	if _, exists = sdb.facilitiesByID[facility.ID]; exists != false {
		return fmt.Errorf("%s/%s (#%d): %w: facility id", system.DbName, facility.DbName, facility.ID, ErrDuplicateEntity)
	}

	for _, existing := range system.facilities {
		if strings.EqualFold(existing.DbName, facility.DbName) {
			return fmt.Errorf("%s/%s (#%d): %w: facility name in system", system.DbName, facility.DbName, facility.ID, ErrDuplicateEntity)
		}
	}

	system.facilities = append(system.facilities, facility)
	sdb.facilitiesByID[facility.ID] = facility

	return nil
}

func (sdb *SystemDatabase) registerFromMessage(message proto.Message, schema *Schema) error {
	switch typed := message.(type) {
	case *gomschema.Commodity:
		return sdb.updateCommodity(typed, schema)

	case *gomschema.System:
		return sdb.updateSystem(typed, schema)

	case *gomschema.Facility:
		return sdb.updateFacility(typed, schema)

	case *gomschema.FacilityListing:
		return sdb.updateFacilityListing(typed, schema)

	default:
		panic("Unknown message type")
	}
}

func (sdb *SystemDatabase) GetSystemByID(id EntityID) (system *System) {
	if system, exists := sdb.systemsByID[id]; exists {
		return system
	}
	return nil
}

func (sdb *SystemDatabase) GetSystem(name string) (system *System) {
	if id, exists := sdb.systemIDs[strings.ToLower(name)]; exists {
		system = sdb.systemsByID[id]
	}
	return
}

func (sdb *SystemDatabase) GetFacilityByID(id EntityID) *Facility {
	if facility, exists := sdb.facilitiesByID[id]; exists {
		return facility
	}
	return nil
}

func writeMessageForId(message proto.Message, schema *Schema) error {
	type Identifiable interface {
		GetId() uint32
	}

	key := make([]byte, 4)
	binary.LittleEndian.PutUint32(key, message.(Identifiable).GetId())

	value, err := proto.Marshal(message)
	if err != nil {
		return err
	}
	return schema.Put(key, value)
}

func (sdb *SystemDatabase) newCommodity(gomItem *gomschema.Commodity) error {
	entity, err := NewDbEntity(int64(gomItem.Id), gomItem.Name)
	if err == nil {
		item, err := NewCommodity(entity, gomItem.CategoryId, gomItem.IsRare, gomItem.IsNonMarketable, gomItem.AverageCr)
		if err == nil {
			err = sdb.registerCommodity(item)
		}
	}
	return err
}

func (sdb *SystemDatabase) newSystem(gomItem *gomschema.System) error {
	entity, err := NewDbEntity(int64(gomItem.Id), gomItem.Name)
	if err != nil {
		return err
	}

	item := NewSystem(entity, Coordinate{gomItem.Position.X, gomItem.Position.Y, gomItem.Position.Z})
	item.TimestampUtc = getTimestamp(gomItem)
	item.Populated = gomItem.GetPopulated()
	item.NeedsPermit = gomItem.GetNeedsPermit()
	item.SecurityLevel = gomItem.GetSecurityLevel()
	item.Government = gomItem.GetGovernment()
	item.Allegiance = gomItem.GetAllegiance()

	err = sdb.registerSystem(item)
	if err == nil {
		sdb.registerSystemToSector(item)
	}

	return err
}

func (sdb *SystemDatabase) newFacility(gomItem *gomschema.Facility) error {
	entity, err := NewDbEntity(int64(gomItem.Id), gomItem.Name)
	if err != nil {
		return err
	}
	system := sdb.GetSystemByID(EntityID(gomItem.SystemId))
	if system == nil {
		return fmt.Errorf("facility %s: %d: %w: sytem id: %d", gomItem.Name, gomItem.Id, ErrUnknownEntity, gomItem.SystemId)
	}
	item, err := NewFacility(entity, system, gomItem.FacilityType, FacilityFeatureMask(gomItem.Features))
	if err == nil {
		item.TimestampUtc = getTimestamp(gomItem)
		item.LsFromStar = gomItem.GetLsFromStar()
		item.Government = gomItem.GetGovernment()
		item.Allegiance = gomItem.GetAllegiance()
		err = sdb.registerFacility(item)
	}
	return err
}

func (sdb *SystemDatabase) newListings(gomItem *gomschema.FacilityListing) error {
	facility := sdb.GetFacilityByID(EntityID(gomItem.GetId()))
	if facility == nil {
		return fmt.Errorf("%w: facility id: %d", ErrUnknownEntity, gomItem.GetId())
	}
	listings := gomItem.GetListings()
	if len(listings) == 0 {
		facility.listings = nil
		return nil
	}
	if facility.listings == nil {
		facility.listings = make(map[EntityID]*Listing, len(listings))
	}
	for _, gomListing := range listings {
		l := Listing{
			CommodityID:  EntityID(gomListing.CommodityId),
			Supply:       gomListing.GetSupplyUnits(),
			StationPays:  gomListing.GetSupplyCredits(),
			Demand:       gomListing.GetDemandUnits(),
			StationAsks:  gomListing.GetDemandCredits(),
			TimestampUtc: getTimestamp(gomListing),
		}
		facility.listings[l.CommodityID] = &l
	}
	return nil
}

func (sdb *SystemDatabase) updateCommodity(item *gomschema.Commodity, schema *Schema) error {
	name := strings.ToLower(item.Name)
	if existing, exists := sdb.commodityIDs[name]; exists {
		if existing != EntityID(item.Id) {
			return fmt.Errorf("commodity %s: %d: name collides with #%d", item.Name, item.Id, existing)
		}
		commodity := sdb.commoditiesByID[existing]
		commodity.DbEntity.DbName = item.Name
		commodity.CategoryID = item.GetCategoryId()
		commodity.IsRare = item.IsRare
		commodity.IsNonMarketable = item.IsNonMarketable
		commodity.AverageCr = item.AverageCr
	} else {
		if err := sdb.newCommodity(item); err != nil {
			return err
		}
	}
	return writeMessageForId(item, schema)
}

func requireNewer(newer, older Timestamped) error {
	if older.GetTimestampUtc() != 0 && newer.GetTimestampUtc() <= older.GetTimestampUtc() {
		return fmt.Errorf("stale update (%v v %v)", newer.GetTimestampUtc(), older.GetTimestampUtc())
	}
	return nil
}

func (sdb *SystemDatabase) updateSystem(item *gomschema.System, schema *Schema) error {
	name := strings.ToLower(item.Name)
	if existing, exists := sdb.systemIDs[name]; exists {
		if existing != EntityID(item.Id) {
			return fmt.Errorf("system %s (%d): name collides with #%d", item.Name, item.Id, existing)
		}
		// Is this an update?
		system := sdb.systemsByID[existing]
		if err := requireNewer(item, system); err != nil {
			log.Printf("%s (%d): %s", item.Name, item.Id, err)
			return nil
		}
		system.DbEntity.DbName = item.Name
		system.TimestampUtc = getTimestamp(item)
		system.position = Coordinate{item.Position.X, item.Position.Y, item.Position.Z }
		system.Populated = item.Populated
		system.NeedsPermit = item.NeedsPermit
		system.SecurityLevel = item.SecurityLevel
		system.Government = item.Government
		system.Allegiance = item.Allegiance
	} else {
		if err := sdb.newSystem(item); err != nil {
			return err
		}
	}
	return writeMessageForId(item, schema)
}

func updateExistingFacility(sdb *SystemDatabase, newSystem *System, oldFacility *Facility, item *gomschema.Facility) error {
	if err := requireNewer(item, oldFacility); err != nil {
		log.Printf("%s (%d): %s", oldFacility.Name(), item.Id, err)
		return err
	}
	if oldFacility.System != newSystem || strings.ToLower(oldFacility.GetName()) != strings.ToLower(item.Name) {
		return errors.New("can't handle facility renames/relocates")
	}

	oldFacility.DbName = item.Name
	oldFacility.System = newSystem
	oldFacility.FacilityType = item.FacilityType
	oldFacility.Features = FacilityFeatureMask(item.Features)
	oldFacility.TimestampUtc = getTimestamp(item)
	oldFacility.LsFromStar = item.GetLsFromStar()
	oldFacility.Government = item.GetGovernment()
	oldFacility.Allegiance = item.GetAllegiance()

	return nil
}

func (sdb *SystemDatabase) updateFacility(item *gomschema.Facility, schema *Schema) (err error) {
	// Does the destination system exist?
	system := sdb.GetSystemByID(EntityID(item.SystemId))
	if system == nil {
		return fmt.Errorf("%w: facility %s (%d): no such system %d", ErrUnknownEntity, item.Name, item.Id, item.SystemId)
	}

	// Does the facility already exist?
	if oldFacility, exists := sdb.facilitiesByID[EntityID(item.Id)]; exists {
		err = updateExistingFacility(sdb, system, oldFacility, item)
	} else {
		err = sdb.newFacility(item)
	}
	if err != nil {
		return err
	}
	return writeMessageForId(item, schema)
}

func (sdb *SystemDatabase) updateFacilityListing(item *gomschema.FacilityListing, schema *Schema) error {
	panic("Not handled")
}
