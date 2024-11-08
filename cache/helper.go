package cache

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

const (
	DefaultExpiration = 30 * time.Minute
)

// GetOrSet gets value from cache, if not exists then sets it
func GetOrSet(key string, result interface{}, fn func() (interface{}, error), tags ...string) error {
	if !isEnabled {
		value, err := fn()
		if err != nil {
			return err
		}
		return convertValue(value, result)
	}

	err := Get(key, result)
	if err == nil {
		return nil
	}

	value, err := fn()
	if err != nil {
		return err
	}

	if err := Set(key, value, DefaultExpiration, tags...); err != nil {
		return err
	}

	return convertValue(value, result)
}

// CacheStruct caches struct type
func CacheStruct[T any](key string, result *T, fn func() (*T, error), tags ...string) error {
	return GetOrSet(key, result, func() (interface{}, error) {
		return fn()
	}, tags...)
}

// CacheMap caches map type
func CacheMap(key string, result *map[string]interface{}, fn func() (map[string]interface{}, error), tags ...string) error {
	return GetOrSet(key, result, func() (interface{}, error) {
		return fn()
	}, tags...)
}

// CacheSyncMap caches sync.Map type
func CacheSyncMap(key string, result *sync.Map, fn func() (*sync.Map, error), tags ...string) error {
	if !isEnabled {
		value, err := fn()
		if err != nil {
			return err
		}
		value.Range(func(k, v interface{}) bool {
			result.Store(k, v)
			return true
		})
		return nil
	}

	var tempMap map[string]interface{}
	err := Get(key, &tempMap)
	if err == nil {
		for k, v := range tempMap {
			result.Store(k, v)
		}
		return nil
	}

	value, err := fn()
	if err != nil {
		return err
	}

	tempMap = make(map[string]interface{})
	value.Range(func(k, v interface{}) bool {
		if keyStr, ok := k.(string); ok {
			tempMap[keyStr] = v
		}
		return true
	})

	if err := Set(key, tempMap, DefaultExpiration, tags...); err != nil {
		return err
	}

	value.Range(func(k, v interface{}) bool {
		result.Store(k, v)
		return true
	})
	return nil
}

// convertValue converts value type
func convertValue(from interface{}, to interface{}) error {
	// If types match, copy values instead of the map itself
	if fromPtr, ok := from.(*sync.Map); ok {
		if toPtr, ok := to.(*sync.Map); ok {
			fromPtr.Range(func(k, v interface{}) bool {
				toPtr.Store(k, v)
				return true
			})
			return nil
		}
	}

	// Otherwise convert through JSON
	data, err := json.Marshal(from)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	return json.Unmarshal(data, to)
}
