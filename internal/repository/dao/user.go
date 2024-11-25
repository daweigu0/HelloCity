package dao

import (
	"HelloCity/internal/domain"
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
)

type UserDao interface {
	FindOrCreateByOpenId(ctx context.Context, openId string) (domain.User, error)
}
type UserDaoHandler struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDaoHandler {
	return &UserDaoHandler{
		db: db,
	}
}
func (h *UserDaoHandler) FindOrCreateByOpenId(ctx context.Context, openId string) (domain.User, error) {
	res := User{
		OpenId: openId,
	}
	err := h.db.WithContext(ctx).Where("open_id = ?", res.OpenId).First(&res).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		log.Println("找不到记录，创建")
		err = h.db.WithContext(ctx).Create(&res).Error
		if err != nil {
			log.Println("创建失败", err)
			return domain.User{}, err
		}
		return domain.User{
			OpenId: openId,
			Uid:    res.ID,
		}, nil
	} else if err != nil {
		log.Println(err)
		return domain.User{}, err
	}
	return domain.User{
		OpenId: openId,
		Uid:    res.ID,
	}, nil
}

type User struct {
	gorm.Model
	OpenId string `gorm:"unique"`
}
