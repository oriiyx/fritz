package definitions

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oriiyx/fritz/app/core/api/common/errhandler"
	l "github.com/oriiyx/fritz/app/core/api/common/log"
	"github.com/oriiyx/fritz/app/core/services/objects/definitions"
	ctxUtil "github.com/oriiyx/fritz/app/core/utils/ctx"
	validatorUtil "github.com/oriiyx/fritz/app/core/utils/validator"
)

// Update is an endpoint that handles updating existing fritz entity TODO
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())
	ID := chi.URLParam(r, EntityIDKey)

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

	validation, err := h.entityBuilder.ValidateExistingDefinition(&req)
	if err != nil {
		h.Logger.Error().Err(err).Msg("Validation of definitions at create entrypoint failed")
		errhandler.BadRequest(w, errhandler.RespFailedToValidateDefinitions)
		return
	}

	if validation != nil {
		errhandler.BadRequest(w, validation)
		return
	}

	/**
	TODO:
		1. Get existing definition
		2. Check changes from existing definition and new one
		3. Create UPDATE TABLE dynamic query and run it
		4. Update the database/schema/entity_*.sql with a fresh schema
		5. Store definition - reuse the function from the create new definition
		6. Delete existing CRUD Operations
		7. Create fresh CRUD Operations with the existing function
	*/

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
