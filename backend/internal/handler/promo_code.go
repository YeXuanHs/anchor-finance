package handler

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// PromoCodeHandler 优惠码处理器
type PromoCodeHandler struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewPromoCodeHandler(db *gorm.DB, log *logger.Logger) *PromoCodeHandler {
	return &PromoCodeHandler{db: db, log: log}
}

// GetList 获取优惠码列表
func (h *PromoCodeHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keywords := c.Query("keywords")
	status := c.Query("status")

	query := h.db.Model(&model.PromoCode{})
	if keywords != "" {
		query = query.Where("code LIKE ?", "%"+keywords+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var list []model.PromoCode
	query.Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list)

	response.Success(c, gin.H{"list": list, "total": total})
}

// GetDetail 获取优惠码详情
func (h *PromoCodeHandler) GetDetail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var promo model.PromoCode
	if err := h.db.First(&promo, id).Error; err != nil {
		response.NotFound(c, "优惠码不存在")
		return
	}
	response.Success(c, promo)
}

// Create 创建优惠码
func (h *PromoCodeHandler) Create(c *gin.Context) {
	var req struct {
		Code           string   `json:"code" binding:"required"`
		Type           string   `json:"type" binding:"required,oneof=percent fixed override free"`
		Value          float64  `json:"value"`
		Cycles         []string `json:"cycles"`
		AppliesTo      []uint   `json:"applies_to"`
		Requires       []uint   `json:"requires"`
		Recurring      bool     `json:"recurring"`
		RecurFor       int      `json:"recur_for"`
		RequiresExist  bool     `json:"requires_exist"`
		MaxTimes       int      `json:"max_times"`
		Lifelong       bool     `json:"lifelong"`
		OneTime        bool     `json:"one_time"`
		OnlyNewClient  bool     `json:"only_new_client"`
		OnlyOldClient  bool     `json:"only_old_client"`
		OncePerClient  bool     `json:"once_per_client"`
		Upgrades       bool     `json:"upgrades"`
		IsDiscount     bool     `json:"is_discount"`
		Notes          string   `json:"notes"`
		StartTime      int64    `json:"start_time" binding:"required"`
		ExpirationTime int64    `json:"expiration_time" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}

	// 检查优惠码是否已存在
	var count int64
	h.db.Model(&model.PromoCode{}).Where("code = ?", req.Code).Count(&count)
	if count > 0 {
		response.BadRequest(c, "优惠码已存在")
		return
	}

	promo := model.PromoCode{
		Code:           req.Code,
		Type:           req.Type,
		Value:          req.Value,
		Cycles:         strings.Join(req.Cycles, ","),
		AppliesTo:      uintSliceToString(req.AppliesTo),
		Requires:       uintSliceToString(req.Requires),
		Recurring:      req.Recurring,
		RecurFor:       req.RecurFor,
		RequiresExist:  req.RequiresExist,
		MaxTimes:       req.MaxTimes,
		Lifelong:       req.Lifelong,
		OneTime:        req.OneTime,
		OnlyNewClient:  req.OnlyNewClient,
		OnlyOldClient:  req.OnlyOldClient,
		OncePerClient:  req.OncePerClient,
		Upgrades:       req.Upgrades,
		IsDiscount:     req.IsDiscount,
		Notes:          req.Notes,
		StartTime:      req.StartTime,
		ExpirationTime: req.ExpirationTime,
		Status:         1,
	}

	if err := h.db.Create(&promo).Error; err != nil {
		response.ServerError(c, "创建失败")
		return
	}
	response.Success(c, promo)
}

// Update 更新优惠码
func (h *PromoCodeHandler) Update(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var promo model.PromoCode
	if err := h.db.First(&promo, id).Error; err != nil {
		response.NotFound(c, "优惠码不存在")
		return
	}

	var req struct {
		Code           string   `json:"code"`
		Type           string   `json:"type"`
		Value          float64  `json:"value"`
		Cycles         []string `json:"cycles"`
		AppliesTo      []uint   `json:"applies_to"`
		Requires       []uint   `json:"requires"`
		Recurring      bool     `json:"recurring"`
		RecurFor       int      `json:"recur_for"`
		RequiresExist  bool     `json:"requires_exist"`
		MaxTimes       int      `json:"max_times"`
		Lifelong       bool     `json:"lifelong"`
		OneTime        bool     `json:"one_time"`
		OnlyNewClient  bool     `json:"only_new_client"`
		OnlyOldClient  bool     `json:"only_old_client"`
		OncePerClient  bool     `json:"once_per_client"`
		Upgrades       bool     `json:"upgrades"`
		IsDiscount     bool     `json:"is_discount"`
		Notes          string   `json:"notes"`
		StartTime      int64    `json:"start_time"`
		ExpirationTime int64    `json:"expiration_time"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}

	// 检查优惠码是否重复
	if req.Code != "" && req.Code != promo.Code {
		var count int64
		h.db.Model(&model.PromoCode{}).Where("code = ? AND id != ?", req.Code, id).Count(&count)
		if count > 0 {
			response.BadRequest(c, "优惠码已存在")
			return
		}
	}

	updates := map[string]interface{}{
		"type":            req.Type,
		"value":           req.Value,
		"cycles":          strings.Join(req.Cycles, ","),
		"applies_to":      uintSliceToString(req.AppliesTo),
		"requires":        uintSliceToString(req.Requires),
		"recurring":       req.Recurring,
		"recur_for":       req.RecurFor,
		"requires_exist":  req.RequiresExist,
		"max_times":       req.MaxTimes,
		"lifelong":        req.Lifelong,
		"one_time":        req.OneTime,
		"only_new_client": req.OnlyNewClient,
		"only_old_client": req.OnlyOldClient,
		"once_per_client": req.OncePerClient,
		"upgrades":        req.Upgrades,
		"is_discount":     req.IsDiscount,
		"notes":           req.Notes,
		"start_time":      req.StartTime,
		"expiration_time": req.ExpirationTime,
	}
	if req.Code != "" {
		updates["code"] = req.Code
	}

	h.db.Model(&promo).Updates(updates)
	response.Success(c, promo)
}

// Delete 删除优惠码
func (h *PromoCodeHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.db.Delete(&model.PromoCode{}, id).Error; err != nil {
		response.ServerError(c, "删除失败")
		return
	}
	response.Success(c, nil)
}

// SetStatus 启用/禁用优惠码
func (h *PromoCodeHandler) SetStatus(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	h.db.Model(&model.PromoCode{}).Where("id = ?", id).Update("status", req.Status)
	response.Success(c, nil)
}

// ExpireImmediately 立即过期
func (h *PromoCodeHandler) ExpireImmediately(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	h.db.Model(&model.PromoCode{}).Where("id = ?", id).Updates(map[string]interface{}{
		"expiration_time": time.Now().Unix(),
		"status":          0,
	})
	response.Success(c, nil)
}

// AutoGenerate 自动生成优惠码
func (h *PromoCodeHandler) AutoGenerate(c *gin.Context) {
	prefix := c.DefaultQuery("prefix", "AF")
	length, _ := strconv.Atoi(c.DefaultQuery("length", "8"))
	code := generatePromoCode(prefix, length)
	response.Success(c, gin.H{"code": code})
}

// GetUsageLogs 获取使用记录
func (h *PromoCodeHandler) GetUsageLogs(c *gin.Context) {
	promoID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var total int64
	h.db.Model(&model.PromoCodeLog{}).Where("promo_id = ?", promoID).Count(&total)

	var logs []model.PromoCodeLog
	h.db.Where("promo_id = ?", promoID).Order("id DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&logs)

	response.Success(c, gin.H{"list": logs, "total": total})
}

// Validate 验证优惠码（用户端使用）
func (h *PromoCodeHandler) Validate(c *gin.Context) {
	code := c.Query("code")
	productID, _ := strconv.ParseUint(c.Query("product_id"), 10, 32)

	var promo model.PromoCode
	if err := h.db.Where("code = ? AND status = 1", code).First(&promo).Error; err != nil {
		response.BadRequest(c, "优惠码无效")
		return
	}

	now := time.Now().Unix()
	if now < promo.StartTime {
		response.BadRequest(c, "优惠码尚未生效")
		return
	}
	if now > promo.ExpirationTime {
		response.BadRequest(c, "优惠码已过期")
		return
	}

	// 检查使用次数
	if promo.MaxTimes > 0 && promo.UsedCount >= promo.MaxTimes {
		response.BadRequest(c, "优惠码已达最大使用次数")
		return
	}

	// 检查产品适用性
	if promo.AppliesTo != "" {
		appliesTo := stringToUintSlice(promo.AppliesTo)
		found := false
		for _, id := range appliesTo {
			if id == uint(productID) {
				found = true
				break
			}
		}
		if !found {
			response.BadRequest(c, "优惠码不适用于此产品")
			return
		}
	}

	response.Success(c, gin.H{
		"valid":    true,
		"type":     promo.Type,
		"value":    promo.Value,
		"recurring": promo.Recurring,
		"recur_for": promo.RecurFor,
	})
}

// generatePromoCode 生成随机优惠码
func generatePromoCode(prefix string, length int) string {
	const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	code := make([]byte, length)
	for i := range code {
		code[i] = chars[r.Intn(len(chars))]
	}
	return prefix + string(code)
}

// uintSliceToString 将uint切片转为逗号分隔字符串
func uintSliceToString(slice []uint) string {
	if len(slice) == 0 {
		return ""
	}
	strs := make([]string, len(slice))
	for i, v := range slice {
		strs[i] = strconv.FormatUint(uint64(v), 10)
	}
	return strings.Join(strs, ",")
}

// stringToUintSlice 将逗号分隔字符串转为uint切片
func stringToUintSlice(s string) []uint {
	if s == "" {
		return nil
	}
	parts := strings.Split(s, ",")
	result := make([]uint, 0, len(parts))
	for _, p := range parts {
		if v, err := strconv.ParseUint(strings.TrimSpace(p), 10, 32); err == nil {
			result = append(result, uint(v))
		}
	}
	return result
}
