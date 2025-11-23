package server

import (
	"bytes"
	"io"
	"net/http"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/mathiasdonoso/dummy/internal/model"
	testutils "github.com/mathiasdonoso/dummy/pkg/test_utils"
)

func TestLocalServerEndpoints(t *testing.T) {
	apiResults := model.ImportResult{
		ServiceName: "api",
		Endpoints: []model.Endpoint{
			{
				Method: "POST",
				Path:   "/api/auth",
				Responses: []model.MockResponse{
					{
						StatusCode:  200,
						Body:        testutils.MustReadFile(t, "test_data/auth-200.json"),
						Headers:     map[string]string{},
						DelayMs:     0,
						RequestBody: "{\"username\": \"username\",\"password\": \"password\"}",
					},
					{
						StatusCode:  400,
						Body:        testutils.MustReadFile(t, "test_data/auth-400.json"),
						Headers:     map[string]string{},
						DelayMs:     0,
						RequestBody: "{\"username\": \"wrong\",\"password\": \"wrong\"}",
					},
				},
				Headers:     map[string]string{},
				QueryParams: map[string]string{},
			},
			{
				Method: "GET",
				Path:   "/api/v2.0/projects",
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
				Method: "GET",
				Path:   "/api/v2.0/projects/someproject/repositories",
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
				Method: "GET",
				Path:   "/api/v2.0/projects/someproject/repositories/somerepository/artifacts",
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
		wantErr bool
	}{
		{"api", apiResults, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer()
			ts := s.StartTest(tt.model)
			defer ts.Close()

			for _, e := range tt.model.Endpoints {
				for _, r := range e.Responses {
					jsonData := []byte(r.RequestBody)

					url := ts.URL + e.Path
					req, err := http.NewRequest(e.Method, url, bytes.NewBuffer(jsonData))
					if err != nil {
						t.Errorf("unexpected error: %v", err)
					}

					res, err := http.DefaultClient.Do(req)
					if err != nil {
						t.Errorf("unexpected error: %v", err)
					}
					defer res.Body.Close()

					if res.StatusCode != r.StatusCode {
						t.Errorf("unexpected status code %d, wanted %d", res.StatusCode, r.StatusCode)
					}

					responseBody, err := io.ReadAll(res.Body)
					if err != nil {
						t.Errorf("unexpected error: %v", err)
					}

					if diff := cmp.Diff(string(r.Body), string(responseBody)); diff != "" {
						t.Errorf("output mismatch (-want +got):\n%s", diff)
					}
				}
			}
		})
	}
}
