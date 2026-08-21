package http

import (
	"context"
	"fmt"
	"go-market/pkg/requestid"
	"log/slog"
	"net/http"
	"strings"

	httpmiddleware "go-market/gateway/internal/transport/http/middleware"
	catalogv1 "go-market/gen/go/catalog/v1"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/protobuf/encoding/protojson"
)

func NewHandler(
	ctx context.Context,
	logger *slog.Logger,
	catalog catalogv1.CatalogServiceClient,
) (http.Handler, error) {
	gatewayMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(
			runtime.MIMEWildcard,
			&runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: false,
				},
			},
		),
		runtime.WithErrorHandler(grpcErrorHandler(logger)),
		runtime.WithIncomingHeaderMatcher(incomingHeaderMatcher),
	)

	if err := catalogv1.RegisterCatalogServiceHandlerClient(
		ctx,
		gatewayMux,
		catalog,
	); err != nil {
		return nil, fmt.Errorf("register catalog gateway handler: %w", err)
	}

	rootMux := http.NewServeMux()
	rootMux.Handle("/", gatewayMux)

	var handler http.Handler

	handler = httpmiddleware.Recovery(logger, rootMux)
	handler = httpmiddleware.Logging(logger, handler)
	handler = httpmiddleware.RequestID(handler)

	return handler, nil
}

func incomingHeaderMatcher(header string) (string, bool) {
	if strings.EqualFold(
		header,
		httpmiddleware.RequestIDHeader,
	) {
		return requestid.GRPCMetadataKey, true
	}

	return runtime.DefaultHeaderMatcher(header)
}
