package definitions

import (
	"encoding/json"
	"net/http"

	"github.com/oriiyx/fritz/app/core/api/common/errhandler"
	l "github.com/oriiyx/fritz/app/core/api/common/log"
	ctxUtil "github.com/oriiyx/fritz/app/core/utils/ctx"
)

// GetExisting get existing definitions in the system
func (h *Handler) GetExisting(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	definitions, err := h.entityBuilder.LoadDefinitionsFromEntityFiles()
	if err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Msg("Failed to get entities")
		errhandler.ServerError(w, errhandler.RespDBDataAccessFailure)
		return
	}

	_ = json.NewEncoder(w).Encode(definitions)
}
