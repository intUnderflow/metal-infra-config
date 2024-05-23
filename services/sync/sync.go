package sync

import (
	"context"
	"errors"
	"github.com/intunderflow/metal-infra-config/entities"
	"github.com/intunderflow/metal-infra-config/pkg/closeable"
	"github.com/intunderflow/metal-infra-config/services/peerdiscovery"
	"go.uber.org/zap"
	"time"
)

const (
	_errFailedToSyncAnyPeers = "failed to sync to any peers"
)

type Sync interface {
	Sync(context.Context) error
	SyncPeriodically(context.Context, time.Duration) closeable.Closeable
}

type syncImpl struct {
	log           *zap.Logger
	config        entities.Config
	peerDiscovery peerdiscovery.PeerDiscovery
	rpc           RPC
}

func NewSync(
	log *zap.Logger,
	config entities.Config,
	peerDiscovery peerdiscovery.PeerDiscovery,
	rpc RPC,
) Sync {
	return &syncImpl{
		log:           log,
		config:        config,
		peerDiscovery: peerDiscovery,
		rpc:           rpc,
	}
}

// Sync contacts peers we are aware of and synchronises configurations from that peer.
// It writes any new values we are not aware of to our config object.
// Sync calls each peer one-by-one, not all at once and uses a bidirectional stream to copy each others configs.
func (s *syncImpl) Sync(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	peers, err := s.peerDiscovery.GetPeers()
	if err != nil {
		return err
	}

	syncedAtLeastOnePeer := false
	for _, peer := range peers {
		client, err := s.rpc.GetClient(peer)
		if err != nil {
			reportSyncFailure(s.log, peer, err)
			continue
		}

		syncSession, err := client.Sync(ctx)
		if err != nil {
			reportSyncFailure(s.log, peer, err)
			continue
		}

		err = s.config.Sync(syncSession)
		if err != nil {
			reportSyncFailure(s.log, peer, err)
			continue
		}

		syncedAtLeastOnePeer = true
	}

	if len(peers) > 0 && !syncedAtLeastOnePeer {
		return errors.New(_errFailedToSyncAnyPeers)
	}

	return nil
}

func (s *syncImpl) SyncPeriodically(ctx context.Context, interval time.Duration) closeable.Closeable {
	closer := closeable.NewCloseable()
	go s.syncPeriodically(ctx, closer, interval)
	return closer
}

func (s *syncImpl) syncPeriodically(ctx context.Context, closer closeable.Closeable, interval time.Duration) {
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		closer.WaitForClosure()
		cancel()
	}()
	for !closer.IsClosed() {
		err := s.Sync(ctx)
		if err != nil {
			s.log.Warn("failed to sync", zap.Error(err))
		}
		t := time.NewTimer(interval)
		select {
		case <-ctx.Done():
			break
		case <-t.C:
			continue
		}
	}
}

func reportSyncFailure(log *zap.Logger, peer entities.Peer, err error) {
	log.Warn("Failed to sync with peer",
		zap.String("peer", string(peer.ID())+"@"+peer.Address()),
		zap.Error(err),
	)
}
