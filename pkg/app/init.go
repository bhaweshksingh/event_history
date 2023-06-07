package app

import (
	"event-history/pkg/config"
	"event-history/pkg/eventinfo"
	"event-history/pkg/http/router"
	"event-history/pkg/http/server"
	"event-history/pkg/reporters"
	"event-history/pkg/repository"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"os"
)

func initHTTPServer(configFile string) {
	config := config.NewConfig(configFile)
	logger := initLogger(config)
	rt := initRouter(config, logger)

	server.NewServer(config, logger, rt).Start()
}

func initRouter(cfg config.Config, logger *zap.Logger) http.Handler {
	eventRepo := initRepository(cfg)
	eventService := initService(eventRepo)

	return router.NewRouter(logger, eventService)
}

func initService(eventRepository repository.EventRepository) eventinfo.Service {
	eventService := eventinfo.NewEventService(eventRepository)

	return eventService
}

func initRepository(cfg config.Config) repository.EventRepository {
	dbConfig := cfg.GetDBConfig()
	dbHandler := repository.NewDBHandler(dbConfig)

	db, err := dbHandler.GetDB()
	if err != nil {
		log.Fatal(err.Error())
	}

	return repository.NewEventRepository(db)
}

func initLogger(cfg config.Config) *zap.Logger {
	return reporters.NewLogger(
		cfg.GetLogConfig().GetLevel(),
		getWriters(cfg.GetLogFileConfig())...,
	)
}

func getWriters(cfg config.LogFileConfig) []io.Writer {
	return []io.Writer{
		os.Stdout,
		reporters.NewExternalLogFile(cfg),
	}
}
