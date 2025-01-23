package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-order/config"
	v1HTTP "github.com/idoyudha/eshop-order/internal/controller/http/v1"
	v1Kafka "github.com/idoyudha/eshop-order/internal/controller/kafka/v1"
	"github.com/idoyudha/eshop-order/internal/event"
	"github.com/idoyudha/eshop-order/internal/usecase"
	"github.com/idoyudha/eshop-order/internal/usecase/commandrepo"
	"github.com/idoyudha/eshop-order/internal/usecase/queryrepo"
	"github.com/idoyudha/eshop-order/pkg/httpserver"
	"github.com/idoyudha/eshop-order/pkg/kafka"
	"github.com/idoyudha/eshop-order/pkg/logger"
	"github.com/idoyudha/eshop-order/pkg/postgresql/postgrecommand"
	"github.com/idoyudha/eshop-order/pkg/postgresql/postgrequery"
	"github.com/idoyudha/eshop-order/pkg/redis"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

	kafkaConsumer, err := kafka.NewKafkaConsumer(cfg.Kafka.Broker)
	if err != nil {
		l.Fatal("app - Run - kafka.NewKafkaConsumer: ", err)
	}
	defer kafkaConsumer.Close()

	kafkaProducer, err := kafka.NewKafkaProducer(cfg.Kafka.Broker)
	if err != nil {
		l.Fatal("app - Run - kafka.NewKafkaProducer: ", err)
	}
	defer kafkaProducer.Close()

	postgreSQLCommand, err := postgrecommand.NewPostgres(cfg.PostgreSQLCommand)
	if err != nil {
		l.Fatal("app - Run - postgrecommand.NewPostgres: ", err)
	}

	postgreSQLQuery, err := postgrequery.NewPostgres(cfg.PostgreSQLQuery)
	if err != nil {
		l.Fatal("app - Run - postgrequery.NewPostgres: ", err)
	}

	redisClient, err := redis.NewRedis(cfg.Redis)
	if err != nil {
		l.Fatal("app - Run - redis.NewRedis: ", err)
	}

	orderCommandUseCase := usecase.NewOrderCommandUseCase(
		commandrepo.NewOrderPostgreCommandRepo(postgreSQLCommand),
		queryrepo.NewOrderPostgreQueryRepo(postgreSQLQuery),
		commandrepo.NewOrderRedisRepo(redisClient),
		kafkaProducer,
		cfg.WarehouseService,
		cfg.ShippingCostService,
		cfg.Constant,
	)

	orderQueryUseCase := usecase.NewOrderQueryUseCase(
		queryrepo.NewOrderPostgreQueryRepo(postgreSQLQuery),
	)

	// HTTP Server
	handler := gin.Default()
	v1HTTP.NewRouter(handler, orderQueryUseCase, orderCommandUseCase, l, cfg.AuthService)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

	// Kafka Consumer
	kafkaErrChan := make(chan error, 1)
	go func() {
		if err := v1Kafka.KafkaNewRouter(orderQueryUseCase, orderCommandUseCase, l, kafkaConsumer, cfg.ProductService); err != nil {
			kafkaErrChan <- err
		}
	}()

	// Redis Consumer
	redisErrChan := make(chan error, 1)
	go func() {
		if err := event.NewRedisScheduledEvents(redisClient, orderCommandUseCase, l, cfg.Constant); err != nil {
			redisErrChan <- err
		}
	}()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		l.Info("app - Run - signal: %s", s.String())
	case err = <-httpServer.Notify():
		l.Error("app - Run - httpServer.Notify: ", err)
	}

	// Shutdown
	err = httpServer.Shutdown()
	if err != nil {
		l.Info("app - Run - httpServer.Shutdown: %s", err)
	}
}
