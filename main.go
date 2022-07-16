package main

import (
	"fmt"
	"github.com/go-ping/ping"
)

func main() {
	pinger, err := ping.NewPinger("www.google.com")
	pinger.SetPrivileged(true)
	if err != nil {
			panic(err)
	}
	pinger.Count = 1
	pinger.Run()
	stats := pinger.Statistics()
	fmt.Printf("loss: %f\n", stats.PacketLoss)
	fmt.Printf("Rtts: %s\n", stats.Rtts)
	fmt.Printf("AvgRtt: %s\n", stats.AvgRtt)
}