package hsbc_home_task_zhangling

import (
	"errors"
	"github.com/1k2222/hsbc-home-task-zhangling/repo"
	"github.com/1k2222/hsbc-home-task-zhangling/service"
	"reflect"
	"testing"
	"time"
)

func initService(tokenExpireTime time.Duration) *service.Service {
	userRepo := repo.NewUser()
	roleRepo := repo.NewRole()
	userRoleRepo := repo.NewUserRoleRepo(userRepo, roleRepo)
	authTokenRepo := repo.NewAuthTokenRepo(userRepo)
	serv := service.NewService(userRepo, roleRepo, userRoleRepo, authTokenRepo, service.OptionTokenExpireTime(tokenExpireTime))
	return serv
}

func assertErrEqual(t *testing.T, actual, expected error) {
	if !errors.Is(expected, actual) {
		t.Errorf("unexpected result, test case: %s, expected err: %v, actual err: %v", t.Name(), actual, expected)
	}
}

func assertDeepEqual(t *testing.T, actual, expected interface{}) {
	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("unexpected result, test case: %s, expected: %v, actual: %v", t.Name(), actual, expected)
	}
}

func TestUserRoleCreateAndDelete(t *testing.T) {
	serv := initService(2 * time.Hour)

	// Create user
	err := serv.CreateUser("User1", "PassWord1")
	assertErrEqual(t, err, nil)
	err = serv.CreateUser("User2", "PassWord2")
	assertErrEqual(t, err, nil)

	// Create existed user
	err = serv.CreateUser("User1", "PassWord3")
	assertErrEqual(t, err, repo.ErrUserAlreadyExists)

	// Delete User
	err = serv.DeleteUser("User1")
	assertErrEqual(t, err, nil)

	// Delete user not existed
	err = serv.DeleteUser("User1")
	assertErrEqual(t, err, repo.ErrUserNotExist)

	// Create role
	err = serv.CreateRole("RoleA")
	assertErrEqual(t, err, nil)
	err = serv.CreateRole("RoleB")
	assertErrEqual(t, err, nil)

	// Create existed role
	err = serv.CreateRole("RoleA")
	assertErrEqual(t, err, repo.ErrRoleAlreadyExists)

	// Delete role
	err = serv.DeleteRole("RoleA")
	assertErrEqual(t, err, nil)

	// Delete role not existed
	err = serv.DeleteRole("RoleA")
	assertErrEqual(t, err, repo.ErrRoleNotExist)
}

func TestAuthenticateAndInvalidate(t *testing.T) {
	serv := initService(2 * time.Hour)

	// Create user
	err := serv.CreateUser("User1", "PassWord1")
	assertErrEqual(t, err, nil)
	err = serv.CreateUser("User2", "PassWord2")
	assertErrEqual(t, err, nil)

	// Right username and password.
	token, err := serv.Authenticate("User1", "PassWord1")
	assertErrEqual(t, err, nil)
	t.Logf("%s, login 'User1', get token: %s", t.Name(), token)

	// Right username, wrong password.
	_, err = serv.Authenticate("User1", "PassWord2")
	assertErrEqual(t, err, service.ErrInvalidUserNameOrPassword)

	// Wrong username.
	_, err = serv.Authenticate("User3", "PassWord2")
	assertErrEqual(t, err, service.ErrInvalidUserNameOrPassword)

	serv.Invalidate("000") // Should happen nothing because token '000' is invalid.
	// Token should be valid.
	_, err = serv.CheckRole(token, "Role1")
	assertErrEqual(t, err, nil)
	_, err = serv.AllRoles(token)
	assertErrEqual(t, err, nil)

	// After delete 'User1', its token should be invalid.
	err = serv.DeleteUser("User1")
	assertErrEqual(t, err, nil)
	_, err = serv.CheckRole(token, "Role1")
	assertErrEqual(t, err, service.ErrInvalidToken)
	_, err = serv.AllRoles(token)
	assertErrEqual(t, err, service.ErrInvalidToken)

	token, err = serv.Authenticate("User2", "PassWord2")
	assertErrEqual(t, err, nil)
	t.Logf("%s, login 'User2', get token: %s", t.Name(), token)

	// After invalidate token of 'User2', the token should be invalid.
	serv.Invalidate(token)
	_, err = serv.CheckRole(token, "Role1")
	assertErrEqual(t, err, service.ErrInvalidToken)
	_, err = serv.AllRoles(token)
	assertErrEqual(t, err, service.ErrInvalidToken)
}

func TestTokenExpireTime(t *testing.T) {
	serv := initService(3 * time.Second)

	// Create user
	err := serv.CreateUser("User1", "PassWord1")
	assertErrEqual(t, err, nil)

	// Right username and password.
	token, err := serv.Authenticate("User1", "PassWord1")
	assertErrEqual(t, err, nil)
	t.Logf("login, get token: %s", token)

	// Not expired.
	time.Sleep(time.Second * 2)
	_, err = serv.CheckRole(token, "Role1")
	assertErrEqual(t, err, nil)
	_, err = serv.AllRoles(token)
	assertErrEqual(t, err, nil)

	// Expired.
	time.Sleep(time.Second * 2)
	_, err = serv.CheckRole(token, "Role1")
	assertErrEqual(t, err, service.ErrTokenExpired)
	_, err = serv.AllRoles(token)
	assertErrEqual(t, err, service.ErrTokenExpired)
}

func TestUserRole(t *testing.T) {
	serv := initService(2 * time.Hour)

	err := serv.CreateUser("User1", "PassWord1")
	assertErrEqual(t, err, nil)
	token, err := serv.Authenticate("User1", "PassWord1")
	assertErrEqual(t, err, nil)
	t.Logf("%s, login 'User1', get token: %s", t.Name(), token)

	// User1 has no roles.
	ok, err := serv.CheckRole(token, "Role1")
	assertErrEqual(t, err, nil)
	assertDeepEqual(t, ok, false)
	allRoles, err := serv.AllRoles(token)
	assertErrEqual(t, err, nil)
	assertDeepEqual(t, allRoles, []string(nil))

	// User or role not exists.
	err = serv.AddRoleToUser("User1", "Role1")
	assertErrEqual(t, err, repo.ErrRoleNotExist)
	err = serv.CreateRole("Role1")
	assertErrEqual(t, err, nil)
	err = serv.CreateRole("Role2")
	assertErrEqual(t, err, nil)
	err = serv.AddRoleToUser("User2", "Role1")
	assertErrEqual(t, err, repo.ErrUserNotExist)

	// User1 was added to Role1 and Role2.
	roles := []string{"Role1", "Role2"}
	for i, role := range roles {
		err = serv.AddRoleToUser("User1", role)
		assertErrEqual(t, err, nil)
		ok, err = serv.CheckRole(token, role)
		assertErrEqual(t, err, nil)
		assertDeepEqual(t, ok, true)
		allRoles, err = serv.AllRoles(token)
		assertErrEqual(t, err, nil)
		assertDeepEqual(t, allRoles, roles[:i+1])
	}

	// Role1 was deleted.
	err = serv.DeleteRole("Role1")
	assertErrEqual(t, err, nil)
	ok, err = serv.CheckRole(token, "Role1")
	assertErrEqual(t, err, nil)
	assertDeepEqual(t, ok, false)
	ok, err = serv.CheckRole(token, "Role2")
	assertErrEqual(t, err, nil)
	assertDeepEqual(t, ok, true)
	allRoles, err = serv.AllRoles(token)
	assertErrEqual(t, err, nil)
	assertDeepEqual(t, allRoles, []string{"Role2"})

	// Role3 was added.
	err = serv.CreateRole("Role3")
	assertErrEqual(t, err, nil)
	err = serv.AddRoleToUser("User1", "Role3")
	ok, err = serv.CheckRole(token, "Role1")
	assertErrEqual(t, err, nil)
	assertDeepEqual(t, ok, false)
	ok, err = serv.CheckRole(token, "Role2")
	assertErrEqual(t, err, nil)
	assertDeepEqual(t, ok, true)
	ok, err = serv.CheckRole(token, "Role3")
	assertErrEqual(t, err, nil)
	assertDeepEqual(t, ok, true)
	allRoles, err = serv.AllRoles(token)
	assertErrEqual(t, err, nil)
	assertDeepEqual(t, allRoles, []string{"Role2", "Role3"})

	// User1 was deleted.
	err = serv.DeleteUser("User1")
	assertErrEqual(t, err, nil)
	_, err = serv.CheckRole(token, "Role1")
	assertErrEqual(t, err, service.ErrInvalidToken)
	_, err = serv.AllRoles(token)
	assertErrEqual(t, err, service.ErrInvalidToken)
}
