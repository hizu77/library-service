package book

import (
	"github.com/hizu77/library-service/internal/model/db"
	"sync"

	"github.com/hizu77/library-service/internal/repository"
)

var _ repository.BookRepository = (*RepositoryImpl)(nil)

type RepositoryImpl struct {
	mx          *sync.RWMutex
	books       map[string]*dbmodel.Book
	authorBooks map[string]map[string]struct{}
}

func NewInMemoryRepository() *RepositoryImpl {
	return &RepositoryImpl{
		mx:          &sync.RWMutex{},
		books:       make(map[string]*dbmodel.Book),
		authorBooks: make(map[string]map[string]struct{}),
	}
}
