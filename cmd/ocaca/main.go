package main

import (
	"log"
	"os"
	"time"
	"encoding/csv"
	"strconv"

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

var ocacaLogger *log.Logger
var now time.Time

func main() {
	file, err := os.OpenFile("ocaca.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	ocacaLogger = log.New(file, "", log.Lmicroseconds)
	log.SetFlags(log.Lmicroseconds)

	wlanCsv, err := os.OpenFile("rtt-wlan.csv", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer wlanCsv.Close()
	ranCsv, err := os.OpenFile("rtt-ran.csv", os.O_RDWR | os.O_CREATE | os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	defer ranCsv.Close()

	wlanCsvWriter := csv.NewWriter(wlanCsv)
	if err := wlanCsvWriter.Write([]string{"sequence", "times[us]", "RTT[us]"}); err != nil {
		panic(err)
	}
	ranCsvWriter := csv.NewWriter(ranCsv)
	if err := ranCsvWriter.Write([]string{"sequence", "times[us]", "RTT[us]"}); err != nil {
		panic(err)
	}

	localExp, remoteExp, err := ospfd.ConnectOspfd()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := localExp.Close(); err != nil {
			ocacaLogger.Printf("localExp.Close failed: %v", err)
		}
	}()
	defer func() {
		if err := remoteExp.Close(); err != nil {
			ocacaLogger.Printf("remoteExp.Close failed: %v", err)
		}
	}()
	
	wlanNic := Peer{remtoeAddress: "10.10.20.3", remoteIf: "wlangre", localIf: "wlangre", localAddress: "localhost", preffered: false}
	ranNic := Peer{remtoeAddress: "10.10.10.3", remoteIf: "rangre", localIf: "rangre", localAddress: "localhost", preffered: true}

	now = time.Now()
	for {
		wlanCH := make(chan time.Duration)
		ranCH := make(chan time.Duration)
		go measureRTT(wlanNic.remtoeAddress, wlanCH, wlanCsvWriter)
		go measureRTT(ranNic.remtoeAddress, ranCH, ranCsvWriter)

		wlanStats := <-wlanCH
		ranStats := <-ranCH

		if wlanStats < ranStats && !wlanNic.preffered {
			ospfd.CostSet(localExp, 1, wlanNic.localIf)
			ospfd.CostSet(localExp, 2, ranNic.localIf)
			ospfd.CostSet(remoteExp, 1, wlanNic.remoteIf)
			ospfd.CostSet(remoteExp, 2, ranNic.remoteIf)

			wlanNic.preffered = true
			ranNic.preffered = false
			ocacaLogger.Println("wlan became priority")
		}
		if ranStats < wlanStats && !ranNic.preffered {
			ospfd.CostSet(localExp, 2, wlanNic.localIf)
			ospfd.CostSet(localExp, 1, ranNic.localIf)
			ospfd.CostSet(remoteExp, 2, wlanNic.remoteIf)
			ospfd.CostSet(remoteExp, 1, ranNic.remoteIf)

			ranNic.preffered = true
			wlanNic.preffered = false
			ocacaLogger.Println("RAN became priority")
		}
	}
}

func measureRTT(addr string, ch chan time.Duration, csvWriter *csv.Writer) {
	pinger, err := ping.NewPinger(addr)
	pinger.SetPrivileged(true)
	pinger.Interval = 5 * time.Millisecond
	if err != nil {
		panic(err)
	}

	pinger.OnRecv = func(pkt *ping.Packet) {
		// TODO: sequence nubmerを数え上げ
		if err = csvWriter.Write([]string{
			strconv.Itoa(pkt.Seq),
			strconv.FormatInt(time.Since(now).Microseconds(), 10),
			strconv.FormatInt(pkt.Rtt.Microseconds(), 10),
		}); err != nil {
			panic(err)
		}
	}

	go pinger.Run()
	time.Sleep(15 * time.Millisecond)
	pinger.Stop()
	csvWriter.Flush()
	
	stats := pinger.Statistics()
	ch <- stats.AvgRtt
}
