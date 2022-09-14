package model

type Schema interface {
	GetCreateSchemas() string
	GetDropSchemas() string
}

type SchemaAgregator struct{}

func (SchemaAgregator) GetCreateSchemas() string {
	offerSchema := GetOfferSchema()

	schema := offerSchema

	return schema
}

func (SchemaAgregator) GetDropSchemas() string {

	schema := `
		drop table offer;
		`

	return schema
}
