package definitions

import (
	"fmt"
	"net/http"
	"os/exec"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
	"github.com/oriiyx/fritz/app/core/api/common/errhandler"
	l "github.com/oriiyx/fritz/app/core/api/common/log"
	"github.com/oriiyx/fritz/app/core/services/objects/definition_builder"
	ctxUtil "github.com/oriiyx/fritz/app/core/utils/ctx"
)

const SQLCGenerateQueriesPath = "database/generated"

// Delete is an endpoint that handles deleting existing fritz entity
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())
	ID := chi.URLParam(r, EntityIDKey)

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

	err = h.Queries.DeleteEntityByClass(r.Context(), definition.ID)
	if err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("definition_id", ID).Msg("Failed to delete entities from the entity table")
		errhandler.ServerError(w, errhandler.RespDBDataRemoveFailure)
		return
	}

	// 3. Find and delete database/fritz/queries_*.sql
	queriesName := fmt.Sprintf("queries_%s.sql", definition.ID)
	err = h.CustomWriter.DeleteFile(definition_builder.EntitiesTableQueriesFilePathTemplate, queriesName)
	if err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("definition_id", ID).Msgf("Failed to delete queries from %s", definition_builder.EntitiesTableQueriesFilePathTemplate)
		errhandler.ServerError(w, errhandler.RespDBDataRemoveFailure)
		return
	}

	// 4.1 Find and delete database/schema/entity_*.sql
	err = h.CustomWriter.DeleteFile(definition_builder.EntitiesTableSchemaFilePathTemplate, fmt.Sprintf("%s.sql", tablename))
	if err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("definition_id", ID).Msgf("Failed to delete schema from %s", definition_builder.EntitiesTableSchemaFilePathTemplate)
		errhandler.ServerError(w, errhandler.RespDBDataRemoveFailure)
		return
	}

	// 4.2 Find and delete database/generated/queries_*.sql.go
	err = h.CustomWriter.DeleteFile(SQLCGenerateQueriesPath, fmt.Sprintf("queries_%s.sql.go", definition.ID))
	if err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("definition_id", ID).Msgf("Failed to delete schema from %s", definition_builder.EntitiesTableSchemaFilePathTemplate)
		errhandler.ServerError(w, errhandler.RespDBDataRemoveFailure)
		return
	}

	// 5. Find and delete the var/entities/definitions/entity_*.json
	filename := h.entityBuilder.CreateEntityDefinitionFileName(definition)
	err = h.CustomWriter.DeleteFile(definition_builder.EntitiesDefinitionsFilePathTemplate, filename)
	if err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("definition_id", ID).Msgf("Failed to delete json from %s", definition_builder.EntitiesDefinitionsFilePathTemplate)
		errhandler.ServerError(w, errhandler.RespDBDataRemoveFailure)
		return
	}

	// 6. Delete adaption code from app/core/services/entities/adapters
	adapterFilename := h.entityBuilder.CreateAdapterFileName(definition)
	err = h.CustomWriter.DeleteFile(definition_builder.EntitiesAdaptersFilePathTemplate, adapterFilename)
	if err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("definition_id", ID).Msgf("Failed to delete adapter code from %s", definition_builder.EntitiesAdaptersFilePathTemplate)
		errhandler.ServerError(w, errhandler.RespDBDataRemoveFailure)
		return
	}

	// 7. Update loader file
	err = h.entityBuilder.UpdateAdapterLoader(definition)
	if err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("definition_id", ID).Msg("Failed to update adapter loader file")
		errhandler.ServerError(w, errhandler.RespDBDataRemoveFailure)
		return
	}

	// 8. Run sqlc generate
	cmd := exec.Command("sqlc", "generate")
	output, err := cmd.CombinedOutput()
	if err != nil {
		h.Logger.Error().
			Err(err).
			Str(l.KeyReqID, reqID).Str("definition_id", ID).
			Str("output", string(output)).
			Msg("SQLC generate failed while trying to delete the definition")
		errhandler.ServerError(w, errhandler.RespDBDataRemoveFailure)
		return
	}

	// 9. Delete entry from the definition_schema table
	err = h.Queries.DeleteDefinitionSchema(r.Context(), definition.ID)
	if err != nil {
		h.Logger.Error().
			Err(err).
			Str(l.KeyReqID, reqID).Str("definition_id", ID).
			Msg("Failed to delete definition schema from the table")
		errhandler.ServerError(w, errhandler.RespDBDataRemoveFailure)
		return
	}

	h.Logger.Info().
		Str("entity_id", definition.ID).
		Msg("SQLC generation completed successfully")

	w.WriteHeader(http.StatusOK)
}
