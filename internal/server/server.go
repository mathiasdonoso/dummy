package server

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"text/tabwriter"
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

	for range 50 {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			conn.Close()
			return nil
		}
		time.Sleep(20 * time.Millisecond)
	}

	return fmt.Errorf("server not responding on port %d", port)
}

// For testing purposes
func (s *server) Start(model model.ImportResult) error {
	mux := s.buildMux(model)

	go func() {
		_ = http.ListenAndServe(fmt.Sprintf(":%d", s.Port), mux)
	}()

	return waitUntilListening(s.Port)
}

func (s *server) StartAndBlock(model model.ImportResult) error {
	mux := s.buildMux(model)

	fmt.Printf("Starting server at localhost:%d\n", s.Port)
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), mux)
}

func printEndpoints(model model.ImportResult) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)

	fmt.Fprintln(w, "METHOD\tPATH\tSTATUS")
	for _, e := range model.Endpoints {
		fmt.Fprintf(
			w,
			"%s\t%s\t%d\n",
			e.Method,
			e.Path,
			e.Response.StatusCode,
		)
	}

	w.Flush()
	fmt.Println()
}

func (s *server) buildMux(model model.ImportResult) *http.ServeMux {
	mux := http.NewServeMux()

	seen := []string{}

outer:
	for _, e := range model.Endpoints {
		endpoint := fmt.Sprintf("%s %s", e.Method, e.Path)

		for _, s := range seen {
			if s == endpoint {
				continue outer
			}
		}

		seen = append(seen, endpoint)

		mux.HandleFunc(endpoint, func(w http.ResponseWriter, r *http.Request) {
			slog.Debug(fmt.Sprintf("configuring endpoint %s", endpoint))

			res := e.Response
			for k, v := range res.Headers {
				w.Header().Set(k, v)
			}
			w.WriteHeader(res.StatusCode)
			w.Write(res.Body)
		})
	}

	printEndpoints(model)

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
