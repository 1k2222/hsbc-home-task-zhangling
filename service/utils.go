package service

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/1k2222/hsbc-home-task-zhangling/entity"
	"math/rand"
	"strconv"
	"time"
)

func hash(input string) string {
	h := md5.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

func hashPassword(username, password string) string {
	return hash(password + hash(username))
}

func generateAuthToken(username string, password string) entity.AuthToken {
	timeNow := time.Now().Unix()
	// Make token relative to password and random number to avoid other people to guess the token.
	token := strconv.FormatInt(timeNow+rand.Int63(), 10) + username + password
	token = hash(token)
	return entity.AuthToken{
		Token:     token,
		UserName:  username,
		CreatedAt: timeNow,
	}
}
