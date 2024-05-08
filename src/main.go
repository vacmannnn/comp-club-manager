package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

type clubInfo struct {
	tables       int
	startTime    time.Duration
	endTime      time.Duration
	pricePerHour int
}

type action struct {
	time     time.Duration
	id       int
	user     string
	tableNum int
}

func main() {
	if len(os.Args) != 2 {
		// TODO
		log.Fatal()
	}
	bytesRead, err := os.ReadFile(os.Args[1])
	if err != nil {
		// TODO
		log.Fatal(err)
	}
	fileContent := string(bytesRead)
	lines := strings.Split(fileContent, "\n")
	club, err := readHeader(lines)
	fmt.Println(club)
	if err != nil {
		// TODO
		log.Fatal(err)
	}

	// TODO: why -1
	for i := 3; i < len(lines)-1; i++ {
		newAction, _ := parseClient(lines[i])
		fmt.Println(newAction)
	}
}

// TODO: return error
func parseClient(str string) (action, error) {
	info := strings.Split(str, " ")
	actionTime, err := parseTime(info[0])
	actionID, err := strconv.Atoi(info[1])

	var table int
	if len(info) > 3 {
		table, err = strconv.Atoi(info[3])
	}
	
	return action{time: actionTime, id: actionID, user: info[2], tableNum: table}, err
}

// TODO: check if stTime > endTime and return first error
func readHeader(lines []string) (clubInfo, error) {
	tables, err := strconv.Atoi(lines[0])

	timeStr := strings.Split(lines[1], " ")
	stTime, err := parseTime(timeStr[0])
	endTime, err := parseTime(timeStr[1])

	pricePerHour, err := strconv.Atoi(lines[2])

	return clubInfo{tables: tables, startTime: stTime, endTime: endTime, pricePerHour: pricePerHour}, err
}

// parseTime expects format "XX:YY", where XX - hours, YY - minutes
func parseTime(str string) (time.Duration, error) {
	return time.ParseDuration(str[0:2] + "h" + str[3:5] + "m")
}
