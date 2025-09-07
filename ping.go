package goping

import (
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type Options struct {
	Count    int
	Interval time.Duration
	Timeout  time.Duration
	Size     int
}

type Stats struct {
	Sent   int
	Recv   int
	MinRTT time.Duration
	AvgRTT time.Duration
	MaxRTT time.Duration
	Loss   float64
	RTTs   []time.Duration
}

type Pinger struct {
	Host    string
	Options Options
}

type PingerOption func(*Options)

func WithCount(c int) PingerOption {
	return func(o *Options) { o.Count = c }
}
func WithInterval(i time.Duration) PingerOption {
	return func(o *Options) { o.Interval = i }
}
func WithTimeout(t time.Duration) PingerOption {
	return func(o *Options) { o.Timeout = t }
}
func WithSize(s int) PingerOption {
	return func(o *Options) { o.Size = s }
}

func NewPinger(host string, setters ...PingerOption) *Pinger {
	opts := Options{
		Count:   5,
		Timeout: 3 * time.Second,
		Size:    64,
	}
	for _, set := range setters {
		set(&opts)
	}
	return &Pinger{
		Host:    host,
		Options: opts,
	}
}

func (p *Pinger) Ping() (*Stats, error) {
	var (
		sent, recv int
		rtts       []time.Duration
	)

	dstIP, err := resolveIP(p.Host)
	if err != nil {
		return nil, fmt.Errorf("resolve error: %w", err)
	}

	if hasICMPPrivilege() {
		sent, recv, rtts, err = runICMP(dstIP, p.Options)
		if err != nil {
			return nil, fmt.Errorf("ICMP ping failed with error: %w", err)
		}
	} else {
		fmt.Println("ICMP not permitted, falling back to TCP")
		sent, recv, rtts, err = runTCP(dstIP, p.Options)
		if err != nil {
			return nil, fmt.Errorf("TCP ping failed with error: %w", err)
		}
	}

	stats := computeStats(sent, recv, rtts)
	return &stats, nil
}

func resolveIP(host string) (net.IP, error) {
	ip := net.ParseIP(host)
	if ip != nil {
		return ip, nil
	}

	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, fmt.Errorf("could not resolve hostname %q: %v", host, err)
	}
	if len(ips) == 0 {
		return nil, fmt.Errorf("no IPs for %s", host)
	}

	return ips[0], nil
}

func hasICMPPrivilege() bool {
	c, err := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		if strings.Contains(err.Error(), "operation not permitted") {
			return false
		}
		// Some other ICMP error → still fail
		fmt.Printf("Failed to open ICMP socket: %v\n", err)
		os.Exit(1)
	}
	defer c.Close()
	return true
}
