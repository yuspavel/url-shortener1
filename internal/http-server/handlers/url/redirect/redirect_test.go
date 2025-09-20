package redirect_test

import (
	"net/http/httptest"
	"testing"
	"url-shortener/internal/http-server/handlers/url/redirect"
	"url-shortener/internal/http-server/handlers/url/redirect/mocks"
	"url-shortener/internal/lib/api"
	"url-shortener/internal/lib/logger/handlers/slogdiscard"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
)

func TestRedirect(t *testing.T) {
	cases := []struct {
		name        string
		alias       string
		url         string
		respMessage string
		mockErr     error
		code        int
	}{
		/*{
			name:  "***Successfull test***",
			alias: "bla bla",
			url:   "https://ya.ru",
		},*/
		{
			name:        "***Invalid request***",
			alias:       "",
			url:         "https://lenta.ru",
			respMessage: "invalid request",
			code:        404,
		},
		/*{
			name:        "***URL not exist***",
			alias:       "RqT7UN",
			respMessage: "url not exist in db",
		},*/
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			urlGetter := mocks.NewURLGetter(t)

			if tc.respMessage == "" || tc.mockErr != nil {
				urlGetter.On("GetURL", tc.alias).Return(tc.url, tc.mockErr)
			}

			router := chi.NewRouter()
			router.Get("/{alias}", redirect.New(slogdiscard.NewDiscardLogger(), urlGetter))

			ts := httptest.NewServer(router)
			defer ts.Close()

			redirectedURL, err := api.GetRedirect(ts.URL + "/" + tc.alias)

			require.NoError(t, err)
			//require.Equal(t,tc.code,)

			assert.Equal(t, tc.url, redirectedURL)

		})
	}
}
