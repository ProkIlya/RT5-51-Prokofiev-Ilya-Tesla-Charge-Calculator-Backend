package ds

import "time"

type TripApplication struct {
	ID              uint       `gorm:"primaryKey"`
	Status          string     `gorm:"default:'черновик';check:status IN ('черновик', 'удалён', 'сформирован', 'завершён', 'отклонён')"`
	CreatorID       uint       `gorm:"not null"`
	StartCharge     *float64   `gorm:"default:null"` // Изменено на указатель Nullable
	RemainingCharge *float64   `gorm:"default:null"` // Изменено на указатель
	CreatedAt       time.Time  `gorm:"not null;default:current_timestamp"`
	SubmittedAt     *time.Time `gorm:"default:null"` // Изменено на указатель
	CompletedAt     *time.Time `gorm:"default:null"` // Изменено на указатель
	ModeratorID     *uint      `gorm:"default:null"` // Изменено на указатель

	Creator   User `gorm:"foreignKey:CreatorID"`
	Moderator User `gorm:"foreignKey:ModeratorID"`
}
