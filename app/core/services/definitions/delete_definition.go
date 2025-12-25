package definitions

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/oriiyx/fritz/app/core/api/common/errhandler"
	l "github.com/oriiyx/fritz/app/core/api/common/log"
	"github.com/oriiyx/fritz/app/core/services/objects/entity_builder"
	ctxUtil "github.com/oriiyx/fritz/app/core/utils/ctx"
)

// Delete is an endpoint that handles deleting existing fritz entity
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())
	ID := chi.URLParam(r, EntityIDKey)

	/**
	TODO
		1. Check if definition with the id exists
		2. Create SQL statements that will delete the entire table from the database
		3. Find and delete database/fritz/queries_*.sql
		4. Find and delete database/schema/entity_*.sql
		5. Find and delete the var/entities/definitions/entity_*.json
	*/

	// 1. get existing definition
	definition, err := h.entityBuilder.LoadDefinitionByID(ID)
	if err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("definition_id", ID).Msg("Failed to load existing definition for update")
		errhandler.BadRequest(w, errhandler.RespDBDataAccessFailure)
		return
	}

	// 2. Create SQL statements that will delete the entire table from the database
	tablename := h.entityBuilder.CreateEntityTableName(definition)
	dropSQL := fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", pgx.Identifier{tablename}.Sanitize())
	_, err = h.DB.Exec(r.Context(), dropSQL)
	if err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("definition_id", ID).Msg("Failed to drop table from the database")
		errhandler.ServerError(w, errhandler.RespDBDataRemoveFailure)
		return
	}

	// 3. Find and delete database/fritz/queries_*.sql
	queriesName := fmt.Sprintf("queries_%s.sql", definition.ID)
	err = h.CustomWriter.DeleteFile(entity_builder.EntitiesTableQueriesFilePathTemplate, queriesName)
	if err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("definition_id", ID).Msgf("Failed to delete queries from %s", entity_builder.EntitiesTableQueriesFilePathTemplate)
		errhandler.ServerError(w, errhandler.RespDBDataRemoveFailure)
		return
	}

	// 4. Find and delete database/schema/entity_*.sql
	err = h.CustomWriter.DeleteFile(entity_builder.EntitiesTableSchemaFilePathTemplate, fmt.Sprintf("%s.sql", tablename))
	if err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("definition_id", ID).Msgf("Failed to delete schema from %s", entity_builder.EntitiesTableSchemaFilePathTemplate)
		errhandler.ServerError(w, errhandler.RespDBDataRemoveFailure)
		return
	}

	// 5. Find and delete the var/entities/definitions/entity_*.json
	filename := h.entityBuilder.CreateEntityDefinitionFileName(definition)
	err = h.CustomWriter.DeleteFile(entity_builder.EntitiesDefinitionsFilePathTemplate, filename)
	if err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("definition_id", ID).Msgf("Failed to delete json from %s", entity_builder.EntitiesDefinitionsFilePathTemplate)
		errhandler.ServerError(w, errhandler.RespDBDataRemoveFailure)
		return
	}

	w.WriteHeader(http.StatusOK)
}
