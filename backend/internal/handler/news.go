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

type NewsHandler struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewNewsHandler(db *gorm.DB, log *logger.Logger) *NewsHandler {
	return &NewsHandler{db: db, log: log}
}

// ---------- Public Endpoints ----------

// GetCategories returns all active news categories.
func (h *NewsHandler) GetCategories(c *gin.Context) {
	var categories []model.NewsCategory
	if err := h.db.Where("is_active = ?", true).Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, categories)
}

// GetList returns paginated published news with optional category and keyword filter.
func (h *NewsHandler) GetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var categoryID uint
	if cid := c.Query("category_id"); cid != "" {
		v, _ := strconv.ParseUint(cid, 10, 64)
		categoryID = uint(v)
	}

	var news []model.News
	var total int64

	query := h.db.Model(&model.News{}).Where("is_published = ?", true)
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR keywords LIKE ? OR summary LIKE ?", q, q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("is_sticky DESC, published_at DESC, id DESC").
		Preload("Category").
		Find(&news).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, news, total, page, pageSize)
}

// GetDetail returns a single published news by slug and increments view count.
func (h *NewsHandler) GetDetail(c *gin.Context) {
	slug := c.Param("slug")
	if slug == "" {
		response.BadRequest(c, "slug is required")
		return
	}

	var news model.News
	if err := h.db.Preload("Category").
		Where("slug = ? AND is_published = ?", slug, true).
		First(&news).Error; err != nil {
		response.NotFound(c, "news not found")
		return
	}

	h.db.Model(&news).UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	response.Success(c, news)
}

// Search searches news articles by keyword.
func (h *NewsHandler) Search(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		keyword = c.Query("keyword")
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	var news []model.News
	var total int64

	query := h.db.Model(&model.News{}).Where("is_published = ?", true)
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR summary LIKE ? OR content LIKE ?", q, q, q)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("created_at DESC").
		Preload("Category").
		Find(&news).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.SuccessPage(c, news, total, page, pageSize)
}

// ---------- Admin Endpoints ----------

// AdminGetCategories returns all news categories including inactive (admin).
func (h *NewsHandler) AdminGetCategories(c *gin.Context) {
	var categories []model.NewsCategory
	if err := h.db.Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, categories)
}

type CreateNewsCategoryRequest struct {
	Name      string `json:"name" binding:"required,max=50"`
	Slug      string `json:"slug" binding:"required,max=50"`
	SortOrder int    `json:"sort_order"`
}

// AdminCreateCategory creates a news category (admin).
func (h *NewsHandler) AdminCreateCategory(c *gin.Context) {
	var req CreateNewsCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	cat := model.NewsCategory{
		Name:      req.Name,
		Slug:      req.Slug,
		SortOrder: req.SortOrder,
		IsActive:  true,
	}
	if err := h.db.Create(&cat).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, cat)
}

type UpdateNewsCategoryRequest struct {
	Name      *string `json:"name"`
	Slug      *string `json:"slug"`
	SortOrder *int    `json:"sort_order"`
	IsActive  *bool   `json:"is_active"`
}

// AdminUpdateCategory updates a news category (admin).
func (h *NewsHandler) AdminUpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}

	var req UpdateNewsCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if len(updates) > 0 {
		if err := h.db.Model(&model.NewsCategory{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			response.ServerError(c, err.Error())
			return
		}
	}
	response.SuccessMsg(c, "category updated")
}

// AdminDeleteCategory deletes a news category (admin).
func (h *NewsHandler) AdminDeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid category id")
		return
	}

	var count int64
	h.db.Model(&model.News{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		response.BadRequest(c, "category has news articles, remove them first")
		return
	}

	if err := h.db.Delete(&model.NewsCategory{}, id).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "category deleted")
}

// AdminGetList returns all news including unpublished (admin).
func (h *NewsHandler) AdminGetList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	var categoryID uint
	if cid := c.Query("category_id"); cid != "" {
		v, _ := strconv.ParseUint(cid, 10, 64)
		categoryID = uint(v)
	}

	var news []model.News
	var total int64

	query := h.db.Model(&model.News{})
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		q := "%" + keyword + "%"
		query = query.Where("title LIKE ? OR keywords LIKE ?", q, q)
	}

	if err := query.Count(&total).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).
		Order("id DESC").
		Preload("Category").
		Find(&news).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, news, total, page, pageSize)
}

// AdminGetDetail returns a single news by ID (admin, includes unpublished).
func (h *NewsHandler) AdminGetDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid news id")
		return
	}

	var news model.News
	if err := h.db.Preload("Category").First(&news, id).Error; err != nil {
		response.NotFound(c, "news not found")
		return
	}
	response.Success(c, news)
}

type CreateNewsRequest struct {
	CategoryID  uint   `json:"category_id" binding:"required"`
	Title       string `json:"title" binding:"required,max=255"`
	Slug        string `json:"slug" binding:"required,max=255"`
	Summary     string `json:"summary" binding:"max=500"`
	Content     string `json:"content"`
	CoverImage  string `json:"cover_image" binding:"max=255"`
	Keywords    string `json:"keywords" binding:"max=255"`
	IsPublished bool   `json:"is_published"`
	IsSticky    bool   `json:"is_sticky"`
	AdminID     uint   `json:"admin_id"`
}

// AdminCreate creates a news article (admin).
func (h *NewsHandler) AdminCreate(c *gin.Context) {
	var req CreateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	news := model.News{
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Slug:        req.Slug,
		Summary:     req.Summary,
		Content:     req.Content,
		CoverImage:  req.CoverImage,
		Keywords:    req.Keywords,
		IsPublished: req.IsPublished,
		IsSticky:    req.IsSticky,
		AdminID:     req.AdminID,
	}

	if req.IsPublished {
		now := time.Now()
		news.PublishedAt = &now
	}

	if err := h.db.Create(&news).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, news)
}

type UpdateNewsRequest struct {
	CategoryID  *uint   `json:"category_id"`
	Title       *string `json:"title"`
	Slug        *string `json:"slug"`
	Summary     *string `json:"summary"`
	Content     *string `json:"content"`
	CoverImage  *string `json:"cover_image"`
	Keywords    *string `json:"keywords"`
	IsPublished *bool   `json:"is_published"`
	IsSticky    *bool   `json:"is_sticky"`
}

// AdminUpdate updates a news article (admin).
func (h *NewsHandler) AdminUpdate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid news id")
		return
	}

	var req UpdateNewsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if req.CategoryID != nil {
		updates["category_id"] = *req.CategoryID
	}
	if req.Title != nil {
		updates["title"] = *req.Title
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Summary != nil {
		updates["summary"] = *req.Summary
	}
	if req.Content != nil {
		updates["content"] = *req.Content
	}
	if req.CoverImage != nil {
		updates["cover_image"] = *req.CoverImage
	}
	if req.Keywords != nil {
		updates["keywords"] = *req.Keywords
	}
	if req.IsPublished != nil {
		updates["is_published"] = *req.IsPublished
	}
	if req.IsSticky != nil {
		updates["is_sticky"] = *req.IsSticky
	}

	if len(updates) > 0 {
		if err := h.db.Model(&model.News{}).Where("id = ?", id).Updates(updates).Error; err != nil {
			response.ServerError(c, err.Error())
			return
		}
	}
	response.SuccessMsg(c, "news updated")
}

// AdminDelete soft-deletes a news article (admin).
func (h *NewsHandler) AdminDelete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid news id")
		return
	}

	if err := h.db.Delete(&model.News{}, id).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "news deleted")
}

// ─── News Admin Extended Methods (from zjmf NewsController) ───

// GetCatsPage returns paginated news categories with filters.
func (h *NewsHandler) GetCatsPage(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	search := c.Query("search")
	status := c.Query("status")
	parentID, _ := strconv.ParseUint(c.Query("parent_id"), 10, 64)

	query := h.db.Model(&model.NewsCategory{})
	if search != "" {
		query = query.Where("title LIKE ?", "%"+search+"%")
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if parentID > 0 {
		query = query.Where("parent_id = ?", parentID)
	}

	var total int64
	query.Count(&total)

	var items []model.NewsCategory
	offset := (page - 1) * pageSize
	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&items)

	response.Success(c, gin.H{"list": items, "meta": gin.H{"total": total, "page": page, "limit": pageSize}})
}

// GetCateList returns all categories grouped by parent.
func (h *NewsHandler) GetCateList(c *gin.Context) {
	status := c.Query("status")
	parentID := c.Query("parent_id")

	query := h.db.Model(&model.NewsCategory{})
	if status != "" {
		query = query.Where("status = ?", status)
	}
	if parentID != "" {
		query = query.Where("parent_id = ?", parentID)
	}

	var list []model.NewsCategory
	query.Order("parent_id ASC, id DESC").Find(&list)

	data := make(map[uint]*map[string]interface{})
	var result []map[string]interface{}
	for _, v := range list {
		if v.ParentID > 0 {
			if parent, ok := data[v.ParentID]; ok {
				children, _ := (*parent)["list"].(map[uint]model.NewsCategory)
				if children == nil {
					children = make(map[uint]model.NewsCategory)
				}
				children[v.ID] = v
				(*parent)["list"] = children
			}
		} else {
			node := map[string]interface{}{
				"id":        v.ID,
				"title":     v.Title,
				"parent_id": v.ParentID,
				"list":      make(map[uint]model.NewsCategory),
			}
			data[v.ID] = &node
			result = append(result, node)
		}
	}
	response.Success(c, result)
}

// GetCatData returns a single category by ID.
func (h *NewsHandler) GetCatData(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if id == 0 {
		response.BadRequest(c, "id is required")
		return
	}

	var cat model.NewsCategory
	if err := h.db.First(&cat, id).Error; err != nil {
		response.NotFound(c, "category not found")
		return
	}
	response.Success(c, cat)
}

// PostEditCat creates or updates a news category.
func (h *NewsHandler) PostEditCat(c *gin.Context) {
	var req struct {
		ID       uint   `json:"id"`
		Title    string `json:"title" binding:"required"`
		ParentID int    `json:"parent_id"`
		Hidden   int    `json:"hidden"`
		Status   int    `json:"status"`
		Sort     int    `json:"sort"`
		Alias    string `json:"alias"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.ID > 0 && req.ID < 3 {
		response.BadRequest(c, "系统分类不可修改")
		return
	}

	if req.Alias != "" {
		var count int64
		query := h.db.Model(&model.NewsCategory{}).Where("alias = ?", req.Alias)
		if req.ID > 0 {
			query = query.Where("id != ?", req.ID)
		}
		query.Count(&count)
		if count > 0 {
			response.BadRequest(c, "别名已被使用")
			return
		}
	}

	if req.ID > 0 {
		updates := map[string]interface{}{
			"title":     req.Title,
			"parent_id": req.ParentID,
			"hidden":    req.Hidden,
			"status":    req.Status,
			"sort":      req.Sort,
		}
		if req.Alias != "" {
			updates["alias"] = req.Alias
		}
		h.db.Model(&model.NewsCategory{}).Where("id = ?", req.ID).Updates(updates)
	} else {
		cat := model.NewsCategory{
			Title:    req.Title,
			ParentID: uint(req.ParentID),
			SortOrder: req.Sort,
			IsActive: req.Status != 0,
		}
		h.db.Create(&cat)
	}
	response.SuccessMsg(c, "category saved")
}

// GetCheckalias checks if a category alias is available.
func (h *NewsHandler) GetCheckalias(c *gin.Context) {
	alias := c.Query("alias")
	if alias == "" {
		response.BadRequest(c, "alias is required")
		return
	}

	var count int64
	query := h.db.Model(&model.NewsCategory{}).Where("alias = ?", alias)
	if id := c.Query("id"); id != "" {
		query = query.Where("id != ?", id)
	}
	query.Count(&count)

	result := 1
	if count > 0 {
		result = 0
	}
	response.Success(c, result)
}

// DeleteCat deletes a news category (system categories 1-2 excluded).
func (h *NewsHandler) DeleteCat(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if id == 0 {
		response.BadRequest(c, "id is required")
		return
	}
	if id <= 2 {
		response.BadRequest(c, "系统分类不可删除")
		return
	}

	// Check if category has news
	var newsCount int64
	h.db.Model(&model.News{}).Where("category_id = ?", id).Count(&newsCount)
	if newsCount > 0 {
		response.BadRequest(c, "分类下有新闻，不可删除")
		return
	}

	if err := h.db.Delete(&model.NewsCategory{}, id).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "category deleted")
}

// GetContent returns a news article's content by ID.
func (h *NewsHandler) GetContent(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if id == 0 {
		response.BadRequest(c, "id is required")
		return
	}

	var news model.News
	if err := h.db.First(&news, id).Error; err != nil {
		response.NotFound(c, "news not found")
		return
	}
	response.Success(c, news)
}

// PostEditContent creates or updates a news article with content.
func (h *NewsHandler) PostEditContent(c *gin.Context) {
	var req struct {
		ID          uint   `json:"id"`
		ParentID    uint   `json:"parent_id" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Keywords    string `json:"keywords"`
		Description string `json:"description"`
		Content     string `json:"content"`
		CoverImage  string `json:"cover_image"`
		Hidden      int    `json:"hidden"`
		Sort        int    `json:"sort"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if req.ID > 0 {
		updates := map[string]interface{}{
			"category_id": req.ParentID,
			"title":       req.Title,
			"keywords":    req.Keywords,
			"summary":     req.Description,
			"content":     req.Content,
			"cover_image": req.CoverImage,
			"is_published": req.Hidden == 0,
			"sort_order":  req.Sort,
		}
		h.db.Model(&model.News{}).Where("id = ?", req.ID).Updates(updates)
	} else {
		news := model.News{
			CategoryID:  req.ParentID,
			Title:       req.Title,
			Keywords:    req.Keywords,
			Summary:     req.Description,
			Content:     req.Content,
			CoverImage:  req.CoverImage,
			IsPublished: req.Hidden == 0,
		}
		now := time.Now()
		if news.IsPublished {
			news.PublishedAt = &now
		}
		h.db.Create(&news)
	}
	response.SuccessMsg(c, "news saved")
}

// DeleteContent deletes a news article.
func (h *NewsHandler) DeleteContent(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if id == 0 {
		response.BadRequest(c, "id is required")
		return
	}

	if err := h.db.Delete(&model.News{}, id).Error; err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "news deleted")
}

// GetGetCustomParam returns custom fields for news.
func (h *NewsHandler) GetGetCustomParam(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	query := h.db.Table("customfields").Where("type = ?", "depot")
	var total int64
	query.Count(&total)

	var items []map[string]interface{}
	offset := (page - 1) * pageSize
	query.Select("id, fieldname").Order("id DESC").Offset(offset).Limit(pageSize).Find(&items)

	// Get values
	for i, item := range items {
		fieldID, _ := item["id"].(uint)
		var value struct{ Value string }
		h.db.Raw("SELECT value FROM customfieldsvalues WHERE fieldid = ?", fieldID).Scan(&value)
		item["value"] = value.Value
		items[i] = item
	}

	response.Success(c, gin.H{"count": total, "list": items})
}

// GetAddCustomParam adds a custom field.
func (h *NewsHandler) GetAddCustomParam(c *gin.Context) {
	var req struct {
		FieldName string `json:"fieldname" binding:"required"`
		Value     string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// Check uniqueness
	var count int64
	h.db.Table("customfields").Where("fieldname = ? AND type = 'depot'", req.FieldName).Count(&count)
	if count > 0 {
		response.BadRequest(c, "该自定义字段已存在")
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		res := tx.Exec("INSERT INTO customfields (type, relid, fieldname, fieldtype, description, create_time, update_time) VALUES ('depot', 0, ?, 'text', '站务自定义字段', ?, ?)",
			req.FieldName, time.Now().Unix(), time.Now().Unix())
		if res.Error != nil {
			return res.Error
		}
		var lid struct{ ID uint }
		tx.Raw("SELECT LAST_INSERT_ID() as id").Scan(&lid)
		tx.Exec("INSERT INTO customfieldsvalues (fieldid, relid, value) VALUES (?, 0, ?)", lid.ID, req.Value)
		return nil
	})
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "custom field added")
}

// GetUpdateCustomParam updates a custom field.
func (h *NewsHandler) GetUpdateCustomParam(c *gin.Context) {
	var req struct {
		ID        uint   `json:"id" binding:"required"`
		FieldName string `json:"fieldname" binding:"required"`
		Value     string `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		tx.Exec("UPDATE customfields SET fieldname = ?, update_time = ? WHERE id = ? AND type = 'depot'",
			req.FieldName, time.Now().Unix(), req.ID)
		tx.Exec("UPDATE customfieldsvalues SET value = ? WHERE fieldid = ?", req.Value, req.ID)
		return nil
	})
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "custom field updated")
}

// GetDelCustomParam deletes a custom field.
func (h *NewsHandler) GetDelCustomParam(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if id == 0 {
		response.BadRequest(c, "id is required")
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		tx.Exec("DELETE FROM customfields WHERE id = ? AND type = 'depot'", id)
		tx.Exec("DELETE FROM customfieldsvalues WHERE fieldid = ?", id)
		return nil
	})
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "custom field deleted")
}

// GetGetCustomUpdateVal returns custom field value for editing.
func (h *NewsHandler) GetGetCustomUpdateVal(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Query("id"), 10, 64)
	if id == 0 {
		response.BadRequest(c, "id is required")
		return
	}

	var field struct {
		ID        uint
		FieldName string
	}
	if err := h.db.Raw("SELECT id, fieldname FROM customfields WHERE id = ?", id).Scan(&field).Error; err != nil {
		response.NotFound(c, "record not found")
		return
	}

	var val struct{ Value string }
	h.db.Raw("SELECT value FROM customfieldsvalues WHERE fieldid = ?", id).Scan(&val)

	response.Success(c, gin.H{
		"id":    field.ID,
		"field": field.FieldName,
		"value": val.Value,
	})
}
