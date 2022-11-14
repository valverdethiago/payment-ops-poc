package main

import (
	"context"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/adapters"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/api"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/config"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/domain"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/infra"
	"github.com/Pauca-Technologies/payment-ops-poc/trio-provider-ms/util"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

var (
	DefaultAccountMapping domain.AccountMapping = domain.AccountMapping{
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
	accountRepository, syncRequestRepository := configureRepositories(database)
	eventDispatcher := configureEventDispatcher(ctx, config)
	eventSubscriberService := configureEventSubscriberService(syncRequestRepository, accountRepository,
		trioClient, eventDispatcher)
	configureConsumer(ctx, config, eventSubscriberService)
	//producer := configureKafkaProducer(ctx, config)
	configureControllers(server)
	return server
}

func configureEventDispatcher(ctx context.Context, config *config.Config) domain.EventDispatcher {
	return adapters.NewEventDispatcherImpl(ctx, []string{config.KafkaBroker}, config.SyncRequestOutputTopic)
}

func configureRepositories(database *mgo.Database) (domain.AccountRepository, domain.SyncRequestRepository) {
	syncRequestRepository := adapters.NewSyncRepositoryMongoDbImpl(database)
	accountRepository := adapters.NewAccountMappingMongoDbRepositoryImpl(database)
	configureDefaultAccountMapping(accountRepository)
	return accountRepository, syncRequestRepository
}

func configureDefaultAccountMapping(accountRepository domain.AccountRepository) {
	AccountMapping, err := accountRepository.FindByAccountId(DefaultAccountMapping.InternalAccountId)
	if err != nil || AccountMapping == nil {
		log.Println("Error searching default account mapping")
		accountRepository.Store(&DefaultAccountMapping)
	}
}

func openDatabaseConnection(config *config.Config) *mgo.Database {
	return util.ConnectToDatabase(config)
}

func configureTrioClient(config *config.Config) domain.TrioClient {
	return adapters.NewTrioHttpClient(config)
}

func configureControllers(server *api.Server) {
	SyncRequestService := configureSyncRequestService()
	BalanceService := configureBalanceService()
	TransactionService := configureTransactionService()
	controller := domain.NewWebHookControllerImpl(SyncRequestService, BalanceService, TransactionService)
	server.ConfigureController(controller)
}

func configureSyncRequestService() domain.SyncRequestService {
	//TODO implement
	return nil
}

func configureBalanceService() domain.BalanceService {
	//TODO implement
	return nil
}

func configureTransactionService() domain.TransactionService {
	//TODO implement
	return nil
}

func configureEventSubscriberService(syncRequestRepository domain.SyncRequestRepository,
	accountRepository domain.AccountRepository, trioClient domain.TrioClient,
	eventDispatcher domain.EventDispatcher) domain.EventSubscriberService {
	return domain.NewEventSubscriberServiceImpl(syncRequestRepository, accountRepository, trioClient, eventDispatcher)
}

func configureConsumer(ctx context.Context, config *config.Config,
	eventSubscriberService domain.EventSubscriberService) *infra.Consumer {
	consumer := infra.NewConsumer(ctx, []string{config.KafkaBroker},
		config.SyncRequestInputTopic, config.KafkaClientId)
	go consumer.StartReading(eventSubscriberService.OnMessageReceive)
	return consumer
}

func startServer(config *config.Config, server *api.Server) {
	err := server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("Failed to start HTTP server")
	}
}
