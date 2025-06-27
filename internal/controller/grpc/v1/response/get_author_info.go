package response

import (
	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/entity"
)

func NewGetAuthorInfo(author *entity.Author) *generated.GetAuthorInfoResponse {
	return &generated.GetAuthorInfoResponse{
		Id:   author.ID,
		Name: author.Name,
	}
}
