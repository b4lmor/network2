package timer

import (
	"time"
)

type Timer struct {
	subscribers []chan byte
	Delay       time.Duration
}

func NewTimer(delay time.Duration) *Timer {
	return &Timer{Delay: delay, subscribers: make([]chan byte, 0)}
}

func (t *Timer) Subscribe() chan byte {
	ch := make(chan byte)
	t.subscribers = append(t.subscribers, ch)
	return ch
}

func (t *Timer) Start() {
	go func() {
		for {
			time.Sleep(t.Delay)
			t.notifySubscribers()
		}
	}()
}

func (t *Timer) notifySubscribers() {
	for _, subscriber := range t.subscribers {
		subscriber <- 0x01
	}
}
