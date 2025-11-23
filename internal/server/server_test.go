package server

import (
	"bytes"
	"encoding/json"
	"fmt"
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
				Response: model.MockResponse{
					StatusCode: 200,
					Body:       testutils.MustReadFile(t, "test_data/auth-200.json"),
					Headers:    map[string]string{},
					DelayMs:    0,
				},
				Headers:     map[string]string{},
				QueryParams: map[string]string{},
				Body:        "{\"username\": \"username\",\"password\": \"password\"}",
			},
			// {
			// 	Method: "POST",
			// 	Path:   "/api/auth",
			// 	Response: model.MockResponse{
			// 		StatusCode: 400,
			// 		Body:       testutils.MustReadFile(t, "test_data/auth-400.json"),
			// 		Headers:    map[string]string{},
			// 		DelayMs:    0,
			// 	},
			// 	Headers:     map[string]string{},
			// 	QueryParams: map[string]string{},
			// 	Body:        "{\"username\": \"wrong\",\"password\": \"wrong\"}",
			// },
			{
				Method: "GET",
				Path:   "/api/v2.0/projects",
				Response: model.MockResponse{
					StatusCode: 200,
					Body:       testutils.MustReadFile(t, "test_data/harbor_response_projects.json"),
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
				Method: "GET",
				Path:   "/api/v2.0/projects/someproject/repositories",
				Response: model.MockResponse{
					StatusCode: 200,
					Body:       testutils.MustReadFile(t, "test_data/harbor_response_repositories.json"),
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
				Method: "GET",
				Path:   "/api/v2.0/projects/someproject/repositories/somerepository/artifacts",
				Response: model.MockResponse{
					StatusCode: 200,
					Body:       testutils.MustReadFile(t, "test_data/harbor_response_artifacts.json"),
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
		name    string
		model   model.ImportResult
		wantErr bool
	}{
		{"api", apiResults, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewServer()
			err := s.Start(tt.model)

			if tt.wantErr && err == nil {
				t.Errorf("expected error but got nil")
			}

			if !tt.wantErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			for _, e := range tt.model.Endpoints {
				jsonData, err := json.Marshal(e.Body)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				req, err := http.NewRequest(e.Method, fmt.Sprintf("%s:%d%s", "http://localhost", s.Port, e.Path), bytes.NewBuffer(jsonData))
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				res, err := http.DefaultClient.Do(req)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if res.StatusCode != e.Response.StatusCode {
					t.Errorf("unexpected status code %d, wanted %d", res.StatusCode, e.Response.StatusCode)
				}

				responseBody, err := io.ReadAll(res.Body)
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}

				if diff := cmp.Diff(string(responseBody), string(e.Response.Body)); diff != "" {
					t.Errorf("output mismatch (-want +got):\n%s", diff)
				}
			}
		})
	}
}
