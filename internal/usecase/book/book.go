package book

import (
	"context"

	"github.com/google/uuid"
	"github.com/hizu77/library-service/internal/entity"
	"go.uber.org/zap"
)

func (u *UseCaseImpl) GetAuthorBooks(ctx context.Context, id string) ([]entity.Book, error) {
	if _, err := u.authorRepository.GetAuthor(ctx, id); err != nil {
		u.logger.Error("authorRepository.GetAuthor", zap.String("id", id))
		return nil, entity.ErrAuthorNotFound
	}

	books, err := u.bookRepository.GetBooksByAuthorID(ctx, id)
	if err != nil {
		return nil, err
	}

	u.logger.Info("GetAuthorBooks", zap.String("ID", id))

	return books, nil
}

func (u *UseCaseImpl) GetBookInfo(ctx context.Context, id string) (entity.Book, error) {
	v, err := u.bookRepository.GetBook(ctx, id)

	if err != nil {
		u.logger.Error("bookRepository.GetBook", zap.Error(err))
		return entity.Book{}, err
	}

	u.logger.Info("GetBookInfo", zap.String("ID", id))

	return v, nil
}

func (u *UseCaseImpl) AddBook(ctx context.Context, book entity.Book) (entity.Book, error) {
	for i := range book.AuthorsIDs {
		_, err := u.authorRepository.GetAuthor(ctx, book.AuthorsIDs[i])

		if err != nil {
			u.logger.Error("authorRepository.GetAuthor", zap.Error(err))
			return entity.Book{}, err
		}
	}

	book.ID = uuid.New().String()

	v, err := u.bookRepository.AddBook(ctx, book)

	if err != nil {
		u.logger.Error("bookRepository.AddBook", zap.Error(err))
		return entity.Book{}, err
	}

	u.logger.Info("AddBook", zap.String("ID", v.ID))

	return v, nil
}

func (u *UseCaseImpl) UpdateBook(ctx context.Context, book entity.Book) (entity.Book, error) {
	for i := range book.AuthorsIDs {
		_, err := u.authorRepository.GetAuthor(ctx, book.AuthorsIDs[i])
		if err != nil {
			u.logger.Error("authorRepository.GetAuthor", zap.Error(err))
			return entity.Book{}, err
		}
	}

	v, err := u.bookRepository.UpdateBook(ctx, book)

	if err != nil {
		u.logger.Error("bookRepository.UpdateBook", zap.Error(err))
		return entity.Book{}, err
	}

	u.logger.Info("UpdateBook", zap.String("ID", v.ID))

	return v, nil
}
