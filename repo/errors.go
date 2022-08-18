package repo

import "errors"

var ErrRoleAlreadyExists = errors.New("role already exists")
var ErrRoleNotExist = errors.New("role not exist")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrUserNotExist = errors.New("user not exist")
