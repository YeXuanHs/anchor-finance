package service

import (
	"errors"
	"fmt"
	"time"

	"anchorfinance/internal/model"
	"anchorfinance/internal/repository"
)

type SupplierService struct {
	supplierRepo *repository.SupplierRepository
}

func NewSupplierService() *SupplierService {
	return &SupplierService{
		supplierRepo: repository.NewSupplierRepository(),
	}
}

// GetSuppliers 获取供应商列表
func (s *SupplierService) GetSuppliers(params *model.SupplierQueryParams) (*model.SupplierListResponse, error) {
	if params.Page < 1 {
		params.Page = 1
	}
	if params.PageSize < 1 || params.PageSize > 100 {
		params.PageSize = 20
	}

	suppliers, total, err := s.supplierRepo.List(params)
	if err != nil {
		return nil, err
	}

	return &model.SupplierListResponse{
		List:     suppliers,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, nil
}

// GetSupplierByID 根据ID获取供应商
func (s *SupplierService) GetSupplierByID(id uint) (*model.Supplier, error) {
	return s.supplierRepo.FindByID(id)
}

// CreateSupplier 创建供应商
func (s *SupplierService) CreateSupplier(req *model.CreateSupplierRequest) (*model.Supplier, error) {
	supplier := &model.Supplier{
		Name:        req.Name,
		APIType:     req.APIType,
		APIURL:      req.APIURL,
		APIKey:      req.APIKey,
		APIPassword: req.APIPassword,
		Description: req.Description,
		Status:      "active",
	}

	if err := s.supplierRepo.Create(supplier); err != nil {
		return nil, err
	}

	return supplier, nil
}

// UpdateSupplier 更新供应商
func (s *SupplierService) UpdateSupplier(id uint, req *model.UpdateSupplierRequest) (*model.Supplier, error) {
	supplier, err := s.supplierRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("供应商不存在")
	}

	if req.Name != "" {
		supplier.Name = req.Name
	}
	if req.APIType != "" {
		supplier.APIType = req.APIType
	}
	if req.APIURL != "" {
		supplier.APIURL = req.APIURL
	}
	if req.APIKey != "" {
		supplier.APIKey = req.APIKey
	}
	if req.APIPassword != "" {
		supplier.APIPassword = req.APIPassword
	}
	if req.Description != "" {
		supplier.Description = req.Description
	}
	if req.Status != "" {
		supplier.Status = req.Status
	}

	if err := s.supplierRepo.Update(supplier); err != nil {
		return nil, err
	}

	return supplier, nil
}

// DeleteSupplier 删除供应商
func (s *SupplierService) DeleteSupplier(id uint) error {
	supplier, err := s.supplierRepo.FindByID(id)
	if err != nil {
		return errors.New("供应商不存在")
	}

	if supplier.ProductsCount > 0 {
		return errors.New("该供应商下还有产品，无法删除")
	}

	return s.supplierRepo.Delete(id)
}

// TestConnection 测试连接
func (s *SupplierService) TestConnection(id uint) (string, error) {
	supplier, err := s.supplierRepo.FindByID(id)
	if err != nil {
		return "", errors.New("供应商不存在")
	}

	switch supplier.APIType {
	case "manual":
		return "手动管理模式无需测试连接", nil
	case "zjmf":
		return s.testZjmfConnection(supplier)
	case "v10":
		return s.testV10Connection(supplier)
	case "anchor":
		return s.testAnchorConnection(supplier)
	default:
		return "", errors.New("不支持的API类型")
	}
}

// testZjmfConnection 测试智简魔方连接
func (s *SupplierService) testZjmfConnection(supplier *model.Supplier) (string, error) {
	if supplier.APIURL == "" || supplier.APIKey == "" {
		return "", errors.New("请先配置API地址和密钥")
	}
	// TODO: 实际调用zjmf API测试
	return "连接成功", nil
}

// testV10Connection 测试V10连接
func (s *SupplierService) testV10Connection(supplier *model.Supplier) (string, error) {
	if supplier.APIURL == "" || supplier.APIKey == "" {
		return "", errors.New("请先配置API地址和密钥")
	}
	// TODO: 实际调用V10 API测试
	return "连接成功", nil
}

// testAnchorConnection 测试锚点财务连接
func (s *SupplierService) testAnchorConnection(supplier *model.Supplier) (string, error) {
	if supplier.APIURL == "" || supplier.APIKey == "" {
		return "", errors.New("请先配置API地址和密钥")
	}
	// TODO: 实际调用锚点财务API测试
	return "连接成功", nil
}

// SyncProducts 同步产品
func (s *SupplierService) SyncProducts(id uint) (int, error) {
	supplier, err := s.supplierRepo.FindByID(id)
	if err != nil {
		return 0, errors.New("供应商不存在")
	}

	if supplier.APIType == "manual" {
		return 0, errors.New("手动管理模式不支持同步")
	}

	var syncedCount int
	switch supplier.APIType {
	case "zjmf":
		syncedCount, err = s.syncZjmfProducts(supplier)
	case "v10":
		syncedCount, err = s.syncV10Products(supplier)
	case "anchor":
		syncedCount, err = s.syncAnchorProducts(supplier)
	}

	if err != nil {
		return 0, err
	}

	// 更新最后同步时间和产品数量
	now := time.Now()
	supplier.LastSyncAt = &now
	supplier.ProductsCount += syncedCount
	s.supplierRepo.Update(supplier)

	return syncedCount, nil
}

// syncZjmfProducts 同步智简魔方产品
func (s *SupplierService) syncZjmfProducts(supplier *model.Supplier) (int, error) {
	// TODO: 实际调用zjmf API同步产品
	fmt.Println("同步zjmf产品...")
	return 0, nil
}

// syncV10Products 同步V10产品
func (s *SupplierService) syncV10Products(supplier *model.Supplier) (int, error) {
	// TODO: 实际调用V10 API同步产品
	fmt.Println("同步V10产品...")
	return 0, nil
}

// syncAnchorProducts 同步锚点财务产品
func (s *SupplierService) syncAnchorProducts(supplier *model.Supplier) (int, error) {
	// TODO: 实际调用锚点财务API同步产品
	fmt.Println("同步锚点财务产品...")
	return 0, nil
}

// ToggleStatus 切换状态
func (s *SupplierService) ToggleStatus(id uint, status string) error {
	supplier, err := s.supplierRepo.FindByID(id)
	if err != nil {
		return errors.New("供应商不存在")
	}

	supplier.Status = status
	return s.supplierRepo.Update(supplier)
}
