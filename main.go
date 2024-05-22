package main

import (
	"github.com/intunderflow/metal-infra-config/entities"
	"github.com/intunderflow/metal-infra-config/handler"
	"github.com/intunderflow/metal-infra-config/proto"
	"go.uber.org/fx"
	"google.golang.org/grpc"
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
		),
		fx.Invoke(
			proto.RegisterMetalInfraConfigServer,
		),
	)
}
