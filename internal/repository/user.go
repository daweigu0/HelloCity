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
	UpdateNonZeroFields(ctx context.Context, user domain.User) error
}
type UserRepositoryHandler struct {
	dao dao.UserDao
}

func (repo *UserRepositoryHandler) UpdateNonZeroFields(ctx context.Context, user domain.User) error {
	// 更新 DB 之后，删除
	err := repo.dao.UpdateById(ctx, repo.toEntity(user))
	if err != nil {
		return err
	}
	return nil
}

func NewUserRepositoryHandler(dao dao.UserDao) *UserRepositoryHandler {
	return &UserRepositoryHandler{dao: dao}
}

func (*UserRepositoryHandler) toDomain(u dao.User) domain.User {
	return domain.User{
		ID:       u.ID,
		OpenID:   u.OpenID,
		Mobile:   u.Mobile,
		Avatar:   u.Avatar,
		NickName: u.NickName,
		Email:    u.Email,
	}
}

func (*UserRepositoryHandler) toEntity(u domain.User) dao.User {
	return dao.User{
		Mobile:   u.Mobile,
		OpenID:   u.OpenID,
		Avatar:   u.Avatar,
		NickName: u.NickName,
		Gender:   u.Gender,
		Email:    u.Email,
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
