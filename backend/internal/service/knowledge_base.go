package service

import (
	"fmt"
	"strings"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type KnowledgeBaseService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewKnowledgeBaseService(db *gorm.DB, log *logger.Logger) *KnowledgeBaseService {
	return &KnowledgeBaseService{db: db, log: log}
}

// ─── 分类 ───

func (s *KnowledgeBaseService) ListCategories(showInactive bool) ([]model.KnowledgeBaseCategory, error) {
	var cats []model.KnowledgeBaseCategory
	q := s.db.Order("sort_order ASC, id ASC")
	if !showInactive {
		q = q.Where("is_active = ?", true)
	}
	if err := q.Find(&cats).Error; err != nil {
		return nil, err
	}
	return cats, nil
}

func (s *KnowledgeBaseService) CreateCategory(cat *model.KnowledgeBaseCategory) error {
	return s.db.Create(cat).Error
}

func (s *KnowledgeBaseService) UpdateCategory(id uint, updates map[string]interface{}) error {
	updates["updated_at"] = gorm.Expr("NOW()")
	return s.db.Model(&model.KnowledgeBaseCategory{}).Where("id = ?", id).Updates(updates).Error
}

func (s *KnowledgeBaseService) DeleteCategory(id uint) error {
	// 检查是否有文章
	var count int64
	s.db.Model(&model.KnowledgeBaseArticle{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return fmt.Errorf("该分类下有 %d 篇文章，请先删除或移动文章", count)
	}
	return s.db.Delete(&model.KnowledgeBaseCategory{}, id).Error
}

// ─── 文章 ───

func (s *KnowledgeBaseService) ListArticles(categoryID uint, page, pageSize int, keyword string) ([]model.KnowledgeBaseArticle, int64, error) {
	var articles []model.KnowledgeBaseArticle
	var total int64

	q := s.db.Model(&model.KnowledgeBaseArticle{})
	if categoryID > 0 {
		q = q.Where("category_id = ?", categoryID)
	}
	if keyword != "" {
		like := "%" + keyword + "%"
		q = q.Where("title LIKE ? OR content LIKE ? OR tags LIKE ? OR keywords LIKE ?", like, like, like, like)
	}

	q.Count(&total)
	offset := (page - 1) * pageSize
	if err := q.Order("sort_order ASC, id DESC").Offset(offset).Limit(pageSize).Find(&articles).Error; err != nil {
		return nil, 0, err
	}
	return articles, total, nil
}

func (s *KnowledgeBaseService) GetArticle(id uint) (*model.KnowledgeBaseArticle, error) {
	var article model.KnowledgeBaseArticle
	if err := s.db.First(&article, id).Error; err != nil {
		return nil, err
	}
	// 增加浏览量
	s.db.Model(&article).UpdateColumn("view_count", gorm.Expr("view_count + 1"))
	return &article, nil
}

func (s *KnowledgeBaseService) CreateArticle(article *model.KnowledgeBaseArticle) error {
	return s.db.Create(article).Error
}

func (s *KnowledgeBaseService) UpdateArticle(id uint, updates map[string]interface{}) error {
	updates["updated_at"] = gorm.Expr("NOW()")
	return s.db.Model(&model.KnowledgeBaseArticle{}).Where("id = ?", id).Updates(updates).Error
}

func (s *KnowledgeBaseService) DeleteArticle(id uint) error {
	return s.db.Delete(&model.KnowledgeBaseArticle{}, id).Error
}

// SearchForAI AI 搜索知识库（用于自动回复）
func (s *KnowledgeBaseService) SearchForAI(query string, limit int) ([]model.KnowledgeBaseArticle, error) {
	var articles []model.KnowledgeBaseArticle

	// 提取关键词进行匹配
	keywords := extractKeywords(query)
	if len(keywords) == 0 {
		return articles, nil
	}

	// 构建搜索条件：匹配标题、内容、关键词、标签
	q := s.db.Model(&model.KnowledgeBaseArticle{}).Where("is_active = ?", true)

	conditions := []string{}
	args := []interface{}{}
	for _, kw := range keywords {
		like := "%" + kw + "%"
		conditions = append(conditions, "(title LIKE ? OR keywords LIKE ? OR tags LIKE ? OR summary LIKE ?)")
		args = append(args, like, like, like, like)
	}

	if len(conditions) > 0 {
		q = q.Where(strings.Join(conditions, " OR "), args...)
	}

	if err := q.Order("is_faq DESC, view_count DESC").Limit(limit).Find(&articles).Error; err != nil {
		return nil, err
	}
	return articles, nil
}

// MarkHelpful 标记文章有帮助
func (s *KnowledgeBaseService) MarkHelpful(id uint) error {
	return s.db.Model(&model.KnowledgeBaseArticle{}).Where("id = ?", id).UpdateColumn("help_count", gorm.Expr("help_count + 1")).Error
}

// extractKeywords 从文本中提取关键词（简单实现）
func extractKeywords(text string) []string {
	// 移除常见停用词
	stopWords := map[string]bool{
		"的": true, "了": true, "是": true, "在": true, "我": true, "有": true,
		"和": true, "就": true, "都": true, "而": true, "及": true, "与": true,
		"或": true, "但": true, "吗": true, "呢": true, "啊": true, "吧": true,
		"the": true, "a": true, "an": true, "is": true, "are": true, "was": true,
		"were": true, "be": true, "been": true, "being": true, "have": true,
		"has": true, "had": true, "do": true, "does": true, "did": true,
		"will": true, "would": true, "could": true, "should": true, "may": true,
		"might": true, "can": true, "shall": true, "to": true, "of": true,
		"in": true, "for": true, "on": true, "with": true, "at": true, "by": true,
		"from": true, "as": true, "into": true, "about": true, "how": true,
		"what": true, "why": true, "when": true, "where": true, "who": true,
		"which": true, "this": true, "that": true, "these": true, "those": true,
		"i": true, "you": true, "he": true, "she": true, "it": true, "we": true,
		"they": true, "me": true, "him": true, "her": true, "us": true, "them": true,
	}

	words := strings.FieldsFunc(text, func(r rune) bool {
		return r == ' ' || r == ',' || r == '。' || r == '，' || r == '！' || r == '？' ||
			r == '.' || r == '!' || r == '?' || r == ';' || r == '；' || r == '\n' || r == '\t'
	})

	var keywords []string
	seen := map[string]bool{}
	for _, w := range words {
		w = strings.TrimSpace(strings.ToLower(w))
		if len(w) < 2 || stopWords[w] || seen[w] {
			continue
		}
		seen[w] = true
		keywords = append(keywords, w)
	}
	return keywords
}
