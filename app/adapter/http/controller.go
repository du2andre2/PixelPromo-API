package http

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"pixelPromo/domain/model"
	"pixelPromo/domain/port"
	"strings"
)

type Controller interface {
	CreateUser(*gin.Context)
	CreateInteraction(*gin.Context)
	GetInteractionByID(*gin.Context)
	GetCommentsByPromotionID(*gin.Context)
	GetInteractionsCountersByPromotionID(*gin.Context)
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
	interactionHandler port.InteractionHandler
	promotionHandler   port.PromotionHandler
	userHandler        port.UserHandler
}

func NewController(
	interactionHandler port.InteractionHandler,
	promotionHandler port.PromotionHandler,
	userHandler port.UserHandler,
) Controller {
	return &controller{
		interactionHandler: interactionHandler,
		promotionHandler:   promotionHandler,
		userHandler:        userHandler,
	}
}

func (r *controller) GetInteractionByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	interaction, err := r.interactionHandler.GetInteractionByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	if interaction == nil {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	ctx.IndentedJSON(http.StatusOK, interaction)
	return
}

func (r *controller) GetInteractionsCountersByPromotionID(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	counters, err := r.interactionHandler.GetInteractionsCountersByPromotionID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	if len(counters) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	ctx.IndentedJSON(http.StatusOK, counters)
	return
}

func (r *controller) GetCommentsByPromotionID(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	promotion, err := r.interactionHandler.GetCommentsByPromotionID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	if promotion == nil {
		ctx.Writer.WriteHeader(http.StatusNoContent)
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

	err = r.interactionHandler.CreateOrUpdateInteraction(ctx, &interaction)
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

	err = r.userHandler.CreateUser(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, user)
}

func (r *controller) UpdateUserPicture(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
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

	err = r.userHandler.UpdateUserPicture(ctx, id, fileBytes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *controller) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	user, err := r.userHandler.GetUserByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	if user == nil {
		ctx.Writer.WriteHeader(http.StatusNoContent)
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

	err = r.promotionHandler.CreatePromotion(ctx, &promotion)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, promotion)
}

func (r *controller) UpdatePromotionImage(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
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

	err = r.promotionHandler.UpdatePromotionImage(ctx, id, fileBytes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *controller) GetPromotionByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	promotion, err := r.promotionHandler.GetPromotionByID(ctx, id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	if promotion == nil {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	ctx.IndentedJSON(http.StatusOK, promotion)
	return
}

func (r *controller) GetPromotions(ctx *gin.Context) {

	categories, _ := ctx.GetQueryArray("category")
	search, _ := ctx.GetQuery("search")

	params := &model.PromotionQuery{
		Search:     search,
		Categories: categories,
	}
	promotions, err := r.promotionHandler.GetPromotions(ctx, params)
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
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	promotions, err := r.promotionHandler.GetPromotionsByCategory(ctx, category)
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

	categories, err := r.promotionHandler.GetCategories(ctx)
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
