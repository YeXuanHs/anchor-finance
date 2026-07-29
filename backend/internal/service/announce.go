package service

import (
	"errors"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/util"
	"anchorfinance/pkg/logger"

	"gorm.io/gorm"
)

type AnnounceService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewAnnounceService(db *gorm.DB, log *logger.Logger) *AnnounceService {
	return &AnnounceService{db: db, log: log}
}

type CreateAnnounceRequest struct {
	Title     string  `json:"title" binding:"required,max=256"`
	Content   string  `json:"content"`
	Type      string  `json:"type" binding:"omitempty,oneof=notice maintenance urgent"`
	Priority  int     `json:"priority"`
	Status    int16   `json:"status"`
	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`
	IsTop     bool    `json:"is_top"`
	AuthorID  uint    `json:"author_id"`
}

func (s *AnnounceService) Create(req CreateAnnounceRequest) (*model.Announce, error) {
	announce := model.Announce{
		Title:    req.Title,
		Content:  req.Content,
		Type:     req.Type,
		Priority: req.Priority,
		Status:   req.Status,
		IsTop:    req.IsTop,
		AuthorID: req.AuthorID,
	}
	if announce.Status == 0 {
		announce.Status = 1
	}
	if announce.Type == "" {
		announce.Type = "notice"
	}

	if req.StartTime != nil && *req.StartTime != "" {
		t, err := util.ParseTime(*req.StartTime)
		if err != nil {
			return nil, errors.New("invalid start_time")
		}
		announce.StartTime = &t
	}
	if req.EndTime != nil && *req.EndTime != "" {
		t, err := util.ParseTime(*req.EndTime)
		if err != nil {
			return nil, errors.New("invalid end_time")
		}
		announce.EndTime = &t
	}

	if err := s.db.Create(&announce).Error; err != nil {
		return nil, err
	}
	s.log.Infof("announce created: id=%d title=%s", announce.ID, announce.Title)
	return &announce, nil
}

func (s *AnnounceService) Update(id uint, updates map[string]interface{}) (*model.Announce, error) {
	var announce model.Announce
	if err := s.db.First(&announce, id).Error; err != nil {
		return nil, err
	}
	if err := s.db.Model(&announce).Updates(updates).Error; err != nil {
		return nil, err
	}
	s.db.First(&announce, id)
	return &announce, nil
}

func (s *AnnounceService) Delete(id uint) error {
	result := s.db.Delete(&model.Announce{}, id)
	if result.RowsAffected == 0 {
		return errors.New("announce not found")
	}
	return result.Error
}

func (s *AnnounceService) GetByID(id uint) (*model.Announce, error) {
	var announce model.Announce
	if err := s.db.First(&announce, id).Error; err != nil {
		return nil, err
	}
	return &announce, nil
}

func (s *AnnounceService) GetList(page, pageSize int, status int, announceType string) ([]model.Announce, int64, error) {
	var items []model.Announce
	var total int64

	query := s.db.Model(&model.Announce{})
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	if announceType != "" {
		query = query.Where("type = ?", announceType)
	}

	query.Count(&total)
	offset, limit := util.Paginate(page, pageSize)
	query.Offset(offset).Limit(limit).Order("is_top DESC, priority DESC, id DESC").Find(&items)
	return items, total, nil
}

func (s *AnnounceService) GetActive() ([]model.Announce, error) {
	var items []model.Announce
	now := time.Now()

	query := s.db.Where("status = 1")
	query = query.Where("(start_time IS NULL OR start_time <= ?)", now)
	query = query.Where("(end_time IS NULL OR end_time >= ?)", now)
	query.Order("is_top DESC, priority DESC, id DESC").Find(&items)

	ids := make([]uint, len(items))
	for i, item := range items {
		ids[i] = item.ID
	}
	if len(ids) > 0 {
		s.db.Model(&model.Announce{}).Where("id IN ?", ids).Update("views", gorm.Expr("views + 1"))
	}
	return items, nil
}
