// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package models

import (
	"time"
)

const TableNamePasswordResetToken = "password_reset_tokens"

// PasswordResetToken mapped from table <password_reset_tokens>
type PasswordResetToken struct {
	Email     string    `gorm:"column:email;primaryKey" json:"email"`
	Token     string    `gorm:"column:token;not null" json:"token"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
}

// TableName PasswordResetToken's table name
func (*PasswordResetToken) TableName() string {
	return TableNamePasswordResetToken
}
