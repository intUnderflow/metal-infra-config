#!/bin/bash

REPO="github.com/intunderflow/metal-infra-config"

go install go.uber.org/mock/mockgen@latest
mockgen -destination "entities/mock/config.go" "$REPO/entities" Config
mockgen -destination "proto/mock/mockmetalinfraconfig.go" "$REPO/proto" InternalSyncClient,InternalSync_SyncClient
mockgen -destination "services/peerdiscovery/mock/peerdiscovery.go" "$REPO/services/peerdiscovery" PeerDiscovery
mockgen -destination "services/sync/mock/rpc.go" "$REPO/services/sync" RPC