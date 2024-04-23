package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"pixelPromo/domain/service"
	"strings"
)

type Controller interface {
	GetUser(*gin.Context)
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
	}

	if user == nil {
		ctx.IndentedJSON(http.StatusNotFound, gin.H{"message": "album not found"})
	}

	userJson, err := json.Marshal(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.Error{Err: err})
	}

	ctx.JSON(http.StatusOK, userJson)
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

	userJson, err := json.Marshal(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.Error{Err: err})
	}

	ctx.JSON(http.StatusOK, userJson)
	return

}
