package fx

import (
	"context"
	mock_closeable "github.com/intunderflow/metal-infra-config/pkg/closeable/mock"
	mock_sync "github.com/intunderflow/metal-infra-config/services/sync/mock"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"testing"
	"time"
)

func Test_Start_StartsSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	syncMock := mock_sync.NewMockSync(ctrl)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	interval := time.Second
	syncMock.EXPECT().SyncPeriodically(ctx, interval).Return(mock_closeable.NewMockCloseable(ctrl))

	fx := NewFX(syncMock)
	fx.Start(ctx, interval)
}

func Test_Stop_StopsSync(t *testing.T) {
	ctrl := gomock.NewController(t)
	syncMock := mock_sync.NewMockSync(ctrl)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	interval := time.Second
	closeableMock := mock_closeable.NewMockCloseable(ctrl)
	closeableMock.EXPECT().Close().Return()
	syncMock.EXPECT().SyncPeriodically(ctx, interval).Return(closeableMock)

	fx := NewFX(syncMock)
	fx.Start(ctx, interval)
	require.NoError(t, fx.Stop())
}

func Test_Stop_WhenAlreadyStopped_ReturnsError(t *testing.T) {
	ctrl := gomock.NewController(t)
	fx := NewFX(mock_sync.NewMockSync(ctrl))
	require.EqualError(t, fx.Stop(), _errAlreadyStoppedOrNeverStarted)
}
