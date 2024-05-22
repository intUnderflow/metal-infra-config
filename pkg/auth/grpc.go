package auth

import (
	"context"
	"errors"
	"github.com/intunderflow/metal-infra-config/entities"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/peer"
)

const (
	// ErrContextNotAGRPCPeer is emitted if the provided context in GetPeerIDFromGRPCContext is not a valid peer.Peer
	ErrContextNotAGRPCPeer = "context is not a GRPC peer"
	_errNoTLSInfoForPeer   = "no TLS info for GRPC peer"
	_errNoCertificate      = "no certificate for GRPC peer"
)

func GetPeerIDFromGRPCContext(ctx context.Context) (entities.PeerID, error) {
	grpcPeer, ok := peer.FromContext(ctx)
	if !ok {
		return "", errors.New(ErrContextNotAGRPCPeer)
	}

	tlsInfo, ok := grpcPeer.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return "", errors.New(_errNoTLSInfoForPeer)
	}

	chains := tlsInfo.State.VerifiedChains
	if len(chains) == 0 {
		return "", errors.New(_errNoCertificate)
	}

	if len(chains[0]) == 0 {
		return "", errors.New(_errNoCertificate)
	}

	return getPeerID(chains[0][0])
}
