package uptime

/*
Author: Zachariah Magee
Date: January 22, 2025
Description: Coding challenge for Software Engineer @ Electric Era Technologies
Challenge Repo: https://gitlab.com/electric-era-public/coding-challenge-charger-uptime
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
/*** Input 1: ***
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

---*** Output 1 ***
0 100
1 0
2 75
*/

// parseFile parses a given file to build mappings of stations to their availability reports.
// It returns a map where the key is the station ID and the value is a slice of AvailabilityReports.
// The file is expected to contain two sections: [Stations] and [Charger Availability Reports].
func ParseFile(filepath string) (map[uint32][]AvailabilityReport, error) {
    // Open the file for reading. Return an error if the file cannot be opened.
    file, err := os.Open(filepath)
    if err != nil {
        return nil, err
    }
    defer file.Close() // Ensure the file is closed when the function exits.
    
    // Initialize the data structures:
    // - `stations`: Maps station IDs to a list of the associated charger's availability reports.
    // - `chargerStations`: Maps charger IDs to their associated station IDs.
    // - `section`: Tracks the current section of the file being processed.
    var (
        stations    = make(map[uint32][]AvailabilityReport) 
        chargerStations    = make(map[uint32]uint32)
        section     string
        scanner     = bufio.NewScanner(file)
    )
    // process the file line by line
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" {
            continue
        }
        // Capture the current section "[Stations]" or "[Charger Availability Report]"
        if strings.HasPrefix(line, "[") {
            section = line
            continue
        }
        // Process "Station" data - Example station: 0 1001 1002
        if section == "[Stations]" {
            parts := strings.Fields(line)
            if len(parts) < 2 {
                return nil, fmt.Errorf("invalid station line: %s", line)
            }
            stationID, _ := strconv.ParseUint(parts[0], 10, 32)
            // map chargerID (key) to the stationID (value) for quick access when creating reports
            for _, charger := range parts[1:] {
                chargerID, _ := strconv.ParseUint(charger, 10, 32)
                chargerStations[uint32(chargerID)] = uint32(stationID)
            }
            continue
        }
        // Process reports and map it to the associated station - Ex. report: 1001 50000 100000 true
        if section == "[Charger Availability Reports]" {
            parts := strings.Fields(line)
            if len(parts) != 4 {
                return nil, fmt.Errorf("invalid report line: %s", line)
            }

            chargerID, _ := strconv.ParseUint(parts[0], 10, 32) // unsigned 32 bit integer
            start, _ := strconv.ParseUint(parts[1], 10, 64) // unsigned 64 bit integer
            end, _ := strconv.ParseUint(parts[2], 10, 64) // unsigned 64 bit integer
            up := parts[3] == "true" || parts[3] != "false"// boolean - up (true), down (false)

            // Find the station ID for the charger and append the report to the station's list of reports.  
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

// calculateUptime calculates the uptime percentage for each station based on its availability reports.
// It returns a map where the key is the station ID and the value is the uptime percentage (0-100).
func CalculateUptime(stations map[uint32][]AvailabilityReport) map[uint32]int {
    // Initialize a map to store the uptime percentage for each station.
    stationUptime := make(map[uint32]int)
    for stationID, reports := range stations {
        // Set the initial uptime percentage to 0.
        stationUptime[stationID] = int(0)
        // Merge overlapping or adjacent availability intervals.
        merged := MergeReports(reports)
        if len(merged) == 0 { 
            continue 
        }
        // Filter out only the intervals where the station was up.
        upIntervals := FilterSlice(merged, func (report AvailabilityReport) bool { return report.Up })
        // Calculate the total uptime by summing the durations of the up intervals.
        uptime := Reduce(upIntervals, 0, func (acc uint64, report AvailabilityReport) uint64 { acc += report.End - report.Start; return acc })
        // Calculate the total time covered by the merged intervals. 
        totalTime := merged[len(merged) - 1].End - merged[0].Start
        fmt.Printf("station: %d total: %d up:%d\n", stationID, uptime, totalTime)
        // Calculate the uptime percentage if the total time is greater than 0.
        if totalTime > 0 {
            stationUptime[stationID] = int((uptime * 100) / totalTime)
        }
    }
    return stationUptime
}

func MergeReports(intervals []AvailabilityReport) []AvailabilityReport {
    if len(intervals) == 0 {
        return nil
    }
    // sort intervals by start time (unless start times are equal, then sort by the end time)
    sort.Slice(intervals, func(i, j int) bool {
        if intervals[i].Start == intervals[j].Start {
            return intervals[i].End < intervals[j].End
        }
        return intervals[i].Start < intervals[j].Start
    })
    
    // initialize merged list with the first interval
    merged := []AvailabilityReport{ intervals[0] }
    //fmt.Printf("intervals: %d, %d\n", intervals[0].Start, intervals[1].End)
    // process remaining intervals
    for _, interval := range intervals[1:] {
        // retrieve the last interval
        last := &merged[len(merged)-1] 
        
        // ensure that both intervals are up
        up := interval.Up && last.Up
        // check if the current up interval overlaps or meets the last up interval
        if up && interval.Start <= last.End {
            // Merge: update the end time of the last interval
            last.End = max(last.End, interval.End)
        } else {
            // No overlap: add the current interval as a new entity
            merged = append(merged, interval)
        }
    }
    return merged
}

// printResults prints the uptime percentages for each station in ascending order of station IDs.
// It takes a map where the key is the station ID and the value is the uptime percentage.
func PrintResults(uptime map[uint32]int) {
    // Create a slice to hold the station IDs for sorting.
    keys := make([]uint32, 0, len(uptime))
    for key := range uptime {
        keys = append(keys, key)
    }
    // Sort the station IDs in ascending order.
    sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
    // Iterate through the sorted station IDs and print the results.
    // Format: "<StationID> <UptimePercentage>"
    for _, key := range keys {
        fmt.Printf("%d %d\n", key, uptime[key])
    }
}

func PrintReports(reports []AvailabilityReport) {
    fmt.Printf("-----Availability Reports-----\n")
    for _, report := range reports {
        fmt.Printf("%+v\n\n", report)
    }
}

// filterSlice filters a slice based on a provided condition function.
// It takes an input slice and a condition function that determines whether an element should be included.
// Returns a new slice containing only the elements that satisfy the condition.
func FilterSlice[T any](input []T, condition func(T) bool) []T {
    result := []T{}
    for _, v:= range input {
        if condition(v) {
            result = append(result, v)
        }
    }
    return result
}

// reduce reduces a slice to a single value using a provided reducer function.
// It takes a slice, an initial value, and a reducer function that combines the accumulated value with each element.
// Returns the final accumulated value.
func Reduce[T any, R any](slice []T, initial R, reducer func(R, T) R) R {
    result := initial
    for _, v := range slice {
        result = reducer(result, v)
    }
    return result
}
