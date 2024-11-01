
func Migrate() error {
	return model.GetDB().AutoMigrate(
  )
}