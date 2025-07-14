package response

import (
	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/entity"
)

func NewRegisterAuthor(author *entity.Author) *generated.RegisterAuthorResponse {
	return &generated.RegisterAuthorResponse{
		Id: author.ID,
	}
}
