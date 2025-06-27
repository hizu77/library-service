package book

import (
	"context"

	"github.com/hizu77/library-service/internal/repository/dbmodel/mapper"

	"github.com/hizu77/library-service/internal/entity"
	"github.com/hizu77/library-service/pkg/utils"
)

func (b *RepositoryImpl) GetBook(_ context.Context, id string) (entity.Book, error) {
	b.mx.RLock()
	defer b.mx.RUnlock()

	if book, ok := b.books[id]; ok {
		return mapper.BookToDomain(*book), nil
	}

	return entity.Book{}, entity.ErrBookNotFound
}

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
