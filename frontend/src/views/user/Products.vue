<template>
  <div class="products-page">
    <div class="page-header">
      <h1 class="page-title">我的产品</h1>
      <el-button type="primary" @click="$router.push('/products')">
        <el-icon><Plus /></el-icon>购买新产品
      </el-button>
    </div>

    <el-radio-group v-model="statusFilter" class="status-filter">
      <el-radio-button value="all">全部</el-radio-button>
      <el-radio-button value="active">运行中</el-radio-button>
      <el-radio-button value="suspended">已暂停</el-radio-button>
      <el-radio-button value="expired">已过期</el-radio-button>
    </el-radio-group>

    <div class="products-grid">
      <el-card
        v-for="product in filteredProducts"
        :key="product.id"
        class="product-card"
        shadow="never"
      >
        <div class="product-header">
          <div class="product-icon-wrap">
            <el-icon :size="22" color="#0056FF"><component :is="product.icon" /></el-icon>
          </div>
          <el-tag :type="getStatusType(product.status)" size="small" effect="light" round>
            {{ product.statusText }}
          </el-tag>
        </div>

        <h3 class="product-name">{{ product.name }}</h3>
        <p class="product-domain">{{ product.domain }}</p>

        <div class="product-specs">
          <div v-for="spec in product.specs" :key="spec.label" class="spec-item">
            <span class="spec-label">{{ spec.label }}</span>
            <span class="spec-value">{{ spec.value }}</span>
          </div>
        </div>

        <div class="product-meta">
          <div class="meta-item">
            <span class="meta-label">到期时间</span>
            <span class="meta-value">{{ product.expiry }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-label">费用</span>
            <span class="meta-value price">¥{{ product.price }}/月</span>
          </div>
        </div>

        <div class="product-actions">
          <el-button size="small" @click="handleManage(product)">管理</el-button>
          <el-button
            v-if="product.status === 'active'"
            size="small"
            type="primary"
            plain
            @click="handleRenew(product)"
          >续费</el-button>
          <el-button
            v-if="product.status === 'suspended'"
            size="small"
            type="warning"
            plain
            @click="handleReactivate(product)"
          >恢复</el-button>
        </div>
      </el-card>
    </div>

    <el-empty v-if="filteredProducts.length === 0" description="暂无产品">
      <el-button type="primary" @click="$router.push('/products')">去购买</el-button>
    </el-empty>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus, Monitor, Connection, Lock, Coin } from '@element-plus/icons-vue'

const statusFilter = ref('all')

interface ProductSpec { label: string; value: string }
interface Product {
  id: number
  name: string
  domain: string
  icon: any
  status: string
  statusText: string
  specs: ProductSpec[]
  expiry: string
  price: string
}

const products = ref<Product[]>([
  {
    id: 1, name: '香港云服务器', domain: 'hk-web-01.anchorfin.com', icon: Monitor,
    status: 'active', statusText: '运行中',
    specs: [
      { label: 'CPU', value: '2核' }, { label: '内存', value: '4GB' },
      { label: '硬盘', value: '50G SSD' }, { label: '带宽', value: '5Mbps' }
    ],
    expiry: '2026-08-25', price: '49'
  },
  {
    id: 2, name: 'OV SSL证书', domain: '*.anchorfin.com', icon: Lock,
    status: 'active', statusText: '运行中',
    specs: [
      { label: '类型', value: '通配符' }, { label: '验证', value: '企业OV' }, { label: '有效期', value: '1年' }
    ],
    expiry: '2027-07-20', price: '16.6'
  },
  {
    id: 3, name: '新加坡 VPS', domain: 'sg-api.anchorfin.com', icon: Connection,
    status: 'suspended', statusText: '已暂停',
    specs: [
      { label: 'CPU', value: '2核' }, { label: '内存', value: '4GB' }, { label: '硬盘', value: '60G NVMe' }
    ],
    expiry: '2026-07-05', price: '35'
  },
  {
    id: 4, name: '域名注册', domain: 'anchorfin.com', icon: Coin,
    status: 'expired', statusText: '已过期',
    specs: [{ label: '后缀', value: '.com' }, { label: '年限', value: '1年' }],
    expiry: '2026-07-10', price: '7.5'
  }
])

const filteredProducts = computed(() => {
  if (statusFilter.value === 'all') return products.value
  return products.value.filter(p => p.status === statusFilter.value)
})

function getStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    active: 'success', suspended: 'warning', expired: 'danger'
  }
  return map[status] || 'info'
}

function handleManage(product: Product) { ElMessage.info(`管理产品：${product.name}`) }
function handleRenew(product: Product) { ElMessage.info(`续费产品：${product.name}`) }
function handleReactivate(product: Product) { ElMessage.info(`恢复产品：${product.name}`) }
</script>

<style scoped>
.products-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.status-filter {
  align-self: flex-start;
}

.products-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 16px;
}

.product-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
  transition: all 0.3s ease;
  background: #fff;
}

.product-card:hover {
  border-color: #0056FF;
  box-shadow: 0 2px 12px rgba(0,86,255,0.1);
  transform: translateY(-2px);
}

.product-card :deep(.el-card__body) {
  padding: 20px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.product-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.product-icon-wrap {
  width: 44px;
  height: 44px;
  border-radius: 10px;
  background: rgba(0,86,255,0.06);
  display: flex;
  align-items: center;
  justify-content: center;
}

.product-name {
  font-size: 17px;
  font-weight: 600;
  color: #303133;
  margin: 0;
}

.product-domain {
  font-size: 13px;
  color: #909399;
  margin: 0;
  font-family: 'Monaco', 'Menlo', monospace;
}

.product-specs {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.spec-item {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  background: #f5f7fa;
  border-radius: 6px;
  font-size: 12px;
}

.spec-label { color: #909399; }
.spec-value { color: #303133; font-weight: 500; }

.product-meta {
  display: flex;
  justify-content: space-between;
  padding-top: 8px;
  border-top: 1px solid #f2f3f5;
}

.meta-item { display: flex; flex-direction: column; gap: 2px; }
.meta-label { font-size: 12px; color: #c0c4cc; }
.meta-value { font-size: 14px; color: #303133; font-weight: 500; }
.meta-value.price { color: #fa8c16; font-weight: 700; }

.product-actions {
  display: flex;
  gap: 8px;
  padding-top: 4px;
}

@media (max-width: 768px) {
  .products-grid { grid-template-columns: 1fr; }
  .page-header { flex-direction: column; gap: 12px; align-items: flex-start; }
}
</style>
