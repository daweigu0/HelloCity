package repository

import (
	"HelloCity/internal/domain"
	"HelloCity/internal/repository/dao"
	"context"
)

var (
	ErrDuplicateUser = dao.ErrDuplicateMobile
	ErrUserNotFound  = dao.ErrRecordNotFound
)

type UserRepository interface {
	FindByOpenId(ctx context.Context, openId string) (domain.User, error)
	FindById(ctx context.Context, uid uint64) (domain.User, error)
	Create(ctx context.Context, user domain.User) error
	Update(ctx context.Context, id uint64, user domain.User) error
}
type UserRepositoryHandler struct {
	dao dao.UserDao
}

func NewUserRepositoryHandler(dao dao.UserDao) *UserRepositoryHandler {
	return &UserRepositoryHandler{dao: dao}
}

func (*UserRepositoryHandler) toDomain(u dao.User) domain.User {
	return domain.User{
		ID:            u.ID,
		OpenID:        u.OpenID,
		Mobile:        u.Mobile,
		Avatar:        u.Avatar,
		NickName:      u.NickName,
		Email:         u.Email,
		ThumbsCount:   u.ThumbsCount,
		FansCount:     u.FansCount,
		FollowerCount: u.FollowersCount,
		Signature:     u.Signature,
		Constellation: u.Constellation,
		Province:      u.Province,
		City:          u.City,
	}
}

func (*UserRepositoryHandler) toEntity(u domain.User) dao.User {
	return dao.User{
		Mobile:         u.Mobile,
		OpenID:         u.OpenID,
		Avatar:         u.Avatar,
		NickName:       u.NickName,
		Gender:         u.Gender,
		Email:          u.Email,
		ThumbsCount:    u.ThumbsCount,
		FansCount:      u.FansCount,
		FollowersCount: u.FollowerCount,
		Signature:      u.Signature,
		Constellation:  u.Constellation,
		Province:       u.Province,
		City:           u.City,
	}
}

func (repo *UserRepositoryHandler) FindByOpenId(ctx context.Context, openId string) (domain.User, error) {
	u, err := repo.dao.FindUserByOpenId(ctx, openId)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepositoryHandler) FindById(ctx context.Context, uid uint64) (domain.User, error) {
	u, err := repo.dao.FindUserById(ctx, uid)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

func (repo *UserRepositoryHandler) Create(ctx context.Context, user domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(user))
}
func (repo *UserRepositoryHandler) Update(ctx context.Context, id uint64, user domain.User) error {
	return repo.dao.Update(ctx, id, repo.toEntity(user))
}
