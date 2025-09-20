package save

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	resp "url-shortener/internal/lib/api/response"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/lib/random"
	"url-shortener/internal/storage"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

const aliasLen = 6

type Request struct {
	Alias string `json:"alias,omitempty"`             //------Помимо мапинга полей, указываются и тэги валидации
	URL   string `json:"url" validate:"required,url"` //------Помимо мапинга полей, указываются и тэги валидации
}

// Эта структура будет возвращаться клиенту, если запрос выполнится успешно, при этом поле Error вложенной структуры resp.Response не будет заполняться.
type Response struct {
	resp.Response
	Alias string `json:"alias"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.5 --name=URLSaver
type URLSaver interface {
	SaveURL(URL, alias string) (int64, error)
}

func New(log *slog.Logger, urlSaver URLSaver) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "save.New"

		log = log.With(
			slog.String("", "\n"),
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		var req Request

		err := render.DecodeJSON(r.Body, &req) //-----------------Декодировали тело запроса в из JSON в структуру Request

		//Проверки на случай ошибки декодирования тела запроса
		if err != nil {
			if errors.Is(err, io.EOF) {
				log.Error("request body is empty")
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, resp.Error("empty request"))
				return
			}

			log.Error("failed to decode request body")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("internal error"))
			return
		}

		log.Info("request body decoded", slog.Any("req", req))

		//Валидация декодированного запроса (структуры Request)
		if err = validator.New().Struct(req); err != nil { //-----Метод Struct проводит валидацию структуры Req на соответствие тегам валидации указанных в описании структуры Request
			log.Error("validation error", sl.Err(err))
			if validateErr, ok := err.(validator.ValidationErrors); ok { //-----Приведение к типу ошибки ValidationErrors. Если удачно, то..
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				render.JSON(w, r, resp.ValidateError(validateErr)) //--..рендерим ответ клиенту из структуры resp.Response с указанием ошибки валидации
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error(err.Error())) //----Иначе рендерим ответ клиенту с текстом общей ошибки
			return
		}

		//Создаем наш алиас, по сути это и есть короткая ссылка (суть сервиса url-shortener)
		alias := req.Alias
		if alias == "" {
			alias = random.NewRandomString(aliasLen)
		}

		id, err := urlSaver.SaveURL(alias, req.URL)
		if err != nil {
			if errors.Is(err, storage.ErrURLExists) {
				log.Info("url already exist", slog.String("url", req.URL))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusFound)
				render.JSON(w, r, resp.Error("url already exist"))
				return
			}

			log.Info("failed to add url")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, resp.Error("failed to add url"))
			return
		}

		log.Info("url is added", slog.Int64("id", id))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		ResponseOK(w, r, alias)
	}
}

// Вспомогательная функция, обеспечивающая рендер ответа удачной вставки URL и возврат алиаса пользователю.
func ResponseOK(w http.ResponseWriter, r *http.Request, alias string) {
	render.JSON(w, r, Response{Response: resp.OK(), Alias: alias})
}
