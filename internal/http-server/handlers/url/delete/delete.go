package remove

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=URLDeleter
type URLDeleter interface {
	DeleteURL(string) error
}

func New(log *slog.Logger, urlDeleter URLDeleter) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "delete.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "del")
		if alias == "" {
			log.Info("invalid alias")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, resp.Error("invalid alias"))
			return
		}

		err := urlDeleter.DeleteURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("removing url error")
			w.WriteHeader(http.StatusNotFound)
			render.JSON(w, r, resp.Error(storage.ErrURLNotFound.Error()))
			return
		}
		if err != nil {
			log.Info("internal error")
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("url " + r.URL.String() + " removed")
		w.WriteHeader(http.StatusOK)     //----устанавливаем статус в заголовок ответа
		render.JSON(w, r, resp.StatusOK) //---рендерим структуру Response ответ в JSON ответ
	}
}
