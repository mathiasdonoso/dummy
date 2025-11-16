package importer

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/mathiasdonoso/dummy/internal/postman"
)

func ImportPostmanHandler(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("missing Postman collection file")
	}

	file := args[0]
	slog.Debug(fmt.Sprintf("importing Postman collection: %s", file))

	d, err := os.ReadFile(file)
	if err != nil {
		return nil
	}

	_, err = postman.Parse(d)
	if err != nil {
		return err
	}

	return nil
}
