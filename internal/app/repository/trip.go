package repository

import (
	"tesla-app/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) GetUserDraft(userID uint) (*ds.TripApplication, error) { // Получение черновика заявки пользователя через ORM
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

func (r *Repository) CreateDraft(userID uint) (*ds.TripApplication, error) { // Создание новой заявки-черновика через ORM
	trip := ds.TripApplication{
		Status:      "черновик",
		CreatorID:   userID,
		StartCharge: 100, // Значение по умолчанию
	}
	err := r.db.Create(&trip).Error
	return &trip, err
}

func (r *Repository) AddScenarioToTrip(tripID, scenarioID uint, duration float64) error { // Добавление или обновление сценария в заявке через ORM
	// Проверяем, есть ли уже такой сценарий в заявке
	var count int64
	r.db.Model(&ds.TripScenario{}).
		Where("trip_application_id = ? AND driving_scenario_id = ?", tripID, scenarioID).
		Count(&count)

	if count > 0 {
		// Если сценарий уже есть, обновляем значение
		return nil
	}

	// Если сценария нет, создаем новый
	tripScenario := ds.TripScenario{
		TripApplicationID: tripID,
		DrivingScenarioID: scenarioID,
		Duration:          duration,
	}
	return r.db.Create(&tripScenario).Error
}

func (r *Repository) GetTripByID(id uint) (*ds.TripApplication, error) { // Получение заявки по ID через ORM
	var trip ds.TripApplication
	err := r.db.Where("id = ?", id).First(&trip).Error
	if err != nil {
		return nil, err
	}
	return &trip, nil
}

func (r *Repository) GetTripScenarios(tripID uint) ([]ds.TripScenario, error) { // Получение всех сценариев для конкретной заявки через ORM
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

func (r *Repository) DeleteTrip(tripID uint) error { // Логическое удаление заявки через SQL-запрос
	return r.db.Exec("UPDATE trip_applications SET status = 'удалён' WHERE id = ?", tripID).Error
}

func (r *Repository) GetTripScenariosCount(tripID uint) int64 { // Получение количества сценариев езды в заявке через ORM
	var count int64
	r.db.Model(&ds.TripScenario{}).Where("trip_application_id = ?", tripID).Count(&count)
	return count
}
