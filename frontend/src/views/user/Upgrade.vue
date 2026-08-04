<template>
  <div class="upgrade-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>产品升降级</span>
        </div>
      </template>

      <el-alert
        title="选择需要升降级的产品，系统将自动计算差价"
        type="info"
        show-icon
        :closable="false"
        style="margin-bottom: 20px"
      />

      <el-table :data="products" style="width: 100%" v-loading="loading">
        <el-table-column prop="product_name" label="当前产品" />
        <el-table-column prop="current_plan" label="当前方案" />
        <el-table-column prop="amount" label="当前金额">
          <template #default="{ row }">
            <span>¥{{ row.amount?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="next_due_date" label="下次续费" />
        <el-table-column label="操作" width="180">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="showUpgradeDialog(row)">升降级</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 升降级对话框 -->
    <el-dialog v-model="showDialog" title="产品升降级" width="600px">
      <div class="upgrade-form">
        <div class="current-info">
          <h4>当前产品信息</h4>
          <p>产品名称：{{ currentProduct?.product_name }}</p>
          <p>当前方案：{{ currentProduct?.current_plan }}</p>
          <p>当前金额：¥{{ currentProduct?.amount?.toFixed(2) }}</p>
        </div>

        <el-divider />

        <h4>选择新方案</h4>
        <el-radio-group v-model="selectedPlan" class="plan-options">
          <el-radio-button
            v-for="plan in availablePlans"
            :key="plan.id"
            :label="plan.id"
          >
            {{ plan.name }} - ¥{{ plan.price?.toFixed(2) }}
          </el-radio-button>
        </el-radio-group>

        <div class="price-diff" v-if="selectedPlan">
          <el-divider />
          <p>差价：<span class="diff-amount" :class="{ positive: priceDiff > 0, negative: priceDiff < 0 }">
            {{ priceDiff > 0 ? '+' : '' }}¥{{ priceDiff?.toFixed(2) }}
          </span></p>
          <p class="tip" v-if="priceDiff > 0">升级需要补差价</p>
          <p class="tip" v-else-if="priceDiff < 0">降级将退还差价到余额</p>
        </div>
      </div>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="submitUpgrade">确认升降级</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const products = ref([])
const showDialog = ref(false)
const currentProduct = ref<any>(null)
const selectedPlan = ref('')
const availablePlans = ref([])

const priceDiff = computed(() => {
  if (!selectedPlan.value || !currentProduct.value) return 0
  const newPlan = availablePlans.value.find(p => p.id === selectedPlan.value)
  return (newPlan?.price || 0) - (currentProduct.value?.amount || 0)
})

const showUpgradeDialog = (product: any) => {
  currentProduct.value = product
  selectedPlan.value = ''
  showDialog.value = true
  const res = await request.get(`/api/v2/upgrades/available/${product.id}`)
  availablePlans.value = res.data.data.plans
}

const submitUpgrade = async () => {
  await request.post('/api/v2/upgrades', { host_id: currentProduct.value.id, target_plan_id: selectedPlan.value })
  ElMessage.success('升降级请求已提交')
  showDialog.value = false
}
</script>

<style scoped lang="scss">
.upgrade-page {
  .upgrade-form {
    .current-info {
      background: #f5f7fa;
      padding: 15px;
      border-radius: 8px;

      h4 {
        margin: 0 0 10px;
      }

      p {
        margin: 5px 0;
        color: #606266;
      }
    }

    .plan-options {
      display: flex;
      flex-wrap: wrap;
      gap: 10px;
    }

    .price-diff {
      text-align: center;

      .diff-amount {
        font-size: 20px;
        font-weight: bold;

        &.positive {
          color: #e6a23c;
        }

        &.negative {
          color: #67c23a;
        }
      }

      .tip {
        color: #909399;
        font-size: 14px;
      }
    }
  }
}
</style>
