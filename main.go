package main

import (
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/go-ping/ping"
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
	wlanNic := Peer{remtoeAddress: "192.168.120.1", remoteIf: "wlan0", localIf: "wlan1", localAddress: "localhost", preffered: false}
	ranNic := Peer{remtoeAddress: "10.10.10.1", remoteIf: "srsgre", localIf: "srsgre", localAddress: "localhost", preffered: true}

	file, err := os.OpenFile("zero.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	log.SetOutput(file)

	// TODO: 無限ループに
	for {
		wlanPinger, err := ping.NewPinger(wlanNic.remtoeAddress)
		wlanPinger.SetPrivileged(true)
		wlanPinger.Count = 5
		if err != nil {
			panic(err)
		}

		ranPinger, err := ping.NewPinger(ranNic.remtoeAddress)
		ranPinger.SetPrivileged(true)
		ranPinger.Count = 5
		if err != nil {
			panic(err)
		}

		// TODO: 並列化
		wlanPinger.Run()
		ranPinger.Run()

		wlanStats := wlanPinger.Statistics()
		ranStats := ranPinger.Statistics()

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

		time.Sleep(1 * time.Second)
	}
}

func setCost(cost int, peer *Peer) error {
	localArgs := []string{"cost-set.sh", "localhost", peer.localIf, strconv.Itoa(cost)}
	if err := exec.Command("/usr/bin/expect", localArgs...).Run(); err != nil {
		// TODO: wrap error
		return err
	}
	log.Printf("cost of %s is set to %d", peer.localAddress, cost)

	remoteArgs := []string{"cost-set.sh", "192.168.120.1", peer.remoteIf, strconv.Itoa(cost)}
	if err := exec.Command("/usr/bin/expect", remoteArgs...).Run(); err != nil {
		return err
	}
	log.Printf("cost of %s is set to %d", peer.remtoeAddress, cost)
	return nil
}
