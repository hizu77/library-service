package book

import (
	"context"

	"github.com/hizu77/library-service/internal/model/db/mapper"

	"github.com/hizu77/library-service/internal/entity"
)

func (b *RepositoryImpl) AddBook(_ context.Context, book entity.Book) (entity.Book, error) {
	b.mx.Lock()
	defer b.mx.Unlock()

	if _, ok := b.books[book.ID]; ok {
		return entity.Book{}, entity.ErrBookAlreadyExists
	}

	dbBook := mapper.BookToDB(book)
	b.books[book.ID] = &dbBook

	for i := range book.AuthorsIDs {
		if _, ok := b.authorBooks[book.AuthorsIDs[i]]; !ok {
			b.authorBooks[book.AuthorsIDs[i]] = make(map[string]struct{})
		}

		b.authorBooks[book.AuthorsIDs[i]][book.ID] = struct{}{}
	}

	return book, nil
}
