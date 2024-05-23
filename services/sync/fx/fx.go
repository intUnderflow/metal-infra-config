package fx

import (
	"context"
	"errors"
	"github.com/intunderflow/metal-infra-config/pkg/closeable"
	servicesync "github.com/intunderflow/metal-infra-config/services/sync"
	"sync"
	"time"
)

const (
	_errAlreadyStoppedOrNeverStarted = "SyncFX was not started or was stopped already"
)

type SyncFX struct {
	mutex   *sync.Mutex
	sync    servicesync.Sync
	control closeable.Closeable
}

func NewFX(syncer servicesync.Sync) *SyncFX {
	return &SyncFX{
		mutex:   &sync.Mutex{},
		sync:    syncer,
		control: nil,
	}
}

func (s *SyncFX) Start(ctx context.Context, interval time.Duration) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.control = s.sync.SyncPeriodically(ctx, interval)
}

func (s *SyncFX) Stop() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if s.control == nil {
		return errors.New(_errAlreadyStoppedOrNeverStarted)
	}
	s.control.Close()
	s.control = nil
	return nil
}
