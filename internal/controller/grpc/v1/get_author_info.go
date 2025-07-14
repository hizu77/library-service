package v1

import (
	"context"

	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/controller/grpc/v1/response"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *ControllerImpl) GetAuthorInfo(ctx context.Context, request *generated.GetAuthorInfoRequest) (*generated.GetAuthorInfoResponse, error) {
	if err := validateGetAuthorInfoRequest(request); err != nil {
		c.zap.Error("validate get author info request", zap.Error(err))
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
