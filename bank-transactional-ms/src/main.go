package main

import (
	"context"
	"database/sql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/adapters"
	"github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/api"
	"github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/config"
	db "github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/db/sqlc"
	"github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/domain"
	"github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/infra"
	"github.com/Pauca-Technologies/payment-ops-poc/bank-transactional-ms/util"
)

var database *sql.DB

func main() {
	config := loadConfig()
	server := configureServer(config)
	startServer(server)
}

func loadConfig() *config.Config {
	return util.LoadEnvConfig(".", "app")
}

func openDatabaseConnection(config *config.Config) *sql.DB {
	return util.ConnectToDatabase(config)
}

func executeMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://./db/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatal("Failed to apply database migrations", err)
	}
	m.Up()
}

func configureAccountRepository(ctx context.Context, querier db.Querier) domain.AccountRepository {
	return adapters.NewAccountRepositoryImpl(querier, ctx)
}

func configureAccountBalanceRepository(ctx context.Context, querier db.Querier) domain.AccountBalanceRepository {
	return adapters.NewAccountBalanceRepositoryImpl(querier, ctx)
}

func configureAccountService(repository domain.AccountRepository) domain.AccountService {
	return domain.NewAccountServiceImpl(repository)
}

func configureAccountBalanceService(repository domain.AccountBalanceRepository) domain.AccountBalanceService {
	return domain.NewAccountBalanceServiceImpl(repository)
}

func configureSyncRequestService(eventDispatcher domain.EventDispatcher) domain.SyncRequestService {
	return domain.NewSyncRequestServiceImpl(eventDispatcher)
}

func configureEventSubscriberService(accountService domain.AccountService,
	accountBalanceService domain.AccountBalanceService,
	syncRequestService domain.SyncRequestService) domain.EventSubscriberService {
	return domain.NewEventSubscriberServiceImpl(accountService, accountBalanceService, syncRequestService)
}

func configureEventDispatcher(ctx context.Context, config *config.Config) domain.EventDispatcher {
	return adapters.NewEventDispatcherImpl(ctx, []string{config.KafkaBroker}, config.SyncRequestOutputTopic)
}

func configureSyncRequestConsumer(ctx context.Context, config *config.Config,
	subscriberService domain.EventSubscriberService) *infra.Consumer {
	return configureConsumer(ctx, config.KafkaBroker,
		config.SyncRequestInputTopic, config.KafkaClientId, subscriberService.OnReceiveSyncRequest)
}

func configureTransactionsUpdateConsumer(ctx context.Context, config *config.Config,
	subscriberService domain.EventSubscriberService) *infra.Consumer {
	return configureConsumer(ctx, config.KafkaBroker, config.TransactionsUpdateTopic,
		config.KafkaClientId, subscriberService.OnReceiveTransactionsUpdate)
}

func configureBalanceUpdateConsumer(ctx context.Context, config *config.Config,
	subscriberService domain.EventSubscriberService) *infra.Consumer {
	return configureConsumer(ctx, config.KafkaBroker, config.BalanceUpdateTopic,
		config.KafkaClientId, subscriberService.OnReceiveBalanceUpdate)
}

func configureConsumer(ctx context.Context, kafkaBroker string, kafkaTopic string,
	kafkaClientId string, callback domain.OnMessageReceive) *infra.Consumer {
	consumer := infra.NewConsumer(ctx, []string{kafkaBroker}, kafkaTopic, kafkaClientId)
	go consumer.StartReading(callback)
	return consumer
}

func configureServer(config *config.Config) *api.Server {
	ctx := context.Background()
	server := api.NewServer(config)
	conn := openDatabaseConnection(config)
	executeMigrations(conn)
	queries := db.New(conn)
	accountRepository := configureAccountRepository(ctx, queries)
	accountBalanceRepository := configureAccountBalanceRepository(ctx, queries)
	accountService := configureAccountService(accountRepository)
	accountBalanceService := configureAccountBalanceService(accountBalanceRepository)
	eventDispatcher := configureEventDispatcher(ctx, config)
	syncRequestService := configureSyncRequestService(eventDispatcher)
	service := configureEventSubscriberService(accountService, accountBalanceService, syncRequestService)
	configureConsumers(ctx, config, service)
	controller := domain.NewTestController()
	server.ConfigureController(controller)
	return server
}

func configureConsumers(ctx context.Context, config *config.Config, service domain.EventSubscriberService) {
	configureSyncRequestConsumer(ctx, config, service)
	configureBalanceUpdateConsumer(ctx, config, service)
	//configureTransactionsUpdateConsumer(ctx, config, service)
}

func startServer(server *api.Server) {
	err := server.Start()
	if err != nil {
		log.Fatal("Failed to start HTTP server")
	}
}
