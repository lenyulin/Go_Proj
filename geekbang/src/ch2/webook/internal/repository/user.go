package repository

import (
	"ch2/webook/internal/domain"
	"ch2/webook/internal/repository/dao"
	"context"
	"github.com/gin-gonic/gin"
	"time"
)

type UserRepository struct {
	dao *dao.UserDAO
}

var (
	ErrduplicateEmail = dao.ErrduplicateEmail
	EreUserNotFound   = dao.EreRecordNotFound
)

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}
func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}
}
func (repo *UserRepository) infoToDomain(u dao.User) domain.User {
	nanoTimestamp := u.Birthday * 1e6
	t := time.Unix(0, nanoTimestamp)
	layout := "2006-01-02"
	dateStr := t.Format(layout)
	return domain.User{
		Id:           u.Id,
		Email:        u.Email,
		NickName:     u.NickName,
		Birthday:     dateStr,
		Introduction: u.Introduction,
	}
}
func (repo *UserRepository) FindInfoByEmail(ctx *gin.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.infoToDomain(u), nil
}

func (repo *UserRepository) Edit(ctx *gin.Context, u domain.User) error {
	//1997-02-03
	t, err := time.Parse("2006-01-01", u.Birthday)
	if err != nil {
		return err
	}
	nanoTimestamp := t.UnixNano() / 1e6
	return repo.dao.UpdateInfo(ctx, dao.User{
		Email:        u.Email,
		Password:     u.Password,
		Birthday:     nanoTimestamp,
		Introduction: u.Introduction,
		NickName:     u.NickName,
	})
}
