package repositories

//go:generate mockgen -package mock -destination=mock/repositories.go . Database

type Database interface {
	Create(value interface{}, omits ...string) error
	Read(value interface{}, sort, search, identifier string) error
	Update(value interface{}, omits ...string) error
	Delete(value interface{}) error

	AppendAssociatons(model interface{}, association string, values interface{}) error
	ReplaceAssociatons(model interface{}, association string, values interface{}) error
	//ReadAssociatons(model interface{}, association string, store interface{}) error
	//DeleteAssociatons(model interface{}, association string) error
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
