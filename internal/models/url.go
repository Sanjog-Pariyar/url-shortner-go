package models

import "time"

type URL struct {
	ID          uint64     `gorm:"primaryKey;autoIncrement" json:"id"`
	ShortCode   string     `gorm:"uniqueIndex;size:20;not null" json:"shortCode"`
	OriginalURL string     `gorm:"type:text;not null" json:"originalUrl"`
	UserID      *uint64    `gorm:"index" json:"userId,omitempty"`
	Clicks      int64      `gorm:"default:0" json:"clicks"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
	CreatedAt   time.Time  `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt   time.Time  `gorm:"autoUpdateTime" json:"updatedAt"`
}

type CreateURLRequest struct {
	URL        string `json:"url" binding:"required"`
	CustomCode string `json:"customCode,omitempty" binding:"omitempty,alphanum,min=3,max=20"`
	ExpiresIn  *int   `json:"expiresIn,omitempty"`
}

type CreateURLResponse struct {
	ShortCode string     `json:"shortCode"`
	ShortURL  string     `json:"shortUrl"`
	ExpiresAt *time.Time `json:"expiresAt,omitempty"`
}

type URLStatsResponse struct {
	ShortCode   string     `json:"shortCode"`
	OriginalURL string     `json:"originalUrl"`
	Clicks      int64      `json:"clicks"`
	CreatedAt   time.Time  `json:"createdAt"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
}
