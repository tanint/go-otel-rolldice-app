package exception

import (
	"fmt"
	"net/http"
)

const (
	OK                           = "SUCCESS"
	DOMAIN_EXCEPTION             = "DOMAIN_EXCEPTION"
	NOT_FOUND_EXCEPTION          = "NOT_FOUND"
	UNEXPECTED_EXCEPTION         = "UNHANDLED_EXCEPTION"
	REQUEST_VALIDATION_EXCEPTION = "REQUEST_VALIDATION_EXCEPTION"
	UNAUTHORIZED_REQUEST         = "UNAUTHORIZED"
	REQUIRES_CREDENTIALS         = "MISSING_CREDENTIALS"
	INVALID_CREDENTIALS          = "INVALID_CREDENTIALS"
	INVALID_ACCESS               = "INVALID_ACCESS"
	CONCURRENCY_EXCEPTION        = "CONCURRENCY_EXCEPTION"
	DOWNSTREAM_EXCEPTION         = "DOWNSTREAM_EXCEPTION"
	UNSUPPORTED_CONTENT          = "UNSUPPORTED_CONTENT"
)

func ReportError(err error, message string) {
	if err != nil {
		fmt.Printf("%s : %v", message, err)
	}
}

const (
	DEFAULT_ERROR_MESSAGE = "unknown error captured"
)

type ApplicationError struct {
	Code    string
	Message string
}

func (ae *ApplicationError) Error() string {
	return fmt.Sprintf("error captured [%s] => %s", ae.Code, ae.Message)
}

var (
	NotFoundException           = ApplicationError{Code: NOT_FOUND_EXCEPTION, Message: "requested resource does not exist"}
	UnhandledException          = ApplicationError{Code: UNEXPECTED_EXCEPTION, Message: "unhandled exception"}
	DownstreamException         = ApplicationError{Code: DOWNSTREAM_EXCEPTION, Message: "downstream system error"}
	RequestValidationException  = ApplicationError{Code: REQUEST_VALIDATION_EXCEPTION, Message: "invalid request"}
	UnauthorisedException       = ApplicationError{Code: UNAUTHORIZED_REQUEST, Message: "unauthorised"}
	InvalidAccessException      = ApplicationError{Code: INVALID_ACCESS, Message: "invalid access"}
	ConcurrencyException        = ApplicationError{Code: CONCURRENCY_EXCEPTION, Message: "concurrency exception"}
	MissingCredentialsException = ApplicationError{Code: REQUIRES_CREDENTIALS, Message: "required credentials"}
	InvalidCredentialException  = ApplicationError{Code: INVALID_CREDENTIALS, Message: "invalid credentials"}
	UnsupportedContentException = ApplicationError{Code: UNSUPPORTED_CONTENT, Message: "unsupported content"}
)

var HttpStatusMap = map[string]int{
	NotFoundException.Code:           http.StatusNotFound,
	UnhandledException.Code:          http.StatusInternalServerError,
	DownstreamException.Code:         http.StatusServiceUnavailable,
	RequestValidationException.Code:  http.StatusBadRequest,
	UnauthorisedException.Code:       http.StatusUnauthorized,
	InvalidAccessException.Code:      http.StatusForbidden,
	ConcurrencyException.Code:        http.StatusConflict,
	MissingCredentialsException.Code: http.StatusUnauthorized,
	InvalidCredentialException.Code:  http.StatusUnauthorized,
	UnsupportedContentException.Code: http.StatusUnsupportedMediaType,
}

func GetHttpStatusForCode(code string) int {
	return HttpStatusMap[code]
}
