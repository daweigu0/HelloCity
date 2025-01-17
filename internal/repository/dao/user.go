package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrDuplicateMobile = errors.New("手机号冲突")
	ErrRecordNotFound  = gorm.ErrRecordNotFound
)

type UserDao interface {
	Insert(ctx context.Context, user User) error
	FindUserByOpenId(ctx context.Context, openId string) (User, error)
	FindUserById(ctx context.Context, id uint64) (User, error)
	Update(ctx context.Context, id uint64, user User) error
	UpdateById(ctx context.Context, entity User) error
}
type GORMUserDao struct {
	db *gorm.DB
}

func (dao *GORMUserDao) UpdateById(ctx context.Context, entity User) error {
	// 这种写法依赖于 GORM 的零值和主键更新特性
	// Update 非零值 WHERE id = ?
	//return dao.db.WithContext(ctx).Updates(&entity).Error
	return dao.db.WithContext(ctx).Model(&entity).Where("id = ?", entity.ID).
		Updates(map[string]any{
			"updated_at": time.Now(),
			"nick_name":  entity.NickName,
			"avatar":     entity.Avatar,
			"gender":     entity.Gender,
		}).Error
}

func NewUserDAO(db *gorm.DB) *GORMUserDao {
	return &GORMUserDao{
		db: db,
	}
}

func (dao *GORMUserDao) Insert(ctx context.Context, user User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	err := dao.db.WithContext(ctx).Create(&user).Error
	if me, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			// 冲入冲突，手机号冲突
			return ErrDuplicateMobile
		}
	}
	return err
}

func (dao *GORMUserDao) FindUserByOpenId(ctx context.Context, openId string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("id = ?", openId).First(&user).Error
	return user, err

}

func (dao *GORMUserDao) FindUserById(ctx context.Context, id uint64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("id = ?", id).First(&user).Omit("password").Error
	return user, err
}
func (dao *GORMUserDao) Update(ctx context.Context, id uint64, user User) error {
	user.UpdatedAt = time.Now()
	err := dao.db.WithContext(ctx).Model(&user).Where("id=?", id).Updates(user).Error
	return err
}

type User struct {
	ID             uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt      time.Time `gorm:"type:datetime(3);default:null" json:"created_at"`
	UpdatedAt      time.Time `gorm:"type:datetime(3);default:null" json:"updated_at"`
	DeletedAt      time.Time `gorm:"type:datetime(3);default:null" json:"deleted_at"`
	Username       string    `gorm:"type:varchar(256);default:null" json:"username"`
	Password       string    `gorm:"type:varchar(256);default:null" json:"password"`
	Email          string    `gorm:"type:varchar(256);default:null" json:"email"`
	Mobile         string    `gorm:"type:varchar(256);default:null;unique" json:"mobile"`
	OpenID         string    `gorm:"type:varchar(256);default:null;unique" json:"open_id"`
	UnionID        string    `gorm:"type:varchar(256);default:null;unique" json:"union_id"`
	NickName       string    `gorm:"type:varchar(256);default:null" json:"nick_name"`
	Gender         string    `gorm:"type:varchar(10);default:null" json:"gender"`
	Avatar         string    `gorm:"type:varchar(256);default:null" json:"avatar"`
	Address        string    `gorm:"type:varchar(256);default:null" json:"address"`
	Longitude      float64   `gorm:"type:decimal(10,7);default:null" json:"longitude"`
	Latitude       float64   `gorm:"type:decimal(10,7);default:null" json:"latitude"`
	Status         int8      `gorm:"type:tinyint;default:null" json:"status"`
	ThumbsCount    int64     `gorm:"type:bigint;default:0" json:"thumbs_count"`
	FansCount      int64     `gorm:"type:bigint;default:0" json:"fans_count"`
	FollowersCount int64     `gorm:"type:bigint;default:0" json:"followers_count"`
	Signature      string    `gorm:"type:varchar(256);default:null" json:"signature"`
	AboutMe        string    `gorm:"type:varchar(256);default:null" json:"about_me"`
	Constellation  int8      `gorm:"type:tinyint;default:null" json:"constellation"`
	Province       string    `gorm:"type:varchar(256);default:null" json:"province"`
	City           string    `gorm:"type:varchar(256);default:null" json:"city"`
}
