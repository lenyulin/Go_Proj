package dao

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrduplicateEmail = errors.New("Email exits")
	EreRecordNotFound = gorm.ErrRecordNotFound
)

type UserDAO struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

type User struct {
	Id           int64  `gorm:"primaryKey,autoIncrement"`
	Email        string `gorm:"unique"`
	Password     string
	Ctime        int64
	Utime        int64
	Birthday     int64
	Introduction string
	NickName     string
}

func (dao *UserDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			return ErrduplicateEmail
		}
	}
	return err
}

func (dao *UserDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func (dao *UserDAO) UpdateInfo(ctx *gin.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Utime = now
	err := dao.db.WithContext(ctx).Where("email=?", u.Email).Updates(&u).Error
	if err != nil {
		return err
	}
	return nil
}
