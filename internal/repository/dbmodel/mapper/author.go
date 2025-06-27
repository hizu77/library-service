package mapper

import (
	"github.com/hizu77/library-service/internal/entity"
	"github.com/hizu77/library-service/internal/repository/dbmodel"
)

func AuthorToDB(author entity.Author) dbmodel.DBAuthor {
	return dbmodel.DBAuthor{
		ID:   author.ID,
		Name: author.Name,
	}
}

func AuthorToDomain(author dbmodel.DBAuthor) entity.Author {
	return entity.Author{
		ID:   author.ID,
		Name: author.Name,
	}
}
