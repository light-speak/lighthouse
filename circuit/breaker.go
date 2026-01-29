package circuit

import (
	"sync"
	"time"
)

// State represents the circuit breaker state
type State int

const (
	StateClosed   State = iota // Normal operation
	StateOpen                  // Circuit breaker is open, requests fail fast
	StateHalfOpen              // Testing if service has recovered
)

func (s State) String() string {
	switch s {
	case StateClosed:
		return "closed"
	case StateOpen:
		return "open"
	case StateHalfOpen:
		return "half-open"
	default:
		return "unknown"
	}
}

// Breaker implements the circuit breaker pattern
type Breaker struct {
	mu          sync.RWMutex
	state       State
	failures    int
	threshold   int           // Number of failures to trigger open state
	timeout     time.Duration // Duration to stay in open state
	lastFailure time.Time
}

// NewBreaker creates a new circuit breaker
// threshold: number of consecutive failures before opening the circuit
// timeout: duration to wait before transitioning from open to half-open
func NewBreaker(threshold int, timeout time.Duration) *Breaker {
	return &Breaker{
		state:     StateClosed,
		threshold: threshold,
		timeout:   timeout,
	}
}

// Allow returns true if the request is allowed to proceed
func (b *Breaker) Allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()

	switch b.state {
	case StateClosed:
		return true
	case StateOpen:
		if time.Since(b.lastFailure) > b.timeout {
			// Transition to half-open state
			b.state = StateHalfOpen
			return true
		}
		return false
	case StateHalfOpen:
		return true
	}
	return false
}

// RecordSuccess records a successful operation
// Resets the circuit breaker to closed state
func (b *Breaker) RecordSuccess() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.failures = 0
	b.state = StateClosed
}

// RecordFailure records a failed operation
// May transition to open state if threshold is exceeded
func (b *Breaker) RecordFailure() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.failures++
	b.lastFailure = time.Now()

	if b.failures >= b.threshold {
		b.state = StateOpen
	}
}

// State returns the current state of the circuit breaker
func (b *Breaker) State() State {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.state
}

// Failures returns the current failure count
func (b *Breaker) Failures() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.failures
}

// Reset resets the circuit breaker to closed state
func (b *Breaker) Reset() {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.failures = 0
	b.state = StateClosed
}
