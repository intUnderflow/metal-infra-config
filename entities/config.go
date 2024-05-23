package entities

import (
	"errors"
	"sync"
)

const (
	_errValueDoesNotExist = "value does not exist"
	// ErrNewValueIsOlderThanCurrentValue is emitted when the new value being set in Set is older than the current value
	// if this error is emitted the value will not be updated
	ErrNewValueIsOlderThanCurrentValue = "new value is older than current value"
)

// Config is a store of versioned values, each value is indexed by a Key and has a specific Version attached to its
// current value
type Config interface {
	GetWithVersion(Key) (ValueWithVersion, error)
	SetWithVersion(Key, string, uint64) error
	List() map[Key]ValueWithVersion
	Delete(Key) error
	Sync(SyncSession) error
}

type configImpl struct {
	mutex  *sync.RWMutex
	values map[Key]*value
}

// Key is the unique key identifying a value in the Config
type Key string

// ValueWithVersion is the struct returned to callers of Config.GetWithVersion, it allows callers to access both the
// underlying value at the given key and the version attached to that value
type ValueWithVersion struct {
	Version uint64
	Value   string
}

type value struct {
	version uint64
	value   string
}

// NewConfig creates a new instance of Config
func NewConfig() Config {
	return &configImpl{
		mutex:  &sync.RWMutex{},
		values: map[Key]*value{},
	}
}

// GetWithVersion returns a ValueWithVersion representing the underlying value at the provided Key, if no value is
// present it returns an error. The ValueWithVersion is a copy of the underlying value and cannot be used to mutate it.
func (c *configImpl) GetWithVersion(key Key) (ValueWithVersion, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	value, ok := c.values[key]
	if !ok {
		return ValueWithVersion{}, errors.New(_errValueDoesNotExist)
	}

	return valueWithVersionFromInternal(value), nil
}

// SetWithVersion sets the value at the given Key to the specified valueString and version.
// If the value in the Config already has a higher ValueVersion then an error is returned.
func (c *configImpl) SetWithVersion(key Key, valueString string, version uint64) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	oldValue, ok := c.values[key]
	if !ok {
		c.values[key] = &value{
			version: version,
			value:   valueString,
		}
		return nil
	}
	if oldValue.version > version {
		return errors.New(ErrNewValueIsOlderThanCurrentValue)
	}
	oldValue.value = valueString
	oldValue.version = version
	return nil
}

// List returns a map of all values in the Config, each value is mapped to its Key with a ValueWithVersion.
func (c *configImpl) List() map[Key]ValueWithVersion {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.listInternal()
}

func (c *configImpl) listInternal() map[Key]ValueWithVersion {
	allValues := map[Key]ValueWithVersion{}
	for k, v := range c.values {
		allValues[k] = valueWithVersionFromInternal(v)
	}
	return allValues
}

// Delete deletes a Key from the config, it returns an error if the Key doesn't exist
func (c *configImpl) Delete(key Key) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, ok := c.values[key]
	if !ok {
		return errors.New(_errValueDoesNotExist)
	}
	delete(c.values, key)
	return nil
}

func valueWithVersionFromInternal(internalValue *value) ValueWithVersion {
	return ValueWithVersion{
		Version: internalValue.version,
		Value:   internalValue.value,
	}
}
