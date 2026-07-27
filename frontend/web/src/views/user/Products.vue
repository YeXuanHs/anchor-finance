<template>
  <div class="products-page">
    <div class="page-header">
      <h1 class="page-title">我的产品</h1>
      <n-input
        v-model:value="searchKey"
        placeholder="搜索产品名称或域名"
        clearable
        class="search-input"
      >
        <template #prefix>
          <n-icon :component="SearchOutline" color="#bfbfbf" />
        </template>
      </n-input>
    </div>

    <n-tabs v-model:value="activeTab" type="line" animated class="filter-tabs">
      <n-tab v-for="tab in statusTabs" :key="tab.value" :name="tab.value">
        {{ tab.label }}
        <n-badge
          v-if="tab.count > 0"
          :value="tab.count"
          :max="99"
          class="tab-badge"
        />
      </n-tab>
    </n-tabs>

    <div class="product-grid">
      <div
        v-for="product in filteredProducts"
        :key="product.id"
        class="product-card"
      >
        <div class="product-card-header">
          <div class="product-icon">
            <n-icon :size="24" :component="getProductIcon(product.type)" color="#1890ff" />
          </div>
          <n-tag :type="getStatusType(product.status)" size="small" round>
            {{ product.statusText }}
          </n-tag>
        </div>

        <div class="product-card-body">
          <h3 class="product-name">{{ product.name }}</h3>
          <div class="product-domain">
            <n-icon :component="GlobeOutline" size="14" />
            <span>{{ product.domain }}</span>
          </div>
          <div class="product-spec">{{ product.spec }}</div>
          <div class="product-expire">
            <n-icon :component="TimeOutline" size="14" />
            <span>下次到期：{{ product.expireDate }}</span>
          </div>
        </div>

        <div class="product-card-footer">
          <n-button
            v-if="product.status === 'active'"
            type="primary"
            size="small"
            @click="handleRenew(product)"
          >
            续费
          </n-button>
          <n-button
            size="small"
            quaternary
            @click="handleUpgrade(product)"
          >
            升降级
          </n-button>
          <n-button
            size="small"
            quaternary
            @click="handleManage(product)"
          >
            管理
          </n-button>
        </div>
      </div>

      <div v-if="filteredProducts.length === 0" class="empty-state">
        <n-icon :size="64" :component="CubeOutline" color="#d9d9d9" />
        <p>暂无产品</p>
        <n-button type="primary" @click="$router.push('/products')">购买产品</n-button>
      </div>
    </div>

    <div v-if="totalPages > 1" class="pagination-wrapper">
      <n-pagination
        v-model:page="currentPage"
        :page-count="totalPages"
        :page-slot="7"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useMessage } from 'naive-ui'
import {
  SearchOutline,
  CubeOutline,
  GlobeOutline,
  TimeOutline,
  ServerOutline,
  HardwareChipOutline,
  DesktopOutline,
  ShieldCheckmarkOutline
} from '@vicons/ionicons5'

const message = useMessage()
const searchKey = ref('')
const activeTab = ref('all')
const currentPage = ref(1)

const statusTabs = computed(() => [
  { label: '全部', value: 'all', count: products.value.length },
  { label: '使用中', value: 'active', count: products.value.filter(p => p.status === 'active').length },
  { label: '已暂停', value: 'suspended', count: products.value.filter(p => p.status === 'suspended').length },
  { label: '待开通', value: 'pending', count: products.value.filter(p => p.status === 'pending').length }
])

interface Product {
  id: string
  name: string
  type: string
  domain: string
  spec: string
  status: string
  statusText: string
  expireDate: string
}

const products = ref<Product[]>([
  {
    id: 'PRD001',
    name: '香港云服务器',
    type: 'ecs',
    domain: 'hk-ecs-01.example.com',
    spec: '2核4G / 50G SSD / 5Mbps',
    status: 'active',
    statusText: '使用中',
    expireDate: '2026-08-25'
  },
  {
    id: 'PRD002',
    name: '美国独立服务器',
    type: 'dedicated',
    domain: 'us-dedicated-01.example.com',
    spec: 'E5-2680v4 / 64G / 1T SSD',
    status: 'active',
    statusText: '使用中',
    expireDate: '2026-10-24'
  },
  {
    id: 'PRD003',
    name: 'OV SSL证书',
    type: 'ssl',
    domain: '*.example.com',
    spec: '单域名 / 企业验证',
    status: 'active',
    statusText: '使用中',
    expireDate: '2027-07-20'
  },
  {
    id: 'PRD004',
    name: '香港 VPS',
    type: 'vps',
    domain: 'hk-vps-01.example.com',
    spec: '1核2G / 30G NVMe / 1Gbps',
    status: 'suspended',
    statusText: '已暂停',
    expireDate: '2026-07-15'
  },
  {
    id: 'PRD005',
    name: '域名 example.com',
    type: 'domain',
    domain: 'example.com',
    spec: '.com / 首年注册',
    status: 'active',
    statusText: '使用中',
    expireDate: '2027-07-10'
  },
  {
    id: 'PRD006',
    name: '新加坡 VPS',
    type: 'vps',
    domain: 'sg-vps-01.example.com',
    spec: '2核4G / 60G NVMe',
    status: 'pending',
    statusText: '待开通',
    expireDate: '--'
  }
])

const filteredProducts = computed(() => {
  let result = products.value
  if (activeTab.value !== 'all') {
    result = result.filter(p => p.status === activeTab.value)
  }
  if (searchKey.value) {
    const key = searchKey.value.toLowerCase()
    result = result.filter(p =>
      p.name.toLowerCase().includes(key) ||
      p.domain.toLowerCase().includes(key)
    )
  }
  return result
})

const totalPages = computed(() => Math.ceil(filteredProducts.value.length / 12))

function getStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'error' | 'default'> = {
    active: 'success',
    suspended: 'warning',
    pending: 'info'
  }
  return map[status] || 'default'
}

function getProductIcon(type: string) {
  const map: Record<string, any> = {
    ecs: ServerOutline,
    vps: HardwareChipOutline,
    dedicated: DesktopOutline,
    ssl: ShieldCheckmarkOutline,
    domain: GlobeOutline
  }
  return map[type] || CubeOutline
}

function handleRenew(product: Product) {
  message.info(`续费产品：${product.name}`)
}

function handleUpgrade(product: Product) {
  message.info(`升降级：${product.name}`)
}

function handleManage(product: Product) {
  message.info(`管理产品：${product.name}`)
}
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
  font-size: 22px;
  font-weight: 700;
  color: #262626;
  margin: 0;
}

.search-input {
  width: 280px;
}

.filter-tabs {
  background: #fff;
  border-radius: 12px;
  padding: 4px 16px;
  border: 1px solid #f0f0f0;
}

.tab-badge {
  margin-left: 6px;
}

.product-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.product-card {
  background: #fff;
  border-radius: 12px;
  border: 1px solid #f0f0f0;
  overflow: hidden;
  transition: all 0.3s ease;
  display: flex;
  flex-direction: column;
}

.product-card:hover {
  border-color: #1890ff;
  box-shadow: 0 4px 16px rgba(24, 144, 255, 0.1);
  transform: translateY(-2px);
}

.product-card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 20px 12px;
}

.product-icon {
  width: 48px;
  height: 48px;
  background: #f0f5ff;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.product-card-body {
  flex: 1;
  padding: 0 20px 16px;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.product-name {
  font-size: 16px;
  font-weight: 600;
  color: #262626;
  margin: 0;
}

.product-domain {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #1890ff;
  font-family: 'Monaco', 'Menlo', monospace;
}

.product-spec {
  font-size: 12px;
  color: #8c8c8c;
  font-family: 'Monaco', 'Menlo', monospace;
}

.product-expire {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #595959;
}

.product-card-footer {
  display: flex;
  gap: 8px;
  padding: 12px 20px;
  border-top: 1px solid #f0f0f0;
  background: #fafafa;
}

.empty-state {
  grid-column: 1 / -1;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
  padding: 80px 0;
  background: #fff;
  border-radius: 12px;
  border: 1px solid #f0f0f0;
}

.empty-state p {
  margin: 0;
  color: #8c8c8c;
  font-size: 14px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  padding: 20px 0;
}

@media (max-width: 1200px) {
  .product-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }

  .search-input {
    width: 100%;
  }

  .product-grid {
    grid-template-columns: 1fr;
  }
}
</style>
