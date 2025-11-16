package cli

import "github.com/mathiasdonoso/dummy/internal/cli/importer"

var CommandRegistry = map[string]*CommandNode{
	"import": {
		Name:        "import",
		Description: "Import external API definitions into dummy templates",
		Subcommands: map[string]*CommandNode{
			"postman": {
				Name:        "postman",
				Description: "Import a Postman collection",
				Handler:     importer.ImportPostmanHandler,
				Subcommands: map[string]*CommandNode{},
				Flags:       map[string]string{},
			},
		},
		Flags: map[string]string{},
	},
}
