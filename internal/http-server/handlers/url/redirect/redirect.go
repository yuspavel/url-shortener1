package redirect

import (
	"errors"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=URLGetter
type URLGetter interface {
	GetURL(string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "redirect.New"

		log := log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		alias := chi.URLParam(r, "rdr") //----Вытаскиваем из запроса параметр alias
		if alias == "" {
			log.Info("alias not provided")
			render.JSON(w, r, resp.Error("invalid request"))
			return
		}

		url, err := urlGetter.GetURL(alias)
		if errors.Is(err, storage.ErrURLNotFound) {
			log.Info("url not exist in db", "alias", alias)
			render.JSON(w, r, resp.Error("url not found"))
			return
		}

		if err != nil {
			log.Info("failed to get url", sl.Err(err))
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("url found", slog.String("url", url))
		http.Redirect(w, r, url, http.StatusFound)

	}
}
