package config

import (
	"fmt"

	vault "github.com/hashicorp/vault/api"
)

type VaultProvider struct {
	client *vault.Client
	path   string // secret/data/myapp
}

func NewVaultProvider(addr, username, password, path string) (*VaultProvider, error) {
	cfg := vault.DefaultConfig()
	cfg.Address = addr

	client, err := vault.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	secret, failure := client.Logical().Write(fmt.Sprintf("auth/userpass/login/%s", username), map[string]interface{}{
		"password": password,
	})
	if failure != nil {
		return nil, failure
	}
	token := secret.Auth.ClientToken
	client.SetToken(token)

	return &VaultProvider{
		client: client,
		path:   path,
	}, nil
}

func (v *VaultProvider) Get(key string) (string, error) {
	secret, err := v.client.Logical().Read(v.path)
	if err != nil || secret == nil {
		return "", fmt.Errorf("vault read error: %w", err)
	}

	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid vault data format")
	}

	val, ok := data[key]
	if !ok {
		return "", fmt.Errorf("secret %s not found", key)
	}

	return fmt.Sprint(val), nil
}
