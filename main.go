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
	wlanNic := Peer{remtoeAddress: "192.168.130.1", remoteIf: "enp0s10", localIf: "enp0s10", localAddress: "localhost"}
	ranNic := Peer{remtoeAddress: "10.10.10.1", remoteIf: "srsgre", localIf: "srsgre", localAddress: "localhost"}

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

		if wlanStats.AvgRtt < ranStats.AvgRtt {
			if err := setCost(1, wlanNic.localAddress, wlanNic.localIf); err != nil {
				panic(err)
			}
			if err := setCost(2, wlanNic.remtoeAddress, wlanNic.remoteIf); err != nil {
				panic(err)
			}
			// TODO: nic2に対してもsetCostを呼ぶ
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
