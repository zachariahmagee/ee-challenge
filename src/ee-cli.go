package main

/*
* Input 1:
[Stations]
0 1001 1002
1 1003
2 1004

[Charger Availability Reports]
1001 0 50000 true
1001 50000 100000 true
1002 50000 100000 true
1003 25000 75000 false
1004 0 50000 true
1004 100000 200000 true

---
* Output 1
0 100
1 0
2 75
*/

import (
    "fmt"
    "os"
    "strconv"
    "strings"
    "bufio"
    "sort"
)

type AvailabilityReport struct {
    ChargerID uint32
    Start uint64
    End uint64
    Up bool
}

// map[uint32][]AvailabilityReports

func main() {
    if len(os.Args) < 2 {
        fmt.Print("Error\n")
        fmt.Fprintf(os.Stderr, "[WARNING] ee-cli expects one argument, the path to a file\nUsage: ee-cli <PATH>\n")
        return
    }

    filepath := os.Args[1]

    stations, err := parseFile(filepath)
    
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error reading file: %s", err.Error())
    }
    
    printResults(calculateUptime(stations))
}

func parseFile(filepath string) (map[uint32][]AvailabilityReport, error) {
    file, err := os.Open(filepath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var (
        stations    = make(map[uint32][]AvailabilityReport)
        chargerStations    = make(map[uint32]uint32)
        section     string
        scanner     = bufio.NewScanner(file)
    )

    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        if strings.HasPrefix(line, "[") {
            section = line
            continue
        }
        // Example station: 0 1001 1002
        if section == "[Stations]" {
            parts := strings.Fields(line)
            if len(parts) < 2 {
                return nil, fmt.Errorf("invalid station line: %s", line)
            }
            stationID, _ := strconv.ParseUint(parts[0], 10, 32)
            for _, charger := range parts[1:] {
                chargerID, _ := strconv.ParseUint(charger, 10, 32)
                chargerStations[uint32(chargerID)] = uint32(stationID)

            }
            continue
        }
        // Example report: 1001 50000 100000 true
        if section == "[Charger Availability Reports]" {
            parts := strings.Fields(line)
            if len(parts) != 4 {
                return nil, fmt.Errorf("invalid report line: %s", line)
            }

            chargerID, _ := strconv.ParseUint(parts[0], 10, 32)
            start, _ := strconv.ParseUint(parts[1], 10, 64)
            end, _ := strconv.ParseUint(parts[2], 10, 64)
            up := parts[3] == "true"

            if stationID, exists := chargerStations[uint32(chargerID)]; exists {
                stations[uint32(stationID)] = append(stations[uint32(stationID)], AvailabilityReport{
                    ChargerID: uint32(chargerID),
                    Start: start,
                    End: end,
                    Up: up,
                })
            }
        }
    }
    return stations, scanner.Err()
}

func calculateUptime(stations map[uint32][]AvailabilityReport) map[uint32]int {
    
    stationUptime := make(map[uint32]int)
    for stationID, reports := range stations {
        stationUptime[stationID] = int(0)
        
        merged := mergeReports(reports)
        if len(merged) == 0 { continue }
        upIntervals := filterSlice(merged, func (report AvailabilityReport) bool { return report.Up })
        uptime := reduce(upIntervals, 0, func (acc uint64, report AvailabilityReport) uint64 { acc += report.End - report.Start; return acc })
        totalTime := merged[len(merged) - 1].End - merged[0].Start
        
        if totalTime > 0 {
            stationUptime[stationID] = int((uptime * 100) / totalTime)
        }
    }
    return stationUptime
}

func mergeReports(intervals []AvailabilityReport) []AvailabilityReport {
    if len(intervals) == 0 {
        return nil
    }
    // sort intervals by start time
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i].Start < intervals[j].Start
    })
    // initialize merged list with the first interval
    merged := []AvailabilityReport{ intervals[0] }
    // process remaining intervals
    for _, interval := range intervals[1:] {
        // If the interval is down, skip it and continue 
        if !interval.Up {
            continue
        }
        // retrieve the last interval
        last := &merged[len(merged)-1] 
        // check if the current interval overlaps or meets the last interval
        if interval.Start <= last.End {
            // Merge: update the end time of the last interval
            last.End = max(last.End, interval.End)
        } else {
            // No overlap: add the current interval as a new entity
            merged = append(merged, interval)
        }
    }
    return merged
}

func printResults(uptime map[uint32]int) {
    keys := make([]uint32, 0, len(uptime))
    for key := range uptime {
        keys = append(keys, key)
    }
    sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

    for _, key := range keys {
        fmt.Printf("%d %d\n", key, uptime[key])
    }
}

func filterSlice[T any](input []T, condition func(T) bool) []T {
    result := []T{}
    for _, v:= range input {
        if condition(v) {
            result = append(result, v)
        }
    }
    return result
}

func reduce[T any, R any](slice []T, initial R, reducer func(R, T) R) R {
    result := initial
    for _, v := range slice {
        result = reducer(result, v)
    }
    return result
}
