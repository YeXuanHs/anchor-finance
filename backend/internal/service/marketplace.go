package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// MarketplaceService 交易市场服务
type MarketplaceService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewMarketplaceService(db *gorm.DB, log *logger.Logger) *MarketplaceService {
	return &MarketplaceService{db: db, log: log}
}

// GetDB 获取数据库连接
func (s *MarketplaceService) GetDB() *gorm.DB {
	return s.db
}

// ─── 配置管理 ───

// GetConfig 获取市场配置
func (s *MarketplaceService) GetConfig() *model.MarketplaceConfig {
	var config model.MarketplaceConfig
	if err := s.db.First(&config).Error; err != nil {
		config = model.MarketplaceConfig{
			ID:              1,
			Enabled:         true,
			FeeRate:         5,
			MinFee:          1,
			MaxListingDays:  30,
			MinHoldDays:     7,
			RequireRealName: false,
			AllowFeeOnly:    true,
			AutoTransfer:    true,
			NotifyEmail:     true,
		}
		s.db.Create(&config)
	}
	return &config
}

// SaveConfig 保存市场配置
func (s *MarketplaceService) SaveConfig(config *model.MarketplaceConfig) error {
	config.ID = 1
	return s.db.Save(config).Error
}

// ─── 挂售管理 ───

// CreateListing 创建挂售
func (s *MarketplaceService) CreateListing(userID uint, hostID uint, sellPrice float64, description string) (*model.MarketplaceListing, error) {
	config := s.GetConfig()
	if !config.Enabled {
		return nil, errors.New("交易市场已关闭")
	}

	// 检查主机是否属于用户
	var host model.Host
	if err := s.db.Where("id = ? AND user_id = ?", hostID, userID).First(&host).Error; err != nil {
		return nil, errors.New("主机不存在或不属于您")
	}

	// 检查持有天数
	holdDays := int(time.Since(host.CreatedAt).Hours() / 24)
	if holdDays < config.MinHoldDays {
		return nil, fmt.Errorf("持有天数不足，需要%d天，当前%d天", config.MinHoldDays, holdDays)
	}

	// 检查是否已挂售
	var existing model.MarketplaceListing
	if err := s.db.Where("host_id = ? AND status = 1", hostID).First(&existing).Error; err == nil {
		return nil, errors.New("该主机已在挂售中")
	}

	// 获取产品信息
	var product model.Product
	s.db.First(&product, host.ProductID)

	// 计算到期时间
	var expiresAt *time.Time
	if host.ExpiresAt != nil {
		expiresAt = host.ExpiresAt
	}

	listing := &model.MarketplaceListing{
		UserID:        userID,
		HostID:        hostID,
		ProductName:   product.Name,
		Description:   description,
		OriginalPrice: product.Price,
		SellPrice:     sellPrice,
		Currency:      "CNY",
		Status:        1,
		ExpiresAt:     expiresAt,
	}

	if err := s.db.Create(listing).Error; err != nil {
		return nil, err
	}

	return listing, nil
}

// UpdateListing 更新挂售
func (s *MarketplaceService) UpdateListing(userID, listingID uint, sellPrice float64, description string) error {
	var listing model.MarketplaceListing
	if err := s.db.Where("id = ? AND user_id = ?", listingID, userID).First(&listing).Error; err != nil {
		return errors.New("挂售不存在")
	}

	if listing.Status != 1 {
		return errors.New("只能修改在售状态的挂售")
	}

	return s.db.Model(&listing).Updates(map[string]interface{}{
		"sell_price":  sellPrice,
		"description": description,
	}).Error
}

// RemoveListing 下架挂售
func (s *MarketplaceService) RemoveListing(userID, listingID uint) error {
	var listing model.MarketplaceListing
	if err := s.db.Where("id = ? AND user_id = ?", listingID, userID).First(&listing).Error; err != nil {
		return errors.New("挂售不存在")
	}

	if listing.Status != 1 {
		return errors.New("只能下架在售状态的挂售")
	}

	return s.db.Model(&listing).Update("status", 3).Error
}

// GetListing 获取挂售详情
func (s *MarketplaceService) GetListing(listingID uint) (*model.MarketplaceListing, error) {
	var listing model.MarketplaceListing
	if err := s.db.Preload("User").Preload("Host").Where("id = ?", listingID).First(&listing).Error; err != nil {
		return nil, errors.New("挂售不存在")
	}

	// 增加浏览量
	s.db.Model(&listing).Update("view_count", gorm.Expr("view_count + 1"))

	return &listing, nil
}

// GetListings 获取挂售列表
func (s *MarketplaceService) GetListings(page, pageSize int, keyword string) ([]model.MarketplaceListing, int64, error) {
	var listings []model.MarketplaceListing
	var total int64

	query := s.db.Model(&model.MarketplaceListing{}).Where("status = 1")

	if keyword != "" {
		query = query.Where("product_name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").
		Preload("User").Preload("Host").Find(&listings).Error; err != nil {
		return nil, 0, err
	}

	return listings, total, nil
}

// GetUserListings 获取用户的挂售列表
func (s *MarketplaceService) GetUserListings(userID uint) ([]model.MarketplaceListing, error) {
	var listings []model.MarketplaceListing
	if err := s.db.Where("user_id = ?", userID).Order("id DESC").
		Preload("Host").Find(&listings).Error; err != nil {
		return nil, err
	}
	return listings, nil
}

// ─── 订单管理 ───

// CreateOrder 创建订单
func (s *MarketplaceService) CreateOrder(buyerID uint, listingID uint, paymentMethod string) (*model.MarketplaceOrder, error) {
	config := s.GetConfig()
	if !config.Enabled {
		return nil, errors.New("交易市场已关闭")
	}

	// 获取挂售信息
	var listing model.MarketplaceListing
	if err := s.db.Preload("User").Where("id = ? AND status = 1", listingID).First(&listing).Error; err != nil {
		return nil, errors.New("商品不存在或已下架")
	}

	if listing.UserID == buyerID {
		return nil, errors.New("不能购买自己的商品")
	}

	// 检查买家实名要求
	if config.RequireRealName {
		var buyer model.User
		s.db.First(&buyer, buyerID)
		if buyer.RealName == "" {
			return nil, errors.New("需要完成实名认证才能购买")
		}
	}

	// 计算费用
	fee := listing.SellPrice * config.FeeRate / 100
	if fee < config.MinFee {
		fee = config.MinFee
	}

	var amount, totalAmount float64
	if paymentMethod == "full" {
		// 全额购买
		amount = listing.SellPrice
		totalAmount = amount + fee
	} else if paymentMethod == "fee_only" && config.AllowFeeOnly {
		// 仅付手续费
		amount = 0
		totalAmount = fee
	} else {
		return nil, errors.New("无效的支付方式")
	}

	// 检查买家余额
	var buyer model.User
	if err := s.db.First(&buyer, buyerID).Error; err != nil {
		return nil, errors.New("用户不存在")
	}
	if buyer.Balance < totalAmount {
		return nil, errors.New("余额不足")
	}

	order := &model.MarketplaceOrder{
		OrderNo:       fmt.Sprintf("MKT%s%04d", time.Now().Format("20060102150405"), time.Now().UnixNano()%10000),
		ListingID:     listingID,
		BuyerID:       buyerID,
		SellerID:      listing.UserID,
		HostID:        listing.HostID,
		Amount:        amount,
		Fee:           fee,
		TotalAmount:   totalAmount,
		PaymentMethod: paymentMethod,
		Status:        0,
	}

	// 创建订单并扣减余额
	err := s.db.Transaction(func(tx *gorm.DB) error {
		// 扣减余额
		result := tx.Model(&model.User{}).Where("id = ? AND balance >= ?", buyerID, totalAmount).
			Update("balance", gorm.Expr("balance - ?", totalAmount))
		if result.RowsAffected == 0 {
			return errors.New("余额不足")
		}

		// 创建订单
		if err := tx.Create(order).Error; err != nil {
			return err
		}

		// 更新挂售状态
		tx.Model(&listing).Update("status", 2) // 已售

		return nil
	})

	if err != nil {
		return nil, err
	}

	// 自动转移
	if config.AutoTransfer {
		go s.processTransfer(order.ID)
	}

	return order, nil
}

// processTransfer 处理转移
func (s *MarketplaceService) processTransfer(orderID uint) {
	var order model.MarketplaceOrder
	if err := s.db.First(&order, orderID).Error; err != nil {
		return
	}

	if order.Status != 1 { // 未支付
		return
	}

	// 更新转移状态为转移中
	s.db.Model(&order).Update("transfer_status", 1)

	// 执行转移
	err := s.db.Table("hosts").Where("id = ?", order.HostID).
		Update("user_id", order.BuyerID).Error

	if err != nil {
		s.db.Model(&order).Updates(map[string]interface{}{
			"transfer_status": 3, // 转移失败
		})
		s.log.Errorf("转移失败: %v", err)
		return
	}

	// 转移成功
	now := time.Now()
	s.db.Model(&order).Updates(map[string]interface{}{
		"status":          2, // 已转移
		"transfer_status": 2, // 转移成功
	})

	// 创建产品转移记录
	transfer := &model.ProductTransfer{
		TransferNo:    fmt.Sprintf("MKTT%s%04d", time.Now().Format("20060102150405"), time.Now().UnixNano()%10000),
		FromUserID:    order.SellerID,
		ToUserID:      order.BuyerID,
		UserProductID: order.HostID,
		Reason:        "交易市场自动转移",
		Status:        1,
		ProcessedAt:   &now,
	}
	s.db.Create(transfer)
	s.db.Model(&order).Update("transfer_id", transfer.ID)
}

// PayOrder 支付订单（手动确认）
func (s *MarketplaceService) PayOrder(userID, orderID uint) error {
	var order model.MarketplaceOrder
	if err := s.db.Where("id = ? AND buyer_id = ?", orderID, userID).First(&order).Error; err != nil {
		return errors.New("订单不存在")
	}

	if order.Status != 0 {
		return errors.New("订单状态无效")
	}

	now := time.Now()
	return s.db.Model(&order).Updates(map[string]interface{}{
		"status": 1,
		"paid_at": &now,
	}).Error
}

// CompleteOrder 完成订单
func (s *MarketplaceService) CompleteOrder(userID, orderID uint) error {
	var order model.MarketplaceOrder
	if err := s.db.First(&order, orderID).Error; err != nil {
		return errors.New("订单不存在")
	}

	// 只有卖家可以确认完成
	if order.SellerID != userID {
		return errors.New("无权限")
	}

	if order.Status != 2 {
		return errors.New("订单状态无效")
	}

	now := time.Now()
	return s.db.Model(&order).Updates(map[string]interface{}{
		"status":        3,
		"completed_at":  &now,
	}).Error
}

// CancelOrder 取消订单
func (s *MarketplaceService) CancelOrder(userID, orderID uint) error {
	var order model.MarketplaceOrder
	if err := s.db.First(&order, orderID).Error; err != nil {
		return errors.New("订单不存在")
	}

	if order.BuyerID != userID && order.SellerID != userID {
		return errors.New("无权限")
	}

	if order.Status != 0 {
		return errors.New("只能取消待支付订单")
	}

	// 退款
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 退还余额
		tx.Model(&model.User{}).Where("id = ?", order.BuyerID).
			Update("balance", gorm.Expr("balance + ?", order.TotalAmount))

		// 恢复挂售状态
		tx.Model(&model.MarketplaceListing{}).Where("id = ?", order.ListingID).
			Update("status", 1)

		// 更新订单状态
		return tx.Model(&order).Update("status", 4).Error
	})
}

// GetBuyerOrders 获取买家订单
func (s *MarketplaceService) GetBuyerOrders(userID uint, page, pageSize int) ([]model.MarketplaceOrder, int64, error) {
	var orders []model.MarketplaceOrder
	var total int64

	query := s.db.Model(&model.MarketplaceOrder{}).Where("buyer_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").
		Preload("Listing").Preload("Seller").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// GetSellerOrders 获取卖家订单
func (s *MarketplaceService) GetSellerOrders(userID uint, page, pageSize int) ([]model.MarketplaceOrder, int64, error) {
	var orders []model.MarketplaceOrder
	var total int64

	query := s.db.Model(&model.MarketplaceOrder{}).Where("seller_id = ?", userID)
	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").
		Preload("Listing").Preload("Buyer").Find(&orders).Error; err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

// ─── 私聊功能 ───

// SendMessage 发送私聊消息
func (s *MarketplaceService) SendMessage(senderID, receiverID, listingID uint, content string) (*model.MarketplaceChat, error) {
	if senderID == receiverID {
		return nil, errors.New("不能给自己发消息")
	}

	message := &model.MarketplaceChat{
		ListingID:  listingID,
		SenderID:   senderID,
		ReceiverID: receiverID,
		Content:    content,
		IsRead:     false,
	}

	if err := s.db.Create(message).Error; err != nil {
		return nil, err
	}

	// 更新或创建会话
	s.updateChatSession(listingID, senderID, receiverID, content)

	return message, nil
}

// updateChatSession 更新聊天会话
func (s *MarketplaceService) updateChatSession(listingID, senderID, receiverID uint, content string) {
	user1, user2 := senderID, receiverID
	if user1 > user2 {
		user1, user2 = user2, user1
	}

	var session model.MarketplaceChatSession
	err := s.db.Where("listing_id = ? AND user1_id = ? AND user2_id = ?", listingID, user1, user2).
		First(&session).Error

	if err == gorm.ErrRecordNotFound {
		session = model.MarketplaceChatSession{
			ListingID:     listingID,
			User1ID:       user1,
			User2ID:       user2,
			LastMessage:   content,
			LastMessageAt: time.Now(),
			UnreadCount:   1,
		}
		s.db.Create(&session)
	} else if err == nil {
		s.db.Model(&session).Updates(map[string]interface{}{
			"last_message":    content,
			"last_message_at": time.Now(),
			"unread_count":    gorm.Expr("unread_count + 1"),
		})
	}
}

// GetChatMessages 获取聊天消息
func (s *MarketplaceService) GetChatMessages(userID, listingID, otherUserID uint, page, pageSize int) ([]model.MarketplaceChat, int64, error) {
	var messages []model.MarketplaceChat
	var total int64

	query := s.db.Model(&model.MarketplaceChat{}).
		Where("listing_id = ? AND ((sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?))",
			listingID, userID, otherUserID, otherUserID, userID)

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id ASC").
		Preload("Sender").Preload("Receiver").Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	// 标记为已读
	s.db.Model(&model.MarketplaceChat{}).
		Where("listing_id = ? AND sender_id = ? AND receiver_id = ? AND is_read = false",
			listingID, otherUserID, userID).
		Update("is_read", true)

	return messages, total, nil
}

// GetChatSessions 获取聊天会话列表
func (s *MarketplaceService) GetChatSessions(userID uint) ([]model.MarketplaceChatSession, error) {
	var sessions []model.MarketplaceChatSession
	if err := s.db.Where("user1_id = ? OR user2_id = ?", userID, userID).
		Order("last_message_at DESC").
		Preload("Listing").Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// GetUnreadCount 获取未读消息数
func (s *MarketplaceService) GetUnreadCount(userID uint) int64 {
	var count int64
	s.db.Model(&model.MarketplaceChat{}).
		Where("receiver_id = ? AND is_read = false", userID).
		Count(&count)
	return count
}
