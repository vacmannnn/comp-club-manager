package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	clientVisit = iota + 1
	clientSeat
	clientWaiting
	clientLeftByHimself
	clientLeft = iota + 7
	clientSeatAfterWaiting
	errorAction
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
		log.Fatal("Expected format: ./app log_file.txt")
	}

	bytesRead, err := os.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	fileContent := string(bytesRead)
	lines := strings.Split(fileContent, "\n")

	club, err := readHeader(lines)
	if err != nil {
		log.Fatal(err) // TODO
	}
	fmt.Println(lines[1][:5])

	var waiting queue
	freeTables := club.totalTables
	tables := make([]tableInfo, club.totalTables+1)
	clients := make(map[string]clientInfo)
	for i := 3; i < len(lines) && len(lines[i]) > 0; i++ {
		fmt.Println(lines[i])
		line := strings.Split(lines[i], " ")
		act, err := parseAct(line)
		if err != nil {
			log.Fatal(err)
		}
		client := clients[act.userName]
		switch act.id {
		case clientVisit:
			if act.time < club.seatTime {
				fmt.Printf("%s %d NotOpenYet\n", line[0], errorAction)
				continue
			}
			if client.statusID != clientLeft && client.statusID != clientLeftByHimself && client.statusID != 0 {
				fmt.Printf("%s %d YouShallNotPass\n", line[0], errorAction)
				continue
			}
			client.statusID = act.id
		case clientSeat:
			if tables[act.tableNum].isBusy {
				fmt.Printf("%s %d PlaceIsBusy\n", line[0], errorAction)
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
			client.statusID = clientSeat
			tables[act.tableNum].isBusy = true
			freeTables--
		case clientWaiting:
			if freeTables > 0 {
				fmt.Printf("%s %d ICanWaitNoLonger!\n", line[0], errorAction)
				continue
			}
			if len(waiting)+1 > club.totalTables {
				client.statusID = clientLeft
				fmt.Printf("%s %d %s\n", line[0], clientLeft, act.userName)
				continue
			}
			waiting.enqueue(act.userName)
		case clientLeftByHimself:
			timeUsed := act.time.Milliseconds() - client.seatTime
			profit := (timeUsed / millisecondsInHour) * club.pricePerHour
			if timeUsed%millisecondsInHour > 0 {
				profit += club.pricePerHour
			}

			tables[client.curTable].timeUsed += timeUsed
			tables[client.curTable].profit += profit
			client.statusID = clientLeftByHimself
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
			abc.statusID = clientSeatAfterWaiting
			clients[clientName] = abc
			fmt.Printf("%s %d %s %d\n", line[0], clientSeatAfterWaiting, act.userName, client.curTable)
		default:
			log.Fatal("Unknown action ID")
		}
		clients[act.userName] = client
	}
	for k, v := range clients {
		if v.statusID == clientLeft || v.statusID == clientLeftByHimself {
			continue
		}
		timeUsed := club.endTime.Milliseconds() - v.seatTime
		tables[v.curTable].timeUsed += timeUsed
		profit := (timeUsed / millisecondsInHour) * club.pricePerHour
		if timeUsed%millisecondsInHour > 0 {
			profit += club.pricePerHour
		}
		tables[v.curTable].profit += profit
		fmt.Printf("%s %d %s\n", lines[1][6:], clientLeft, k)
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

func parseAct(input []string) (action, error) {
	actionTime, err := parseTime(input[0])
	if err != nil {
		return action{}, fmt.Errorf("incorrect time format, %s", err)
	}

	actionID, err := strconv.Atoi(input[1])
	if err != nil {
		return action{}, fmt.Errorf("incorrect action ID, %s", err)
	}

	var table int
	if len(input) > 3 {
		table, err = strconv.Atoi(input[3])
		if table < 1 {
			err = errors.New("table num should be greater then 0")
		}
	}
	if err != nil {
		return action{}, fmt.Errorf("incorrect table number, %s", err)
	}

	return action{time: actionTime, id: actionID, userName: input[2], tableNum: table}, err
}

func readHeader(lines []string) (clubInfo, error) {
	tables, err := strconv.Atoi(lines[0])
	if tables < 1 {
		err = errors.New("total tables num should be greater then 0")
	}
	if err != nil {
		return clubInfo{}, fmt.Errorf("incorrect tables number, %s", err)
	}

	timeStr := strings.Split(lines[1], " ")
	stTime, err := parseTime(timeStr[0])
	if err != nil {
		return clubInfo{}, fmt.Errorf("incorrect time format, %s", err)
	}
	endTime, err := parseTime(timeStr[1])
	if err != nil {
		return clubInfo{}, fmt.Errorf("incorrect time format, %s", err)
	}

	pricePerHour, err := strconv.Atoi(lines[2])
	if pricePerHour < 1 {
		err = errors.New("total tables num should be greater then 0")
	}
	if err != nil {
		return clubInfo{}, fmt.Errorf("incorrect price number, %s, it should be positive integer number", err)
	}

	return clubInfo{totalTables: tables, seatTime: stTime, endTime: endTime, pricePerHour: int64(pricePerHour)}, err
}

// parseTime expects format "XX:YY", where XX - hours, YY - minutes
func parseTime(str string) (time.Duration, error) {
	return time.ParseDuration(str[0:2] + "h" + str[3:5] + "m")
}
