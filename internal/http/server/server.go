package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"os"
	"os/signal"
	"vehicles/config"
	"vehicles/infra/logger"
	"time"
)

type Server struct {
	Echo *echo.Echo
}

func (s *Server) Start() {
	e := s.Echo

	go func() {
		e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.App().Port)))
	}()

	GracefulShutdown(e)
}

func New() *Server {
	return &Server{Echo: echo.New()}
}

// GracefulShutdown server will gracefully shut down within 5 sec
func GracefulShutdown(e *echo.Echo) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	logger.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = e.Shutdown(ctx)
	logger.Info("server shutdowns gracefully")
}
