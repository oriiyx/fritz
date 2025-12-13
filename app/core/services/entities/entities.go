package entities

import (
	"encoding/json"
	"net/http"

	"github.com/oriiyx/fritz/app/core/api/base"
	"github.com/oriiyx/fritz/app/core/api/common/errhandler"
	l "github.com/oriiyx/fritz/app/core/api/common/log"
	"github.com/oriiyx/fritz/app/core/services/objects/definitions"
	"github.com/oriiyx/fritz/app/core/services/objects/entity_builder"
	ctxUtil "github.com/oriiyx/fritz/app/core/utils/ctx"
	validatorUtil "github.com/oriiyx/fritz/app/core/utils/validator"
)

type Handler struct {
	*base.HandlerController

	entityBuilder *entity_builder.EntityBuilder
}

func NewEntitiesHandler(ctrl *base.HandlerController) *Handler {
	eb := entity_builder.NewEntityBuilder(ctrl.Logger, ctrl.DB, ctrl.CustomWriter)

	return &Handler{
		HandlerController: ctrl,
		entityBuilder:     eb,
	}
}

// CreateEntity is an endpoint that handles creating of a new fritz entity
func (h *Handler) CreateEntity(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

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

	err = h.entityBuilder.CreateEntityTable(r.Context(), &req)
	if err != nil {
		h.Logger.Error().Err(err).Interface("definition", req).Msg("Failed to create entity table")
		errhandler.ServerError(w, errhandler.RespDBDataInsertFailure)
		return
	}

	err = h.entityBuilder.StoreDefinitionIntoEntityFile(&req)
	if err != nil {
		h.Logger.Error().Err(err).Interface("definition", req).Msg("Failed to store definitions into entity .json file")
		errhandler.ServerError(w, errhandler.RespDBDataInsertFailure)
		return
	}

	w.WriteHeader(http.StatusOK)
}
