package main

import (
	"fmt"
	"time"
	"os"
	"github.com/go-ping/ping"
	"log"
)

// TODO: Peerよりいい名前募
type Peer struct {
	remtoeAddress string
	remoteIf string
	localAddress string
	localIf string
}

func main() {
	nic1 := Peer{remtoeAddress: "192.168.130.1", remoteIf: "enp0s10", localIf: "enp0s10", localAddress: "localhost"}
	nic2 := Peer{remtoeAddress: "10.10.10.1", remoteIf: "srsgre", localIf: "srsgre", localAddress: "localhost"}

	file, err := os.OpenFile("zero.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
    if err != nil {
        panic(err)
    }
    defer file.Close()
	log.SetOutput(file)
	
	// TODO: 無限ループに
	for i := 0; i < 3; i++ {
		localPinger, err := ping.NewPinger(nic1.remtoeAddress)
		localPinger.SetPrivileged(true)
		localPinger.Count = 3
		if err != nil {
			panic(err)
		}
		
		remotePinger, err := ping.NewPinger(nic2.remtoeAddress)
		remotePinger.SetPrivileged(true)
		remotePinger.Count = 3
		if err != nil {
				panic(err)
		}

		// TODO: 並列化
		localPinger.Run()
		remotePinger.Run()
	
		localStats := localPinger.Statistics()
		fmt.Printf("loss1: %f\n", localStats.PacketLoss)
		fmt.Printf("AvgRtt1: %s\n", localStats.AvgRtt)

		remoteStats := remotePinger.Statistics()
		fmt.Printf("loss2: %f\n", remoteStats.PacketLoss)
		fmt.Printf("AvgRtt2: %s\n", remoteStats.AvgRtt)

		if localStats.AvgRtt < remoteStats.AvgRtt {
			if err := setCost(1, nic1.localAddress, nic1.localIf); err != nil {
				panic(err)
			}
			if err := setCost(2, nic1.remtoeAddress, nic1.remoteIf); err != nil {
				panic(err)
			}
		}

		fmt.Println()
		time.Sleep(1 * time.Second)
	}
}

func setCost(cost int, addr string, nic string) error {
	// TODO: exec shell script
	log.Printf("cost of %s is set to %d", addr, cost)
	return nil
}
