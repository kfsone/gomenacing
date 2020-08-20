package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func jsonImportWrapper(schema *Schema, source string, fields []string, convertFn func([]gjson.Result) (interface{}, error)) error {
	defer func() { failOnError(schema.Close()) }()
	filename := DataFilePath(*EddbPath, source)
	loaded := 0
	var data []byte

	errorsCh, err := ImportJsonlFile(filename, fields, func(jsonLine *JsonEntry) (err error) {
		item, err := convertFn(jsonLine.Results)
		if err == nil {
			if data, err = json.Marshal(item); err == nil {
				if err = schema.Put([]byte(jsonLine.Results[0].Raw), data); err == nil {
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
	schema, err := db.Systems()
	if err == nil {
		err = jsonImportWrapper(schema, EddbSystems, systemFields, func(results []gjson.Result) (interface{}, error) { return NewSystemFromJson(results) })
	}
	return err
}

func importFacilities(db *Database) error {
	schema, err := db.Facilities()
	if err == nil {
		err = jsonImportWrapper(schema, EddbFacilities, facilityFields, func(results []gjson.Result) (interface{}, error) { return NewFacilityFromJson(results) })
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

	schema, err := db.Commodities()
	if err != nil {
		return err
	}
	defer func() { failOnError(schema.Close()) }()

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
			if err = schema.Put([]byte(value.Map()["id"].Raw), data); err == nil {
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

func readCommaSeparatedValues(from io.Reader, rowCh chan<- []string) error {
	defer close(rowCh)
	scanner := bufio.NewScanner(from)
	for scanner.Scan() {
		rowCh <- strings.Split(scanner.Text(), ",")
	}
	return scanner.Err()
}

func getIndexes(fieldNames []string, headers []string) ([]int, error) {
	indexes := make([]int, len(fieldNames))
	for fieldIdx, fieldName := range fieldNames {
		found := false
		for hdrIdx, hdrName := range headers {
			if fieldName == hdrName {
				indexes[fieldIdx] = hdrIdx
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("missing '%s' column", fieldName)
		}
	}
	return indexes, nil
}

func importListings(db *Database) error {
	lineCh := make(chan []string)
	filename := DataFilePath(*EddbPath, EddbListings)
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer func() { failOnError(file.Close()) }()

	schema, err := db.Listings()
	if err != nil {
		return err
	}
	defer func() { failOnError(schema.Close()) }()

	go func() { failOnError(readCommaSeparatedValues(bufio.NewReader(file), lineCh)) }()

	var fieldNames = []string{"station_id", "commodity_id", "supply", "demand", "buy_price", "sell_price", "collected_at"}
	var fieldIndexes []int
	// Map headers <-> expected column positions
	select {
	case headers, ok := <-lineCh:
		if !ok {
			return errors.New("empty listings file")
		}
		fieldIndexes, err = getIndexes(fieldNames, headers)
		if err != nil {
			return fmt.Errorf("%s: %w", EddbListings, err)
		}
	}

	type StationListing struct {
		stationId   EntityID
		commodityId EntityID
		listing     string
	}

	listingCh := make(chan *StationListing, 4)
	go func() {
		defer close(listingCh)
		for line := range lineCh {
			// Get the stationID
			row := make([]string, len(fieldIndexes))
			for fieldIdx, rowIdx := range fieldIndexes {
				row[fieldIdx] = line[rowIdx]
			}
			stationId, err := strconv.ParseInt(row[0], 10, 64)
			if err != nil {
				if err = FilterError(err); err != nil {
					FilterError(fmt.Errorf("bad station id: %w", err))
				}
				continue
			}
			if stationId <= 0 || stationId >= 1<<32 {
				FilterError(fmt.Errorf("%w: invalid station id: %s", ErrUnknownEntity, row[0]))
				continue
			}
			commodityId, err := strconv.ParseInt(row[1], 10, 64)
			if err != nil {
				FilterError(err)
				continue
			}
			if commodityId <= 0 || commodityId >= 1<<32 {
				FilterError(fmt.Errorf("%w: invalid commodity id: %s", ErrUnknownEntity, row[1]))
				continue
			}
			listingCh <- &StationListing{EntityID(stationId), EntityID(commodityId), strings.Join(row, ",")}
		}
	}()

	loaded := 0
	stationData := make(map[EntityID]map[EntityID]string, 80000)
	for listing := range listingCh {
		listings, exists := stationData[listing.stationId]
		if exists {
			listings[listing.commodityId] = listing.listing
		} else {
			stationData[listing.stationId] = make(map[EntityID]string, 100)
			stationData[listing.stationId][listing.commodityId] = listing.listing
		}
		loaded += 1
	}
	// Consolidate and write
	fmt.Printf("%s: imported %d items for %d stations\n", EddbListings, loaded, len(stationData))
	go func() {
		file, err := os.Create("listings.json")
		if err != nil {
			panic(err)
		}
		data, err := json.Marshal(stationData)
		if err != nil {
			panic(err)
		}
		_, err = file.Write(data)
		failOnError(err)
		failOnError(file.Close())
	}()
	for stationId, listingMap := range stationData {
		values := ""
		for _, value := range listingMap {
			values += value
			values += " "
		}
		key := fmt.Sprintf("%d", stationId)
		err := schema.Put([]byte(key), []byte(values))
		fmt.Fprintf(file, "key:[%s], value:[%s]\n", key, values)
		if err != nil {
			panic(err)
		}
	}
	fmt.Printf("%s: wrote %d stations\n", EddbListings, len(stationData))

	return nil
}
