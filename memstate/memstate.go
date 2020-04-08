package main

import (
	"encoding/json"
	"fmt"
	"runtime"
	"time"
)

// Monitor opa
type Monitor struct {
	//	Alloc,
	//	TotalAlloc,
	Sys uint64
	//	Mallocs,
	//	Frees,
	//	LiveObjects,
	//	PauseTotalNs uint64
	//	NumGC       uint32
	//	NumGorutine int
	CPU int
}

// NewMonitor opa
func NewMonitor(duration int) {

	var m Monitor
	var rtm runtime.MemStats
	var interval = time.Duration(duration) * time.Second
	for {

		<-time.After(interval)

		runtime.ReadMemStats(&rtm)
		//	m.NumGorutine = runtime.NumGoroutine()

		//		m.Alloc = rtm.Alloc
		//		m.TotalAlloc = rtm.TotalAlloc
		m.Sys = rtm.Sys
		//		m.Mallocs = rtm.Mallocs
		//		m.Frees = rtm.Frees
		m.CPU = runtime.NumCPU()

		//		m.LiveObjects = m.Mallocs - m.Frees

		//		m.PauseTotalNs = rtm.PauseTotalNs
		//		m.NumGC = rtm.NumGC
		b, _ := json.Marshal(m)
		fmt.Println(string(b))

	}
}
func main() {
	NewMonitor(10)
}
