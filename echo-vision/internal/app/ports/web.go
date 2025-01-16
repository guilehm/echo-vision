package ports

type ApiListResponse[T any] struct {
	Results []*T `json:"results"`
	// TODO: implement cursor
}
