package health

import (
	"os"
	"sync"
	"time"

	"github.com/light-speak/lighthouse/logs"
)

// SelfHealConfig configures the self-healing behavior
type SelfHealConfig struct {
	// CheckInterval is the interval between health checks
	CheckInterval time.Duration
	// FailureThreshold is the number of consecutive failures before triggering self-heal
	FailureThreshold int
	// GracePeriod is the wait time before exiting to allow cleanup
	GracePeriod time.Duration
	// Enabled determines if self-healing is active
	Enabled bool
}

// DefaultSelfHealConfig returns the default configuration
func DefaultSelfHealConfig() *SelfHealConfig {
	return &SelfHealConfig{
		CheckInterval:    30 * time.Second,
		FailureThreshold: 5, // 5 consecutive failures (~2.5 minutes)
		GracePeriod:      10 * time.Second,
		Enabled:          true,
	}
}

// SelfHealer monitors service health and triggers process exit on sustained failures
// The external process manager (systemd/docker/supervisor) is expected to restart the service
type SelfHealer struct {
	config       *SelfHealConfig
	healthChecks []func() error
	failureCount int
	mu           sync.Mutex
	stop         chan struct{}
	stopped      bool
	exitFunc     func() // Custom exit function, defaults to os.Exit(1)
}

// NewSelfHealer creates a new self-healer with the given configuration
func NewSelfHealer(config *SelfHealConfig) *SelfHealer {
	if config == nil {
		config = DefaultSelfHealConfig()
	}
	return &SelfHealer{
		config:   config,
		stop:     make(chan struct{}),
		exitFunc: func() { os.Exit(1) },
	}
}

// SetExitFunc sets a custom exit function (useful for testing)
func (s *SelfHealer) SetExitFunc(f func()) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.exitFunc = f
}

// RegisterCheck registers a health check function
// The function should return nil if healthy, or an error if unhealthy
func (s *SelfHealer) RegisterCheck(check func() error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.healthChecks = append(s.healthChecks, check)
}

// Start begins the self-heal monitoring
func (s *SelfHealer) Start() {
	if !s.config.Enabled {
		logs.Info().Msg("self-heal is disabled")
		return
	}

	go func() {
		ticker := time.NewTicker(s.config.CheckInterval)
		defer ticker.Stop()

		logs.Info().
			Dur("check_interval", s.config.CheckInterval).
			Int("failure_threshold", s.config.FailureThreshold).
			Msg("self-heal monitoring started")

		for {
			select {
			case <-s.stop:
				logs.Debug().Msg("self-heal monitoring stopped")
				return
			case <-ticker.C:
				s.runChecks()
			}
		}
	}()
}

func (s *SelfHealer) runChecks() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopped {
		return
	}

	// Execute all health checks
	allHealthy := true
	for _, check := range s.healthChecks {
		if err := check(); err != nil {
			logs.Warn().Err(err).Msg("self-heal check failed")
			allHealthy = false
			break
		}
	}

	if allHealthy {
		// Reset failure count
		if s.failureCount > 0 {
			logs.Info().Int("previous_failures", s.failureCount).Msg("service recovered")
		}
		s.failureCount = 0
		return
	}

	s.failureCount++
	logs.Warn().
		Int("failure_count", s.failureCount).
		Int("threshold", s.config.FailureThreshold).
		Msg("self-heal: health check failed")

	if s.failureCount >= s.config.FailureThreshold {
		s.triggerSelfHeal()
	}
}

func (s *SelfHealer) triggerSelfHeal() {
	logs.Error().
		Int("failure_count", s.failureCount).
		Dur("grace_period", s.config.GracePeriod).
		Msg("self-heal triggered: service will exit for restart")

	s.stopped = true

	// Release the lock before sleeping
	s.mu.Unlock()

	// Wait grace period for cleanup
	time.Sleep(s.config.GracePeriod)

	// Exit process, external process manager will restart
	s.exitFunc()
}

// Stop stops the self-heal monitoring
func (s *SelfHealer) Stop() {
	s.mu.Lock()
	if s.stopped {
		s.mu.Unlock()
		return
	}
	s.stopped = true
	s.mu.Unlock()

	close(s.stop)
}

// FailureCount returns the current consecutive failure count
func (s *SelfHealer) FailureCount() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.failureCount
}

// IsEnabled returns whether self-heal is enabled
func (s *SelfHealer) IsEnabled() bool {
	return s.config.Enabled
}
