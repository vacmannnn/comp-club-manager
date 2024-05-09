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
	seatTime     time.Duration
	endTime      time.Duration
	pricePerHour int64
}

type tableInfo struct {
	timeUsed int64
	profit   int64
	isBusy   bool
}

type action struct {
	time     time.Duration
	id       int
	userName string
	tableNum int
}

type clientInfo struct {
	seatTime int64
	curTable int
	statusID int
	valid    bool
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

const millisecondsInHour = 3_600_000

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
	tables := make([]tableInfo, club.totalTables+1)
	clients := make(map[string]clientInfo)
	for i := 3; i < len(lines)-1; i++ {
		fmt.Println(lines[i])
		line := strings.Split(lines[i], " ")
		act, err := parseAct(line)
		if err != nil {
			continue // TODO
		}
		client := clients[act.userName]
		switch act.id {
		case 1:
			if act.time < club.seatTime {
				fmt.Printf("%s 13 NotOpenYet\n", line[0]) // TODO : show correct time
				continue
			}
			if client.statusID != 4 && client.statusID != 0 {
				fmt.Printf("%s 13 YouShallNotPass\n", line[0])
				continue
			}
			client.statusID = act.id
		case 2:
			if tables[act.tableNum].isBusy {
				fmt.Printf("%s 13 PlaceIsBusy\n", line[0])
				continue
			}
			if client.statusID%10 == 2 {
				tables[client.curTable].isBusy = false
				tables[act.tableNum].isBusy = true
				tables[client.curTable].timeUsed += act.time.Milliseconds() - client.seatTime
				client.seatTime = act.time.Milliseconds()
				client.curTable = act.tableNum
			}
			client.seatTime = act.time.Milliseconds()
			client.curTable = act.tableNum
			client.statusID = 2
			tables[act.tableNum].isBusy = true
			freeTables--
		case 3:
			if freeTables > 0 {
				fmt.Printf("%s 13 ICanWaitNoLonger!\n", line[0])
				continue
			}
			if len(waiting)+1 > club.totalTables {
				client.statusID = 11
				fmt.Printf("%s 11 %s\n", line[0], act.userName)
				continue
			}
			waiting.enqueue(act.userName)
		case 4:
			// TODO: Добавить время, но еще округлить часы и добавить к цене за счет этого
			timeUsed := act.time.Milliseconds() - client.seatTime
			profit := (timeUsed / millisecondsInHour) * club.pricePerHour
			if timeUsed%millisecondsInHour > 0 {
				profit += club.pricePerHour
			}

			tables[client.curTable].timeUsed += timeUsed
			tables[client.curTable].profit += profit
			client.statusID = 4
			if waiting.len() == 0 {
				freeTables++
				tables[client.curTable].isBusy = false
				clients[act.userName] = client
				continue
			}
			clientName := waiting.dequeue()
			abc := clients[clientName]
			abc.curTable = client.curTable
			abc.seatTime = act.time.Milliseconds()
			abc.statusID = 12
			clients[clientName] = abc
			fmt.Printf("%s 12 %s %d\n", line[0], act.userName, client.curTable)
		default:
			// TODO
		}
		clients[act.userName] = client
	}
	for k, v := range clients {
		if v.statusID == 4 || v.statusID == 11 {
			continue
		}
		timeUsed := club.endTime.Milliseconds() - v.seatTime
		tables[v.curTable].timeUsed += timeUsed
		profit := (timeUsed / millisecondsInHour) * club.pricePerHour
		if timeUsed%millisecondsInHour > 0 {
			profit += club.pricePerHour
		}
		tables[v.curTable].profit += profit
		fmt.Printf("%s 11 %s\n", lines[1][6:], k)
	}
	fmt.Println(lines[1][6:])
	for i, v := range tables {
		if i == 0 {
			continue
		}

		var t time.Time
		t = t.Add(time.Duration(v.timeUsed) * time.Millisecond)

		fmt.Printf("%d %d %s\n", i, v.profit, t.Format("15:04"))
	}
}

// TODO: return error
func parseAct(input []string) (action, error) {
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

	return clubInfo{totalTables: tables, seatTime: stTime, endTime: endTime, pricePerHour: int64(pricePerHour)}, err
}

// parseTime expects format "XX:YY", where XX - hours, YY - minutes
func parseTime(str string) (time.Duration, error) {
	return time.ParseDuration(str[0:2] + "h" + str[3:5] + "m")
}
