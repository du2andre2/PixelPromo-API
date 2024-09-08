package controller

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"pixelPromo/domain/model"
	"pixelPromo/domain/service"
	"strings"
)

type Controller interface {
	CreateUser(*gin.Context)
	CreateInteraction(*gin.Context)
	GetInteractionByID(*gin.Context)
	UpdateUserPicture(ctx *gin.Context)
	GetUserByID(ctx *gin.Context)
	CreatePromotion(*gin.Context)
	UpdatePromotionImage(ctx *gin.Context)
	GetPromotionByID(ctx *gin.Context)
	GetPromotions(ctx *gin.Context)
	GetPromotionByCategory(ctx *gin.Context)
	GetCategories(ctx *gin.Context)
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

func (r *controller) GetInteractionByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	promotion, err := r.repository.GetInteractionByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	if promotion == nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	ctx.IndentedJSON(http.StatusOK, promotion)
	return
}

func (r *controller) CreateInteraction(ctx *gin.Context) {

	var interaction model.PromotionInteraction
	err := ctx.ShouldBindJSON(&interaction)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	err = r.repository.CreateInteraction(ctx, &interaction)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, interaction)
}

func (r *controller) CreateUser(ctx *gin.Context) {

	var user model.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	err = r.repository.CreateUser(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, user)
}

func (r *controller) UpdateUserPicture(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	f, fh, err := ctx.Request.FormFile("picture")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	if fh.Size <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.Error{Err: errors.New("file length <= 0")})
		return
	}

	buffer := make([]byte, fh.Size)
	f.Read(buffer)

	fileBytes := bytes.NewReader(buffer)

	err = r.repository.UpdateUserPicture(ctx, id, fileBytes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *controller) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	user, err := r.repository.GetUserByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	if user == nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	ctx.IndentedJSON(http.StatusOK, user)
}

func (r *controller) CreatePromotion(ctx *gin.Context) {

	var promotion model.Promotion
	err := ctx.ShouldBindJSON(&promotion)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	err = r.repository.CreatePromotion(ctx, &promotion)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, promotion)
}

func (r *controller) UpdatePromotionImage(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	f, fh, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	if fh.Size <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.Error{Err: errors.New("file length <= 0")})
		return
	}

	buffer := make([]byte, fh.Size)
	f.Read(buffer)

	fileBytes := bytes.NewReader(buffer)

	err = r.repository.UpdatePromotionImage(ctx, id, fileBytes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *controller) GetPromotionByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	promotion, err := r.repository.GetPromotionByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	if promotion == nil {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	ctx.IndentedJSON(http.StatusOK, promotion)
	return
}

func (r *controller) GetPromotions(ctx *gin.Context) {

	categories, _ := ctx.GetQueryArray("category")
	search, _ := ctx.GetQuery("search")

	params := model.PromotionQuery{
		Search:     search,
		Categories: categories,
	}
	promotions, err := r.repository.GetPromotions(ctx, params)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	if promotions == nil || len(promotions) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	ctx.IndentedJSON(http.StatusOK, promotions)
	return
}

func (r *controller) GetPromotionByCategory(ctx *gin.Context) {
	category := ctx.Param("category")

	if len(strings.TrimSpace(category)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	promotions, err := r.repository.GetPromotionByCategory(ctx, category)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	if promotions == nil || len(promotions) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	ctx.IndentedJSON(http.StatusOK, promotions)
	return
}

func (r *controller) GetCategories(ctx *gin.Context) {

	categories, err := r.repository.GetCategories(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	if categories == nil || len(categories) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	ctx.IndentedJSON(http.StatusOK, categories)
	return
}
