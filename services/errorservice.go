package services

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	isError bool   `json:"error"`
}

func (e *AppError) Error() string {
	return e.Message
}
