package repository

import (
	"strings"
	"tesla-app/internal/app/ds"
)

func (r *Repository) GetScenarios() ([]*ds.DrivingScenario, error) {
	var scenarios []*ds.DrivingScenario
	err := r.db.Where("status = ?", "действует").Find(&scenarios).Error
	if err != nil {
		return nil, err
	}
	return scenarios, nil
}

func (r *Repository) GetScenarioByID(id uint) (*ds.DrivingScenario, error) {
	var scenario ds.DrivingScenario
	err := r.db.Where("id = ? AND status = ?", id, "действует").First(&scenario).Error
	if err != nil {
		return nil, err
	}
	return &scenario, nil
}

func (r *Repository) SearchScenarios(query string) ([]*ds.DrivingScenario, error) {
	var scenarios []*ds.DrivingScenario
	searchQuery := "%" + strings.ToLower(query) + "%"
	err := r.db.Where("(LOWER(name) LIKE ? OR LOWER(description) LIKE ?) AND status = ?",
		searchQuery, searchQuery, "действует").Find(&scenarios).Error
	if err != nil {
		return nil, err
	}
	return scenarios, nil
}
