package peerdiscovery

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/intunderflow/metal-infra-config/entities"
	"github.com/stretchr/testify/require"
	"testing"
)

func Test_GetPeers_ReturnsCorrectPeers(t *testing.T) {
	config := entities.NewConfig()
	for i := 0; i < 100; i++ {
		require.NoError(t, config.SetWithVersion(entities.Key(uuid.NewString()), uuid.NewString(), 1))
	}

	peerDiscovery := NewPeerDiscovery(config)

	require.NoError(t, config.SetWithVersion("c8690439-293c-494d-9d26-257d64402b42._internal.peer_address", "peer_address_1", 1))
	require.NoError(t, config.SetWithVersion("c8690439-293c-494d-9d26-257d64402b42._internal.peer_last_seen", "1337", 1))

	require.NoError(t, config.SetWithVersion("4ec63e94-4653-4536-bd10-9bab71d17ade._internal.peer_address", "peer_address_2", 1))
	require.NoError(t, config.SetWithVersion("4ec63e94-4653-4536-bd10-9bab71d17ade._internal.peer_last_seen", "9001", 1))

	peers, err := peerDiscovery.GetPeers()
	require.NoError(t, err)
	require.Len(t, peers, 2)
	for _, peer := range peers {
		if peer.ID() == "c8690439-293c-494d-9d26-257d64402b42" {
			require.Equal(t, "peer_address_1", peer.Address())
			require.Equal(t, int64(1337), peer.LastSeen().Unix())
		} else if peer.ID() == "4ec63e94-4653-4536-bd10-9bab71d17ade" {
			require.Equal(t, "peer_address_2", peer.Address())
			require.Equal(t, int64(9001), peer.LastSeen().Unix())
		} else {
			t.Error(fmt.Sprintf("unexpected peer ID %s", peer.ID()))
		}
	}
}

func Test_GetPeers_WhenLastSeenInvalid_ReturnsErrorButStillReturnsOtherPeers(t *testing.T) {
	config := entities.NewConfig()
	for i := 0; i < 100; i++ {
		require.NoError(t, config.SetWithVersion(entities.Key(uuid.NewString()), uuid.NewString(), 1))
	}

	peerDiscovery := NewPeerDiscovery(config)

	require.NoError(t, config.SetWithVersion("c8690439-293c-494d-9d26-257d64402b42._internal.peer_address", "peer_address_1", 1))
	require.NoError(t, config.SetWithVersion("c8690439-293c-494d-9d26-257d64402b42._internal.peer_last_seen", "1337", 1))

	require.NoError(t, config.SetWithVersion("4ec63e94-4653-4536-bd10-9bab71d17ade._internal.peer_address", "peer_address_2", 1))
	require.NoError(t, config.SetWithVersion("4ec63e94-4653-4536-bd10-9bab71d17ade._internal.peer_last_seen", "abc", 1))

	peers, err := peerDiscovery.GetPeers()
	require.ErrorContains(t, err, _errInvalidPeerLastSeen)
	require.Len(t, peers, 1)
	require.Equal(t, entities.PeerID("c8690439-293c-494d-9d26-257d64402b42"), peers[0].ID())
	require.Equal(t, "peer_address_1", peers[0].Address())
	require.Equal(t, int64(1337), peers[0].LastSeen().Unix())
}
