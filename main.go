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
	return fx.Provide(
		entities.NewConfig,
		handler.NewHandler,
		grpc.NewServer,
		proto.RegisterMetalInfraConfigServer,
	)
}
