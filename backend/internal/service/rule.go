package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/pkg/logger"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type RuleService struct {
	db  *gorm.DB
	log *logger.Logger
}

func NewRuleService(db *gorm.DB, log *logger.Logger) *RuleService {
	return &RuleService{db: db, log: log}
}

type CreateRuleRequest struct {
	Name        string                 `json:"name" binding:"required,max=128"`
	Code        string                 `json:"code" binding:"omitempty,max=64"`
	Type        string                 `json:"type" binding:"required,oneof=pricing promotion notification limit filter routing"`
	Condition   map[string]interface{} `json:"condition" binding:"required"`
	Action      map[string]interface{} `json:"action" binding:"required"`
	Priority    int                    `json:"priority"`
	Description string                 `json:"description"`
	StartDate   *time.Time             `json:"start_date"`
	EndDate     *time.Time             `json:"end_date"`
	Extra       map[string]interface{} `json:"extra"`
}

type UpdateRuleRequest struct {
	Name        string                 `json:"name" binding:"omitempty,max=128"`
	Type        string                 `json:"type" binding:"omitempty,oneof=pricing promotion notification limit filter routing"`
	Condition   map[string]interface{} `json:"condition"`
	Action      map[string]interface{} `json:"action"`
	Priority    *int                   `json:"priority"`
	Status      *int16                 `json:"status"`
	Description string                 `json:"description"`
	StartDate   *time.Time             `json:"start_date"`
	EndDate     *time.Time             `json:"end_date"`
	Extra       map[string]interface{} `json:"extra"`
}

type TestRuleRequest struct {
	TestData map[string]interface{} `json:"test_data" binding:"required"`
}

// List returns paginated rules with filters.
func (s *RuleService) List(page, pageSize int, ruleType string, status *int16, keyword string) ([]model.Rule, int64, error) {
	var rules []model.Rule
	var total int64

	query := s.db.Model(&model.Rule{})
	if ruleType != "" {
		query = query.Where("type = ?", ruleType)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := rulePaginate(page, pageSize)
	if err := query.Offset(offset).Limit(limit).Order("priority DESC, id DESC").Find(&rules).Error; err != nil {
		return nil, 0, err
	}
	return rules, total, nil
}

// GetByID returns a single rule by ID.
func (s *RuleService) GetByID(id uint) (*model.Rule, error) {
	var rule model.Rule
	if err := s.db.First(&rule, id).Error; err != nil {
		return nil, err
	}
	return &rule, nil
}

// Create creates a new rule.
func (s *RuleService) Create(req CreateRuleRequest) (*model.Rule, error) {
	if req.Code != "" {
		var existing model.Rule
		if err := s.db.Where("code = ?", req.Code).First(&existing).Error; err == nil {
			return nil, errors.New("rule code already exists")
		}
	}

	rule := &model.Rule{
		Name:        req.Name,
		Code:        req.Code,
		Type:        req.Type,
		Condition:   toJSONMap(req.Condition),
		Action:      toJSONMap(req.Action),
		Priority:    req.Priority,
		Status:      1,
		Description: req.Description,
		StartDate:   req.StartDate,
		EndDate:     req.EndDate,
		Extra:       toJSONMap(req.Extra),
		Version:     1,
	}

	if err := s.db.Create(rule).Error; err != nil {
		return nil, err
	}

	s.log.Infof("rule created: %s (id=%d, type=%s)", rule.Name, rule.ID, rule.Type)
	return rule, nil
}

// Update updates a rule.
func (s *RuleService) Update(id uint, req UpdateRuleRequest) (*model.Rule, error) {
	var rule model.Rule
	if err := s.db.First(&rule, id).Error; err != nil {
		return nil, err
	}
	if rule.IsSystem {
		return nil, errors.New("system rule cannot be modified")
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Type != "" {
		updates["type"] = req.Type
	}
	if req.Condition != nil {
		updates["condition"] = toJSONMap(req.Condition)
	}
	if req.Action != nil {
		updates["action"] = toJSONMap(req.Action)
	}
	if req.Priority != nil {
		updates["priority"] = *req.Priority
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.StartDate != nil {
		updates["start_date"] = req.StartDate
	}
	if req.EndDate != nil {
		updates["end_date"] = req.EndDate
	}
	if req.Extra != nil {
		updates["extra"] = toJSONMap(req.Extra)
	}

	// Increment version on condition/action change
	if req.Condition != nil || req.Action != nil {
		updates["version"] = rule.Version + 1
	}

	if len(updates) > 0 {
		if err := s.db.Model(&rule).Updates(updates).Error; err != nil {
			return nil, err
		}
	}

	s.log.Infof("rule updated: %s (id=%d)", rule.Name, rule.ID)
	return s.GetByID(id)
}

// Delete deletes a rule (system rules excluded).
func (s *RuleService) Delete(id uint) error {
	var rule model.Rule
	if err := s.db.First(&rule, id).Error; err != nil {
		return err
	}
	if rule.IsSystem {
		return errors.New("system rule cannot be deleted")
	}
	return s.db.Delete(&rule).Error
}

// Enable enables a rule.
func (s *RuleService) Enable(id uint) error {
	var rule model.Rule
	if err := s.db.First(&rule, id).Error; err != nil {
		return err
	}
	if rule.Status == 1 {
		return errors.New("rule is already enabled")
	}
	return s.db.Model(&rule).Update("status", 1).Error
}

// Disable disables a rule.
func (s *RuleService) Disable(id uint) error {
	var rule model.Rule
	if err := s.db.First(&rule, id).Error; err != nil {
		return err
	}
	if rule.Status == 0 {
		return errors.New("rule is already disabled")
	}
	return s.db.Model(&rule).Update("status", 0).Error
}

// Test tests a rule against provided test data without executing actions.
func (s *RuleService) Test(id uint, req TestRuleRequest) (*model.RuleTestResult, error) {
	rule, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}

	start := time.Now()
	result := &model.RuleTestResult{}

	condition, err := s.parseJSON(rule.Condition)
	if err != nil {
		result.Error = fmt.Sprintf("condition parse error: %v", err)
		return result, nil
	}

	matched, details := s.evaluateCondition(condition, req.TestData)
	result.Matched = matched
	result.Duration = time.Since(start).Milliseconds()

	if matched {
		actions, err := s.parseJSON(rule.Action)
		if err != nil {
			result.Error = fmt.Sprintf("action parse error: %v", err)
			return result, nil
		}
		result.Actions = s.prepareActions(actions)
	}

	result.Details = details

	// Log the test
	logEntry := &model.RuleLog{
		RuleID:    rule.ID,
		MatchData: toJSONMap(req.TestData),
		Result:    toJSONMap(map[string]interface{}{"matched": matched}),
		Duration:  result.Duration,
		Success:   result.Error == "",
		ErrorMsg:  result.Error,
	}
	s.db.Create(logEntry)

	return result, nil
}

// GetActiveRules returns all active rules for a specific type, ordered by priority.
func (s *RuleService) GetActiveRules(ruleType string) ([]model.Rule, error) {
	var rules []model.Rule
	query := s.db.Where("status = 1").
		Order("priority DESC, id ASC")

	if ruleType != "" {
		query = query.Where("type = ?", ruleType)
	}

	// Filter by date range
	now := time.Now()
	query = query.Where("(start_date IS NULL OR start_date <= ?) AND (end_date IS NULL OR end_date >= ?)", now, now)

	if err := query.Find(&rules).Error; err != nil {
		return nil, err
	}
	return rules, nil
}

// MatchRules finds all matching rules for given data and type.
func (s *RuleService) MatchRules(ruleType string, data map[string]interface{}) ([]model.Rule, error) {
	rules, err := s.GetActiveRules(ruleType)
	if err != nil {
		return nil, err
	}

	var matched []model.Rule
	for _, rule := range rules {
		condition, err := s.parseJSON(rule.Condition)
		if err != nil {
			continue
		}
		ok, _ := s.evaluateCondition(condition, data)
		if ok {
			matched = append(matched, rule)

			// Update hit count
			now := time.Now()
			s.db.Model(&rule).Updates(map[string]interface{}{
				"hit_count":   rule.HitCount + 1,
				"last_hit_at": &now,
			})
		}
	}
	return matched, nil
}

// ExecuteRule executes a single rule's actions.
func (s *RuleService) ExecuteRule(id uint, data map[string]interface{}) (map[string]interface{}, error) {
	rule, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if rule.Status != 1 {
		return nil, errors.New("rule is not active")
	}

	condition, err := s.parseJSON(rule.Condition)
	if err != nil {
		return nil, fmt.Errorf("condition parse error: %v", err)
	}

	matched, _ := s.evaluateCondition(condition, data)
	if !matched {
		return nil, errors.New("rule condition not matched")
	}

	actions, err := s.parseJSON(rule.Action)
	if err != nil {
		return nil, fmt.Errorf("action parse error: %v", err)
	}

	// Execute actions
	execResults := make(map[string]interface{})
	start := time.Now()

	actionList := s.prepareActions(actions)
	for _, action := range actionList {
		actionType, _ := action["type"].(string)
		switch actionType {
		case "discount":
			// Apply pricing discount: {type:"discount", field:"price", value:0.8, mode:"multiply"}
			execResults["discount"] = action
			s.log.Infof("rule action: discount applied %+v", action)
		case "set_field":
			// Set a field on the matched entity: {type:"set_field", field:"status", value:1}
			execResults["set_field"] = action
			s.log.Infof("rule action: set_field %+v", action)
		case "send_notification":
			// Trigger notification: {type:"send_notification", channel:"email", template:"order_confirmed"}
			execResults["send_notification"] = action
			s.log.Infof("rule action: send_notification %+v", action)
		case "add_tag":
			// Add tag to entity: {type:"add_tag", tag:"vip"}
			execResults["add_tag"] = action
			s.log.Infof("rule action: add_tag %+v", action)
		case "block":
			// Block operation: {type:"block", reason:"exceeds limit"}
			execResults["block"] = action
			s.log.Infof("rule action: block %+v", action)
		case "webhook":
			// Call external webhook: {type:"webhook", url:"https://...", method:"POST"}
			execResults["webhook"] = action
			s.log.Infof("rule action: webhook %+v", action)
		default:
			// Unknown action type, log and store as-is
			execResults[actionType] = action
			s.log.Infof("rule action: unknown type=%s %+v", actionType, action)
		}
	}

	duration := time.Since(start).Milliseconds()

	// Log execution
	logEntry := &model.RuleLog{
		RuleID:    rule.ID,
		MatchData: toJSONMap(data),
		Actions:   toJSONMap(actions),
		Result:    toJSONMap(execResults),
		Duration:  duration,
		Success:   true,
	}
	s.db.Create(logEntry)

	// Update hit count
	now := time.Now()
	s.db.Model(&rule).Updates(map[string]interface{}{
		"hit_count":   rule.HitCount + 1,
		"last_hit_at": &now,
	})

	s.log.Infof("rule executed: %s (id=%d, duration=%dms)", rule.Name, rule.ID, duration)
	return execResults, nil
}

// GetRuleLogs returns execution logs for a rule.
func (s *RuleService) GetRuleLogs(page, pageSize int, ruleID *uint, success *bool) ([]model.RuleLog, int64, error) {
	var logs []model.RuleLog
	var total int64

	query := s.db.Model(&model.RuleLog{})
	if ruleID != nil {
		query = query.Where("rule_id = ?", *ruleID)
	}
	if success != nil {
		query = query.Where("success = ?", *success)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset, limit := rulePaginate(page, pageSize)
	if err := query.Preload("Rule").Offset(offset).Limit(limit).Order("id DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}
	return logs, total, nil
}

// parseJSON parses a datatypes.JSON into a map.
func (s *RuleService) parseJSON(data datatypes.JSON) (map[string]interface{}, error) {
	if data == nil {
		return nil, errors.New("empty data")
	}
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// evaluateCondition evaluates a condition against test data.
func (s *RuleService) evaluateCondition(condition map[string]interface{}, data map[string]interface{}) (bool, map[string]interface{}) {
	details := make(map[string]interface{})
	matched := true

	for key, expected := range condition {
		actual, exists := data[key]
		details[key] = map[string]interface{}{
			"expected": expected,
			"actual":   actual,
			"exists":   exists,
		}

		if !exists {
			matched = false
			continue
		}

		// Check operator-based conditions
		if condMap, ok := expected.(map[string]interface{}); ok {
			for op, val := range condMap {
				switch op {
				case "eq":
					if fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", val) {
						matched = false
					}
				case "ne":
					if fmt.Sprintf("%v", actual) == fmt.Sprintf("%v", val) {
						matched = false
					}
				case "gt":
					if !compareNumeric(actual, val, ">") {
						matched = false
					}
				case "gte":
					if !compareNumeric(actual, val, ">=") {
						matched = false
					}
				case "lt":
					if !compareNumeric(actual, val, "<") {
						matched = false
					}
				case "lte":
					if !compareNumeric(actual, val, "<=") {
						matched = false
					}
				case "in":
					if arr, ok := val.([]interface{}); ok {
						found := false
						for _, v := range arr {
							if fmt.Sprintf("%v", actual) == fmt.Sprintf("%v", v) {
								found = true
								break
							}
						}
						if !found {
							matched = false
						}
					}
				case "contains":
					if str, ok := actual.(string); ok {
						if sub, ok := val.(string); ok {
							if !contains(str, sub) {
								matched = false
							}
						}
					}
				}
			}
		} else {
			// Simple equality check
			if fmt.Sprintf("%v", actual) != fmt.Sprintf("%v", expected) {
				matched = false
			}
		}
	}

	return matched, details
}

// prepareActions prepares actions from a parsed action map.
func (s *RuleService) prepareActions(actions map[string]interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	if actionList, ok := actions["list"].([]interface{}); ok {
		for _, a := range actionList {
			if actionMap, ok := a.(map[string]interface{}); ok {
				result = append(result, actionMap)
			}
		}
	} else {
		result = append(result, actions)
	}
	return result
}

func compareNumeric(a, b interface{}, op string) bool {
	af := toFloat64(a)
	bf := toFloat64(b)
	switch op {
	case ">":
		return af > bf
	case ">=":
		return af >= bf
	case "<":
		return af < bf
	case "<=":
		return af <= bf
	}
	return false
}

func toFloat64(v interface{}) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case json.Number:
		f, _ := val.Float64()
		return f
	}
	return 0
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

func toJSONMap(v interface{}) datatypes.JSON {
	if v == nil {
		return nil
	}
	data, _ := json.Marshal(v)
	return datatypes.JSON(data)
}

func rulePaginate(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	pageSize = int(math.Min(float64(pageSize), 100))
	return (page - 1) * pageSize, pageSize
}
