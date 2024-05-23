package handler

import (
	"github.com/intunderflow/metal-infra-config/entities"
	"github.com/intunderflow/metal-infra-config/proto"
)

func (h *Handler) Sync(stream proto.MetalInfraConfig_SyncServer) error {
	return h.config.Sync(entities.NewSyncSessionFromServer(stream))
}
