package queue

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type item struct {
	identifier string
	fn         func()
	finish     atomic.Bool
	t          *time.Timer
}

type TimeQueue struct {
	timeout int
	q       []*item
	imap    map[string]*item
	size    int32
	m       sync.RWMutex
	r, w    int32
	len     int32
	notify  chan struct{}
}

func NewTimeoutQueue(timeout int, size int32) *TimeQueue {
	return &TimeQueue{
		timeout: timeout,
		size:    size,
		q:       make([]*item, size),
		imap:    make(map[string]*item),
		notify:  make(chan struct{}),
	}
}

func (q *TimeQueue) Push(id string, fn func()) {
	q.m.Lock()
	defer q.m.Unlock()
	if q.len == q.size {
		q.grow()
	}
	i := &item{identifier: id, fn: fn, t: time.NewTimer(time.Duration(q.timeout) * time.Millisecond)}
	q.q[q.w] = i
	q.imap[id] = i
	q.w += 1
	if q.w == q.size {
		q.w = 0
	}
	q.len += 1
	select {
	case q.notify <- struct{}{}:
	default:
	}
}

func (q *TimeQueue) pop() *item {
	q.m.Lock()
	defer q.m.Unlock()
	i := q.q[q.r]
	q.r += 1
	if q.r == q.size {
		q.r = 0
	}
	q.len -= 1
	return i
}

func (q *TimeQueue) Finish(identifier string) {
	q.m.RLock()
	i, ok := q.imap[identifier]
	q.m.RUnlock()
	if ok {
		q.m.Lock()
		delete(q.imap, i.identifier)
		q.m.Unlock()
		i.finish.Store(true)
	}
}

func (q *TimeQueue) IsFinished(identifier string) bool {
	q.m.RLock()
	defer q.m.RUnlock()
	i, ok := q.imap[identifier]
	if ok {
		return i.finish.Load()
	}
	return true
}

func (q *TimeQueue) grow() {
	tmp := make([]*item, len(q.q)*2)
	copy(tmp, q.q)
	q.q = tmp
	q.size = int32(len(q.q))
	q.w = q.len
}

func (q *TimeQueue) Run(ctx context.Context) {
	for {
		if q.len == 0 {
			<-q.notify
		}
		i := q.pop()
		if i.finish.Load() {
			continue
		}
		select {
		case <-i.t.C:
			go i.fn()
		case <-ctx.Done():
			return
		}
	}
}
