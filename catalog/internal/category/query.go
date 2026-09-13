package category

import "github.com/FranciscoBarao/catalog/internal/listopt"

// QuerySchema is the sort/filter allowlist for category list endpoints.
// Public names are JSON field names.
var QuerySchema = listopt.Schema{
	"id":         listopt.IntField("id", listopt.Sortable, listopt.Filterable),
	"slug":       listopt.StringField("slug", listopt.Sortable, listopt.Filterable),
	"name":       listopt.StringField("name", listopt.Sortable, listopt.Filterable),
	"created_at": listopt.SortColumn("created_at"),
	"updated_at": listopt.SortColumn("updated_at"),
}
