package v1

import (
	"context"

	"github.com/hizu77/library-service/internal/controller/grpc/v1/response"

	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (c *ControllerImpl) AddBook(ctx context.Context, request *generated.AddBookRequest) (*generated.AddBookResponse, error) {
	if err := validateAddBookRequest(request); err != nil {
		c.zap.Error("validate add book request", zap.Error(err))
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

func (c *ControllerImpl) UpdateBook(ctx context.Context, request *generated.UpdateBookRequest) (*generated.UpdateBookResponse, error) {
	if err := validateUpdateBookRequest(request); err != nil {
		c.zap.Error("validate update book request", zap.Error(err))
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

func (c *ControllerImpl) GetBookInfo(ctx context.Context, request *generated.GetBookInfoRequest) (*generated.GetBookInfoResponse, error) {
	if err := validateGetBookInfoRequest(request); err != nil {
		c.zap.Error("validate get book info request", zap.Error(err))
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
