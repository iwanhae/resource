package resource

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	MIMEApplicationJSON = "application/json"
	HeaderContentType   = "Content-Type"
)

type Validator interface {
	ValidateCreate(ctx Context) error
	ValidateUpdate(ctx Context, id string) error
}

type ResourceList[T Validator] struct {
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

type List[T Validator] func(ctx Context, offset int, limit int) ([]T, error)
type Create[T Validator] func(ctx Context, resource T) (T, error)
type Update[T Validator] func(ctx Context, id string, resource T) (T, error)
type Get[T Validator] func(ctx Context, id string) (T, error)
type Delete[T Validator] func(ctx Context, id string) error

type Resource[T Validator] struct {
	name   string
	plural string

	// base path for the resource endpoints e.g. /api/v1/
	base string

	list   List[T]
	create Create[T]
	get    Get[T]
	update Update[T]
	delete Delete[T]

	// default limits for list requests
	defaultLimits int
}

func New[T Validator]() *Resource[T] {
	return &Resource[T]{
		defaultLimits: 10,
	}
}

func (r *Resource[T]) Name(name string) *Resource[T] {
	r.name = name
	return r
}

func (r *Resource[T]) Plural(plural string) *Resource[T] {
	r.plural = plural
	return r
}

func (r *Resource[T]) List(f List[T]) *Resource[T] {
	r.list = f
	return r
}

func (r *Resource[T]) Create(f Create[T]) *Resource[T] {
	r.create = f
	return r
}

func (r *Resource[T]) Get(f Get[T]) *Resource[T] {
	r.get = f
	return r
}

func (r *Resource[T]) Update(f Update[T]) *Resource[T] {
	r.update = f
	return r
}

func (r *Resource[T]) Delete(f Delete[T]) *Resource[T] {
	r.delete = f
	return r
}

func (b *Resource[T]) RegisterMux(mux *http.ServeMux) *Resource[T] {
	if b.list != nil {
		pattern := fmt.Sprintf("GET %s/%s", b.base, b.plural)
		mux.HandleFunc(pattern, b.handlerList)
	}
	if b.create != nil {
		pattern := fmt.Sprintf("POST %s/%s", b.base, b.plural)
		mux.HandleFunc(pattern, b.handlerCreate)
	}
	if b.get != nil {
		pattern := fmt.Sprintf("GET %s/%s/{%s}", b.base, b.plural, b.pathID())
		mux.HandleFunc(pattern, b.handlerGet)
	}
	if b.update != nil {
		pattern := fmt.Sprintf("PUT %s/%s/{%s}", b.base, b.plural, b.pathID())
		mux.HandleFunc(pattern, b.handlerUpdate)
	}
	if b.delete != nil {
		pattern := fmt.Sprintf("DELETE %s/%s/{%s}", b.base, b.plural, b.pathID())
		mux.HandleFunc(pattern, b.handlerDelete)
	}

	return b
}

func (b *Resource[T]) handlerList(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r)
	limit, err := parseParamsInt(r, "limit", b.defaultLimits)
	if err != nil {
		JSONError(w, http.StatusBadRequest, err)
		return
	}
	offset, err := parseParamsInt(r, "offset", 0)
	if err != nil {
		JSONError(w, http.StatusBadRequest, err)
		return
	}
	result, err := b.list(ctx, offset, limit)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err)
		return
	}
	JSON(w, http.StatusOK, ResourceList[T]{
		Items: result,
		Metadata: Metadata{
			Offset: offset,
			Limit:  limit,
		},
	})
}

func (b *Resource[T]) pathID() string {
	return fmt.Sprintf("%sId", b.name)
}

func (b *Resource[T]) handlerCreate(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r)
	body := make([]T, 1)[0]
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		JSONError(w, http.StatusBadRequest, err)
		return
	}
	if err := body.ValidateCreate(ctx); err != nil {
		JSONError(w, http.StatusBadRequest, err)
		return
	}
	result, err := b.create(ctx, body)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err)
		return
	}
	JSON(w, http.StatusCreated, result)
}

func (b *Resource[T]) handlerGet(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r)
	id := r.PathValue(b.pathID())
	if id == "" {
		JSONError(w, http.StatusBadRequest, fmt.Errorf("missing resource-id"))
		return
	}
	result, err := b.get(ctx, id)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err)
		return
	}
	JSON(w, http.StatusOK, result)
}

func (b *Resource[T]) handlerUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r)
	id := r.PathValue(b.pathID())
	if id == "" {
		JSONError(w, http.StatusBadRequest, fmt.Errorf("missing resource-id"))
		return
	}
	body := make([]T, 1)[0]
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		JSONError(w, http.StatusBadRequest, err)
		return
	}
	if err := body.ValidateUpdate(ctx, id); err != nil {
		JSONError(w, http.StatusBadRequest, err)
		return
	}
	result, err := b.update(ctx, id, body)
	if err != nil {
		JSONError(w, http.StatusInternalServerError, err)
		return
	}
	JSON(w, http.StatusOK, result)
}

func (b *Resource[T]) handlerDelete(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r)
	id := r.PathValue(b.pathID())
	if id == "" {
		JSONError(w, http.StatusBadRequest, fmt.Errorf("missing resource-id"))
		return
	}
	if err := b.delete(ctx, id); err != nil {
		JSONError(w, http.StatusInternalServerError, err)
		return
	}
	JSON(w, http.StatusNoContent, nil)
}
