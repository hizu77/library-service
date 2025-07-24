package mapper

import (
	"github.com/hizu77/library-service/internal/entity"
	"github.com/hizu77/library-service/internal/model/db"
)

func BookToDB(book entity.Book) dbmodel.Book {
	cp := make([]string, len(book.AuthorsIDs))
	copy(cp, book.AuthorsIDs)

	return dbmodel.Book{
		ID:         book.ID,
		Name:       book.Name,
		AuthorsIDs: cp,
	}
}

func BookToDomain(book dbmodel.Book) entity.Book {
	cp := make([]string, len(book.AuthorsIDs))
	copy(cp, book.AuthorsIDs)

	return entity.Book{
		ID:         book.ID,
		Name:       book.Name,
		AuthorsIDs: cp,
	}
}
