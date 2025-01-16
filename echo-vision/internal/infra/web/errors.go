package web

type apiError struct {
	Status  int
	Message string
}

func (e apiError) Error() string {
	return e.Message
}
