package author

import (
	"github.com/hizu77/library-service/internal/repository/dbmodel"
	"sync"

	"github.com/hizu77/library-service/internal/repository"
)

var _ repository.AuthorRepository = (*RepositoryImpl)(nil)

type RepositoryImpl struct {
	mx      *sync.RWMutex
	authors map[string]*dbmodel.DBAuthor
}

func NewInMemoryRepository() *RepositoryImpl {
	return &RepositoryImpl{
		mx:      new(sync.RWMutex),
		authors: make(map[string]*dbmodel.DBAuthor),
	}
}
