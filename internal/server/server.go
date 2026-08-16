package server

import (
	"ab_system/internal/configs"
	"context"
	"net/http"
	"time"
)

// ИСПОЛЬЗОВАНО СО СТАРОГО ПРОЕКТА Grace full shut down

type Server struct {
	httpServer *http.Server
}

func (s *Server) Run(port string, serverParams configs.ServerParams, handler http.Handler) error {

	s.httpServer = &http.Server{
		Addr:           serverParams.Addr + port,
		Handler:        handler,
		MaxHeaderBytes: serverParams.MaxHeaderMBs * 1024 * 1024,
		ReadTimeout:    time.Duration(serverParams.ReadTimeoutSeconds) * time.Second,
		WriteTimeout:   time.Duration(serverParams.WriteTimeoutSeconds) * time.Second,
	}
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	return s.httpServer.Shutdown(timeoutCtx)
}
