package repository

import (
	"strings"
	"tesla-app/internal/app/ds"
)

func (r *Repository) GetScenarios() ([]*ds.DrivingScenario, error) {
	var ss []*ds.DrivingScenario
	err := r.db.Where("status = ?", "действует").Find(&ss).Error
	return ss, err
}

func (r *Repository) GetScenariosWithFilters(name, typ string) ([]*ds.DrivingScenario, error) {
	var ss []*ds.DrivingScenario
	q := r.db.Where("status = ?", "действует")
	if name != "" {
		q = q.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(name)+"%")
	}
	if typ != "" {
		q = q.Where("type = ?", typ)
	}
	err := q.Find(&ss).Error
	return ss, err
}

func (r *Repository) GetScenarioByID(id uint) (*ds.DrivingScenario, error) {
	var s ds.DrivingScenario
	err := r.db.Where("id = ? AND status = ?", id, "действует").First(&s).Error
	return &s, err
}

func (r *Repository) CreateScenario(s *ds.DrivingScenario) error {
	return r.db.Create(s).Error
}

func (r *Repository) UpdateScenario(s *ds.DrivingScenario) error {
	return r.db.Save(s).Error
}

func (r *Repository) DeleteScenario(id uint) error {
	return r.db.Model(&ds.DrivingScenario{}).Where("id = ?", id).Update("status", "удалён").Error
}
