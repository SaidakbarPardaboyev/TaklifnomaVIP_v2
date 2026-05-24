package config

type Config struct {
	ApiServerBindAddr string
	AllowedOrigin     string
	Timezone          string

	MongoConnectionString string
	MongoDatabaseName     string

	MySQLUsername string
	MySQLPassword string
	MySQLAddr     string
	MySQLDatabase string

	RedisAddr     string
	RedisUsername string
	RedisPassword string

	TelegramBotToken string

	PaymeMerchantID         string
	PaymeApiKeyProd         string
	PaymeApiKeyStaging      string
	PaymeRedirectionLink    string
	InvitationPrice         float64
}
