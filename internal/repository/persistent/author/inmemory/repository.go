package author

import (
	"sync"

	dbmodel "github.com/hizu77/library-service/internal/model/db"

	"github.com/hizu77/library-service/internal/repository"
)

var _ repository.AuthorRepository = (*RepositoryImpl)(nil)

type RepositoryImpl struct {
	mx      *sync.RWMutex
	authors map[string]*dbmodel.Author
}

func NewInMemoryRepository() *RepositoryImpl {
	return &RepositoryImpl{
		mx:      new(sync.RWMutex),
		authors: make(map[string]*dbmodel.Author),
	}
}
