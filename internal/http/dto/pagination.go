package dto

type PageResponse[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Size  int   `json:"size"`
}

type OrdinaryResponse[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
}
