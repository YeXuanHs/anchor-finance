<template>
  <div class="batch-renew-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>批量续费</span>
          <span class="selected-count">已选择 {{ selectedProducts.length }} 个产品</span>
        </div>
      </template>

      <el-table
        ref="tableRef"
        :data="products"
        style="width: 100%"
        @selection-change="handleSelectionChange"
      >
        <el-table-column type="selection" width="55" :selectable="isSelectable" />
        <el-table-column prop="name" label="产品名称" />
        <el-table-column prop="domain" label="标识/域名" />
        <el-table-column prop="expiry" label="到期时间">
          <template #default="{ row }">
            <span :class="{ 'text-danger': isExpired(row.expiry) }">{{ row.expiry }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="current_price" label="当前价格">
          <template #default="{ row }">
            <span>¥{{ row.current_price?.toFixed(2) }}/月</span>
          </template>
        </el-table-column>
        <el-table-column label="续费周期" width="200">
          <template #default="{ row }">
            <el-select v-model="row.renewCycle" size="small">
              <el-option label="1个月" :value="1" />
              <el-option label="3个月" :value="3" />
              <el-option label="6个月" :value="6" />
              <el-option label="1年" :value="12" />
              <el-option label="2年" :value="24" />
              <el-option label="3年" :value="36" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column label="续费金额">
          <template #default="{ row }">
            <span class="amount">¥{{ (row.current_price * row.renewCycle).toFixed(2) }}</span>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card v-if="selectedProducts.length > 0" class="page-card summary-card">
      <div class="summary-content">
        <div class="summary-info">
          <div class="summary-item">
            <span class="label">已选产品：</span>
            <span class="value">{{ selectedProducts.length }} 个</span>
          </div>
          <div class="summary-item">
            <span class="label">续费总额：</span>
            <span class="value total">¥{{ totalAmount.toFixed(2) }}</span>
          </div>
        </div>
        <div class="summary-actions">
          <el-button size="large" @click="clearSelection">清空选择</el-button>
          <el-button type="primary" size="large" :loading="submitting" @click="handleBatchRenew">
            立即续费
          </el-button>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

interface Product {
  id: number
  name: string
  domain: string
  expiry: string
  current_price: number
  renewCycle: number
  status: string
}

const tableRef = ref()
const products = ref<Product[]>([])
const selectedProducts = ref<Product[]>([])
const submitting = ref(false)
const renewCycle = ref(1)

const isSelectable = (row: Product) => {
  return row.status !== 'cancelled'
}

const isExpired = (date: string) => {
  return new Date(date) < new Date()
}

const totalAmount = computed(() => {
  return selectedProducts.value.reduce((sum, p) => sum + p.current_price * p.renewCycle, 0)
})

const handleSelectionChange = (selection: Product[]) => {
  selectedProducts.value = selection
}

const clearSelection = () => {
  tableRef.value?.clearSelection()
}

const handleBatchRenew = async () => {
  if (selectedProducts.value.length === 0) {
    ElMessage.warning('请选择要续费的产品')
    return
  }

  try {
    await ElMessageBox.confirm(
      `确认续费 ${selectedProducts.value.length} 个产品，总计 ¥${totalAmount.value.toFixed(2)}？`,
      '确认续费',
      { type: 'info' }
    )

    submitting.value = true
    await request.post('/api/v1/multi-renew', { host_ids: selectedProducts.value.map(p => p.id), cycle: renewCycle.value })
    ElMessage.success('续费订单已创建')
  } catch {
    // 用户取消
  } finally {
    submitting.value = false
  }
}

const loadProducts = async () => {
  const res = await request.get('/api/v1/user/products', { params: { status: 'active' } })
  products.value = res.data.data.list.map((p: any) => ({ ...p, renewCycle: 1 }))
}

onMounted(() => {
  loadProducts()
})
</script>

<style scoped lang="scss">
.batch-renew-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-card {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;

    .selected-count {
      font-size: 14px;
      color: #409eff;
    }
  }
}

.text-danger {
  color: #f56c6c;
}

.amount {
  color: #f56c6c;
  font-weight: bold;
}

.summary-card {
  position: sticky;
  bottom: 0;
  z-index: 10;
  border-radius: 12px;
  box-shadow: 0 -4px 12px rgba(0, 0, 0, 0.1);
}

.summary-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.summary-info {
  display: flex;
  gap: 32px;

  .summary-item {
    .label {
      color: #909399;
    }

    .value {
      font-weight: 600;
      color: #303133;

      &.total {
        color: #f56c6c;
        font-size: 20px;
      }
    }
  }
}

.summary-actions {
  display: flex;
  gap: 12px;
}
</style>
