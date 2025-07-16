package book

import (
	"context"
	"github.com/hizu77/library-service/internal/entity"
)

func (r *RepositoryImpl) GetBooksByAuthorID(ctx context.Context, authorID string) ([]entity.Book, error) {
	//TODO implement me
	panic("implement me")
}
