package service

import (
	"time"

	"github.com/anchor-finance/backend/pkg/logger"

	"gorm.io/gorm"
)

type Product struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:128;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Price       float64        `gorm:"type:decimal(12,2);not null" json:"price"`
	Period      int            `gorm:"not null;comment:duration in days" json:"period"`
	PeriodUnit  string         `gorm:"size:16;default:day;comment:day/month/year" json:"period_unit"`
	Category    string         `gorm:"size:64" json:"category"`
	Stock       int            `gorm:"default:-1;comment:-1=unlimited" json:"stock"`
	Sort        int            `gorm:"default:0" json:"sort"`
	Status      int            `gorm:"default:1;comment:1=active 0=disabled" json:"status"`
	Config      string         `gorm:"type:text;comment:json config" json:"config"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserProduct struct {
	ID         uint       `gorm:"primaryKey" json:"id"`
	UserID     uint       `gorm:"index;not null" json:"user_id"`
	ProductID  uint       `gorm:"not null" json:"product_id"`
	Product    Product    `gorm:"foreignKey:ProductID" json:"product"`
	OrderID    uint       `json:"order_id"`
	OrderNo    string     `gorm:"size:64" json:"order_no"`
	StartAt    time.Time  `json:"start_at"`
	ExpireAt   time.Time  `gorm:"index" json:"expire_at"`
	Status     int        `gorm:"default:1;comment:1=active 2=expired 3=cancelled" json:"status"`
	AutoRenew  bool       `gorm:"default:false" json:"auto_renew"`
	Remark     string     `gorm:"size:256" json:"remark"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
}

type ProductService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewProductService(db *gorm.DB, log *logger.Logger) *ProductService {
	return &ProductService{db: db, log: log}
}

type CreateProductRequest struct {
	Name        string  `json:"name" binding:"required,max=128"`
	Description string  `json:"description"`
	Price       float64 `json:"price" binding:"required,gt=0"`
	Period      int     `json:"period" binding:"required,gt=0"`
	PeriodUnit  string  `json:"period_unit" binding:"omitempty,oneof=day month year"`
	Category    string  `json:"category"`
	Stock       int     `json:"stock"`
	Sort        int     `json:"sort"`
	Config      string  `json:"config"`
}

type UpdateProductRequest struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Price       *float64 `json:"price"`
	Period      *int     `json:"period"`
	PeriodUnit  *string  `json:"period_unit"`
	Category    *string  `json:"category"`
	Stock       *int     `json:"stock"`
	Sort        *int     `json:"sort"`
	Status      *int     `json:"status"`
	Config      *string  `json:"config"`
}

// GetList returns active products with pagination.
func (s *ProductService) GetList(page, pageSize int, category string) ([]Product, int64, error) {
	var products []Product
	var total int64

	query := s.db.Model(&Product{}).Where("status = ?", 1)
	if category != "" {
		query = query.Where("category = ?", category)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("sort ASC, id DESC").Find(&products).Error; err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// GetByID returns a single product.
func (s *ProductService) GetByID(id uint) (*Product, error) {
	var product Product
	if err := s.db.First(&product, id).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// Create adds a new product.
func (s *ProductService) Create(req CreateProductRequest) (*Product, error) {
	product := &Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Period:      req.Period,
		PeriodUnit:  req.PeriodUnit,
		Category:    req.Category,
		Stock:       req.Stock,
		Sort:        req.Sort,
		Status:      1,
		Config:      req.Config,
	}
	if product.PeriodUnit == "" {
		product.PeriodUnit = "day"
	}
	if err := s.db.Create(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

// Update modifies an existing product.
func (s *ProductService) Update(id uint, req UpdateProductRequest) (*Product, error) {
	product, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Price != nil {
		updates["price"] = *req.Price
	}
	if req.Period != nil {
		updates["period"] = *req.Period
	}
	if req.PeriodUnit != nil {
		updates["period_unit"] = *req.PeriodUnit
	}
	if req.Category != nil {
		updates["category"] = *req.Category
	}
	if req.Stock != nil {
		updates["stock"] = *req.Stock
	}
	if req.Sort != nil {
		updates["sort"] = *req.Sort
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Config != nil {
		updates["config"] = *req.Config
	}

	if err := s.db.Model(product).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

// Delete soft-deletes a product.
func (s *ProductService) Delete(id uint) error {
	return s.db.Delete(&Product{}, id).Error
}

// GetUserProducts returns active products owned by a user.
func (s *ProductService) GetUserProducts(userID uint) ([]UserProduct, error) {
	var ups []UserProduct
	if err := s.db.Preload("Product").
		Where("user_id = ? AND status = 1", userID).
		Order("expire_at DESC").
		Find(&ups).Error; err != nil {
		return nil, err
	}
	return ups, nil
}
