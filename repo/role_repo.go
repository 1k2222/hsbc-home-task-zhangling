package repo

import (
	"github.com/1k2222/hsbc-home-task-zhangling/entity"
)

type RoleRepo struct {
	data map[string]entity.Role
}

func NewRole() *RoleRepo {
	return &RoleRepo{data: make(map[string]entity.Role)}
}

func (r *RoleRepo) Add(role entity.Role) error {
	if _, ok := r.data[role.Name]; ok {
		return ErrRoleAlreadyExists
	}
	r.data[role.Name] = role
	return nil
}

func (r *RoleRepo) Delete(name string) error {
	if _, ok := r.data[name]; !ok {
		return ErrRoleNotExist
	}
	delete(r.data, name)
	return nil
}

func (r *RoleRepo) Get(name string) (entity.Role, bool) {
	role, ok := r.data[name]
	return role, ok
}
