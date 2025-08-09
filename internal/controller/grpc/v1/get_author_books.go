package v1

import (
	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/controller/grpc/v1/response"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *ControllerImpl) GetAuthorBooks(request *generated.GetAuthorBooksRequest, g grpc.ServerStreamingServer[generated.Book]) error {
	ctx, span := tracer.Start(g.Context(), "HandleGetAuthorBooks")
	defer span.End()

	if err := validateGetAuthorBooksRequest(request); err != nil {
		span.RecordError(err)

		c.zap.Error(
			"Validate get author books request",
			zap.Error(err),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
		)

		return err
	}

	books, err := c.bookUseCase.GetAuthorBooks(ctx, request.GetAuthorId())
	if err != nil {
		return c.convertErr(err)
	}

	for _, book := range books {
		if err = g.Send(response.NewGetAuthorBooks(&book)); err != nil {
			return c.convertErr(err)
		}
	}

	return nil
}

func validateGetAuthorBooksRequest(request *generated.GetAuthorBooksRequest) error {
	if err := request.ValidateAll(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}
