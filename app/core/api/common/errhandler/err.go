package errhandler

import (
	"net/http"
)

var (
	RespFailedUUIDDecode              = []byte(`{"error": "cannot decode UUID"}`)
	RespDBDataInsertFailure           = []byte(`{"error": "db data insert failure"}`)
	RespDBDataAccessFailure           = []byte(`{"error": "db data access failure"}`)
	RespDBDataUpdateFailure           = []byte(`{"error": "db data update failure"}`)
	RespDBDataRemoveFailure           = []byte(`{"error": "db data remove failure"}`)
	RespDBDataEmptyFieldCstrViolation = []byte(`{"error": "empty field constraint violation"}`)
	RespDBDataUniqueCstrViolation     = []byte(`{"error": "unique constraint violation"}`)

	RespProcessFailure = []byte(`{"error": "internal processing failure"}`)
	RespProcessTimeout = []byte(`{"error": "processing timeout"}`)

	RespJSONEncodeFailure = []byte(`{"error": "json encode failure"}`)
	RespJSONDecodeFailure = []byte(`{"error": "json decode failure"}`)

	RespInvalidURLParamID             = []byte(`{"error": "invalid url param-id"}`)
	RespInvalidURLQueryParamValue     = []byte(`{"error": "invalid url query-param value"}`)
	RespInvalidIssuerIDParamValue     = []byte(`{"error": "invalid issuer_id value"}`)
	RespInvalidCredentialIDParamValue = []byte(`{"error": "invalid id value"}`)
	RespInvalidRequestBody            = []byte(`{"error": "invalid request body"}`)
	RespUnauthorized                  = []byte(`{"error": "unauthorized"}`)
	RespRestrictedAccess              = []byte(`{"error": "restricted access"}`)

	RespUserEmailAlreadyExists = []byte(`{"error": ""user with your email already exists""}`)

	RespMissingRequiredFields = []byte(`{"error": "missing required field: key_id, email, or salt"}`)
	RespDecryptionFailed      = []byte(`{"error": "key decryption failed"}`)

	RespFailedToValidateDefinitions = []byte(`{"error": "failed to validate definitions"}`)
	RespEntityNameAlreadyExists     = []byte(`{"error": "entity name already exists"}`)
	RespEntityIDAlreadyExists       = []byte(`{"error": "entity id already exists"}`)
)

type Error struct {
	Error string `json:"error"`
}

type Errors struct {
	Errors []string `json:"errors"`
}

func ServerError(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusInternalServerError)
	_, _ = w.Write(error)
}

func BadRequest(w http.ResponseWriter, error []byte) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(error)
}

func ValidationErrors(w http.ResponseWriter, resp []byte) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(resp)
}

func Unauthorized(w http.ResponseWriter, resp []byte) {
	w.WriteHeader(http.StatusBadRequest)
	_, _ = w.Write(resp)
}

func Forbidden(w http.ResponseWriter, resp []byte) {
	w.WriteHeader(http.StatusForbidden)
	_, _ = w.Write(resp)
}

func NotAcceptable(w http.ResponseWriter, resp []byte) {
	w.WriteHeader(http.StatusNotAcceptable)
	_, _ = w.Write(resp)
}
