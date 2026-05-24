package caching

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"saidakbar.origin/caching/models"
)

const clientCredentialsDatabaseIndex = 12

type ClientCredentialsCache interface {
	GetClientCredentials(clientID string, clientSecret string) (clientCredentials *models.ClientCredentials, err error)
}

type clientCredentialsCache struct {
	client *redis.Client
}

func (c clientCredentialsCache) GetClientCredentials(clientID string, clientSecret string) (clientCredential *models.ClientCredentials, err error) {
	var cachedValue string

	if cachedValue, err = c.client.Get(context.TODO(), getClientCredentialsKey(clientID, clientSecret)).Result(); err != nil {
		return
	}

	if err = json.Unmarshal([]byte(cachedValue), &clientCredential); err != nil {
		return
	}

	return
}

func getClientCredentialsKey(clientID string, clientSecret string) string {
	return fmt.Sprintf("cc_%s_%s", clientID, clientSecret)
}

func NewClientCredentialCache(address, username, password string) ClientCredentialsCache {
	return &clientCredentialsCache{
		client: redis.NewClient(&redis.Options{
			Addr:     address,
			Username: username,
			Password: password,
			DB:       clientCredentialsDatabaseIndex,
		}),
	}
}
