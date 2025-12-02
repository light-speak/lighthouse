package health

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"runtime"
	"sync"
	"time"

	"github.com/light-speak/lighthouse/databases"
)

// HealthStatus 健康状态
type HealthStatus struct {
	Status    string                 `json:"status"`     // healthy | degraded | unhealthy
	Timestamp time.Time              `json:"timestamp"`
	Checks    map[string]CheckResult `json:"checks"`
}

// CheckResult 单项检查结果
type CheckResult struct {
	Status  string `json:"status"`            // healthy | unhealthy
	Message string `json:"message,omitempty"`
	Latency string `json:"latency,omitempty"` // 响应时间
}

// Config 健康检查配置
type Config struct {
	// 数据库连接数阈值（超过则降级）
	DBMaxOpenConnsThreshold float64
	// 内存使用阈值（超过则降级），单位 MB
	MemoryThresholdMB uint64
	// 数据库 ping 超时时间
	DBPingTimeout time.Duration
}

var (
	config     *Config
	configOnce sync.Once
)

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	return &Config{
		DBMaxOpenConnsThreshold: 0.8,  // 80% 连接数使用率
		MemoryThresholdMB:       1024, // 1GB
		DBPingTimeout:           3 * time.Second,
	}
}

// SetConfig 设置配置
func SetConfig(c *Config) {
	configOnce.Do(func() {
		config = c
	})
}

func getConfig() *Config {
	if config == nil {
		config = DefaultConfig()
	}
	return config
}

// ReadinessHandler 就绪检查处理器
// 检查服务是否可以接收流量
func ReadinessHandler(w http.ResponseWriter, r *http.Request) {
	cfg := getConfig()
	status := &HealthStatus{
		Status:    "healthy",
		Timestamp: time.Now(),
		Checks:    make(map[string]CheckResult),
	}

	// 检查数据库
	dbCheck := checkDatabase(cfg)
	status.Checks["database"] = dbCheck
	if dbCheck.Status == "unhealthy" {
		status.Status = "unhealthy"
	}

	// 检查内存
	memCheck := checkMemory(cfg)
	status.Checks["memory"] = memCheck
	if memCheck.Status == "unhealthy" && status.Status == "healthy" {
		status.Status = "degraded"
	}

	// 检查数据库连接池使用率
	poolCheck := checkDBPool(cfg)
	status.Checks["db_pool"] = poolCheck
	if poolCheck.Status == "unhealthy" && status.Status == "healthy" {
		status.Status = "degraded"
	}

	// 设置响应
	w.Header().Set("Content-Type", "application/json")
	if status.Status == "unhealthy" {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else if status.Status == "degraded" {
		w.WriteHeader(http.StatusOK) // 降级但仍可用
	} else {
		w.WriteHeader(http.StatusOK)
	}

	json.NewEncoder(w).Encode(status)
}

// LivenessHandler 存活检查处理器
// 只检查进程是否存活
func LivenessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "alive",
		"timestamp": time.Now(),
	})
}

// checkDatabase 检查数据库连接
func checkDatabase(cfg *Config) CheckResult {
	if databases.LightDatabaseClient == nil || !databases.LightDatabaseClient.Completed {
		return CheckResult{
			Status:  "unhealthy",
			Message: "database not initialized",
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), cfg.DBPingTimeout)
	defer cancel()

	start := time.Now()
	db, err := databases.LightDatabaseClient.GetDB(ctx)
	if err != nil {
		return CheckResult{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return CheckResult{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		return CheckResult{
			Status:  "unhealthy",
			Message: "database ping failed: " + err.Error(),
		}
	}

	return CheckResult{
		Status:  "healthy",
		Latency: time.Since(start).String(),
	}
}

// checkMemory 检查内存使用
func checkMemory(cfg *Config) CheckResult {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)

	allocMB := m.Alloc / 1024 / 1024
	if allocMB > cfg.MemoryThresholdMB {
		return CheckResult{
			Status:  "unhealthy",
			Message: "memory usage too high",
		}
	}

	return CheckResult{
		Status:  "healthy",
		Message: formatBytes(m.Alloc),
	}
}

// checkDBPool 检查数据库连接池使用率
func checkDBPool(cfg *Config) CheckResult {
	if databases.LightDatabaseClient == nil || !databases.LightDatabaseClient.Completed {
		return CheckResult{
			Status:  "unhealthy",
			Message: "database not initialized",
		}
	}

	ctx := context.Background()
	db, err := databases.LightDatabaseClient.GetDB(ctx)
	if err != nil {
		return CheckResult{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return CheckResult{
			Status:  "unhealthy",
			Message: err.Error(),
		}
	}

	stats := sqlDB.Stats()
	maxOpen := stats.MaxOpenConnections
	inUse := stats.InUse

	if maxOpen > 0 {
		usage := float64(inUse) / float64(maxOpen)
		if usage > cfg.DBMaxOpenConnsThreshold {
			return CheckResult{
				Status:  "unhealthy",
				Message: "connection pool usage too high",
			}
		}
	}

	return CheckResult{
		Status:  "healthy",
		Message: formatPoolStats(stats.InUse, stats.Idle, stats.MaxOpenConnections),
	}
}

func formatBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}

func formatPoolStats(inUse, idle, maxOpen int) string {
	return fmt.Sprintf("in_use=%d idle=%d max=%d", inUse, idle, maxOpen)
}
