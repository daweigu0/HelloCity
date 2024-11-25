package repository

import (
	"HelloCity/internal/domain"
	"HelloCity/internal/repository/dao"
	"context"
)

type UserRepository interface {
	FindByOpenId(ctx context.Context, openId string) (domain.User, error)
}
type UserRepositoryHandler struct {
	dao dao.UserDao
}

func NewUserRepositoryHandler(dao dao.UserDao) *UserRepositoryHandler {
	return &UserRepositoryHandler{dao: dao}
}
func (h *UserRepositoryHandler) FindByOpenId(ctx context.Context, openId string) (domain.User, error) {
	return h.dao.FindOrCreateByOpenId(ctx, openId)
}
