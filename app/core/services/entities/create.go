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

// CreateEntityRequest - metadata only (no data field)
type CreateEntityRequest struct {
	ParentID  *string `json:"parent_id,omitempty"`
	Key       string  `json:"key" validate:"required,max=255"`
	Path      string  `json:"path" validate:"required"`
	Type      string  `json:"type,omitempty"`
	Published bool    `json:"published"`
}

// SaveEntityRequest - for saving actual entity data
type SaveEntityRequest struct {
	Data map[string]interface{} `json:"data" validate:"required"`
}

// CreateEntity creates a new entity instance (metadata only, no data yet)
func (h *Handler) CreateEntity(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())
	classID := chi.URLParam(r, DefinitionIDKey)

	var req CreateEntityRequest
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

	// Verify adapter exists (validate class ID)
	_, err := adapters.Get(classID)
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

	// Create entity record in entities table ONLY (no data table entry yet)
	entityParams := db.CreateEntityParams{
		EntityClass: classID,
		ParentID:    parentID,
		OKey:        req.Key,
		OPath:       req.Path,
		OType:       entityType,
		Published:   req.Published,
		HasData:     false, // Mark as having no data yet
		CreatedBy:   userID,
		UpdatedBy:   userID,
	}

	entity, err := h.Queries.CreateEntity(r.Context(), entityParams)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to create entity record")
		errhandler.ServerError(w, errhandler.RespDBDataInsertFailure)
		return
	}

	h.Logger.Info().
		Str("entity_id", entity.ID.String()).
		Str("class_id", classID).
		Str("key", req.Key).
		Msg("Entity metadata created successfully (no data yet)")

	response := map[string]interface{}{
		"entity": entity,
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(response)
}

// SaveEntity saves data for an existing entity (can be called multiple times)
func (h *Handler) SaveEntity(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())
	classID := chi.URLParam(r, DefinitionIDKey)
	entityID := chi.URLParam(r, "entity_id")

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

	var entityUUID pgtype.UUID
	if err := entityUUID.Scan(entityID); err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Msg("Invalid entity_id")
		errhandler.BadRequest(w, []byte(`{"error": "invalid entity_id"}`))
		return
	}

	// Verify entity exists and matches class
	entity, err := h.Queries.GetEntityByID(r.Context(), entityUUID)
	if err != nil {
		h.Logger.Error().Err(err).Str("entity_id", entityID).Msg("Entity not found")
		errhandler.BadRequest(w, []byte(`{"error": "entity not found"}`))
		return
	}

	if entity.EntityClass != classID {
		h.Logger.Error().
			Str("expected_class", classID).
			Str("actual_class", entity.EntityClass).
			Msg("Entity class mismatch")
		errhandler.BadRequest(w, []byte(`{"error": "entity class mismatch"}`))
		return
	}

	// Check if entity already has data (determines create vs update)
	var result interface{}
	if !entity.HasData {
		// First time saving - CREATE in data table
		result, err = adapter.Create(r.Context(), entity.ID, req.Data)
		if err != nil {
			h.Logger.Error().Err(err).Msg("Failed to create entity data")
			errhandler.ServerError(w, errhandler.RespDBDataInsertFailure)
			return
		}

		// Update has_data flag
		_, err = h.Queries.UpdateEntity(r.Context(), db.UpdateEntityParams{
			ID:        entity.ID,
			ParentID:  entity.ParentID,
			OKey:      entity.OKey,
			OPath:     entity.OPath,
			Published: entity.Published,
			HasData:   true, // Mark as having data now
			UpdatedBy: entity.UpdatedBy,
		})
		if err != nil {
			h.Logger.Error().Err(err).Msg("Failed to update has_data flag")
			// Non-fatal - data was saved, flag update failed
		}

		h.Logger.Info().
			Str("entity_id", entityID).
			Str("class_id", classID).
			Msg("Entity data created (first save)")
	} else {
		// Subsequent saves - UPDATE in data table
		result, err = adapter.Update(r.Context(), entity.ID, req.Data)
		if err != nil {
			h.Logger.Error().Err(err).Msg("Failed to update entity data")
			errhandler.ServerError(w, errhandler.RespDBDataUpdateFailure)
			return
		}

		h.Logger.Info().
			Str("entity_id", entityID).
			Str("class_id", classID).
			Msg("Entity data updated")
	}

	response := map[string]interface{}{
		"entity": entity,
		"data":   result,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
