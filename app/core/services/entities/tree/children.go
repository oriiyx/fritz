package tree

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oriiyx/fritz/app/core/api/common/errhandler"
	l "github.com/oriiyx/fritz/app/core/api/common/log"
	ctxUtil "github.com/oriiyx/fritz/app/core/utils/ctx"
	validatorUtil "github.com/oriiyx/fritz/app/core/utils/validator"
	db "github.com/oriiyx/fritz/database/generated"
)

type ChildrenParams struct {
	Limit        uint   `json:"limit" validate:"required,gt=0,lte=1000"`
	Offset       uint   `json:"offset" validate:"required"`
	ParentID     string `json:"parent_id" validate:"required,uuid4"`
	DefinitionID string `json:"definition_id,omitempty"`
}

type GetChildrenResponse struct {
	Items   []db.GetEntityChildrenRow `json:"items"`
	Total   int64                     `json:"total"`
	Limit   int                       `json:"limit"`
	Offset  int                       `json:"offset"`
	HasMore bool                      `json:"has_more"`
}

func (h *Handler) Children(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	var req ChildrenParams
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

	var parentID pgtype.UUID
	if err := parentID.Scan(req.ParentID); err != nil {
		h.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Msg("Invalid entity id")
		errhandler.BadRequest(w, []byte(`{"error": "invalid entity id"}`))
		return
	}

	children, err := h.Queries.GetEntityChildren(r.Context(), db.GetEntityChildrenParams{
		ParentID: parentID,
		Limit:    int32(req.Limit),
		Offset:   int32(req.Offset),
	})
	if err != nil {
		h.Logger.Error().Str(l.KeyReqID, reqID).Err(err).Str("parent_id", req.ParentID).Uint("limit", req.Limit).Uint("offset", req.Offset).Msg("Failed to fetch tree children")
		errhandler.ServerError(w, errhandler.RespDBDataAccessFailure)
		return
	}

	total, err := h.Queries.GetEntityChildrenCount(r.Context(), parentID)
	_ = json.NewEncoder(w).Encode(GetChildrenResponse{
		Items:   children,
		Total:   total,
		Limit:   int(req.Limit),
		Offset:  int(req.Offset),
		HasMore: (int64(req.Offset) + int64(req.Limit)) < total,
	})
}
