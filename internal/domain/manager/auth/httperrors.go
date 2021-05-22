package auth

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

type Response struct {
	Err            error             `json:"-"` // низкоуровневая ошибка исполнения
	HTTPStatusCode int               `json:"-"` // HTTP статус код
	ErrorMessage   *Details          `json:"error"`
	Validation     map[string]string `json:"validation,omitempty"` // ошибки валидации
}

type Details struct {
	StatusText  string `json:"status"`            // сообщение пользовательского уровня
	Penalty     int    `json:"penalty"`           // пенальти по времени на восстановление
	AppCode     int64  `json:"code,omitempty"`    // application-определенный код ошибки
	MessageText string `json:"message,omitempty"` // application-level сообщение, для дебага
}

func (e *Response) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func TooManyRequests(err error, penalty int) render.Renderer {
	return &Response{
		Err:            err,
		HTTPStatusCode: http.StatusTooManyRequests,
		ErrorMessage: &Details{
			AppCode:     http.StatusTooManyRequests,
			Penalty:     penalty,
			StatusText:  fmt.Sprintf("Too many requests, try in %d seconds", penalty),
			MessageText: err.Error(),
		},
	}
}
