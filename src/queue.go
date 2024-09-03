package main

import (
	"sync"
)

var queue = Queue{
	messages:    make(map[string][]string),
	subscribers: make(map[string][]chan string),
}

type Queue struct {
	sync.Mutex

	messages    map[string][]string
	subscribers map[string][]chan string
}

func (q *Queue) Push(qName, msg string) {
	q.Lock()
	defer q.Unlock()

	if len(q.subscribers[qName]) != 0 {
		defer close(q.subscribers[qName][0])
		q.subscribers[qName][0] <- msg
		q.subscribers[qName] = q.subscribers[qName][1:]
	} else {
		q.messages[qName] = append(q.messages[qName], msg)
	}
}

func (q *Queue) Pop(qName string) string {
	q.Lock()
	defer q.Unlock()

	messages, ok := q.messages[qName]
	if !ok || len(messages) == 0 {
		return ""
	}
	q.messages[qName] = messages[1:]

	return messages[0]
}

func (q *Queue) PopWait(qName string) chan string {
	q.Lock()

	messages, ok := q.messages[qName]
	if ok && len(messages) != 0 {
		ch := make(chan string)
		go func() {
			ch <- messages[0]
		}()
		q.messages[qName] = messages[1:]
		q.Unlock()

		return ch
	}

	ch := make(chan string)
	q.subscribers[qName] = append(q.subscribers[qName], ch)
	q.Unlock()

	return ch
}

func (q *Queue) RemoveSub(qName string, ch chan string) {
	q.Lock()
	defer q.Unlock()
	defer close(ch)

	oldSubs := q.subscribers[qName]
	newSubs := make([]chan string, 0)

	for _, el := range oldSubs {
		if el != ch {
			newSubs = append(newSubs, el)
		}
	}

	q.subscribers[qName] = newSubs
}
