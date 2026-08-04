<template>
  <div class="products-page">
    <div class="page-header">
      <h1 class="page-title">{{ $t('myProducts.title') }}</h1>
      <el-button type="primary" @click="$router.push('/products')">
        <el-icon><Plus /></el-icon>{{ $t('myProducts.buyNew') }}
      </el-button>
    </div>

    <el-radio-group v-model="statusFilter" class="status-filter">
      <el-radio-button value="all">{{ $t('helpCommon.all') }}</el-radio-button>
      <el-radio-button value="active">{{ $t('myProducts.active') }}</el-radio-button>
      <el-radio-button value="suspended">{{ $t('myProducts.suspended') }}</el-radio-button>
      <el-radio-button value="expired">{{ $t('myProducts.expired') }}</el-radio-button>
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
            <span class="meta-label">{{ $t('myProducts.dueDate') }}</span>
            <span class="meta-value">{{ product.expiry }}</span>
          </div>
          <div class="meta-item">
            <span class="meta-label">{{ $t('helpCommon.cost') }}</span>
            <span class="meta-value price">¥{{ product.price }}/月</span>
          </div>
        </div>

        <div class="product-actions">
          <el-button size="small" @click="handleManage(product)">{{ $t('myProducts.manage') }}</el-button>
          <el-button
            v-if="product.status === 'active'"
            size="small"
            type="primary"
            plain
            @click="handleRenew(product)"
          >{{ $t('myProducts.renew') }}</el-button>
          <el-button
            v-if="product.status === 'suspended'"
            size="small"
            type="warning"
            plain
            @click="handleReactivate(product)"
          >{{ $t('helpCommon.reactivate') }}</el-button>
        </div>
      </el-card>
    </div>

    <el-empty v-if="filteredProducts.length === 0" :description="$t('myProducts.noProducts')">
      <el-button type="primary" @click="$router.push('/products')">{{ $t('helpCommon.goToBuy') }}</el-button>
    </el-empty>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage } from 'element-plus'
import { Plus, Monitor, Connection, Lock, Coin } from '@element-plus/icons-vue'
import request from '@/utils/request'

const { t } = useI18n()
const router = useRouter()

const statusFilter = ref('all')
const loading = ref(false)

const iconMap: Record<string, any> = { Monitor, Connection, Lock, Coin }

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

const products = ref<Product[]>([])

onMounted(async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v2/user/products')
    const list = data.data?.list || data.list || data.data || []
    products.value = list.map((p: any) => ({ ...p, icon: iconMap[p.icon] || Monitor }))
  } catch (e) { console.error(e) } finally { loading.value = false }
})

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

function handleManage(product: Product) {
  router.push(`/user/products/${product.id}`)
}
function handleRenew(product: Product) { ElMessage.info(`${t('helpCommon.renewProduct')}${product.name}`) }
function handleReactivate(product: Product) { ElMessage.info(`${t('helpCommon.reactivateProduct')}${product.name}`) }
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
