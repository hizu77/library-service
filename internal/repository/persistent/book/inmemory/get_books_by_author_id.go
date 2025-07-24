package book

import (
	"context"
	"github.com/hizu77/library-service/internal/model/db/mapper"

	"github.com/hizu77/library-service/internal/entity"
	"github.com/hizu77/library-service/pkg/utils"
)

func (b *RepositoryImpl) GetBooksByAuthorID(_ context.Context, authorID string) ([]entity.Book, error) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	if books, ok := b.authorBooks[authorID]; ok {
		booksIDs := utils.MapToKeySlice(books)
		bookEntities := make([]entity.Book, 0, len(booksIDs))

		for i := range booksIDs {
			bookEntities = append(bookEntities, mapper.BookToDomain(*b.books[booksIDs[i]]))
		}

		return bookEntities, nil
	}

	return []entity.Book{}, nil
}
