package goping

import (
	"fmt"
	"net"
	"time"
)

func runTCP(dstIP net.IP, opts Options) (int, int, []time.Duration, error) {
	var (
		sent, recv int
		rtts       []time.Duration
	)

	for i := range make([]struct{}, opts.Count) {
		sent++
		start := time.Now()
		c, err := net.DialTimeout("tcp", dstIP.String()+":80", opts.Timeout)
		if err != nil {
			return sent, recv, rtts, err
		}

		rtt := time.Since(start)
		recv++
		rtts = append(rtts, rtt)
		fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n", opts.Size, dstIP.String(), i, rtt)
		_ = c.Close()
		time.Sleep(opts.Interval)
	}

	return sent, recv, rtts, nil
}
