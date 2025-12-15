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
)

type ReadEntityRequest struct {
	ID string `json:"id" validate:"required"`
}

// ReadEntity is an endpoint that handles reading entity
func (h *Handler) ReadEntity(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())
	classID := chi.URLParam(r, ClassIDKey)

	var req ReadEntityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Msg("Failed to decode request")
		errhandler.BadRequest(w, errhandler.RespInvalidRequestBody)
		return
	}

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

	entity, err := h.Queries.GetEntityByID(r.Context(), entityID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to read entity record")
		errhandler.ServerError(w, errhandler.RespDBDataInsertFailure)
		return
	}

	// Create the entity data using the adapter
	result, err := adapter.Read(r.Context(), entity.ID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to read entity data")
		// Rollback: delete the entity record
		// TODO: Consider using transactions
		errhandler.ServerError(w, errhandler.RespDBDataInsertFailure)
		return
	}

	response := map[string]interface{}{
		"entity": entity,
		"data":   result,
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
