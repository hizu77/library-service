package book

import (
	"context"
	"github.com/hizu77/library-service/internal/model/db/mapper"

	"github.com/hizu77/library-service/internal/entity"
)

func (b *RepositoryImpl) UpdateBook(_ context.Context, book entity.Book) (entity.Book, error) {
	b.mx.Lock()
	defer b.mx.Unlock()

	v, ok := b.books[book.ID]
	if !ok {
		return entity.Book{}, entity.ErrBookNotFound
	}

	for i := range v.AuthorsIDs {
		if _, ok := b.authorBooks[v.AuthorsIDs[i]]; !ok {
			return entity.Book{}, entity.ErrAuthorNotFound
		}

		delete(b.authorBooks[v.AuthorsIDs[i]], book.ID)
	}

	newAuthors := book.AuthorsIDs
	for i := range newAuthors {
		if _, ok := b.authorBooks[newAuthors[i]]; !ok {
			b.authorBooks[newAuthors[i]] = make(map[string]struct{})
		}

		b.authorBooks[newAuthors[i]][book.ID] = struct{}{}
	}

	dbBook := mapper.BookToDB(book)
	b.books[book.ID] = &dbBook

	return book, nil
}
