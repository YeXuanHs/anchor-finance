package service

import (
	"errors"
	"time"

	"anchorfinance/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// ProductInterflow 产品关联/互通配置
type ProductInterflow struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"index;not null" json:"user_id"`
	SourceProductID uint         `gorm:"index;not null" json:"source_product_id"`
	TargetProductID uint         `gorm:"index;not null" json:"target_product_id"`
	RelationType  string         `gorm:"type:varchar(32);not null;default:'link'" json:"relation_type"` // link/sync/depend/share
	Config        datatypes.JSON `gorm:"type:json" json:"config"`
	Status        int16          `gorm:"type:smallint;default:1;not null" json:"status"` // 1=启用 0=禁用
	Remark        string         `gorm:"type:text" json:"remark"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

type InterflowService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewInterflowService(db *gorm.DB, log *logger.Logger) *InterflowService {
	return &InterflowService{db: db, log: log}
}

type CreateInterflowRequest struct {
	SourceProductID uint           `json:"source_product_id" binding:"required"`
	TargetProductID uint           `json:"target_product_id" binding:"required"`
	RelationType    string         `json:"relation_type" binding:"omitempty,oneof=link sync depend share"`
	Config          datatypes.JSON `json:"config"`
	Remark          string         `json:"remark" binding:"omitempty,max=500"`
}

type UpdateInterflowRequest struct {
	RelationType string         `json:"relation_type" binding:"omitempty,oneof=link sync depend share"`
	Config       datatypes.JSON `json:"config"`
	Status       *int16         `json:"status" binding:"omitempty,oneof=0 1"`
	Remark       string         `json:"remark" binding:"omitempty,max=500"`
}

// Create creates a new product interflow relation.
func (s *InterflowService) Create(userID uint, req CreateInterflowRequest) (*ProductInterflow, error) {
	if req.SourceProductID == req.TargetProductID {
		return nil, errors.New("source and target product cannot be the same")
	}

	// Verify both products belong to the user
	var count int64
	s.db.Table("user_products").Where("user_id = ? AND (product_id = ? OR product_id = ?)",
		userID, req.SourceProductID, req.TargetProductID).Count(&count)
	if count < 2 {
		return nil, errors.New("both products must belong to you")
	}

	// Check for existing relation
	var existing ProductInterflow
	err := s.db.Where("user_id = ? AND source_product_id = ? AND target_product_id = ? AND deleted_at IS NULL",
		userID, req.SourceProductID, req.TargetProductID).First(&existing).Error
	if err == nil {
		return nil, errors.New("relation already exists")
	}

	relationType := req.RelationType
	if relationType == "" {
		relationType = "link"
	}

	interflow := &ProductInterflow{
		UserID:          userID,
		SourceProductID: req.SourceProductID,
		TargetProductID: req.TargetProductID,
		RelationType:    relationType,
		Config:          req.Config,
		Status:          1,
		Remark:          req.Remark,
	}

	if err := s.db.Create(interflow).Error; err != nil {
		return nil, err
	}

	s.log.Infof("product interflow created: user=%d source=%d target=%d", userID, req.SourceProductID, req.TargetProductID)
	return interflow, nil
}

// GetByID returns an interflow relation by ID.
func (s *InterflowService) GetByID(userID, id uint) (*ProductInterflow, error) {
	var interflow ProductInterflow
	if err := s.db.Where("id = ? AND user_id = ?", id, userID).First(&interflow).Error; err != nil {
		return nil, err
	}
	return &interflow, nil
}

// GetList returns paginated interflow relations for a user.
func (s *InterflowService) GetList(userID uint, page, pageSize int, relationType string) ([]ProductInterflow, int64, error) {
	var items []ProductInterflow
	var total int64

	query := s.db.Model(&ProductInterflow{}).Where("user_id = ?", userID)
	if relationType != "" {
		query = query.Where("relation_type = ?", relationType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := Paginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("id DESC").Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

// GetByProduct returns all interflow relations for a specific product.
func (s *InterflowService) GetByProduct(userID, productID uint) ([]ProductInterflow, error) {
	var items []ProductInterflow
	if err := s.db.Where("user_id = ? AND (source_product_id = ? OR target_product_id = ?)",
		userID, productID, productID).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// Update updates an interflow relation.
func (s *InterflowService) Update(userID, id uint, req UpdateInterflowRequest) (*ProductInterflow, error) {
	interflow, err := s.GetByID(userID, id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.RelationType != "" {
		updates["relation_type"] = req.RelationType
	}
	if req.Config != nil {
		updates["config"] = req.Config
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Remark != "" {
		updates["remark"] = req.Remark
	}

	if len(updates) > 0 {
		if err := s.db.Model(interflow).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	return s.GetByID(userID, id)
}

// Delete soft-deletes an interflow relation.
func (s *InterflowService) Delete(userID, id uint) error {
	result := s.db.Where("id = ? AND user_id = ?", id, userID).Delete(&ProductInterflow{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("relation not found")
	}
	return nil
}

// ToggleStatus toggles the status of an interflow relation.
func (s *InterflowService) ToggleStatus(userID, id uint) error {
	interflow, err := s.GetByID(userID, id)
	if err != nil {
		return err
	}

	newStatus := int16(1)
	if interflow.Status == 1 {
		newStatus = 0
	}
	return s.db.Model(interflow).Update("status", newStatus).Error
}
