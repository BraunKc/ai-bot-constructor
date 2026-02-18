package autherrors

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AppError struct {
	HTTPStatus int
	Message    string
	Err        error
}

func (ae *AppError) Error() string {
	if ae.Err != nil {
		return fmt.Sprintf("msg: %s err: %v", ae.Message, ae.Err)
	}

	return ae.Message
}

func GRPCToHTTPError(err error) error {
	if err == nil {
		return nil
	}

	st, ok := status.FromError(err)
	if !ok {
		return &AppError{
			HTTPStatus: http.StatusInternalServerError,
			Message:    "internal error",
			Err:        err,
		}
	}

	code := st.Code()
	msg := st.Message()

	switch code {
	case codes.AlreadyExists:
		return &AppError{HTTPStatus: http.StatusConflict, Message: msg, Err: err}
	case codes.NotFound:
		return &AppError{HTTPStatus: http.StatusNotFound, Message: msg, Err: err}
	case codes.DataLoss,
		codes.Internal:
		return &AppError{HTTPStatus: http.StatusInternalServerError, Message: "internal error", Err: err}
	case codes.PermissionDenied:
		return &AppError{HTTPStatus: http.StatusForbidden, Message: msg, Err: err}
	case codes.InvalidArgument:
		return &AppError{HTTPStatus: http.StatusBadRequest, Message: msg, Err: err}
	case codes.Unauthenticated:
		return &AppError{HTTPStatus: http.StatusUnauthorized, Message: msg, Err: err}
	default:
		return &AppError{HTTPStatus: http.StatusInternalServerError, Message: "internal error", Err: err}
	}
}
