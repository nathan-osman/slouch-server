package server

import (
	"net"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

// Server provides the web interface for the application.
type Server struct {
	listener    net.Listener
	templateSet *pongo2.TemplateSet
	log         *logrus.Entry
	stoppedChan chan bool
}

// New creates a new server instance.
func New(cfg *Config) (*Server, error) {
	l, err := net.Listen("tcp", cfg.Addr)
	if err != nil {
		return nil, err
	}
	var (
		r = mux.NewRouter()
		s = &Server{
			listener:    l,
			templateSet: pongo2.NewSet("", &b0xLoader{}),
			log:         logrus.WithField("context", "server"),
			stoppedChan: make(chan bool),
		}
		server = http.Server{
			Handler: r,
		}
	)
	go func() {
		defer close(s.stoppedChan)
		defer s.log.Info("server has stopped")
		s.log.Info("server has started")
		if err := server.Serve(l); err != nil {
			s.log.Error(err.Error())
		}
	}()
	return s, nil
}

// Close shuts down the web server.
func (s *Server) Close() {
	s.listener.Close()
	<-s.stoppedChan
}
