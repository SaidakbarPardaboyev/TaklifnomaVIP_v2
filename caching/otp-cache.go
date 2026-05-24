package caching

import (
	"context"

	"github.com/go-redis/redis/v8"
	"saidakbar.origin/core"
)

const (
	otpKeyPrefix     = "otp:"
	otpDatabaseIndex = 0
)

type OtpCache interface {
	SetOTP(phone, code string) error
	GetOTP(phone string) (string, error)
	DeleteOTP(phone string) error
}

type otpRedisClient struct {
	client *redis.Client
}

func NewOtpCache(host, username, password string) OtpCache {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Username: username,
		Password: password,
		DB:       otpDatabaseIndex,
	})
	return &otpRedisClient{client: client}
}

func (r *otpRedisClient) SetOTP(phone, code string) error {
	return r.client.Set(context.TODO(), otpKeyPrefix+phone, code, core.OtpTTL).Err()
}

func (r *otpRedisClient) GetOTP(phone string) (string, error) {
	return r.client.Get(context.TODO(), otpKeyPrefix+phone).Result()
}

func (r *otpRedisClient) DeleteOTP(phone string) error {
	return r.client.Del(context.TODO(), otpKeyPrefix+phone).Err()
}
