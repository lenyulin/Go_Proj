package service

import (
	"ch2/webook/internal/domain"
	"ch2/webook/internal/repository"
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

var (
	ErrduplicateEmail        = repository.ErrduplicateEmail
	ErrInvalidUserOrPassword = errors.New("User Or Password Error")
)

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (svc *UserService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repo.Create(ctx, u)
}

func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.EreUserNotFound {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	if err != nil {
		return domain.User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}
	return u, err
}

func (svc *UserService) Profile(ctx *gin.Context, email string) (domain.User, error) {
	u, err := svc.repo.FindInfoByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}

func (svc *UserService) Edit(ctx *gin.Context, user domain.User) error {
	_, err := svc.repo.FindByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	return svc.repo.Edit(ctx, user)
}
