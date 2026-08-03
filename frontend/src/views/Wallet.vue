<template>
  <div class="wallet-page">
    <div class="wallet-header-section">
      <div class="balance-card">
        <div class="balance-card-inner">
          <div class="balance-info">
            <div class="balance-label">
              <n-icon :component="WalletOutline" size="20" />
              <span>账户余额</span>
            </div>
            <div class="balance-amount">
              <span class="currency">¥</span>
              <span class="amount">{{ balance.toFixed(2) }}</span>
            </div>
            <div class="balance-actions">
              <n-button type="primary" size="large" @click="showRechargeModal = true">
                <template #icon><n-icon :component="AddCircleOutline" /></template>
                充值
              </n-button>
              <n-button size="large" @click="showWithdrawModal = true">
                <template #icon><n-icon :component="ArrowUpOutline" /></template>
                提现
              </n-button>
            </div>
          </div>
          <div class="balance-decoration">
            <div class="deco-circle deco-circle-1"></div>
            <div class="deco-circle deco-circle-2"></div>
            <div class="deco-circle deco-circle-3"></div>
          </div>
        </div>
      </div>

      <div class="stats-row">
        <n-card class="stat-card" :bordered="false">
          <n-statistic label="本月充值" :value="monthlyRecharge" precision="2" prefix="¥" />
        </n-card>
        <n-card class="stat-card" :bordered="false">
          <n-statistic label="本月消费" :value="monthlyExpense" precision="2" prefix="¥" />
        </n-card>
        <n-card class="stat-card" :bordered="false">
          <n-statistic label="冻结金额" :value="frozenAmount" precision="2" prefix="¥" />
        </n-card>
      </div>
    </div>

    <n-card class="records-card" :bordered="false">
      <template #header>
        <div class="records-header">
          <span class="records-title">余额变动记录</span>
          <n-input
            v-model:value="searchKeyword"
            placeholder="搜索备注..."
            clearable
            style="width: 240px"
          >
            <template #prefix>
              <n-icon :component="SearchOutline" color="#1890ff" />
            </template>
          </n-input>
        </div>
      </template>

      <n-tabs v-model:value="activeType" type="segment" animated class="type-tabs">
        <n-tab-pane name="all" tab="全部" />
        <n-tab-pane name="recharge" tab="充值" />
        <n-tab-pane name="expense" tab="消费" />
        <n-tab-pane name="refund" tab="退款" />
        <n-tab-pane name="withdraw" tab="提现" />
      </n-tabs>

      <n-data-table
        :columns="columns"
        :data="filteredRecords"
        :bordered="false"
        :single-line="false"
        :pagination="pagination"
        class="records-table"
      />
    </n-card>

    <!-- 充值弹窗 -->
    <n-modal v-model:show="showRechargeModal" preset="card" title="账户充值" style="width: 480px" :bordered="false">
      <n-form ref="rechargeFormRef" :model="rechargeForm" :rules="rechargeRules">
        <n-form-item label="充值金额" path="amount">
          <n-input-number
            v-model:value="rechargeForm.amount"
            :min="1"
            :max="50000"
            :precision="2"
            placeholder="请输入充值金额"
            size="large"
            style="width: 100%"
          >
            <template #prefix>¥</template>
          </n-input-number>
        </n-form-item>

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

        <n-form-item label="支付方式" path="payMethod">
          <n-radio-group v-model:value="rechargeForm.payMethod">
            <div class="pay-methods">
              <n-radio value="alipay">
                <div class="pay-method-label">
                  <svg viewBox="0 0 24 24" width="20" height="20" fill="#1677FF">
                    <path d="M21.422 15.358c-3.492 1.464-6.924-0.096-9.348-1.776 0.78-2.016 1.32-4.392 1.32-6.78 0-1.176-0.18-2.292-0.456-3.324h-3.12v4.536h5.64v1.8H11.04v1.224c1.872 0.864 3.708 1.8 5.604 1.8 1.416 0 2.856-0.624 3.78-1.62v3.936zM12.012 2.484C6.48 2.484 2.004 6.96 2.004 12.492s4.476 10.008 10.008 10.008 10.008-4.476 10.008-10.008-4.476-10.008-10.008-10.008zm0 18c-4.416 0-8-3.584-8-8s3.584-8 8-8 8 3.584 8 8-3.584 8-8 8z"/>
                  </svg>
                  支付宝
                </div>
              </n-radio>
              <n-radio value="wechat">
                <div class="pay-method-label">
                  <svg viewBox="0 0 24 24" width="20" height="20" fill="#07c160">
                    <path d="M8.691 2.188C3.891 2.188 0 5.476 0 9.53c0 2.212 1.17 4.203 3.002 5.55a.59.59 0 0 1 .213.665l-.39 1.48c-.019.07-.048.141-.048.213 0 .163.13.295.29.295a.326.326 0 0 0 .167-.054l1.903-1.114a.864.864 0 0 1 .717-.098 10.16 10.16 0 0 0 2.837.403c.276 0 .543-.027.811-.05-.857-2.578.157-4.972 1.932-6.446 1.703-1.415 3.882-1.98 5.853-1.838-.576-3.583-4.196-6.348-8.596-6.348z"/>
                  </svg>
                  微信支付
                </div>
              </n-radio>
              <n-radio value="bank">
                <div class="pay-method-label">
                  <n-icon :component="CardOutline" size="20" color="#1890ff" />
                  银行卡
                </div>
              </n-radio>
            </div>
          </n-radio-group>
        </n-form-item>
      </n-form>

      <template #footer>
        <div class="modal-footer">
          <n-button @click="showRechargeModal = false">取消</n-button>
          <n-button type="primary" :loading="recharging" @click="handleRecharge">
            确认充值
          </n-button>
        </div>
      </template>
    </n-modal>

    <!-- 提现弹窗 -->
    <n-modal v-model:show="showWithdrawModal" preset="card" title="申请提现" style="width: 480px" :bordered="false">
      <n-form ref="withdrawFormRef" :model="withdrawForm" :rules="withdrawRules">
        <n-form-item label="可提现余额">
          <div class="available-balance">¥{{ balance.toFixed(2) }}</div>
        </n-form-item>

        <n-form-item label="提现金额" path="amount">
          <n-input-number
            v-model:value="withdrawForm.amount"
            :min="1"
            :max="balance"
            :precision="2"
            placeholder="请输入提现金额"
            size="large"
            style="width: 100%"
          >
            <template #prefix>¥</template>
          </n-input-number>
        </n-form-item>

        <n-form-item label="提现账户" path="account">
          <n-input
            v-model:value="withdrawForm.account"
            placeholder="请输入银行卡号或支付宝账号"
            size="large"
          >
            <template #prefix>
              <n-icon :component="CardOutline" color="#1890ff" />
            </template>
          </n-input>
        </n-form-item>

        <n-form-item label="提现备注" path="remark">
          <n-input
            v-model:value="withdrawForm.remark"
            type="textarea"
            placeholder="选填"
            :rows="2"
          />
        </n-form-item>
      </n-form>

      <template #footer>
        <div class="modal-footer">
          <n-button @click="showWithdrawModal = false">取消</n-button>
          <n-button type="primary" :loading="withdrawing" @click="handleWithdraw">
            提交申请
          </n-button>
        </div>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, h, onMounted } from 'vue'
import request from '@/utils/request'
import { useMessage } from 'naive-ui'
import type { FormInst, FormRules, DataTableColumns } from 'naive-ui'
import { NTag } from 'naive-ui'
import {
  WalletOutline,
  AddCircleOutline,
  ArrowUpOutline,
  SearchOutline,
  CardOutline
} from '@vicons/ionicons5'

const message = useMessage()

const balance = ref(0)
const monthlyRecharge = ref(0)
const monthlyExpense = ref(0)
const frozenAmount = ref(0)
const searchKeyword = ref('')
const activeType = ref('all')

const showRechargeModal = ref(false)
const showWithdrawModal = ref(false)
const recharging = ref(false)
const withdrawing = ref(false)

const rechargeFormRef = ref<FormInst | null>(null)
const withdrawFormRef = ref<FormInst | null>(null)

const quickAmounts = [50, 100, 200, 500, 1000, 2000]

const rechargeForm = ref({
  amount: null as number | null,
  payMethod: 'alipay'
})

const withdrawForm = ref({
  amount: null as number | null,
  account: '',
  remark: ''
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

interface RecordItem {
  id: number
  time: string
  type: string
  typeLabel: string
  amount: number
  balance: number
  remark: string
}

const records = ref<RecordItem[]>([])

const loading = ref(false)

const fetchWalletData = async () => {
  loading.value = true
  try {
    const [balanceRes, recordsRes] = await Promise.all([
      request.get('/api/v2/balance'),
      request.get('/api/v2/balance/logs', { params: { page: pagination.page, page_size: pagination.pageSize } })
    ])
    if (balanceRes?.data?.data) {
      const d = balanceRes.data.data
      balance.value = d.balance || 0
      monthlyRecharge.value = d.monthly_recharge || 0
      monthlyExpense.value = d.monthly_expense || 0
      frozenAmount.value = d.frozen || 0
    }
    if (recordsRes?.data?.data) {
      records.value = recordsRes.data.data.list || recordsRes.data.data || []
      pagination.itemCount = recordsRes.data.data.total || records.value.length
    }
  } catch (e) {
    console.error('Failed to fetch wallet data:', e)
  } finally {
    loading.value = false
  }
}

onMounted(fetchWalletData)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => { pagination.page = page },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
  }
})

const typeColorMap: Record<string, string> = {
  recharge: 'success',
  expense: 'error',
  refund: 'info',
  withdraw: 'warning'
}

const columns: DataTableColumns<RecordItem> = [
  { title: '时间', key: 'time', width: 180 },
  {
    title: '类型',
    key: 'type',
    width: 100,
    render(row) {
      return h(NTag, { type: typeColorMap[row.type] as any, size: 'small', bordered: false }, { default: () => row.typeLabel })
    }
  },
  {
    title: '金额',
    key: 'amount',
    width: 120,
    render(row) {
      const color = row.amount > 0 ? '#52c41a' : '#ff4d4f'
      const prefix = row.amount > 0 ? '+' : ''
      return h('span', { style: { color, fontWeight: 600 } }, `${prefix}${row.amount.toFixed(2)}`)
    }
  },
  {
    title: '余额',
    key: 'balance',
    width: 120,
    render(row) {
      return h('span', {}, `¥${row.balance.toFixed(2)}`)
    }
  },
  { title: '备注', key: 'remark', ellipsis: { tooltip: true } }
]

const filteredRecords = computed(() => {
  let list = records.value
  if (activeType.value !== 'all') {
    list = list.filter(r => r.type === activeType.value)
  }
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    list = list.filter(r => r.remark.toLowerCase().includes(kw))
  }
  return list
})

async function handleRecharge() {
  try {
    await rechargeFormRef.value?.validate()
    recharging.value = true
    const res = await request.post('/api/v2/balance/recharge', { amount: rechargeForm.value.amount, payment_method: rechargeForm.value.payMethod })
    balance.value += rechargeForm.value.amount || 0
    message.success('充值成功')
    showRechargeModal.value = false
    rechargeForm.value.amount = null
  } catch {
    message.error('请正确填写充值信息')
  } finally {
    recharging.value = false
  }
}

async function handleWithdraw() {
  try {
    await withdrawFormRef.value?.validate()
    if ((withdrawForm.value.amount || 0) > balance.value) {
      message.error('提现金额不能超过余额')
      return
    }
    withdrawing.value = true
    await request.post('/api/v2/balance/withdraw', { amount: withdrawForm.value.amount, method: withdrawForm.value.account })
    message.success('提现申请已提交')
    showWithdrawModal.value = false
    withdrawForm.value.amount = null
    withdrawForm.value.account = ''
    withdrawForm.value.remark = ''
  } catch {
    message.error('请正确填写提现信息')
  } finally {
    withdrawing.value = false
  }
}
</script>

<style scoped>
.wallet-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 24px;
}

.wallet-header-section {
  margin-bottom: 24px;
}

.balance-card {
  margin-bottom: 20px;
}

.balance-card-inner {
  background: linear-gradient(135deg, #1890ff, #096dd9);
  border-radius: 16px;
  padding: 36px 40px;
  color: #fff;
  position: relative;
  overflow: hidden;
}

.balance-info {
  position: relative;
  z-index: 1;
}

.balance-label {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 15px;
  opacity: 0.9;
  margin-bottom: 12px;
}

.balance-amount {
  display: flex;
  align-items: baseline;
  gap: 4px;
  margin-bottom: 24px;
}

.balance-amount .currency {
  font-size: 28px;
  font-weight: 600;
}

.balance-amount .amount {
  font-size: 48px;
  font-weight: 700;
  letter-spacing: 1px;
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
  background: rgba(255, 255, 255, 0.95);
  color: #1890ff;
  border: none;
}

.balance-actions .n-button:last-child {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.4);
}

.balance-decoration {
  position: absolute;
  top: 0;
  right: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
}

.deco-circle {
  position: absolute;
  border-radius: 50%;
  opacity: 0.1;
}

.deco-circle-1 {
  width: 300px;
  height: 300px;
  background: #fff;
  top: -100px;
  right: -50px;
}

.deco-circle-2 {
  width: 200px;
  height: 200px;
  background: #fff;
  bottom: -60px;
  right: 150px;
}

.deco-circle-3 {
  width: 120px;
  height: 120px;
  background: #fff;
  top: 20px;
  right: 250px;
  opacity: 0.06;
}

.stats-row {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
}

.stat-card {
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(24, 144, 255, 0.06);
  transition: box-shadow 0.3s, transform 0.3s;
}

.stat-card:hover {
  box-shadow: 0 4px 20px rgba(24, 144, 255, 0.12);
  transform: translateY(-2px);
}

.records-card {
  border-radius: 12px;
  box-shadow: 0 2px 12px rgba(24, 144, 255, 0.06);
}

.records-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.records-title {
  font-size: 16px;
  font-weight: 600;
  color: #262626;
}

.type-tabs {
  margin-bottom: 16px;
}

.records-table {
  margin-top: 8px;
}

.quick-amounts {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  margin-bottom: 16px;
}

.quick-amounts .n-button {
  min-width: 80px;
}

.pay-methods {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.pay-method-label {
  display: flex;
  align-items: center;
  gap: 8px;
}

.available-balance {
  font-size: 20px;
  font-weight: 600;
  color: #1890ff;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
}
</style>
