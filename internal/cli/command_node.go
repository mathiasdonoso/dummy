package cli

type CommandNode struct {
	Name        string
	Description string
	Handler     func(args []string) error
	Subcommands map[string]*CommandNode
	Flags       map[string]string
}
