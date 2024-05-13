package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

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

func parseHeader(lines []string) (clubInfo, error) {
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

	return clubInfo{totalTables: tables, openTime: stTime, closeTime: endTime, pricePerHour: int64(pricePerHour)}, err
}

// parseTime expects format "XX:YY", where XX - hours, YY - minutes
func parseTime(str string) (time.Duration, error) {
	t, err := time.ParseDuration(str[0:2] + "h" + str[3:5] + "m")
	fmt.Println(t)
	if err != nil {
		return t, err
	}
	if t.Milliseconds() > time.Hour.Milliseconds()*24 {
		return t, errors.New("time should be less then 24 hour")
	}
	return t, nil
}
