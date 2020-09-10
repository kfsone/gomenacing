package main

import (
	gom "github.com/kfsone/gomenacing/pkg/gomschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewSystem(t *testing.T) {
	system := NewSystem(DbEntity{ID: 111, DbName: "system of a down"}, Coordinate{X: 1.0, Y: 2.0, Z: 3.0})

	require.NotNil(t, system)
	assert.Equal(t, EntityID(111), system.ID)
	assert.Equal(t, "system of a down", system.DbName)
	assert.Equal(t, Coordinate{1.0, 2.0, 3.0}, system.position)

	assert.Equal(t, uint64(0), system.TimestampUtc)
	assert.False(t, system.Populated)
	assert.False(t, system.NeedsPermit)
	assert.Equal(t, gom.SecurityLevel_SecurityNone, system.SecurityLevel)
	assert.Equal(t, gom.GovernmentType_GovNone, system.Government)
	assert.Equal(t, gom.AllegianceType_AllegNone, system.Allegiance)
	assert.Nil(t, system.facilities)
}

func TestSystem_GetFacility(t *testing.T) {
	system := NewSystem(DbEntity{ID: 100, DbName: "SYSTEM"}, Coordinate{X: 1.0, Y: 2.0, Z: 3.0})
	t.Run("Facility lookup with no facilities", func(t *testing.T) {
		assert.Nil(t, system.GetFacility(""))
		assert.Nil(t, system.GetFacility("facility"))
	})

	t.Run("Facility lookup one facility", func(t *testing.T) {
		facility := &Facility{DbEntity: DbEntity{ID: 200, DbName: "Facility 1"}}
		system.facilities = append(system.facilities, facility)

		assert.Nil(t, system.GetFacility(""))
		assert.Equal(t, facility, system.GetFacility("FacIliTy 1"))
		assert.Nil(t, system.GetFacility("Facility 2"))
	})

	t.Run("Facility lookup multiple facilities", func(t *testing.T) {
		system.facilities = append(system.facilities, &Facility{DbEntity: DbEntity{ID: 111, DbName: "Facile 3"}})
		system.facilities = append(system.facilities, &Facility{DbEntity: DbEntity{ID: 400, DbName: "Das facility Ein"}})
		facility := &Facility{DbEntity: DbEntity{ID: 201, DbName: "facility 2"}}
		system.facilities = append(system.facilities, facility)
		system.facilities = append(system.facilities, &Facility{DbEntity: DbEntity{ID: 202, DbName: "Fin"}})

		assert.Nil(t, system.GetFacility(""))
		assert.NotNil(t, system.GetFacility("facility 1"))
		assert.NotEqual(t, facility, system.GetFacility("facility 1"))
		assert.Equal(t, facility, system.GetFacility("faCIliTY 2"))
		assert.Nil(t, system.GetFacility("Facility 3"))
	})
}

func TestSystem_Name(t *testing.T) {
	system := System{DbEntity: DbEntity{DbName: "System #1"}}
	assert.Equal(t, "System #1", system.Name())

	madeSystem := NewSystem(system.DbEntity, Coordinate{})
	assert.Equal(t, "System #1", madeSystem.Name())
}

func TestSystem_String(t *testing.T) {
	system := System{DbEntity: DbEntity{DbName: "System PQY Z.1+2"}}
	assert.Equal(t, "System PQY Z.1+2", system.String())

	madeSystem := NewSystem(system.DbEntity, Coordinate{})
	assert.Equal(t, "System PQY Z.1+2", madeSystem.String())
}

func TestSystem_Position(t *testing.T) {
	system := System{}
	assert.Equal(t, &system.position, system.Position())
}

func TestSystem_GetDbId(t *testing.T) {
	system := System{}
	assert.Equal(t, "00000000", system.GetDbId())

	system.DbEntity.ID = 0x010ab
	assert.Equal(t, "000010ab", system.GetDbId())

	system.DbEntity.DbName = "SomethingOrOther"
	assert.Equal(t, "000010ab", system.GetDbId())
}

func TestSystem_GetTimestampUtc(t *testing.T) {
	system := System{}
	assert.EqualValues(t, 0, system.GetTimestampUtc())
	system.TimestampUtc = 4999
	assert.EqualValues(t, 4999, system.GetTimestampUtc())
}
