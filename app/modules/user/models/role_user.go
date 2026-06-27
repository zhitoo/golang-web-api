package models

type RoleUser struct {
	Role   Role `gorm:"foreignKey:RoleId;constraint:OnUpdate:NO ACTION;OnDelete:NO ACTION"`
	User   User `gorm:"foreignKey:UserId;constraint:OnUpdate:NO ACTION;OnDelete:NO ACTION"`
	RoleId int
	UserId int
}

func (RoleUser) TableName() string {
	return "role_user"
}

/*
user -> n roles
role -> n users

n users <-> n roles

users
roles
user_role -> user_id, role_id

*/
