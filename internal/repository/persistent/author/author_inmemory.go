package author

import (
	"context"
	"github.com/hizu77/library-service/internal/repository/dbmodel/mapper"

	"github.com/hizu77/library-service/internal/entity"
)

func (a *RepositoryImpl) GetAuthor(ctx context.Context, id string) (entity.Author, error) {
	a.mx.RLock()
	defer a.mx.RUnlock()

	if author, ok := a.authors[id]; ok {
		return mapper.AuthorToDomain(*author), nil
	} else {
		return entity.Author{}, entity.ErrAuthorNotFound
	}
}

func (a *RepositoryImpl) AddAuthor(ctx context.Context, author entity.Author) (entity.Author, error) {
	a.mx.Lock()
	defer a.mx.Unlock()

	if _, ok := a.authors[author.ID]; ok {
		return entity.Author{}, entity.ErrAuthorAlreadyExists
	}

	dbAuthor := mapper.AuthorToDB(author)
	a.authors[author.ID] = &dbAuthor

	return author, nil
}

func (a *RepositoryImpl) UpdateAuthor(ctx context.Context, author entity.Author) (entity.Author, error) {
	a.mx.Lock()
	defer a.mx.Unlock()

	if _, ok := a.authors[author.ID]; !ok {
		return entity.Author{}, entity.ErrAuthorNotFound
	}

	dbAuthor := mapper.AuthorToDB(author)
	a.authors[author.ID] = &dbAuthor

	return author, nil
}
