package handler

import (
	"github.com/intunderflow/metal-infra-config/entities"
	"github.com/intunderflow/metal-infra-config/proto"
)

func NewHandler(config *entities.Config) proto.MetalInfraConfigServer {
	return &Handler{
		config: config,
	}
}

type Handler struct {
	config *entities.Config
	proto.UnimplementedMetalInfraConfigServer
}
