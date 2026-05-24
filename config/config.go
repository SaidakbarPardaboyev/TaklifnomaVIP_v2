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

	// TODO: add host fields for each external API client your service depends on
	// Example: InternalExampleApiHost string
}
