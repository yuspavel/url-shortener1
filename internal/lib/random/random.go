package random

import (
	"math/rand"
	"time"
)

func NewRandomString(size int) string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ" +
		"abcdefghijklmnopqrstuvwxyz" +
		"0123456789")

	b := make([]rune, size)

	for i := range b {
		b[i] = chars[rnd.Intn(len(chars))] //----Дословно: в каждый элемент b[i] присваиваем случайное число из диапазона чисел генератора,
		// --------------------------------------выбирая из него каждый раз разное подмножество чисел равное len(chars)
	}
	return string(b)
}
