package model

type Mechanism struct {
	Name       string      `gorm:"primarykey" json:"name" valid:"alphanum, maxstringlength(30)`
	Boardgames []Boardgame `gorm:"many2many:boardgame_mechanisms;" json:"-"`
}

func NewMechanism(name string) Mechanism {
	return Mechanism{
		Name: name,
	}
}

func (mechanism *Mechanism) UpdateMechanism(name string) {
	if name != "" {
		mechanism.Name = name
	}
}

// Getters
func (mechanism Mechanism) GetName() string {
	return mechanism.Name
}
