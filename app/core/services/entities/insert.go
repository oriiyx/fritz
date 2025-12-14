package entities

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oriiyx/fritz/app/core/api/common/errhandler"
	l "github.com/oriiyx/fritz/app/core/api/common/log"
	"github.com/oriiyx/fritz/app/core/services/entities/adapters"
	ctxUtil "github.com/oriiyx/fritz/app/core/utils/ctx"
	db "github.com/oriiyx/fritz/database/generated"
)

type InsertEntityRequest struct {
	ParentID  *string                `json:"parent_id,omitempty"`
	Key       string                 `json:"key" validate:"required,max=255"`
	Path      string                 `json:"path" validate:"required"`
	Type      string                 `json:"type,omitempty"`
	Published bool                   `json:"published"`
	Data      map[string]interface{} `json:"data" validate:"required"`
}

// InsertEntity creates a new entity instance
func (h *Handler) InsertEntity(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())
	classID := chi.URLParam(r, ClassIDKey)

	var req InsertEntityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Msg("Failed to decode request")
		errhandler.BadRequest(w, errhandler.RespInvalidRequestBody)
		return
	}

	// Validate request
	if err := h.Validator.Struct(req); err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Msg("Validation failed")
		errhandler.BadRequest(w, errhandler.RespInvalidRequestBody)
		return
	}

	// Get the adapter for this entity class
	adapter, err := adapters.Get(classID)
	if err != nil {
		h.Logger.Error().Err(err).Str("class_id", classID).Msg("Unknown entity class")
		errhandler.BadRequest(w, []byte(`{"error": "unknown entity class"}`))
		return
	}

	// Prepare entity creation params
	entityType := "object"
	if req.Type != "" {
		entityType = req.Type
	}

	var parentID pgtype.UUID
	if req.ParentID != nil {
		if err := parentID.Scan(*req.ParentID); err != nil {
			h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Msg("Invalid parent_id")
			errhandler.BadRequest(w, []byte(`{"error": "invalid parent_id"}`))
			return
		}
	}

	// TODO: Get user ID from session/context
	var userID pgtype.UUID

	// Create entity record in entities table
	entityParams := db.CreateEntityParams{
		EntityClass: classID,
		ParentID:    parentID,
		OKey:        req.Key,
		OPath:       req.Path,
		OType:       entityType,
		Published:   req.Published,
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	entity, err := h.Queries.CreateEntity(r.Context(), entityParams)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to create entity record")
		errhandler.ServerError(w, errhandler.RespDBDataInsertFailure)
		return
	}

	// Create the entity data using the adapter
	result, err := adapter.Create(r.Context(), entity.ID, req.Data)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to create entity data")
		// Rollback: delete the entity record
		// TODO: Consider using transactions
		errhandler.ServerError(w, errhandler.RespDBDataInsertFailure)
		return
	}

	response := map[string]interface{}{
		"entity": entity,
		"data":   result,
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}
