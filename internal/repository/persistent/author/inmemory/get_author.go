package author

import (
	"context"

	"github.com/hizu77/library-service/internal/model/db/mapper"

	"github.com/hizu77/library-service/internal/entity"
)

func (a *RepositoryImpl) GetAuthor(_ context.Context, id string) (entity.Author, error) {
	a.mx.RLock()
	defer a.mx.RUnlock()

	if author, ok := a.authors[id]; ok {
		return mapper.AuthorToDomain(*author), nil
	}

	return entity.Author{}, entity.ErrAuthorNotFound
}
