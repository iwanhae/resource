package resource

import (
	"encoding/json"
	"net/http"
)

func (b *Builder[T]) handlerList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
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
	ctx := r.Context()
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
