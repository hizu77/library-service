package v1

import (
	"context"

	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/controller/grpc/v1/response"
	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *ControllerImpl) RegisterAuthor(ctx context.Context, request *generated.RegisterAuthorRequest) (*generated.RegisterAuthorResponse, error) {
	ctx, span := tracer.Start(ctx, "HandleRegisterAuthor")
	defer span.End()

	if err := validateRegisterAuthorRequest(request); err != nil {
		span.RecordError(err)

		c.zap.Error(
			"Validate register author request",
			zap.Error(err),
			zap.String("trace_id", span.SpanContext().TraceID().String()),
		)

		return nil, err
	}

	author, err := c.authorUseCase.RegisterAuthor(ctx, entity.Author{
		Name: request.GetName(),
	})
	if err != nil {
		return nil, c.convertErr(err)
	}

	return response.NewRegisterAuthor(&author), nil
}

func validateRegisterAuthorRequest(request *generated.RegisterAuthorRequest) error {
	if err := request.ValidateAll(); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}

	return nil
}
