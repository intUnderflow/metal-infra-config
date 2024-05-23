package closeable

import "sync/atomic"

type Closeable interface {
	Close()
	WaitForClosure()
	IsClosed() bool
}

type closeableImpl struct {
	closed  *atomic.Bool
	channel chan interface{}
}

func NewCloseable() Closeable {
	return &closeableImpl{
		closed:  &atomic.Bool{},
		channel: make(chan interface{}),
	}
}

func (c *closeableImpl) Close() {
	wasCalledBefore := c.closed.Swap(true)
	if !wasCalledBefore {
		close(c.channel)
	}
}

func (c *closeableImpl) WaitForClosure() {
	<-c.channel
}

func (c *closeableImpl) IsClosed() bool {
	return c.closed.Load()
}
