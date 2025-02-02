package models

import "github.com/google/uuid"

type User struct{
	Id uuid.UUID	`gorm:"type:uuid;default:uuid_generate_v4();" json:"id"`
	Nickname string `json:"nickname" gorm:"not null;unique"`
	Login string 	`json:"login" gorm:"not null;unique"`
	Email string 	`json:"email" gorm:"not null"`
	Password []byte	`json:"-" gorm:"not null"`
}