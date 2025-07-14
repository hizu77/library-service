package author

import (
	"context"

	"github.com/hizu77/library-service/internal/entity"
	"github.com/hizu77/library-service/internal/repository/dbmodel/mapper"
)

func (a *RepositoryImpl) AddAuthor(_ context.Context, author entity.Author) (entity.Author, error) {
	a.mx.Lock()
	defer a.mx.Unlock()

	if _, ok := a.authors[author.ID]; ok {
		return entity.Author{}, entity.ErrAuthorAlreadyExists
	}

	dbAuthor := mapper.AuthorToDB(author)
	a.authors[author.ID] = &dbAuthor

	return author, nil
}
