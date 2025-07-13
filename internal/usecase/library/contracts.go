package usecase

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
)

type (
	Author interface {
		GetAuthorInfo(ctx context.Context, id string) (*entity.Author, error)
		GetAuthorBooks(ctx context.Context, id string) ([]*entity.Author, error)
		RegisterAuthor(ctx context.Context, author *entity.Author) (*entity.Author, error)
		ChangeAuthorInfo(ctx context.Context, author *entity.Author) (*entity.Author, error)
	}

	Book interface {
		GetBookInfo(ctx context.Context, id string) (*entity.Book, error)
		AddBook(ctx context.Context, book *entity.Book) (*entity.Book, error)
		UpdateBook(ctx context.Context, book *entity.Book) (*entity.Book, error)
	}
)
