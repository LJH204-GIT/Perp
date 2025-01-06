package main

import (
	"time"
)

func main() {
	start := time.Now()
	AssignWork(ReadFile("data.csv"))
	//AssignWork(ReadFile("a.txt"))
	allTime := time.Since(start)
	println("\033[1F\033[s\033[1000F" + "\033[36m\033[2K" + "                           已完成! 用时:" + allTime.String() + "\033[0m\033[u")
	time.Sleep(9999 * time.Second)
}
