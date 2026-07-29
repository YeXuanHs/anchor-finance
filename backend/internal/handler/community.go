package handler

import (
	"strconv"

	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"

	"github.com/gin-gonic/gin"
)

type CommunityHandler struct {
	svc *service.CommunityService
	log *logger.Logger
}

func NewCommunityHandler(svc *service.CommunityService, log *logger.Logger) *CommunityHandler {
	return &CommunityHandler{svc: svc, log: log}
}

// GetPostList returns paginated community posts.
func (h *CommunityHandler) GetPostList(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	category := c.Query("category")
	keyword := c.Query("keyword")

	posts, total, err := h.svc.GetPostList(page, pageSize, category, keyword)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, posts, total, page, pageSize)
}

// GetPost returns a single community post.
func (h *CommunityHandler) GetPost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	post, err := h.svc.GetPostByID(uint(id))
	if err != nil {
		response.NotFound(c, "post not found")
		return
	}
	response.Success(c, post)
}

// CreatePost creates a new community post.
func (h *CommunityHandler) CreatePost(c *gin.Context) {
	var req struct {
		Title    string `json:"title" binding:"required"`
		Content  string `json:"content" binding:"required"`
		Category string `json:"category"`
		Tags     string `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	post := &model.CommunityPost{
		UserID:   c.GetUint("user_id"),
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Tags:     req.Tags,
		Status:   1,
	}

	if err := h.svc.CreatePost(post); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, post)
}

// UpdatePost updates a community post.
func (h *CommunityHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	var req struct {
		Title    string `json:"title"`
		Content  string `json:"content"`
		Category string `json:"category"`
		Tags     string `json:"tags"`
		Status   *int16 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	updates := make(map[string]interface{})
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.Content != "" {
		updates["content"] = req.Content
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.Tags != "" {
		updates["tags"] = req.Tags
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if err := h.svc.UpdatePost(uint(id), updates); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "post updated")
}

// DeletePost deletes a community post.
func (h *CommunityHandler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	if err := h.svc.DeletePost(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "post deleted")
}

// GetComments returns comments for a post.
func (h *CommunityHandler) GetComments(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	comments, total, err := h.svc.GetComments(uint(postID), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessPage(c, comments, total, page, pageSize)
}

// CreateComment creates a new comment.
func (h *CommunityHandler) CreateComment(c *gin.Context) {
	postID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid post id")
		return
	}

	var req struct {
		Content  string `json:"content" binding:"required"`
		ParentID *uint  `json:"parent_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	comment := &model.CommunityComment{
		PostID:   uint(postID),
		UserID:   c.GetUint("user_id"),
		ParentID: req.ParentID,
		Content:  req.Content,
		Status:   1,
	}

	if err := h.svc.CreateComment(comment); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.Success(c, comment)
}

// DeleteComment deletes a comment.
func (h *CommunityHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(c, "invalid comment id")
		return
	}

	if err := h.svc.DeleteComment(uint(id)); err != nil {
		response.ServerError(c, err.Error())
		return
	}
	response.SuccessMsg(c, "comment deleted")
}
