package bean

import "sync"

type Container struct {
	beans map[any]any
	lock  sync.RWMutex
}

func NewContainer() *Container {
	return &Container{
		beans: map[any]any{},
	}
}

func (c *Container) Get(key any) (value any, ok bool) {
	c.lock.RLock()
	value, ok = c.beans[key]
	c.lock.RUnlock()
	return
}

func (c *Container) Set(key any, value any) {
	c.lock.Lock()
	c.beans[key] = value
	c.lock.Unlock()
}

func (c *Container) Delete(key any) {
	c.lock.Lock()
	delete(c.beans, key)
	c.lock.Unlock()
}
