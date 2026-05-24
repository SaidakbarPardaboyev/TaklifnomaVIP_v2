package config

import (
	"fmt"
	"os"
)

type EnvProvider struct{}

func NewEnvProvider() *EnvProvider {
	return &EnvProvider{}
}

func (e *EnvProvider) Get(key string) (string, error) {
	val, ok := os.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("env var %s not set", key)
	}
	return val, nil
}
