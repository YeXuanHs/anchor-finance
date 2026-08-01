package model

// AIShoppingConfig AI购物助手配置
// 对应 mahiru_ai_shopping 的插件配置
// 简单配置：启用状态、API地址、密钥、模型、浮窗标题、欢迎语、系统提示词
type AIShoppingConfig struct {
	ID             uint   `gorm:"primaryKey" json:"id"`
	Key            string `gorm:"type:varchar(50);uniqueIndex;not null" json:"key"`
	Value          string `gorm:"type:text" json:"value"`
}

// AIShoppingChatLog AI购物聊天记录
type AIShoppingChatLog struct {
	ID        uint   `gorm:"primaryKey" json:"id"`
	SessionID string `gorm:"type:varchar(64);index" json:"session_id"`
	UserID    uint   `gorm:"index;default:0" json:"user_id"`
	Role      string `gorm:"type:varchar(20);not null;comment:user/assistant/system" json:"role"`
	Content   string `gorm:"type:text;not null" json:"content"`
	CreatedAt int64  `gorm:"autoCreateTime" json:"created_at"`
}

// ProductCatalogTool 商品目录工具定义
// 用于 AI 工具调用搜索商品
type ProductCatalogTool struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Description string `gorm:"type:varchar(500)" json:"description"`
	Parameters  string `gorm:"type:text;comment:JSON Schema" json:"parameters"`
	IsActive    bool   `gorm:"default:true" json:"is_active"`
}
