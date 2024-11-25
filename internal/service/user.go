package service

import (
	"HelloCity/internal/domain"
	"HelloCity/internal/repository"
	"context"
)

type UserService interface {
	Login(ctx context.Context, openId string) (domain.User, error)
}
type UserServiceHandler struct {
	repo repository.UserRepository
}

func NewUserServiceHandler(repo repository.UserRepository) *UserServiceHandler {
	return &UserServiceHandler{repo}
}
func (h *UserServiceHandler) Login(ctx context.Context, openId string) (domain.User, error) {
	return h.repo.FindByOpenId(ctx, openId)
}
