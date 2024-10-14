package model

import "time"

type LoaderConfig[T any] struct {
	MaxBatch int
	Wait     time.Duration
	Fetch    func(keys []int64) ([]T, []error)
}
