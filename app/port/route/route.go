package route

import (
	"github.com/gin-gonic/gin"
	"pixelPromo/port/controller"
)

type Route interface {
	Setup(severMux *gin.Engine)
}

type route struct {
	controller controller.Controller
}

func NewRoute(
	controller controller.Controller,
) Route {
	return &route{
		controller: controller,
	}
}

func (r *route) Setup(router *gin.Engine) {

	router.GET("/user", r.controller.GetUser)
	router.GET("/user/:id", r.controller.GetUser)

	return
}
