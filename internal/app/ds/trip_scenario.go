package ds

type TripScenario struct {
	TripApplicationID uint    `gorm:"primaryKey"`
	DrivingScenarioID uint    `gorm:"primaryKey"`
	Duration          float64 `gorm:"not null"`

	TripApplication TripApplication `gorm:"foreignKey:TripApplicationID"`
	DrivingScenario DrivingScenario `gorm:"foreignKey:DrivingScenarioID"`
}
