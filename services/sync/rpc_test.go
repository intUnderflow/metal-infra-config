package sync

import (
	"github.com/intunderflow/metal-infra-config/entities"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"time"
)

func Test_GetClient_WithClientCreationFailure_ReturnsError(t *testing.T) {
	rpc := NewRPC([]grpc.DialOption{})
	_, err := rpc.GetClient(entities.NewPeer("foobar", ":", time.Now()))
	require.ErrorContains(t, err, "grpc")
}

func Test_GetClient_WithValidPeer_ReturnsClient(t *testing.T) {
	rpc := NewRPC([]grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	client, err := rpc.GetClient(entities.NewPeer("foobar", "1.2.3.4:5678", time.Now()))
	require.NoError(t, err)
	require.NotNil(t, client)
}
