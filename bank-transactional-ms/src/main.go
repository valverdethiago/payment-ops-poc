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

func configureServer(config *config.Config) *api.Server {
	ctx := context.Background()
	server := api.NewServer(config)
	conn := openDatabaseConnection(config)
	executeMigrations(conn)
	queries := db.New(conn)
	accountRepository, accountBalanceRepository, transactionRepository := configureRepositories(ctx, queries)
	accountService, accountBalanceService, transactionService := configureServices(accountRepository,
		accountBalanceRepository, transactionRepository)
	configureEventServices(ctx, config, accountService, accountBalanceService, transactionService)
	configureControllers(*server, accountService, transactionService)
	return server
}

func configureControllers(server api.Server, accountService domain.AccountService,
	transactionService domain.TransactionService) {
	server.ConfigureController(domain.NewTestController())
	server.ConfigureController(domain.NewAccountControllerImpl(accountService, transactionService))
}

func configureRepositories(ctx context.Context, querier db.Querier) (domain.AccountRepository, domain.AccountBalanceRepository,
	domain.TransactionRepository) {
	return adapters.NewAccountRepositoryImpl(querier, ctx),
		adapters.NewAccountBalanceRepositoryImpl(querier, ctx),
		adapters.NewTransactionRepositoryImpl(querier, ctx)
}

func configureServices(accountRepository domain.AccountRepository,
	accountBalanceRepository domain.AccountBalanceRepository,
	transactionRepository domain.TransactionRepository) (domain.AccountService,
	domain.BalanceService, domain.TransactionService) {
	accountBalanceService := domain.NewAccountBalanceServiceImpl(accountBalanceRepository)
	transactionService := domain.NewTransactionServiceImpl(transactionRepository)
	accountService := domain.NewAccountServiceImpl(accountRepository, accountBalanceService)
	return accountService,
		accountBalanceService,
		transactionService
}

func configureEventServices(ctx context.Context, config *config.Config,
	accountService domain.AccountService, accountBalanceService domain.BalanceService,
	transactionService domain.TransactionService) {
	eventDispatcher := configureEventDispatcher(ctx, config)
	syncRequestService := configureSyncRequestService(eventDispatcher)
	eventSubscriberServices := configureEventSubscriberService(accountService, accountBalanceService, transactionService,
		syncRequestService)
	configureEventConsumers(ctx, config, eventSubscriberServices)
}

func configureSyncRequestService(eventDispatcher domain.EventDispatcher) domain.SyncRequestService {
	return domain.NewSyncRequestServiceImpl(eventDispatcher)
}

func configureEventSubscriberService(accountService domain.AccountService,
	accountBalanceService domain.BalanceService,
	transactionService domain.TransactionService,
	syncRequestService domain.SyncRequestService) domain.EventSubscriberService {
	return domain.NewEventSubscriberServiceImpl(accountService, accountBalanceService, transactionService,
		syncRequestService)
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

func configureEventConsumers(ctx context.Context, config *config.Config, service domain.EventSubscriberService) {
	configureSyncRequestConsumer(ctx, config, service)
	configureBalanceUpdateConsumer(ctx, config, service)
	configureTransactionsUpdateConsumer(ctx, config, service)
}

func startServer(server *api.Server) {
	err := server.Start()
	if err != nil {
		log.Fatal("Failed to start HTTP server")
	}
}
