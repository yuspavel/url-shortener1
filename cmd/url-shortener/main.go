package main

import (
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"url-shortener/internal/config"
	remove "url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/get"
	"url-shortener/internal/http-server/handlers/url/redirect"
	"url-shortener/internal/http-server/handlers/url/save"
	mdv "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/storage/sqlite"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)
	log = log.With(slog.String("ENV", cfg.Env))                          //-------К каждому сообщению  логгер будет добавлять идентификатор той среды, которое его создала
	log.Info("initializing server", slog.String("address", cfg.Address)) //-------При старте сообщения будет выводиться адрес сервера
	log.Debug("logger debug mode enabled")                               //-------Включение режима debug

	storage, err := sqlite.New(cfg.StoragePath)

	if err != nil {
		log.Error("failed to initialize storage:", err)
	}
	stop := make(chan os.Signal, 1)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(mdv.New(log)) //------------Собственная реализация логера
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	//Создаем отдельный подмаршрут /url, в котором все маршруты будут использовать базовую HTTP аутентификацию
	router.Route("/url", func(r chi.Router) {
		r.Use(middleware.BasicAuth("url-shortener", map[string]string{
			cfg.HTTPServer.User: cfg.HTTPServer.Password,
		}))

		r.Post("/save", save.New(log, storage)) //------------------/url/save
	})

	//Остальные маршруты остаются
	//router.Post("/save", save.New(log, storage))
	router.Get("/{rdr}", redirect.New(log, storage))
	router.Delete("/{del}", remove.New(log, storage))
	router.Get("/{get}", get.New(log, storage))

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.Timeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Error("failed to start server")

		}
	}()

	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<-stop

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case "envLocal":
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "envDev":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case "envProd":
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}
