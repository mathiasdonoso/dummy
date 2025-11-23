package model

type ImportResult struct {
	ServiceName string
	Endpoints   []Endpoint
}

type Endpoint struct {
	Method      string
	Path        string
	Body        string
	Description string
	Response    MockResponse
	Headers     map[string]string
	QueryParams map[string]string
}

type MockResponse struct {
	StatusCode int
	Body       []byte
	Headers    map[string]string
	DelayMs    int
}
