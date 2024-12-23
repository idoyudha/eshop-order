package app

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/idoyudha/eshop-order/config"
	v1HTTP "github.com/idoyudha/eshop-order/internal/controller/http/v1"
	"github.com/idoyudha/eshop-order/internal/usecase"
	"github.com/idoyudha/eshop-order/internal/usecase/commandrepo"
	"github.com/idoyudha/eshop-order/internal/usecase/queryrepo"
	"github.com/idoyudha/eshop-order/pkg/httpserver"
	"github.com/idoyudha/eshop-order/pkg/kafka"
	"github.com/idoyudha/eshop-order/pkg/logger"
	"github.com/idoyudha/eshop-order/pkg/postgresql/postgrecommand"
	"github.com/idoyudha/eshop-order/pkg/postgresql/postgrequery"
)

func Run(cfg *config.Config) {
	l := logger.New(cfg.Log.Level)

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

	orderCommandUseCase := usecase.NewOrderCommandUseCase(
		commandrepo.NewOrderPostgreCommandRepo(postgreSQLCommand),
		kafkaProducer,
		cfg.WarehouseService,
	)

	orderQueryUseCase := usecase.NewOrderQueryUseCase(
		queryrepo.NewOrderPostgreCommandRepo(postgreSQLQuery),
	)

	// HTTP Server
	handler := gin.Default()
	v1HTTP.NewRouter(handler, orderQueryUseCase, orderCommandUseCase, l, cfg.AuthService)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))

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
