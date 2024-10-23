package port

import (
	"context"
	"io"
)

type Storage interface {
	UploadUserPicture(context.Context, string, io.Reader) (string, error)
	UploadPromotionImage(context.Context, string, io.Reader) (string, error)
}
