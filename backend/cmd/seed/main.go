package main

import (
	"fmt"
	"log"

	"github.com/YeXuanHs/anchor-finance/config"
	"github.com/YeXuanHs/anchor-finance/internal/database"
	"github.com/YeXuanHs/anchor-finance/internal/model"
	"github.com/YeXuanHs/anchor-finance/internal/service"
)

func main() {
	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	database.Init(&cfg.Database)

	// 自动迁移
	db := database.GetDB()
	err := db.AutoMigrate(
		&model.User{},
		&model.Admin{},
		&model.Role{},
		&model.Order{},
		&model.Invoice{},
		&model.Service{},
		&model.Ticket{},
		&model.TicketDepartment{},
		&model.TicketStatus{},
		&model.Product{},
		&model.ProductGroup{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	// 创建默认角色
	role := model.Role{
		Name:        "超级管理员",
		Description: "系统超级管理员，拥有所有权限",
		IsSuper:     true,
	}
	db.FirstOrCreate(&role, model.Role{Name: "超级管理员"})

	// 创建默认管理员
	hashedPassword, err := service.HashPassword("admin123")
	if err != nil {
		log.Fatalf("Failed to hash password: %v", err)
	}

	admin := model.Admin{
		Username:     "admin",
		Email:        "admin@anchor-finance.com",
		PasswordHash: hashedPassword,
		RealName:     "系统管理员",
		RoleID:       role.ID,
		Status:       "active",
	}

	result := db.FirstOrCreate(&admin, model.Admin{Username: "admin"})
	if result.RowsAffected > 0 {
		fmt.Println("✅ 管理员账号创建成功")
		fmt.Println("   用户名: admin")
		fmt.Println("   密码: admin123")
	} else {
		fmt.Println("ℹ️  管理员账号已存在")
	}

	// 创建测试用户
	testPassword, _ := service.HashPassword("123456")
	testUser := model.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: testPassword,
		Phone:        "13800138000",
		Company:      "测试公司",
		Status:       "active",
	}
	db.FirstOrCreate(&testUser, model.User{Username: "testuser"})

	fmt.Println("✅ 测试用户创建成功")

	// 创建默认工单部门
	departments := []model.TicketDepartment{
		{Name: "技术支持", SortOrder: 1},
		{Name: "财务部门", SortOrder: 2},
		{Name: "销售部门", SortOrder: 3},
	}
	for _, dept := range departments {
		db.FirstOrCreate(&dept, model.TicketDepartment{Name: dept.Name})
	}
	fmt.Println("✅ 默认工单部门创建成功")

	// 创建默认工单状态
	statuses := []model.TicketStatus{
		{Value: "open", Label: "开启", SortOrder: 1},
		{Value: "pending", Label: "待回复", SortOrder: 2},
		{Value: "closed", Label: "已关闭", SortOrder: 3},
	}
	for _, status := range statuses {
		db.FirstOrCreate(&status, model.TicketStatus{Value: status.Value})
	}
	fmt.Println("✅ 默认工单状态创建成功")

	// 创建默认产品分组
	groups := []model.ProductGroup{
		{Name: "云服务器", SortOrder: 1},
		{Name: "虚拟主机", SortOrder: 2},
		{Name: "域名", SortOrder: 3},
	}
	for _, group := range groups {
		db.FirstOrCreate(&group, model.ProductGroup{Name: group.Name})
	}
	fmt.Println("✅ 默认产品分组创建成功")

	fmt.Println("\n🎉 数据初始化完成！")
}
