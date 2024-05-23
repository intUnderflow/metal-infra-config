#!/bin/bash

REPO="github.com/intunderflow/metal-infra-config"

go install go.uber.org/mock/mockgen@latest
mockgen -destination "entities/mock/mock.go" "$REPO/entities" Config
mockgen -destination "pkg/closeable/mock/mock.go" "$REPO/pkg/closeable" Closeable
mockgen -destination "proto/mock/mock.go" "$REPO/proto" MetalInfraConfigClient,MetalInfraConfig_SyncClient,MetalInfraConfig_SyncServer
mockgen -destination "services/peerdiscovery/mock/mock.go" "$REPO/services/peerdiscovery" PeerDiscovery
mockgen -destination "services/sync/mock/mock.go" "$REPO/services/sync" Sync,RPC