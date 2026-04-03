package services

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/sanjog-pariyar/url-shorten-api/internal/models"
	"github.com/sanjog-pariyar/url-shorten-api/internal/repository"
)

var (
	ErrInvalidURL      = errors.New("invalid URL")
	ErrShortCodeExists = errors.New("short code already exists")
	ErrURLExpired      = errors.New("URL has expired")
)

type URLService interface {
	CreateShortURL(req *models.CreateURLRequest, userID *uint64) (*models.CreateURLResponse, error)
	GetOriginalURL(shortCode string) (string, error)
	GetStats(shortCode string) (*models.URLStatsResponse, error)
	DeleteURL(shortCode string, userID *uint64) error
}

type urlService struct {
	repo    repository.URLRepository
	baseURL string
}

func NewURLService(repo repository.URLRepository, baseURL string) URLService {
	return &urlService{
		repo:    repo,
		baseURL: baseURL,
	}
}

func (s *urlService) CreateShortURL(req *models.CreateURLRequest, userID *uint64) (*models.CreateURLResponse, error) {
	urlToValidate := req.URL
	if !strings.HasPrefix(req.URL, "http://") && !strings.HasPrefix(req.URL, "https://") {
		urlToValidate = "https://" + req.URL
	}

	if _, err := url.ParseRequestURI(urlToValidate); err != nil {
		return nil, ErrInvalidURL
	}

	shortCode := req.CustomCode
	if shortCode == "" {
		shortCode = s.generateShortCode()
	}

	var expiresAt *time.Time
	if req.ExpiresIn != nil && *req.ExpiresIn > 0 {
		t := time.Now().Add(time.Duration(*req.ExpiresIn) * time.Hour)
		expiresAt = &t
	}

	url := &models.URL{
		ShortCode:   shortCode,
		OriginalURL: urlToValidate,
		UserID:      userID,
		ExpiresAt:   expiresAt,
	}

	if err := s.repo.Create(url); err != nil {
		return nil, err
	}

	return &models.CreateURLResponse{
		ShortCode: shortCode,
		ShortURL:  s.baseURL + "/" + shortCode,
		ExpiresAt: expiresAt,
	}, nil
}

func (s *urlService) GetOriginalURL(shortCode string) (string, error) {
	url, err := s.repo.FindByShortCode(shortCode)
	if err != nil {
		return "", err
	}

	if url.ExpiresAt != nil && url.ExpiresAt.Before(time.Now()) {
		return "", ErrURLExpired
	}

	if err := s.repo.IncrementClicks(shortCode); err != nil {
		return url.OriginalURL, nil
	}

	return url.OriginalURL, nil
}

func (s *urlService) GetStats(shortCode string) (*models.URLStatsResponse, error) {
	url, err := s.repo.FindByShortCode(shortCode)
	if err != nil {
		return nil, err
	}

	return &models.URLStatsResponse{
		ShortCode:   url.ShortCode,
		OriginalURL: url.OriginalURL,
		Clicks:      url.Clicks,
		CreatedAt:   url.CreatedAt,
		ExpiresAt:   url.ExpiresAt,
	}, nil
}

func (s *urlService) DeleteURL(shortCode string, userID *uint64) error {
	url, err := s.repo.FindByShortCode(shortCode)
	if err != nil {
		return err
	}

	if userID != nil && url.UserID != nil && *userID != *url.UserID {
		return errors.New("unauthorized")
	}

	return s.repo.Delete(shortCode)
}

func (s *urlService) generateShortCode() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return base64.URLEncoding.EncodeToString(bytes)[:8]
}