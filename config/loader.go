package config

import "fmt"

const (
	configApiServerBindAddr = "API_SERVER_BIND_ADDR"
	configAllowedOrigin     = "ALLOWED_ORIGIN"
	configTimezone          = "TIMEZONE"

	configMongoConnectionString = "MONGO_CONNECTION_STRING"
	configMongoDatabaseName     = "MONGO_DATABASE_NAME"

	configMySQLUsername = "MYSQL_USERNAME"
	configMySQLPassword = "MYSQL_PASSWORD"
	configMySQLAddr     = "MYSQL_ADDR"
	configMySQLDatabase = "MYSQL_DATABASE"

	configRedisAddr     = "REDIS_ADDR"
	configRedisUsername = "REDIS_USERNAME"
	configRedisPassword = "REDIS_PASSWORD"

	configTelegramBotToken = "TELEGRAM_BOT_TOKEN"

	configPaymeMerchantID      = "PAYME_MERCHANT_ID"
	configPaymeApiKeyProd      = "PAYME_API_KEY_PROD"
	configPaymeApiKeyStaging   = "PAYME_API_KEY_STAGING"
	configPaymeRedirectionLink = "PAYME_REDIRECTION_LINK"
	configInvitationPrice      = "INVITATION_PRICE"
)

func Load(provider *EnvProvider) (cfg *Config, err error) {
	cfg = new(Config)

	if cfg.ApiServerBindAddr, err = provider.Get(configApiServerBindAddr); err != nil {
		return
	}
	if cfg.AllowedOrigin, err = provider.Get(configAllowedOrigin); err != nil {
		return
	}
	if cfg.Timezone, err = provider.Get(configTimezone); err != nil {
		return
	}

	if cfg.MongoConnectionString, err = provider.Get(configMongoConnectionString); err != nil {
		return
	}
	if cfg.MongoDatabaseName, err = provider.Get(configMongoDatabaseName); err != nil {
		return
	}

	if cfg.MySQLUsername, err = provider.Get(configMySQLUsername); err != nil {
		return
	}
	if cfg.MySQLPassword, err = provider.Get(configMySQLPassword); err != nil {
		return
	}
	if cfg.MySQLAddr, err = provider.Get(configMySQLAddr); err != nil {
		return
	}
	if cfg.MySQLDatabase, err = provider.Get(configMySQLDatabase); err != nil {
		return
	}

	if cfg.RedisAddr, err = provider.Get(configRedisAddr); err != nil {
		return
	}
	if cfg.RedisUsername, err = provider.Get(configRedisUsername); err != nil {
		return
	}
	if cfg.RedisPassword, err = provider.Get(configRedisPassword); err != nil {
		return
	}

	if cfg.TelegramBotToken, err = provider.Get(configTelegramBotToken); err != nil {
		return
	}

	if cfg.PaymeMerchantID, err = provider.Get(configPaymeMerchantID); err != nil {
		return
	}
	if cfg.PaymeApiKeyProd, err = provider.Get(configPaymeApiKeyProd); err != nil {
		return
	}
	if cfg.PaymeApiKeyStaging, err = provider.Get(configPaymeApiKeyStaging); err != nil {
		return
	}
	if cfg.PaymeRedirectionLink, err = provider.Get(configPaymeRedirectionLink); err != nil {
		return
	}

	var priceStr string
	if priceStr, err = provider.Get(configInvitationPrice); err != nil {
		return
	}
	if _, err = fmt.Sscanf(priceStr, "%f", &cfg.InvitationPrice); err != nil {
		err = fmt.Errorf("invalid %s: %w", configInvitationPrice, err)
		return
	}

	return cfg, nil
}
