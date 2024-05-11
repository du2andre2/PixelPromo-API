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
	UpdateUserPicture(ctx *gin.Context)
	GetUserByID(ctx *gin.Context)
	CreatePromotion(*gin.Context)
	UpdatePromotionImage(ctx *gin.Context)
	GetPromotionByID(ctx *gin.Context)
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

func (r *controller) CreateUser(ctx *gin.Context) {

	var user model.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.Error{Err: err})
		return
	}

	err = r.repository.CreateUser(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.Error{Err: err})
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *controller) UpdateUserPicture(ctx *gin.Context) {
	id := ctx.Param("id")

	if len(strings.TrimSpace(id)) == 0 {
		ctx.Writer.WriteHeader(http.StatusNotFound)
		return
	}

	f, fh, err := ctx.Request.FormFile("picture")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.Error{Err: err})
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
		ctx.JSON(http.StatusInternalServerError, gin.Error{Err: err})
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

func (r *controller) CreatePromotion(ctx *gin.Context) {

	var promotion model.Promotion
	err := ctx.ShouldBindJSON(&promotion)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.Error{Err: err})
		return
	}

	err = r.repository.CreatePromotion(ctx, &promotion)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.Error{Err: err})
		return
	}

	ctx.Status(http.StatusOK)
}

func (r *controller) UpdatePromotionImage(ctx *gin.Context) {
	f, fh, err := ctx.Request.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.Error{Err: err})
		return
	}

	if fh.Size <= 0 {
		ctx.JSON(http.StatusBadRequest, gin.Error{Err: errors.New("file length <= 0")})
		return
	}

	buffer := make([]byte, fh.Size)
	f.Read(buffer)

	fileBytes := bytes.NewReader(buffer)

	err = r.repository.UpdatePromotionImage(ctx, fileBytes)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.Error{Err: err})
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

	user, err := r.repository.GetPromotionByID(id)
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

//func (r *controller) PutFile(ctx *gin.Context) {
//	if !strings.Contains(ctx.GetHeader("Content-Type"), "multipart/form-data") {
//		r.handleFileInBody(ctx)
//		return
//	}
//
//	r.handleFileInForm(ctx)
//}
//
//func (r *controller) handleFileInBody(ctx *gin.Context) {
//
//	if ctx.Request.ContentLength <= 0 {
//		ctx.JSON(http.StatusBadRequest, gin.Error{Err: errors.New("content length <= 0")})
//		return
//	}
//
//	f, err := getFile("")
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.Error{Err: err})
//		return
//	}
//	defer f.Close()
//
//	written, err := io.Copy(f, ctx.Request.Body)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.Error{Err: err})
//		return
//	}
//
//	ctx.Status(http.StatusOK)
//
//	log.Println("Written", written)
//}
//
//func (r *controller) handleFileInForm(ctx *gin.Context) {
//	f, fh, err := ctx.Request.FormFile("file")
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.Error{Err: err})
//		return
//	}
//
//	if fh.Size <= 0 {
//		ctx.JSON(http.StatusBadRequest, gin.Error{Err: errors.New("file length <= 0")})
//		return
//	}
//
//	buffer := make([]byte, fh.Size)
//	f.Read(buffer)
//
//	fileBytes := bytes.NewReader(buffer)
//
//	err = r.repository.PutOffer(ctx, fileBytes)
//	if err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.Error{Err: err})
//		return
//	}
//
//	ctx.Status(http.StatusOK)
//
//}
//
//func getFile(fname string) (*os.File, error) {
//	var fileName string
//
//	now := time.Now()
//	if fname != "" {
//		fileName = strconv.Itoa(int(now.Unix())) + "_" + fname
//	} else {
//		fileName = "temp_" + strconv.Itoa(int(now.Unix())) + ".txt"
//	}
//
//	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
//	if err != nil {
//		log.Println("create file error", err)
//		return nil, err
//	}
//
//	return f, nil
//
//}
