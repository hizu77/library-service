package response

import (
	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/entity"
)

func NewGetBookInfo(book *entity.Book) *generated.GetBookInfoResponse {
	return &generated.GetBookInfoResponse{
		Book: &generated.Book{
			Id:       book.ID,
			Name:     book.Name,
			AuthorId: book.AuthorsIDs,
		},
	}
}
