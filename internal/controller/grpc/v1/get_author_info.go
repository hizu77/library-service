package v1

import (
	"context"
	"time"

	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/controller/grpc/v1/response"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const getAuthorInfoEndpoint = "GetAuthorInfo"

func (c *ControllerImpl) GetAuthorInfo(ctx context.Context, request *generated.GetAuthorInfoRequest) (*generated.GetAuthorInfoResponse, error) {
	start := time.Now()

	EndpointRequests.
		WithLabelValues(getAuthorInfoEndpoint).
		Inc()

	defer func() {
		EndpointLatency.
			WithLabelValues(getAuthorInfoEndpoint).
			Observe(float64(time.Since(start).Seconds()))
	}()

	ctx, span := tracer.Start(ctx, "GetAuthorInfo")
	defer span.End()

	if err := validateGetAuthorInfoRequest(request); err != nil {
		span.RecordError(err)

		c.zap.Error(
			"Validate get author info request",
			zap.Error(err),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
		)

		return nil, err
	}

	author, err := c.authorUseCase.GetAuthorInfo(ctx, request.GetId())
	if err != nil {
		return nil, c.convertErr(err)
	}

	return response.NewGetAuthorInfo(&author), nil
}

func validateGetAuthorInfoRequest(request *generated.GetAuthorInfoRequest) error {
	if err := request.ValidateAll(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}
