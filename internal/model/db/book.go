package dbmodel

type Book struct {
	ID         string
	Name       string
	AuthorsIDs []string
}
