package postman

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/mathiasdonoso/dummy/internal/model"
)

func Parse(data []byte) (*model.ImportResult, error) {
	var p PostmanJSON
	err := json.Unmarshal(data, &p)
	if err != nil {
		return &model.ImportResult{}, err
	}

	slog.Debug(fmt.Sprintf("top level items found for %s: %d", p.Info.Name, len(p.Item)))
	var r model.ImportResult
	r.ServiceName = p.Info.Name
	r.Endpoints = make([]model.Endpoint, len(p.Item))

	if len(p.Item) == 0 {
		return &model.ImportResult{}, fmt.Errorf("postman collection has no items to parse")
	}

	for i, it := range p.Item {
		responses := make([]model.MockResponse, len(it.Response))

		for idx, resp := range it.Response {
			responses[idx] = model.MockResponse{
				StatusCode: resp.Code,
				Body:       []byte(resp.Body),
				Headers:    map[string]string{},
				DelayMs:    0,
			}
		}

		u, err := url.Parse(it.Request.Url.Raw)
		if err != nil {
			return &model.ImportResult{}, err
		}
		querySplit := strings.Split(u.RawQuery, "&")
		query := make(map[string]string)
		for _, q := range querySplit {
			kv := strings.Split(q, "=")
			query[kv[0]] = kv[1]
		}

		r.Endpoints[i] = model.Endpoint{
			Method:      it.Request.Method,
			Path:        u.Path,
			Description: "",
			Responses:   responses,
			Headers:     map[string]string{},
			QueryParams: query,
		}
	}

	slog.Debug(fmt.Sprintf("parsing completed for %s, total of endpoints: %d", r.ServiceName, len(r.Endpoints)))

	return &r, nil
}

func Validate(data []byte) error {
	return nil
}
