package http

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	httpmiddleware "go-market/services/gateway/internal/transport/http/middleware"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func grpcErrorHandler(logger *slog.Logger) runtime.ErrorHandlerFunc {
	return func(
		_ context.Context,
		_ *runtime.ServeMux,
		_ runtime.Marshaler,
		w http.ResponseWriter,
		r *http.Request,
		err error,
	) {
		grpcStatus := status.Convert(err)
		httpStatus := runtime.HTTPStatusFromCode(grpcStatus.Code())

		var httpStatusError *runtime.HTTPStatusError
		if errors.As(err, &httpStatusError) {
			httpStatus = httpStatusError.HTTPStatus
			grpcStatus = status.Convert(httpStatusError.Err)
		}

		message := publicErrorMessage(grpcStatus)
		httpmiddleware.RecordError(r.Context(), err)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(httpStatus)

		if encodeErr := json.NewEncoder(w).Encode(errorResponse{
			Code:    grpcStatus.Code().String(),
			Message: message,
		}); encodeErr != nil {
			logger.ErrorContext(
				r.Context(),
				"failed to encode gateway error response",
				slog.Any("error", encodeErr),
			)
		}
	}
}

func publicErrorMessage(grpcStatus *status.Status) string {
	switch grpcStatus.Code() {
	case codes.Internal, codes.Unknown, codes.DataLoss:
		return "internal server error"
	case codes.Unavailable:
		return "service unavailable"
	default:
		return grpcStatus.Message()
	}
}
