package entities

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oriiyx/fritz/app/core/api/base"
	"github.com/oriiyx/fritz/app/core/api/common/errhandler"
	l "github.com/oriiyx/fritz/app/core/api/common/log"
	"github.com/oriiyx/fritz/app/core/services/objects/definition_builder"
	ctxUtil "github.com/oriiyx/fritz/app/core/utils/ctx"
	validatorUtil "github.com/oriiyx/fritz/app/core/utils/validator"
)

const (
	DefinitionIDKey = "definition_id"
	EntityIDKey     = "entity_id"
)

type Handler struct {
	*base.HandlerController

	entityBuilder *definition_builder.Builder
}

func New(ctrl *base.HandlerController) *Handler {
	eb := definition_builder.NewDefinitionsBuilder(ctrl.Logger, ctrl.DB, ctrl.CustomWriter)

	return &Handler{
		HandlerController: ctrl,
		entityBuilder:     eb,
	}
}

type GetEntityDataRequest struct {
	ID string `json:"id" validate:"required"`
}

// GetEntityData is an endpoint that handles reading entity
func (h *Handler) GetEntityData(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	var req GetEntityDataRequest
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

	var entityID pgtype.UUID
	if err := entityID.Scan(req.ID); err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Msg("Invalid entity id")
		errhandler.BadRequest(w, []byte(`{"error": "invalid entity id"}`))
		return
	}

	entity, err := h.Queries.GetEntityByID(r.Context(), entityID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to read entity record")
		errhandler.ServerError(w, errhandler.RespDBDataAccessFailure)
		return
	}

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(entity)
}
