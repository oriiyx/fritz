package entities

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oriiyx/fritz/app/core/api/common/errhandler"
	l "github.com/oriiyx/fritz/app/core/api/common/log"
	ctxUtil "github.com/oriiyx/fritz/app/core/utils/ctx"
	validatorUtil "github.com/oriiyx/fritz/app/core/utils/validator"
)

type DeleteEntityRequest struct {
	ID string `json:"id" validate:"required"`
}

// DeleteEntity is an endpoint that handles deleting entity
func (h *Handler) DeleteEntity(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	var req DeleteEntityRequest
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

	err := h.Queries.DeleteEntity(r.Context(), entityID)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Failed to delete entity record")
		errhandler.ServerError(w, errhandler.RespDBDataInsertFailure)
		return
	}

	w.WriteHeader(http.StatusOK)
}
