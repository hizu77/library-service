package response

import (
	generated "github.com/hizu77/library-service/generated/api/library"
	"github.com/hizu77/library-service/internal/entity"
)

func NewChangeAuthorInfo(_ *entity.Author) *generated.ChangeAuthorInfoResponse {
	return &generated.ChangeAuthorInfoResponse{}
}
