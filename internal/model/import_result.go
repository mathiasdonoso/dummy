package model

type ImportResult struct {
	ServiceName string
	Endpoints   []Endpoint
}

type Endpoint struct {
	Method string
	Path   string
	// Body        string
	Description string
	Responses   []MockResponse
	Headers     map[string]string
	QueryParams map[string]string
}

type MockResponse struct {
	RequestBody string
	StatusCode  int
	Body        []byte
	Headers     map[string]string
	DelayMs     int
}
