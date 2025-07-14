package response

import (
	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/entity"
)

func NewUpdateBook(_ *entity.Book) *generated.UpdateBookResponse {
	return &generated.UpdateBookResponse{}
}
