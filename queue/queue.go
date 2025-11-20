package queue

import (
	"context"
	"errors"
	"sync"

	"github.com/hibiken/asynq"
	"github.com/light-speak/lighthouse/logs"
)

type JobConfig struct {
	Name     string
	Priority int
	Executor Executor
}

var (
	JobMutex     sync.RWMutex
	JobConfigMap = map[string]JobConfig{}
	client       *asynq.Client
	clientOnce   sync.Once
)

type Executor interface {
	Execute(ctx context.Context, task *asynq.Task) error
}

func StartQueue() error {
	if !LightQueueConfig.Enable {
		return errors.New("queue is not enabled")
	}

	JobMutex.RLock()
	defer JobMutex.RUnlock()

	if len(JobConfigMap) == 0 {
		return errors.New("no job config found")
	}

	queuePriority := map[string]int{}
	for _, job := range JobConfigMap {
		if job.Priority <= 0 {
			return errors.New("job priority must be greater than 0")
		}
		if job.Executor == nil {
			return errors.New("job executor must be not nil")
		}
		queuePriority[job.Name] = job.Priority
	}

	concurrency := 8

	srv := asynq.NewServer(
		getRedisConfig(),
		asynq.Config{Concurrency: concurrency, Queues: queuePriority},
	)

	mux := asynq.NewServeMux()
	for _, job := range JobConfigMap {
		mux.HandleFunc(job.Name, job.Executor.Execute)
	}

	logs.Info().Int("concurrency", concurrency).Int("jobs", len(JobConfigMap)).Msg("starting queue server")

	if err := srv.Run(mux); err != nil {
		logs.Error().Err(err).Msg("queue failed to run")
		return err
	}

	return nil
}

func getRedisConfig() *asynq.RedisClientOpt {
	return &asynq.RedisClientOpt{
		Addr:     LightQueueConfig.Host + ":" + LightQueueConfig.Port,
		Password: LightQueueConfig.Password,
		DB:       LightQueueConfig.DB,
	}
}

func GetClient() (*asynq.Client, error) {
	if !LightQueueConfig.Enable {
		return nil, errors.New("queue is not enabled")
	}
	clientOnce.Do(func() {
		client = asynq.NewClient(getRedisConfig())
		logs.Info().Msg("queue client initialized")
	})
	if client == nil {
		return nil, errors.New("queue client not initialized")
	}
	return client, nil
}

func RegisterJob(name string, config JobConfig) {
	JobMutex.Lock()
	defer JobMutex.Unlock()
	JobConfigMap[name] = config
}

func CloseClient() error {
	if client != nil {
		return client.Close()
	}
	return nil
}
