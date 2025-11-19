package postman

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mathiasdonoso/dummy/internal/model"
	testutils "github.com/mathiasdonoso/dummy/pkg/test_utils"
)

func TestParsingPostmanCollection(t *testing.T) {
	portainerApiResult := model.ImportResult{
		ServiceName: "portainer api",
		Endpoints: []model.Endpoint{
			{
				Method:      "POST",
				Path:        "/api/auth",
				Description: "",
				Response: model.MockResponse{
					StatusCode: 200,
					Body:       testutils.MustReadFile(t, "./test_data/portainer_api/responses/auth-200.json"),
					Headers:    map[string]string{},
					DelayMs:    0,
				},
				Headers:     map[string]string{},
				QueryParams: map[string]string{},
			},
		},
	}

	harborApiResult := model.ImportResult{
		ServiceName: "harbor api",
		Endpoints: []model.Endpoint{
			{
				Method:      "GET",
				Path:        "/api/v2.0/projects",
				Description: "",
				Response: model.MockResponse{
					StatusCode: 200,
					Body:       testutils.MustReadFile(t, "./test_data/harbor_api/responses/projects-200.json"),
					Headers:    map[string]string{},
					DelayMs:    0,
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
				Response: model.MockResponse{
					StatusCode: 200,
					Body:       testutils.MustReadFile(t, "./test_data/harbor_api/responses/repositories-200.json"),
					Headers:    map[string]string{},
					DelayMs:    0,
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
				Response: model.MockResponse{
					StatusCode: 200,
					Body:       testutils.MustReadFile(t, "./test_data/harbor_api/responses/artifacts-200.json"),
					Headers:    map[string]string{},
					DelayMs:    0,
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
		name      string
		inputFile string
		want      *model.ImportResult
		wantErr   bool
	}{
		{"harbor api basic import", "./test_data/harbor_api/harbor.postman_collection.json", &harborApiResult, false},
		{"portainer api basic import", "./test_data/portainer_api/portainer.postman_collection.json", &portainerApiResult, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := testutils.MustReadFile(t, tt.inputFile)
			r, err := Parse(input)

			if tt.wantErr && err == nil {
				t.Errorf("expected error but got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.wantErr && err == nil {
				if len(tt.want.Endpoints) == 0 || len(r.Endpoints) == 0 {
					t.Errorf("expected endpoint to be more than 0")
				}

				if len(tt.want.Endpoints) != len(r.Endpoints) {
					t.Errorf("expected endpoints length to be %d but got %d", len(tt.want.Endpoints), len(r.Endpoints))
				}

				// Strip trailing newline introduced by POSIX text-file semantics.
				// Editors routinely append a final '\n' even when it's not visible,
				// but Postman response bodies don't include it—so we normalize here.
				for i := range tt.want.Endpoints {
					tt.want.Endpoints[i].Response.Body = bytes.TrimRight(tt.want.Endpoints[i].Response.Body, "\n")
				}

				if diff := cmp.Diff(tt.want.Endpoints, r.Endpoints); diff != "" {
					t.Errorf("output mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
