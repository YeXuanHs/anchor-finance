<template>
  <div class="transaction-page">
    <!-- 统计卡片 -->
    <div class="stats-row">
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon income">
            <el-icon :size="24"><Top /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">总收入</div>
            <div class="stat-value">¥{{ stats.totalIncome.toFixed(2) }}</div>
          </div>
        </div>
      </el-card>
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon expense">
            <el-icon :size="24"><Bottom /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">总支出</div>
            <div class="stat-value">¥{{ stats.totalExpense.toFixed(2) }}</div>
          </div>
        </div>
      </el-card>
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon recharge">
            <el-icon :size="24"><Wallet /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">本月充值</div>
            <div class="stat-value">¥{{ stats.monthlyRecharge.toFixed(2) }}</div>
          </div>
        </div>
      </el-card>
      <el-card class="stat-card" shadow="hover">
        <div class="stat-content">
          <div class="stat-icon withdraw">
            <el-icon :size="24"><CreditCard /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-label">本月提现</div>
            <div class="stat-value">¥{{ stats.monthlyWithdraw.toFixed(2) }}</div>
          </div>
        </div>
      </el-card>
    </div>

    <!-- Tabs 切换 -->
    <el-card class="main-card">
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="充值记录" name="recharge">
          <RechargeRecord ref="rechargeRef" />
        </el-tab-pane>
        <el-tab-pane label="退款记录" name="refund">
          <RefundRecord ref="refundRef" />
        </el-tab-pane>
        <el-tab-pane label="提现记录" name="withdraw">
          <WithdrawRecord ref="withdrawRef" />
        </el-tab-pane>
        <el-tab-pane label="信用记录" name="credit">
          <CreditRecord ref="creditRef" />
        </el-tab-pane>
        <el-tab-pane label="余额变动" name="accounts">
          <AccountsRecord ref="accountsRef" />
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Top, Bottom, Wallet, CreditCard } from '@element-plus/icons-vue'
import request from '@/utils/request'
import RechargeRecord from './RechargeRecord.vue'
import RefundRecord from './RefundRecord.vue'
import WithdrawRecord from './WithdrawRecord.vue'
import CreditRecord from './CreditRecord.vue'
import AccountsRecord from './AccountsRecord.vue'

const activeTab = ref('recharge')

const stats = reactive({
  totalIncome: 58620.00,
  totalExpense: 32450.80,
  monthlyRecharge: 8500.00,
  monthlyWithdraw: 3000.00
})

const rechargeRef = ref()
const refundRef = ref()
const withdrawRef = ref()
const creditRef = ref()
const accountsRef = ref()

const handleTabChange = (tab: string | number) => {
  // 切换tab时可触发子组件刷新
}

onMounted(async () => {
  try {
    const { data } = await request.get('/api/v1/balances/logs')
    if (data?.data) {
      stats.totalIncome = data.data.totalIncome || 0
      stats.totalExpense = data.data.totalExpense || 0
      stats.monthlyRecharge = data.data.monthlyRecharge || 0
      stats.monthlyWithdraw = data.data.monthlyWithdraw || 0
    }
  } catch (e) {
    console.error('Failed to fetch transaction stats:', e)
  }
})
</script>

<style scoped lang="scss">
.transaction-page {
  .stats-row {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 16px;
    margin-bottom: 20px;
  }

  .stat-card {
    border-radius: 12px;
    border: none;

    :deep(.el-card__body) {
      padding: 20px;
    }
  }

  .stat-content {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .stat-icon {
    width: 56px;
    height: 56px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;

    &.income {
      background: linear-gradient(135deg, #52c41a, #389e0d);
    }

    &.expense {
      background: linear-gradient(135deg, #ff4d4f, #cf1322);
    }

    &.recharge {
      background: linear-gradient(135deg, #1890ff, #096dd9);
    }

    &.withdraw {
      background: linear-gradient(135deg, #faad14, #d48806);
    }
  }

  .stat-info {
    .stat-label {
      font-size: 14px;
      color: #8c8c8c;
      margin-bottom: 4px;
    }

    .stat-value {
      font-size: 24px;
      font-weight: 600;
      color: #262626;
    }
  }

  .main-card {
    border-radius: 12px;
    border: none;

    :deep(.el-tabs__header) {
      margin-bottom: 20px;
    }

    :deep(.el-tabs__item) {
      font-size: 15px;
      padding: 0 24px;
      height: 48px;
      line-height: 48px;
    }
  }
}
</style>
