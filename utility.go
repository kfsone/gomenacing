package main

import (
	"bytes"
	gom "github.com/kfsone/gomenacing/pkg/gomschema"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Miscellaneous utility functions.

// ConditionallyOr will or `add` with `base` if predicate is true, without a branch.
// https://gcc.godbolt.org/z/PGbz77
func ConditionallyOrFeatures(base, add FacilityFeatureMask, predicate bool) FacilityFeatureMask {
	var value FacilityFeatureMask
	if predicate {
		value = add
	}
	return base | value
}

// Captures any log output produced by a test function and returns it as a string.
func captureLog(t *testing.T, test func(t *testing.T)) []string {
	// https://stackoverflow.com/questions/44119951/how-to-check-a-log-output-in-go-test
	var buf bytes.Buffer
	log.SetOutput(&buf)
	defer func() {
		log.SetOutput(os.Stderr)
	}()
	test(t)
	if buf.Len() == 0 {
		return nil
	}
	return strings.Split(strings.TrimSpace(buf.String()), "\n")
}

// ensureDirectory creates a directory and its parents if it does not already exist.
func ensureDirectory(path string) (created bool, err error) {
	// Check whether it exists, and if it does check it's a directory.
	info, err := os.Stat(path)
	if err == nil && !info.IsDir() {
		err = os.ErrExist
	}
	if os.IsNotExist(err) {
		if err = os.MkdirAll(path, 0750); err == nil {
			created = true
		}
	}
	return
}

func GetTestDir(path ...string) TestDir {
	pathname := filepath.Join(path...)
	if dirname, err := ioutil.TempDir(pathname, "menace"); err != nil {
		panic(err)
	} else {
		return TestDir(dirname)
	}
}

func (td *TestDir) Close() {
	if err := os.RemoveAll(string(*td)); err != nil {
		panic(err)
	}
}

func stringToFeaturePad(padSize string) FacilityFeatureMask {
	if len(padSize) == 1 {
		switch strings.ToUpper(padSize)[0] {
		case 'L':
			return FeatLargePad
		case 'M':
			return FeatMediumPad
		case 'S':
			return FeatSmallPad
		default:
			break
		}
	}
	return FacilityFeatureMask(0)
}

func FeatureMaskToPadSize(mask FacilityFeatureMask) gom.PadSize {
	if (mask & FeatLargePad) != 0 {
		return gom.PadSize_PadLarge
	}
	if (mask & FeatMediumPad) != 0 {
		return gom.PadSize_PadMedium
	}
	if (mask & FeatSmallPad) != 0 {
		return gom.PadSize_PadSmall
	}
	return gom.PadSize_PadNone
}

func FeatureMaskToServices(mask FacilityFeatureMask, services *gom.Services) {
	// This syntax appear to avoid branching.
	services.HasMarket = (mask & FeatMarket) != 0
	services.HasBlackMarket = (mask & FeatBlackMarket) != 0
	services.HasCommodities = (mask & FeatCommodities) != 0
	services.HasDocking = (mask & FeatDocking) != 0
	services.HasOutfitting = (mask & FeatOutfitting) != 0
	services.IsPlanetary = (mask & FeatPlanetary) != 0
	services.HasRearm = (mask & FeatRearm) != 0
	services.HasRefuel = (mask & FeatRefuel) != 0
	services.HasRepair = (mask & FeatRepair) != 0
	services.HasShipyard = (mask & FeatShipyard) != 0
}

func ServicesToFeatures(services *gom.Services, pad gom.PadSize) (mask FacilityFeatureMask) {
	// Perform boolean operations without branching.
	mask = ConditionallyOrFeatures(mask, FeatMarket, services.HasMarket)
	mask = ConditionallyOrFeatures(mask, FeatBlackMarket, services.HasBlackMarket)
	mask = ConditionallyOrFeatures(mask, FeatCommodities, services.HasCommodities)
	mask = ConditionallyOrFeatures(mask, FeatDocking, services.HasDocking)
	mask = ConditionallyOrFeatures(mask, FeatLargePad, pad == gom.PadSize_PadLarge)
	mask = ConditionallyOrFeatures(mask, FeatMediumPad, pad == gom.PadSize_PadMedium)
	mask = ConditionallyOrFeatures(mask, FeatOutfitting, services.HasOutfitting)
	mask = ConditionallyOrFeatures(mask, FeatPlanetary, services.IsPlanetary)
	mask = ConditionallyOrFeatures(mask, FeatRearm, services.HasRearm)
	mask = ConditionallyOrFeatures(mask, FeatRefuel, services.HasRefuel)
	mask = ConditionallyOrFeatures(mask, FeatRepair, services.HasRepair)
	mask = ConditionallyOrFeatures(mask, FeatShipyard, services.HasShipyard)
	mask = ConditionallyOrFeatures(mask, FeatSmallPad, pad == gom.PadSize_PadSmall)
	return mask
}
