package auth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/oriiyx/fritz/app/core/api/common/errhandler"
	l "github.com/oriiyx/fritz/app/core/api/common/log"
	ctxUtil "github.com/oriiyx/fritz/app/core/utils/ctx"
	validatorUtil "github.com/oriiyx/fritz/app/core/utils/validator"
	db "github.com/oriiyx/fritz/database/generated"
	"github.com/rs/xid"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (a *API) Login(w http.ResponseWriter, r *http.Request) {
	reqID := ctxUtil.RequestID(r.Context())

	var req NewUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		a.Logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
		errhandler.BadRequest(w, errhandler.RespJSONDecodeFailure)
		return
	}

	if err = a.Validator.Struct(req); err != nil {
		respBody, err := json.Marshal(validatorUtil.ToErrResponse(err))
		if err != nil {
			a.Logger.Error().Str(l.KeyReqID, reqID).Err(err).Msg("")
			errhandler.ServerError(w, errhandler.RespJSONDecodeFailure)
			return
		}

		errhandler.ValidationErrors(w, respBody)
		return
	}

	user, err := a.Queries.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		a.Logger.Info().Str(l.KeyReqID, reqID).Str("email", req.Email).Str("ip", r.RemoteAddr).Msg("Attempted login with existing email - no email exists")
		errhandler.BadRequest(w, errhandler.RespUnauthorized)
		return
	}

	isOkay, err := a.Queries.VerifyPassword(r.Context(), db.VerifyPasswordParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		a.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("email", req.Email).Msg("Could not check if user pass is okay")
		errhandler.ServerError(w, errhandler.RespDBDataInsertFailure)
		return
	}

	if !isOkay {
		a.Logger.Info().Str(l.KeyReqID, reqID).Str("email", req.Email).Str("ip", r.RemoteAddr).Msg("Password is not valid")
		errhandler.BadRequest(w, errhandler.RespUnauthorized)
		return
	}

	deviceID := fmt.Sprintf("web-%s", xid.New().String())
	session, err := a.SessionManager.CreateSession(r.Context(), user.ID, deviceID)
	if err != nil {
		a.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("email", req.Email).Msg("Password is not valid")
		errhandler.BadRequest(w, errhandler.RespUnauthorized)
		return
	}

	w.Header().Set("Access-Control-Allow-Headers", "Set-Cookie")
	w.Header().Set("access-control-expose-headers", "Set-Cookie")
	http.SetCookie(w, &http.Cookie{
		Name:   a.EnvConf.Session.SessionCookieName,
		Value:  session.ID.String(),
		MaxAge: int(a.EnvConf.Session.SessionDuration.Seconds()),
		Path:   "/",
		Secure: a.EnvConf.IsProduction,
	})

	a.Logger.Info().Str(l.KeyReqID, reqID).Str("email", req.Email).Msg("User logged in successfully")
	http.Redirect(w, r, "http://localhost:3333/dashboard", http.StatusFound)
}
