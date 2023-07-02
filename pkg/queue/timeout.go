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

type Queue struct {
	timeout int
	q       []*item
	imap    map[string]*item
	size    int32
	m       sync.RWMutex
	r, w    int32
	len     int32
	notify  chan struct{}
}

func NewTimeoutQueue(timeout int, size int32) *Queue {
	return &Queue{
		timeout: timeout,
		size:    size,
		q:       make([]*item, size),
		imap:    make(map[string]*item),
		notify:  make(chan struct{}),
	}
}

func (q *Queue) Push(id string, fn func()) {
	q.m.Lock()
	defer q.m.Unlock()
	if q.len == q.size {
		q.grow()
	}
	i := &item{identifier: id, fn: fn, t: time.NewTimer(time.Duration(q.timeout) * time.Second)}
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

func (q *Queue) pop() *item {
	q.m.Lock()
	defer q.m.Unlock()
	i := q.q[q.r]
	delete(q.imap, i.identifier)
	q.r += 1
	if q.r == q.size {
		q.r = 0
	}
	q.len -= 1
	return i
}

func (q *Queue) Finish(identifier string) {
	q.m.RLock()
	defer q.m.RUnlock()
	i, ok := q.imap[identifier]
	if ok {
		i.finish.Store(true)
	}
}

func (q *Queue) grow() {
	tmp := make([]*item, len(q.q)*2)
	copy(tmp, q.q)
	q.q = tmp
	q.size = int32(len(q.q))
	q.w = q.len
}

func (q *Queue) Run(ctx context.Context) {
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
			i.fn()
		case <-ctx.Done():
			return
		}
	}
}
