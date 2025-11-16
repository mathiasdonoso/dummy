package server

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/mathiasdonoso/dummy/internal/model"
)

type OptFunc func(*Opts)

type Opts struct {
	Port int
}

const DEFAULT_PORT = 8000

func defaultOpts() Opts {
	port := DEFAULT_PORT

	for {
		used, _ := isPortInUse(port)
		if !used {
			break
		}
		port++
	}

	return Opts{
		Port: port,
	}
}

func WithPort(port int) OptFunc {
	return func(opts *Opts) {
		opts.Port = port
	}
}

type server struct {
	Opts
}

func NewServer(opts ...OptFunc) *server {
	o := defaultOpts()

	for _, fn := range opts {
		fn(&o)
	}

	return &server{
		Opts: o,
	}
}

func waitUntilListening(port int) error {
	addr := fmt.Sprintf("localhost:%d", port)

	for i := 0; i < 50; i++ { // retry for ~1 second
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(20 * time.Millisecond)
	}

	return fmt.Errorf("server not responding on port %d", port)
}

func (s *server) Start(model model.ImportResult) error {
	mux := s.buildMux(model)

	go func() {
		_ = http.ListenAndServe(fmt.Sprintf(":%d", s.Port), mux)
	}()

	return waitUntilListening(s.Port)
}

func (s *server) StartAndBlock(model model.ImportResult) error {
	mux := s.buildMux(model)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), mux)
}

func (s *server) buildMux(model model.ImportResult) *http.ServeMux {
	mux := http.NewServeMux()

	for _, ep := range model.Endpoints {
		epCopy := ep
		endpoint := fmt.Sprintf("%s %s", epCopy.Method, epCopy.Path)
		mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
			slog.Debug(fmt.Sprintf("configuring endpoint %s", endpoint))

			res := epCopy.Responses[0]
			for k, v := range res.Headers {
				w.Header().Set(k, v)
			}
			w.WriteHeader(res.StatusCode)
			w.Write(res.Body)
		})
	}

	return mux
}

func isPortInUse(port int) (bool, error) {
	addr := fmt.Sprintf(":%d", port)
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return true, nil
	}

	_ = l.Close()
	return false, nil
}
