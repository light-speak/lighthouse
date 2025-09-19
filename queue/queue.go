package queue

import (
	"context"
	"errors"

	"github.com/hibiken/asynq"
	"github.com/light-speak/lighthouse/logs"
)

type JobConfig struct {
	Name     string
	Priority int
	Executor Executor
}

var (
	JobConfigMap = map[string]JobConfig{}
	client       *asynq.Client
)

type Executor interface {
	Execute(ctx context.Context, task *asynq.Task) error
}

func StartQueue() error {
	if !LightQueueConfig.Enable {
		return errors.New("queue is not enabled")
	}

	queuePriority := map[string]int{}
	concurrency := 0
	for _, job := range JobConfigMap {
		queuePriority[job.Name] = job.Priority
		concurrency += job.Priority
	}

	srv := asynq.NewServer(
		getRedisConfig(),
		asynq.Config{Concurrency: concurrency, Queues: queuePriority},
	)

	mux := asynq.NewServeMux()
	for _, job := range JobConfigMap {
		mux.HandleFunc(job.Name, job.Executor.Execute)
	}

	if err := srv.Run(mux); err != nil {
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

func GetClient() *asynq.Client {
	if !LightQueueConfig.Enable {
		return nil
	}
	if client == nil {
		logs.Warn().Msg("init queue client")
		client = asynq.NewClient(getRedisConfig())
	}
	return client
}
