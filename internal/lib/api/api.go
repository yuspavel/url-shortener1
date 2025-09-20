package api

import (
	"errors"
	"fmt"
	"net/http"
)

var ErrInvalidStatusCode = errors.New("invalid status code")

func GetRedirect(url string) (string, error) {
	const op = "api.GetRedirect"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse //-----------------------------------Переопределение редиректа клиента. Говорим ему, что никуда редиректить не надо, просто верни последний ответ запроса
		},
	}

	resp, err := client.Get(url) //------Выполнение запроса по переданному в функцию url
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound { //-----Проверка StatusCode ответа на совпадение с кодом http.StatusFound
		return "", fmt.Errorf("%s: %v, %d", op, ErrInvalidStatusCode, resp.StatusCode)
	}

	return resp.Header.Get("Location"), nil //----------------Возврат URL редиректа в результат функции
}
