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
	FindUserByOpenID(ctx context.Context, openId string) (domain.User, error)
	SignUp(ctx context.Context, user domain.User) error
	UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error
	Profile(ctx context.Context, id uint64) (domain.User, error)
	Edit(ctx context.Context, id uint64, user domain.User) error
}
type userService struct {
	repo repository.UserRepository
}

func (svc *userService) FindUserByOpenID(ctx context.Context, openId string) (domain.User, error) {
	return svc.repo.FindByOpenId(ctx, openId)
}

func (svc *userService) UpdateNonSensitiveInfo(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateNonZeroFields(ctx, user)
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) Login(ctx context.Context, openId string) (domain.User, error) {
	u, err := svc.repo.FindByOpenId(ctx, openId)
	if err == repository.ErrUserNotFound {
		return domain.User{}, ErrInvalidUser
	}
	return u, err
}

func (svc *userService) FindUserByID(ctx context.Context, id uint64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}

func (svc *userService) SignUp(ctx context.Context, user domain.User) error {
	return svc.repo.Create(ctx, user)
}
func (svc *userService) Profile(ctx context.Context, id uint64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}
func (svc *userService) Edit(ctx context.Context, id uint64, user domain.User) error {
	return svc.repo.Update(ctx, id, user)
}
