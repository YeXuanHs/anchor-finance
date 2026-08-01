package handler

import (
	"anchorfinance/internal/model"
	"anchorfinance/internal/service"
	"anchorfinance/pkg/logger"
	"anchorfinance/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

// LanguageHandler 语言处理器
type LanguageHandler struct {
	svc *service.LanguageService
	log *logger.Logger
}

func NewLanguageHandler(svc *service.LanguageService, log *logger.Logger) *LanguageHandler {
	return &LanguageHandler{svc: svc, log: log}
}

// GetLanguages 获取语言列表
func (h *LanguageHandler) GetLanguages(c *gin.Context) {
	langs, err := h.svc.GetLanguages()
	if err != nil {
		response.ServerError(c, "获取语言列表失败")
		return
	}
	response.Success(c, langs)
}

// GetActiveLanguages 获取启用的语言
func (h *LanguageHandler) GetActiveLanguages(c *gin.Context) {
	langs, err := h.svc.GetActiveLanguages()
	if err != nil {
		response.ServerError(c, "获取语言列表失败")
		return
	}
	response.Success(c, langs)
}

// CreateLanguage 创建语言
func (h *LanguageHandler) CreateLanguage(c *gin.Context) {
	var lang model.Language
	if err := c.ShouldBindJSON(&lang); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.CreateLanguage(&lang); err != nil {
		response.ServerError(c, "创建失败")
		return
	}
	response.Success(c, lang)
}

// UpdateLanguage 更新语言
func (h *LanguageHandler) UpdateLanguage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var updates map[string]interface{}
	if err := c.ShouldBindJSON(&updates); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.UpdateLanguage(uint(id), updates); err != nil {
		response.ServerError(c, "更新失败")
		return
	}
	response.Success(c, nil)
}

// DeleteLanguage 删除语言
func (h *LanguageHandler) DeleteLanguage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.DeleteLanguage(uint(id)); err != nil {
		response.ServerError(c, "删除失败")
		return
	}
	response.Success(c, nil)
}

// SetDefaultLanguage 设置默认语言
func (h *LanguageHandler) SetDefaultLanguage(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := h.svc.SetDefaultLanguage(uint(id)); err != nil {
		response.ServerError(c, "设置失败")
		return
	}
	response.Success(c, nil)
}

// GetTranslations 获取翻译
func (h *LanguageHandler) GetTranslations(c *gin.Context) {
	langCode := c.Param("code")
	module := c.Query("module")
	
	var translations map[string]string
	var err error
	
	if module != "" {
		translations, err = h.svc.GetTranslationsByModule(langCode, module)
	} else {
		translations, err = h.svc.GetTranslations(langCode)
	}
	
	if err != nil {
		response.ServerError(c, "获取翻译失败")
		return
	}
	response.Success(c, translations)
}

// SaveTranslations 保存翻译
func (h *LanguageHandler) SaveTranslations(c *gin.Context) {
	langCode := c.Param("code")
	var translations map[string]string
	if err := c.ShouldBindJSON(&translations); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	if err := h.svc.SaveTranslations(langCode, translations); err != nil {
		response.ServerError(c, "保存失败")
		return
	}
	response.Success(c, nil)
}

// ImportTranslations 导入翻译
func (h *LanguageHandler) ImportTranslations(c *gin.Context) {
	langCode := c.Param("code")
	module := c.Query("module")
	
	var data map[string]string
	if err := c.ShouldBindJSON(&data); err != nil {
		response.BadRequest(c, "参数错误")
		return
	}
	
	if err := h.svc.ImportFromMap(langCode, module, data); err != nil {
		response.ServerError(c, "导入失败")
		return
	}
	response.Success(c, nil)
}

// GetLangKeys 获取语言键列表
func (h *LanguageHandler) GetLangKeys(c *gin.Context) {
	module := c.Query("module")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))
	
	keys, total, err := h.svc.GetLangKeys(module, page, pageSize)
	if err != nil {
		response.ServerError(c, "获取失败")
		return
	}
	response.Success(c, gin.H{"list": keys, "total": total})
}
