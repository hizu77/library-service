package repository

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
)

//go:generate ../../bin/mockgen -source=contracts.go -destination=../usecase/mocks_repo_test.go -package=usecase_test

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
