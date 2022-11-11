package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

var evalLogger *log.Logger

func main() {	
	file, err := os.OpenFile("evalution.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	evalLogger = log.New(file, "", log.Lmicroseconds)

	done := make(chan struct{})
	timer := time.After(10 * time.Second)

	go func ()  {
		evalLogger.Println("start")
		args := []string{"-c192.168.1.1", "-b1M", "-i1", "-t20"}
		if err := exec.Command("/usr/bin/iperf3", args...).Run(); err != nil {
			// TODO: wrap error
			panic(err)
		}
		evalLogger.Println("done")
		close(done)
	}()

	for {
		select {
		case <- done:
			args := []string{"qdisc", "del", "dev", "wlan0", "root"}
			if err := exec.Command("/usr/sbin/tc", args...).Run(); err != nil {
				// TODO: wrap error
				panic(err)
			}
			return
		case <-timer:
			args := []string{"qdisc", "add", "dev", "wlan0", "root", "handle", "1:0", "netem", "delay", "100ms"}
			if err := exec.Command("/usr/sbin/tc", args...).Run(); err != nil {
				// TODO: wrap error
				panic(err)
			}
			evalLogger.Println("exec tc cmd")
		}
	}
}
