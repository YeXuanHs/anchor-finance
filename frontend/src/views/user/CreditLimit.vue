<template>
  <div class="credit-limit-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>信用额度</span>
          <el-button type="primary" @click="applyCredit">申请提额</el-button>
        </div>
      </template>

      <div class="credit-overview">
        <el-row :gutter="20">
          <el-col :span="8">
            <div class="credit-card">
              <div class="credit-label">总额度</div>
              <div class="credit-value">¥{{ creditInfo.total?.toFixed(2) || '0.00' }}</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="credit-card">
              <div class="credit-label">已用额度</div>
              <div class="credit-value used">¥{{ creditInfo.used?.toFixed(2) || '0.00' }}</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="credit-card">
              <div class="credit-label">可用额度</div>
              <div class="credit-value available">¥{{ creditInfo.available?.toFixed(2) || '0.00' }}</div>
            </div>
          </el-col>
        </el-row>
      </div>

      <el-divider />

      <h3>额度使用记录</h3>
      <el-table :data="creditLogs" style="width: 100%">
        <el-table-column prop="created_at" label="时间" />
        <el-table-column prop="type" label="类型">
          <template #default="{ row }">
            <el-tag :type="row.type === 'increase' ? 'success' : 'danger'">
              {{ row.type === 'increase' ? '增加' : '减少' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额">
          <template #default="{ row }">
            <span :class="row.type === 'increase' ? 'text-success' : 'text-danger'">
              {{ row.type === 'increase' ? '+' : '-' }}¥{{ row.amount?.toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="原因" />
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'

const creditInfo = ref({
  total: 10000,
  used: 3500,
  available: 6500
})

const creditLogs = ref([])

const applyCredit = () => {
  // TODO: 申请提额
}
</script>

<style scoped lang="scss">
.credit-limit-page {
  .credit-overview {
    margin-bottom: 20px;
  }

  .credit-card {
    background: #f5f7fa;
    padding: 20px;
    border-radius: 8px;
    text-align: center;

    .credit-label {
      color: #909399;
      font-size: 14px;
      margin-bottom: 8px;
    }

    .credit-value {
      font-size: 24px;
      font-weight: bold;
      color: #303133;

      &.used {
        color: #e6a23c;
      }

      &.available {
        color: #67c23a;
      }
    }
  }

  .text-success {
    color: #67c23a;
  }

  .text-danger {
    color: #f56c6c;
  }
}
</style>
