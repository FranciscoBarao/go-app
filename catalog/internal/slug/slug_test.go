package slug

import "testing"

func TestFromName(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{"Isaac Childres", "isaac-childres"},
		{"Area Control / Area Influence", "area-control-area-influence"},
		{"  Catan  ", "catan"},
		{"", ""},
	}
	for _, tt := range tests {
		if got := FromName(tt.in); got != tt.want {
			t.Errorf("FromName(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestWithSuffix(t *testing.T) {
	if got := WithSuffix("catan", 2); got != "catan-2" {
		t.Errorf("WithSuffix = %q, want catan-2", got)
	}
}
