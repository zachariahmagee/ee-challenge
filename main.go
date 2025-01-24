package main

/*
Author: Zachariah Magee
Date: January 22, 2025
Description: Coding challenge for Software Engineer @ Electric Era Technologies
Challenge Repo: https://gitlab.com/electric-era-public/coding-challenge-charger-uptime
*/

import (
    "fmt"
    "os"
    "ee-challenge/pkg/uptime"

)

func main() {
    if len(os.Args) < 2 {
        fmt.Print("Error\n")
        fmt.Fprintf(os.Stderr, "[WARNING] ee-cli expects one argument, the path to a file\nUsage: ee-cli <PATH>\n")
        return
    }

    filepath := os.Args[1]

    stations, err := uptime.ParseFile(filepath)
    
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error reading file: %s", err.Error())
    }
    if len(stations) == 0 {
        fmt.Printf("Error\n")
        return
    }
    
    uptime.PrintResults(uptime.CalculateUptime(stations))
}


