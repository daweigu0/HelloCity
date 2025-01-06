package dao

import (
	"HelloCity/internal/domain"
	"context"
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
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
		OpenID: openId,
	}
	err := h.db.WithContext(ctx).Where("open_id = ?", res.OpenID).First(&res).Error
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
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt time.Time `gorm:"type:datetime(3);default:null" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime(3);default:null" json:"updated_at"`
	DeletedAt time.Time `gorm:"type:datetime(3);default:null" json:"deleted_at"`
	Username  string    `gorm:"type:varchar(256);default:null" json:"username"`
	Password  string    `gorm:"type:varchar(256);default:null" json:"password"`
	Email     string    `gorm:"type:varchar(256);default:null" json:"email"`
	Mobile    string    `gorm:"type:varchar(256);default:null" json:"mobile"`
	OpenID    string    `gorm:"type:varchar(256);default:null" json:"open_id"`
	UnionID   string    `gorm:"type:varchar(256);default:null" json:"union_id"`
	NickName  string    `gorm:"type:varchar(256);default:null" json:"nick_name"`
	Gender    string    `gorm:"type:varchar(10);default:null" json:"gender"`
	Avatar    string    `gorm:"type:varchar(256);default:null" json:"avatar"`
	Address   string    `gorm:"type:varchar(256);default:null" json:"address"`
	Longitude float64   `gorm:"type:decimal(10,7);default:null" json:"longitude"`
	Latitude  float64   `gorm:"type:decimal(10,7);default:null" json:"latitude"`
	Status    int8      `gorm:"type:tinyint;default:null" json:"status"`
}
