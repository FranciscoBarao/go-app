package category

import "github.com/FranciscoBarao/catalog/internal/listopt"

// QueryAllowlist is the sort/filter allowlist for category list endpoints.
// Public names are JSON field names.
var QueryAllowlist = listopt.Allowlist{
	"id":         listopt.IntField("id", listopt.Sortable, listopt.Filterable),
	"slug":       listopt.StringField("slug", listopt.Sortable, listopt.Filterable),
	"name":       listopt.StringField("name", listopt.Sortable, listopt.Filterable),
	"created_at": listopt.SortOnly("created_at"),
	"updated_at": listopt.SortOnly("updated_at"),
}
