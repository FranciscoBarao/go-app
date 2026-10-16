package contributor

import "testing"

func TestRoleValid(t *testing.T) {
	tests := []struct {
		name  string
		role  Role
		valid bool
	}{
		{name: "designer", role: RoleDesigner, valid: true},
		{name: "artist", role: RoleArtist, valid: true},
		{name: "publisher", role: RolePublisher, valid: true},
		{name: "developer", role: RoleDeveloper, valid: true},
		{name: "graphic designer", role: RoleGraphicDesigner, valid: true},
		{name: "unknown", role: Role("editor"), valid: false},
		{name: "empty", role: Role(""), valid: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.role.Valid(); got != tt.valid {
				t.Fatalf("Role(%q).Valid() = %v, want %v", tt.role, got, tt.valid)
			}
		})
	}
}
