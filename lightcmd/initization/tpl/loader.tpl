var migrateModels = []interface{}{}

func main() {
	option := gormschema.WithConfig(&gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		IgnoreRelationshipsWhenMigrating:         true,
	})
	stmts, err := gormschema.New("mysql", option).Load(migrateModels...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
