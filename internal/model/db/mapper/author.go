package mapper

import (
	"github.com/hizu77/library-service/internal/entity"
	"github.com/hizu77/library-service/internal/model/db"
)

func AuthorToDB(author entity.Author) dbmodel.Author {
	return dbmodel.Author{
		ID:   author.ID,
		Name: author.Name,
	}
}

func AuthorToDomain(author dbmodel.Author) entity.Author {
	return entity.Author{
		ID:   author.ID,
		Name: author.Name,
	}
}
