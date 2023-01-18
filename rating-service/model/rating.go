package model

import "github.com/gofrs/uuid"

type Rating struct {
	CustomBase          `swaggerignore:"true"`
	Username            string `json:"username" db:"username" valid:"required,alphanum,maxstringlength(50)" gorm:"index:unique_rating,unique"`
	Reference_namespace string `json:"reference_namespace" db:"reference_namespace" gorm:"index:unique_rating,unique" valid:"required,alpha,maxstringlength(50)"`
	Reference_id        string `json:"reference_id" db:"reference_id" gorm:"index:unique_rating,unique" valid:"required,uuid"`
	Value               int    `json:"value" db:"value" valid:"required,int,range(0|10)"`
}

func (rating *Rating) GetId() uuid.UUID {
	return rating.ID
}
