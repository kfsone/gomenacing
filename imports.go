package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/akrylysov/pogreb"
	"github.com/tidwall/gjson"
)

func importWrapper(store *pogreb.DB, source string, fields []string, convertFn func([]gjson.Result) (interface{}, error)) error {
	defer func() { failOnError(store.Close()) }()
	filename := DataFilePath(*EddbPath, source)
	loaded := 0
	var data []byte

	errorsCh, err := ImportJsonlFile(filename, fields, func(jsonLine *JsonEntry) (err error) {
		item, err := convertFn(jsonLine.Results)
		if err == nil {
			if data, err = json.Marshal(item); err == nil {
				if err = store.Put([]byte(jsonLine.Results[0].Raw), data); err == nil {
					loaded += 1
					return nil
				}
			}
		}
		return fmt.Errorf("%s:%d: %w", source, jsonLine.LineNo, err)
	})
	if err == nil {
		err = countErrors(errorsCh)
	}
	if err != nil {
		return err
	}

	fmt.Printf("%s: imported %d items\n", source, loaded)
	return nil
}

func importSystems(db *Database) error {
	store, err := db.Systems()
	if err == nil {
		err = importWrapper(store, EddbSystems, systemFields, func(results []gjson.Result) (interface{}, error) { return NewSystemFromJson(results) })
	}
	return err
}

func importFacilities(db *Database) error {
	store, err := db.Facilities()
	if err == nil {
		err = importWrapper(store, EddbFacilities, facilityFields, func(results []gjson.Result) (interface{}, error) { return NewFacilityFromJson(results) })
	}
	return err
}

func importCommodities(db *Database) error {
	// Read the json.
	filename := DataFilePath(*EddbPath, EddbCommodities)
	text, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	result := gjson.ParseBytes(text)
	// Check what we have is an array
	if !result.IsArray() || len(result.Array()) <= 0 {
		log.Println("No commodities to import")
		return nil
	}

	store, err := db.Commodities()
	if err != nil {
		return err
	}
	defer func() { failOnError(store.Close()) }()

	var hadError bool
	var loaded int64
	result.ForEach(func(_, value gjson.Result) bool {
		item, err := NewCommodityFromJsonMap(value)
		if err != nil {
			if FilterError(err) != nil {
				log.Print(err)
				hadError = true
			}
			return true
		}
		if data, err := json.Marshal(item); err == nil {
			if err = store.Put([]byte(value.Map()["id"].Raw), data); err == nil {
				loaded += 1
				return true
			}
		}
		return false
	})

	if hadError {
		return errors.New("errors importing commodities")
	}

	fmt.Printf("%s: imported %d items\n", EddbCommodities, loaded)
	return nil
}
