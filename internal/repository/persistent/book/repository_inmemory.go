package book

import (
	"sync"

	"github.com/hizu77/library-service/internal/repository/dbmodel"

	"github.com/hizu77/library-service/internal/repository"
)

var _ repository.BookRepository = (*RepositoryImpl)(nil)

type RepositoryImpl struct {
	mx          *sync.RWMutex
	books       map[string]*dbmodel.DBBook
	authorBooks map[string]map[string]struct{}
}

func NewInMemoryRepository() *RepositoryImpl {
	return &RepositoryImpl{
		mx:          &sync.RWMutex{},
		books:       make(map[string]*dbmodel.DBBook),
		authorBooks: make(map[string]map[string]struct{}),
	}
}
