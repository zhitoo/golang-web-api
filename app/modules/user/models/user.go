package models

type User struct {
	BaseModel
	Username     string `gorm:"type:string;size:20;not null;unique"`
	FirstName    string `gorm:"type:string;size:30;null"`
	LastName     string `gorm:"type:string;size:30;null"`
	MobileNumber string `gorm:"type:string;size:11;null;unique;default:null"`
	Email        string `gorm:"type:string;size:60;null;unique;default:null"`
	Password     string `gorm:"type:string;size:128;not null"`
	Enabled      bool   `gorm:"default:true"`
	UserRoles    *[]RoleUser
}

/*
user -> n roles
role -> n users

n users <-> n roles

users
roles
user_role -> user_id, role_id

*/
