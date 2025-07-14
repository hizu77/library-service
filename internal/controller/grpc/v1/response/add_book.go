package response

import (
	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/entity"
)

func NewAddBook(book *entity.Book) *generated.AddBookResponse {
	return &generated.AddBookResponse{
		Book: &generated.Book{
			Id:       book.ID,
			Name:     book.Name,
			AuthorId: book.AuthorsIDs,
		},
	}
}
