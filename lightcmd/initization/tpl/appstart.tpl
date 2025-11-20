type Start struct{}

func (c *Start) Name() string {
	return "app:start"
}

func (c *Start) Usage() string {
	return "This is a command to Start the service"
}

func (c *Start) Args() []*cmd.CommandArg {
	return []*cmd.CommandArg{}
}

func (c *Start) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		server.StartService()
		return nil
	}
}

func (c *Start) OnExit() func() {
	return func() {
		logs.Info().Msg("shutting down gracefully...")

		// 关闭数据库连接
		if databases.LightDatabaseClient != nil {
			databases.LightDatabaseClient.CloseConnections()
		}

		// 关闭 Redis 连接
		if lr, err := redis.GetLightRedis(); err == nil && lr != nil {
			if err := lr.Close(); err != nil {
				logs.Error().Err(err).Msg("failed to close redis")
			} else {
				logs.Info().Msg("redis connection closed")
			}
		}

		// 关闭队列客户端
		if err := queue.CloseClient(); err != nil {
			logs.Error().Err(err).Msg("failed to close queue client")
		} else {
			logs.Info().Msg("queue client closed")
		}

		// 关闭 NATS 连接
		if broker := messaging.GetBroker(); broker != nil {
				broker.Close()
				logs.Info().Msg("messaging broker closed")
		}

		logs.Info().Msg("shutdown complete")
	}
}

func init() {
	AddCommand(&Start{})
}
