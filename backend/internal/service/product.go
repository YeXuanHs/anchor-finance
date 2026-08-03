package service

import (
	"fmt"
	"time"

	"anchorfinance/pkg/logger"

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

// GetHotProducts returns the top selling active products.
func (s *ProductService) GetHotProducts(limit int) ([]Product, error) {
	var products []Product
	if err := s.db.Where("status = ?", 1).
		Order("sort ASC, id DESC").
		Limit(limit).
		Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
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

// --- Product Group ---

// ProductGroup represents a product category/group.
type ProductGroup struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"size:128;not null" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	Slug        string `gorm:"size:128;uniqueIndex" json:"slug"`
	ParentID    *uint  `gorm:"index" json:"parent_id"`
	Icon        string `gorm:"size:512" json:"icon"`
	SortOrder   int    `gorm:"default:0" json:"sort_order"`
	Status      int16  `gorm:"default:1;not null" json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (ProductGroup) TableName() string { return "product_groups" }

type CreateGroupRequest struct {
	Name        string `json:"name" binding:"required,max=128"`
	Description string `json:"description"`
	Slug        string `json:"slug" binding:"omitempty,max=128"`
	ParentID    *uint  `json:"parent_id"`
	Icon        string `json:"icon"`
	SortOrder   int    `json:"sort_order"`
}

type UpdateGroupRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Slug        *string `json:"slug"`
	ParentID    *uint   `json:"parent_id"`
	Icon        *string `json:"icon"`
	SortOrder   *int    `json:"sort_order"`
	Status      *int16  `json:"status"`
}

// GetGroups returns all product groups.
func (s *ProductService) GetGroups() ([]ProductGroup, error) {
	var groups []ProductGroup
	if err := s.db.Order("sort_order ASC, id ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// GetGroupByID returns a single product group by ID.
func (s *ProductService) GetGroupByID(id uint) (*ProductGroup, error) {
	var group ProductGroup
	if err := s.db.First(&group, id).Error; err != nil {
		return nil, err
	}
	return &group, nil
}

// CreateGroup creates a new product group.
func (s *ProductService) CreateGroup(req CreateGroupRequest) (*ProductGroup, error) {
	group := &ProductGroup{
		Name:        req.Name,
		Description: req.Description,
		Slug:        req.Slug,
		ParentID:    req.ParentID,
		Icon:        req.Icon,
		SortOrder:   req.SortOrder,
		Status:      1,
	}
	if err := s.db.Create(group).Error; err != nil {
		return nil, err
	}
	return group, nil
}

// UpdateGroup updates a product group.
func (s *ProductService) UpdateGroup(id uint, req UpdateGroupRequest) (*ProductGroup, error) {
	group, err := s.GetGroupByID(id)
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
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.ParentID != nil {
		updates["parent_id"] = *req.ParentID
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if len(updates) > 0 {
		if err := s.db.Model(group).Updates(updates).Error; err != nil {
			return nil, err
		}
	}
	return s.GetGroupByID(id)
}

// DeleteGroup deletes a product group.
func (s *ProductService) DeleteGroup(id uint) error {
	return s.db.Delete(&ProductGroup{}, id).Error
}

// ==================== First Group (一级分组) ====================

// FirstGroupDTO is a lightweight struct for first group queries.
type FirstGroupDTO struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	Hidden     bool   `json:"hidden"`
	SortOrder  int    `json:"sort_order"`
	UpstreamID *uint  `json:"upstream_id"`
	CreatedAt  time.Time `json:"created_at"`
}

func (FirstGroupDTO) TableName() string { return "product_first_groups" }

type CreateFirstGroupRequest struct {
	Name       string `json:"name" binding:"required,max=128"`
	Hidden     bool   `json:"hidden"`
	UpstreamID *uint  `json:"upstream_id"`
}

type UpdateFirstGroupRequest struct {
	Name       *string `json:"name"`
	Hidden     *bool   `json:"hidden"`
	UpstreamID *uint   `json:"upstream_id"`
}

// GetFirstGroups returns all first-level groups.
func (s *ProductService) GetFirstGroups() ([]FirstGroupDTO, error) {
	var groups []FirstGroupDTO
	if err := s.db.Table("product_first_groups").Order("sort_order ASC, id ASC").Find(&groups).Error; err != nil {
		return nil, err
	}
	return groups, nil
}

// CreateFirstGroup creates a new first-level group.
func (s *ProductService) CreateFirstGroup(req CreateFirstGroupRequest) (*FirstGroupDTO, error) {
	group := map[string]interface{}{
		"name":        req.Name,
		"hidden":      req.Hidden,
		"sort_order":  0,
		"upstream_id": req.UpstreamID,
		"created_at":  time.Now(),
		"updated_at":  time.Now(),
	}
	if err := s.db.Table("product_first_groups").Create(&group).Error; err != nil {
		return nil, err
	}
	var result FirstGroupDTO
	s.db.Table("product_first_groups").Order("id DESC").First(&result)
	return &result, nil
}

// UpdateFirstGroup updates a first-level group.
func (s *ProductService) UpdateFirstGroup(id uint, req UpdateFirstGroupRequest) error {
	updates := map[string]interface{}{"updated_at": time.Now()}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Hidden != nil {
		updates["hidden"] = *req.Hidden
	}
	if req.UpstreamID != nil {
		updates["upstream_id"] = *req.UpstreamID
	}
	return s.db.Table("product_first_groups").Where("id = ?", id).Updates(updates).Error
}

// DeleteFirstGroup deletes a first-level group if it has no child groups.
func (s *ProductService) DeleteFirstGroup(id uint) error {
	if id == 1 {
		return fmt.Errorf("cannot delete the default group")
	}
	var count int64
	s.db.Table("product_groups").Where("first_group_id = ?", id).Count(&count)
	if count > 0 {
		return fmt.Errorf("cannot delete group with child groups")
	}
	return s.db.Table("product_first_groups").Where("id = ?", id).Delete(nil).Error
}

// ==================== Sort Management ====================

type SortItem struct {
	ID       uint `json:"id"`
	SortOrder int  `json:"sort_order"`
}

type BatchSortRequest struct {
	Items []SortItem `json:"items" binding:"required,min=1"`
}

// UpdateFirstGroupSort batch-updates first-level group sort orders.
func (s *ProductService) UpdateFirstGroupSort(items []SortItem) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Table("product_first_groups").Where("id = ?", item.ID).
				Update("sort_order", item.SortOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateGroupSort batch-updates second-level group sort orders.
func (s *ProductService) UpdateGroupSort(items []SortItem) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Table("product_groups").Where("id = ?", item.ID).
				Update("sort_order", item.SortOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// UpdateProductSort batch-updates product sort orders.
func (s *ProductService) UpdateProductSort(items []SortItem) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range items {
			if err := tx.Table("products").Where("id = ?", item.ID).
				Update("sort_order", item.SortOrder).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// ==================== Duplicate Product ====================

type DuplicateProductRequest struct {
	ExistingProductID uint   `json:"existing_product_id" binding:"required"`
	NewProductName    string `json:"new_product_name" binding:"required,max=256"`
}

// DuplicateProduct copies a product along with its pricing, config links, custom fields, and downloads.
func (s *ProductService) DuplicateProduct(req DuplicateProductRequest) (*Product, error) {
	src, err := s.GetByID(req.ExistingProductID)
	if err != nil {
		return nil, fmt.Errorf("source product not found: %w", err)
	}

	var newProduct *Product
	err = s.db.Transaction(func(tx *gorm.DB) error {
		// 1. Copy the product record (exclude ID, timestamps)
		newProduct = &Product{
			GroupID:       src.GroupID,
			Name:          req.NewProductName,
			Description:   src.Description,
			Content:       src.Content,
			Price:         src.Price,
			OriginalPrice: src.OriginalPrice,
			Currency:      src.Currency,
			BillingCycle:  src.BillingCycle,
			SetupFee:      src.SetupFee,
			Stock:         src.Stock,
			SalesCount:    0,
			Type:          src.Type,
			AutoSetup:     src.AutoSetup,
			StockControl:  src.StockControl,
			QuantityMin:   src.QuantityMin,
			QuantityMax:   src.QuantityMax,
			TrialDays:     src.TrialDays,
			SortOrder:     0,
			Status:        0, // 默认下架，待管理员审核
			Hidden:        true,
			Download:      src.Download,
			Featured:      false,
			Image:         src.Image,
			Images:        src.Images,
			ConfigOptions: src.ConfigOptions,
			Metadata:      src.Metadata,
			Tags:          src.Tags,
		}
		if err := tx.Create(newProduct).Error; err != nil {
			return err
		}

		// 2. Copy pricing records
		var pricings []ProductPricing
		if err := tx.Where("product_id = ?", src.ID).Find(&pricings).Error; err == nil && len(pricings) > 0 {
			newPricings := make([]ProductPricing, len(pricings))
			for i, p := range pricings {
				newPricings[i] = ProductPricing{
					ProductID: newProduct.ID,
					Cycle:     p.Cycle,
					Price:     p.Price,
					SetupFee:  p.SetupFee,
					Currency:  p.Currency,
					SortOrder: p.SortOrder,
				}
			}
			tx.Create(&newPricings)
		}

		// 3. Copy config option links (ProductConfigOptionLink)
		var configLinks []struct {
			GroupID   uint `gorm:"column:group_id"`
			SortOrder int  `gorm:"column:sort_order"`
		}
		if err := tx.Table("product_config_option_links").Where("product_id = ?", src.ID).
			Find(&configLinks).Error; err == nil && len(configLinks) > 0 {
			newLinks := make([]map[string]interface{}, len(configLinks))
			for i, cl := range configLinks {
				newLinks[i] = map[string]interface{}{
					"product_id": newProduct.ID,
					"group_id":   cl.GroupID,
					"sort_order": cl.SortOrder,
					"created_at": time.Now(),
					"updated_at": time.Now(),
				}
			}
			tx.Table("product_config_option_links").Create(&newLinks)
		}

		// 4. Copy custom fields
		var customFields []struct {
			ID       uint   `gorm:"column:id"`
			Name     string `gorm:"column:name"`
			Label    string `gorm:"column:label"`
			Type     string `gorm:"column:type"`
			Group    string `gorm:"column:group"`
			Required bool   `gorm:"column:required"`
			SortOrder int   `gorm:"column:sort_order"`
		}
		if err := tx.Table("custom_fields").Where("`group` = 'product' AND id IN (?)",
			tx.Table("custom_field_values").Select("field_id").Where("owner_id = ? AND owner_type = 'product'", src.ID),
		).Find(&customFields).Error; err == nil && len(customFields) > 0 {
			for _, cf := range customFields {
				newField := map[string]interface{}{
					"name":       cf.Name,
					"label":      cf.Label,
					"type":       cf.Type,
					"group":      "product",
					"required":   cf.Required,
					"sort_order": cf.SortOrder,
					"enabled":    true,
					"created_at": time.Now(),
					"updated_at": time.Now(),
				}
				tx.Table("custom_fields").Create(&newField)
			}
		}

		// 5. Copy product-download associations
		var downloads []ProductDownload
		if err := tx.Where("product_id = ?", src.ID).Find(&downloads).Error; err == nil && len(downloads) > 0 {
			newDownloads := make([]ProductDownload, len(downloads))
			for i, d := range downloads {
				newDownloads[i] = ProductDownload{
					ProductID:  newProduct.ID,
					DownloadID: d.DownloadID,
				}
			}
			tx.Create(&newDownloads)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return newProduct, nil
}

// ==================== Edit Stock ====================

type EditStockRequest struct {
	Stock int `json:"stock" binding:"required"`
}

// EditStock updates a product's stock quantity.
func (s *ProductService) EditStock(id uint, stock int) error {
	return s.db.Table("products").Where("id = ?", id).Update("stock", stock).Error
}

// ==================== Batch Operations ====================

type BatchUpdateRequest struct {
	IDs    []uint                 `json:"ids" binding:"required,min=1"`
	Fields map[string]interface{} `json:"fields" binding:"required"`
}

type BatchDeleteRequest struct {
	IDs []uint `json:"ids" binding:"required,min=1"`
}

// BatchUpdate updates multiple products at once.
func (s *ProductService) BatchUpdate(ids []uint, fields map[string]interface{}) error {
	allowedFields := map[string]bool{
		"name": true, "description": true, "price": true, "status": true,
		"hidden": true, "featured": true, "stock": true, "sort_order": true,
		"group_id": true, "type": true, "auto_setup": true, "stock_control": true,
	}
	safeFields := map[string]interface{}{}
	for k, v := range fields {
		if allowedFields[k] {
			safeFields[k] = v
		}
	}
	if len(safeFields) == 0 {
		return fmt.Errorf("no valid fields to update")
	}
	safeFields["updated_at"] = time.Now()
	return s.db.Table("products").Where("id IN ?", ids).Updates(safeFields).Error
}

// BatchDelete soft-deletes multiple products.
func (s *ProductService) BatchDelete(ids []uint) error {
	return s.db.Table("products").Where("id IN ?", ids).
		Update("deleted_at", time.Now()).Error
}

// ==================== Check Alias ====================

type CheckAliasRequest struct {
	Alias string `json:"alias" binding:"required"`
	ID    uint   `json:"id"` // 排除自身
}

// CheckAlias checks if a group alias is already in use.
func (s *ProductService) CheckAlias(alias string, excludeID uint) (bool, error) {
	var count int64
	query := s.db.Table("product_groups").Where("alias = ?", alias)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	if err := query.Count(&count).Error; err != nil {
		return false, err
	}
	return count == 0, nil
}

// ==================== Delete Custom Field ====================

// DeleteCustomField deletes a custom field by ID.
func (s *ProductService) DeleteCustomField(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// Delete associated values first
		if err := tx.Table("custom_field_values").Where("field_id = ?", id).Delete(nil).Error; err != nil {
			return err
		}
		// Delete the field
		return tx.Table("custom_fields").Where("id = ?", id).Delete(nil).Error
	})
}

// ==================== Download Management ====================

type ManageDownloadsRequest struct {
	ProductID  uint `json:"product_id" binding:"required"`
	AddID      uint `json:"add_id"`      // 要添加的下载文件ID
	RemoveID   uint `json:"remove_id"`   // 要移除的下载文件ID
}

type AddDownloadFileRequest struct {
	ProductID   uint   `json:"product_id" binding:"required"`
	CategoryID  uint   `json:"category_id" binding:"required"`
	Title       string `json:"title" binding:"required,max=255"`
	Description string `json:"description"`
}

// ManageDownloads adds or removes download file associations for a product.
func (s *ProductService) ManageDownloads(req ManageDownloadsRequest) (string, error) {
	if req.AddID > 0 {
		// Check if already linked
		var count int64
		s.db.Table("product_downloads").Where("product_id = ? AND download_id = ?", req.ProductID, req.AddID).Count(&count)
		if count > 0 {
			return "already linked", nil
		}
		link := map[string]interface{}{
			"product_id":  req.ProductID,
			"download_id": req.AddID,
			"created_at":  time.Now(),
		}
		if err := s.db.Table("product_downloads").Create(&link).Error; err != nil {
			return "", err
		}
		return "added", nil
	}
	if req.RemoveID > 0 {
		if err := s.db.Table("product_downloads").Where("product_id = ? AND download_id = ?", req.ProductID, req.RemoveID).Delete(nil).Error; err != nil {
			return "", err
		}
		return "removed", nil
	}
	return "", fmt.Errorf("either add_id or remove_id is required")
}

// GetProductDownloads returns download files linked to a product.
func (s *ProductService) GetProductDownloads(productID uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	err := s.db.Table("download_files df").
		Select("df.id, df.title, df.description, df.file_path, df.file_size, df.file_type, df.download_count").
		Joins("JOIN product_downloads pd ON pd.download_id = df.id").
		Where("pd.product_id = ?", productID).
		Where("df.deleted_at IS NULL").
		Scan(&results).Error
	return results, err
}

// AddDownloadFile creates a new download file and links it to a product.
func (s *ProductService) AddDownloadFile(req AddDownloadFileRequest) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		file := map[string]interface{}{
			"category_id": req.CategoryID,
			"title":       req.Title,
			"description": req.Description,
			"file_path":   "",
			"file_size":   0,
			"is_published": true,
			"created_at":  time.Now(),
			"updated_at":  time.Now(),
		}
		if err := tx.Table("download_files").Create(&file).Error; err != nil {
			return err
		}
		// Get the new file ID
		var newID uint
		tx.Table("download_files").Select("MAX(id)").Scan(&newID)
		link := map[string]interface{}{
			"product_id":  req.ProductID,
			"download_id": newID,
			"created_at":  time.Now(),
		}
		return tx.Table("product_downloads").Create(&link).Error
	})
}

// ==================== Discount List ====================

type DiscountDTO struct {
	ID           uint    `json:"id"`
	ProductID    uint    `json:"product_id"`
	Cycle        string  `json:"cycle"`
	Price        float64 `json:"price"`
	SetupFee     float64 `json:"setup_fee"`
	Currency     string  `json:"currency"`
}

// GetDiscountList returns pricing/discount records for a product.
func (s *ProductService) GetDiscountList(productID uint) ([]DiscountDTO, error) {
	var list []DiscountDTO
	if err := s.db.Table("product_pricings").
		Where("product_id = ?", productID).
		Order("sort_order ASC, id ASC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// ==================== EditResProduct Fields ====================

// UpdateProductFields 直接更新产品字段（用于 EditResProduct）
func (s *ProductService) UpdateProductFields(id uint, updates map[string]interface{}) error {
	updates["updated_at"] = time.Now()
	return s.db.Table("products").Where("id = ?", id).Updates(updates).Error
}

// GetDB returns the database instance for direct queries.
func (s *ProductService) GetDB() *gorm.DB {
	return s.db
}

// ==================== Get Upstream Price ====================

// GetUpstreamPrice fetches upstream pricing for a product via its linked upstream provider.
func (s *ProductService) GetUpstreamPrice(productID uint) (map[string]interface{}, error) {
	// Find the upstream mapping for this product
	var mapping struct {
		UpstreamID      uint   `gorm:"column:upstream_id"`
		RemoteProductID string `gorm:"column:remote_product_id"`
		Config          string `gorm:"column:config"`
	}
	if err := s.db.Table("upstream_products").
		Where("local_product_id = ?", productID).First(&mapping).Error; err != nil {
		return nil, fmt.Errorf("no upstream mapping found for product %d", productID)
	}

	// Find the upstream provider
	var provider struct {
		ID     uint   `gorm:"column:id"`
		Name   string `gorm:"column:name"`
		Type   string `gorm:"column:type"`
		APIURL string `gorm:"column:api_url"`
		APIKey string `gorm:"column:api_key"`
	}
	if err := s.db.Table("upstream_providers").
		Where("id = ? AND is_active = true", mapping.UpstreamID).First(&provider).Error; err != nil {
		return nil, fmt.Errorf("upstream provider not found")
	}

	result := map[string]interface{}{
		"provider_id":       provider.ID,
		"provider_name":     provider.Name,
		"provider_type":     provider.Type,
		"remote_product_id": mapping.RemoteProductID,
		"config":            mapping.Config,
	}
	return result, nil
}
