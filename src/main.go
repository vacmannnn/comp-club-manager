package main

import (
	"fmt"
	"log"
	"os"
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
	errorID
)

type clubInfo struct {
	totalTables  int
	openTime     time.Duration
	closeTime    time.Duration
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
	//valid    bool
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

	club, err := parseHeader(lines)
	if err != nil {
		log.Fatal(err) // TODO
	}
	fmt.Println(lines[1][:5])

	var waiting queue
	var previousActTime time.Duration
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
		if act.tableNum > club.totalTables {
			log.Fatal("table number should be lower than total tables")
		}
		if act.time.Milliseconds() < previousActTime.Milliseconds() ||
			act.time.Milliseconds() > club.closeTime.Milliseconds() {
			log.Fatal("incorrect amount of time")
		}
		previousActTime = act.time
		client := clients[act.userName]

		switch act.id {
		case clientVisit:
			if act.time < club.openTime {
				fmt.Printf("%s %d NotOpenYet\n", line[0], errorID)
				break
			}
			if client.statusID != clientLeft && client.statusID != clientLeftByHimself && client.statusID != 0 {
				fmt.Printf("%s %d YouShallNotPass\n", line[0], errorID)
				break
			}
			client.statusID = clientVisit

		case clientSeat:
			if tables[act.tableNum].isBusy {
				fmt.Printf("%s %d PlaceIsBusy\n", line[0], errorID)
				break
			}
			// the client moved to another table
			if client.statusID%10 == 2 {
				tables[client.curTable].isBusy = false
				tables[act.tableNum].isBusy = true
				profit, timeUsed := calcProfit(client.seatTime, act.time.Milliseconds(), club.pricePerHour)
				tables[client.curTable].timeUsed += timeUsed
				tables[client.curTable].profit += profit
			} else {
				freeTables--
			}
			client.seatTime = act.time.Milliseconds()
			client.curTable = act.tableNum
			client.statusID = clientSeat

		case clientWaiting:
			if freeTables > 0 {
				fmt.Printf("%s %d ICanWaitNoLonger!\n", line[0], errorID)
				break
			}
			if len(waiting)+1 > club.totalTables {
				client.statusID = clientLeft
				fmt.Println(line[0], clientLeft, act.userName)
				break
			}
			waiting.enqueue(act.userName)

		case clientLeftByHimself:
			profit, timeUsed := calcProfit(client.seatTime, act.time.Milliseconds(), club.pricePerHour)
			tables[client.curTable].timeUsed += timeUsed
			tables[client.curTable].profit += profit
			client.statusID = clientLeftByHimself

			if waiting.len() == 0 {
				freeTables++
				tables[client.curTable].isBusy = false
				clients[act.userName] = client
				break
			}

			clientName := waiting.dequeue()
			waitedClient := clientInfo{curTable: client.curTable, seatTime: act.time.Milliseconds(),
				statusID: clientSeatAfterWaiting}
			clients[clientName] = waitedClient
			fmt.Println(line[0], clientSeatAfterWaiting, act.userName, client.curTable)

		default:
			log.Fatal("Unknown action ID")
		}
		clients[act.userName] = client
	}

	// check if all clients left
	for k, v := range clients {
		if v.statusID == clientLeft || v.statusID == clientLeftByHimself {
			continue
		}
		// possible case if client was in club, but he waited all the time and didn't seat at all
		if v.statusID%10 == 2 {
			profit, timeUsed := calcProfit(v.seatTime, club.closeTime.Milliseconds(), club.pricePerHour)
			tables[v.curTable].timeUsed += timeUsed
			tables[v.curTable].profit += profit
		}
		v.statusID = clientLeft
		fmt.Println(lines[1][6:], clientLeft, k)
	}

	fmt.Println(lines[1][6:])
	for i, v := range tables {
		if i == 0 {
			continue
		}

		var t time.Time
		t = t.Add(time.Duration(v.timeUsed) * time.Millisecond)

		fmt.Println(i, v.profit, t.Format("15:04"))
	}
}
