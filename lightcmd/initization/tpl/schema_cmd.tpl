type Schema struct{}

func (c *Schema) Name() string {
	return "schema"
}

func (c *Schema) Usage() string {
	return "Export schema.graphql file"
}

func (c *Schema) Args() []*cmd.CommandArg {
	return []*cmd.CommandArg{}
}

func (c *Schema) Action() func(flagValues map[string]any) error {
	return func(flagValues map[string]any) error {
		cwd, err := os.Getwd()
		if err != nil {
			logs.Error().Err(err).Msg("failed to get current working directory")
			return err
		}

		// Get ExecutableSchema
		rs := &resolver.Resolver{LDB: nil}
		es := graph.NewExecutableSchema(graph.Config{Resolvers: rs})

		// Get AST Schema
		schema := es.Schema()
		if schema == nil {
			logs.Error().Msg("failed to get schema")
			return nil
		}

		// Format to SDL
		var buf bytes.Buffer
		f := formatter.NewFormatter(&buf)
		f.FormatSchema(schema)

		// Write to file
		gqlFile := filepath.Join(cwd, "schema.graphql")
		if err := os.WriteFile(gqlFile, buf.Bytes(), 0644); err != nil {
			logs.Error().Err(err).Msgf("failed to write file: %s", gqlFile)
			return err
		}

		logs.Info().Msgf("Schema exported to: %s", gqlFile)
		return nil
	}
}

func (c *Schema) OnExit() func() {
	return func() {}
}

func init() {
	AddCommand(&Schema{})
}
