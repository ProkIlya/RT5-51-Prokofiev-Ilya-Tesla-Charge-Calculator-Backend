package repository

import "tesla-app/internal/app/ds"

func (r *Repository) AddScenarioToTrip(tripID, scID uint, dur float64) error {
	var cnt int64
	r.db.Model(&ds.TripScenario{}).
		Where("trip_application_id = ? AND driving_scenario_id = ?", tripID, scID).
		Count(&cnt)
	if cnt > 0 {
		return r.db.Model(&ds.TripScenario{}).
			Where("trip_application_id = ? AND driving_scenario_id = ?", tripID, scID).
			Update("duration", dur).Error
	}
	ts := ds.TripScenario{TripApplicationID: tripID, DrivingScenarioID: scID, Duration: &dur}
	return r.db.Create(&ts).Error
}

func (r *Repository) GetTripScenarios(tripID uint) ([]ds.TripScenario, error) {
	var list []ds.TripScenario
	err := r.db.Preload("DrivingScenario").Where("trip_application_id = ?", tripID).Find(&list).Error
	return list, err
}

func (r *Repository) GetTripScenariosCount(tripID uint) int64 {
	var cnt int64
	r.db.Model(&ds.TripScenario{}).Where("trip_application_id = ?", tripID).Count(&cnt)
	return cnt
}

func (r *Repository) RemoveScenarioFromTrip(tripID, scID uint) error {
	return r.db.Where("trip_application_id = ? AND driving_scenario_id = ?", tripID, scID).
		Delete(&ds.TripScenario{}).Error
}

func (r *Repository) UpdateTripScenario(tripID, scID uint, dur float64) error {
	return r.db.Model(&ds.TripScenario{}).
		Where("trip_application_id = ? AND driving_scenario_id = ?", tripID, scID).
		Update("duration", dur).Error
}
