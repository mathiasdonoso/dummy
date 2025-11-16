package templates

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mathiasdonoso/dummy/internal/model"
)

type TemplateGenerator struct {
	Path string
}

func (tg *TemplateGenerator) Build(m *model.ImportResult) error {
	for _, e := range m.Endpoints {
		ep := filepath.Join(tg.Path, e.Method+"_", e.Path)
		fmt.Printf("ep: %s\n", ep)
		if err := os.MkdirAll(ep, os.ModePerm); err != nil {
			return err
		}

		fr, err := os.Create(fmt.Sprintf("%s/response.json", ep))
		if err != nil {
			return err
		}
		fr.Close()

		fm, err := os.Create(fmt.Sprintf("%s/meta.yaml", ep))
		if err != nil {
			return err
		}
		fm.Close()
	}

	return nil
}
