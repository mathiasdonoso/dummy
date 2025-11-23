package cli

import "github.com/mathiasdonoso/dummy/internal/cli/importer"

var CommandRegistry = map[string]*CommandNode{
	"run": {
		Name:        "run",
		Description: "Starts a mock server using the format specified by the chosen subcommand",
		Subcommands: map[string]*CommandNode{
			"postman": {
				Name:        "postman",
				Description: "Starts a a mock server directly from a Postman Collection file",
				Handler:     importer.PostmanHandler,
				Subcommands: map[string]*CommandNode{},
				Flags:       map[string]string{},
			},
		},
		Flags: map[string]string{},
	},
}
