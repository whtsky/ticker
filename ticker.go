package ticker

import (
	"context"
	"time"
)

type ImmediateTicker struct {
	ticker *time.Ticker
	C      <-chan time.Time
	ctx    context.Context
	cancel context.CancelFunc
}

func NewImmediateTickerWithContext(
	ctx context.Context, d time.Duration,
) *ImmediateTicker {
	ch := make(chan time.Time, 1)
	ctx, cancel := context.WithCancel(ctx)
	ticker := &ImmediateTicker{
		ticker: time.NewTicker(d),
		C:      ch,
		ctx:    ctx,
		cancel: cancel,
	}

	go func() {
		select {
		case <-ticker.ctx.Done():
			ticker.ticker.Stop()
			return
		default:
			ch <- time.Now()
		}
		for {
			select {
			case t := <-ticker.ticker.C:
				select {
				case ch <- t:
				case <-ticker.ctx.Done():
					ticker.ticker.Stop()
					return
				}
			case <-ticker.ctx.Done():
				ticker.ticker.Stop()
				return
			}
		}
	}()

	return ticker
}

func NewImmediateTicker(d time.Duration) *ImmediateTicker {
	return NewImmediateTickerWithContext(context.Background(), d)
}

func (t *ImmediateTicker) Stop() {
	t.cancel()
}
