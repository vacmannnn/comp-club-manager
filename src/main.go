package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// TODO - именованные константы для ID

type clubInfo struct {
	totalTables  int
	startTime    time.Duration
	endTime      time.Duration
	pricePerHour int
}

type action struct {
	time     time.Duration
	id       int
	userName string
	tableNum int
}

type clientInfo struct {
	startTime time.Duration
	endTime   time.Duration
	curTable  int
	statusID  int
	valid     bool
}

type queue []string

func (q *queue) enqueue(d string) {
	*q = append(*q, d)
}

func (q *queue) dequeue() string {
	dequeued := (*q)[0]
	*q = (*q)[1:]
	return dequeued
}

func (q *queue) len() int {
	return len(*q)
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal() // TODO
	}

	bytesRead, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err) // TODO
	}
	fileContent := string(bytesRead)
	lines := strings.Split(fileContent, "\n")

	club, err := readHeader(lines)
	// fmt.Println(club)
	if err != nil {
		log.Fatal(err) // TODO
	}
	fmt.Println(lines[1][:5])

	// TODO: why -1
	var waiting queue
	freeTables := club.totalTables
	busyTables := make([]bool, club.totalTables)
	clients := make(map[string]clientInfo)
	for i := 3; i < len(lines)-1; i++ {
		fmt.Println(lines[i])
		line := strings.Split(lines[i], " ")
		act, err := parseClient(line)
		if err != nil {
			continue // TODO
		}
		client := clients[act.userName]
		switch act.id {
		case 1:
			if act.time < club.startTime {
				fmt.Printf("%s 13 NotOpenYet\n", line[0]) // TODO : show correct time
				continue
			}
			if client.statusID%4 != 0 {
				fmt.Printf("%s 13 YouShallNotPass\n", line[0])
				continue
			}
			client.statusID = act.id
		case 2:
			client.startTime = act.time
			client.curTable = act.tableNum - 1
			busyTables[act.tableNum-1] = true
			freeTables--
		case 3:
			waiting.enqueue(act.userName)
		case 4:
			client.endTime = act.time
			busyTables[client.curTable] = false
			freeTables++
		default:
			// TODO
		}
		clients[act.userName] = client
	}
	fmt.Printf("%v\n", clients)
	fmt.Println(busyTables)
}

// TODO: return error
func parseClient(input []string) (action, error) {
	actionTime, err := parseTime(input[0])
	actionID, err := strconv.Atoi(input[1])

	var table int
	if len(input) > 3 {
		table, err = strconv.Atoi(input[3])
	}

	return action{time: actionTime, id: actionID, userName: input[2], tableNum: table}, err
}

// TODO: check if stTime > endTime and return first error
func readHeader(lines []string) (clubInfo, error) {
	tables, err := strconv.Atoi(lines[0])

	timeStr := strings.Split(lines[1], " ")
	stTime, err := parseTime(timeStr[0])
	endTime, err := parseTime(timeStr[1])

	pricePerHour, err := strconv.Atoi(lines[2])

	return clubInfo{totalTables: tables, startTime: stTime, endTime: endTime, pricePerHour: pricePerHour}, err
}

// parseTime expects format "XX:YY", where XX - hours, YY - minutes
func parseTime(str string) (time.Duration, error) {
	return time.ParseDuration(str[0:2] + "h" + str[3:5] + "m")
}
