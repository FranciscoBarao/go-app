package model

type Schema interface {
	GetCreateSchemas() string
	GetDropSchemas() string
}

type SchemaAgregator string

func (sa SchemaAgregator) GetCreateSchemas() string {
	offerSchema := GetOfferSchema()

	schema := offerSchema

	return schema
}

func (sa SchemaAgregator) GetDropSchemas() string {

	schema := `
		drop table person;
		`

	return schema
}
