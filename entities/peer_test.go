package entities

import (
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func Test_Peer_ID(t *testing.T) {
	peer := NewPeer("foobar", "example.com:1234", time.Unix(1337, 0))
	require.Equal(t, PeerID("foobar"), peer.ID())
}

func Test_Peer_Address(t *testing.T) {
	peer := NewPeer("foobar", "example.com:1234", time.Unix(1337, 0))
	require.Equal(t, "example.com:1234", peer.Address())
}

func Test_Peer_LastSeen(t *testing.T) {
	peer := NewPeer("foobar", "example.com:1234", time.Unix(1337, 0))
	require.Equal(t, int64(1337), peer.LastSeen().Unix())
}
