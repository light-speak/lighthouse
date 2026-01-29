var selfHealer *health.SelfHealer

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
		// 启动自愈监控
		selfHealer = health.NewSelfHealer(nil) // 使用默认配置

		// 注册关键依赖的健康检查
		selfHealer.RegisterCheck(func() error {
			// 数据库检查
			if databases.LightDatabaseClient == nil || !databases.LightDatabaseClient.Completed {
				return errors.New("database not available")
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			db, _ := databases.LightDatabaseClient.GetDB(ctx)
			sqlDB, _ := db.DB()
			return sqlDB.PingContext(ctx)
		})

		selfHealer.RegisterCheck(func() error {
			// Redis 检查（如果启用）
			if !redis.LightRedisConfig.Enable {
				return nil
			}
			client, err := redis.GetLightRedis()
			if err != nil {
				return err
			}
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			_, err = client.Client.Ping(ctx).Result()
			return err
		})

		selfHealer.Start()

		server.StartService()
		return nil
	}
}

func (c *Start) OnExit() func() {
	return func() {
		logs.Info().Msg("shutting down gracefully...")

		// 停止自愈监控
		if selfHealer != nil {
			selfHealer.Stop()
		}

		// 添加总超时保护
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		done := make(chan struct{})
		go func() {
			// 按顺序关闭各组件
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

			close(done)
		}()

		select {
		case <-done:
			logs.Info().Msg("graceful shutdown completed")
		case <-ctx.Done():
			logs.Warn().Msg("graceful shutdown timeout, forcing exit")
		}
	}
}

func init() {
	AddCommand(&Start{})
}
