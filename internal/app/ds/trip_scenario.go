package ds

type TripScenario struct {
	TripApplicationID uint     `gorm:"primaryKey"`
	DrivingScenarioID uint     `gorm:"primaryKey"`
	Duration          *float64 `gorm:"default:null"` // Изменено на указатель и nullable

	TripApplication TripApplication `gorm:"foreignKey:TripApplicationID"`
	DrivingScenario DrivingScenario `gorm:"foreignKey:DrivingScenarioID"`
}
