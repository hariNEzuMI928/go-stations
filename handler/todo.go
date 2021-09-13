package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/TechBowl-japan/go-stations/model"
	"github.com/TechBowl-japan/go-stations/service"
)

// A TODOHandler implements handling REST endpoints.
type TODOHandler struct {
	svc *service.TODOService
}

// NewTODOHandler returns TODOHandler based http.Handler.
func NewTODOHandler(svc *service.TODOService) *TODOHandler {
	return &TODOHandler{
		svc: svc,
	}
}

// Create handles the endpoint that creates the TODO.
func (h *TODOHandler) Create(ctx context.Context, req *model.CreateTODORequest) (*model.CreateTODOResponse, error) {
	res, err := h.svc.CreateTODO(ctx, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	return &model.CreateTODOResponse{TODO: *res}, nil
}

// Read handles the endpoint that reads the TODOs.
func (h *TODOHandler) Read(ctx context.Context, req *model.ReadTODORequest) (*model.ReadTODOResponse, error) {
	res, err := h.svc.ReadTODO(ctx, req.PrevID, req.Size)
	if err != nil {
		return &model.ReadTODOResponse{}, err
	}

	return &model.ReadTODOResponse{TODOs: res}, nil
}

// Update handles the endpoint that updates the TODO.
func (h *TODOHandler) Update(ctx context.Context, req *model.UpdateTODORequest) (*model.UpdateTODOResponse, error) {
	res, err := h.svc.UpdateTODO(ctx, req.ID, req.Subject, req.Description)
	if err != nil {
		return nil, err
	}

	return &model.UpdateTODOResponse{TODO: *res}, nil
}

// Delete handles the endpoint that deletes the TODOs.
func (h *TODOHandler) Delete(ctx context.Context, req *model.DeleteTODORequest) (*model.DeleteTODOResponse, error) {
	_ = h.svc.DeleteTODO(ctx, nil)
	return &model.DeleteTODOResponse{}, nil
}

// ServeHTTP implements http.Handler interface.
func (h *TODOHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		req := model.ReadTODORequest{}
		var err error

		req.PrevID, err = strconv.ParseInt(r.URL.Query().Get("prev_id"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		req.Size, err = strconv.ParseInt(r.URL.Query().Get("size"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		res, err := h.Read(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(res)
	case http.MethodPost:
		req := model.CreateTODORequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return
		}
		if req.Subject == "" {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res, err := h.Create(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(res)
	case http.MethodPut:
		req := model.UpdateTODORequest{}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			return
		}
		if req.Subject == "" || req.ID == 0 {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		res, err := h.Update(r.Context(), &req)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		json.NewEncoder(w).Encode(res)
	}
}
