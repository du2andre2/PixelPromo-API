package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"pixelPromo/domain/service"
	"strings"
)

type Controller interface {
	GetUser(*gin.Context)
	GetUserByID(ctx *gin.Context)
}

type controller struct {
	repository service.Repository
}

func NewController(
	repository service.Repository,
) Controller {
	return &controller{
		repository: repository,
	}
}

func (r *controller) GetUser(ctx *gin.Context) {
	user, err := r.repository.GetUser()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.Error{Err: err})
		return
	}

	if user == nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
		return
	}

	ctx.IndentedJSON(http.StatusOK, user)
	return

}

func (r *controller) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	user, err := r.repository.GetUserByID(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.Error{Err: err})
		return
	}

	if user == nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	ctx.IndentedJSON(http.StatusOK, user)
	return

}
