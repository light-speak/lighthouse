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
func GetOrSet(key string, result interface{}, fn func() (interface{}, error)) error {
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

	if err := Set(key, value, DefaultExpiration); err != nil {
		return err
	}

	return convertValue(value, result)
}

// CacheStruct caches struct type
func CacheStruct[T any](key string, result *T, fn func() (*T, error)) error {
	return GetOrSet(key, result, func() (interface{}, error) {
		return fn()
	})
}

// CacheMap caches map type
func CacheMap(key string, result *map[string]interface{}, fn func() (map[string]interface{}, error)) error {
	return GetOrSet(key, result, func() (interface{}, error) {
		return fn()
	})
}

// CacheSyncMap caches sync.Map type
func CacheSyncMap(key string, result *sync.Map, fn func() (*sync.Map, error)) error {
	if !isEnabled {
		value, err := fn()
		if err != nil {
			return err
		}
		// Copy values instead of the map itself
		value.Range(func(k, v interface{}) bool {
			result.Store(k, v)
			return true
		})
		return nil
	}

	// Try to get from cache first
	var tempMap map[string]interface{}
	err := Get(key, &tempMap)
	if err == nil {
		// Convert to sync.Map
		for k, v := range tempMap {
			result.Store(k, v)
		}
		return nil
	}

	// If cache miss, execute function to get data
	value, err := fn()
	if err != nil {
		return err
	}

	// Convert sync.Map to regular map for caching
	tempMap = make(map[string]interface{})
	value.Range(func(k, v interface{}) bool {
		if keyStr, ok := k.(string); ok {
			tempMap[keyStr] = v
		}
		return true
	})

	// Store in cache
	if err := Set(key, tempMap, DefaultExpiration); err != nil {
		return err
	}

	// Copy values instead of the map itself
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
