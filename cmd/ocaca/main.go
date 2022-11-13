package main

import (
	"log"
	"os"
	"time"

	"github.com/go-ping/ping"
	"github.com/oonya/quagga-ospf-autocost/usecase/ospfd"
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
	localExp, remoteExp, err := ospfd.ConnectOspfd()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := localExp.Close(); err != nil {
			log.Printf("localExp.Close failed: %v", err)
		}
	}()
	defer func() {
		if err := remoteExp.Close(); err != nil {
			log.Printf("remoteExp.Close failed: %v", err)
		}
	}()
	
	wlanNic := Peer{remtoeAddress: "10.10.20.1", remoteIf: "wlangre", localIf: "wlangre", localAddress: "localhost", preffered: false}
	ranNic := Peer{remtoeAddress: "10.10.10.1", remoteIf: "srsgre", localIf: "srsgre", localAddress: "localhost", preffered: true}

	file, err := os.OpenFile("ocaca.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	log.SetOutput(file)
	log.SetFlags(log.Lmicroseconds)

	for {
		wlanCH := make(chan time.Duration)
		ranCH := make(chan time.Duration)
		go measureRTT(wlanNic.remtoeAddress, wlanCH)
		go measureRTT(ranNic.remtoeAddress, ranCH)

		wlanStats := <-wlanCH
		ranStats := <-ranCH

		if wlanStats < ranStats && !wlanNic.preffered {
			ospfd.CostSet(localExp, 1, wlanNic.localIf)
			ospfd.CostSet(localExp, 2, ranNic.localIf)
			ospfd.CostSet(remoteExp, 1, wlanNic.remoteIf)
			ospfd.CostSet(remoteExp, 2, ranNic.remoteIf)

			wlanNic.preffered = true
			ranNic.preffered = false
			log.Println("wlan became priority")
		}
		if ranStats < wlanStats && !ranNic.preffered {
			ospfd.CostSet(localExp, 2, wlanNic.localIf)
			ospfd.CostSet(localExp, 1, ranNic.localIf)
			ospfd.CostSet(remoteExp, 2, wlanNic.remoteIf)
			ospfd.CostSet(remoteExp, 1, ranNic.remoteIf)

			ranNic.preffered = true
			wlanNic.preffered = false
			log.Println("RAN became priority")
		}
	}
}

func measureRTT(addr string, ch chan time.Duration) {
	pingger, err := ping.NewPinger(addr)
	pingger.SetPrivileged(true)
	pingger.Count = 2
	if err != nil {
		panic(err)
	}

	pingger.Run()
	stats := pingger.Statistics()
	ch <- stats.AvgRtt
}