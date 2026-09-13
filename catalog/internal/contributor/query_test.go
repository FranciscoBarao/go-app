package contributor

import (
	"maps"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/listopt"
)

func TestQuerySchema(t *testing.T) {
	want := listopt.Schema{
		"id":         listopt.IntField("id", listopt.Sortable, listopt.Filterable),
		"slug":       listopt.StringField("slug", listopt.Sortable, listopt.Filterable),
		"name":       listopt.StringField("name", listopt.Sortable, listopt.Filterable),
		"created_at": listopt.SortColumn("created_at"),
		"updated_at": listopt.SortColumn("updated_at"),
	}
	if !maps.Equal(QuerySchema, want) {
		t.Fatalf("QuerySchema = %#v\nwant %#v", QuerySchema, want)
	}
	for _, k := range []string{"deleted_at", "bio", "bgg_id", "description", "boardgame_id"} {
		if _, ok := QuerySchema[k]; ok {
			t.Errorf("omitted field %q is in schema", k)
		}
	}
}
