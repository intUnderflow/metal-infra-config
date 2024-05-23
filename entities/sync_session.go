package entities

import "github.com/intunderflow/metal-infra-config/proto"

type SyncSession interface {
	Send(*proto.SyncRecord) error
	Recv() (*proto.SyncRecord, error)
	CloseSend() error
}

// proto.MetalInfraConfigServer_SyncClient already satisfies the SyncSession interface

type syncSessionServerImpl struct {
	server proto.MetalInfraConfig_SyncServer
}

func NewSyncSessionFromServer(server proto.MetalInfraConfig_SyncServer) SyncSession {
	return &syncSessionServerImpl{
		server: server,
	}
}

func (s *syncSessionServerImpl) Send(record *proto.SyncRecord) error {
	return s.server.Send(record)
}

func (s *syncSessionServerImpl) Recv() (*proto.SyncRecord, error) {
	return s.server.Recv()
}

func (s *syncSessionServerImpl) CloseSend() error {
	// Server does not have an equivalent of CloseSend
	return nil
}
