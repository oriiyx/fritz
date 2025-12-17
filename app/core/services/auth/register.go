package auth

import (
	"encoding/json"
	"net/http"

	"github.com/oriiyx/fritz/app/core/api/common/errhandler"
	l "github.com/oriiyx/fritz/app/core/api/common/log"
	ctxUtil "github.com/oriiyx/fritz/app/core/utils/ctx"
	validatorUtil "github.com/oriiyx/fritz/app/core/utils/validator"
	db "github.com/oriiyx/fritz/database/generated"
)

type NewUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

func (a *API) Register(w http.ResponseWriter, r *http.Request) {
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

	_, err = a.Queries.GetUserByEmail(r.Context(), req.Email)
	if err == nil {
		a.Logger.Warn().Str(l.KeyReqID, reqID).Str("email", req.Email).Str("ip", r.RemoteAddr).Msg("Attempted registration with existing email")
		errhandler.BadRequest(w, errhandler.RespUserEmailAlreadyExists)
		return
	}

	_, err = a.Queries.CreateUser(r.Context(), db.CreateUserParams{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		a.Logger.Error().Err(err).Str(l.KeyReqID, reqID).Str("email", req.Email).Msg("Could not register a new user")
		errhandler.ServerError(w, errhandler.RespDBDataInsertFailure)
		return
	}

	a.Logger.Info().Str(l.KeyReqID, reqID).Str("email", req.Email).Msg("User registered successfully")
}
