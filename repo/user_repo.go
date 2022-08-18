package repo

import (
	"github.com/1k2222/hsbc-home-task-zhangling/entity"
)

type UserRepo struct {
	data map[string]entity.User
}

func NewUser() *UserRepo {
	return &UserRepo{data: make(map[string]entity.User)}
}

func (u *UserRepo) Add(newUser entity.User) error {
	if _, ok := u.data[newUser.Name]; ok {
		return ErrUserAlreadyExists
	}
	u.data[newUser.Name] = newUser
	return nil
}

func (u *UserRepo) Delete(name string) error {
	if _, ok := u.data[name]; !ok {
		return ErrUserNotExist
	}
	delete(u.data, name)
	return nil
}

func (u *UserRepo) Get(name string) (entity.User, bool) {
	user, ok := u.data[name]
	return user, ok
}
