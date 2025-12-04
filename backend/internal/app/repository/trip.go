package repository

import (
	"tesla-app/internal/app/ds"
	"time"

	"gorm.io/gorm"
)

func (r *Repository) GetUserDraft(userID uint) (*ds.TripApplication, error) {
	var t ds.TripApplication
	err := r.db.Where("creator_id = ? AND status = ?", userID, "черновик").First(&t).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &t, err
}

func (r *Repository) CreateDraft(userID uint) (*ds.TripApplication, error) {
	t := ds.TripApplication{
		CreatorID:   userID,
		Status:      "черновик",
		StartCharge: nil, // или можно установить значение по умолчанию, если нужно
	}
	err := r.db.Create(&t).Error
	return &t, err
}

func (r *Repository) GetTripByIDWithDetails(id uint) (*ds.TripApplication, error) {
	var t ds.TripApplication
	err := r.db.Preload("Creator").Preload("Moderator").Where("id = ?", id).First(&t).Error
	return &t, err
}

func (r *Repository) GetTripsWithFilters(status string, from, to *time.Time, userID *uint) ([]*ds.TripApplication, error) { // изменил
	var ts []*ds.TripApplication
	q := r.db.Preload("Creator").Preload("Moderator").Where("status != ? AND status != ?", "удалён", "черновик")

	if userID != nil {
		q = q.Where("creator_id = ?", *userID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if from != nil {
		q = q.Where("submitted_at >= ?", *from)
	}
	if to != nil {
		q = q.Where("submitted_at <= ?", *to)
	}
	err := q.Find(&ts).Error
	return ts, err
}

func (r *Repository) UpdateTrip(t *ds.TripApplication) error {
	return r.db.Save(t).Error
}

func (r *Repository) DeleteTrip(id uint) error {
	return r.db.Model(&ds.TripApplication{}).Where("id = ?", id).Update("status", "удалён").Error
}

func (r *Repository) GetTripByID(id uint) (*ds.TripApplication, error) {
	var trip ds.TripApplication
	err := r.db.Where("id = ?", id).First(&trip).Error
	if err != nil {
		return nil, err
	}
	return &trip, nil
}
func (r *Repository) CreateTrip(t *ds.TripApplication) error {
	return r.db.Create(t).Error
}
