package service

import "time"

func OptionTokenExpireTime(duration time.Duration) ServiceOptions {
	return func(s *Service) {
		s.tokenExpireTime = int64(duration) / int64(time.Second)
	}
}
