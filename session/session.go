package session

import (
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/Workiva/go-datastructures/trie/ctrie"
)

// node
type node struct {
	latestAccess int64
	value        interface{}
}

// nodePool is a pool of node
var nodePool = sync.Pool{
	New: func() interface{} {
		return &node{}
	},
}

// Local is a session of app scope
type Local struct {
	data    ctrie.Ctrie
	ttl     int64
	count   int64
	prevGC  int64
	counter chan struct{}
}

// NewLocal create new Local instance
func NewLocal() *Local {
	return &Local{
		data:    *ctrie.New(nil),
		counter: make(chan struct{}, 1),
	}
}

// IncreaseCount is increasing count of local session
func (l *Local) IncreaseCount() {
	atomic.AddInt64(&l.count, 1)
	v := atomic.LoadInt64(&l.count)
	p := atomic.LoadInt64(&l.prevGC)
	if p*3/2 <= v {
		l.counter <- struct{}{}
	}
}

// DecreaseCount is decreasing count of local session
func (l *Local) DecreaseCount() {
	atomic.AddInt64(&l.count, -1)
}

// StartGC is starting GC process
func StartGC(localSession *Local) {
	gc := func() {
		cancel := make(chan struct{}, 1)
		iter := localSession.data.Iterator(cancel)
		for p := range iter {
			n, ok := p.Value.(*node)
			if !ok {
				continue
			}
			if n.latestAccess+localSession.ttl < time.Now().Unix() {
				localSession.DecreaseCount()
				localSession.data.Remove(p.Key)
				nodePool.Put(n)
			}
		}
		atomic.StoreInt64(&localSession.prevGC, atomic.LoadInt64(&localSession.count))
	}
	go func() {
		interval := time.NewTicker(time.Duration(localSession.ttl) * 125 / 100)
		for {
			select {
			case <-localSession.counter:
				gc()
			case <-interval.C:
				gc()
			}
		}
	}()
}

// SetTTL set Local's ttl to given duration
func (l *Local) SetTTL(duration time.Duration) {
	l.ttl = duration.Milliseconds()
}

// errNotExists is Not Exists Error value
var errNotExists = errors.New("not exists")

// IsNotExists is a function to check given error is Not Exists Error
func IsNotExists(err error) bool {
	return errors.Is(err, errNotExists)
}

// errNotMatchType is Not Match Type Error value
var errNotMatchType = errors.New("not match type")

// IsNotMatchType is a function to check given error is Not Match Type Error
func IsNotMatchType(err error) bool {
	return errors.Is(err, errNotMatchType)
}

// errTimeout is Time Out Error value
var errTimeout = errors.New("time out")

// IsTimeout is a function to check given error is Time Out Error
func IsTimeout(err error) bool {
	return errors.Is(err, errTimeout)
}

// GetLocal is getting the value of given key.
// if not exists, return Not Exists Error.
// if not match type with generic, return Not Match Type Error.
func GetLocal[T any](localSession *Local, key []byte) (*T, error) {
	data, ok := localSession.data.Lookup(key)
	if !ok {
		return nil, errNotExists
	}
	n, ok := data.(*node)
	if !ok {
		return nil, errNotMatchType
	}
	if n.latestAccess+localSession.ttl >= time.Now().Unix() {
		localSession.DecreaseCount()
		localSession.data.Remove(key)
		nodePool.Put(n)
		return nil, errTimeout
	}
	t, ok := n.value.(*T)
	if !ok {
		return nil, errNotMatchType
	}
	n.latestAccess = time.Now().Unix()
	localSession.IncreaseCount()
	localSession.data.Insert(key, n)
	return t, nil
}

// errAlreadyExists is Already Exists Error value
var errAlreadyExists = errors.New("already exists")

// IsAlreadyExists is a function to check given error is Already Exists Error
func IsAlreadyExists(err error) bool {
	return errors.Is(err, errAlreadyExists)
}

// SetLocal is setting value of key into local session.
// if already exists key value pair in local session, return Already Exists Error.
func SetLocal[T any](localSession *Local, key []byte, value T) error {
	if d, ok := localSession.data.Lookup(key); ok {
		n, ok := d.(*node)
		if !ok {
			localSession.DecreaseCount()
			localSession.data.Remove(key)
		}
		if n.latestAccess+localSession.ttl < time.Now().Unix() {
			return errAlreadyExists
		}
		localSession.DecreaseCount()
		localSession.data.Remove(key)
		nodePool.Put(n)
	}
	n := nodePool.Get().(*node)
	n.latestAccess = time.Now().Unix()
	n.value = &value
	localSession.IncreaseCount()
	localSession.data.Insert(key, n)
	return nil
}

// RemoveLocal is deleting key value pair in local session.
// if not exists, return Not Exists Error.
// if not match type with generic, return Not Match Type Error.
func RemoveLocal[T any](localSession *Local, key []byte) (*T, error) {
	data, ok := localSession.data.Remove(key)
	if !ok {
		return nil, errNotExists
	}
	n, ok := data.(*node)
	if !ok {
		return nil, errNotMatchType
	}
	t, ok := n.value.(*T)
	if !ok {
		return nil, errNotMatchType
	}
	localSession.DecreaseCount()
	nodePool.Put(n)
	return t, nil
}
