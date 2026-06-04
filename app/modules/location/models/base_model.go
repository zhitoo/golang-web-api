package models

import (
	"database/sql"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	Id int `gorm:"primaryKey"`

	CreatedAt  time.Time     `gotm:"type:TIMESTAMP with time zone;not null"`
	ModifiedAt *sql.NullTime `gotm:"type:TIMESTAMP with time zone;null"`
	DeletedAt  *sql.NullTime `gotm:"type:TIMESTAMP with time zone;null"`

	CreatedBy  *sql.NullInt64 `gorm:"null"`
	ModifiedBy *sql.NullInt64 `gorm:"null"`
	DeletedBy  *sql.NullInt64 `gorm:"null"`
}

func (bm *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value("UserId")
	var userId = &sql.NullInt64{Valid: false}
	if value != nil {
		userId.Valid = true
		uid, _ := strconv.Atoi(value.(string))
		userId.Int64 = int64(uid)
	}
	bm.CreatedAt = time.Now().UTC()
	bm.CreatedBy = userId
	return nil
}

func (bm *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value("UserId")
	var userId = &sql.NullInt64{Valid: false}
	if value != nil {
		userId.Valid = true
		uid, _ := strconv.Atoi(value.(string))
		userId.Int64 = int64(uid)
	}
	bm.ModifiedAt = &sql.NullTime{Time: time.Now().UTC(), Valid: true} //
	bm.ModifiedBy = userId

	return nil
}

func (bm *BaseModel) BeforeDelete(tx *gorm.DB) (err error) {
	value := tx.Statement.Context.Value("UserId")
	var userId = &sql.NullInt64{Valid: false}
	if value != nil {
		userId.Valid = true
		uid, _ := strconv.Atoi(value.(string))
		userId.Int64 = int64(uid)
	}
	bm.DeletedAt = &sql.NullTime{Time: time.Now().UTC(), Valid: true} //
	bm.DeletedBy = userId

	return nil
}
