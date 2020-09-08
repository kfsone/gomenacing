package main

import (
	"fmt"
	"io"
	"math"
	"sort"
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
			return list[index] + (list[index + 1] - list[index]) / 2
		} else {
			return list[index]
		}
	} else if len(list) == 1 {
		return list[0]
	}
	return 0
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
		fmt.Fprintf(o, "- %s: %d\n", list[0])
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

func countFacilities(system *System, predicate func(*Facility) bool) (count int) {
	for _, facility := range system.facilities {
		var val int
		if predicate(facility) {
			val = 1
		}
		count += val
	}
	return
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

func getBounds(sdb *SystemDatabase) (x1, x2, y1, y2, z1, z2 int){
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

func (sdb *SystemDatabase) Stats(o io.Writer) {
	fmt.Fprintf(o, "Commodities: %d\n", len(sdb.commoditiesByID))
	fmt.Fprintf(o, "Systems: %d\n", len(sdb.systemsByID))

	if len(sdb.systemIDs) > 0 {
		analyzeSystems(o, "w/facilities", sdb.systemsByID, func(s *System) int { return len(s.facilities) })
		analyzeSystems(o, "planetary", sdb.systemsByID, func(s *System) int { return countFacilities(s, func (f *Facility) bool {
			return f.HasFeatures(FeatPlanetary)
		})})
		populated := countSystems(sdb, func(s *System) bool { return s.Populated })
		fmt.Fprintf(o, "- Populated: %d (%.2f%%)\n", populated, percentage(populated, len(sdb.systemsByID)))
		permits := countSystems(sdb, func(s *System) bool { return s.NeedsPermit })
		fmt.Fprintf(o, "- Need permit: %d (%.2f%%)\n", permits, percentage(permits, len(sdb.systemsByID)))
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