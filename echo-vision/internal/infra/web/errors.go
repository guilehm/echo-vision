package web

type apiError struct {
	Status  int
	Message string
}

func (e apiError) Error() string {
	return e.Message
}

func newApiError(status int, message string) *apiError {
	return &apiError{
		Status:  status,
		Message: message,
	}
}
