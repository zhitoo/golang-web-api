package models

type Country struct {
	BaseModel
	Name   string `gorm:"size:30;type:string;not null;unique"`
	Cities *[]City
}
