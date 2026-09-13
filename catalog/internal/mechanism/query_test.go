package mechanism

import (
	"maps"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/listopt"
)

func TestQueryAllowlist(t *testing.T) {
	want := listopt.Allowlist{
		"id":         listopt.IntField("id", listopt.Sortable, listopt.Filterable),
		"slug":       listopt.StringField("slug", listopt.Sortable, listopt.Filterable),
		"name":       listopt.StringField("name", listopt.Sortable, listopt.Filterable),
		"created_at": listopt.SortOnly("created_at"),
		"updated_at": listopt.SortOnly("updated_at"),
	}
	if !maps.Equal(QueryAllowlist, want) {
		t.Fatalf("QueryAllowlist = %#v\nwant %#v", QueryAllowlist, want)
	}
	for _, k := range []string{"deleted_at", "bgg_id", "description", "bio", "boardgame_id"} {
		if _, ok := QueryAllowlist[k]; ok {
			t.Errorf("omitted field %q is in schema", k)
		}
	}
}
