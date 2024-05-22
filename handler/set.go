package handler

import (
	"context"
	"errors"
	"github.com/intunderflow/metal-infra-config/entities"
	"github.com/intunderflow/metal-infra-config/pkg/auth"
	"github.com/intunderflow/metal-infra-config/proto"
)

const (
	_errKeyRequired     = "key is required"
	_errValueRequired   = "value is required"
	_errVersionRequired = "version is required"
)

func (h *Handler) Set(ctx context.Context, req *proto.SetRequest) (*proto.SetResponse, error) {
	if req.Key == "" {
		return nil, errors.New(_errKeyRequired)
	}
	if req.Value == "" {
		return nil, errors.New(_errValueRequired)
	}
	if req.Version == 0 {
		return nil, errors.New(_errVersionRequired)
	}

	peerID, err := auth.GetPeerIDFromGRPCContext(ctx)
	if err != nil {
		return nil, err
	}
	
	keyInConfig := configKeyForPeer(peerID, req.Key)
	err = h.config.SetWithVersion(keyInConfig, req.Value, req.Version)
	if err != nil {
		return nil, err
	}

	return &proto.SetResponse{}, nil
}

func configKeyForPeer(peerID entities.PeerID, key string) entities.Key {
	return entities.Key(string(peerID) + "." + key)
}
