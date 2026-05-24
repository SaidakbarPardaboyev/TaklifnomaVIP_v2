package caching

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"saidakbar.origin/caching/models"
)

const (
	tokenDatabaseIndex = 0
)

type tokenRedisClient struct {
	client *redis.Client
}

type TokenCache interface {
	GetAccount(token string) (*models.Account, error)
	GetClient(token string) (*models.ClientToken, error)
	SetAccount(token string, account *models.Account, ttl time.Duration) error
}

func (r *tokenRedisClient) GetAccount(token string) (*models.Account, error) {
	val, err := r.client.Get(context.TODO(), token).Result()

	if err != nil {
		return nil, err
	}

	var account models.Account

	if err = json.Unmarshal([]byte(val), &account); err != nil {
		return nil, err
	}

	return &account, nil

}

func (r *tokenRedisClient) GetClient(token string) (*models.ClientToken, error) {
	val, err := r.client.Get(context.TODO(), token).Result()

	if err != nil {
		return nil, err
	}

	var client models.ClientToken

	if err = json.Unmarshal([]byte(val), &client); err != nil {
		return nil, err
	}

	return &client, nil

}

func (r *tokenRedisClient) SetAccount(token string, account *models.Account, ttl time.Duration) error {
	data, err := json.Marshal(account)
	if err != nil {
		return err
	}
	return r.client.Set(context.TODO(), token, data, ttl).Err()
}

func NewClient(host string, username string, password string) TokenCache {
	client := redis.NewClient(&redis.Options{
		Addr:     host,
		Username: username,
		Password: password,
		DB:       tokenDatabaseIndex,
	})
	return &tokenRedisClient{client: client}
}
