package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mathiasdonoso/dummy/internal/model"
	testutils "github.com/mathiasdonoso/dummy/pkg/test_utils"
)

func TestProjectGeneration(t *testing.T) {
	harborResult := model.ImportResult{
		ServiceName: "harbor api",
		Endpoints: []model.Endpoint{
			{
				Method:      "GET",
				Path:        "/api/v2.0/projects",
				Description: "",
				Responses: []model.MockResponse{
					{
						StatusCode: 200,
						Body:       testutils.MustReadFile(t, "test_data/harbor_response_projects.json"),
						Headers:    map[string]string{},
						DelayMs:    0,
					},
				},
				Headers: map[string]string{},
				QueryParams: map[string]string{
					"page":      "1",
					"page_size": "100",
				},
			},
			{
				Method:      "GET",
				Path:        "/api/v2.0/projects/onboarding/repositories",
				Description: "",
				Responses: []model.MockResponse{
					{
						StatusCode: 200,
						Body:       testutils.MustReadFile(t, "test_data/harbor_response_repositories.json"),
						Headers:    map[string]string{},
						DelayMs:    0,
					},
				},
				Headers: map[string]string{},
				QueryParams: map[string]string{
					"page":      "1",
					"page_size": "100",
				},
			},
			{
				Method:      "GET",
				Path:        "/api/v2.0/projects/onboarding/repositories/ng-ui-mx/artifacts",
				Description: "",
				Responses: []model.MockResponse{
					{
						StatusCode: 200,
						Body:       testutils.MustReadFile(t, "test_data/harbor_response_artifacts.json"),
						Headers:    map[string]string{},
						DelayMs:    0,
					},
				},
				Headers: map[string]string{},
				QueryParams: map[string]string{
					"page":      "1",
					"page_size": "100",
				},
			},
		},
	}

	tmpDir := t.TempDir()
	tests := []struct {
		name    string
		model   model.ImportResult
		wantErr bool
	}{
		{"harbor api basic import", harborResult, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tg := TemplateGenerator{
				Path: tmpDir,
			}
			err := tg.Build(&tt.model)

			if tt.wantErr && err == nil {
				t.Errorf("expected error but got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.wantErr && err == nil {
				if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
					t.Errorf("subfolder was not created")
				}

				for _, e := range tt.model.Endpoints {
					folderName := fmt.Sprintf("%s_%s", e.Method, e.Path)
					if _, err := os.Stat(filepath.Join(tmpDir, folderName, "response.json")); os.IsNotExist(err) {
						t.Errorf(fmt.Sprintf("%s/response.json was not created", folderName))
					}

					if _, err := os.Stat(filepath.Join(tmpDir, folderName, "meta.yaml")); os.IsNotExist(err) {
						t.Errorf(fmt.Sprintf("%s/meta.yaml was not created", folderName))
					}
				}
			}
		})
	}
}
