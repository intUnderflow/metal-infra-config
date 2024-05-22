package sync

import (
	"github.com/intunderflow/metal-infra-config/entities"
	"github.com/intunderflow/metal-infra-config/proto"
	"google.golang.org/grpc"
)

type RPC struct {
	dialOptions []grpc.DialOption
}

func NewRPC(dialOptions []grpc.DialOption) *RPC {
	return &RPC{
		dialOptions: dialOptions,
	}
}

func (r *RPC) getClient(peer entities.Peer) (proto.InternalSyncClient, error) {
	clientConn, err := grpc.NewClient(peer.Address(), r.dialOptions...)
	if err != nil {
		return nil, err
	}

	return proto.NewInternalSyncClient(clientConn), nil
}
