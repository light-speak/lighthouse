package model

import (
	"fmt"
	"sync"
	"time"

	"github.com/light-speak/lighthouse/errors"
	"gorm.io/gorm"
)

// LoaderConfig defines the configuration for a data loader
type LoaderConfig[T ModelInterface] struct {
	MaxBatch int                                // Maximum number of keys to batch in a single fetch
	Wait     time.Duration                      // Amount of time to wait before processing a batch
	Fetch    func(keys []int64) ([]*T, []error) // Function to fetch data for a batch of keys
}

// NewLoader creates a new Loader with the given configuration
func GetLoader[T ModelInterface](db *gorm.DB, config *LoaderConfig[T]) *Loader[T] {
	if config == nil {
		config = &LoaderConfig[T]{
			MaxBatch: 100,
			Wait:     1 * time.Millisecond,
			Fetch: func(keys []int64) ([]*T, []error) {
				var data []*T
				if err := db.Find(&data, keys).Error; err != nil {
					return nil, nil
				}
				dataById := map[int64]*T{}
				for _, d := range data {
					dataById[(*d).GetId()] = d
				}
				result := make([]*T, len(keys))
				for i, key := range keys {
					result[i] = dataById[key]
				}
				return result, nil
			},
		}
	}
	return &Loader[T]{
		maxBatch: config.MaxBatch,
		wait:     config.Wait,
		fetch:    config.Fetch,
		cache:    make(map[int64]*T), // Initialize cache
	}
}

// Loader represents a data loader that can batch and cache requests
type Loader[T ModelInterface] struct {
	wait     time.Duration                      // Amount of time to wait before processing a batch
	maxBatch int                                // Maximum number of keys to batch in a single fetch
	fetch    func(keys []int64) ([]*T, []error) // Function to fetch data for a batch of keys

	cache map[int64]*T    // Cache to store fetched data
	batch *LoaderBatch[T] // Current batch of requests
	mu    sync.Mutex      // Mutex to protect concurrent access
}

// LoaderBatch represents a batch of data loader requests
type LoaderBatch[T ModelInterface] struct {
	keys    []int64
	results []*T
	errors  []error
	closing bool
	done    chan struct{}
}

// Load loads a single item by key
func (l *Loader[T]) Load(key int64) (T, error) {
	ptr, _ := l.loadThunk(key)()
	return *ptr, nil
}

// loadThunk creates a thunk for loading a single item
func (l *Loader[T]) loadThunk(key int64) func() (*T, error) {
	l.mu.Lock()
	if d, ok := l.cache[key]; ok {
		l.mu.Unlock()
		return func() (*T, error) { return d, nil }
	}
	if l.batch == nil {
		l.batch = &LoaderBatch[T]{done: make(chan struct{})}
	}
	batch := l.batch
	pos := batch.keyIndex(l, key)
	l.mu.Unlock()

	return func() (*T, error) {
		<-batch.done

		if pos < len(batch.results) {
			res := batch.results[pos]
			if len(batch.errors) == 0 || batch.errors[pos] == nil {
				l.set(key, res)
				return res, nil
			}
			return res, batch.errors[pos]
		}

		return nil, &errors.DataloaderError{Msg: fmt.Sprintf("key %d not found", key)}
	}
}

// LoadAll loads multiple items by keys
func (l *Loader[T]) LoadAll(keys []int64) ([]*T, []error) {
	thunks := make([]func() (*T, error), len(keys))
	for i, key := range keys {
		thunks[i] = l.loadThunk(key)
	}

	results := make([]*T, len(keys))
	errors := make([]error, len(keys))
	for i, thunk := range thunks {
		results[i], errors[i] = thunk()
	}
	return results, errors
}

// LoadAllThunk creates a thunk for loading multiple items
func (l *Loader[T]) LoadAllThunk(keys []int64) func() ([]*T, []error) {
	thunks := make([]func() (*T, error), len(keys))
	for i, key := range keys {
		thunks[i] = l.loadThunk(key)
	}

	return func() ([]*T, []error) {
		results := make([]*T, len(keys))
		errors := make([]error, len(keys))
		for i, thunk := range thunks {
			results[i], errors[i] = thunk()
		}
		return results, errors
	}
}

// Prime adds the provided key-value pair to the cache
func (l *Loader[T]) Prime(key int64, value *T) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.cache[key]; !ok {
		copy := *value
		l.unsafeSet(key, &copy)
		return true
	}
	return false
}

// Clear removes the value associated with key from the cache
func (l *Loader[T]) Clear(key int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.cache, key)
}

// set adds the provided key-value pair to the cache
func (l *Loader[T]) set(key int64, value *T) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.unsafeSet(key, value)
}

// unsafeSet adds the provided key-value pair to the cache without locking
func (l *Loader[T]) unsafeSet(key int64, value *T) {
	l.cache[key] = value
}

// keyIndex returns the index of the given key in the batch, adding it if not present
func (b *LoaderBatch[T]) keyIndex(l *Loader[T], key int64) int {
	for i, k := range b.keys {
		if k == key {
			return i
		}
	}
	pos := len(b.keys)
	b.keys = append(b.keys, key)
	if pos == 0 {
		go b.startTimer(l)
	}
	if l.maxBatch != 0 && pos >= l.maxBatch-1 {
		if !b.closing {
			b.closing = true
			l.batch = nil
			go b.end(l)
		}
	}
	return pos
}

// startTimer starts the timer for processing the batch
func (b *LoaderBatch[T]) startTimer(l *Loader[T]) {
	time.Sleep(l.wait)
	l.mu.Lock()
	defer l.mu.Unlock()

	if b.closing {
		return
	}
	l.batch = nil
	go b.end(l)
}

// end processes the batch
func (b *LoaderBatch[T]) end(l *Loader[T]) {
	b.results, b.errors = l.fetch(b.keys)
	close(b.done)
}

// ClearAll removes all entries from the cache
func (l *Loader[T]) ClearAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cache = make(map[int64]*T)
}

// CacheSize returns the number of entries in the cache
func (l *Loader[T]) CacheSize() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.cache)
}
