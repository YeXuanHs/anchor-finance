package service

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"anchorfinance/pkg/logger"
)

// InvoiceItem 账单项
type InvoiceItem struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	InvoiceID   uint           `gorm:"index;not null" json:"invoice_id"`
	Type        string         `gorm:"size:32;not null;default:product" json:"type"` // product/fee/discount/credit/custom
	RelID       uint           `gorm:"index" json:"rel_id"`
	RelType     string         `gorm:"size:32" json:"rel_type"`
	Description string         `gorm:"size:512;not null" json:"description"`
	Quantity    int            `gorm:"default:1" json:"quantity"`
	UnitPrice   float64        `gorm:"type:decimal(20,4);not null" json:"unit_price"`
	Discount    float64        `gorm:"type:decimal(20,4);default:0" json:"discount"`
	Tax         float64        `gorm:"type:decimal(20,4);default:0" json:"tax"`
	Total       float64        `gorm:"type:decimal(20,4);not null" json:"total"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type InvoiceItemDiscount struct {
	ID             uint    `gorm:"primaryKey" json:"id"`
	ItemID         uint    `gorm:"index;not null" json:"item_id"`
	PromotionID    *uint   `gorm:"index" json:"promotion_id"`
	Type           string  `gorm:"size:32;not null" json:"type"` // coupon/promotion/bulk/manual
	DiscountType   string  `gorm:"size:16;not null" json:"discount_type"` // amount/percent
	Value          float64 `gorm:"type:decimal(12,2);not null" json:"value"`
	AppliedAmount  float64 `gorm:"type:decimal(12,2);not null" json:"applied_amount"`
	Description    string  `gorm:"size:256" json:"description"`
	CreatedAt      time.Time `json:"created_at"`
}

type InvoiceItemTax struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	ItemID      uint    `gorm:"index;not null" json:"item_id"`
	TaxName     string  `gorm:"size:64;not null" json:"tax_name"`
	TaxRate     float64 `gorm:"type:decimal(5,4);not null" json:"tax_rate"`
	TaxAmount   float64 `gorm:"type:decimal(12,2);not null" json:"tax_amount"`
	Description string  `gorm:"size:256" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type InvoiceItemService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewInvoiceItemService(db *gorm.DB, log *logger.Logger) *InvoiceItemService {
	return &InvoiceItemService{db: db, log: log}
}

type CreateInvoiceItemRequest struct {
	InvoiceID   uint    `json:"invoice_id" binding:"required"`
	Type        string  `json:"type" binding:"omitempty,oneof=product fee discount credit custom"`
	RelID       uint    `json:"rel_id"`
	RelType     string  `json:"rel_type"`
	Description string  `json:"description" binding:"required,max=512"`
	Quantity    int     `json:"quantity" binding:"omitempty,gte=1"`
	UnitPrice   float64 `json:"unit_price" binding:"required"`
	Discount    float64 `json:"discount"`
	Tax         float64 `json:"tax"`
	SortOrder   int     `json:"sort_order"`
}

type UpdateInvoiceItemRequest struct {
	Type        *string  `json:"type"`
	RelID       *uint    `json:"rel_id"`
	RelType     *string  `json:"rel_type"`
	Description *string  `json:"description"`
	Quantity    *int     `json:"quantity"`
	UnitPrice   *float64 `json:"unit_price"`
	Discount    *float64 `json:"discount"`
	Tax         *float64 `json:"tax"`
	SortOrder   *int     `json:"sort_order"`
}

type BatchCreateInvoiceItemRequest struct {
	InvoiceID uint                      `json:"invoice_id" binding:"required"`
	Items     []CreateInvoiceItemRequest `json:"items" binding:"required,min=1"`
}

// Create creates a new invoice item.
func (s *InvoiceItemService) Create(req CreateInvoiceItemRequest) (*InvoiceItem, error) {
	quantity := req.Quantity
	if quantity < 1 {
		quantity = 1
	}

	itemType := req.Type
	if itemType == "" {
		itemType = "product"
	}

	total := float64(quantity)*req.UnitPrice - req.Discount + req.Tax

	item := &InvoiceItem{
		InvoiceID:   req.InvoiceID,
		Type:        itemType,
		RelID:       req.RelID,
		RelType:     req.RelType,
		Description: req.Description,
		Quantity:    quantity,
		UnitPrice:   req.UnitPrice,
		Discount:    req.Discount,
		Tax:         req.Tax,
		Total:       total,
		SortOrder:   req.SortOrder,
	}

	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}

	s.log.Infof("invoice item created: invoice=%d desc=%s total=%.2f", req.InvoiceID, req.Description, total)
	return item, nil
}

// GetByID returns a single invoice item by ID.
func (s *InvoiceItemService) GetByID(id uint) (*InvoiceItem, error) {
	var item InvoiceItem
	if err := s.db.First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// Update modifies an existing invoice item.
func (s *InvoiceItemService) Update(id uint, req UpdateInvoiceItemRequest) (*InvoiceItem, error) {
	item, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := map[string]interface{}{}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.RelID != nil {
		updates["rel_id"] = *req.RelID
	}
	if req.RelType != nil {
		updates["rel_type"] = *req.RelType
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Quantity != nil {
		updates["quantity"] = *req.Quantity
	}
	if req.UnitPrice != nil {
		updates["unit_price"] = *req.UnitPrice
	}
	if req.Discount != nil {
		updates["discount"] = *req.Discount
	}
	if req.Tax != nil {
		updates["tax"] = *req.Tax
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	// Recalculate total if relevant fields changed
	if req.Quantity != nil || req.UnitPrice != nil || req.Discount != nil || req.Tax != nil {
		quantity := item.Quantity
		if req.Quantity != nil {
			quantity = *req.Quantity
		}
		unitPrice := item.UnitPrice
		if req.UnitPrice != nil {
			unitPrice = *req.UnitPrice
		}
		discount := item.Discount
		if req.Discount != nil {
			discount = *req.Discount
		}
		tax := item.Tax
		if req.Tax != nil {
			tax = *req.Tax
		}
		updates["total"] = float64(quantity)*unitPrice - discount + tax
	}

	if err := s.db.Model(item).Updates(updates).Error; err != nil {
		return nil, err
	}
	return s.GetByID(id)
}

// Delete soft-deletes an invoice item.
func (s *InvoiceItemService) Delete(id uint) error {
	return s.db.Delete(&InvoiceItem{}, id).Error
}

// GetByInvoiceID returns all items for an invoice.
func (s *InvoiceItemService) GetByInvoiceID(invoiceID uint) ([]InvoiceItem, error) {
	var items []InvoiceItem
	if err := s.db.Where("invoice_id = ?", invoiceID).Order("sort_order ASC, id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// BatchCreate creates multiple invoice items in a single transaction.
func (s *InvoiceItemService) BatchCreate(req BatchCreateInvoiceItemRequest) ([]InvoiceItem, error) {
	var items []InvoiceItem

	for _, itemReq := range req.Items {
		quantity := itemReq.Quantity
		if quantity < 1 {
			quantity = 1
		}

		itemType := itemReq.Type
		if itemType == "" {
			itemType = "product"
		}

		total := float64(quantity)*itemReq.UnitPrice - itemReq.Discount + itemReq.Tax

		item := InvoiceItem{
			InvoiceID:   req.InvoiceID,
			Type:        itemType,
			RelID:       itemReq.RelID,
			RelType:     itemReq.RelType,
			Description: itemReq.Description,
			Quantity:    quantity,
			UnitPrice:   itemReq.UnitPrice,
			Discount:    itemReq.Discount,
			Tax:         itemReq.Tax,
			Total:       total,
			SortOrder:   itemReq.SortOrder,
		}
		items = append(items, item)
	}

	if err := s.db.CreateInBatches(items, 100).Error; err != nil {
		return nil, err
	}

	s.log.Infof("batch invoice items created: invoice=%d count=%d", req.InvoiceID, len(items))
	return items, nil
}

// BatchDelete deletes multiple invoice items.
func (s *InvoiceItemService) BatchDelete(ids []uint) error {
	if len(ids) == 0 {
		return errors.New("empty ids")
	}
	return s.db.Delete(&InvoiceItem{}, ids).Error
}

// BatchUpdate updates multiple invoice items.
func (s *InvoiceItemService) BatchUpdate(ids []uint, req UpdateInvoiceItemRequest) error {
	if len(ids) == 0 {
		return errors.New("empty ids")
	}

	updates := map[string]interface{}{}
	if req.Type != nil {
		updates["type"] = *req.Type
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}

	if len(updates) == 0 {
		return nil
	}

	return s.db.Model(&InvoiceItem{}).Where("id IN ?", ids).Updates(updates).Error
}

// CalculateInvoiceTotal calculates the total for an invoice based on its items.
func (s *InvoiceItemService) CalculateInvoiceTotal(invoiceID uint) (map[string]float64, error) {
	var items []InvoiceItem
	if err := s.db.Where("invoice_id = ?", invoiceID).Find(&items).Error; err != nil {
		return nil, err
	}

	var subtotal, totalDiscount, totalTax float64
	for _, item := range items {
		subtotal += float64(item.Quantity) * item.UnitPrice
		totalDiscount += item.Discount
		totalTax += item.Tax
	}

	total := subtotal - totalDiscount + totalTax

	return map[string]float64{
		"sub_total":      subtotal,
		"total_discount": totalDiscount,
		"total_tax":      totalTax,
		"total":          total,
		"item_count":     float64(len(items)),
	}, nil
}

// AddDiscount adds a discount to an invoice item.
func (s *InvoiceItemService) AddDiscount(itemID uint, discountType string, value float64, description string) error {
	item, err := s.GetByID(itemID)
	if err != nil {
		return err
	}

	var appliedAmount float64
	if discountType == "amount" {
		appliedAmount = value
	} else {
		appliedAmount = float64(item.Quantity) * item.UnitPrice * value / 100
	}

	discount := &InvoiceItemDiscount{
		ItemID:        itemID,
		Type:          "manual",
		DiscountType:  discountType,
		Value:         value,
		AppliedAmount: appliedAmount,
		Description:   description,
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(discount).Error; err != nil {
			return err
		}
		newDiscount := item.Discount + appliedAmount
		newTotal := float64(item.Quantity)*item.UnitPrice - newDiscount + item.Tax
		return tx.Model(item).Updates(map[string]interface{}{
			"discount": newDiscount,
			"total":    newTotal,
		}).Error
	})
}

// AddTax adds a tax to an invoice item.
func (s *InvoiceItemService) AddTax(itemID uint, taxName string, taxRate float64, description string) error {
	item, err := s.GetByID(itemID)
	if err != nil {
		return err
	}

	taxAmount := (float64(item.Quantity)*item.UnitPrice - item.Discount) * taxRate

	tax := &InvoiceItemTax{
		ItemID:      itemID,
		TaxName:     taxName,
		TaxRate:     taxRate,
		TaxAmount:   taxAmount,
		Description: description,
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(tax).Error; err != nil {
			return err
		}
		newTax := item.Tax + taxAmount
		newTotal := float64(item.Quantity)*item.UnitPrice - item.Discount + newTax
		return tx.Model(item).Updates(map[string]interface{}{
			"tax":   newTax,
			"total": newTotal,
		}).Error
	})
}
