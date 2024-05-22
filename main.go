package main

import (
	"fmt"
	"github.com/intunderflow/metal-infra-config/entities"
	"github.com/intunderflow/metal-infra-config/handler"
	"github.com/intunderflow/metal-infra-config/pkg/port"
	"github.com/intunderflow/metal-infra-config/proto"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"net"
	"os"
)

func main() {
	fx.New(opts()).Run()
}

func opts() fx.Option {
	return fx.Options(
		fx.Provide(
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
		),
		fx.Invoke(
			proto.RegisterMetalInfraConfigServer,
			func(server *grpc.Server, listener net.Listener) error {
				return server.Serve(listener)
			},
		),
	)
}
