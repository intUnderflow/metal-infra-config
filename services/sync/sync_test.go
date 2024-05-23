package sync

import (
	"context"
	"errors"
	"github.com/intunderflow/metal-infra-config/entities"
	mock_entities "github.com/intunderflow/metal-infra-config/entities/mock"
	"github.com/intunderflow/metal-infra-config/proto"
	mock_proto "github.com/intunderflow/metal-infra-config/proto/mock"
	mock_peerdiscovery "github.com/intunderflow/metal-infra-config/services/peerdiscovery/mock"
	mock_sync "github.com/intunderflow/metal-infra-config/services/sync/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap/zaptest"
	"sync"
	"testing"
	"time"
)

type syncTestFixtures struct {
	ctrl          *gomock.Controller
	sync          Sync
	config        *mock_entities.MockConfig
	peerDiscovery *mock_peerdiscovery.MockPeerDiscovery
	rpc           *mock_sync.MockRPC
}

func newTestFixtures(t *testing.T) *syncTestFixtures {
	ctrl := gomock.NewController(t)
	config := mock_entities.NewMockConfig(ctrl)
	peerDiscovery := mock_peerdiscovery.NewMockPeerDiscovery(ctrl)
	rpc := mock_sync.NewMockRPC(ctrl)
	sync := NewSync(
		zaptest.NewLogger(t),
		config,
		peerDiscovery,
		rpc,
	)
	return &syncTestFixtures{
		ctrl:          ctrl,
		sync:          sync,
		config:        config,
		peerDiscovery: peerDiscovery,
		rpc:           rpc,
	}
}

func Test_Sync_WhenGetPeersFails_ErrorIsReturned(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f := newTestFixtures(t)
	expectedErr := errors.New("foobar")
	f.peerDiscovery.EXPECT().GetPeers().Return(nil, expectedErr)

	err := f.sync.Sync(ctx)
	require.ErrorContains(t, err, expectedErr.Error())
}

func Test_Sync_WhenGetClientFails_ErrorIsReturned(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f := newTestFixtures(t)
	expectedErr := errors.New("foobar")
	peer := entities.NewPeer("foobar", "1.2.3.4:5678", time.Now())
	f.peerDiscovery.EXPECT().GetPeers().Return([]entities.Peer{
		peer,
	}, nil)
	f.rpc.EXPECT().GetClient(peer).Return(nil, expectedErr)

	err := f.sync.Sync(ctx)
	require.EqualError(t, err, _errFailedToSyncAnyPeers)
}

func Test_Sync_WhenStartSyncSessionFails_ErrorIsReturned(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f := newTestFixtures(t)
	expectedErr := errors.New("foobar")
	peer := entities.NewPeer("foobar", "1.2.3.4:5678", time.Now())
	f.peerDiscovery.EXPECT().GetPeers().Return([]entities.Peer{
		peer,
	}, nil)
	mockClient := mock_proto.NewMockMetalInfraConfigClient(f.ctrl)
	f.rpc.EXPECT().GetClient(peer).Return(mockClient, nil)
	mockClient.EXPECT().Sync(gomock.Any(), gomock.Any()).Return(nil, expectedErr)

	err := f.sync.Sync(ctx)
	require.EqualError(t, err, _errFailedToSyncAnyPeers)
}

func Test_Sync_WhenConfigSyncFails_ErrorIsReturned(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f := newTestFixtures(t)
	expectedErr := errors.New("foobar")
	peer := entities.NewPeer("foobar", "1.2.3.4:5678", time.Now())
	f.peerDiscovery.EXPECT().GetPeers().Return([]entities.Peer{
		peer,
	}, nil)
	mockClient := mock_proto.NewMockMetalInfraConfigClient(f.ctrl)
	f.rpc.EXPECT().GetClient(peer).Return(mockClient, nil)
	syncSession := mock_entities.NewMockSyncSession(f.ctrl)
	mockClient.EXPECT().Sync(gomock.Any(), gomock.Any()).Return(syncSession, nil)
	f.config.EXPECT().Sync(syncSession).Return(expectedErr)

	err := f.sync.Sync(ctx)
	require.EqualError(t, err, _errFailedToSyncAnyPeers)
}

func Test_Sync_WhenAtLeastOneClientSucceeds_NoErrorIsReturned(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f := newTestFixtures(t)
	peer1 := entities.NewPeer("foobar", "1.2.3.4:5678", time.Now())
	peer2 := entities.NewPeer("barbaz", "5.6.7.8:910", time.Now())
	f.peerDiscovery.EXPECT().GetPeers().Return([]entities.Peer{
		peer1, peer2,
	}, nil)
	mockClientPeer1 := mock_proto.NewMockMetalInfraConfigClient(f.ctrl)
	mockClientPeer1.EXPECT().Sync(gomock.Any(), gomock.Any()).Return(nil, errors.New("foo"))
	mockClientPeer2 := mock_proto.NewMockMetalInfraConfigClient(f.ctrl)
	syncSession := mock_proto.NewMockMetalInfraConfig_SyncClient(f.ctrl)
	mockClientPeer2.EXPECT().Sync(gomock.Any(), gomock.Any()).Return(syncSession, nil)
	f.rpc.EXPECT().GetClient(peer1).Return(mockClientPeer1, nil)
	f.rpc.EXPECT().GetClient(peer2).Return(mockClientPeer2, nil)
	f.config.EXPECT().Sync(syncSession).Return(nil)

	err := f.sync.Sync(ctx)
	require.NoError(t, err)
}

func Test_Sync_WhenAllClientsSucceed_NoErrorIsReturned(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f := newTestFixtures(t)
	peer1 := entities.NewPeer("foobar", "1.2.3.4:5678", time.Now())
	peer2 := entities.NewPeer("barbaz", "5.6.7.8:910", time.Now())
	f.peerDiscovery.EXPECT().GetPeers().Return([]entities.Peer{
		peer1, peer2,
	}, nil)
	mockClientPeer1 := mock_proto.NewMockMetalInfraConfigClient(f.ctrl)
	syncSessionPeer1 := mock_proto.NewMockMetalInfraConfig_SyncClient(f.ctrl)
	mockClientPeer1.EXPECT().Sync(gomock.Any(), gomock.Any()).Return(syncSessionPeer1, nil)
	mockClientPeer2 := mock_proto.NewMockMetalInfraConfigClient(f.ctrl)
	syncSessionPeer2 := mock_proto.NewMockMetalInfraConfig_SyncClient(f.ctrl)
	mockClientPeer2.EXPECT().Sync(gomock.Any(), gomock.Any()).Return(syncSessionPeer2, nil)
	f.rpc.EXPECT().GetClient(peer1).Return(mockClientPeer1, nil)
	f.rpc.EXPECT().GetClient(peer2).Return(mockClientPeer2, nil)
	f.config.EXPECT().Sync(syncSessionPeer1).Return(nil)
	f.config.EXPECT().Sync(syncSessionPeer2).Return(nil)

	err := f.sync.Sync(ctx)
	require.NoError(t, err)
}

func Test_SyncPeriodically(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	f := newTestFixtures(t)
	peer1 := entities.NewPeer("foobar", "1.2.3.4:5678", time.Now())
	peer2 := entities.NewPeer("barbaz", "5.6.7.8:910", time.Now())
	f.peerDiscovery.EXPECT().GetPeers().Return([]entities.Peer{
		peer1, peer2,
	}, nil)
	mockClientPeer1 := mock_proto.NewMockMetalInfraConfigClient(f.ctrl)
	syncSessionPeer1 := mock_entities.NewMockSyncSession(f.ctrl)
	mockClientPeer1.EXPECT().Sync(gomock.Any(), gomock.Any()).Return(syncSessionPeer1, nil)
	mockClientPeer2 := mock_proto.NewMockMetalInfraConfigClient(f.ctrl)
	syncSessionPeer2 := mock_entities.NewMockSyncSession(f.ctrl)
	mockClientPeer2.EXPECT().Sync(gomock.Any(), gomock.Any()).Return(syncSessionPeer2, nil)
	f.rpc.EXPECT().GetClient(peer1).Return(mockClientPeer1, nil)
	f.rpc.EXPECT().GetClient(peer2).Return(mockClientPeer2, nil)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	f.config.EXPECT().Sync(syncSessionPeer1).DoAndReturn(func(_ proto.MetalInfraConfig_SyncClient) error {
		wg.Done()
		return nil
	})
	f.config.EXPECT().Sync(syncSessionPeer2).DoAndReturn(func(_ proto.MetalInfraConfig_SyncClient) error {
		wg.Done()
		return nil
	})

	closer := f.sync.SyncPeriodically(ctx, time.Millisecond)
	wg.Wait()
	closer.Close()
}
