package repository

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
)

//go:generate ../../bin/mockgen -destination=../usecase/mock/author.go -package=mock github.com/hizu77/library-service/internal/repository AuthorRepository
//go:generate ../../bin/mockgen -destination=../usecase/mock/book.go -package=mock github.com/hizu77/library-service/internal/repository BookRepository

type (
	AuthorRepository interface {
		GetAuthor(ctx context.Context, id string) (entity.Author, error)
		AddAuthor(ctx context.Context, author entity.Author) (entity.Author, error)
		UpdateAuthor(ctx context.Context, author entity.Author) (entity.Author, error)
	}

	BookRepository interface {
		GetBook(ctx context.Context, id string) (entity.Book, error)
		GetBooksByAuthorID(ctx context.Context, authorID string) ([]entity.Book, error)
		UpdateBook(ctx context.Context, book entity.Book) (entity.Book, error)
		AddBook(ctx context.Context, book entity.Book) (entity.Book, error)
	}
)
