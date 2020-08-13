package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestNewSystem(t *testing.T) {
	t.Run("Reject bad parameters", func(t *testing.T) {
		system, err := NewSystem(0, "x", Coordinate{}, false)
		assert.Nil(t, system)
		if assert.Error(t, err) {
			assert.Equal(t, "invalid system id: 0", err.Error())
		}

		system, err = NewSystem(1<<32, "x", Coordinate{}, false)
		assert.Nil(t, system)
		if assert.Error(t, err) {
			assert.Equal(t, fmt.Errorf("invalid system id (too large): %d", 1<<32), err)
		}

		system, err = NewSystem(1, "", Coordinate{}, false)
		assert.Nil(t, system)
		if assert.Error(t, err) {
			assert.Equal(t, "empty system name", err.Error())
		}

		system, err = NewSystem((1<<32)-1, "  ", Coordinate{}, false)
		assert.Nil(t, system)
		if assert.Error(t, err) {
			assert.Equal(t, "empty system name", err.Error())
		}
	})

	system, err := NewSystem(111, "system of a down", Coordinate{1.0, 2.0, 3.0}, true)
	require.Nil(t, err)
	assert.Equal(t, EntityID(111), system.Id)
	assert.Equal(t, "SYSTEM OF A DOWN", system.DbName)
	assert.Equal(t, 1.0, system.Position.X)
	assert.Equal(t, 2.0, system.Position.Y)
	assert.Equal(t, 3.0, system.Position.Z)
	assert.True(t, system.Permit)
	assert.Empty(t, system.Facilities)
}

func TestSystem_Distance(t *testing.T) {
	first, err := NewSystem(123, "first", Coordinate{1.0, -2.0, 5.0}, false)
	require.Nil(t, err)
	assert.Equal(t, Square(0), first.Distance(first))

	second, err := NewSystem(201, "second", Coordinate{0, 3.3, 9.9}, false)
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
	system, err := NewSystem(100, "system", Coordinate{1.0, 2.0, 3.0}, false)
	require.Nil(t, err)
	t.Run("Facility lookup with no facilities", func(t *testing.T) {
		assert.Nil(t, system.GetFacility(""))
		assert.Nil(t, system.GetFacility("facility"))
	})

	t.Run("Facility lookup one facility", func(t *testing.T) {
		facility, err := system.NewFacility(200, "facility 1", FacilityFeatureMask(0))
		assert.NotNil(t, facility)
		assert.Nil(t, err)
		system.Facilities = append(system.Facilities, facility)

		assert.Nil(t, system.GetFacility(""))
		assert.Equal(t, facility, system.GetFacility("FacIliTy 1"))
		assert.Nil(t, system.GetFacility("Facility 2"))
	})

	t.Run("Facility lookup two facilities", func(t *testing.T) {
		facility, err := system.NewFacility(201, "facility 2", FacilityFeatureMask(0))
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
	system, err := NewSystem(100, "system1", Coordinate{}, false)
	require.Nil(t, err)
	assert.Equal(t, "SYSTEM1", system.Name(-1))
	assert.Equal(t, "SYSTEM1", system.Name(0))
	assert.Equal(t, "SYSTEM1", system.Name(1))
	assert.Equal(t, "SYSTEM1", system.Name(2))
	assert.Equal(t, "SYSTEM1", system.Name(999))
}

func TestSystem_NewFacility(t *testing.T) {
	system, err := NewSystem(100, "system", Coordinate{}, false)
	require.Nil(t, err)
	assert.Empty(t, system.Facilities)

	t.Run("Reject bad values", func(t *testing.T) {
		facility, err := system.NewFacility(0, "", 0)
		assert.Nil(t, facility)
		if assert.Error(t, err) {
			assert.Equal(t, "invalid facility id: 0", err.Error())
		}

		facility, err = system.NewFacility(1<<32, "", 0)
		assert.Nil(t, facility)
		if assert.Error(t, err) {
			assert.Equal(t, fmt.Errorf("invalid facility id: %d", 1<<32), err)
		}

		facility, err = system.NewFacility(1<<32-1, "", 0)
		assert.Nil(t, facility)
		if assert.Error(t, err) {
			assert.Equal(t, "invalid (empty) facility name", err.Error())
		}

		facility, err = system.NewFacility(1, " \t ", 0)
		assert.Nil(t, facility)
		if assert.Error(t, err) {
			assert.Equal(t, "invalid (empty) facility name", err.Error())
		}
	})

	t.Run("Create genuine facility", func(t *testing.T) {
		features := FeatBlackMarket | FeatSmallPad
		facility, err := system.NewFacility(111, "first", features)
		assert.Nil(t, err)
		assert.NotNil(t, facility)
		assert.Equal(t, EntityID(111), facility.Id)
		assert.Equal(t, "FIRST", facility.DbName)
		assert.Equal(t, system, facility.System)
		assert.Equal(t, features, facility.Features)

		assert.Empty(t, system.Facilities)
	})
}

func TestSystem_String(t *testing.T) {
	system, err := NewSystem(100, "system1", Coordinate{}, false)
	require.Nil(t, err)
	assert.Equal(t, "SYSTEM1", system.String())
	assert.Equal(t, "SYSTEM1", fmt.Sprintf("%s", system))
}
