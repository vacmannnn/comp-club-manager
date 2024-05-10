package main

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

func calcProfit(seatTime, leftTime, pricePerHour int64) (int64, int64) {
	timeUsed := leftTime - seatTime
	profit := (timeUsed / millisecondsInHour) * pricePerHour
	if timeUsed%millisecondsInHour > 0 {
		profit += pricePerHour
	}
	return profit, timeUsed
}
