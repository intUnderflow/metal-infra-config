package entities

import (
	"errors"
	mock_entities "github.com/intunderflow/metal-infra-config/entities/mock"
	"github.com/intunderflow/metal-infra-config/proto"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"io"
	"testing"
)

func Test_Sync_WhenErrorSendingRecords_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSession := mock_entities.NewMockSyncSession(ctrl)
	expectedErr := errors.New("foobar")
	mockSession.EXPECT().Send(gomock.Any()).Return(expectedErr)

	config := NewConfig()
	require.NoError(t, config.SetWithVersion("abc", "def", 1))

	err := config.Sync(mockSession)
	require.EqualError(t, err, expectedErr.Error())
}

func Test_Sync_WhenErrorClosingSend_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSession := mock_entities.NewMockSyncSession(ctrl)
	expectedErr := errors.New("foobar")
	mockSession.EXPECT().Send(gomock.Any()).Times(2).Return(nil)
	mockSession.EXPECT().CloseSend().Return(expectedErr)

	config := NewConfig()
	require.NoError(t, config.SetWithVersion("abc", "def", 1))
	require.NoError(t, config.SetWithVersion("xyz", "abc", 1))

	err := config.Sync(mockSession)
	require.EqualError(t, err, expectedErr.Error())
}

func Test_Sync_WhenErrorReceivingRecords_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSession := mock_entities.NewMockSyncSession(ctrl)
	expectedErr := errors.New("foobar")
	mockSession.EXPECT().Send(gomock.Any()).Times(2).Return(nil)
	mockSession.EXPECT().CloseSend().Return(nil)
	mockSession.EXPECT().Recv().Return(nil, expectedErr)

	config := NewConfig()
	require.NoError(t, config.SetWithVersion("abc", "def", 1))
	require.NoError(t, config.SetWithVersion("xyz", "abc", 1))

	err := config.Sync(mockSession)
	require.EqualError(t, err, expectedErr.Error())
}

func Test_Sync_SynchronisesOverSession(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockSession := mock_entities.NewMockSyncSession(ctrl)
	mockSession.EXPECT().Send(gomock.Any()).Times(2).Return(nil)
	mockSession.EXPECT().CloseSend().Return(nil)
	// New value that we expect to be populated
	mockSession.EXPECT().Recv().Return(&proto.SyncRecord{
		Key:     "newkey",
		Value:   "newvalue",
		Version: 6,
	}, nil)
	// Update to existing value that we expect to be populated
	mockSession.EXPECT().Recv().Return(&proto.SyncRecord{
		Key:     "abc",
		Value:   "aaa",
		Version: 2,
	}, nil)
	// Older version of existing value that we expect to be ignored
	mockSession.EXPECT().Recv().Return(&proto.SyncRecord{
		Key:     "xyz",
		Value:   "bbb",
		Version: 3,
	}, nil)
	// Signal for no more messages
	mockSession.EXPECT().Recv().Return(nil, io.EOF)

	config := NewConfig()
	require.NoError(t, config.SetWithVersion("abc", "def", 1))
	require.NoError(t, config.SetWithVersion("xyz", "abc", 4))

	err := config.Sync(mockSession)
	require.NoError(t, err)

	value, err := config.GetWithVersion("newkey")
	require.NoError(t, err)
	require.Equal(t, "newvalue", value.Value)
	require.Equal(t, uint64(6), value.Version)

	value, err = config.GetWithVersion("abc")
	require.NoError(t, err)
	require.Equal(t, "aaa", value.Value)
	require.Equal(t, uint64(2), value.Version)

	value, err = config.GetWithVersion("xyz")
	require.NoError(t, err)
	require.Equal(t, "abc", value.Value)
	require.Equal(t, uint64(4), value.Version)
}
