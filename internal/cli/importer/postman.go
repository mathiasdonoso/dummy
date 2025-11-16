package importer

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/mathiasdonoso/dummy/internal/postman"
	"github.com/mathiasdonoso/dummy/internal/server"
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

	m, err := postman.Parse(d)
	if err != nil {
		return err
	}

	s := server.NewServer()
	err = s.StartAndBlock(*m)
	if err != nil {
		return err
	}

	fmt.Printf("Server running at localhost:%d\n", s.Port)

	return nil
}
