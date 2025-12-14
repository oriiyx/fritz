package entities

import (
	"encoding/json"
	"net/http"

	"github.com/oriiyx/fritz/app/core/api/common/errhandler"
	l "github.com/oriiyx/fritz/app/core/api/common/log"
	"github.com/oriiyx/fritz/app/core/services/objects/definitions"
	ctxUtil "github.com/oriiyx/fritz/app/core/utils/ctx"
	validatorUtil "github.com/oriiyx/fritz/app/core/utils/validator"
)

// ReadEntity is an endpoint that handles reading entity
func (h *Handler) ReadEntity(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())
	// classID := chi.URLParam(r, ClassIDKey)

	var req definitions.EntityDefinition
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

	validation, err := h.entityBuilder.ValidateNewDefinition(&req)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Validation of definitions at create entrypoint failed")
		errhandler.BadRequest(w, errhandler.RespFailedToValidateDefinitions)
		return
	}

	if validation != nil {
		errhandler.BadRequest(w, validation)
		return
	}

	w.WriteHeader(http.StatusOK)
}
