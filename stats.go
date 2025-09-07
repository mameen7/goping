package goping

import "time"

func computeStats(sent, recv int, rtts []time.Duration) Stats {
	if recv == 0 {
		loss := 0.0
		if sent > 0 {
			loss = 100.0 // all lost
		}
		return Stats{
			Sent: sent,
			Recv: recv,
			Loss: loss,
			RTTs: rtts,
		}
	}

	var min, max, sum time.Duration
	for i, rtt := range rtts {
		if i == 0 || rtt < min {
			min = rtt
		}
		if rtt > max {
			max = rtt
		}
		sum += rtt
	}

	return Stats{
		Sent:   sent,
		Recv:   recv,
		MinRTT: min,
		AvgRTT: sum / time.Duration(len(rtts)),
		MaxRTT: max,
		Loss:   100.0 * float64(sent-recv) / float64(sent),
		RTTs:   rtts,
	}
}
