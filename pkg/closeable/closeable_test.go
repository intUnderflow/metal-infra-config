package closeable

import (
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
	"time"
)

func Test_Close_ClosesCloseable(t *testing.T) {
	closeable := NewCloseable()
	closeable.Close()
	closeable.WaitForClosure()
}

func Test_Close_CanBeCalledMultipleTimes(t *testing.T) {
	closeable := NewCloseable()
	closeable.Close()
	closeable.Close()
	closeable.WaitForClosure()
}

func Test_WaitForClosure_BlocksIfNotClosed(t *testing.T) {
	closeable := NewCloseable()
	waited := &sync.Mutex{}
	waited.Lock()
	go func() {
		closeable.WaitForClosure()
		waited.Unlock()
	}()
	time.Sleep(time.Millisecond)
	require.False(t, waited.TryLock())
	closeable.Close()
	waited.Lock()
}

func Test_WaitForClosure_BlocksMultipleCallersIfNotClosed(t *testing.T) {
	closeable := NewCloseable()
	waited1 := &sync.Mutex{}
	waited1.Lock()
	waited2 := &sync.Mutex{}
	waited2.Lock()
	go func() {
		closeable.WaitForClosure()
		waited1.Unlock()
	}()
	go func() {
		closeable.WaitForClosure()
		waited2.Unlock()
	}()
	time.Sleep(time.Millisecond)
	require.False(t, waited1.TryLock())
	require.False(t, waited2.TryLock())
	closeable.Close()
	waited1.Lock()
	waited2.Lock()
}

func Test_IsClosed_ReportsClosureStatus(t *testing.T) {
	closeable := NewCloseable()
	require.False(t, closeable.IsClosed())
	closeable.Close()
	require.True(t, closeable.IsClosed())
}
