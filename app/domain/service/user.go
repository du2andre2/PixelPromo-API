package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"pixelPromo/config"
	"pixelPromo/domain/model"
	"pixelPromo/domain/port"
	"regexp"
	"sort"
	"strings"
	"time"
)

func NewUserService(
	rp port.Repository,
	st port.Storage,
	cfg *config.Config,
	log config.Logger,
) port.UserHandler {
	return &userService{
		rp:  rp,
		st:  st,
		cfg: cfg,
		log: log,
	}
}

type userService struct {
	rp  port.Repository
	st  port.Storage
	cfg *config.Config
	log config.Logger
}

func (r *userService) CreateUser(ctx context.Context, user *model.User) error {

	err := r.validUser(user)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	user.CreatedAt = time.Now()
	user.ID = fmt.Sprintf("%d", user.CreatedAt.UnixNano())

	if err = r.rp.CreateOrUpdateUser(ctx, user); err != nil {
		r.log.Error(err.Error())
		return err
	}

	r.log.Debug("user created")
	return nil
}

func (r *userService) UpdateUserPicture(ctx context.Context, id string, image io.Reader) error {

	user, err := r.rp.GetUserByID(ctx, id)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	if user == nil {
		err = errors.New("user not found")
		r.log.Error(err.Error())
		return err
	}

	url, err := r.st.UploadUserPicture(ctx, fmt.Sprintf("%s.jpg", id), image)
	if err != nil {
		r.log.Error(err.Error())
		return err
	}

	user.PictureUrl = url

	if err = r.rp.CreateOrUpdateUser(ctx, user); err != nil {
		r.log.Error(err.Error())
		return err
	}

	r.log.Debug("picture uploaded")
	return nil
}

func (r *userService) Login(ctx context.Context, login *model.Login) (*model.User, error) {
	err := r.validLogin(login)
	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}
	user, err := r.rp.GetUserByEmailAndPassword(ctx, login.Email, login.Password)
	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (r *userService) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := r.rp.GetUserByID(ctx, id)
	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (r *userService) GetUserRank(ctx context.Context, limit int) ([]model.User, error) {

	initDate := time.Now().Add((24 * 7 * time.Hour) * -1)
	scoreList, err := r.rp.GetAllUserScoreByTime(ctx, initDate)
	if err != nil {
		r.log.Error(err.Error())
		return nil, err
	}

	usersRank := make(map[string]int)
	for _, score := range scoreList {
		usersRank[score.UserID] += score.Points
	}

	keys := make([]string, 0, len(usersRank))

	for key := range usersRank {
		keys = append(keys, key)
	}
	sort.SliceStable(keys, func(i, j int) bool {
		return usersRank[keys[i]] > usersRank[keys[j]]
	})

	users := make([]model.User, 0, limit)
	for i, k := range keys {
		if i > limit {
			break
		}
		user, err := r.rp.GetUserByID(ctx, k)
		if err != nil {
			r.log.Error(err.Error())
			return nil, err
		}
		users = append(users, *user)
	}

	return users, nil
}

func (r *userService) validUser(user *model.User) error {
	if user == nil {
		return errors.New("user is nil")
	}
	if len(strings.TrimSpace(user.Email)) == 0 {
		return errors.New("email is empty")
	}
	if len(strings.TrimSpace(user.Name)) == 0 {
		return errors.New("name is empty")
	}
	if len(strings.TrimSpace(user.Password)) == 0 {
		return errors.New("password is empty")
	}
	if !isEmailValid(user.Email) {
		return errors.New("user email is invalid")
	}
	return nil
}

func (r *userService) validLogin(login *model.Login) error {
	if login == nil {
		return errors.New("login is nil")
	}
	if len(strings.TrimSpace(login.Email)) == 0 {
		return errors.New("email is empty")
	}
	if len(strings.TrimSpace(login.Password)) == 0 {
		return errors.New("password is empty")
	}
	if !isEmailValid(login.Email) {
		return errors.New("user email is invalid")
	}
	return nil
}

func isEmailValid(e string) bool {
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return emailRegex.MatchString(e)
}
