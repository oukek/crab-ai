package chains

import (
	"sync"
	"time"
)

type Middleware[C any] func(handlerFunc HandlerFunc[C]) HandlerFunc[C]

type HandlerFunc[C any] func(ctx C) error

type Context[C any] struct {
	Index    int
	handlers []HandlerFunc[C]

	mu   sync.RWMutex
	Keys map[string]interface{}
}

func NewContext[C any]() *Context[C] {
	return &Context[C]{
		Index: -1,
	}
}

func (c *Context[C]) Reset(i int) {
	c.mu.RLock()
	c.Index = i
	c.mu.RUnlock()
}

func (c *Context[C]) Next(nc C) error {
	c.Index++
	for c.Index < len(c.handlers) {
		err := c.handlers[c.Index](nc)
		if err != nil {
			return err
		}
		c.Index++
	}
	return nil
}
func (c *Context[C]) Abort() {
	c.Index = len(c.handlers)
}

func (c *Context[C]) Use(handlers ...HandlerFunc[C]) {
	c.handlers = append(c.handlers, handlers...)
}

func (c *Context[C]) Get(key string) (value any, exists bool) {
	c.mu.RLock()
	value, exists = c.Keys[key]
	c.mu.RUnlock()
	return
}

func (c *Context[C]) GetWithDefault(key string, defaultValue any) any {
	if value, exists := c.Get(key); exists {
		return value
	}
	return defaultValue
}

func (c *Context[C]) Set(key string, value any) {

	c.mu.Lock()
	if c.Keys == nil {
		c.Keys = make(map[string]interface{})
	}

	c.Keys[key] = value
	c.mu.Unlock()
}

func (c *Context[C]) Deadline() (deadline time.Time, ok bool) {
	return
}

func (c *Context[C]) Done() <-chan struct{} {
	return nil
}

func (c *Context[C]) Err() error {
	return nil
}

func (c *Context[C]) Value(key interface{}) interface{} {
	if keyAsString, ok := key.(string); ok {
		val, _ := c.Get(keyAsString)
		return val
	}
	return nil
}
