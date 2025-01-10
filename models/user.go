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
	ID        int       `json:"id" gorm:"primaryKey;autoIncrement"`
	Rating    int64     `json:"rating" gorm:"type:bigint"`
	Content   string    `json:"content" gorm:"type:longtext"`
	CreatedAt time.Time `json:"created_at" gorm:"type:datetime(3);autoCreateTime"`
	UserID    int       `json:"user_id" gorm:"column:user_id;not null"`   // Pastikan ini int
	AnimeID   uint      `json:"anime_id" gorm:"column:anime_id;not null"` // Pastikan ini uint
}

func (Review) TableName() string {
	return "reviews_new" // Menggunakan nama tabel yang sudah ada
}

type Favorite struct {
	ID         uint64    `json:"id" gorm:"primaryKey"` // ID favorit
	AnimeID    uint64    `json:"anime_id"`             // ID anime yang difavoritkan
	AnimeTitle string    `json:"animeTitle"`           // Judul anime
	UserID     int       `json:"-"`                    // ID pengguna (tidak dikembalikan dalam respons)
	CreatedAt  time.Time `json:"created_at"`           // Waktu pembuatan
}

func (Favorite) TableName() string {
	return "favorite" // Menggunakan nama tabel yang sudah ada
}

type Anime struct {
	ID            uint    `json:"id" gorm:"primaryKey"`
	Title         string  `json:"title" gorm:"not null"`
	Description   string  `json:"description"`
	Genre         string  `json:"genre"`
	ReleaseDate   string  `json:"releaseDate"`
	CreatedBy     uint    `json:"createdBy"`
	AverageRating float64 `json:"average_rating"`
}

func GetUserByID(DB *gorm.DB, userID int) (User, error) {
	var user User
	err := DB.Preload("Reviews").Preload("Favorites").First(&user, userID).Error
	return user, err
}
