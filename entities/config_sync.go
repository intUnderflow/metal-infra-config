package entities

import (
	"errors"
	"github.com/intunderflow/metal-infra-config/proto"
	"io"
)

func (c *configImpl) Sync(session SyncSession) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	values := c.listInternal()
	for key, record := range values {
		err := session.Send(&proto.SyncRecord{
			Key:     string(key),
			Value:   record.Value,
			Version: record.Version,
		})
		if err != nil {
			return err
		}
	}
	err := session.CloseSend()
	if err != nil {
		return err
	}
	for {
		record, err := session.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return err
		}
		newKey := protoToKey(record)
		newValue := protoToValue(record)
		existingValue, ok := c.values[newKey]
		if !ok || existingValue.version < newValue.version {
			c.values[newKey] = newValue
		}
	}
	return nil
}

func protoToKey(record *proto.SyncRecord) Key {
	return Key(record.Key)
}

func protoToValue(record *proto.SyncRecord) *value {
	return &value{
		version: record.Version,
		value:   record.Value,
	}
}
