package base

import (
	"context"
	"time"

	"server/global"
)

type JwtService struct{}

// GetRedisJWT 获取jwt
func (jwtService *JwtService) GetRedisJWT(username string) (redisJWT string, err error) {
	redisJWT, err = global.Redis.Get(context.Background(), username).Result()
	return redisJWT, err
}

// SetRedisJWT jwt存入redis并设置过期时间
func (jwtService *JwtService) SetRedisJWT(username string, jwt string) (err error) {
	// 此处过期时间等于jwt过期时间
	err = global.Redis.Set(context.Background(), username, jwt, time.Duration(global.Config.JWT.ExpiresTime)*time.Second).Err()
	return err
}
