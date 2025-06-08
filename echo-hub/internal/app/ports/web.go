package ports

type ApiListResponse[T any] struct {
	Results    []*T   `json:"results"`
	NextCursor string `json:"nextCursor"`
}
