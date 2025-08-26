func StartService() {
	port := configs.Config.Port
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := databases.LightDatabaseClient.GetSlaveDB(ctx)
	if err != nil {
		logs.Error().Err(err).Msg("failed to get slave db")
		return
	}

	cfg := graph.Config{
		Resolvers: &resolver.Resolver{
			LDB: databases.LightDatabaseClient,
		},
	}
	cfg.Directives.Auth = auth.AuthDirective

	srv := handler.New(graph.NewExecutableSchema(cfg))
	srv.AddTransport(transport.Websocket{KeepAlivePingInterval: 10 * time.Second})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.MultipartForm{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{Cache: lru.New[string](100)})

	srv.SetRecoverFunc(func(ctx context.Context, err interface{}) error {
		logs.Error().Interface("error", err).Msg("panic")
		return gqlerror.Errorf("服务器繁忙")
	})
	srv.SetErrorPresenter(lighterr.ErrorPresenter)

	router := routers.NewRouter()
	router.Use(auth.Middleware())
	router.Use(dataloader.Middleware(db))
	router.Handle("/", playground.ApolloSandboxHandler("GraphQL playground", "/query"))
	router.Handle("/query", srv)

	logs.Info().Msgf("connect to http://localhost:%s/ for GraphQL playground", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		logs.Error().Err(err).Msg("failed to start server")
	}
}
