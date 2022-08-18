package repo

import (
	"github.com/1k2222/hsbc-home-task-zhangling/entity"
)

type AuthTokenRepo struct {
	userRepo *UserRepo

	userNameIndex map[string]entity.AuthToken
	tokenIndex    map[string]entity.AuthToken
}

func NewAuthTokenRepo(userRepo *UserRepo) *AuthTokenRepo {
	return &AuthTokenRepo{
		userRepo: userRepo,

		userNameIndex: make(map[string]entity.AuthToken),
		tokenIndex:    make(map[string]entity.AuthToken),
	}
}

func (a *AuthTokenRepo) Add(authToken entity.AuthToken) {
	if _, ok := a.userRepo.Get(authToken.UserName); !ok {
		return
	}
	a.tokenIndex[authToken.Token] = authToken
	a.userNameIndex[authToken.UserName] = authToken
}

// DeleteToken Do nothing if token not exist to satisfy idempotence.
func (a *AuthTokenRepo) DeleteUser(name string) {
	authToken, ok := a.userNameIndex[name]
	if !ok {
		return
	}
	delete(a.userNameIndex, authToken.UserName)
	delete(a.tokenIndex, authToken.Token)
	return
}

func (a *AuthTokenRepo) DeleteToken(token string) {
	authToken, ok := a.tokenIndex[token]
	if !ok {
		return
	}
	delete(a.userNameIndex, authToken.UserName)
	delete(a.tokenIndex, authToken.Token)
}

func (a *AuthTokenRepo) GetToken(token string) (entity.AuthToken, bool) {
	authToken, ok := a.tokenIndex[token]
	return authToken, ok
}
