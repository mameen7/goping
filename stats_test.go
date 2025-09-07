package goping

import (
	"testing"
	"time"
)

func TestComputeStats_NoPackets(t *testing.T) {
	stats := computeStats(0, 0, nil)
	if stats.Sent != 0 || stats.Recv != 0 {
		t.Errorf("expected 0 sent/recv, got %d/%d", stats.Sent, stats.Recv)
	}
	if stats.Loss != 0.0 {
		t.Errorf("expected 0 loss, got %.2f", stats.Loss)
	}
	if stats.AvgRTT != 0 {
		t.Errorf("expected avg=0, got %v", stats.AvgRTT)
	}
}

func TestComputeStats_AllReceived(t *testing.T) {
	rtts := []time.Duration{
		10 * time.Millisecond,
		20 * time.Millisecond,
		30 * time.Millisecond,
	}

	stats := computeStats(5, 5, rtts)

	if stats.Sent != 5 || stats.Recv != 5 {
		t.Errorf("expected 5 sent/recv, got %d/%d", stats.Sent, stats.Recv)
	}
	if stats.MinRTT != 10*time.Millisecond {
		t.Errorf("expected min=10ms, got %v", stats.MinRTT)
	}
	if stats.MaxRTT != 30*time.Millisecond {
		t.Errorf("expected max=30ms, got %v", stats.MaxRTT)
	}
	if stats.AvgRTT != 20*time.Millisecond {
		t.Errorf("expected avg=20ms, got %v", stats.AvgRTT)
	}
	if stats.Loss != 0.0 {
		t.Errorf("expected loss=0%%, got %.1f", stats.Loss)
	}
}

func TestComputeStats_PacketLoss(t *testing.T) {
	rtts := []time.Duration{
		5 * time.Millisecond,
	}

	stats := computeStats(5, 3, rtts)

	if stats.Sent != 5 || stats.Recv != 3 {
		t.Errorf("expected 5 sent, 3 recv, got %d/%d", stats.Sent, stats.Recv)
	}
	if stats.MinRTT != 5*time.Millisecond || stats.MaxRTT != 5*time.Millisecond {
		t.Errorf("expected min=max=5ms, got %v/%v", stats.MinRTT, stats.MaxRTT)
	}
	if stats.AvgRTT != 5*time.Millisecond {
		t.Errorf("expected avg=5ms, got %v", stats.AvgRTT)
	}
	if stats.Loss != 40 && stats.Loss != 66.6666666667 {
		t.Errorf("expected ~40%% loss, got %.2f", stats.Loss)
	}
}
