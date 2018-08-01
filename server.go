package sm

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	// Logging
	"github.com/unrolled/logger"

	// Routing
	"github.com/julienschmidt/httprouter"
)

// Options ...
type Options struct {
}

// Server ...
type Server struct {
	bind   string
	server *http.Server

	// Router
	router *httprouter.Router

	// Logger
	logger *logger.Logger
}

// ListenAndServe ...
func (s *Server) ListenAndServe() {
	log.Fatal(s.server.ListenAndServe())
}

func (s *Server) AddRoute(method, path string, handler http.Handler) {
	s.router.Handler(method, path, handler)
}

func (s *Server) Shutdown() {
	if err := s.server.Shutdown(context.Background()); err != nil {
		log.Errorf("error shutting down server: %v", err)
	}
}

func (s *Server) initRoutes() {
	s.router.GET("/", s.IndexHandler())
	s.router.POST("/create", s.CreateHandler())
	s.router.POST("/close/:id", s.CloseHandler())
	s.router.GET("/search", s.SearchHandler())
	s.router.GET("/search/:id", s.SearchHandler())
}

// NewServer ...
func NewServer(bind string, options *Options) *Server {
	router := httprouter.New()

	server := &Server{
		server: &http.Server{
			Addr: bind,
			Handler: logger.New(logger.Options{
				Prefix:               "sm",
				RemoteAddressHeaders: []string{"X-Forwarded-For"},
			}).Handler(router),
		},

		// Router
		router: router,
	}

	server.initRoutes()

	return server
}
