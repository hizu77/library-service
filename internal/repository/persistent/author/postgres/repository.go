package author

import (
	"github.com/hizu77/library-service/internal/repository"
	"github.com/hizu77/library-service/pkg/postgres"
)

var _ repository.AuthorRepository = (*RepositoryImpl)(nil)

type RepositoryImpl struct {
	*postgres.Postgres
}

func New(postgres *postgres.Postgres) *RepositoryImpl {
	return &RepositoryImpl{postgres}
}
