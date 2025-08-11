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

const updateAuthorEndpoint = "UpdateAuthor"

func (c *ControllerImpl) UpdateBook(ctx context.Context, request *generated.UpdateBookRequest) (*generated.UpdateBookResponse, error) {
	start := time.Now()

	EndpointRequests.
		WithLabelValues(updateAuthorEndpoint).
		Inc()

	defer func() {
		EndpointLatency.
			WithLabelValues(updateAuthorEndpoint).
			Observe(float64(time.Since(start).Seconds()))
	}()

	ctx, span := tracer.Start(ctx, "HandleUpdateBook")
	defer span.End()

	if err := validateUpdateBookRequest(request); err != nil {
		span.RecordError(err)

		c.zap.Error(
			"Validate update book request",
			zap.Error(err),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
		)

		return nil, err
	}

	book, err := c.bookUseCase.UpdateBook(ctx, entity.Book{
		ID:         request.GetId(),
		Name:       request.GetName(),
		AuthorsIDs: request.GetAuthorIds(),
	})
	if err != nil {
		return nil, c.convertErr(err)
	}

	return response.NewUpdateBook(&book), nil
}

func validateUpdateBookRequest(request *generated.UpdateBookRequest) error {
	if err := request.ValidateAll(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}
