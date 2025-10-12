package ds

type User struct {
	ID          uint   `gorm:"primaryKey"`
	Login       string `gorm:"unique;not null;size:150"`
	Password    string `gorm:"not null;size:128"`
	IsModerator bool   `gorm:"default:false"`
}
