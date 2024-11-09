package http

import (
	"bytes"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"net/http"
	"pixelPromo/domain/model"
	"pixelPromo/domain/port"
	"strconv"
	"strings"
	"time"
)

type Controller struct {
	handler port.Handler
}

func NewController(
	handler port.Handler,
) *Controller {
	return &Controller{
		handler: handler,
	}
}

func (r *Controller) GetInteractionStatisticsByPromotionID(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	counters, err := r.handler.GetInteractionStatisticsByPromotionID(ctx, id)
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

func (r *Controller) GetCommentsByPromotionID(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	promotion, err := r.handler.GetCommentsByPromotionID(ctx, id)
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

func (r *Controller) CreateInteraction(ctx *gin.Context) {

	var interaction model.PromotionInteraction
	err := ctx.ShouldBindJSON(&interaction)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	err = r.handler.CreateInteraction(ctx, &interaction)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, interaction)
}

func (r *Controller) CreateUser(ctx *gin.Context) {

	var user model.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	err = r.handler.CreateUser(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusCreated, user)
}

func (r *Controller) UpdateUserPicture(ctx *gin.Context) {
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

	err = r.handler.UpdateUserPicture(ctx, id, fileBytes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *Controller) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	user, err := r.handler.GetUserByID(ctx, id)
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

func (r *Controller) GetUserRank(ctx *gin.Context) {
	limitStr, _ := ctx.GetQuery("limit")

	if len(strings.TrimSpace(limitStr)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
	}

	users, err := r.handler.GetUserRank(ctx, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	if users == nil {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	ctx.IndentedJSON(http.StatusOK, users)
}

func (r *Controller) Login(ctx *gin.Context) {
	var login model.Login
	err := ctx.ShouldBindJSON(&login)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	user, err := r.handler.Login(ctx, &login)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	if user == nil {
		ctx.Writer.WriteHeader(http.StatusUnauthorized)
		return
	}

	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &Claims{
		Username: login.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
		return
	}

	type loginStr struct {
		Token string      `json:"token"`
		User  *model.User `json:"user"`
	}

	response := loginStr{
		Token: tokenString,
		User:  user,
	}
	ctx.JSON(http.StatusOK, response)
}

func (r *Controller) CreatePromotion(ctx *gin.Context) {

	var promotion model.Promotion
	err := ctx.ShouldBindJSON(&promotion)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	err = r.handler.CreatePromotion(ctx, &promotion)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	ctx.IndentedJSON(http.StatusOK, promotion)
}

func (r *Controller) UpdatePromotionImage(ctx *gin.Context) {
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

	err = r.handler.UpdatePromotionImage(ctx, id, fileBytes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"Err": err.Error()})
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *Controller) GetPromotionByID(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	promotion, err := r.handler.GetPromotionByID(ctx, id)
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

func (r *Controller) GetFavoritesPromotionsByUserID(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNoContent)
		return
	}

	promotion, err := r.handler.GetFavoritesPromotionsByUserID(ctx, id)
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

func (r *Controller) GetPromotions(ctx *gin.Context) {

	categories, _ := ctx.GetQueryArray("category")
	search, _ := ctx.GetQuery("search")
	limit, _ := ctx.GetQuery("limit")
	var limitInt int
	if limit != "" {
		limitInt, _ = strconv.Atoi(limit)
	}

	params := &model.PromotionQuery{
		Search:     search,
		Categories: categories,
		Limit:      int32(limitInt),
	}
	promotions, err := r.handler.GetPromotions(ctx, params)
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

func (r *Controller) GetCategories(ctx *gin.Context) {

	categories, err := r.handler.GetCategories(ctx)
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
