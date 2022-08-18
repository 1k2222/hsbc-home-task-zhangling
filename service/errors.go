package service

import "errors"

var ErrInvalidUserNameOrPassword = errors.New("invalid username or password")
var ErrInvalidToken = errors.New("invalid token")
var ErrTokenExpired = errors.New("token expired")
