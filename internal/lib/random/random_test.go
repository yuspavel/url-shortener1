package random

import (
	"testing"
	"time"

	"github.com/go-playground/assert/v2"
)

func TestRandomString(t *testing.T) {

	tests := []struct {
		name string
		size int
	}{
		{
			name: "Test size=1",
			size: 1,
		},
		{
			name: "Test size=5",
			size: 5,
		},
		{
			name: "Test size=8",
			size: 8,
		},
		{
			name: "Test size=15",
			size: 15,
		},
		{
			name: "Test size 23",
			size: 23,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str1 := NewRandomString(tt.size)
			time.Sleep(1 * time.Nanosecond) //---Засыпание установлено, т.к. последовательный вызов функции NewRandomString слишком быстрый, что не обеспечивает уникальность строк
			str2 := NewRandomString(tt.size)

			assert.Equal(t, len(str1), tt.size) //--Проверка на одинаковый размер сгенерированной строки и размера передаваемого в качестве параметра функции генерации
			assert.Equal(t, len(str2), tt.size)

			assert.Equal(t, len(str1), len(str2)) //--Проверка, что две сгенерированные строки имеют одинаковый размер
			assert.NotEqual(t, str1, str2)        //--Проверка, что сгенерированные строки не совпадают

		})
	}
}
