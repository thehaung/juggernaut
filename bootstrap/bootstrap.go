package bootstrap

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/thehaung/juggernaut/config"
	"github.com/thehaung/juggernaut/delivery/http"
	tgHandler "github.com/thehaung/juggernaut/delivery/telegram"
	"github.com/thehaung/juggernaut/internal/httpserver"
	"github.com/thehaung/juggernaut/internal/logger"
	"github.com/thehaung/juggernaut/internal/telegram"
	"github.com/thehaung/juggernaut/repository"
	"github.com/thehaung/juggernaut/usecase"

	"os"
	"os/signal"
	"syscall"
)

func Run(conf *config.Config) {
	// Third Party
	telegramService := telegram.NewTelegram(conf)

	// Repository
	juggernautRepository := repository.NewJuggernautRepository()

	// UseCase
	juggernautUseCase := usecase.NewJuggernautUseCase(juggernautRepository)

	// Handler
	tgHandler.NewTelegramHandler(juggernautUseCase, telegramService.GetTelegramBotAPI(), telegramService.GetUpdateChannel())

	// Router
	router := chi.NewRouter()
	http.NewRoute(router, conf, juggernautUseCase)

	// HttpServer
	newHttpServer(router, conf)
}

func newHttpServer(router *chi.Mux, conf *config.Config) {
	log := logger.GetLogger()
	httpServer := httpserver.New(router, httpserver.Port(conf.HttpServer.Port))
	// Waiting signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Warnf("app - Run - signal: %s", s.String())
	case err := <-httpServer.Notify():
		log.Fatal(fmt.Errorf("app - Run - httpServer.Notify: %w", err))
	}

	// Shutdown
	err := httpServer.Shutdown()
	if err != nil {
		log.Fatal(fmt.Errorf("app - Run - httpServer.Shutdown: %w", err))
	}
}
