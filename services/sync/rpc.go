package sync

import (
	"github.com/intunderflow/metal-infra-config/entities"
	"github.com/intunderflow/metal-infra-config/proto"
	"google.golang.org/grpc"
)

type RPC interface {
	GetClient(peer entities.Peer) (proto.InternalSyncClient, error)
}

type rpcImpl struct {
	dialOptions []grpc.DialOption
}

func NewRPC(dialOptions []grpc.DialOption) RPC {
	return &rpcImpl{
		dialOptions: dialOptions,
	}
}

func (r *rpcImpl) GetClient(peer entities.Peer) (proto.InternalSyncClient, error) {
	clientConn, err := grpc.NewClient(peer.Address(), r.dialOptions...)
	if err != nil {
		return nil, err
	}

	return proto.NewInternalSyncClient(clientConn), nil
}
