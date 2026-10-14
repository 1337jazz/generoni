package item

import (
	"net/http"
	"strconv"

	"github.com/1337jazz/generoni/internal/adapter/http/respond"
	"github.com/1337jazz/generoni/internal/core/port"
)

type ItemHandler struct {
	service port.ItemService
}

func NewItemHandler(service port.ItemService) *ItemHandler {
	return &ItemHandler{service: service}
}

func (h *ItemHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := validateGetItem(r.PathValue("id"))
	if err != nil {
		respond.Error(w, r, err)
		return
	}

	found, err := h.service.Get(r.Context(), id)
	if err != nil {
		respond.Error(w, r, err)
		return
	}

	respond.JSON(w, r, http.StatusOK, itemToResponse(found))
}

func (h *ItemHandler) List(w http.ResponseWriter, r *http.Request) {
	page, err := validateListItems(r.URL.Query())
	if err != nil {
		respond.Error(w, r, err)
		return
	}

	items, total, err := h.service.List(r.Context(), page)
	if err != nil {
		respond.Error(w, r, err)
		return
	}

	respond.JSON(w, r, http.StatusOK, itemPageToResponse(items, total))
}

func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createItemRequest
	if err := respond.Decode(w, r, &req); err != nil {
		respond.Error(w, r, err)
		return
	}

	created, err := h.service.Create(r.Context(), req.Title)
	if err != nil {
		respond.Error(w, r, err)
		return
	}

	w.Header().Set("Location", "/items/"+strconv.FormatInt(created.ID, 10))
	respond.JSON(w, r, http.StatusCreated, itemToResponse(created))
}
