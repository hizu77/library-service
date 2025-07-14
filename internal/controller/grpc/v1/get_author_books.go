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
	if err := validateGetAuthorBooksRequest(request); err != nil {
		c.zap.Error("validate get author books request", zap.Error(err))
		return err
	}

	books, err := c.bookUseCase.GetAuthorBooks(g.Context(), request.GetAuthorId())

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
