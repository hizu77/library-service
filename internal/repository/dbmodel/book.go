package dbmodel

type DBBook struct {
	ID         string
	Name       string
	AuthorsIDs []string
}
