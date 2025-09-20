package remove_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	remove "url-shortener/internal/http-server/handlers/url/delete"
	rem "url-shortener/internal/http-server/handlers/url/delete/mocks"
	"url-shortener/internal/http-server/handlers/url/save"
	sav "url-shortener/internal/http-server/handlers/url/save/mocks"

	"url-shortener/internal/lib/logger/handlers/slogdiscard"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRemoveUrl(t *testing.T) {
	cases := []struct {
		name       string
		alias      string
		url        string
		respStatus int
		respMsg    string
		mockErr    error
	}{
		{
			name:       "Success test",
			url:        "http://ya.ru",
			alias:      "QwertY",
			respStatus: http.StatusOK,
		},
		{
			name:       "invalid alias",
			url:        "http://kaz.kz",
			alias:      "",
			respStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			//Для проверки кейсов на удаление, сначала создадим алиас
			urlSaver := sav.NewURLSaver(t) //----------создал mock для URLSaver интерфейса

			urlSaver.On("SaveURL", tc.url, mock.AnythingOfType("string")).Return(int64(1), tc.mockErr).Once() //----настроил ожидания для URLSaver

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaver)        //------------создал хэндлер для urlSaver
			body := fmt.Sprintf(`{"url": "%s","alias": "%s"}`, tc.url, tc.alias) //------------создал JSON тело запроса для URLSaver

			r, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(body))) //----создал новый запрос для отправки JSON тела методом POST
			require.NoError(t, err)

			w := httptest.NewRecorder() //-----------создал recorder для записи ответа URLSaver
			handler.ServeHTTP(w, r)     //------выполнение запроса и получение ответа

			var save_resp save.Response

			err = json.Unmarshal(w.Body.Bytes(), &save_resp)
			require.NoError(t, err)

			//Следующий шаг, пробуем его удалять
			urlDeleter := rem.NewURLDeleter(t) //----------создал mock для URLDeleter интерфейса
			if tc.respMsg == "" || tc.mockErr != nil {
				urlDeleter.On("DeleteURL", tc.alias).Return(tc.mockErr).Once() //--------------настроил ожидания для URLDeleter
			}

			router := chi.NewRouter()
			router.Delete("/{alias}", remove.New(slogdiscard.NewDiscardLogger(), urlDeleter))

			ts := httptest.NewServer(router)
			defer ts.Close()

			req, err := http.NewRequest(http.MethodDelete, ts.URL+"/"+tc.alias, nil)
			require.NoError(t, err)

			client := http.Client{}
			resp, err := client.Do(req)
			require.NoError(t, err)
			defer resp.Body.Close()

			require.Equal(t, tc.respStatus, resp.StatusCode)
		})
	}
}
