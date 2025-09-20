package get

import (
	"errors"
	"log/slog"
	"net/http"
	"url-shortener/internal/lib/api/response"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type URLGetter interface {
	GetURL(string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "get.New"
		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "get")
		if alias == "" {
			log.Info("invalid alias")
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, response.Error("invalid alias"))
			return
		}

		url, err := urlGetter.GetURL(alias)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("url not found")
				w.WriteHeader(http.StatusNotFound)
				render.JSON(w, r, response.Error("_url not found"))
				return
			}

			log.Info("internal server error")
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, response.Error("internal server error"))
			return
		}

		log.Info("url found")
		w.WriteHeader(http.StatusOK)
		render.JSON(w, r, url)
	}
}
