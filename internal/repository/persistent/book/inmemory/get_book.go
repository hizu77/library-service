package book

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
	"github.com/hizu77/library-service/internal/repository/dbmodel/mapper"
)

func (b *RepositoryImpl) GetBook(_ context.Context, id string) (entity.Book, error) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	if book, ok := b.books[id]; ok {
		return mapper.BookToDomain(*book), nil
	}

	return entity.Book{}, entity.ErrBookNotFound
}
