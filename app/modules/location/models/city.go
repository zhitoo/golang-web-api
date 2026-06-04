package models

type City struct {
	BaseModel
	Name      string `gorm:"size:30;type:string;not null;"`
	CountryId int64
	Country   Country `gorm:"foreignKey:CountryId"`
}
