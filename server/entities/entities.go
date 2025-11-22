package entities

import (
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/oriiyx/fritz/api/common/errhandler"
	l "github.com/oriiyx/fritz/api/common/log"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/oriiyx/fritz/objects/definitions"
	"github.com/oriiyx/fritz/objects/entity_builder"
	ctxUtil "github.com/oriiyx/fritz/utils/ctx"
	validatorUtil "github.com/oriiyx/fritz/utils/validator"
	"github.com/rs/zerolog"
)

type Handler struct {
	logger    *zerolog.Logger
	queries   *db.Queries
	validator *validator.Validate

	entityBuilder *entity_builder.EntityBuilder
}

func NewEntitiesHandler(queries *db.Queries, validator *validator.Validate, logger *zerolog.Logger) *Handler {
	loggerWithService := logger.With().Str("service", "entities").Logger()
	entityBuilder := entity_builder.EntityBuilder{}

	return &Handler{
		logger:        &loggerWithService,
		queries:       queries,
		validator:     validator,
		entityBuilder: &entityBuilder,
	}
}

// CreateEntity is an endpoint that handles creating of a new fritz entity
func (h *Handler) CreateEntity(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	var req definitions.EntityDefinition
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error().Err(err).Str(l.KeyReqID, reqID).Msg("Failed to decode request")
		errhandler.BadRequest(w, errhandler.RespInvalidRequestBody)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			h.logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("Failed to marshal validation errors")
			errhandler.ServerError(w, errhandler.RespJSONEncodeFailure)
			return
		}
		errhandler.ValidationErrors(w, respBody)
		return
	}

	validation, err := h.entityBuilder.ValidateNewDefinition(&req)
	if err != nil {
		h.logger.Error().Err(err).Msg("Validation of definitions at create entrypoint failed")
		errhandler.BadRequest(w, errhandler.RespFailedToValidateDefinitions)
		return
	}

	if validation != nil {
		errhandler.BadRequest(w, validation)
		return
	}

	err = h.entityBuilder.StoreDefinitionIntoEntityFile(&req)
	if err != nil {
		h.logger.Error().Err(err).Interface("definition", req).Msg("Failed to store definitions into entity .json file")
		return
	}

	w.WriteHeader(http.StatusOK)
}
