package auth

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oriiyx/fritz/app/core/api/common/errhandler"
	l "github.com/oriiyx/fritz/app/core/api/common/log"
	ctxUtil "github.com/oriiyx/fritz/app/core/utils/ctx"
)

type FrontedUser struct {
	ID          pgtype.UUID `json:"id"`
	Email       string      `json:"email"`
	DisplayName pgtype.Text `json:"display_name"`
}

func (a *API) MeHandler(w http.ResponseWriter, r *http.Request) {
	session := ctxUtil.GetSession(r.Context())
	reqID := ctxUtil.RequestID(r.Context())

	user, err := a.Queries.GetUserByID(r.Context(), session.UserIdentityID)
	if err != nil {
		a.Logger.Error().Str(l.KeyReqID, reqID).Err(err).Interface("user_id", session.UserIdentityID).Msg("Failed to get user by ID")
		errhandler.ServerError(w, errhandler.RespInvalidRequestBody)
		return
	}

	returnUser := FrontedUser{
		ID:          user.ID,
		Email:       user.Email,
		DisplayName: user.FullName,
	}

	if err = json.NewEncoder(w).Encode(returnUser); err != nil {
		a.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Msg("Failed to encode me response")
		errhandler.ServerError(w, errhandler.RespProcessFailure)
		return
	}
}
