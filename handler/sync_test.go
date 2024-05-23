package handler

import (
	"github.com/intunderflow/metal-infra-config/entities"
	"github.com/intunderflow/metal-infra-config/proto"
	mock_proto "github.com/intunderflow/metal-infra-config/proto/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"testing"
)

func Test_Sync_SynchronisesValues(t *testing.T) {
	config := entities.NewConfig()
	require.NoError(t, config.SetWithVersion("foo", "bar", 2))

	handler := NewHandler(config)
	mockSession := mock_proto.NewMockMetalInfraConfig_SyncServer(gomock.NewController(t))
	mockSession.EXPECT().Send(&proto.SyncRecord{
		Key:     "foo",
		Value:   "bar",
		Version: 2,
	}).Return(nil)
	mockSession.EXPECT().Recv().Return(&proto.SyncRecord{
		Key:     "baz",
		Value:   "aaa",
		Version: 1,
	}, nil)
	mockSession.EXPECT().Recv().Return(&proto.SyncRecord{
		Key:     "foo",
		Value:   "old",
		Version: 1,
	}, nil)
	mockSession.EXPECT().Recv().Return(nil, io.EOF)

	err := handler.Sync(mockSession)
	require.NoError(t, err)
	value, err := config.GetWithVersion("foo")
	require.NoError(t, err)
	require.Equal(t, value.Value, "bar")
	require.Equal(t, value.Version, uint64(2))
	value, err = config.GetWithVersion("baz")
	require.NoError(t, err)
	require.Equal(t, value.Value, "aaa")
	require.Equal(t, value.Version, uint64(1))
}
