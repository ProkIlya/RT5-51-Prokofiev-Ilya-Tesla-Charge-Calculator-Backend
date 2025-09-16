package ds

import "time"

type TripApplication struct {
	ID              uint    `gorm:"primaryKey"`
	Status          string  `gorm:"default:'черновик';check:status IN ('черновик', 'удалён', 'сформирован', 'завершён', 'отклонён')"`
	CreatorID       uint    `gorm:"not null"`
	StartCharge     float64 `gorm:"not null;default:0"`
	RemainingCharge float64
	CreatedAt       time.Time `gorm:"not null;default:current_timestamp"`
	SubmittedAt     time.Time
	CompletedAt     time.Time
	ModeratorID     uint `gorm:"default:null"`

	Creator   User `gorm:"foreignKey:CreatorID"`
	Moderator User `gorm:"foreignKey:ModeratorID"`
	//TripScenarios []TripScenario `gorm:"foreignKey:TripApplicationID"` // Добавленная связь
}
