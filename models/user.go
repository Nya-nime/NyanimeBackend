package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        int        `json:"id" gorm:"primaryKey"`
	Username  string     `json:"username" gorm:"not null"`
	Email     string     `json:"email" gorm:"unique;not null"`
	Password  string     `json:"password" gorm:"not null"`
	Role      string     `json:"role" gorm:"not null"`
	Reviews   []Review   `json:"reviews" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Favorites []Favorite `json:"favorites" gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type Review struct {
	ID         int     `json:"id" gorm:"primaryKey"`
	AnimeTitle string  `json:"animeTitle"`
	Rating     float64 `json:"rating"` // Ganti menjadi float64 jika diperlukan
	Content    string  `json:"content"`
	UserID     int     `json:"-" gorm:"constraint:OnDelete:CASCADE"` // Tidak dikembalikan dalam response JSON
}
type Favorite struct {
	ID         int    `json:"id" gorm:"primaryKey"`
	AnimeTitle string `json:"animeTitle"`
	UserID     int    `json:"-"`
}

type Anime struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null"`
	Description string    `json:"description"`
	Genre       string    `json:"genre"`
	ReleaseDate time.Time `json:"releaseDate"`
	CreatedBy   uint      `json:"createdBy"`
}

func GetUserByID(DB *gorm.DB, userID int) (User, error) {
	var user User
	err := DB.Preload("Reviews").Preload("Favorites").First(&user, userID).Error
	return user, err
}
