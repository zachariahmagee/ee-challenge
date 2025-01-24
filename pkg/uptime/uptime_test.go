package uptime

import (
    "testing"
    "reflect"
)

func TestMergeReports(t *testing.T) {
    reports := []AvailabilityReport{
        {ChargerID: 1001, Start: 100, End: 200, Up: true},
        {ChargerID: 1001, Start: 150, End: 300, Up: true},
        {ChargerID: 1001, Start: 400, End: 500, Up: false},
    }

    merged := MergeReports(reports)
    expected := []AvailabilityReport{
        {ChargerID: 1001, Start: 100, End: 300, Up: true},
        {ChargerID: 1001, Start: 400, End: 500, Up: false},
    }

    if !reflect.DeepEqual(merged, expected) {
        t.Errorf("Expected %v, but got %v", expected, merged)
    }
}


func TestParseFile_ValidInput(t *testing.T) {
    // Arrange
    filepath := "../../test/data/test-input-1.txt"

    // Act
    result, err := ParseFile(filepath)

    // Assert
    if err != nil {
        t.Fatalf("Expected no error, but got: %v", err)
    }
    if len(result) == 0 {
        t.Errorf("Expected stations to be parsed, but got empty result")
    }
}
