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

	if len(p.Item) == 0 {
		return &model.ImportResult{}, fmt.Errorf("postman collection has no items to parse")
	}

	r.Endpoints = []model.Endpoint{}
	for _, it := range p.Item {
		for _, resp := range it.Response {
			response := model.MockResponse{
				StatusCode: resp.Code,
				Body:       []byte(resp.Body),
				Headers:    map[string]string{},
				DelayMs:    0,
			}

			// TODO: refactor the url parser
			uf := strings.ReplaceAll(resp.Originalrequest.Url.Raw, "{{", "")
			uf = strings.ReplaceAll(uf, "}}", "")
			u, err := url.Parse(uf)
			if err != nil {
				return &model.ImportResult{}, err
			}
			querySplit := strings.Split(u.RawQuery, "&")
			query := make(map[string]string)
			for _, q := range querySplit {
				kv := strings.Split(q, "=")
				if len(kv) == 2 {
					query[kv[0]] = kv[1]
				}
			}

			r.Endpoints = append(r.Endpoints, model.Endpoint{
				Method:      resp.Originalrequest.Method,
				Path:        u.Path,
				Description: "",
				Response:    response,
				Headers:     map[string]string{},
				QueryParams: query,
			})
		}
	}

	slog.Debug(fmt.Sprintf("parsing completed for %s, total of endpoints: %d", r.ServiceName, len(r.Endpoints)))

	return &r, nil
}
