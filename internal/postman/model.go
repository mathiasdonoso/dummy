package postman

type Query struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Raw struct {
	Language string `json:"language"`
}

type Options struct {
	Raw Raw `json:"raw"`
}

type Body struct {
	Mode    string  `json:"mode"`
	Raw     string  `json:"raw"`
	Options Options `json:"options"`
}

type Originalrequest struct {
	Method string `json:"method"`
	Header []any  `json:"header"`
	Url    Url    `json:"url"`
	Body   Body   `json:"body"`
}

type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Response struct {
	Name                   string          `json:"name"`
	Originalrequest        Originalrequest `json:"originalRequest"`
	Status                 string          `json:"status"`
	Code                   int             `json:"code"`
	PostmanPreviewlanguage string          `json:"_postman_previewlanguage"`
	Header                 []Header        `json:"header"`
	Cookie                 []any           `json:"cookie"`
	Body                   string          `json:"body"`
}

type Info struct {
	PostmanId      string `json:"_postman_id"`
	Name           string `json:"name"`
	Schema         string `json:"schema"`
	ExporterId     string `json:"_exporter_id"`
	CollectionLink string `json:"_collection_link"`
}

type Basic struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

type Auth struct {
	Type  string  `json:"type"`
	Basic []Basic `json:"basic"`
}

type Url struct {
	Raw      string   `json:"raw"`
	Protocol string   `json:"protocol"`
	Host     []string `json:"host"`
	Path     []string `json:"path"`
	Query    []Query  `json:"query"`
}

type Request struct {
	Auth   Auth   `json:"auth"`
	Method string `json:"method"`
	Header []any  `json:"header"`
	Url    Url    `json:"url"`
}

type Item struct {
	Name     string     `json:"name"`
	Request  Request    `json:"request"`
	Response []Response `json:"response"`
}

type PostmanJSON struct {
	Info Info   `json:"info"`
	Item []Item `json:"item"`
}
