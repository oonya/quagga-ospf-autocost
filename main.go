package main

import (
	"fmt"
	"github.com/go-ping/ping"
	"log"
	"os"
	"time"
)

// TODO: Peerよりいい名前募
type Peer struct {
	remtoeAddress string
	remoteIf      string
	localAddress  string
	localIf       string
	preffered     bool
}

func main() {
	wlanNic := Peer{remtoeAddress: "192.168.130.1", remoteIf: "enp0s10", localIf: "enp0s10", localAddress: "localhost", preffered: false}
	ranNic := Peer{remtoeAddress: "10.10.10.1", remoteIf: "srsgre", localIf: "srsgre", localAddress: "localhost", preffered: true}

	file, err := os.OpenFile("zero.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	log.SetOutput(file)

	// TODO: 無限ループに
	for i := 0; i < 3; i++ {
		wlanPinger, err := ping.NewPinger(wlanNic.remtoeAddress)
		wlanPinger.SetPrivileged(true)
		wlanPinger.Count = 1
		if err != nil {
			panic(err)
		}

		ranPinger, err := ping.NewPinger(ranNic.remtoeAddress)
		ranPinger.SetPrivileged(true)
		ranPinger.Count = 1
		if err != nil {
			panic(err)
		}

		// TODO: 並列化
		wlanPinger.Run()
		ranPinger.Run()

		wlanStats := wlanPinger.Statistics()
		fmt.Printf("loss1: %f\n", wlanStats.PacketLoss)
		fmt.Printf("AvgRtt1: %s\n", wlanStats.AvgRtt)

		ranStats := ranPinger.Statistics()
		fmt.Printf("loss2: %f\n", ranStats.PacketLoss)
		fmt.Printf("AvgRtt2: %s\n", ranStats.AvgRtt)

		if wlanStats.AvgRtt < ranStats.AvgRtt && !wlanNic.preffered {
			if err = setCost(1, &wlanNic); err != nil {
				log.Panicf("cost of wlan couldnt set to %d\n", 1)
			}
			if err = setCost(2, &ranNic); err != nil {
				log.Panicf("cost of ran couldnt set to %d\n", 2)
			}
			wlanNic.preffered = true
			ranNic.preffered = false
			log.Println("wlan became priority")
		}
		if ranStats.AvgRtt < wlanStats.AvgRtt && !ranNic.preffered {
			if err = setCost(1, &ranNic); err != nil {
				log.Panicf("cost of ran couldnt set to %d\n", 1)
			}
			if err = setCost(2, &wlanNic); err != nil {
				log.Panicf("cost of wlan couldnt set to %d\n", 2)
			}
			ranNic.preffered = true
			wlanNic.preffered = false
			log.Println("RAN became priority")
		}

		fmt.Println()
		time.Sleep(1 * time.Second)
	}
}

func setCost(cost int, peer *Peer) error {
	// TODO: exec shell script
	log.Printf("cost of %s is set to %d", peer.localAddress, cost)
	log.Printf("cost of %s is set to %d", peer.remtoeAddress, cost)
	return nil
}
