package mechanism

// Mechanism represents a board game mechanism.
type Mechanism struct {
	Name string `gorm:"primarykey" json:"name" valid:"alphanum, maxstringlength(30)"`
}

// NewMechanism creates a new Mechanism with the given name.
func NewMechanism(name string) *Mechanism {
	return &Mechanism{
		Name: name,
	}
}

// UpdateMechanism updates the mechanism name if non-empty.
func (mechanism *Mechanism) UpdateMechanism(name string) {
	if name != "" {
		mechanism.Name = name
	}
}

// GetName returns the mechanism name.
func (mechanism Mechanism) GetName() string {
	return mechanism.Name
}
