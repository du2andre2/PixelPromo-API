package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"pixelPromo/domain/model"
	"regexp"
	"sort"
	"strings"
	"time"
)

func (s *service) CreateUser(ctx context.Context, user *model.User) error {

	err := s.validUser(user)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	user.CreatedAt = time.Now()
	user.ID = fmt.Sprintf("%d", user.CreatedAt.UnixNano())

	if err = s.rp.CreateOrUpdateUser(ctx, user); err != nil {
		s.log.Error(err.Error())
		return err
	}

	s.log.Debug("user created")
	return nil
}

func (s *service) UpdateUserPicture(ctx context.Context, id string, image io.Reader) error {

	user, err := s.rp.GetUserByID(ctx, id)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	if user == nil {
		err = errors.New("user not found")
		s.log.Error(err.Error())
		return err
	}

	url, err := s.st.UploadUserPicture(ctx, fmt.Sprintf("%s.jpg", id), image)
	if err != nil {
		s.log.Error(err.Error())
		return err
	}

	user.PictureUrl = url

	if err = s.rp.CreateOrUpdateUser(ctx, user); err != nil {
		s.log.Error(err.Error())
		return err
	}

	s.log.Debug("picture uploaded")
	return nil
}

func (s *service) Login(ctx context.Context, login *model.Login) (*model.User, error) {
	err := s.validLogin(login)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}
	user, err := s.rp.GetUserByEmailAndPassword(ctx, login.Email, login.Password)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (s *service) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user, err := s.rp.GetUserByID(ctx, id)
	if err != nil {
		s.log.Error(err.Error())
		return nil, err
	}

	return user, nil
}

func (s *service) GetUserRank(ctx context.Context, limit int) ([]model.User, error) {

	initDate := time.Now().Add((24 * 7 * time.Hour) * -1)
	scoreList, err := s.rp.GetAllUserScoreByTime(ctx, initDate)
	if err != nil {
		s.log.Error(err.Error())
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
		if i == limit {
			break
		}
		user, err := s.rp.GetUserByID(ctx, k)
		if err != nil {
			s.log.Error(err.Error())
			return nil, err
		}
		user.TotalScore = usersRank[k]
		users = append(users, *user)
	}

	return users, nil
}

func (s *service) validUser(user *model.User) error {
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

func (s *service) validLogin(login *model.Login) error {
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
