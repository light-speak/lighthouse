func main() {
	c := &commands.Command{}
	if err := cmd.Run(c, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(1)
	}
}
