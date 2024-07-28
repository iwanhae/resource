package resource

import "context"

const (
	MIMEApplicationJSON = "application/json"
	HeaderContentType   = "Content-Type"
)

type Resource interface {
	ValidateCreate(ctx Context) error
	ValidateUpdate(ctx Context, id string) error
}

type ResourceList[T Resource] struct {
	Items    []T `json:"items"`
	Metadata `json:"metadata"`
}

type Metadata struct {
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type ErrorResponse struct {
	Message string      `json:"message"`
	Code    int         `json:"code"`
	Error   interface{} `json:"error"`
}

type Context struct {
	context.Context
}

type List[T Resource] func(ctx Context, offset int, limit int) ([]T, error)
type Create[T Resource] func(ctx Context, resource T) (T, error)
type Update[T Resource] func(ctx Context, id string, resource T) (T, error)
type Get[T Resource] func(ctx Context, id string) (T, error)
type Delete[T Resource] func(ctx Context, id string) error
