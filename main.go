package main

import (
	"context"
	"fmt"
	"github.com/intunderflow/metal-infra-config/entities"
	"github.com/intunderflow/metal-infra-config/handler"
	"github.com/intunderflow/metal-infra-config/pkg/port"
	"github.com/intunderflow/metal-infra-config/proto"
	"github.com/intunderflow/metal-infra-config/services/peerdiscovery"
	"github.com/intunderflow/metal-infra-config/services/sync"
	syncfx "github.com/intunderflow/metal-infra-config/services/sync/fx"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
	"os"
	"time"
)

func main() {
	fx.New(opts()).Run()
}

func opts() fx.Option {
	return fx.Options(
		fx.Provide(
			context.Background,
			zap.NewProduction,
			entities.NewConfig,
			handler.NewHandler,
			grpc.NewServer,
			func(g *grpc.Server) grpc.ServiceRegistrar {
				return g
			},
			func() (port.Port, error) {
				return port.ParsePort(os.Getenv("PORT"))
			},
			func(p port.Port) (net.Listener, error) {
				return net.Listen("tcp", fmt.Sprintf(":%d", p))
			},
			func() []grpc.DialOption {
				return []grpc.DialOption{}
			},
			sync.NewRPC,
			sync.NewSync,
			syncfx.NewFX,
			peerdiscovery.NewPeerDiscovery,
		),
		fx.Invoke(
			proto.RegisterMetalInfraConfigServer,
			func(lc fx.Lifecycle, server *grpc.Server, listener net.Listener) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						return server.Serve(listener)
					},
					OnStop: func(ctx context.Context) error {
						server.Stop()
						return nil
					},
				})
			},
			func(lc fx.Lifecycle, syncFX *syncfx.SyncFX) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						syncFX.Start(ctx, time.Second*20)
						return nil
					},
					OnStop: func(_ context.Context) error {
						return syncFX.Stop()
					},
				})
			},
		),
	)
}
