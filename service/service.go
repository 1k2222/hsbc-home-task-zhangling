package service

import (
	"github.com/1k2222/hsbc-home-task-zhangling/entity"
	"github.com/1k2222/hsbc-home-task-zhangling/repo"
	"sort"
	"strings"
	"time"
)

type Service struct {
	authTokenRepo *repo.AuthTokenRepo
	roleRepo      *repo.RoleRepo
	userRepo      *repo.UserRepo
	userRoleRepo  *repo.UserRoleRepo

	tokenExpireTime int64
}

type ServiceOptions func(s *Service)

func NewService(userRepo *repo.UserRepo, roleRepo *repo.RoleRepo, userRoleRepo *repo.UserRoleRepo, authTokenRepo *repo.AuthTokenRepo, opts ...ServiceOptions) *Service {
	ret := &Service{
		authTokenRepo: authTokenRepo,
		roleRepo:      roleRepo,
		userRepo:      userRepo,
		userRoleRepo:  userRoleRepo,
	}

	ret.ApplyOptions(opts...)
	return ret
}

func (s *Service) ApplyOptions(opts ...ServiceOptions) {
	for _, opt := range opts {
		opt(s)
	}
}

func (s *Service) CreateUser(name string, password string) error {
	user := entity.User{
		Name:     name,
		Password: hashPassword(name, password), // add salt to make password much harder to decipher by rainbow table.
	}
	return s.userRepo.Add(user)
}

func (s *Service) DeleteUser(name string) error {
	s.userRoleRepo.DeleteUser(name)
	s.authTokenRepo.DeleteUser(name)
	return s.userRepo.Delete(name)
}

func (s *Service) CreateRole(name string) error {
	role := entity.Role{Name: name}
	return s.roleRepo.Add(role)
}

func (s *Service) DeleteRole(name string) error {
	s.userRoleRepo.DeleteUserRole(name)
	return s.roleRepo.Delete(name)
}

func (s *Service) AddRoleToUser(userName, roleName string) error {
	return s.userRoleRepo.AddRoleToUser(userName, roleName)
}

func (s *Service) Authenticate(username string, password string) (string, error) {
	user, ok := s.userRepo.Get(username)
	if !ok {
		return "", ErrInvalidUserNameOrPassword
	}
	if hashPassword(username, password) != user.Password {
		return "", ErrInvalidUserNameOrPassword
	}
	authToken := generateAuthToken(username, password)
	s.authTokenRepo.Add(authToken)
	return authToken.Token, nil
}

func (s *Service) Invalidate(token string) {
	s.authTokenRepo.DeleteToken(token)
}

func (s *Service) CheckRole(token string, roleName string) (bool, error) {
	authToken, err := s.getAndCheckToken(token)
	if err != nil {
		return false, err
	}
	return s.userRoleRepo.UserBelongsToRole(authToken.UserName, roleName), nil
}

func (s *Service) AllRoles(token string) ([]string, error) {
	authToken, err := s.getAndCheckToken(token)
	if err != nil {
		return nil, err
	}
	allRoles := s.userRoleRepo.GetRolesOfUser(authToken.UserName)
	// Sort by lexicographical order to optimize user experience
	sort.Slice(allRoles, func(i, j int) bool {
		return strings.Compare(allRoles[i], allRoles[j]) < 0
	})
	return allRoles, nil
}

func (s *Service) getAndCheckToken(token string) (entity.AuthToken, error) {
	authToken, ok := s.authTokenRepo.GetToken(token)
	if !ok {
		return entity.AuthToken{}, ErrInvalidToken
	}
	if time.Now().Unix()-authToken.CreatedAt > s.tokenExpireTime {
		return entity.AuthToken{}, ErrTokenExpired
	}
	return authToken, nil
}
