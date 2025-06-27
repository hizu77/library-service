package response

import (
	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/entity"
)

func NewGetAuthorBooks(entity *entity.Book) *generated.Book {
	return &generated.Book{
		Id:       entity.ID,
		Name:     entity.Name,
		AuthorId: entity.AuthorsIDs,
	}
}
