package time

import (
	"context"
	"time"
)

// Timer with non-blocking reset
type Timer struct {
	C      chan time.Time
	d      time.Duration
	f      func()
	timer  *time.Timer
	cancel context.CancelFunc
}

func NewTimer(d time.Duration) *Timer {
	t := &Timer{
		C:      make(chan time.Time, 1),
		d:      d,
		timer:  time.NewTimer(d),
		cancel: func() {}, // mock
	}
	t.start()
	return t
}

func AfterFunc(d time.Duration, f func()) *Timer {
	t := &Timer{
		d:      d,
		f:      f,
		timer:  time.NewTimer(d),
		cancel: func() {}, // mock
	}
	t.start()
	return t
}

func (t *Timer) start() {
	ctx, cancel := context.WithCancel(context.Background())
	t.cancel = cancel
	go t.run(ctx)
}

func (t *Timer) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case tt := <-t.timer.C:
			if t.C != nil {
				t.C <- tt
			}
			if t.f != nil {
				t.f()
			}
		}
	}
}

func (t *Timer) Reset() {
	t.Stop()
	t.timer = time.NewTimer(t.d)
	go t.start()
}

func (t *Timer) Stop() {
	t.cancel()
	t.timer.Stop()
}
