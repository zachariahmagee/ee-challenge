package test

import (
    "testing"
    "reflect"
    "ee-challenge/pkg/uptime"
)



func TestIntegration(t *testing.T) {
    filepath := "./data/test-input-1.txt"

    stations, err := uptime.ParseFile(filepath)
    if err != nil {
        t.Fatalf("Failed to parse file: %v", err)
    }

    uptime := uptime.CalculateUptime(stations)
    expected := map[uint32]int{
        0: 50,
        1: 75,
    }

    if !reflect.DeepEqual(uptime, expected) {
        t.Errorf("Expected %v, but got %v", expected, uptime)
    }
}
