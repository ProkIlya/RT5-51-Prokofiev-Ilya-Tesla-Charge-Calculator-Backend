package handler

import "tesla-app/internal/app/ds"

// GetCurrentUser возвращает текущего пользователя (singleton)
// Должен совпадать с одним из пользователей в БД из fill.sql
func GetCurrentUser() *ds.User {
	// Возвращаем пользователя с ID=1, который есть в fill.sql
	return &ds.User{
		ID:          1,
		Login:       "user",
		Password:    "user",
		IsModerator: false,
	}
}

// Функция для получения модератора (для тестирования)
func GetModeratorUser() *ds.User {
	return &ds.User{
		ID:          2,
		Login:       "admin", // Модератор из БД
		Password:    "admin",
		IsModerator: true,
	}
}
