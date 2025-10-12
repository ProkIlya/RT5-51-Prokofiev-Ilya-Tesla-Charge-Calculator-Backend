package repository

import (
	"tesla-app/internal/app/ds"

	"gorm.io/gorm"
)

func (r *Repository) CreateUser(u *ds.User) error {
	return r.db.Create(u).Error
}

func (r *Repository) GetUserByLogin(login string) (*ds.User, error) {
	var u ds.User
	err := r.db.Where("login = ?", login).First(&u).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &u, err
}

func (r *Repository) UpdateUser(u *ds.User) error {
	return r.db.Save(u).Error
}

func (r *Repository) GetUserByID(id uint) (*ds.User, error) {
	var u ds.User
	err := r.db.Where("id = ?", id).First(&u).Error
	return &u, err
}
