package goping

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func runICMP(dstIP net.IP, opts Options) (int, int, []time.Duration, error) {
	ipv := 4
	icmpReqType := icmp.Type(ipv4.ICMPTypeEcho)
	icmpRepType := icmp.Type(ipv4.ICMPTypeEchoReply)
	netw := networkForOS(4)

	if dstIP.To4() == nil {
		ipv = 6
		icmpReqType = ipv6.ICMPTypeEchoRequest
		icmpRepType = ipv6.ICMPTypeEchoReply
		netw = networkForOS(6)
	}

	c, err := icmp.ListenPacket(netw, "")
	if err != nil {
		return 0, 0, nil, fmt.Errorf("listen error on %q: %w", netw, err)
	}

	defer c.Close()

	dst := &net.IPAddr{IP: dstIP}
	pid := os.Getpid() & 0xffff
	seq := 0

	var sent, recv int
	var rtts []time.Duration

	readBuf := make([]byte, 1500)

	for i := range make([]struct{}, opts.Count) {
		seq++
		size := opts.Size
		payload := make([]byte, size)

		if size < 8 {
			size = 8
			payload = make([]byte, size) // must hold timestamp
		}
		binary.BigEndian.PutUint64(payload[:8], uint64(time.Now().UnixNano()))

		msg := icmp.Message{
			Type: icmpReqType,
			Code: 0,
			Body: &icmp.Echo{
				ID:   pid,
				Seq:  seq,
				Data: payload,
			},
		}
		b, _ := msg.Marshal(nil)

		start := time.Now()
		if _, err := c.WriteTo(b, dst); err != nil {
			continue
		}
		sent++

		_ = c.SetReadDeadline(time.Now().Add(opts.Timeout))

		for {
			n, _, err := c.ReadFrom(readBuf)
			if err != nil {
				break
			}
			rtt := time.Since(start)
			fmt.Printf("%d bytes from %s: icmp_seq=%d time=%v\n", size, dstIP.String(), i, rtt)

			// https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
			var proto int
			if ipv == 4 {
				proto = 1
			} else {
				proto = 58
			}

			rm, err := icmp.ParseMessage(proto, readBuf[:n])
			if err != nil {
				continue
			}
			if rb, ok := rm.Body.(*icmp.Echo); ok {
				if rm.Type == icmpRepType && rb.ID == pid && rb.Seq == seq {
					if len(rb.Data) >= 8 {
						ts := time.Unix(0, int64(binary.BigEndian.Uint64(rb.Data[:8])))
						rtt = time.Since(ts)
					}
					recv++
					rtts = append(rtts, rtt)
					break
				}
			}
			if time.Since(start) > opts.Timeout {
				break
			}
		}
		if i != opts.Count-1 {
			time.Sleep(opts.Interval)
		}
	}
	return sent, recv, rtts, nil
}

func networkForOS(ipVersion int) string {
	if ipVersion == 4 {
		return "ip4:icmp"
	}
	return "ip6:ipv6-icmp"
}
