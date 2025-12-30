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
	validatorUtil "github.com/oriiyx/fritz/app/core/utils/validator"
	db "github.com/oriiyx/fritz/database/generated"
)

type SaveEntityRequest struct {
	ID        string                 `json:"id" validate:"required"`
	ParentID  *string                `json:"parent_id,omitempty"`
	Key       string                 `json:"key" validate:"required,max=255"`
	Path      string                 `json:"path" validate:"required"`
	Type      string                 `json:"type,omitempty"`
	Published bool                   `json:"published"`
	Data      map[string]interface{} `json:"data" validate:"required"`
}

// SaveEntity is an endpoint that handles saving entity
func (h *Handler) SaveEntity(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())
	classID := chi.URLParam(r, DefinitionIDKey)

	var req SaveEntityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Msg("Failed to decode request")
		errhandler.BadRequest(w, errhandler.RespInvalidRequestBody)
		return
	}

	// Validate request
	if err := h.Validator.Struct(req); err != nil {
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			h.Logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("Failed to marshal validation errors")
			errhandler.ServerError(w, errhandler.RespJSONEncodeFailure)
			return
		}
		errhandler.ValidationErrors(w, respBody)
		return
	}

	// Get the adapter for this entity class
	adapter, err := adapters.Get(classID)
	if err != nil {
		h.Logger.Error().Err(err).Str("class_id", classID).Msg("Unknown entity class")
		errhandler.BadRequest(w, []byte(`{"error": "unknown entity class"}`))
		return
	}

	var entityID pgtype.UUID
	if err := entityID.Scan(req.ID); err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Msg("Invalid entity id")
		errhandler.BadRequest(w, []byte(`{"error": "invalid entity id"}`))
		return
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
	entityParams := db.UpdateEntityParams{
		ID:        entityID,
		ParentID:  parentID,
		OKey:      req.Key,
		OPath:     req.Path,
		Published: req.Published,
		HasData:   true,
		UpdatedBy: userID,
	}

	entity, err := h.Queries.UpdateEntity(r.Context(), entityParams)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to update entity record")
		errhandler.ServerError(w, errhandler.RespDBDataInsertFailure)
		return
	}

	// Create the entity data using the adapter
	result, err := adapter.Update(r.Context(), entity.ID, req.Data)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to update entity data")
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
