package save_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/http-server/handlers/url/save/mocks"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

/*func TestSaveUrl(t *testing.T) {
	cases := []struct {
		name    string
		url     string
		alias   string
		respErr string
		mockErr error
	}{
		{
			name:  "Success test",
			url:   "https://ya.ru",
			alias: "my_alias",
		},
		{
			name:  "No url",
			url:   "https://lenta.ru",
			alias: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlSaver := mocks.NewURLSaver(t)

			if tc.respErr == "" || tc.mockErr != nil { //---вне зависимости от результата (была ошибка или нет) формируется ожидание от сервиса
				urlSaver.On("SaveURL", tc.url, mock.AnythingOfType("string")).Return(int64(1), tc.mockErr).Once() //дословно: mock urlSaver вызывает метод SaveURL, с параметрами url и любого string значения в качестве алиаса
				//--------------------------------------------------------------------------------------------------в качестве возвращаемого результата значение int64 и ошибка. Метод SaveURL выполняется один раз.
			}

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaver) //---Создание хэндлера с методом save.New, с передачей ему в качестве параметра, mock'а urlSaver и пустого логера

			input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, tc.url, tc.alias) //Создание тела запроса (JSON)

			req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input))) //Создание нового запроса http request, с передачей ему тела запроса
			require.NoError(t, err)

			resp := httptest.NewRecorder() //Создание Recorder'а, для записи ответа

			handler.ServeHTTP(resp, req) //---Выполнение запроса req и запись ответа в resp

			require.Equal(t, resp.Code, http.StatusOK)

			body := resp.Body.String()

			var response save.Response

			require.NoError(t, json.Unmarshal([]byte(body), &response))
			require.Equal(t, tc.respErr, response.Error)

		})
	}
}*/

func TestSaveSuccess(t *testing.T) {
	urlSaver := mocks.NewURLSaver(t)

	var (
		err   error
		url   string = "https://ya.ru"
		alias string = "bla bla"
	)
	urlSaver.On("SaveURL", url, mock.AnythingOfType("string")).Return(int64(1), err).Once()

	handler := save.New(slogdiscard.NewDiscardLogger(), urlSaver)

	input := fmt.Sprintf(`{"url": "%s", "alias": "%s"}`, url, alias) //---формируем тело запроса в формате JSON

	req, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(input)))
	require.NoError(t, err)

	resp := httptest.NewRecorder()
	require.Equal(t, http.StatusOK, resp.Code)

	handler.ServeHTTP(resp, req)

	var responseStruct save.Response

	body := resp.Body.String() //-----Извлекаем Body из http ответа

	err = json.Unmarshal([]byte(body), &responseStruct)
	require.NoError(t, err)
	urlSaver.AssertNumberOfCalls(t, "SaveURL", 1)
}

func TestSaveFailure(t *testing.T) {
	cases := []struct {
		name     string
		url      string
		alias    string
		body     string
		cnt      int
		id       int64
		respInfo string
		mockErr  error
	}{
		{
			name:  "Success test",
			url:   "https://mail.ru",
			alias: "mail_alias",
			//body:     `{"alias": "mail_alias", "url": "https://mail.ru"}`,
			cnt:      1,
			id:       2,
			respInfo: "OK",
		},
		/*{
			name:     "Empty request body",
			url:      "http://ya.ru",
			alias:    "",
			body:     "",
			cnt:      1,
			respInfo: "empty request",
		},
		{
			name:     "Decode failed",
			url:      "",
			body:     `{"url": "https://ya.ru", alias: ""}`,
			cnt:      1,
			respInfo: "failed to decode request",
		},
		{
			name:     "URL required",
			body:     `{"url": "","alias": "my_alias"}`,
			cnt:      1,
			respInfo: "field URL is a required field",
		},
		{
			name:     "Wrong URL",
			body:     `{"url": "http//ya.ru", "alias": "my_alias"}`,
			cnt:      1,
			respInfo: "field URL is not a URL",
		},*/
		/*{
			name:    "URL exist",
			body:    `{"url": "http://ya.ru", "alias": "my_alias"}`,
			cnt:     2,
			respInfo: "url already exist",
		},*/
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			urlSaver := mocks.NewURLSaver(t)

			if tt.mockErr != nil {
				urlSaver.On("SaveURL", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(mock.AnythingOfType("int64"), tt.mockErr).Times(tt.cnt)
			}
			//urlSaver.Mock.Test(t)

			handler := save.New(slogdiscard.NewDiscardLogger(), urlSaver)
			body := fmt.Sprintf(`{"url": "%s","alias": "%s"}`, tt.url, tt.alias)

			r, err := http.NewRequest(http.MethodPost, "/save", bytes.NewReader([]byte(body)))
			require.NoError(t, err)

			w := httptest.NewRecorder()

			handler.ServeHTTP(w, r)

			var resp save.Response

			err = json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			require.NoError(t, tt.mockErr)
			require.Equal(t, tt.respInfo, resp.Error)
			//require.Equal(t, int64(1), tt.id)

			//require.Equal(t, response.StatusError, resp.Status)
			//require.Equal(t, 1)
			//urlSaver.AssertExpectations(t)
			//urlSaver.AssertNumberOfCalls(t, "SaveURL", tt.cnt)
			//fmt.Println(urlSaver.Calls)

		})
	}

}
