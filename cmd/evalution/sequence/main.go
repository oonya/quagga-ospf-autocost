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
	timer := time.After(3 * time.Second)

	go execIperf3(done)
	evalLogger.Println("===iperf3===")

	go func() {
		args := []string{"netns", "exec", "server", "/usr/bin/make", "evalRTT"}
		if err := exec.Command("/usr/sbin/ip", args...).Run(); err != nil {
			// TODO: wrap error
			panic(err)
		}
	}()

	for {
		select {
		case <- done:
			args := []string{"qdisc", "change", "dev", "rangre", "root", "handle", "1:0", "netem", "delay", "5ms"}
			if err := exec.Command("/usr/sbin/tc", args...).Run(); err != nil {
				// TODO: wrap error
				panic(err)
			}
			time.Sleep(1 * time.Second)
			return
		case <-timer:
			args := []string{"qdisc", "change", "dev", "rangre", "root", "handle", "1:0", "netem", "delay", "1ms"}
			if err := exec.Command("/usr/sbin/tc", args...).Run(); err != nil {
				// TODO: wrap error
				panic(err)
			}
			evalLogger.Println("exec tc cmd")
		}
	}
}

func execIperf3(done chan struct{}) {
	evalLogger.Println("iperf3 start")
	args := []string{"-c192.168.3.4", "-b80M", "-i1", "-t6", "-B192.168.1.1"}
	if output, err := exec.Command("/usr/bin/iperf3", args...).CombinedOutput(); err != nil {
		// TODO: wrap error
		panic(err)
	} else {
		evalLogger.Printf("%s", output)
	}
	evalLogger.Println("iperf3 done")
	close(done)
}
