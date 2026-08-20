package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type KnowledgeBaseHandler struct {
	svc *service.KnowledgeBaseService
	log *logger.Logger
}

func NewKnowledgeBaseHandler(svc *service.KnowledgeBaseService, log *logger.Logger) *KnowledgeBaseHandler {
	return &KnowledgeBaseHandler{svc: svc, log: log}
}

// ─── 分类 ───

// ListCategories 获取知识库分类列表
func (h *KnowledgeBaseHandler) ListCategories(c *gin.Context) {
	showInactive := c.Query("show_inactive") == "true"
	cats, err := h.svc.ListCategories(showInactive)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, cats)
}

// CreateCategory 创建分类
func (h *KnowledgeBaseHandler) CreateCategory(c *gin.Context) {
	var cat model.KnowledgeBaseCategory
	if err := c.ShouldBindJSON(&cat); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.CreateCategory(&cat); err != nil {
		response.ServerError(c, "创建失败")
		return
	}
	response.Success(c, cat)
}

// UpdateCategory 更新分类
func (h *KnowledgeBaseHandler) UpdateCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.UpdateCategory(uint(id), updates); err != nil {
		response.ServerError(c, "更新失败")
		return
	}
	response.SuccessMsg(c, "更新成功")
}

// DeleteCategory 删除分类
func (h *KnowledgeBaseHandler) DeleteCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err := h.svc.DeleteCategory(uint(id)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// ─── 文章 ───

// ListArticles 获取文章列表
func (h *KnowledgeBaseHandler) ListArticles(c *gin.Context) {
	categoryID, _ := strconv.ParseUint(c.Query("category_id"), 10, 64)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	articles, total, err := h.svc.ListArticles(uint(categoryID), page, pageSize, keyword)
	if err != nil {
		response.ServerError(c, "查询失败")
		return
	}
	response.Success(c, gin.H{
		"items": articles,
		"total": total,
		"page":  page,
	})
}

// GetArticle 获取文章详情
func (h *KnowledgeBaseHandler) GetArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	article, err := h.svc.GetArticle(uint(id))
	if err != nil {
		response.NotFound(c, "文章未找到")
		return
	}
	response.Success(c, article)
}

// CreateArticle 创建文章
func (h *KnowledgeBaseHandler) CreateArticle(c *gin.Context) {
	var article model.KnowledgeBaseArticle
	if err := c.ShouldBindJSON(&article); err != nil {
		response.BadRequest(c, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.CreateArticle(&article); err != nil {
		response.ServerError(c, "创建失败")
		return
	}
	response.Success(c, article)
}

// UpdateArticle 更新文章
func (h *KnowledgeBaseHandler) UpdateArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.UpdateArticle(uint(id), updates); err != nil {
		response.ServerError(c, "更新失败")
		return
	}
	response.SuccessMsg(c, "更新成功")
}

// DeleteArticle 删除文章
func (h *KnowledgeBaseHandler) DeleteArticle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err := h.svc.DeleteArticle(uint(id)); err != nil {
		response.ServerError(c, "删除失败")
		return
	}
	response.SuccessMsg(c, "删除成功")
}

// MarkHelpful 标记有帮助
func (h *KnowledgeBaseHandler) MarkHelpful(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "无效的ID")
		return
	}
	if err := h.svc.MarkHelpful(uint(id)); err != nil {
		response.ServerError(c, "操作失败")
		return
	}
	response.SuccessMsg(c, "感谢您的反馈")
}

// SearchPublic 前台搜索知识库
func (h *KnowledgeBaseHandler) SearchPublic(c *gin.Context) {
	keyword := c.Query("q")
	if keyword == "" {
		response.BadRequest(c, "请输入搜索关键词")
		return
	}
	articles, total, err := h.svc.ListArticles(0, 1, 20, keyword)
	if err != nil {
		response.ServerError(c, "搜索失败")
		return
	}
	response.Success(c, gin.H{
		"items": articles,
		"total": total,
	})
}
