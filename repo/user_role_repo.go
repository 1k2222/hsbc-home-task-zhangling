package repo

type UserRoleRepo struct {
	userRepo *UserRepo
	roleRepo *RoleRepo

	userRoles map[string]map[string]struct{}
}

func NewUserRoleRepo(userRepo *UserRepo, roleRepo *RoleRepo) *UserRoleRepo {
	return &UserRoleRepo{
		userRepo: userRepo,
		roleRepo: roleRepo,

		userRoles: make(map[string]map[string]struct{}),
	}
}

func (u *UserRoleRepo) AddRoleToUser(userName, roleName string) error {
	if _, ok := u.userRepo.Get(userName); !ok {
		return ErrUserNotExist
	}
	if _, ok := u.roleRepo.Get(roleName); !ok {
		return ErrRoleNotExist
	}
	if u.userRoles[userName] == nil {
		u.userRoles[userName] = make(map[string]struct{})
	}
	u.userRoles[userName][roleName] = struct{}{}
	return nil
}

func (u *UserRoleRepo) UserBelongsToRole(userName, roleName string) bool {
	_, ok := u.userRoles[userName][roleName]
	return ok
}

func (u *UserRoleRepo) GetRolesOfUser(userName string) []string {
	var ret []string
	for role := range u.userRoles[userName] {
		ret = append(ret, role)
	}
	return ret
}

func (u *UserRoleRepo) DeleteUser(name string) {
	delete(u.userRoles, name)
}

func (u *UserRoleRepo) DeleteUserRole(name string) {
	for _, roles := range u.userRoles {
		delete(roles, name)
	}
}
