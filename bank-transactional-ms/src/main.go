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

func configureAccountService(repository domain.AccountRepository) domain.AccountService {
	return domain.NewAccountServiceImpl(repository)
}

func configureSyncRequestService(eventDispatcher domain.EventDispatcher) domain.SyncRequestService {
	return domain.NewSyncRequestServiceImpl(eventDispatcher)
}

func configureEventSubscriberService(service domain.AccountService, syncRequestService domain.SyncRequestService) domain.EventSubscriberService {
	return domain.NewEventSubscriberServiceImpl(service, syncRequestService)
}

func configureEventDispatcher(ctx context.Context, config *config.Config) domain.EventDispatcher {
	return adapters.NewEventDispatcherImpl(ctx, []string{config.KafkaBroker}, config.SyncRequestOutputTopic)
}

func configureConsumer(ctx context.Context, config *config.Config) *infra.Consumer {
	return infra.NewConsumer(ctx, []string{config.KafkaBroker}, config.SyncRequestInputTopic, config.KafkaClientId)
}

func configureServer(config *config.Config) *api.Server {
	ctx := context.Background()
	server := api.NewServer(config)
	conn := openDatabaseConnection(config)
	executeMigrations(conn)
	queries := db.New(conn)
	accountRepository := configureAccountRepository(ctx, queries)
	accountService := configureAccountService(accountRepository)
	eventDispatcher := configureEventDispatcher(ctx, config)
	syncRequestService := configureSyncRequestService(eventDispatcher)
	service := configureEventSubscriberService(accountService, syncRequestService)
	consumer := configureConsumer(ctx, config)
	consumer.StartReading(service.OnMessageReceive)
	controller := domain.NewTestController()
	server.ConfigureController(controller)
	return server
}

func startServer(server *api.Server) {
	err := server.Start()
	if err != nil {
		log.Fatal("Failed to start HTTP server")
	}
}
