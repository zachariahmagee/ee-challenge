# Availability Report Processor

This application reads charger availability data from a file, processes it to calculate uptime percentages for stations, and outputs the results in a sorted format.

## Features 

	•	Parses a configuration file with station and charger availability reports.
	•	Calculates uptime percentages for each station based on availability reports.
	•	Outputs the results in a sorted format for easy readability.

## Requirements

    •	Go (Golang) v1.18+
	•	Compatible with major operating systems (Linux, macOS, Windows).

---

## How to Compile

1.	Ensure that you have Go installed. You can download it from [the Go website](https://go.dev/dl/).
2.	Clone this repository to your local machine:
```bash
git clone https://github.com/your-repo/ee-challenge.git
cd ee-challenge
```

3. Build the application using the go build command:
```sh
go build -o station-uptime src/main.go
```
This will create an executable named station-uptime in the current directory.

---

## How to Run

1. Prepare an input file in the following format:
```
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
```

2. Run the application by providing the path to your input file:
```sh
./station-uptime path/to/input-file.txt
```

3. The output will be a list of station IDs and the associated uptime percentages:
```
0 100
1 0
2 75
```
---

## Input File Format

The input file consists of two sections:
	1.	Stations:
	•	Lists station IDs and their associated charger IDs.
	•	Format: <StationID> <ChargerID1> <ChargerID2> ...
	2.	Charger Availability Reports:
	•	Lists charger availability data, including start and end times of availability intervals.
	•	Format: <ChargerID> <StartTime> <EndTime> <Up (true/false)>

---

## Tests

Run tests with 
```sh
go test ./...

```
