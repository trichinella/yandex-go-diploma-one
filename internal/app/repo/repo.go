package repo

import "diploma1/internal/app/entity"

type UserRepositoryInterface interface {
	AddUser(user entity.User) (*entity.User, error)
	UserByLogin(login string) (*entity.User, error)
}
