package main

import (
	"log"
	"os"
	"time"
	"strconv"
	"github.com/go-ping/ping"
	"encoding/csv"
)

var evalLogger *log.Logger

func main() {
	file, err := os.OpenFile("evalution.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	evalLogger = log.New(file, "", log.Lmicroseconds)
	evalLogger.Println("===ping===")
	execPing()
}

func execPing() {
	file, err := os.OpenFile("rtt_ocaca.csv", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	csvWriter := csv.NewWriter(file)
	csvWriter.Write([]string{"sequence", "times[us]", "RTT[us]"})
	
	pinger, err := ping.NewPinger("192.168.1.1")
	if err != nil {
		panic(err)
	}
	pinger.SetPrivileged(true)
	pinger.Interval = 10 * time.Millisecond

	now := time.Now()

	pinger.OnRecv = func(pkt *ping.Packet) {
		// evalLogger.Printf("%d bytes from %s: icmp_seq=%d time=%v\n", pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)

		csvWriter.Write([]string{
			strconv.Itoa(pkt.Seq),
			strconv.FormatInt(time.Since(now).Microseconds(), 10),
			strconv.FormatInt(pkt.Rtt.Microseconds(), 10),
		})
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
		evalLogger.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n", stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
		evalLogger.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n", stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}
	evalLogger.Println("ping start")
	go pinger.Run()

	time.Sleep(6 * time.Second)
	pinger.Stop()
	csvWriter.Flush()
	evalLogger.Println("ping done")
}
