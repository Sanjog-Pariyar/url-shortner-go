package repository

import (
	"errors"
	"time"

	"github.com/sanjog-pariyar/url-shorten-api/internal/models"

	"gorm.io/gorm"
)

var ErrURLNotFound = errors.New("URL not found")

type URLRepository interface {
	Create(url *models.URL) error
	FindByShortCode(code string) (*models.URL, error)
	IncrementClicks(code string) error
	Delete(code string) error
	FindByUserID(userID uint64) ([]models.URL, error)
}

type urlRepository struct {
	db *gorm.DB
}

func NewURLRepository(db *gorm.DB) URLRepository {
	return &urlRepository{db: db}
}

func (r *urlRepository) Create(url *models.URL) error {
	return r.db.Create(url).Error
}

func (r *urlRepository) FindByShortCode(code string) (*models.URL, error) {
	var url models.URL
	err := r.db.Where("short_code = ?", code).First(&url).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrURLNotFound
		}
		return nil, err
	}
	return &url, nil
}

func (r *urlRepository) IncrementClicks(code string) error {
	return r.db.Model(&models.URL{}).
		Where("short_code = ?", code).
		UpdateColumn("clicks", gorm.Expr("clicks + ?", 1)).
		Error
}

func (r *urlRepository) Delete(code string) error {
	result := r.db.Where("short_code = ?", code).Delete(&models.URL{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrURLNotFound
	}
	return nil
}

func (r *urlRepository) FindByUserID(userID uint64) ([]models.URL, error) {
	var urls []models.URL
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&urls).Error
	return urls, err
}

func (r *urlRepository) IsShortCodeExists(code string) (bool, error) {
	var count int64
	err := r.db.Model(&models.URL{}).Where("short_code = ?", code).Count(&count).Error
	return count > 0, err
}

func (r *urlRepository) DeleteExpired() error {
	return r.db.Where("expires_at IS NOT NULL AND expires_at < ?", time.Now()).Delete(&models.URL{}).Error
}
