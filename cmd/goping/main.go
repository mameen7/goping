package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/mameen7/goping"
)

func main() {
	var count int
	var interval time.Duration
	var timeout time.Duration
	var size int

	flag.IntVar(&count, "c", 5, "number of echo request to send")
	flag.DurationVar(&interval, "i", time.Second, "interval between sends")
	flag.DurationVar(&timeout, "t", 10*time.Second, "per packet timeout")
	flag.IntVar(&size, "s", 64, "echo packet size")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Fprintf(os.Stderr, "usage: %s [options] host\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}

	host := flag.Arg(0)
	p := goping.NewPinger(
		host,
		goping.WithCount(count),
		goping.WithInterval(interval),
		goping.WithTimeout(timeout),
		goping.WithSize(size),
	)
	stats, err := p.Ping()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ping error: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("\n--- %s ping statistics ---\n", host)
	fmt.Printf("%d packets transmitted, %d received, %.1f%% packet loss\n", stats.Sent, stats.Recv, stats.Loss)
	if stats.Recv > 0 {
		fmt.Printf("RTT MIN/AVG/MAX = %v/%v/%v\n",
			stats.MinRTT, stats.AvgRTT, stats.MaxRTT)
	}
}
