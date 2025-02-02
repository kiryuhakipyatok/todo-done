package models

import "github.com/google/uuid"

type Tasks struct{
	Task_id uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();" json:"task_id"`
	User_id uuid.UUID `gorm:"type:uuid;" json:"user_id"`
	Time string `json:"time" gorm:"default:'Время не указано'"`
	Task string `json:"task" gorm:"not null;"`
}

