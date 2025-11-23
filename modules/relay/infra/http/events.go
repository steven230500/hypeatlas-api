package http

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/steven230500/hypeatlas-api/domain/entities"
	in "github.com/steven230500/hypeatlas-api/modules/relay/domain/ports/in"
	sharedhttp "github.com/steven230500/hypeatlas-api/shared/http"
)

type EventHandler struct{ svc in.Service }

func NewEventHandler(s in.Service) *EventHandler { return &EventHandler{svc: s} }

func (h *EventHandler) Register(r chi.Router) {
	r.Route("/relay/events", func(r chi.Router) {
		r.Get("/", h.list)
		r.Get("/{slug}", h.get)
		r.Post("/", h.create)
		r.Put("/{slug}", h.update)
		r.Delete("/{slug}", h.delete)
	})
}

// @Summary      Listar eventos
// @Tags         relay
// @Param        game   query string false "val|lol"
// @Param        league query string false "Filtrar por liga"
// @Param        status query string false "upcoming|live|past"
// @Param        limit  query int    false "Number of items per page (default: 20, max: 100)"
// @Param        offset query int    false "Number of items to skip (default: 0)"
// @Produce      json
// @Success      200 {object} sharedhttp.PaginatedResponse
// @Failure      500 {string} string "error"
// @Router       /v1/relay/events [get]
func (h *EventHandler) list(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	params := sharedhttp.ParsePaginationParams(r)

	events, total, err := h.svc.ListEvents(
		r.Context(),
		q.Get("game"),
		q.Get("league"),
		q.Get("status"),
		params.Limit,
		params.Offset,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(sharedhttp.PaginatedResponse{
		Items: events,
		Pagination: sharedhttp.PaginationMeta{
			Limit:   params.Limit,
			Offset:  params.Offset,
			Total:   total,
			HasMore: params.Offset+params.Limit < total,
		},
	})
}

// @Summary      Obtener evento por slug
// @Tags         relay
// @Param        slug path string true "Event slug"
// @Produce      json
// @Success      200 {object} entities.Event
// @Failure      404 {string} string "event not found"
// @Router       /v1/relay/events/{slug} [get]
func (h *EventHandler) get(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	event, err := h.svc.GetEvent(r.Context(), slug)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(event)
}

// @Summary      Crear evento
// @Tags         relay
// @Accept       json
// @Produce      json
// @Param        event body entities.Event true "Event data"
// @Success      201 {object} entities.Event
// @Failure      400 {string} string "validation error"
// @Router       /v1/relay/events [post]
func (h *EventHandler) create(w http.ResponseWriter, r *http.Request) {
	var event entities.Event
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.CreateEvent(r.Context(), &event); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(event)
}

// @Summary      Actualizar evento
// @Tags         relay
// @Accept       json
// @Produce      json
// @Param        slug path string true "Event slug"
// @Param        updates body map[string]interface{} true "Fields to update"
// @Success      200 {string} string "ok"
// @Failure      400 {string} string "validation error"
// @Router       /v1/relay/events/{slug} [put]
func (h *EventHandler) update(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.svc.UpdateEvent(r.Context(), slug, updates); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// @Summary      Eliminar evento
// @Tags         relay
// @Param        slug path string true "Event slug"
// @Success      200 {string} string "ok"
// @Failure      400 {string} string "error"
// @Router       /v1/relay/events/{slug} [delete]
func (h *EventHandler) delete(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	if err := h.svc.DeleteEvent(r.Context(), slug); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}
