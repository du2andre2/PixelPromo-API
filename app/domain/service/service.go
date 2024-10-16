package service

import (
	"pixelPromo/config"
	"pixelPromo/domain/port"
)

func NewService(
	rp port.Repository,
	cfg *config.Config,
	st port.Storage,
	log config.Logger,
) port.Handler {
	return &service{
		rp:  rp,
		cfg: cfg,
		st:  st,
		log: log,
	}
}

type service struct {
	rp  port.Repository
	cfg *config.Config
	st  port.Storage
	log config.Logger
}
