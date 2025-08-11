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

const getBookInfoEndpoint = "GetBookInfo"

func (c *ControllerImpl) GetBookInfo(ctx context.Context, request *generated.GetBookInfoRequest) (*generated.GetBookInfoResponse, error) {
	start := time.Now()

	EndpointRequests.
		WithLabelValues(getBookInfoEndpoint).
		Inc()

	defer func() {
		EndpointLatency.
			WithLabelValues(getBookInfoEndpoint).
			Observe(float64(time.Since(start).Seconds()))
	}()

	ctx, span := tracer.Start(ctx, "HandleGetBookInfo")
	defer span.End()

	if err := validateGetBookInfoRequest(request); err != nil {
		span.RecordError(err)

		c.zap.Error(
			"Validate get book info request",
			zap.Error(err),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
		)

		return nil, err
	}

	book, err := c.bookUseCase.GetBookInfo(ctx, request.GetId())
	if err != nil {
		return nil, c.convertErr(err)
	}

	return response.NewGetBookInfo(&book), nil
}

func validateGetBookInfoRequest(request *generated.GetBookInfoRequest) error {
	if err := request.ValidateAll(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}
