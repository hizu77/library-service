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

const addBookEndpoint = "AddBook"

func (c *ControllerImpl) AddBook(ctx context.Context, request *generated.AddBookRequest) (*generated.AddBookResponse, error) {
	start := time.Now()

	EndpointRequests.
		WithLabelValues(addBookEndpoint).
		Inc()

	defer func() {
		EndpointLatency.
			WithLabelValues(addBookEndpoint).
			Observe(float64(time.Since(start).Seconds()))
	}()

	ctx, span := tracer.Start(ctx, "HandleAddBook")
	defer span.End()

	if err := validateAddBookRequest(request); err != nil {
		span.RecordError(err)

		c.zap.Error(
			"Validate add book request",
			zap.Error(err),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
		)

		return nil, err
	}

	book, err := c.bookUseCase.AddBook(ctx, entity.Book{
		Name:       request.GetName(),
		AuthorsIDs: request.GetAuthorIds(),
	})
	if err != nil {
		return nil, c.convertErr(err)
	}

	return response.NewAddBook(&book), nil
}

func validateAddBookRequest(request *generated.AddBookRequest) error {
	if err := request.ValidateAll(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}
