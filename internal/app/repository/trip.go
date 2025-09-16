package repository

import (
	"tesla-app/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) GetUserDraft(userID uint) (*ds.TripApplication, error) {
	var trip ds.TripApplication
	err := r.db.Where("creator_id = ? AND status = ?", userID, "черновик").First(&trip).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &trip, nil
}

func (r *Repository) CreateDraft(userID uint) (*ds.TripApplication, error) {
	trip := ds.TripApplication{
		Status:      "черновик",
		CreatorID:   userID,
		StartCharge: 100, // Значение по умолчанию
	}
	err := r.db.Create(&trip).Error
	return &trip, err
}

func (r *Repository) AddScenarioToTrip(tripID, scenarioID uint, value float64) error {
	// Проверяем, есть ли уже такой сценарий в заявке
	var count int64
	r.db.Model(&ds.TripScenario{}).
		Where("trip_application_id = ? AND driving_scenario_id = ?", tripID, scenarioID).
		Count(&count)

	if count > 0 {
		// Если сценарий уже есть, обновляем значение
		return r.db.Model(&ds.TripScenario{}).
			Where("trip_application_id = ? AND driving_scenario_id = ?", tripID, scenarioID).
			Update("value", value).Error
	}

	// Если сценария нет, создаем новый
	tripScenario := ds.TripScenario{
		TripApplicationID: tripID,
		DrivingScenarioID: scenarioID,
		Value:             value,
	}
	return r.db.Create(&tripScenario).Error
}

func (r *Repository) GetTripByID(id uint) (*ds.TripApplication, error) {
	var trip ds.TripApplication
	err := r.db.Where("id = ?", id).First(&trip).Error
	if err != nil {
		return nil, err
	}
	return &trip, nil
}

func (r *Repository) GetTripScenarios(tripID uint) ([]ds.TripScenario, error) {
	var tripScenarios []ds.TripScenario
	err := r.db.
		Preload("DrivingScenario"). // Добавляем предзагрузку связанных данных
		Where("trip_application_id = ?", tripID).
		Find(&tripScenarios).Error
	if err != nil {
		return nil, err
	}
	return tripScenarios, nil
}

func (r *Repository) DeleteTrip(tripID uint) error {
	return r.db.Exec("UPDATE trip_applications SET status = 'удалён' WHERE id = ?", tripID).Error
}

func (r *Repository) GetTripScenariosCount(tripID uint) int64 {
	var count int64
	r.db.Model(&ds.TripScenario{}).Where("trip_application_id = ?", tripID).Count(&count)
	return count
}
