package service

import (
	"errors"
	"math"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type AppStoreService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewAppStoreService(db *gorm.DB, log *logger.Logger) *AppStoreService {
	return &AppStoreService{db: db, log: log}
}

// ---------- App ----------

// GetList returns paginated apps with optional category and keyword filter.
func (s *AppStoreService) GetList(page, pageSize int, category, keyword string) ([]model.App, int64, error) {
	var apps []model.App
	var total int64

	query := s.db.Model(&model.App{}).Where("status = ?", 1)
	if category != "" {
		query = query.Where("category = ?", category)
	}
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR description LIKE ? OR author LIKE ?", q, q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("featured DESC, sort_order ASC, downloads DESC").
		Find(&apps).Error; err != nil {
		return nil, 0, err
	}
	return apps, total, nil
}

// GetByID returns a single app by ID.
func (s *AppStoreService) GetByID(id uint) (*model.App, error) {
	var app model.App
	if err := s.db.First(&app, id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

// GetBySlug returns a single app by slug.
func (s *AppStoreService) GetBySlug(slug string) (*model.App, error) {
	var app model.App
	if err := s.db.Where("slug = ?", slug).First(&app).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

// Install installs an app for a user.
func (s *AppStoreService) Install(userID, appID uint) (*model.AppInstall, error) {
	var app model.App
	if err := s.db.First(&app, appID).Error; err != nil {
		return nil, errors.New("app not found")
	}
	if app.Status != 1 {
		return nil, errors.New("app is not available")
	}

	var existing model.AppInstall
	if err := s.db.Where("user_id = ? AND app_id = ? AND status = 1", userID, appID).
		First(&existing).Error; err == nil {
		return nil, errors.New("app already installed")
	}

	install := &model.AppInstall{
		UserID:  userID,
		AppID:   appID,
		Version: app.Version,
		Status:  1,
	}
	if err := s.db.Create(install).Error; err != nil {
		return nil, err
	}

	s.db.Model(&model.App{}).Where("id = ?", appID).
		UpdateColumn("downloads", gorm.Expr("downloads + 1"))

	s.db.Preload("App").First(install, install.ID)
	return install, nil
}

// Uninstall marks an app installation as uninstalled.
func (s *AppStoreService) Uninstall(userID, appID uint) error {
	return s.db.Model(&model.AppInstall{}).
		Where("user_id = ? AND app_id = ? AND status = 1", userID, appID).
		Update("status", 3).Error
}

// Update updates an installed app to the latest version.
func (s *AppStoreService) Update(userID, appID uint) (*model.AppInstall, error) {
	var install model.AppInstall
	if err := s.db.Where("user_id = ? AND app_id = ? AND status = 1", userID, appID).
		First(&install).Error; err != nil {
		return nil, errors.New("app not installed")
	}

	var app model.App
	if err := s.db.First(&app, appID).Error; err != nil {
		return nil, err
	}
	if install.Version == app.Version {
		return nil, errors.New("already on latest version")
	}

	install.Version = app.Version
	install.Status = 1
	if err := s.db.Save(&install).Error; err != nil {
		return nil, err
	}

	s.db.Preload("App").First(&install, install.ID)
	return &install, nil
}

// GetInstalledApps returns all installed apps for a user.
func (s *AppStoreService) GetInstalledApps(userID uint) ([]model.AppInstall, error) {
	var installs []model.AppInstall
	if err := s.db.Preload("App").
		Where("user_id = ? AND status = 1", userID).
		Order("updated_at DESC").
		Find(&installs).Error; err != nil {
		return nil, err
	}
	return installs, nil
}

// ---------- Review ----------

// CreateReview creates or updates a review for an app.
func (s *AppStoreService) CreateReview(userID, appID uint, rating int8, content string) (*model.AppReview, error) {
	if rating < 1 || rating > 5 {
		return nil, errors.New("rating must be between 1 and 5")
	}

	var existing model.AppReview
	if err := s.db.Where("user_id = ? AND app_id = ?", userID, appID).
		First(&existing).Error; err == nil {
		existing.Rating = rating
		existing.Content = content
		if err := s.db.Save(&existing).Error; err != nil {
			return nil, err
		}
		s.recalcRating(appID)
		return &existing, nil
	}

	review := &model.AppReview{
		UserID:  userID,
		AppID:   appID,
		Rating:  rating,
		Content: content,
	}
	if err := s.db.Create(review).Error; err != nil {
		return nil, err
	}
	s.recalcRating(appID)
	return review, nil
}

// GetReviews returns paginated reviews for an app.
func (s *AppStoreService) GetReviews(appID uint, page, pageSize int) ([]model.AppReview, int64, error) {
	var reviews []model.AppReview
	var total int64

	query := s.db.Model(&model.AppReview{}).Where("app_id = ?", appID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("id DESC").
		Find(&reviews).Error; err != nil {
		return nil, 0, err
	}
	return reviews, total, nil
}

// recalcRating recalculates the average rating and count for an app.
func (s *AppStoreService) recalcRating(appID uint) {
	var count int64
	var avg float64
	s.db.Model(&model.AppReview{}).Where("app_id = ?", appID).Count(&count)
	if count > 0 {
		s.db.Model(&model.AppReview{}).Where("app_id = ?", appID).
			Select("AVG(rating)").Scan(&avg)
		avg = math.Round(avg*100) / 100
	}
	s.db.Model(&model.App{}).Where("id = ?", appID).Updates(map[string]interface{}{
		"rating":       avg,
		"rating_count": count,
	})
}

// ---------- Admin ----------

// AdminGetList returns all apps including inactive with pagination.
func (s *AppStoreService) AdminGetList(page, pageSize int, keyword string) ([]model.App, int64, error) {
	var apps []model.App
	var total int64

	query := s.db.Model(&model.App{})
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("name LIKE ? OR author LIKE ?", q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).
		Order("id DESC").
		Find(&apps).Error; err != nil {
		return nil, 0, err
	}
	return apps, total, nil
}

// Create creates a new app (admin).
func (s *AppStoreService) Create(req CreateAppRequest) (*model.App, error) {
	app := &model.App{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Icon:        req.Icon,
		Banner:      req.Banner,
		Version:     req.Version,
		Author:      req.Author,
		Website:     req.Website,
		Price:       req.Price,
		Currency:    req.Currency,
		Category:    req.Category,
		Tags:        req.Tags,
		Homepage:    req.Homepage,
		Repository:  req.Repository,
		MinVersion:  req.MinVersion,
		MaxVersion:  req.MaxVersion,
		Status:      req.Status,
		SortOrder:   req.SortOrder,
		Featured:    req.Featured,
		Verified:    req.Verified,
		Screenshots: req.Screenshots,
		Metadata:    req.Metadata,
	}
	if app.Currency == "" {
		app.Currency = "CNY"
	}
	if app.Status == 0 {
		app.Status = 1
	}
	if err := s.db.Create(app).Error; err != nil {
		return nil, err
	}
	return app, nil
}

// Update updates an app (admin).
func (s *AppStoreService) Update(id uint, req UpdateAppRequest) (*model.App, error) {
	var app model.App
	if err := s.db.First(&app, id).Error; err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}
	if req.Banner != nil {
		updates["banner"] = *req.Banner
	}
	if req.Version != nil {
		updates["version"] = *req.Version
	}
	if req.Author != nil {
		updates["author"] = *req.Author
	}
	if req.Website != nil {
		updates["website"] = *req.Website
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Currency != nil {
		updates["currency"] = *req.Currency
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Featured != nil {
		updates["featured"] = *req.Featured
	}
	if req.Verified != nil {
		updates["verified"] = *req.Verified
	}

	if len(updates) > 0 {
		if err := s.db.Model(&app).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	if err := s.db.First(&app, id).Error; err != nil {
		return nil, err
	}
	return &app, nil
}

// Delete soft-deletes an app (admin).
func (s *AppStoreService) Delete(id uint) error {
	return s.db.Delete(&model.App{}, id).Error
}

// ---------- Request DTOs ----------

type CreateAppRequest struct {
	Name        string  `json:"name" binding:"required,max=128"`
	Slug        string  `json:"slug" binding:"required,max=128"`
	Description string  `json:"description"`
	Icon        string  `json:"icon"`
	Banner      string  `json:"banner"`
	Version     string  `json:"version" binding:"required"`
	Author      string  `json:"author"`
	Website     string  `json:"website"`
	Price       float64 `json:"price"`
	Currency    string  `json:"currency"`
	Category    string  `json:"category"`
	Tags        string  `json:"tags"`
	Homepage    string  `json:"homepage"`
	Repository  string  `json:"repository"`
	MinVersion  string  `json:"min_version"`
	MaxVersion  string  `json:"max_version"`
	Status      int16   `json:"status"`
	SortOrder   int     `json:"sort_order"`
	Featured    bool    `json:"featured"`
	Verified    bool    `json:"verified"`
	Screenshots string  `json:"screenshots"`
	Metadata    string  `json:"metadata"`
}

type UpdateAppRequest struct {
	Name        *string  `json:"name"`
	Slug        *string  `json:"slug"`
	Description *string  `json:"description"`
	Icon        *string  `json:"icon"`
	Banner      *string  `json:"banner"`
	Version     *string  `json:"version"`
	Author      *string  `json:"author"`
	Website     *string  `json:"website"`
	Price       *float64 `json:"price"`
	Currency    *string  `json:"currency"`
	Category    *string  `json:"category"`
	Status      *int16   `json:"status"`
	SortOrder   *int     `json:"sort_order"`
	Featured    *bool    `json:"featured"`
	Verified    *bool    `json:"verified"`
}

type CreateReviewRequest struct {
	Rating  int8   `json:"rating" binding:"required,gte=1,lte=5"`
	Content string `json:"content"`
}
