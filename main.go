package main

import (
	"fmt"
	"github.com/go-ping/ping"
)

// TODO: Peerよりいい名前募
type Peer struct {
	remtoeAddress string
	remoteIf string
	localAddress string
	localIf string
}

func main() {
	nic1 := Peer{remtoeAddress: "192.168.130.3", remoteIf: "enp0s10", localIf: "enp0s10", localAddress: "localhost"}
	nic2 := Peer{remtoeAddress: "10.10.10.3", remoteIf: "srsgre", localIf: "srsgre", localAddress: "localhost"}

	localPinger, err := ping.NewPinger(nic1.remtoeAddress)
	localPinger.SetPrivileged(true)
	localPinger.Count = 1
	if err != nil {
			panic(err)
	}

	remotePinger, err := ping.NewPinger(nic2.remtoeAddress)
	remotePinger.SetPrivileged(true)
	remotePinger.Count = 1
	if err != nil {
			panic(err)
	}

	// TODO: 無限ループに
	for i := 0; i < 3; i++ {
		localPinger.Run()
		remotePinger.Run()
	
		// TODO: 並列化
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
		// TODO: loging & sleep
	}
}

func setCost(cost int, addr string, nic string) error {
	// TODO: exec shell script
	return nil
}