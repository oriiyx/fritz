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

	// 1. get existing definition
	existingDefinition, err := h.entityBuilder.LoadDefinitionByID(ID)
	if err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("definition_id", ID).Msg("Failed to load existing definition for update")
		errhandler.ServerError(w, errhandler.RespDBDataAccessFailure)
		return
	}

	// 2. compare differences
	changeset, err := h.entityBuilder.CompareDefinitions(existingDefinition, &req)
	if err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("definition_id", ID).Msg("Failed to compare old and new definition for update")
		errhandler.ServerError(w, errhandler.RespDBDataAccessFailure)
		return
	}

	// 3. create UPDATE TABLE dynamic query and run it
	tablename := h.entityBuilder.CreateEntityTableName(existingDefinition)
	err = h.entityBuilder.UpdateTableFromChangeset(changeset, tablename, r.Context())
	if err != nil {
		h.Logger.Error().Err(err).Interface("definition", req).Msg("Failed to create update table queries")
		errhandler.ServerError(w, errhandler.RespDBDataUpdateFailure)
		return
	}

	// 4. Update the database/schema/entity_*.sql with a fresh schema
	err = h.entityBuilder.StoreDefinitionIntoEntityFile(&req)
	if err != nil {
		h.Logger.Error().Err(err).Interface("definition", req).Msg("Failed to store definitions into entity .json file")
		errhandler.ServerError(w, errhandler.RespDBDataInsertFailure)
		return
	}

	// 5. Update the database/schema/entity_*.sql with a fresh schema
	_, err = h.entityBuilder.CreateEntityTable(r.Context(), &req)
	if err != nil {
		h.Logger.Error().Err(err).Interface("definition", req).Msg("Failed to create entity table")
		errhandler.ServerError(w, errhandler.RespDBDataInsertFailure)
		return
	}

	// 6 & 7. Update the database/schema/entity_*.sql with a fresh schema
	err = h.entityBuilder.CreateCrudOperations(tablename, &req)
	if err != nil {
		h.Logger.Error().Err(err).Interface("definition", req).Msg("Failed to create crud operations")
		errhandler.ServerError(w, errhandler.RespDBDataInsertFailure)
		return
	}

	w.WriteHeader(http.StatusOK)
}
