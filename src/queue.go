package main

import (
	"sync"
	"time"
)

type Queue struct {
	sync.RWMutex

	channels map[string]chan string
}

func (q *Queue) Push(key, msg string) {
	ch := q.getChan(key)
	if ch == nil {
		ch = q.createChan(key)
	}

	go func() {
		ch <- msg
	}()
}

func (q *Queue) Pop(key string) string {
	ch := q.getChan(key)
	if ch == nil {
		return ""
	}

	select {
	case msg := <-ch:
		return msg
	default:
		return ""
	}
}

func (q *Queue) PopWait(key string, ttl time.Duration) string {
	ch := q.getChan(key)
	if ch == nil {
		ch = q.createChan(key)
	}

	select {
	case msg := <-ch:
		return msg
	case <-time.After(ttl * time.Second):
		return ""
	}
}

func (q *Queue) getChan(key string) chan string {
	q.RLock()
	ch := q.channels[key]
	q.RUnlock()

	return ch
}

func (q *Queue) createChan(key string) chan string {
	ch := make(chan string)
	q.Lock()
	q.channels[key] = ch
	q.Unlock()

	return ch
}
