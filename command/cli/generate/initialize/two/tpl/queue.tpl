type JobConfig struct {
	Name     string
	Priority int
	Handler  func(ctx bCtx.Context, t *asynq.Task) error
}

var (
	JobConfigMap = map[string]JobConfig{}
)

func StartQueue() {

	queuePriority := map[string]int{}
	for _, job := range JobConfigMap {
		queuePriority[job.Name] = job.Priority
	}

	srv := asynq.NewServer(
		asynq.RedisClientOpt{Addr: env.LighthouseConfig.Redis.Host, Password: env.LighthouseConfig.Redis.Password, DB: env.LighthouseConfig.Redis.Db},
		asynq.Config{Concurrency: 20, Queues: queuePriority},
	)

	mux := asynq.NewServeMux()
	for _, job := range JobConfigMap {
		mux.HandleFunc(job.Name, job.Handler)
	}

	if err := srv.Run(mux); err != nil {
		panic(err)
	}
}

