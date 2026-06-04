package models

type Role struct {
	BaseModel
	Name      string `gorm:"type:string;size:15;not null;unique"`
	RoleUsers *[]RoleUser
}

/*
user -> n roles
role -> n users

n users <-> n roles

users
roles
user_role -> user_id, role_id

*/
