package repositories

//go:generate mockgen --build_flags=--mod=mod -package repositories -destination=database_mock.go . Database

type Database interface {
	Create(value interface{}) error
	Read(value interface{}, sort, search, identifier string) error
	Update(value interface{}) error
	Delete(value interface{}) error
	ReplaceAssociatons(model interface{}, association string, values interface{}) error
}

// Repositories contains all the repo structs
type Repositories struct {
	BoardgameRepository *BoardgameRepository
	TagRepository       *TagRepository
	CategoryRepository  *CategoryRepository
	MechanismRepository *MechanismRepository
}

// InitRepositories should be called in main.go
func InitRepositories(db Database) *Repositories {
	boardgameRepository := NewBoardgameRepository(db)
	tagRepository := NewTagRepository(db)
	categoryRepository := NewCategoryRepository(db)
	mechanismRepository := NewMechanismRepository(db)

	return &Repositories{
		BoardgameRepository: boardgameRepository,
		TagRepository:       tagRepository,
		CategoryRepository:  categoryRepository,
		MechanismRepository: mechanismRepository,
	}
}
