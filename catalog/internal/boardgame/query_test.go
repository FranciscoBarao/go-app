package boardgame

import (
	"maps"
	"testing"

	"github.com/FranciscoBarao/catalog/internal/listopt"
)

func TestQuerySchema(t *testing.T) {
	want := listopt.Schema{
		"id":             listopt.IntField("id", listopt.Sortable, listopt.Filterable),
		"slug":           listopt.StringField("slug", listopt.Sortable, listopt.Filterable),
		"name":           listopt.StringField("name", listopt.Sortable, listopt.Filterable),
		"year_published": listopt.IntField("year_published", listopt.Sortable, listopt.Filterable),
		"min_players":    listopt.IntField("min_players", listopt.Sortable, listopt.Filterable),
		"max_players":    listopt.IntField("max_players", listopt.Sortable, listopt.Filterable),
		"min_play_time":  listopt.IntField("min_play_time", listopt.Sortable, listopt.Filterable),
		"max_play_time":  listopt.IntField("max_play_time", listopt.Sortable, listopt.Filterable),
		"min_age":        listopt.IntField("min_age", listopt.Sortable, listopt.Filterable),
		"created_at":     listopt.SortColumn("created_at"),
		"updated_at":     listopt.SortColumn("updated_at"),
	}
	if !maps.Equal(QuerySchema, want) {
		t.Fatalf("QuerySchema = %#v\nwant %#v", QuerySchema, want)
	}
	for _, k := range []string{
		"deleted_at", "description", "bgg_id", "boardgame_id",
		"categories", "mechanisms", "contributions", "ratings", "expansions",
		"minplayers", "bio",
	} {
		if _, ok := QuerySchema[k]; ok {
			t.Errorf("omitted field %q is in schema", k)
		}
	}
}
