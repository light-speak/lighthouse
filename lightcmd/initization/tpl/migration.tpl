type Migration struct{}

func (c *Migration) Name() string {
	return "migration:apply"
}

func (c *Migration) Usage() string {
	return "This is a command to apply database migrations"
}

func (c *Migration) Args() []*cmd.CommandArg {
	return []*cmd.CommandArg{
		{
			Name:     "env",
			Usage:    "The environment to apply the migrations to",
			Required: true,
			Default:  "dev",
		},
	}
}

func (c *Migration) Action() func(flagValues map[string]interface{}) error {
	return func(flagValues map[string]interface{}) error {
		cwd, err := os.Getwd()
		if err != nil {
			logs.Error().Err(err).Msg("failed to get current working directory")
			return err
		}
		logs.Info().Msgf("Current working directory: %s", cwd)

		workdir, err := atlasexec.NewWorkingDir(
			atlasexec.WithMigrations(os.DirFS("migrations")),
			atlasexec.WithAtlasHCLPath(cwd+"/atlas.hcl"),
		)
		if err != nil {
			logs.Error().Err(err).Msg("failed to create working directory")
			return err
		}
		defer workdir.Close()

		client, err := atlasexec.NewClient(workdir.Path(), "atlas")
		if err != nil {
			logs.Error().Err(err).Msg("failed to create atlas client")
			return err
		}

		logs.Info().Msgf("workdir: %s", workdir.Path())

		env, err := cmd.GetStringArg(flagValues, "env")
		if err != nil {
			logs.Error().Err(err).Msg("failed to get environment")
			return err
		}
		logs.Info().Msgf("env: %s", *env)

		res, err := client.MigrateApply(context.Background(), &atlasexec.MigrateApplyParams{
			Env:    *env,
			DirURL: "file://" + cwd + "/migrations",
		})
		if err != nil {
			logs.Error().Err(err).Msg("failed to apply migrations")
			return err
		}

		logs.Info().Msgf("Applied migrations: %v", res.Applied)
		return nil
	}
}

func (c *Migration) OnExit() func() {
	return func() {}
}

func init() {
	AddCommand(&Migration{})
}
