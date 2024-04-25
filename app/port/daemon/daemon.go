package daemon

import (
	"pixelPromo/domain/service"
)

type Daemon interface {
	Run()
	Stop()
}

type daemon struct {
	repository service.Repository
}

func NewDaemon(
	repository service.Repository,
) Daemon {
	return &daemon{
		repository: repository,
	}
}

func (r *daemon) Run() {

}

func (r *daemon) Stop() {

}
