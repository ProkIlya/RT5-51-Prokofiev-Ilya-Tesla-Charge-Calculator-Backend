package ds

type DrivingScenario struct {
	ID               uint    `gorm:"primaryKey"`
	Name             string  `gorm:"not null;size:255"`
	Description      string  `gorm:"type:text"`
	Status           string  `gorm:"default:'действует';check:status IN ('действует', 'удалён')"`
	ImageURL         string  `gorm:"size:500"`
	Type             string  `gorm:"not null;check:type IN ('дорога', 'комфорт')"`
	SystemConsuption float64 `gorm:"default:0"`
	Speed            float64 `gorm:"default:0"`
	AeroCoeff        float64 `gorm:"default:0"`
	RollingCoeff     float64 `gorm:"default:0"`
}
