package model

import (
	"fmt"
	"sync"
	"time"

	"github.com/light-speak/lighthouse/errors"
	"gorm.io/gorm"
)

var (
	loaderCache = make(map[string]interface{})
	loaderMutex sync.Mutex
)

// Initialize loader cache cleanup
func init() {
	// Run cleanup every hour
	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for range ticker.C {
			cleanupLoaderCache()
		}
	}()
}

// cleanupLoaderCache removes expired entries from the loader cache
func cleanupLoaderCache() {
	loaderMutex.Lock()
	defer loaderMutex.Unlock()

	// Clear all entries and let them be recreated on next access
	loaderCache = make(map[string]interface{})
}

// LoaderConfig defines the configuration for a data loader
type LoaderConfig[K comparable] struct {
	MaxBatch int                                                  // Maximum number of keys to batch in a single fetch
	Wait     time.Duration                                        // Amount of time to wait before processing a batch
	Fetch    func(keys []K) ([][]map[string]interface{}, []error) // Function to fetch data for a batch of keys

	Field string
}

func GetLoader[K comparable](db *gorm.DB, table string, field string, config ...*LoaderConfig[K]) *Loader[K] {
	if len(config) > 0 {
		field = config[0].Field
	}
	loaderKey := fmt.Sprintf("%T-%v", table, field)
	loaderMutex.Lock()
	defer loaderMutex.Unlock()
	if loader, ok := loaderCache[loaderKey]; ok {
		return loader.(*Loader[K])
	}
	loader := getLoader(db, table, field, config...)
	loaderCache[loaderKey] = loader
	return loader
}

// NewLoader creates a new Loader with the given configuration
func getLoader[K comparable](db *gorm.DB, table string, field string, c ...*LoaderConfig[K]) *Loader[K] {
	config := &LoaderConfig[K]{
		MaxBatch: 100,
		Wait:     10 * time.Millisecond,
		Field:    field,
	}

	if len(c) == 0 {
		fetch := func(keys []K) ([][]map[string]interface{}, []error) {
			var data []map[string]interface{}
			if err := db.Table(table).Where(fmt.Sprintf("%s IN (?)", config.Field), keys).Find(&data).Error; err != nil {
				return nil, []error{err}
			}

			dataByKey := make(map[K][]map[string]interface{})
			for _, d := range data {
				if _, ok := dataByKey[d[config.Field].(K)]; !ok {
					dataByKey[d[config.Field].(K)] = make([]map[string]interface{}, 0)
				}
				dataByKey[d[config.Field].(K)] = append(dataByKey[d[config.Field].(K)], d)
			}

			result := make([][]map[string]interface{}, len(keys))
			for i, key := range keys {
				if val, ok := dataByKey[key]; ok {
					result[i] = val
				} else {
					result[i] = nil
				}
			}
			return result, nil
		}
		config.Fetch = fetch
	} else {
		config = c[0]
	}

	return &Loader[K]{
		maxBatch: config.MaxBatch,
		wait:     config.Wait,
		fetch:    config.Fetch,
		cache:    make(map[K][]map[string]interface{}), // Initialize cache
		Field:    config.Field,
		expiry:   make(map[K]time.Time),
	}
}

// Loader represents a data loader that can batch and cache requests
type Loader[K comparable] struct {
	wait     time.Duration                                        // Amount of time to wait before processing a batch
	maxBatch int                                                  // Maximum number of keys to batch in a single fetch
	fetch    func(keys []K) ([][]map[string]interface{}, []error) // Function to fetch data for a batch of keys

	cache  map[K][]map[string]interface{} // Cache to store fetched data
	batch  *LoaderBatch[K]                // Current batch of requests
	expiry map[K]time.Time
	mu     sync.Mutex // Mutex to protect concurrent access
	Field  string
}

// LoaderBatch represents a batch of data loader requests
type LoaderBatch[K comparable] struct {
	keys    []K
	results [][]map[string]interface{}
	errors  []error
	closing bool
	done    chan struct{}
}

func (l *Loader[K]) updateExpiry(key K) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.expiry[key] = time.Now().Add(300 * time.Millisecond)
	go l.startExpiryTimer(key)
}

func (l *Loader[K]) startExpiryTimer(key K) {
	time.Sleep(300 * time.Millisecond)
	l.mu.Lock()
	defer l.mu.Unlock()
	if expiryTime, ok := l.expiry[key]; ok && time.Now().After(expiryTime) {
		delete(l.cache, key)
		delete(l.expiry, key)
	}
}

// Load loads a single item by key
func (l *Loader[K]) Load(key K) (map[string]interface{}, error) {
	ptr, err := l.loadThunk(key)()
	if err != nil {
		return nil, err
	}
	l.updateExpiry(key)
	return ptr, nil
}

// loadThunk creates a thunk for loading a single item
func (l *Loader[K]) loadThunk(key K) func() (map[string]interface{}, error) {
	l.mu.Lock()
	if d, ok := l.cache[key]; ok {
		l.mu.Unlock()
		return func() (map[string]interface{}, error) { return d[0], nil }
	}

	if l.batch == nil {
		l.batch = &LoaderBatch[K]{done: make(chan struct{})}
	}
	batch := l.batch
	pos := batch.keyIndex(l, key)
	l.mu.Unlock()
	return func() (map[string]interface{}, error) {
		<-batch.done

		if pos < len(batch.results) {
			res := batch.results[pos]
			if len(batch.errors) == 0 || batch.errors[pos] == nil {
				l.set(key, res)
				return res[0], nil
			}
			return res[0], batch.errors[pos]
		}

		return nil, &errors.DataloaderError{Msg: fmt.Sprintf("key %v not found", key)}
	}
}

// LoadList loads a list of items by key
func (l *Loader[K]) LoadList(key K) ([]map[string]interface{}, error) {
	ptr, err := l.loadListThunk(key)()
	if err != nil {
		return nil, err
	}
	l.updateExpiry(key)
	return ptr, nil
}

func (l *Loader[K]) loadListThunk(key K) func() ([]map[string]interface{}, error) {
	l.mu.Lock()
	if d, ok := l.cache[key]; ok {
		l.mu.Unlock()
		return func() ([]map[string]interface{}, error) { return d, nil }
	}

	if l.batch == nil {
		l.batch = &LoaderBatch[K]{done: make(chan struct{})}
	}
	batch := l.batch
	pos := batch.keyIndex(l, key)
	l.mu.Unlock()
	return func() ([]map[string]interface{}, error) {
		<-batch.done

		if pos < len(batch.results) {
			res := batch.results[pos]
			if len(batch.errors) == 0 || batch.errors[pos] == nil {
				l.set(key, res)
				return res, nil
			}
			return res, batch.errors[pos]
		}

		return nil, &errors.DataloaderError{Msg: fmt.Sprintf("key %v not found", key)}
	}
}

// LoadAll loads multiple items by keys
func (l *Loader[K]) LoadAll(keys []K) ([]map[string]interface{}, []error) {
	thunks := make([]func() (map[string]interface{}, error), len(keys))
	for i, key := range keys {
		thunks[i] = l.loadThunk(key)
	}

	results := make([]map[string]interface{}, len(keys))
	errors := make([]error, len(keys))
	for i, thunk := range thunks {
		results[i], errors[i] = thunk()
	}
	return results, errors
}

// Prime adds the provided key-value pair to the cache
func (l *Loader[K]) Prime(key K, value []map[string]interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if _, ok := l.cache[key]; !ok {
		l.unsafeSet(key, value)
		return true
	}
	return false
}

// Clear removes the value associated with key from the cache
func (l *Loader[K]) Clear(key K) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.cache, key)
}

// set adds the provided key-value pair to the cache
func (l *Loader[K]) set(key K, value []map[string]interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.unsafeSet(key, value)
}

// unsafeSet adds the provided key-value pair to the cache without locking
func (l *Loader[K]) unsafeSet(key K, value []map[string]interface{}) {
	l.cache[key] = value
}

// keyIndex returns the index of the given key in the batch, adding it if not present
func (b *LoaderBatch[K]) keyIndex(l *Loader[K], key K) int {
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
func (b *LoaderBatch[K]) startTimer(l *Loader[K]) {
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
func (b *LoaderBatch[K]) end(l *Loader[K]) {
	b.results, b.errors = l.fetch(b.keys)
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cache = make(map[K][]map[string]interface{})
	close(b.done)
}

// ClearAll removes all entries from the cache
func (l *Loader[K]) ClearAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.cache = make(map[K][]map[string]interface{})
}

// CacheSize returns the number of entries in the cache
func (l *Loader[K]) CacheSize() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return len(l.cache)
}
