package entities

import (
	"errors"
	"github.com/intunderflow/metal-infra-config/proto"
	mock_proto "github.com/intunderflow/metal-infra-config/proto/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
)

func Test_NewSyncSessionFromServer_Send(t *testing.T) {
	ctrl := gomock.NewController(t)
	protoRecord := &proto.SyncRecord{
		Key:     "foo",
		Value:   "bar",
		Version: 5,
	}
	expectedErr := errors.New("foobar")
	protoSession := mock_proto.NewMockMetalInfraConfig_SyncServer(ctrl)
	protoSession.EXPECT().Send(protoRecord).Return(expectedErr)
	session := NewSyncSessionFromServer(protoSession)
	require.EqualError(t, session.Send(protoRecord), expectedErr.Error())
}

func Test_NewSyncSessionFromServer_Recv(t *testing.T) {
	ctrl := gomock.NewController(t)
	protoRecord := &proto.SyncRecord{
		Key:     "foo",
		Value:   "bar",
		Version: 5,
	}
	protoSession := mock_proto.NewMockMetalInfraConfig_SyncServer(ctrl)
	protoSession.EXPECT().Recv().Return(protoRecord, nil)
	session := NewSyncSessionFromServer(protoSession)
	got, err := session.Recv()
	require.NoError(t, err)
	require.EqualValues(t, protoRecord, got)
}

func Test_NewSyncSessionFromServer_CloseSend(t *testing.T) {
	ctrl := gomock.NewController(t)
	session := NewSyncSessionFromServer(mock_proto.NewMockMetalInfraConfig_SyncServer(ctrl))
	require.NoError(t, session.CloseSend())
}
