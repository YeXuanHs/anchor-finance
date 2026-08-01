package service

import (
	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

// LanguageService 语言服务
type LanguageService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewLanguageService(db *gorm.DB, log *logger.Logger) *LanguageService {
	return &LanguageService{db: db, log: log}
}

// GetLanguages 获取所有语言
func (s *LanguageService) GetLanguages() ([]model.Language, error) {
	var langs []model.Language
	err := s.db.Order("is_default DESC, code ASC").Find(&langs).Error
	return langs, err
}

// GetActiveLanguages 获取启用的语言
func (s *LanguageService) GetActiveLanguages() ([]model.Language, error) {
	var langs []model.Language
	err := s.db.Where("status = 1").Order("is_default DESC, code ASC").Find(&langs).Error
	return langs, err
}

// CreateLanguage 创建语言
func (s *LanguageService) CreateLanguage(lang *model.Language) error {
	return s.db.Create(lang).Error
}

// UpdateLanguage 更新语言
func (s *LanguageService) UpdateLanguage(id uint, updates map[string]interface{}) error {
	return s.db.Model(&model.Language{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteLanguage 删除语言
func (s *LanguageService) DeleteLanguage(id uint) error {
	return s.db.Delete(&model.Language{}, id).Error
}

// SetDefaultLanguage 设置默认语言
func (s *LanguageService) SetDefaultLanguage(id uint) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		// 取消所有默认
		if err := tx.Model(&model.Language{}).Where("is_default = true").Update("is_default", false).Error; err != nil {
			return err
		}
		// 设置新的默认
		return tx.Model(&model.Language{}).Where("id = ?", id).Update("is_default", true).Error
	})
}

// GetTranslations 获取指定语言的所有翻译
func (s *LanguageService) GetTranslations(langCode string) (map[string]string, error) {
	var results []struct {
		Key   string
		Value string
	}
	
	err := s.db.Table("lang_translations t").
		Select("k.key, t.value").
		Joins("JOIN lang_keys k ON k.id = t.key_id").
		Where("t.lang_code = ?", langCode).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	
	translations := make(map[string]string, len(results))
	for _, r := range results {
		translations[r.Key] = r.Value
	}
	return translations, nil
}

// GetTranslationsByModule 获取指定模块的翻译
func (s *LanguageService) GetTranslationsByModule(langCode, module string) (map[string]string, error) {
	var results []struct {
		Key   string
		Value string
	}
	
	err := s.db.Table("lang_translations t").
		Select("k.key, t.value").
		Joins("JOIN lang_keys k ON k.id = t.key_id").
		Where("t.lang_code = ? AND k.module = ?", langCode, module).
		Scan(&results).Error
	if err != nil {
		return nil, err
	}
	
	translations := make(map[string]string, len(results))
	for _, r := range results {
		translations[r.Key] = r.Value
	}
	return translations, err
}

// SaveTranslations 批量保存翻译
func (s *LanguageService) SaveTranslations(langCode string, translations map[string]string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for key, value := range translations {
			// 获取或创建key
			var langKey model.LangKey
			if err := tx.Where("`key` = ?", key).FirstOrCreate(&langKey, model.LangKey{Key: key}).Error; err != nil {
				return err
			}
			
			// 更新或创建翻译
			var trans model.LangTranslation
			result := tx.Where("key_id = ? AND lang_code = ?", langKey.ID, langCode).
				Assign(model.LangTranslation{Value: value}).
				FirstOrCreate(&trans)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

// ImportFromMap 从map导入翻译
func (s *LanguageService) ImportFromMap(langCode, module string, data map[string]string) error {
	return s.db.Transaction(func(tx *gorm.DB) error {
		for key, value := range data {
			// 获取或创建key
			var langKey model.LangKey
			if err := tx.Where("`key` = ?", key).FirstOrCreate(&langKey, model.LangKey{Key: key, Module: module}).Error; err != nil {
				return err
			}
			
			// 更新或创建翻译
			var trans model.LangTranslation
			result := tx.Where("key_id = ? AND lang_code = ?", langKey.ID, langCode).
				Assign(model.LangTranslation{Value: value}).
				FirstOrCreate(&trans)
			if result.Error != nil {
				return result.Error
			}
		}
		return nil
	})
}

// GetLangKeys 获取语言键列表
func (s *LanguageService) GetLangKeys(module string, page, pageSize int) ([]model.LangKey, int64, error) {
	var keys []model.LangKey
	var total int64
	
	query := s.db.Model(&model.LangKey{})
	if module != "" {
		query = query.Where("module = ?", module)
	}
	
	query.Count(&total)
	err := query.Order("id ASC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&keys).Error
	return keys, total, err
}
