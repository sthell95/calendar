package main

import (
	"context"
	"fmt"
	"os"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-lib/metrics"

	"os/signal"

	"github.com/spf13/viper"

	"calendar.com/config"
	"calendar.com/pkg/controller"
	"calendar.com/pkg/domain/repository"
	"calendar.com/pkg/domain/service"
	"calendar.com/pkg/logger"
	"calendar.com/pkg/storage"

	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
)

func main() {
	fmt.Println("-> Running application")

	if err := initConfig(); err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "main")
	}

	grace := make(chan os.Signal, 1)
	signal.Notify(grace, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		logger.NewLogger().Write(logger.Error, fmt.Sprintf("system call: %+v", <-grace), "main")
		cancel()
	}()
	db := storage.NewDB(ctx)

	storages := storage.Storage{Gorm: db}
	eventRepository := repository.NewEventRepository(storages)
	eventService := service.NewEventService(eventRepository)
	userRepository := repository.NewUserRepository(storages)
	authService := service.NewAuthService(userRepository)
	c := controller.NewController(eventService, authService)
	handlers := new(config.Handlers)
	handlers.NewHandler(*c)

	cfg := jaegercfg.Configuration{
		ServiceName: "calendar",
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	jLogger := jaegerlog.StdLogger
	jMetricsFactory := metrics.NullFactory

	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
		jaegercfg.Metrics(jMetricsFactory),
	)
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "serve")
	}

	defer closer.Close()

	opentracing.SetGlobalTracer(tracer)

	err = config.Run(ctx, handlers.NewRouter())
	if err != nil {
		logger.NewLogger().Write(logger.Error, err.Error(), "serve")
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
