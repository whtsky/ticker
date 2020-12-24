package ticker

import (
	"context"
	"testing"
	"time"

	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestTick(t *testing.T) {
	const tickCount = 10
	delta := 200 * time.Millisecond
	maxDuration := delta * 2
	ticker := NewImmediateTicker(delta)
	// immediate one
	<-ticker.C
	lastTick := time.Now()
	for i := 0; i < tickCount; i++ {
		<-ticker.C
		duration := time.Since(lastTick)
		if duration > maxDuration {
			t.Fatalf(
				"tick took %s, expected < %s",
				duration, maxDuration,
			)
		}
		lastTick = time.Now()
	}
	ticker.Stop()

	select {
	case <-ticker.C:
		t.Fatal("tick did not stop")
	default:
	}
}

func TestStopTicker(t *testing.T) {
	duration := 100 * time.Millisecond
	ticker := NewImmediateTicker(duration)
	<-ticker.C
	ticker.Stop()
	select {
	case <-ticker.C:
		t.Fatal("tick did not stop")
	case <-time.After(duration):
		return
	}
}

func TestTickerWithStoppedContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ticker := NewImmediateTickerWithContext(ctx, time.Microsecond)
	select {
	case <-ticker.C:
		t.Fatal("tick started")
	default:
	}
}
