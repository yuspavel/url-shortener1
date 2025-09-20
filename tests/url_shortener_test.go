package tests

import (
	"net/http"
	"net/url"
	"testing"
	"url-shortener/internal/http-server/handlers/url/save"
	"url-shortener/internal/lib/random"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
)

const host = "localhost:8082"

func TestSaveURLShortener(t *testing.T) {
	u := url.URL{Scheme: "http",
		Host: host}

	alias := random.NewRandomString(6)
	url := gofakeit.URL()
	req := save.Request{Alias: alias, URL: url}

	e := httpexpect.Default(t, u.String())
	e.POST("/url/save").WithJSON(req).WithBasicAuth("my_user", "my_pass").Expect().Status(200).JSON().Object().ContainsKey("alias")

	//e.POST("/url/save").WithJSON(req).WithBasicAuth("my_user", "my_pass").Expect().Status(200).Body().NotEmpty()

	e.POST("/url/save").WithJSON(req).WithBasicAuth("my_user", "my_pas").Expect().Status(http.StatusUnauthorized)
	e.POST("/url/sav").WithJSON(req).WithBasicAuth("my_user", "my_pass").Expect().Status(http.StatusNotFound)
	e.GET("/url/save").WithJSON(req).WithBasicAuth("my_user", "my_pass").Expect().Status(http.StatusMethodNotAllowed)

	e.POST("/url/save").WithJSON(nil).WithBasicAuth("my_user", "my_pass").Expect().Status(http.StatusBadRequest).Body().IsEqual("{\"status\":\"Error\",\"error\":\"field URL is a required field\"}\n")

	alias = random.NewRandomString(6)
	url = gofakeit.URL()
	req = save.Request{Alias: random.NewRandomString(6), URL: url}
	e.POST("/url/save").WithJSON(req).WithBasicAuth("my_user", "my_pass").Expect().Status(http.StatusOK).JSON().Object().HasValue("status", "OK").HasValue("alias", alias)

}
