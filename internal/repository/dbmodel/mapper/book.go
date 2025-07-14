package mapper

import (
	"github.com/hizu77/library-service/internal/entity"
	"github.com/hizu77/library-service/internal/repository/dbmodel"
)

func BookToDB(book entity.Book) dbmodel.DBBook {
	cp := make([]string, len(book.AuthorsIDs))
	copy(cp, book.AuthorsIDs)

	return dbmodel.DBBook{
		ID:         book.ID,
		Name:       book.Name,
		AuthorsIDs: cp,
	}
}

func BookToDomain(book dbmodel.DBBook) entity.Book {
	cp := make([]string, len(book.AuthorsIDs))
	copy(cp, book.AuthorsIDs)

	return entity.Book{
		ID:         book.ID,
		Name:       book.Name,
		AuthorsIDs: cp,
	}
}
