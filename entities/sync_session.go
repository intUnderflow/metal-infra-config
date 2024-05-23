package entities

import "github.com/intunderflow/metal-infra-config/proto"

type SyncSession interface {
	Send(*proto.SyncRecord) error
	Recv() (*proto.SyncRecord, error)
	CloseSend() error
}

type syncSessionClientImpl struct {
	client proto.MetalInfraConfig_SyncClient
}

func NewSyncSessionFromClient(client proto.MetalInfraConfig_SyncClient) SyncSession {
	return &syncSessionClientImpl{
		client: client,
	}
}

func (s *syncSessionClientImpl) Send(record *proto.SyncRecord) error {
	return s.client.Send(record)
}

func (s *syncSessionClientImpl) Recv() (*proto.SyncRecord, error) {
	return s.client.Recv()
}

func (s *syncSessionClientImpl) CloseSend() error {
	return s.client.CloseSend()
}
