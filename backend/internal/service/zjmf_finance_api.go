package service

import (
	"crypto/md5"
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// ZJMFFinanceAPIService 魔方兼容API服务
type ZJMFFinanceAPIService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewZJMFFinanceAPIService(db *gorm.DB, log *logger.Logger) *ZJMFFinanceAPIService {
	return &ZJMFFinanceAPIService{db: db, log: log}
}

// AuthRequest 认证请求
type AuthRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse 认证响应
type AuthResponse struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	UserID    uint      `json:"user_id"`
}

// ProductResponse 产品响应
type ProductResponse struct {
	ID       uint    `json:"id"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
	Stock    int     `json:"stock"`
}

// OrderRequest 创建订单请求
type OrderRequest struct {
	UserID    uint `json:"user_id" binding:"required"`
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity"`
}

// Authenticate 用户认证
func (s *ZJMFFinanceAPIService) Authenticate(req AuthRequest) (*AuthResponse, error) {
	var user model.User
	if err := s.db.Where("email = ?", req.Email).First(&user).Error; err != nil {
		return nil, errors.New("invalid credentials")
	}

	hash := fmt.Sprintf("%x", md5.Sum([]byte(req.Password+user.Salt)))
	if hash != user.PasswordHash {
		return nil, errors.New("invalid credentials")
	}

	if user.Status != 1 {
		return nil, errors.New("account disabled")
	}

	expiresAt := time.Now().Add(24 * time.Hour)
	token := fmt.Sprintf("%x", md5.Sum([]byte(fmt.Sprintf("%d:%d:%s", user.ID, time.Now().UnixNano(), user.Email))))

	s.db.Model(&user).Updates(map[string]interface{}{
		"last_login_at": time.Now(),
	})

	s.log.Infof("zjmf api auth: user=%d", user.ID)
	return &AuthResponse{
		Token:     token,
		ExpiresAt: expiresAt,
		UserID:    user.ID,
	}, nil
}

// GetProducts 获取产品列表
func (s *ZJMFFinanceAPIService) GetProducts(page, pageSize int) ([]ProductResponse, int64, error) {
	var products []model.Product
	var total int64

	query := s.db.Model(&model.Product{}).Where("status = ?", 1)
	query.Count(&total)
	offset, limit := util.Paginate(page, pageSize)
	query.Offset(offset).Limit(limit).Order("id ASC").Find(&products)

	result := make([]ProductResponse, len(products))
	for i, p := range products {
		result[i] = ProductResponse{
			ID:       p.ID,
			Name:     p.Name,
			Price:    p.Price,
			Currency: p.Currency,
			Stock:    p.Stock,
		}
	}
	return result, total, nil
}

// CreateOrder 创建订单
func (s *ZJMFFinanceAPIService) CreateOrder(req OrderRequest) (*model.Order, error) {
	var user model.User
	if err := s.db.First(&user, req.UserID).Error; err != nil {
		return nil, errors.New("user not found")
	}

	var product model.Product
	if err := s.db.First(&product, req.ProductID).Error; err != nil {
		return nil, errors.New("product not found")
	}

	quantity := req.Quantity
	if quantity <= 0 {
		quantity = 1
	}

	price, _ := product.Price.Float64()
	totalAmount := price * float64(quantity)

	order := model.Order{
		OrderNo:  util.GenerateOrderNo(),
		UserID:   req.UserID,
		Amount:   product.Price,
		Total:    totalAmount,
		Status:   0,
		Quantity: quantity,
		Currency: product.Currency,
	}
	if err := s.db.Create(&order).Error; err != nil {
		return nil, err
	}

	s.log.Infof("zjmf api order created: id=%d user=%d total=%.2f", order.ID, req.UserID, totalAmount)
	return &order, nil
}

// GetInvoice 获取发票
func (s *ZJMFFinanceAPIService) GetInvoice(id uint) (*model.Invoice, error) {
	var invoice model.Invoice
	if err := s.db.First(&invoice, id).Error; err != nil {
		return nil, err
	}
	return &invoice, nil
}

// PayInvoice 支付发票
func (s *ZJMFFinanceAPIService) PayInvoice(id uint) (*model.Invoice, error) {
	var invoice model.Invoice
	if err := s.db.First(&invoice, id).Error; err != nil {
		return nil, err
	}
	if invoice.Status == 1 {
		return nil, errors.New("invoice already paid")
	}

	now := time.Now()
	updates := map[string]interface{}{
		"status":  1,
		"paid_at": &now,
	}
	if err := s.db.Model(&invoice).Updates(updates).Error; err != nil {
		return nil, err
	}
	s.db.First(&invoice, id)
	s.log.Infof("zjmf api invoice paid: id=%d", id)
	return &invoice, nil
}
