package definitions

import (
	"encoding/json"
	"net/http"

	"github.com/oriiyx/fritz/app/core/api/base"
	"github.com/oriiyx/fritz/app/core/api/common/errhandler"
	l "github.com/oriiyx/fritz/app/core/api/common/log"
	"github.com/oriiyx/fritz/app/core/services/objects/definition_builder"
	"github.com/oriiyx/fritz/app/core/services/objects/definitions"
	ctxUtil "github.com/oriiyx/fritz/app/core/utils/ctx"
	validatorUtil "github.com/oriiyx/fritz/app/core/utils/validator"
)

type Handler struct {
	*base.HandlerController

	entityBuilder *definition_builder.Builder
}

const EntityIDKey = "id"

func New(ctrl *base.HandlerController) *Handler {
	eb := definition_builder.NewDefinitionsBuilder(ctrl.Logger, ctrl.DB, ctrl.CustomWriter)

	return &Handler{
		HandlerController: ctrl,
		entityBuilder:     eb,
	}
}

// Create is an endpoint that handles creating of a new fritz entity
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
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

	tablename, err := h.entityBuilder.CreateEntityTable(r.Context(), &req)
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

	err = h.entityBuilder.CreateCrudOperations(tablename, &req)
	if err != nil {
		h.Logger.Error().Err(err).Interface("definition", req).Msg("Failed to create crud operations")
		errhandler.ServerError(w, errhandler.RespDBDataInsertFailure)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetDataComponentTypes returns all available data component types
func (h *Handler) GetDataComponentTypes(w http.ResponseWriter, r *http.Request) {
	_ = json.NewEncoder(w).Encode(definitions.GetAllDataComponents())
}
