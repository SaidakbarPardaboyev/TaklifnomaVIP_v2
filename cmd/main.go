package main

import (
	"log"
	"strings"

	apiserver "saidakbar.origin/api-server"
	"saidakbar.origin/caching"
	"saidakbar.origin/config"
	mongo_db "saidakbar.origin/db/mongo"
	mysql_db "saidakbar.origin/db/mysql"
	"saidakbar.origin/plugins"
	"saidakbar.origin/repository"
	account_service "saidakbar.origin/services/account"
	telegram_service "saidakbar.origin/services/telegram"
	templete_mongo "saidakbar.origin/services/templete-mongo"
	templete_mysql "saidakbar.origin/services/templete-mysql"

	"github.com/joho/godotenv"
)

var configs *config.Config

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, reading from environment")
	}

	var failure error
	if configs, failure = config.Load(config.NewEnvProvider()); failure != nil {
		log.Fatal(failure)
	}

	plugins.SetTimeZoneName(configs.Timezone)
}

func main() {
	mongoDB := mongo_db.NewDatabase(configs.MongoConnectionString, configs.MongoDatabaseName)

	mysqlDB, err := mysql_db.New(configs.MySQLUsername, configs.MySQLPassword, configs.MySQLAddr, configs.MySQLDatabase)
	if err != nil {
		log.Fatal(err)
	}
	if err = mysqlDB.Migrate(); err != nil {
		log.Fatal(err)
	}

	tokenCache := caching.NewClient(configs.RedisAddr, configs.RedisUsername, configs.RedisPassword)
	clientCredentialsCache := caching.NewClientCredentialCache(configs.RedisAddr, configs.RedisUsername, configs.RedisPassword)
	otpCache := caching.NewOtpCache(configs.RedisAddr, configs.RedisUsername, configs.RedisPassword)

	var templeteMongoRepository = repository.NewTempleteMongoRepository(mongoDB)
	var templatMysqlRepository = repository.NewTemplateMysqlRepository(mysqlDB)
	var accountRepository = repository.NewAccountRepository(mysqlDB)

	var templeteMongoService = templete_mongo.NewTempleteMongoService(templeteMongoRepository)
	var templatMysqlService = templete_mysql.NewTemplateMysqlService(templatMysqlRepository)
	var accountSvc = account_service.NewService(accountRepository, otpCache, tokenCache)
	// TODO: add more services here

	telegramSvc, err := telegram_service.NewService(configs.TelegramBotToken, accountRepository, otpCache)
	if err != nil {
		log.Fatal(err)
	}
	go telegramSvc.StartPolling()

	server := apiserver.NewApiServer(
		configs.ApiServerBindAddr,
		strings.Split(configs.AllowedOrigin, ";"),
		tokenCache,
		clientCredentialsCache,
		accountSvc,
		templeteMongoService,
		templatMysqlService,
	)

	log.Fatalln(server.Start())
}
