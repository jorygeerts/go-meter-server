package main

import "time"

type Measurement struct {
	Collection Collection
	Id int
	Value int
	Date time.Time
}
