func main() {
	cmd := &cmd.Command{}
	if err := command.Run(cmd, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
