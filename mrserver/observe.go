package mrserver

import (
	"net/http"
	"time"
)

type (
	// RequestStat - интерфейс для сбора и отправки статистики HTTP-запросов.
	RequestStat interface {
		Enabled() bool
		Emit(r *http.Request, body []byte, size int, responseBody []byte, responseSize int, duration time.Duration, status int)
	}

	// nopRequestStat - заглушка, реализующая интерфейс RequestStat.
	// Игнорирует данные для сбора статистики.
	nopRequestStat struct{}
)

// NopRequestStat - создаёт RequestStat, который игнорирует все данные для сбора статистики.
func NopRequestStat() RequestStat {
	return nopRequestStat{}
}

// Enabled - всегда возвращает false.
func (rs nopRequestStat) Enabled() bool {
	return false
}

// Emit - имитирует сбор данных.
func (rs nopRequestStat) Emit(_ *http.Request, _ []byte, _ int, _ []byte, _ int, _ time.Duration, _ int) {
}
