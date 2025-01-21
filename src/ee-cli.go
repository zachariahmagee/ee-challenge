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

//var stations map[uint32][]uint32


func main() {
    if len(os.Args) < 1 {
        fmt.Println("[WARNING] ee-cli expects one argument, the path to a file\nUsage: ee-cli <PATH>")
        return
    }

    filepath := os.Args[1]


    stations, reports, err := parseFile(filepath)

    printResults(calculateUptime(stations, reports))

    if err != nil {
        fmt.Println("Error reading file", err)
    }

}

func parseFile(filepath string) (map[uint32][]uint32, []AvailabilityReport, error) {
    file, err := os.Open(filepath)
    if err != nil {
        return nil, nil, err
    }
    defer file.Close()

    var (
        stations    = make(map[uint32][]uint32)
        reports     []AvailabilityReport
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
                return nil, nil, fmt.Errorf("invaluid station line: %s", line)
            }
            stationID, _ := strconv.ParseUint(parts[0], 10, 32)
            for _, charger := range parts[1:] {
                chargerID, _ := strconv.ParseUint(charger, 10, 32)
                stations[uint32(stationID)] = append(stations[uint32(stationID)], uint32(chargerID))
            }
            continue
        }
        // Example report: 1001 50000 100000 true
        if section == "[Charger Availability Reports]" {
            parts := strings.Fields(line)
            if len(parts) != 4 {
                return nil, nil, fmt.Errorf("invalid report line: %s", line)
            }

            chargerID, _ := strconv.ParseUint(parts[0], 10, 32)
            start, _ := strconv.ParseUint(parts[1], 10, 64)
            end, _ := strconv.ParseUint(parts[2], 10, 64)
            up := parts[3] == "true"
            reports = append(reports, AvailabilityReport{
                ChargerID: uint32(chargerID),
                Start: start,
                End: end,
                Up: up,
            })
        }
    }
    return stations, reports, scanner.Err()
}

func calculateUptime(stations map[uint32][]uint32, reports []AvailabilityReport) map[uint32]int {
    chargerUptime := make(map[uint32][][2]uint64) // Charger ID -> uptime intervals [chargerID][]{start, end} 
    chargerDowntime := make(map[uint32][][2]uint64)
    for _, report := range reports {
        if report.Up {
            chargerUptime[report.ChargerID] = append(chargerUptime[report.ChargerID], [2]uint64{report.Start, report.End})
        } else {
            chargerDowntime[report.ChargerID] = append(chargerUptime[report.ChargerID], [2]uint64{report.Start, report.End})
        }
    }

    stationUptime := make(map[uint32]int)
    for stationID, chargerIDs := range stations {
        stationUptime[stationID] = int(0)
        totalUptime := uint64(0)
        totalTime := uint64(0)
        // fmt.Printf("Station ID: %s")
        for _, chargerID := range chargerIDs {
            intervals := mergeIntervals(chargerUptime[chargerID])
            //fmt.Printf("%d: %d, %d", chargerID, intervals[0], intervals[1])
            for _, interval := range intervals {
                totalUptime += interval[1] - interval[0]
                // fmt.Printf("interval %d %d", chargerID, totalUptime)
            }
            // fmt.Printf("total %d %d", chargerID, totalUptime)
        }
        if len(chargerIDs) > 0 {
            totalTime = calculateTotalTime(chargerUptime, chargerDowntime, chargerIDs)
        }
        if totalTime > 0 {
            stationUptime[stationID] = int((totalUptime * 100) / totalTime)
            stationUptime[stationID] = min(stationUptime[stationID], 100)
        }
    }
    return stationUptime
}

func mergeIntervals(intervals [][2]uint64) [][2]uint64 {
    if len(intervals) == 0 {
        return nil
    }
    // sort intervals by start time
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i][0] < intervals[j][0]
    })
    // initialize merged list with the first interval
    merged := [][2]uint64{intervals[0]}
    // process remaining intervals
    for _, interval := range intervals[1:] {
        // retrieve the last interval
        last := &merged[len(merged)-1]
        // check if the current interval overlaps or meets the last interval
        if interval[0] <= last[1] {
            // Merge: update the end time of the last interval
            last[1] = max(last[1], interval[1])
        } else {
            // No overlap: add the current interval as a new entity
            merged = append(merged, interval)
        }
    }
    return merged
}

func calculateTotalTime(chargerUptime map[uint32][][2]uint64, chargerDowntime map[uint32][][2]uint64, chargerIDs []uint32) uint64 {

    
    var intervals [][2]uint64
        
    for _, chargerID := range chargerIDs {
        if chargerIntervals, exists := chargerUptime[chargerID]; exists {
            intervals = append(intervals, chargerIntervals...)
        }
        if chargerIntervals, exists := chargerDowntime[chargerID]; exists {
            intervals = append(intervals, chargerIntervals...)
        }
    }

    if len(intervals) == 0 {
        return 0
    }

    merged := mergeIntervals(intervals)
    totalTime := merged[len(merged) - 1][1] - merged[0][0]

    return totalTime

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
