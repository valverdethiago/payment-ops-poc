package main

import (
	"context"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/adapters"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/api"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/config"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/domain"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/events"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/infra"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

var (
	DefaultAccountMapping = domain.Account{
		ID:                bson.NewObjectId(),
		InternalAccountId: "9b2272fe-53ad-4d5d-bfaf-ce2f1cf27ccf",
		ProviderAccountId: "6501663c-3a19-47df-9a2a-bc0796b702fd",
	}
)

func main() {
	configuration := loadConfig()
	server := configureServer(configuration)
	startServer(configuration, server)
}

func loadConfig() *config.Config {
	return util.LoadEnvConfig(".", "app")
}

func configureServer(config *config.Config) *api.Server {
	ctx := context.Background()
	server := api.NewServer(config)
	database := openDatabaseConnection(config)
	trioClient := configureTrioClient(config)
	accountRepository, transactionRepository, syncRequestRepository := configureRepositories(database)
	eventDispatcher := configureEventDispatcher(ctx, config)
	syncRequestService, balanceService, transactionService := configureServices(eventDispatcher,
		accountRepository, transactionRepository, syncRequestRepository, trioClient)
	eventSubscriberService := configureEventSubscriberService(syncRequestService, accountRepository,
		trioClient)
	configureConsumer(ctx, config, eventSubscriberService)
	configureControllers(server, syncRequestService, balanceService, transactionService, accountRepository)
	return server
}

func configureEventDispatcher(ctx context.Context, config *config.Config) domain.EventDispatcher {
	return adapters.NewEventDispatcherImpl(ctx, []string{config.KafkaBroker},
		config.SyncRequestOutputTopic, config.BalanceUpdateTopic, config.TransactionsUpdateTopic)
}

func configureRepositories(database *mgo.Database) (domain.AccountRepository,
	domain.TransactionRepository, domain.SyncRequestRepository) {
	syncRequestRepository := adapters.NewSyncRepositoryMongoDbImpl(database)
	accountRepository := adapters.NewAccountMongoDbRepositoryImpl(database)
	transactionRepository := adapters.NewTransactionMongoDbRepositoryImpl(database)
	configureDefaultAccountMapping(accountRepository)
	return accountRepository, transactionRepository, syncRequestRepository
}

func configureDefaultAccountMapping(accountRepository domain.AccountRepository) {
	AccountMapping, err := accountRepository.FindByInternalAccountId(DefaultAccountMapping.InternalAccountId)
	if err != nil || AccountMapping == nil {
		log.Println("Error searching default account mapping")
		accountRepository.Insert(&DefaultAccountMapping)
	}
}

func openDatabaseConnection(config *config.Config) *mgo.Database {
	return util.ConnectToDatabase(config)
}

func configureTrioClient(config *config.Config) domain.TrioClient {
	return adapters.NewTrioHttpClient(config)
}

func configureServices(dispatcher domain.EventDispatcher,
	accountRepository domain.AccountRepository,
	transactionRepository domain.TransactionRepository,
	syncRequestRepository domain.SyncRequestRepository,
	trioClient domain.TrioClient) (domain.SyncRequestService,
	domain.BalanceService, domain.TransactionService) {
	syncRequestService := configureSyncRequestService(dispatcher, syncRequestRepository)
	return syncRequestService,
		configureBalanceService(dispatcher, accountRepository, syncRequestService),
		configureTransactionService(dispatcher, syncRequestService, accountRepository,
			transactionRepository, trioClient)
}

func configureControllers(server *api.Server, syncRequestService domain.SyncRequestService,
	balanceService domain.BalanceService, transactionService domain.TransactionService,
	accountRepository domain.AccountRepository) {
	controller := domain.NewWebHookControllerImpl(syncRequestService, balanceService,
		transactionService, accountRepository)
	server.ConfigureController(controller)
}

func configureSyncRequestService(dispatcher domain.EventDispatcher,
	syncRequestRepository domain.SyncRequestRepository) domain.SyncRequestService {
	return domain.NewSyncRequestServiceImpl(dispatcher, syncRequestRepository)
}

func configureBalanceService(dispatcher domain.EventDispatcher,
	accountRepository domain.AccountRepository, syncRequestService domain.SyncRequestService) domain.BalanceService {
	return domain.NewBalanceServiceImpl(dispatcher, accountRepository, syncRequestService)
}

func configureTransactionService(dispatcher domain.EventDispatcher,
	syncRequestService domain.SyncRequestService,
	accountRepository domain.AccountRepository,
	transactionRepository domain.TransactionRepository,
	trioClient domain.TrioClient) domain.TransactionService {
	return domain.NewTransactionServiceImpl(dispatcher, syncRequestService,
		accountRepository, transactionRepository, trioClient)
}

func configureEventSubscriberService(syncRequestService domain.SyncRequestService,
	accountRepository domain.AccountRepository, trioClient domain.TrioClient) events.EventSubscriberService {
	return events.NewEventSubscriberServiceImpl(syncRequestService, accountRepository, trioClient)
}

func configureConsumer(ctx context.Context, config *config.Config,
	eventSubscriberService events.EventSubscriberService) *infra.Consumer {
	consumer := infra.NewConsumer(ctx, []string{config.KafkaBroker},
		config.SyncRequestInputTopic, config.KafkaClientId)
	go consumer.StartReading(eventSubscriberService.OnReceiveSyncRequest)
	return consumer
}

func startServer(config *config.Config, server *api.Server) {
	err := server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Failed to start HTTP server")
	}
}
