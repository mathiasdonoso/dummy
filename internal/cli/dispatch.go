package cli

import (
	"errors"
	"fmt"
)

func Dispatch(args []string) error {
	if len(args) == 0 {
		return errors.New("no command provided")
	}

	node, ok := CommandRegistry[args[0]]
	if !ok {
		return fmt.Errorf("unknown command: %s", args[0])
	}

	return walk(node, args[1:])
}

func walk(node *CommandNode, args []string) error {
	if len(args) == 0 {
		if node.Handler == nil {
			return fmt.Errorf("command '%s' requires subcommand", node.Name)
		}
		return node.Handler([]string{})
	}

	next := args[0]
	sub, ok := node.Subcommands[next]
	if ok {
		return walk(sub, args[1:])
	}

	if node.Handler == nil {
		return fmt.Errorf("unknown subcommand '%s' for '%s'", next, node.Name)
	}

	return node.Handler(args)
}
