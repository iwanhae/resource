package resource

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func newContext(r *http.Request) Context {
	return Context{
		Context: r.Context(),
	}
}

func (b *Builder[T]) handlerList(w http.ResponseWriter, r *http.Request) {
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

func (b *Builder[T]) handlerCreate(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r)
	body := make([]T, 1)[0]
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
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

func (b *Builder[T]) handlerGet(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r)
	id := r.PathValue("resource-id")
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

func (b *Builder[T]) handlerUpdate(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r)
	id := r.PathValue("resource-id")
	if id == "" {
		JSONError(w, http.StatusBadRequest, fmt.Errorf("missing resource-id"))
		return
	}
	body := make([]T, 1)[0]
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
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

func (b *Builder[T]) handlerDelete(w http.ResponseWriter, r *http.Request) {
	ctx := newContext(r)
	id := r.PathValue("resource-id")
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
