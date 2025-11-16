package templates

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mathiasdonoso/dummy/internal/model"
	testutils "github.com/mathiasdonoso/dummy/pkg/test_utils"
)

type endpointConfig struct {
	jsonFileContent string
	metaFileContent string
}

func assertFileContent(t *testing.T, gotPath, wantPath string) {
	t.Helper()

	got, err := os.ReadFile(gotPath)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", gotPath, err)
	}

	want, err := os.ReadFile(wantPath)
	if err != nil {
		t.Fatalf("failed to read file %s: %v", wantPath, err)
	}

	if diff := cmp.Diff(string(want), string(got)); diff != "" {
		t.Errorf("content mismatch for %s (-want +got):\n%s", gotPath, diff)
	}
}

func TestProjectGeneration(t *testing.T) {
	tmpDir := t.TempDir()
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
				Path:        "/api/v2.0/projects/someproject/repositories",
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
				Path:        "/api/v2.0/projects/someproject/repositories/somerepository/artifacts",
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

	tests := []struct {
		name    string
		model   model.ImportResult
		want    map[string]endpointConfig
		wantErr bool
	}{
		{
			"harbor api basic import",
			harborResult,
			map[string]endpointConfig{
				"GET_/api/v2.0/projects": {
					jsonFileContent: "./test_data/result_projects_response.json",
					metaFileContent: "./test_data/result_projects_meta.yaml",
				},
				"GET_/api/v2.0/projects/someproject/repositories": {
					jsonFileContent: "./test_data/result_repositories_response.json",
					metaFileContent: "./test_data/result_repositories_meta.yaml",
				},
				"GET_/api/v2.0/projects/someproject/repositories/somerepository/artifacts": {
					jsonFileContent: "./test_data/result_artifacts_response.json",
					metaFileContent: "./test_data/result_artifacts_meta.yaml",
				},
			},
			false,
		},
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
				for _, e := range tt.model.Endpoints {
					folderName := fmt.Sprintf("%s_%s", e.Method, e.Path)
					responseFile := filepath.Join(tmpDir, folderName, "response.json")
					metaFile := filepath.Join(tmpDir, folderName, "meta.yaml")

					assertFileContent(t, responseFile, tt.want[folderName].jsonFileContent)
					assertFileContent(t, metaFile, tt.want[folderName].metaFileContent)
				}
			}
		})
	}
}
