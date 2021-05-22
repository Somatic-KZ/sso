package models

type ErrorResponse struct {
	Err            error             `json:"-"` // низкоуровневая ошибка исполнения
	HTTPStatusCode int               `json:"-"` // HTTP статус код
	ErrorMessage   *ErrorDetails     `json:"error"`
	Validation     map[string]string `json:"validation,omitempty"` // ошибки валидации
}

type ErrorDetails struct {
	StatusText  string `json:"status"`            // сообщение пользовательского уровня
	AppCode     int64  `json:"code,omitempty"`    // application-определенный код ошибки
	MessageText string `json:"message,omitempty"` // application-level сообщение, для дебага
}
