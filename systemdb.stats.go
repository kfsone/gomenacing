package main

import (
	"fmt"
	"github.com/kfsone/gomenacing/pkg/gomschema"
	"io"
	"math"
	"sort"
	"strings"
)

func sum(list []int) (total int64, average float64) {
	for _, value := range list {
		total += int64(value)
	}
	if len(list) > 0 {
		average = float64(total) / float64(len(list))
	}
	return
}

func percentile(ptile float64, list []int) int {
	if len(list) > 1 {
		// A percentile is a value for which p% of the samples in list fall at or below.
		point := float64(len(list)) * ptile
		index := int(point)
		// When percentile point is not a whole value, we calculate a value midway between
		// the index below and above the point.
		if math.Ceil(point) != math.Floor(point) {
			return list[index] + (list[index+1]-list[index])/2
		} else {
			return list[index]
		}
	} else if len(list) == 1 {
		return list[0]
	}
	return 0
}

func percentage(value, of int) (result float64) {
	if of > 0 {
		result = float64(value) * 100. / float64(of)
	}
	return
}

func average(value, of int) (result float64) {
	if of > 0 {
		return float64(value) / float64(of)
	}
	return 0.
}

func getBounds(sdb *SystemDatabase) (x1, x2, y1, y2, z1, z2 int) {
	if len(sdb.sectors) == 0 {
		return 0, 0, 0, 0, 0, 0
	}
	xs := make([]int, 0, len(sdb.sectors))
	ys := make([]int, 0, len(sdb.sectors))
	zs := make([]int, 0, len(sdb.sectors))
	for key, _ := range sdb.sectors {
		xs = append(xs, key.X)
		ys = append(ys, key.Y)
		zs = append(zs, key.Z)
	}
	sort.Ints(xs)
	sort.Ints(ys)
	sort.Ints(zs)

	return xs[0], xs[len(xs)-1], ys[0], ys[len(ys)-1], zs[0], zs[len(zs)-1]
}

func produceStats (source map[string]int, total int) []string {
	type Stat struct {
		name       string
		count      int
		percentage float64
	}

	stats := make([]Stat, 0, len(source))
	for key, count := range source {
		stats = append(stats, Stat{
			key, count, percentage(count, total),
		})
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].count > stats[j].count
	})

	result := make([]string, len(source))
	for idx, stat := range stats {
		result[idx] = fmt.Sprintf("%s: %d (%.2f%%)", stat.name, stat.count, stat.percentage)
	}

	return result
}

func analyzeSystems(o io.Writer, label string, systems map[EntityID]*System, predicate func(*System) int) {
	list := make([]int, 0, len(systems))
	countDist := make(map[int]int)
	for _, system := range systems {
		if count := predicate(system); count > 0 {
			list = append(list, count)
			countDist[count]++
		}
	}

	if len(list) == 0 {
		return
	}
	if len(list) == 1 {
		fmt.Fprintf(o, "- %s: %d\n", label, list[0])
		return
	}

	median := 0
	for value, count := range countDist {
		if count > countDist[median] {
			median = value
		}
	}

	sort.Ints(list)
	min, max := list[0], list[len(list)-1]
	total, avg := sum(list)
	p95 := percentile(.95, list)
	fmt.Fprintf(o, "- Min/Med/Avg/P95/Max of %d %s: %d/%d/%.2f/%d/%d\n", total, label, min, median, avg, p95, max)
}

func countSystems(sdb *SystemDatabase, predicate func(*System) bool) (count int) {
	for _, system := range sdb.systemsByID {
		var val int
		if predicate(system) {
			val = 1
		}
		count += val
	}
	return
}

func reportOnCommodities(o io.Writer, sdb *SystemDatabase) {
	fmt.Fprintf(o, "Commodities: %d\n", len(sdb.commoditiesByID))
	if len(sdb.commoditiesByID) > 0 {
		min, max := ^EntityID(0), EntityID(0)
		for id := range sdb.commoditiesByID {
			if id < min {
				min = id
			}
			if id > max {
				max = id
			}
		}
		fmt.Printf("- IDs: %d..%d\n", min, max)
	}
}

func reportOnSystems(o io.Writer, sdb *SystemDatabase) {
	fmt.Fprintf(o, "Systems: %d\n", len(sdb.systemsByID))
	total := len(sdb.systemsByID)
	if total == 0 {
		return
	}

	populated, permits := 0, 0
	govtDistrib := make(map[string]int, 32)
	allegDistrib := make(map[string]int, 32)
	securityDistrib := make(map[string]int, 32)
	for _, system := range sdb.systemsByID {
		if system.Populated {
			populated++
		}
		if system.NeedsPermit {
			permits++
		}
		govtDistrib[gomschema.GovernmentType_name[int32(system.Government)]]++
		allegDistrib[gomschema.AllegianceType_name[int32(system.Allegiance)]]++
		securityDistrib[gomschema.FacilityType_name[int32(system.SecurityLevel)]]++
	}

	analyzeSystems(o, "w/facilities", sdb.systemsByID, func(s *System) int { return len(s.facilities) })
	fmt.Fprintf(o, "- Populated: %d (%.2f%%)\n", populated, percentage(populated, len(sdb.systemsByID)))
	fmt.Fprintf(o, "- Need permit: %d (%.2f%%)\n", permits, percentage(permits, len(sdb.systemsByID)))
	stats := produceStats(govtDistrib, total)
	fmt.Fprintf(o, "- Governments: %s\n", strings.Join(stats, ", "))
	stats = produceStats(allegDistrib, total)
	fmt.Fprintf(o, "- Allegiances: %s\n", strings.Join(stats, ", "))
	stats = produceStats(securityDistrib, total)
	fmt.Fprintf(o, "- Security Levels: %s\n", strings.Join(stats, ", "))
}

func reportOnSectors(o io.Writer, sdb *SystemDatabase) {
	if len(sdb.sectors) == 0 {
		fmt.Fprintf(o, "Sectors: 0\n")
		return
	}

	x1, x2, y1, y2, z1, z2 := getBounds(sdb)
	fmt.Fprintf(o, "Sectors: %d [(%d-%d),(%d-%d),(%d-%d)]\n", len(sdb.sectors), x1, x2, y1, y2, z1, z2)
	populations := make([]int, 0, len(sdb.sectors))
	totalSecPop := 0
	for _, sector := range sdb.sectors {
		populations = append(populations, len(sector))
		totalSecPop += len(sector)
	}
	avg := average(totalSecPop, len(sdb.sectors))
	sort.Ints(populations)
	if len(populations) > 0 {
		fmt.Fprintf(o, "- Min/Avg/P95/Max Population: %d/%.2f/%d/%d\n", populations[0], avg, percentile(.95, populations), populations[len(populations)-1])
	}
}

func reportOnFacilities(o io.Writer, sdb *SystemDatabase) {
	fmt.Fprintf(o, "Facilities: %d\n", len(sdb.facilitiesByID))
	total := len(sdb.facilitiesByID)
	if total == 0 {
		return
	}

	withCommodities, withListings, planetary := 0, 0, 0
	govtDistrib := make(map[string]int, 32)
	allegDistrib := make(map[string]int, 32)
	typeDistrib := make(map[string]int, 32)
	padDistrib := make(map[string]int, 32)
	for _, facility := range sdb.facilitiesByID {
		if facility.HasFeatures(FeatCommodities) {
			withCommodities++
		}
		if len(facility.listings) > 0 {
			withListings++
		}
		if facility.HasFeatures(FeatPlanetary) {
			planetary++
		}
		govtDistrib[gomschema.GovernmentType_name[int32(facility.Government)]]++
		allegDistrib[gomschema.AllegianceType_name[int32(facility.Allegiance)]]++
		typeDistrib[gomschema.FacilityType_name[int32(facility.FacilityType)]]++
		if facility.SupportsPadSize(FeatLargePad) {
			padDistrib["Large"]++
		} else if facility.SupportsPadSize(FeatMediumPad) {
			padDistrib["Medium"]++
		} else if facility.SupportsPadSize(FeatSmallPad) {
			padDistrib["Small"]++
		} else {
			padDistrib["Unknown"]++
		}
	}

	fmt.Fprintf(o, "- Marked 'Has Commodities': %d (%.2f%%)\n", withCommodities, percentage(withCommodities, total))
	fmt.Fprintf(o, "- Known Listings: %d (%.2f%%)\n", withListings, percentage(withListings, total))

	stats := produceStats(typeDistrib, total)
	fmt.Fprintf(o, "- Types: %s\n", strings.Join(stats, ", "))
	stats = produceStats(govtDistrib, total)
	fmt.Fprintf(o, "- Governments: %s\n", strings.Join(stats, ", "))
	stats = produceStats(allegDistrib, total)
	fmt.Fprintf(o, "- Allegiances: %s\n", strings.Join(stats, ", "))
}

func (sdb *SystemDatabase) Stats(o io.Writer) {
	reportOnCommodities(o, sdb)
	reportOnSystems(o, sdb)
	reportOnSectors(o, sdb)
	reportOnFacilities(o, sdb)
}
