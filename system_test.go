package main

import (
	gom "github.com/kfsone/gomenacing/pkg/gomschema"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestNewSystem(t *testing.T) {
	now := time.Now()
	system := NewSystem(
		DbEntity{ID: 111, DbName: "system of a down"},
		now,
		Coordinate{X: 1.0, Y: 2.0, Z: 3.0},
		true,
		false,
		gom.SecurityLevel_SecurityHigh,
		gom.GovernmentType_GovFeudal,
		gom.AllegianceType_AllegEmpire)

	require.NotNil(t, system)
	assert.Equal(t, EntityID(111), system.ID)
	assert.Equal(t, "system of a down", system.DbName)
	assert.Equal(t, now, system.TimestampUtc)
	assert.Equal(t, 1.0, system.Position.X)
	assert.Equal(t, 2.0, system.Position.Y)
	assert.Equal(t, 3.0, system.Position.Z)
	assert.True(t, system.Populated)
	assert.False(t, system.NeedsPermit)
	assert.Equal(t, gom.SecurityLevel_SecurityHigh, system.SecurityLevel)
	assert.Equal(t, gom.GovernmentType_GovFeudal, system.Government)
	assert.Equal(t, gom.AllegianceType_AllegEmpire, system.Allegiance)
	assert.Nil(t, system.facilities)
}

func TestSystem_Distance(t *testing.T) {
	first := System{Position: Coordinate{X: 1.0, Y: -2.0, Z: 5.0}}
	assert.Equal(t, Square(0), first.Distance(&first))

	second := System{Position: Coordinate{Y: 3.3, Z: 9.9}}
	assert.Equal(t, Square(0), second.Distance(&second))

	xDelta := first.Position.X - second.Position.X
	yDelta := first.Position.Y - second.Position.Y
	zDelta := first.Position.Z - second.Position.Z
	distanceSq := Square(xDelta*xDelta + yDelta*yDelta + zDelta*zDelta)
	assert.Equal(t, distanceSq, first.Distance(&second))
	assert.Equal(t, distanceSq, second.Distance(&first))
}

//func TestSystem_GetFacility(t *testing.T) {
//	system, err := NewSystem(DbEntity{ID: 100, DbName: "SYSTEM"}, Coordinate{X: 1.0, Y: 2.0, Z: 3.0}, false)
//	require.Nil(t, err)
//	t.Run("Facility lookup with no facilities", func(t *testing.T) {
//		assert.Nil(t, system.GetFacility(""))
//		assert.Nil(t, system.GetFacility("facility"))
//	})
//
//	t.Run("Facility lookup one facility", func(t *testing.T) {
//		facility, err := NewFacility(DbEntity{ID: 200, DbName: "Facility 1"}, system, FacilityFeatureMask(0))
//		assert.NotNil(t, facility)
//		assert.Nil(t, err)
//		system.Facilities = append(system.Facilities, facility)
//
//		assert.Nil(t, system.GetFacility(""))
//		assert.Equal(t, facility, system.GetFacility("FacIliTy 1"))
//		assert.Nil(t, system.GetFacility("Facility 2"))
//	})
//
//	t.Run("Facility lookup two facilities", func(t *testing.T) {
//		facility, err := NewFacility(DbEntity{ID: 201, DbName: "facility 2"}, system, FacilityFeatureMask(0))
//		assert.NotNil(t, facility)
//		assert.Nil(t, err)
//		system.Facilities = append(system.Facilities, facility)
//
//		assert.Nil(t, system.GetFacility(""))
//		assert.NotNil(t, system.GetFacility("facility 1"))
//		assert.NotEqual(t, facility, system.GetFacility("facility 1"))
//		assert.Equal(t, facility, system.GetFacility("faCIliTY 2"))
//		assert.Nil(t, system.GetFacility("Facility 3"))
//	})
//}

func TestSystem_Name(t *testing.T) {
	system := System{DbEntity: DbEntity{DbName: "System #1"}}
	assert.Equal(t, "System #1", system.Name())

	madeSystem := NewSystem(system.DbEntity, time.Now(), Coordinate{}, true, true, gom.SecurityLevel_SecurityLow, gom.GovernmentType_GovPatronage, gom.AllegianceType_AllegPilotsFederation)
	assert.Equal(t, "System #1", madeSystem.Name())
}

func TestSystem_String(t *testing.T) {
	system := System{DbEntity: DbEntity{DbName: "System PQY Z.1+2"}}
	assert.Equal(t, "System PQY Z.1+2", system.String())

	madeSystem := NewSystem(system.DbEntity, time.Now(), Coordinate{}, true, true, gom.SecurityLevel_SecurityLow, gom.GovernmentType_GovPatronage, gom.AllegianceType_AllegPilotsFederation)
	assert.Equal(t, "System PQY Z.1+2", madeSystem.String())
}
