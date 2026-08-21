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
	fmt.Println("   用户名: testuser")
	fmt.Println("   密码: 123456")

	// 创建测试订单
	testOrder := model.Order{
		UserID:      testUser.ID,
		OrderNo:     "ORD20240001",
		ProductID:   1,
		ProductName: "云服务器 2核4G",
		Quantity:    1,
		Amount:      99.00,
		Status:      "paid",
	}
	db.FirstOrCreate(&testOrder, model.Order{OrderNo: "ORD20240001"})

	fmt.Println("✅ 测试订单创建成功")

	fmt.Println("\n🎉 数据初始化完成！")
}
