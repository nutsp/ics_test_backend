package server

import (
	"app/config"

	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type Server struct {
	fiber *fiber.App
	db    *sqlx.DB
	cfg   *config.Config
}

func NewServer(db *sqlx.DB, cfg *config.Config) *Server {
	return &Server{
		fiber: fiber.New(),
		db:    db,
		cfg:   cfg,
	}
}

func (s *Server) Run() error {
	if err := s.MapHandlers(s.fiber); err != nil {
		return err
	}

	s.fiber.Listen(":3000")
	return nil
}
