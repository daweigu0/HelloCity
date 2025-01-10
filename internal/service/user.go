package service

import (
	"HelloCity/internal/domain"
	"HelloCity/internal/repository"
	"context"
	"errors"
)

var (
	ErrDuplicateMobile = repository.ErrDuplicateUser
	ErrInvalidUser     = errors.New("用户不存在")
)

type UserService interface {
	Login(ctx context.Context, openId string) (domain.User, error)
	FindUserByID(ctx context.Context, id uint64) (domain.User, error)
	SignUp(ctx context.Context, user domain.User) error
}
type UserServiceHandler struct {
	repo repository.UserRepository
}

func NewUserServiceHandler(repo repository.UserRepository) *UserServiceHandler {
	return &UserServiceHandler{repo}
}

func (svc *UserServiceHandler) Login(ctx context.Context, openId string) (domain.User, error) {
	u, err := svc.repo.FindByOpenId(ctx, openId)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUser
	}
	return u, err
}

func (svc *UserServiceHandler) FindUserByID(ctx context.Context, id uint64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}

func (svc *UserServiceHandler) SignUp(ctx context.Context, user domain.User) error {
	return svc.repo.Create(ctx, user)
}
