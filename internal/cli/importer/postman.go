package importer

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/mathiasdonoso/dummy/internal/postman"
	"github.com/mathiasdonoso/dummy/internal/templates"
)

func ImportPostmanHandler(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing Postman collection file")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	file := args[0]
	slog.Debug(fmt.Sprintf("importing Postman collection: %s", file))

	d, err := os.ReadFile(file)
	if err != nil {
		return nil
	}

	result, err := postman.Parse(d)
	if err != nil {
		return err
	}

	tg := templates.TemplateGenerator{
		Path: filepath.Join(home, ".dummy", result.ServiceName),
	}
	err = tg.Build(result)
	if err != nil {
		return err
	}

	return nil
}
