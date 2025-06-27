package v1

import (
	"context"

	"github.com/hizu77/library-service/internal/controller/grpc/v1/response"

	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *ControllerImpl) RegisterAuthor(ctx context.Context, request *generated.RegisterAuthorRequest) (*generated.RegisterAuthorResponse, error) {
	if err := validateRegisterAuthorRequest(request); err != nil {
		c.zap.Error("validate register author request", zap.Error(err))
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

func (c *ControllerImpl) ChangeAuthorInfo(ctx context.Context, request *generated.ChangeAuthorInfoRequest) (*generated.ChangeAuthorInfoResponse, error) {
	if err := validateChangeAuthorInfoRequest(request); err != nil {
		c.zap.Error("validate change author info request", zap.Error(err))
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
