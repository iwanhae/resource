package resource

import (
	"fmt"
	"net/http"
)

type Builder[T Resource] struct {
	name string

	list   List[T]
	create Create[T]
	get    Get[T]
	update Update[T]
	delete Delete[T]

	defaultLimits int
}

func New[T Resource]() *Builder[T] {
	return &Builder[T]{
		defaultLimits: 10,
	}
}

func (r *Builder[T]) Name(name string) *Builder[T] {
	r.name = name
	return r
}

func (r *Builder[T]) List(f List[T]) *Builder[T] {
	r.list = f
	return r
}

func (r *Builder[T]) Create(f Create[T]) *Builder[T] {
	r.create = f
	return r
}

func (r *Builder[T]) Get(f Get[T]) *Builder[T] {
	r.get = f
	return r
}

func (r *Builder[T]) Update(f Update[T]) *Builder[T] {
	r.update = f
	return r
}

func (r *Builder[T]) Delete(f Delete[T]) *Builder[T] {
	r.delete = f
	return r
}

func (b *Builder[T]) Mux() *http.ServeMux {
	mux := http.NewServeMux()

	if b.list != nil {
		pattern := fmt.Sprintf("GET /%s/", b.name)
		mux.HandleFunc(pattern, b.handlerList)
	}
	if b.create != nil {
		pattern := fmt.Sprintf("POST /%s/", b.name)
		mux.HandleFunc(pattern, b.handlerCreate)
	}
	if b.get != nil {
		pattern := fmt.Sprintf("GET /%s/{resource-id}", b.name)
		mux.HandleFunc(pattern, b.handlerGet)
	}
	if b.update != nil {
		pattern := fmt.Sprintf("PUT /%s/{resource-id}", b.name)
		mux.HandleFunc(pattern, b.handlerUpdate)
	}
	if b.delete != nil {
		pattern := fmt.Sprintf("DELETE /%s/{resource-id}", b.name)
		mux.HandleFunc(pattern, b.handlerDelete)
	}

	return mux
}
