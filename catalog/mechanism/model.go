package mechanism

type Mechanism struct {
	Name string `gorm:"primarykey" json:"name" valid:"alphanum, maxstringlength(30)"`
}

func NewMechanism(name string) *Mechanism {
	return &Mechanism{
		Name: name,
	}
}

func (mechanism *Mechanism) UpdateMechanism(name string) {
	if name != "" {
		mechanism.Name = name
	}
}

func (mechanism Mechanism) GetName() string {
	return mechanism.Name
}
