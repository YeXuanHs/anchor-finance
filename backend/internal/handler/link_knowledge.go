package handler

import (
	"strconv"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// LinkKnowledgeHandler handles link knowledge base HTTP requests.
type LinkKnowledgeHandler struct {
	db  *gorm.DB
	log *logger.Logger
}

// NewLinkKnowledgeHandler creates a new LinkKnowledgeHandler.
func NewLinkKnowledgeHandler(db *gorm.DB, log *logger.Logger) *LinkKnowledgeHandler {
	return &LinkKnowledgeHandler{db: db, log: log}
}

// Index returns paginated knowledge list with filters.
// GET /admin/link-knowledge
func (h *LinkKnowledgeHandler) Index(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	query := h.db.Model(&model.LinkKnowledge{})

	if title := c.Query("title"); title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if linkCause := c.Query("link_cause"); linkCause != "" {
		query = query.Where("link_cause = ?", linkCause)
	}
	if typ := c.Query("type"); typ != "" {
		query = query.Where("type = ?", typ)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var items []model.LinkKnowledge
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&items).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	// Enrich with keywords and link cause names
	type EnrichedKnowledge struct {
		model.LinkKnowledge
		Keywords     string `json:"keywords"`
		TypeName     string `json:"type_name"`
		LevelViewName string `json:"level_view_name"`
	}

	typeNames := map[string]string{
		"1": "文本回复",
		"2": "图片回复",
	}

	enriched := make([]EnrichedKnowledge, len(items))
	for i, item := range items {
		enriched[i] = EnrichedKnowledge{
			LinkKnowledge: item,
			TypeName:      typeNames[item.Type],
		}

		// Get keywords
		var keywords []model.LinkKeyword
		h.db.Where("belong = ? AND relid = ?", "knowledge", item.ID).Find(&keywords)
		kwStr := ""
		for j, kw := range keywords {
			if j > 0 {
				kwStr += ","
			}
			kwStr += kw.Keyword
		}
		enriched[i].Keywords = kwStr

		// Get link cause name
		if item.LinkCause > 0 {
			var cause model.LinkCause
			if err := h.db.First(&cause, item.LinkCause).Error; err == nil {
				enriched[i].LevelViewName = cause.Name
			}
		}
	}

	response.SuccessPage(c, enriched, total, page, pageSize)
}

// Edit returns a single knowledge item for editing.
// GET /admin/link-knowledge/:id/edit
func (h *LinkKnowledgeHandler) Edit(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var item model.LinkKnowledge
	if err := h.db.First(&item, id).Error; err != nil {
		response.NotFound(c, "knowledge not found")
		return
	}

	// Get keywords
	var keywords []model.LinkKeyword
	h.db.Where("belong = ? AND relid = ?", "knowledge", item.ID).Find(&keywords)
	kwStr := ""
	for j, kw := range keywords {
		if j > 0 {
			kwStr += ","
		}
		kwStr += kw.Keyword
	}

	typeNames := map[string]string{
		"1": "文本回复",
		"2": "图片回复",
	}

	response.Success(c, gin.H{
		"data": gin.H{
			"id":         item.ID,
			"title":      item.Title,
			"link_cause": item.LinkCause,
			"type":       item.Type,
			"status":     item.Status,
			"reply":      item.Content,
			"keywords":   kwStr,
		},
		"type": typeNames,
	})
}

// Save updates an existing knowledge item.
// PUT /admin/link-knowledge/:id
func (h *LinkKnowledgeHandler) Save(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	var req struct {
		Title     string `json:"title" binding:"required"`
		LinkCause uint   `json:"link_cause" binding:"required"`
		Type      string `json:"type" binding:"required"`
		Status    int16  `json:"status"`
		Reply     string `json:"reply" binding:"required"`
		Keywords  string `json:"keyword"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Verify exists
	var item model.LinkKnowledge
	if err := h.db.First(&item, id).Error; err != nil {
		response.NotFound(c, "knowledge not found")
		return
	}

	// Verify link cause exists
	var cause model.LinkCause
	if err := h.db.First(&cause, req.LinkCause).Error; err != nil {
		response.BadRequest(c, "link cause not found")
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		updates := map[string]interface{}{
			"title":      req.Title,
			"link_cause": req.LinkCause,
			"type":       req.Type,
			"status":     req.Status,
			"content":    req.Reply,
			"updated_at": time.Now(),
		}
		if err := tx.Model(&item).Updates(updates).Error; err != nil {
			return err
		}

		// Update keywords
		if req.Keywords != "" {
			tx.Where("belong = ? AND relid = ?", "knowledge", id).Delete(&model.LinkKeyword{})
			keywords := splitKeywords(req.Keywords)
			for _, kw := range keywords {
				tx.Create(&model.LinkKeyword{
					Keyword:    kw,
					Belong:     "knowledge",
					RelID:      uint(id),
					Status:     1,
					CreateTime: time.Now().Unix(),
				})
			}
		}
		return nil
	})

	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "knowledge updated")
}

// Create returns available types for creating knowledge.
// GET /admin/link-knowledge/create
func (h *LinkKnowledgeHandler) Create(c *gin.Context) {
	typeNames := map[string]string{
		"1": "文本回复",
		"2": "图片回复",
	}
	response.Success(c, gin.H{"type": typeNames})
}

// Add creates a new knowledge item.
// POST /admin/link-knowledge
func (h *LinkKnowledgeHandler) Add(c *gin.Context) {
	var req struct {
		Title     string `json:"title" binding:"required"`
		LinkCause uint   `json:"link_cause" binding:"required"`
		Type      string `json:"type" binding:"required"`
		Status    int16  `json:"status"`
		Reply     string `json:"reply" binding:"required"`
		Keywords  string `json:"keyword"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Verify link cause exists
	var cause model.LinkCause
	if err := h.db.First(&cause, req.LinkCause).Error; err != nil {
		response.BadRequest(c, "link cause not found")
		return
	}

	item := model.LinkKnowledge{
		Title:     req.Title,
		LinkCause: req.LinkCause,
		Type:      req.Type,
		Status:    req.Status,
		Content:   req.Reply,
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&item).Error; err != nil {
			return err
		}

		// Save keywords
		if req.Keywords != "" {
			keywords := splitKeywords(req.Keywords)
			for _, kw := range keywords {
				tx.Create(&model.LinkKeyword{
					Keyword:    kw,
					Belong:     "knowledge",
					RelID:      item.ID,
					Status:     1,
					CreateTime: time.Now().Unix(),
				})
			}
		}
		return nil
	})

	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "knowledge created")
}

// Delete deletes a knowledge item.
// DELETE /admin/link-knowledge/:id
func (h *LinkKnowledgeHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid id")
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		// Delete keywords
		if err := tx.Where("belong = ? AND relid = ?", "knowledge", id).Delete(&model.LinkKeyword{}).Error; err != nil {
			return err
		}
		// Delete knowledge
		if err := tx.Delete(&model.LinkKnowledge{}, id).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessMsg(c, "knowledge deleted")
}

// splitKeywords splits comma-separated keywords and trims whitespace.
func splitKeywords(s string) []string {
	var result []string
	current := ""
	for _, ch := range s {
		if ch == ',' {
			if current != "" {
				result = append(result, current)
			}
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}
