func main() {
  AddCustomImport()
	cmd := &cmd.Command{}
	if err := command.Run(cmd, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}

func AddCustomImport() {
	packages := []string{
		"cmd",
		"schema",
		"service",
		"models",
		"resolver",
	}
	for _, pkgName := range packages {
		template.AddImportRegex(fmt.Sprintf(`%s\.`, pkgName), fmt.Sprintf("%s/%s", "{{ .Module }}", pkgName), "")
	}
}