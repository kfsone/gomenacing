package main

import (
	"fmt"
	"testing"

	"github.com/tidwall/gjson"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestNewSystem(t *testing.T) {
	system, err := NewSystem(DbEntity{111, "system of a down"}, Coordinate{1.0, 2.0, 3.0}, true)
	require.Nil(t, err)
	assert.Equal(t, EntityID(111), system.ID)
	assert.Equal(t, "system of a down", system.DbName)
	assert.Equal(t, 1.0, system.Position.X)
	assert.Equal(t, 2.0, system.Position.Y)
	assert.Equal(t, 3.0, system.Position.Z)
	assert.True(t, system.Permit)
	assert.Empty(t, system.Facilities)
}

func TestSystem_NewSystemFromJson(t *testing.T) {
	json := "{\"id\":1,\"name\":\"sol\",\"x\":0,\"y\":0.0,\"z\":0.00,\"needs_permit\":true}"
	results := gjson.GetMany(json, systemFields...)
	require.Len(t, systemFields, len(results))

	system, err := NewSystemFromJson(results)
	require.Nil(t, err)
	require.NotNil(t, system)

	assert.Equal(t, EntityID(1), system.ID)
	assert.Equal(t, "sol", system.DbName)
	assert.Equal(t, Coordinate{0, 0, 0}, system.Position)
	assert.Equal(t, true, system.Permit)
}

func TestSystem_Distance(t *testing.T) {
	first, err := NewSystem(DbEntity{123, "first"}, Coordinate{1.0, -2.0, 5.0}, false)
	require.Nil(t, err)
	assert.Equal(t, Square(0), first.Distance(first))

	second, err := NewSystem(DbEntity{201, "second"}, Coordinate{0, 3.3, 9.9}, false)
	require.Nil(t, err)
	assert.Equal(t, Square(0), second.Distance(second))

	xDelta := first.Position.X - second.Position.X
	yDelta := first.Position.Y - second.Position.Y
	zDelta := first.Position.Z - second.Position.Z
	distanceSq := Square(xDelta*xDelta + yDelta*yDelta + zDelta*zDelta)
	assert.Equal(t, distanceSq, first.Distance(second))
	assert.Equal(t, distanceSq, second.Distance(first))
}

func TestSystem_GetFacility(t *testing.T) {
	system, err := NewSystem(DbEntity{100, "SYSTEM"}, Coordinate{1.0, 2.0, 3.0}, false)
	require.Nil(t, err)
	t.Run("Facility lookup with no facilities", func(t *testing.T) {
		assert.Nil(t, system.GetFacility(""))
		assert.Nil(t, system.GetFacility("facility"))
	})

	t.Run("Facility lookup one facility", func(t *testing.T) {
		facility, err := NewFacility(DbEntity{200, "Facility 1"}, system, FacilityFeatureMask(0))
		assert.NotNil(t, facility)
		assert.Nil(t, err)
		system.Facilities = append(system.Facilities, facility)

		assert.Nil(t, system.GetFacility(""))
		assert.Equal(t, facility, system.GetFacility("FacIliTy 1"))
		assert.Nil(t, system.GetFacility("Facility 2"))
	})

	t.Run("Facility lookup two facilities", func(t *testing.T) {
		facility, err := NewFacility(DbEntity{201, "facility 2"}, system, FacilityFeatureMask(0))
		assert.NotNil(t, facility)
		assert.Nil(t, err)
		system.Facilities = append(system.Facilities, facility)

		assert.Nil(t, system.GetFacility(""))
		assert.NotNil(t, system.GetFacility("facility 1"))
		assert.NotEqual(t, facility, system.GetFacility("facility 1"))
		assert.Equal(t, facility, system.GetFacility("faCIliTY 2"))
		assert.Nil(t, system.GetFacility("Facility 3"))
	})
}

func TestSystem_Name(t *testing.T) {
	system, err := NewSystem(DbEntity{100, "System1"}, Coordinate{}, false)
	require.Nil(t, err)
	assert.Equal(t, "System1", system.Name())
}

func TestSystem_String(t *testing.T) {
	system, err := NewSystem(DbEntity{100, "System1"}, Coordinate{}, false)
	require.Nil(t, err)
	assert.Equal(t, "System1", system.String())
	assert.Equal(t, "System1", fmt.Sprintf("%s", system))
}
