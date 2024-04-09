package server

import (
	"github.com/Grey0520/s3proxy/internal/config"
	"github.com/Grey0520/s3proxy/internal/storage"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Echo    *echo.Echo
	Storage *storage.StorageProvider
	Config  *config.Config
}

func NewServer(cfg *config.Config) *Server {
	stg, err := storage.NewStorageProvider(*cfg)
	if err != nil {
		panic(err)
	}
	return &Server{
		Echo:    echo.New(),
		Storage: &stg,
		Config:  cfg,
	}
}

func (server *Server) Start(addr string) error {
	return server.Echo.Start(addr)
}
