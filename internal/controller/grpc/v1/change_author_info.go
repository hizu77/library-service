package v1

import (
	"context"
	"time"

	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/controller/grpc/v1/response"
	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const changeAuthorInfoEndpoint = "ChangeAuthorInfo"

func (c *ControllerImpl) ChangeAuthorInfo(ctx context.Context, request *generated.ChangeAuthorInfoRequest) (*generated.ChangeAuthorInfoResponse, error) {
	start := time.Now()

	EndpointRequests.
		WithLabelValues(changeAuthorInfoEndpoint).
		Inc()

	defer func() {
		EndpointLatency.
			WithLabelValues(changeAuthorInfoEndpoint).
			Observe(time.Since(start).Seconds())
	}()

	ctx, span := tracer.Start(ctx, "HandleChangeAuthorInfo")
	defer span.End()

	if err := validateChangeAuthorInfoRequest(request); err != nil {
		span.RecordError(err)

		c.zap.Error(
			"Validate change author info request",
			zap.Error(err),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
		)

		return nil, err
	}

	author, err := c.authorUseCase.ChangeAuthorInfo(ctx, entity.Author{
		ID:   request.GetId(),
		Name: request.GetName(),
	})
	if err != nil {
		return nil, c.convertErr(err)
	}

	return response.NewChangeAuthorInfo(&author), nil
}

func validateChangeAuthorInfoRequest(request *generated.ChangeAuthorInfoRequest) error {
	if err := request.ValidateAll(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}
