package main

import (
	"log"
	"os"
	"os/exec"
	"time"
	"flag"
	"github.com/go-ping/ping"
)

var evalLogger *log.Logger

func main() {
	flag.Parse()
    cmd := flag.Arg(0)
	
	file, err := os.OpenFile("evalution.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	evalLogger = log.New(file, "", log.Lmicroseconds)

	done := make(chan struct{})
	timer := time.After(3 * time.Second)

	switch cmd {
	case "iperf3":
		go execIperf3(done)
		evalLogger.Println("===iperf3===")
	case "ping":
		go execPing(done)
		evalLogger.Println("===ping===")
	default:
		panic("no flag")
	}

	for {
		select {
		case <- done:
			args := []string{"qdisc", "change", "dev", "wlangre", "root", "handle", "1:0", "netem", "delay", "1ms"}
			if err := exec.Command("/usr/sbin/tc", args...).Run(); err != nil {
				// TODO: wrap error
				panic(err)
			}
			return
		case <-timer:
			args := []string{"qdisc", "change", "dev", "wlangre", "root", "handle", "1:0", "netem", "delay", "5ms"}
			if err := exec.Command("/usr/sbin/tc", args...).Run(); err != nil {
				// TODO: wrap error
				panic(err)
			}
			evalLogger.Println("exec tc cmd")
		}
	}
}

func execIperf3(done chan struct{}) {
	evalLogger.Println("start")
	args := []string{"-c192.168.3.4", "-b100M", "-i1", "-t6"}
	if err := exec.Command("/usr/bin/iperf3", args...).Run(); err != nil {
		// TODO: wrap error
		panic(err)
	}
	evalLogger.Println("done")
	close(done)
}

func execPing(done chan struct{}) {
	pinger, err := ping.NewPinger("192.168.3.3")
	if err != nil {
		panic(err)
	}
	pinger.SetPrivileged(true)
	pinger.Interval = 100 * time.Millisecond
	
	timer := time.After(6 * time.Second)

	pinger.OnRecv = func(pkt *ping.Packet) {
			evalLogger.Printf("%d bytes from %s: icmp_seq=%d time=%v\n",
					pkt.Nbytes, pkt.IPAddr, pkt.Seq, pkt.Rtt)
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
			evalLogger.Printf("\n--- %s ping statistics ---\n", stats.Addr)
			evalLogger.Printf("%d packets transmitted, %d packets received, %v%% packet loss\n",
					stats.PacketsSent, stats.PacketsRecv, stats.PacketLoss)
			evalLogger.Printf("round-trip min/avg/max/stddev = %v/%v/%v/%v\n",
					stats.MinRtt, stats.AvgRtt, stats.MaxRtt, stats.StdDevRtt)
	}

	evalLogger.Printf("PING %s (%s):\n", pinger.Addr(), pinger.IPAddr())

	evalLogger.Println("start")
	go pinger.Run()

	<-timer
	pinger.Stop()
	evalLogger.Println(pinger.Statistics().AvgRtt)
	evalLogger.Println("done")
	close(done)
}
