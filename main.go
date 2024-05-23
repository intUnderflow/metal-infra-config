package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
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
	"google.golang.org/grpc/credentials"
	"net"
	"os"
	"time"
)

const (
	_envPort            = "PORT"
	_envTLSRootCAFile   = "TLS_ROOT_CA_FILE"
	_envTLSNodeCertFile = "TLS_NODE_CA_FILE"
	_envTLSNodeKeyFile  = "TLS_NODE_KEY_FILE"
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
			// Server initialisation
			func(g *grpc.Server) grpc.ServiceRegistrar {
				return g
			},
			func() (port.Port, error) {
				return port.ParsePort(os.Getenv(_envPort))
			},
			func(p port.Port) (net.Listener, error) {
				return net.Listen("tcp", fmt.Sprintf(":%d", p))
			},
			func(opts []grpc.ServerOption) *grpc.Server {
				return grpc.NewServer(opts...)
			},

			// TLS root CA
			func() (*x509.CertPool, error) {
				pemClientCA, err := os.ReadFile(os.Getenv(_envTLSRootCAFile))
				if err != nil {
					return nil, err
				}

				certPool := x509.NewCertPool()
				if !certPool.AppendCertsFromPEM(pemClientCA) {
					return nil, fmt.Errorf("failed to add client CA's certificate")
				}

				return certPool, nil
			},
			// Server TLS
			func() (tls.Certificate, error) {
				return tls.LoadX509KeyPair(os.Getenv(_envTLSNodeCertFile), os.Getenv(_envTLSNodeKeyFile))
			},
			func(cert tls.Certificate, rootCA *x509.CertPool) []grpc.ServerOption {
				return []grpc.ServerOption{
					grpc.Creds(credentials.NewTLS(&tls.Config{
						Certificates: []tls.Certificate{cert},
						ClientAuth:   tls.RequireAndVerifyClientCert,
						RootCAs:      rootCA,
					})),
				}
			},
			// Client TLS
			func() ([]grpc.DialOption, error) {
				tlsCredentials, err := credentials.NewClientTLSFromFile(os.Getenv(_envTLSNodeCertFile), "")
				if err != nil {
					return nil, err
				}
				return []grpc.DialOption{
					grpc.WithTransportCredentials(tlsCredentials),
				}, nil
			},

			// Services
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
