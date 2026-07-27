<template>
  <div class="wallet-page">
    <div class="page-header">
      <h2 class="page-title">我的钱包</h2>
      <p class="page-desc">管理您的余额和交易记录</p>
    </div>

    <!-- 余额卡片 -->
    <n-card class="balance-card" :bordered="false">
      <div class="balance-content">
        <div class="balance-info">
          <n-statistic label="账户余额 (CNY)">
            <div class="balance-amount">
              <span class="currency">¥</span>
              <span class="amount">{{ balance.toFixed(2) }}</span>
            </div>
          </n-statistic>
        </div>
        <div class="balance-actions">
          <n-button type="primary" size="large" @click="showRecharge = true">
            <template #icon>
              <n-icon :component="WalletOutline" />
            </template>
            充值
          </n-button>
          <n-button size="large" @click="showWithdraw = true">
            <template #icon>
              <n-icon :component="CashOutline" />
            </template>
            提现
          </n-button>
        </div>
      </div>
    </n-card>

    <!-- 充值弹窗 -->
    <n-modal v-model:show="showRecharge" preset="card" title="账户充值" style="width: 440px" :bordered="false">
      <n-form ref="rechargeFormRef" :model="rechargeForm" :rules="rechargeRules">
        <n-form-item path="amount" label="充值金额">
          <n-input-number
            v-model:value="rechargeForm.amount"
            :min="1"
            :max="50000"
            placeholder="请输入充值金额"
            size="large"
            style="width: 100%"
          >
            <template #prefix>¥</template>
          </n-input-number>
        </n-form-item>

        <n-form-item label="快捷金额">
          <div class="quick-amounts">
            <n-button
              v-for="amount in quickAmounts"
              :key="amount"
              :type="rechargeForm.amount === amount ? 'primary' : 'default'"
              @click="rechargeForm.amount = amount"
            >
              ¥{{ amount }}
            </n-button>
          </div>
        </n-form-item>

        <n-form-item path="payMethod" label="支付方式">
          <n-radio-group v-model:value="rechargeForm.payMethod">
            <n-space>
              <n-radio value="alipay">
                <div class="pay-method-label">
                  <span class="pay-icon alipay-icon">支</span>
                  支付宝
                </div>
              </n-radio>
              <n-radio value="wechat">
                <div class="pay-method-label">
                  <span class="pay-icon wechat-icon">微</span>
                  微信支付
                </div>
              </n-radio>
              <n-radio value="bank">
                <div class="pay-method-label">
                  <span class="pay-icon bank-icon">银</span>
                  银行卡
                </div>
              </n-radio>
            </n-space>
          </n-radio-group>
        </n-form-item>
      </n-form>

      <template #footer>
        <div class="modal-footer">
          <n-button @click="showRecharge = false">取消</n-button>
          <n-button type="primary" :loading="recharging" @click="handleRecharge">
            确认充值
          </n-button>
        </div>
      </template>
    </n-modal>

    <!-- 提现弹窗 -->
    <n-modal v-model:show="showWithdraw" preset="card" title="申请提现" style="width: 440px" :bordered="false">
      <n-form ref="withdrawFormRef" :model="withdrawForm" :rules="withdrawRules">
        <n-form-item path="amount" label="提现金额">
          <n-input-number
            v-model:value="withdrawForm.amount"
            :min="1"
            :max="balance"
            placeholder="请输入提现金额"
            size="large"
            style="width: 100%"
          >
            <template #prefix>¥</template>
          </n-input-number>
        </n-form-item>

        <div class="withdraw-hint">
          可用余额：<span class="highlight">¥{{ balance.toFixed(2) }}</span>
        </div>

        <n-form-item path="account" label="提现账户">
          <n-input
            v-model:value="withdrawForm.account"
            placeholder="请输入银行卡号或支付宝账号"
            size="large"
          />
        </n-form-item>
      </n-form>

      <template #footer>
        <div class="modal-footer">
          <n-button @click="showWithdraw = false">取消</n-button>
          <n-button type="primary" :loading="withdrawing" @click="handleWithdraw">
            提交申请
          </n-button>
        </div>
      </template>
    </n-modal>

    <!-- 交易记录 -->
    <n-card title="交易记录" class="records-card" :bordered="false">
      <n-data-table
        :columns="columns"
        :data="records"
        :pagination="pagination"
        :bordered="false"
        :single-line="false"
        striped
      />
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, h } from 'vue'
import { useMessage } from 'naive-ui'
import type { FormInst, FormRules, DataTableColumns } from 'naive-ui'
import { NTag } from 'naive-ui'
import { WalletOutline, CashOutline } from '@vicons/ionicons5'

const message = useMessage()

const balance = ref(12580.50)
const showRecharge = ref(false)
const showWithdraw = ref(false)
const recharging = ref(false)
const withdrawing = ref(false)

const rechargeFormRef = ref<FormInst | null>(null)
const withdrawFormRef = ref<FormInst | null>(null)

const quickAmounts = [50, 100, 200, 500, 1000, 2000]

const rechargeForm = ref({
  amount: 100 as number | null,
  payMethod: 'alipay'
})

const withdrawForm = ref({
  amount: null as number | null,
  account: ''
})

const rechargeRules: FormRules = {
  amount: [
    { required: true, type: 'number', message: '请输入充值金额', trigger: 'blur' }
  ],
  payMethod: { required: true, message: '请选择支付方式', trigger: 'change' }
}

const withdrawRules: FormRules = {
  amount: [
    { required: true, type: 'number', message: '请输入提现金额', trigger: 'blur' }
  ],
  account: { required: true, message: '请输入提现账户', trigger: 'blur' }
}

interface Record {
  id: number
  time: string
  type: 'recharge' | 'consume' | 'refund' | 'withdraw'
  amount: number
  balance: number
  remark: string
}

const typeMap: Record<string, { label: string; type: 'success' | 'warning' | 'error' | 'info' }> = {
  recharge: { label: '充值', type: 'success' },
  consume: { label: '消费', type: 'warning' },
  refund: { label: '退款', type: 'info' },
  withdraw: { label: '提现', type: 'error' }
}

const columns: DataTableColumns<Record> = [
  { title: '时间', key: 'time', width: 180 },
  {
    title: '类型',
    key: 'type',
    width: 100,
    render(row) {
      const config = typeMap[row.type]
      return h(NTag, { type: config.type, size: 'small', round: true }, { default: () => config.label })
    }
  },
  {
    title: '金额',
    key: 'amount',
    width: 140,
    render(row) {
      const isPositive = ['recharge', 'refund'].includes(row.type)
      const prefix = isPositive ? '+' : '-'
      const color = isPositive ? '#52c41a' : '#ff4d4f'
      return h('span', { style: { color, fontWeight: 600 } }, `${prefix}¥${Math.abs(row.amount).toFixed(2)}`)
    }
  },
  {
    title: '余额',
    key: 'balance',
    width: 140,
    render(row) {
      return h('span', { style: { fontWeight: 500 } }, `¥${row.balance.toFixed(2)}`)
    }
  },
  { title: '备注', key: 'remark', ellipsis: { tooltip: true } }
]

const records = ref<Record[]>([
  { id: 1, time: '2026-07-27 09:30:00', type: 'recharge', amount: 500, balance: 12580.50, remark: '支付宝充值' },
  { id: 2, time: '2026-07-26 15:20:00', type: 'consume', amount: 299, balance: 12080.50, remark: '购买产品 - 企业基础版' },
  { id: 3, time: '2026-07-25 10:00:00', type: 'refund', amount: 99, balance: 12379.50, remark: '订单退款 #ORD20260725' },
  { id: 4, time: '2026-07-24 14:15:00', type: 'consume', amount: 199, balance: 12280.50, remark: '续费服务' },
  { id: 5, time: '2026-07-23 08:00:00', type: 'withdraw', amount: 1000, balance: 12479.50, remark: '提现至工商银行' },
  { id: 6, time: '2026-07-22 16:45:00', type: 'recharge', amount: 2000, balance: 13479.50, remark: '微信充值' },
  { id: 7, time: '2026-07-21 11:30:00', type: 'consume', amount: 599, balance: 11479.50, remark: '购买产品 - 高级版' },
  { id: 8, time: '2026-07-20 09:00:00', type: 'recharge', amount: 1000, balance: 12078.50, remark: '支付宝充值' }
])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => { pagination.page = page },
  onUpdatePageSize: (size: number) => {
    pagination.pageSize = size
    pagination.page = 1
  }
})

async function handleRecharge() {
  try {
    await rechargeFormRef.value?.validate()
    recharging.value = true
    // TODO: 调用充值API
    await new Promise(resolve => setTimeout(resolve, 1500))
    balance.value += rechargeForm.value.amount || 0
    message.success('充值成功')
    showRecharge.value = false
    rechargeForm.value.amount = 100
  } catch {
    message.error('请填写完整信息')
  } finally {
    recharging.value = false
  }
}

async function handleWithdraw() {
  try {
    await withdrawFormRef.value?.validate()
    withdrawing.value = true
    // TODO: 调用提现API
    await new Promise(resolve => setTimeout(resolve, 1500))
    message.success('提现申请已提交，请等待审核')
    showWithdraw.value = false
    withdrawForm.value.amount = null
    withdrawForm.value.account = ''
  } catch {
    message.error('请填写完整信息')
  } finally {
    withdrawing.value = false
  }
}
</script>

<style scoped>
.wallet-page {
  max-width: 960px;
  margin: 0 auto;
  padding: 24px;
}

.page-header {
  margin-bottom: 24px;
}

.page-title {
  font-size: 24px;
  font-weight: 700;
  color: #1a1a1a;
  margin: 0 0 4px;
}

.page-desc {
  color: #8c8c8c;
  font-size: 14px;
  margin: 0;
}

.balance-card {
  margin-bottom: 24px;
  border-radius: 12px;
  background: linear-gradient(135deg, #1890ff, #096dd9);
  color: #fff;
}

.balance-card :deep(.n-card__content) {
  padding: 32px;
}

.balance-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.balance-info :deep(.n-statistic .n-statistic__label) {
  color: rgba(255, 255, 255, 0.85);
  font-size: 14px;
}

.balance-amount {
  display: flex;
  align-items: baseline;
}

.balance-amount .currency {
  font-size: 28px;
  font-weight: 600;
  margin-right: 4px;
}

.balance-amount .amount {
  font-size: 42px;
  font-weight: 700;
  letter-spacing: -1px;
}

.balance-actions {
  display: flex;
  gap: 12px;
}

.balance-actions .n-button {
  border-radius: 10px;
  font-weight: 500;
}

.balance-actions .n-button:first-child {
  background: rgba(255, 255, 255, 0.2);
  border: 1px solid rgba(255, 255, 255, 0.3);
  backdrop-filter: blur(10px);
}

.balance-actions .n-button:first-child:hover {
  background: rgba(255, 255, 255, 0.3);
}

.balance-actions .n-button:last-child {
  background: #fff;
  color: #1890ff;
  border: none;
}

.balance-actions .n-button:last-child:hover {
  background: #f0f7ff;
}

.quick-amounts {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.quick-amounts .n-button {
  border-radius: 8px;
}

.pay-method-label {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pay-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 700;
  color: #fff;
}

.alipay-icon {
  background: #1677ff;
}

.wechat-icon {
  background: #07c160;
}

.bank-icon {
  background: #f5222d;
}

.withdraw-hint {
  text-align: right;
  font-size: 13px;
  color: #8c8c8c;
  margin-bottom: 16px;
}

.withdraw-hint .highlight {
  color: #1890ff;
  font-weight: 600;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}

.records-card {
  border-radius: 12px;
}

.records-card :deep(.n-card-header__main) {
  font-weight: 600;
}
</style>
