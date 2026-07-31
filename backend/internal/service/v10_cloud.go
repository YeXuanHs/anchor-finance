package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// V10CloudService handles cloud server product browsing, configuration, cart, order and host management.
type V10CloudService struct {
	db    *gorm.DB
	log   *logger.Logger
	oSvc  *OrderService
	cSvc  *CouponService
}

func NewV10CloudService(db *gorm.DB, log *logger.Logger, oSvc *OrderService, cSvc *CouponService) *V10CloudService {
	return &V10CloudService{db: db, log: log, oSvc: oSvc, cSvc: cSvc}
}

// ─────────────────────────── Request / Response DTOs ───────────────────────────

type CloudProductListReq struct {
	GroupID  uint `form:"group_id"`
	Page     int  `form:"page"`
	PageSize int  `form:"page_size"`
}

type ConfigFilter struct {
	CPU      *int    `json:"cpu"`
	MemoryMB *int    `json:"memory_mb"`
	DiskGB   *int    `json:"disk_gb"`
	Region   *string `json:"region"`
}

type CalculatePriceReq struct {
	ProductID uint           `json:"product_id" binding:"required"`
	Config    datatypes.JSON `json:"config"`
	Cycle     string         `json:"cycle" binding:"required"`
	Quantity  int            `json:"quantity"`
}

type CartItemReq struct {
	ProductID uint           `json:"product_id" binding:"required"`
	Config    datatypes.JSON `json:"config"`
	Cycle     string         `json:"cycle" binding:"required"`
	Quantity  int            `json:"quantity"`
}

type UpdateCartItemReq struct {
	Config   datatypes.JSON `json:"config"`
	Cycle    string         `json:"cycle"`
	Quantity int            `json:"quantity"`
}

type CreateOrderReq struct {
	CartItemIDs []uint `json:"cart_item_ids" binding:"required"`
	CouponCode  string `json:"coupon_code"`
}

type PayOrderReq struct {
	PaymentMethod string `json:"payment_method" binding:"required"`
}

type PriceBreakdown struct {
	BasePrice  float64            `json:"base_price"`
	Addons     map[string]float64 `json:"addons"`
	Discount   float64            `json:"discount"`
	Total      float64            `json:"total"`
	CycleLabel string             `json:"cycle_label"`
}

// ─────────────────────────── Product Browsing ───────────────────────────

// GetProductList lists available cloud products filtered by group.
func (s *V10CloudService) GetProductList(groupID uint, page, pageSize int) ([]model.V10CloudProduct, int64, error) {
	var products []model.V10CloudProduct
	var total int64

	q := s.db.Model(&model.V10CloudProduct{}).Where("enabled = ?", true)
	if groupID > 0 {
		q = q.Where("group_id = ?", groupID)
	}

	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	if err := q.Offset(offset).Limit(pageSize).Order("id ASC").Find(&products).Error; err != nil {
		return nil, 0, err
	}
	return products, total, nil
}

// GetProductDetail returns a single cloud product by ID.
func (s *V10CloudService) GetProductDetail(id uint) (*model.V10CloudProduct, error) {
	var p model.V10CloudProduct
	if err := s.db.First(&p, id).Error; err != nil {
		return nil, errors.New("product not found")
	}
	return &p, nil
}

// GetRegions lists all enabled regions.
func (s *V10CloudService) GetRegions() ([]model.V10CloudRegion, error) {
	var regions []model.V10CloudRegion
	if err := s.db.Where("enabled = ?", true).Order("sort_order ASC, id ASC").Find(&regions).Error; err != nil {
		return nil, err
	}
	return regions, nil
}

// GetOSTypes lists all enabled OS types.
func (s *V10CloudService) GetOSTypes() ([]model.V10CloudOSType, error) {
	var osTypes []model.V10CloudOSType
	if err := s.db.Where("enabled = ?", true).Order("sort_order ASC, id ASC").Find(&osTypes).Error; err != nil {
		return nil, err
	}
	return osTypes, nil
}

// ─────────────────────────── Configuration ───────────────────────────

// GetConfigOptions returns configurable options for a product.
func (s *V10CloudService) GetConfigOptions(productID uint) ([]model.V10CloudConfigOption, error) {
	var opts []model.V10CloudConfigOption
	if err := s.db.Where("product_id = ? AND enabled = ?", productID, true).
		Order("sort_order ASC, id ASC").Find(&opts).Error; err != nil {
		return nil, err
	}
	return opts, nil
}

// CalculatePrice computes the total price for a given product + config + cycle + quantity.
func (s *V10CloudService) CalculatePrice(req CalculatePriceReq) (*PriceBreakdown, error) {
	product, err := s.GetProductDetail(req.ProductID)
	if err != nil {
		return nil, err
	}

	if req.Quantity < 1 {
		req.Quantity = 1
	}

	basePrice := product.Price
	cycleMultiplier := cycleToMultiplier(req.Cycle)
	basePrice *= cycleMultiplier

	breakdown := &PriceBreakdown{
		BasePrice:  basePrice,
		Addons:     make(map[string]float64),
		CycleLabel: cycleLabel(req.Cycle),
	}

	// Apply config addon pricing
	if req.Config != nil {
		var configMap map[string]interface{}
		if err := json.Unmarshal(req.Config, &configMap); err == nil {
			// Look up config options for addon pricing
			var opts []model.V10CloudConfigOption
			s.db.Where("product_id = ? AND enabled = ?", req.ProductID, true).Find(&opts)

			optMap := make(map[string]model.V10CloudConfigOption)
			for _, o := range opts {
				optMap[o.Value] = o
			}

			for key, val := range configMap {
				valStr := fmt.Sprintf("%v", val)
				if opt, ok := optMap[valStr]; ok && opt.Price > 0 {
					addonPrice := opt.Price * cycleMultiplier
					breakdown.Addons[key] = addonPrice
					breakdown.Total += addonPrice
				}
			}
		}
	}

	breakdown.Total += basePrice
	breakdown.Total *= float64(req.Quantity)

	return breakdown, nil
}

// GetLinkAgeList returns cascading config options based on a parent selection.
func (s *V10CloudService) GetLinkAgeList(productID uint, parentID *uint) ([]model.V10CloudConfigOption, error) {
	var opts []model.V10CloudConfigOption
	q := s.db.Where("product_id = ? AND enabled = ?", productID, true)
	if parentID != nil {
		q = q.Where("parent_id = ?", *parentID)
	} else {
		q = q.Where("parent_id IS NULL")
	}
	if err := q.Order("sort_order ASC, id ASC").Find(&opts).Error; err != nil {
		return nil, err
	}
	return opts, nil
}

// FilterConfigOptions filters available options by given criteria.
func (s *V10CloudService) FilterConfigOptions(productID uint, filters ConfigFilter) ([]model.V10CloudConfigOption, error) {
	var opts []model.V10CloudConfigOption
	q := s.db.Where("product_id = ? AND enabled = ?", productID, true)

	if filters.CPU != nil {
		q = q.Where("type = ? AND value = ?", "cpu", fmt.Sprintf("%d", *filters.CPU))
	}
	if filters.MemoryMB != nil {
		q = q.Where("type = ? AND value = ?", "memory", fmt.Sprintf("%d", *filters.MemoryMB))
	}
	if filters.DiskGB != nil {
		q = q.Where("type = ? AND value = ?", "disk", fmt.Sprintf("%d", *filters.DiskGB))
	}
	if filters.Region != nil {
		q = q.Where("type = ? AND value = ?", "region", *filters.Region)
	}

	if err := q.Order("sort_order ASC, id ASC").Find(&opts).Error; err != nil {
		return nil, err
	}
	return opts, nil
}

// ─────────────────────────── Cart Operations ───────────────────────────

// AddToCart adds a configured cloud product to the user's cart.
func (s *V10CloudService) AddToCart(userID, productID uint, config datatypes.JSON, cycle string, qty int) (*model.V10CloudOrder, error) {
	if qty < 1 {
		qty = 1
	}

	// Verify product exists and is enabled
	product, err := s.GetProductDetail(productID)
	if err != nil {
		return nil, err
	}
	if !product.Enabled {
		return nil, errors.New("product is disabled")
	}

	// Check stock
	if product.Stock >= 0 && product.Stock < qty {
		return nil, errors.New("insufficient stock")
	}

	// Check for existing cart item with same product + cycle + config
	var existing model.V10CloudOrder
	err = s.db.Where("user_id = ? AND product_id = ? AND cycle = ? AND status = ?",
		userID, productID, cycle, "cart").First(&existing).Error

	if err == nil {
		existing.Quantity += qty
		if config != nil {
			existing.Config = config
		}
		if err := s.db.Save(&existing).Error; err != nil {
			return nil, err
		}
		return &existing, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Calculate unit price
	priceReq := CalculatePriceReq{
		ProductID: productID,
		Config:    config,
		Cycle:     cycle,
		Quantity:  1,
	}
	breakdown, err := s.CalculatePrice(priceReq)
	if err != nil {
		return nil, err
	}

	item := &model.V10CloudOrder{
		UserID:     userID,
		ProductID:  productID,
		Cycle:      cycle,
		Quantity:   qty,
		Config:     config,
		TotalPrice: breakdown.Total * float64(qty),
		Status:     "cart",
	}

	if err := s.db.Create(item).Error; err != nil {
		return nil, err
	}
	s.log.Infof("v10 cloud cart item added: user=%d product=%d cycle=%s", userID, productID, cycle)
	return item, nil
}

// UpdateCartItem updates a cart item's config, cycle, or quantity.
func (s *V10CloudService) UpdateCartItem(userID, itemID uint, req UpdateCartItemReq) (*model.V10CloudOrder, error) {
	var item model.V10CloudOrder
	if err := s.db.Where("id = ? AND user_id = ? AND status = ?", itemID, userID, "cart").First(&item).Error; err != nil {
		return nil, errors.New("cart item not found")
	}

	if req.Quantity > 0 {
		item.Quantity = req.Quantity
	}
	if req.Config != nil {
		item.Config = req.Config
	}
	if req.Cycle != "" {
		item.Cycle = req.Cycle
	}

	// Recalculate price
	priceReq := CalculatePriceReq{
		ProductID: item.ProductID,
		Config:    item.Config,
		Cycle:     item.Cycle,
		Quantity:  1,
	}
	breakdown, err := s.CalculatePrice(priceReq)
	if err != nil {
		return nil, err
	}
	item.TotalPrice = breakdown.Total * float64(item.Quantity)

	if err := s.db.Save(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

// GetCartSummary returns cart items with totals.
func (s *V10CloudService) GetCartSummary(userID uint) (map[string]interface{}, error) {
	items, err := s.GetCartItems(userID)
	if err != nil {
		return nil, err
	}

	var subTotal float64
	for _, item := range items {
		subTotal += item.TotalPrice
	}

	summary := map[string]interface{}{
		"items":      items,
		"item_count": len(items),
		"sub_total":  subTotal,
		"total":      subTotal,
		"currency":   "CNY",
	}
	return summary, nil
}

// GetCartItems returns all cart items for a user.
func (s *V10CloudService) GetCartItems(userID uint) ([]model.V10CloudOrder, error) {
	var items []model.V10CloudOrder
	if err := s.db.Where("user_id = ? AND status = ?", userID, "cart").
		Order("created_at DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

// SettleCart prepares cart for checkout by validating all items.
func (s *V10CloudService) SettleCart(userID uint) ([]model.V10CloudOrder, error) {
	items, err := s.GetCartItems(userID)
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New("cart is empty")
	}

	for i, item := range items {
		// Verify product still exists and is enabled
		product, err := s.GetProductDetail(item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("product %d not found", item.ProductID)
		}
		if !product.Enabled {
			return nil, fmt.Errorf("product %s is disabled", product.Name)
		}
		// Verify stock
		if product.Stock >= 0 && product.Stock < item.Quantity {
			return nil, fmt.Errorf("product %s has insufficient stock", product.Name)
		}

		// Recalculate price to catch any changes
		priceReq := CalculatePriceReq{
			ProductID: item.ProductID,
			Config:    item.Config,
			Cycle:     item.Cycle,
			Quantity:  1,
		}
		breakdown, err := s.CalculatePrice(priceReq)
		if err != nil {
			return nil, err
		}
		items[i].TotalPrice = breakdown.Total * float64(item.Quantity)
		if err := s.db.Model(&items[i]).Update("total_price", items[i].TotalPrice).Error; err != nil {
			return nil, err
		}
	}

	return items, nil
}

// ─────────────────────────── Order Flow ───────────────────────────

// CreateOrder creates orders from cart items, optionally applying a coupon.
func (s *V10CloudService) CreateOrder(userID uint, cartItemIDs []uint, couponCode string) ([]model.V10CloudOrder, error) {
	// Fetch cart items
	var cartItems []model.V10CloudOrder
	if err := s.db.Where("id IN ? AND user_id = ? AND status = ?", cartItemIDs, userID, "cart").
		Find(&cartItems).Error; err != nil {
		return nil, err
	}
	if len(cartItems) == 0 {
		return nil, errors.New("no valid cart items found")
	}

	var orders []model.V10CloudOrder
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range cartItems {
			totalPrice := item.TotalPrice

			// Apply coupon if provided
			if couponCode != "" && s.cSvc != nil {
				discount, _, err := s.cSvc.Validate(couponCode, userID, item.ProductID, totalPrice)
				if err == nil {
					totalPrice -= discount
					if totalPrice < 0 {
						totalPrice = 0
					}
				}
			}

			order := model.V10CloudOrder{
				UserID:     userID,
				ProductID:  item.ProductID,
				Cycle:      item.Cycle,
				Quantity:   item.Quantity,
				Config:     item.Config,
				TotalPrice: totalPrice,
				Status:     "pending",
			}

			if err := tx.Create(&order).Error; err != nil {
				return err
			}

			// Update stock
			var product model.V10CloudProduct
			if err := tx.First(&product, item.ProductID).Error; err == nil && product.Stock >= 0 {
				if err := tx.Model(&product).Update("stock", gorm.Expr("stock - ?", item.Quantity)).Error; err != nil {
					return err
				}
			}

			orders = append(orders, order)
			s.log.Infof("v10 cloud order created: user=%d product=%d total=%.2f", userID, item.ProductID, totalPrice)
		}

		// Remove cart items
		if err := tx.Where("id IN ?", cartItemIDs).Delete(&model.V10CloudOrder{}).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return orders, nil
}

// GetOrderDetail returns a single order by ID.
func (s *V10CloudService) GetOrderDetail(orderID uint) (*model.V10CloudOrder, error) {
	var order model.V10CloudOrder
	if err := s.db.First(&order, orderID).Error; err != nil {
		return nil, errors.New("order not found")
	}
	return &order, nil
}

// PayOrder processes payment for an order.
func (s *V10CloudService) PayOrder(orderID uint, paymentMethod string) (*model.V10CloudOrder, error) {
	var order model.V10CloudOrder
	if err := s.db.First(&order, orderID).Error; err != nil {
		return nil, errors.New("order not found")
	}

	if order.Status != "pending" {
		return nil, fmt.Errorf("order is in %s status, cannot pay", order.Status)
	}

	order.Status = "paid"
	if err := s.db.Save(&order).Error; err != nil {
		return nil, err
	}

	s.log.Infof("v10 cloud order paid: id=%d method=%s", orderID, paymentMethod)
	return &order, nil
}

// ─────────────────────────── Host Management ───────────────────────────

// GetHostInfo returns host information.
func (s *V10CloudService) GetHostInfo(hostID uint) (map[string]interface{}, error) {
	var host struct {
		ID            uint            `json:"id"`
		Hostname      string          `json:"hostname"`
		IP            string          `json:"ip"`
		IPv6          string          `json:"ipv6"`
		OS            string          `json:"os"`
		CPU           string          `json:"cpu"`
		CPUCores      int             `json:"cpu_cores"`
		MemoryMB      int             `json:"memory_mb"`
		DiskSizeGB    int             `json:"disk_size_gb"`
		BandwidthMbps int             `json:"bandwidth_mbps"`
		TrafficGB     int             `json:"traffic_gb"`
		Location      string          `json:"location"`
		Status        int16           `json:"status"`
		ExpiredAt     *time.Time      `json:"expired_at"`
		Config        datatypes.JSON  `json:"config"`
	}

	if err := s.db.Table("hosts").Where("id = ?", hostID).First(&host).Error; err != nil {
		return nil, errors.New("host not found")
	}

	result := map[string]interface{}{
		"id":             host.ID,
		"hostname":       host.Hostname,
		"ip":             host.IP,
		"ipv6":           host.IPv6,
		"os":             host.OS,
		"cpu":            host.CPU,
		"cpu_cores":      host.CPUCores,
		"memory_mb":      host.MemoryMB,
		"disk_size_gb":   host.DiskSizeGB,
		"bandwidth_mbps": host.BandwidthMbps,
		"traffic_gb":     host.TrafficGB,
		"location":       host.Location,
		"status":         host.Status,
		"expired_at":     host.ExpiredAt,
		"config":         host.Config,
	}
	return result, nil
}

// GetHostConfig returns the host configuration.
func (s *V10CloudService) GetHostConfig(hostID uint) (datatypes.JSON, error) {
	var config datatypes.JSON
	if err := s.db.Table("hosts").Where("id = ?", hostID).Pluck("config", &config).Error; err != nil {
		return nil, errors.New("host not found")
	}
	return config, nil
}

// GetTrafficUsage returns traffic usage stats for a host.
func (s *V10CloudService) GetTrafficUsage(hostID uint) (map[string]interface{}, error) {
	// Verify host exists
	var count int64
	s.db.Table("hosts").Where("id = ?", hostID).Count(&count)
	if count == 0 {
		return nil, errors.New("host not found")
	}

	// Return traffic usage info from metadata
	var metadata datatypes.JSON
	if err := s.db.Table("hosts").Where("id = ?", hostID).Pluck("metadata", &metadata).Error; err != nil {
		return nil, err
	}

	var metaMap map[string]interface{}
	usedGB := 0.0
	totalGB := 0.0
	if metadata != nil {
		if err := json.Unmarshal(metadata, &metaMap); err == nil {
			if v, ok := metaMap["traffic_used_gb"]; ok {
				usedGB, _ = v.(float64)
			}
			if v, ok := metaMap["traffic_total_gb"]; ok {
				totalGB, _ = v.(float64)
			}
		}
	}

	result := map[string]interface{}{
		"host_id":      hostID,
		"used_gb":      usedGB,
		"total_gb":     totalGB,
		"reset_date":   time.Now().AddDate(0, 1, 0).Format("2006-01-02"),
		"billing_type": "monthly",
	}
	return result, nil
}

// GetOSList returns available OS options for reinstall.
func (s *V10CloudService) GetOSList(hostID uint) ([]model.V10CloudOSType, error) {
	// Verify host exists
	var count int64
	s.db.Table("hosts").Where("id = ?", hostID).Count(&count)
	if count == 0 {
		return nil, errors.New("host not found")
	}

	var osTypes []model.V10CloudOSType
	if err := s.db.Where("enabled = ?", true).Order("sort_order ASC, id ASC").Find(&osTypes).Error; err != nil {
		return nil, err
	}
	return osTypes, nil
}

// ─────────────────────────── Helpers ───────────────────────────

func cycleToMultiplier(cycle string) float64 {
	switch cycle {
	case "monthly":
		return 1
	case "quarterly":
		return 3
	case "semi-annually":
		return 6
	case "annually":
		return 12
	case "biennially":
		return 24
	case "triennially":
		return 36
	default:
		return 1
	}
}

func cycleLabel(cycle string) string {
	switch cycle {
	case "monthly":
		return "月付"
	case "quarterly":
		return "季付"
	case "semi-annually":
		return "半年付"
	case "annually":
		return "年付"
	case "biennially":
		return "两年付"
	case "triennially":
		return "三年付"
	default:
		return cycle
	}
}

// GetDB returns the underlying gorm.DB.
func (s *V10CloudService) GetDB() *gorm.DB {
	return s.db
}
