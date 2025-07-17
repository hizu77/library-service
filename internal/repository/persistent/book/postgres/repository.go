package book

import (
	"github.com/hizu77/library-service/internal/repository"
	"github.com/hizu77/library-service/pkg/postgres"
)

var _ repository.BookRepository = (*RepositoryImpl)(nil)

type RepositoryImpl struct {
	*postgres.Postgres
}

func New(postgres *postgres.Postgres) *RepositoryImpl {
	return &RepositoryImpl{postgres}
}
