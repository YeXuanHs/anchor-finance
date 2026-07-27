<template>
  <div class="admin-page">
    <el-card shadow="hover">
      <el-tabs v-model="activeTab">
        <!-- Payment Methods -->
        <el-tab-pane label="支付方式" name="methods">
          <template #label>
            <span>支付方式</span>
          </template>
          <div class="tab-toolbar">
            <el-input v-model="searchKeyword" placeholder="搜索支付名称" clearable style="width: 220px">
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
            <el-select v-model="filterType" placeholder="支付类型" clearable style="width: 130px">
              <el-option v-for="o in typeOptions" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
            <el-button type="primary" @click="openPaymentModal()">
              <el-icon><Plus /></el-icon>添加支付方式
            </el-button>
          </div>
          <el-table :data="filteredPayments" v-loading="loading" stripe size="small">
            <el-table-column prop="id" label="ID" width="60" sortable />
            <el-table-column prop="name" label="名称" show-overflow-tooltip />
            <el-table-column prop="type" label="类型" width="100">
              <template #default="{ row }">
                <el-tag size="small">{{ typeNameMap[row.type] || row.type }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="feeRate" label="手续费率" width="100" sortable>
              <template #default="{ row }">
                <span style="font-weight: 600; color: #0056FF">{{ row.feeRate }}%</span>
              </template>
            </el-table-column>
            <el-table-column prop="totalAmount" label="累计金额" width="120" sortable>
              <template #default="{ row }">
                <span style="font-weight: 600; color: #52c41a">¥{{ row.totalAmount.toLocaleString() }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="totalCount" label="交易笔数" width="100" sortable />
            <el-table-column prop="sort" label="排序" width="70" sortable />
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-switch :model-value="row.status === 'active'" size="small" @change="handleTogglePayment(row)" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="130" fixed="right">
              <template #default="{ row }">
                <el-button text type="primary" :icon="Edit" @click="openPaymentModal(row)" />
                <el-popconfirm title="确认删除该支付方式？" @confirm="handleDeletePayment(row.id)">
                  <template #reference>
                    <el-button text type="danger" :icon="Delete" />
                  </template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- Payment Records -->
        <el-tab-pane label="支付记录" name="records">
          <div class="tab-toolbar">
            <el-input v-model="recordSearchKeyword" placeholder="搜索订单号/用户名" clearable style="width: 220px">
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
            <el-date-picker v-model="recordDateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" style="width: 280px" clearable />
          </div>
          <el-table :data="filteredRecords" v-loading="recordsLoading" stripe size="small">
            <el-table-column prop="id" label="ID" width="60" sortable />
            <el-table-column prop="orderNo" label="订单号" width="180" show-overflow-tooltip />
            <el-table-column prop="username" label="用户名" width="100" />
            <el-table-column prop="paymentName" label="支付方式" width="120" />
            <el-table-column prop="amount" label="金额" width="120" sortable>
              <template #default="{ row }">
                <span style="font-weight: 600; color: #52c41a">¥{{ row.amount.toFixed(2) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="fee" label="手续费" width="100">
              <template #default="{ row }">
                <span style="color: #f5a623">¥{{ row.fee.toFixed(2) }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="recordStatusMap[row.status]?.type as any" size="small">
                  {{ recordStatusMap[row.status]?.label }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="时间" width="180" sortable />
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- Payment Dialog -->
    <el-dialog v-model="paymentModalVisible" :title="editingPayment ? '编辑支付方式' : '添加支付方式'" width="640px" destroy-on-close>
      <el-form ref="paymentFormRef" :model="paymentForm" :rules="paymentRules" label-width="100px">
        <el-form-item label="支付名称" prop="name">
          <el-input v-model="paymentForm.name" placeholder="请输入支付名称" />
        </el-form-item>
        <el-form-item label="支付类型" prop="type">
          <el-select v-model="paymentForm.type" placeholder="选择支付类型" style="width: 100%">
            <el-option v-for="o in typeOptions" :key="o.value" :label="o.label" :value="o.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="手续费率(%)" prop="feeRate">
          <el-input-number v-model="paymentForm.feeRate" :min="0" :max="100" :precision="2" style="width: 100%">
            <template #suffix>%</template>
          </el-input-number>
        </el-form-item>
        <el-form-item label="App ID" prop="appId">
          <el-input v-model="paymentForm.appId" placeholder="支付网关 App ID" />
        </el-form-item>
        <el-form-item label="App Secret" prop="appSecret">
          <el-input v-model="paymentForm.appSecret" type="password" show-password placeholder="支付网关 App Secret" />
        </el-form-item>
        <el-form-item label="回调URL">
          <el-input v-model="paymentForm.notifyUrl" placeholder="支付回调通知URL" />
        </el-form-item>
        <el-form-item label="排序">
          <el-input-number v-model="paymentForm.sort" :min="0" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="paymentForm.status" active-value="active" inactive-value="inactive" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="paymentModalVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handlePaymentSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Plus, Edit, Delete } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'

const activeTab = ref('methods')
const loading = ref(false)
const recordsLoading = ref(false)
const submitting = ref(false)
const paymentModalVisible = ref(false)
const searchKeyword = ref('')
const filterType = ref<string | null>(null)
const paymentFormRef = ref<FormInstance>()
const editingPayment = ref<any>(null)
const recordSearchKeyword = ref('')
const recordDateRange = ref<any>(null)

const typeOptions = [
  { label: '支付宝', value: 'alipay' }, { label: '微信支付', value: 'wechat' },
  { label: 'PayPal', value: 'paypal' }, { label: 'Stripe', value: 'stripe' },
  { label: '虎皮椒', value: 'hupijiao' }, { label: '易支付', value: 'epay' },
  { label: 'USDT', value: 'usdt' }, { label: '余额', value: 'balance' },
]
const typeNameMap: Record<string, string> = { alipay: '支付宝', wechat: '微信支付', paypal: 'PayPal', stripe: 'Stripe', hupijiao: '虎皮椒', epay: '易支付', usdt: 'USDT', balance: '余额' }

const payments = ref([
  { id: 1, name: '支付宝', type: 'alipay', feeRate: 0.6, appId: '2021000000000001', status: 'active', sort: 1, totalAmount: 125680.50, totalCount: 1256 },
  { id: 2, name: '微信支付', type: 'wechat', feeRate: 0.6, appId: 'wx1234567890abcdef', status: 'active', sort: 2, totalAmount: 98450.00, totalCount: 1589 },
  { id: 3, name: 'PayPal', type: 'paypal', feeRate: 3.5, appId: 'AYSq3RDGrs99999999999999', status: 'active', sort: 3, totalAmount: 45200.75, totalCount: 234 },
  { id: 4, name: 'Stripe', type: 'stripe', feeRate: 2.9, appId: 'pk_live_1234567890', status: 'active', sort: 4, totalAmount: 78900.00, totalCount: 456 },
  { id: 5, name: '虎皮椒', type: 'hupijiao', feeRate: 1.0, appId: '100001', status: 'inactive', sort: 5, totalAmount: 12300.00, totalCount: 89 },
  { id: 6, name: '易支付', type: 'epay', feeRate: 0.5, appId: '1001', status: 'active', sort: 6, totalAmount: 56700.25, totalCount: 678 },
  { id: 7, name: 'USDT(TRC20)', type: 'usdt', feeRate: 0, appId: 'TRX1234567890', status: 'active', sort: 7, totalAmount: 234500.00, totalCount: 123 },
  { id: 8, name: '余额支付', type: 'balance', feeRate: 0, appId: '-', status: 'active', sort: 8, totalAmount: 45600.00, totalCount: 567 },
])

const filteredPayments = computed(() => {
  return payments.value.filter((p) => {
    if (searchKeyword.value.trim() && !p.name.toLowerCase().includes(searchKeyword.value.trim().toLowerCase())) return false
    if (filterType.value && p.type !== filterType.value) return false
    return true
  })
})

const records = ref([
  { id: 1, orderNo: 'ORD20260727001', username: 'user001', paymentName: '支付宝', amount: 299.00, fee: 1.79, status: 'success', createdAt: '2026-07-27 10:15:32' },
  { id: 2, orderNo: 'ORD20260727002', username: 'user002', paymentName: '微信支付', amount: 599.00, fee: 3.59, status: 'success', createdAt: '2026-07-27 10:20:45' },
  { id: 3, orderNo: 'ORD20260727003', username: 'user003', paymentName: 'PayPal', amount: 1299.00, fee: 45.47, status: 'success', createdAt: '2026-07-27 10:25:18' },
  { id: 4, orderNo: 'ORD20260727004', username: 'user004', paymentName: 'Stripe', amount: 89.00, fee: 2.58, status: 'pending', createdAt: '2026-07-27 10:30:56' },
  { id: 5, orderNo: 'ORD20260727005', username: 'user005', paymentName: '支付宝', amount: 399.00, fee: 2.39, status: 'failed', createdAt: '2026-07-27 10:35:22' },
  { id: 6, orderNo: 'ORD20260727006', username: 'user006', paymentName: 'USDT(TRC20)', amount: 199.00, fee: 0, status: 'success', createdAt: '2026-07-27 10:40:11' },
])

const filteredRecords = computed(() => {
  return records.value.filter((r) => {
    if (recordSearchKeyword.value.trim()) {
      const kw = recordSearchKeyword.value.trim().toLowerCase()
      if (!r.orderNo.toLowerCase().includes(kw) && !r.username.toLowerCase().includes(kw)) return false
    }
    return true
  })
})

const recordStatusMap: Record<string, { label: string; type: string }> = {
  success: { label: '成功', type: 'success' }, pending: { label: '待处理', type: 'warning' },
  failed: { label: '失败', type: 'danger' }, refunded: { label: '已退款', type: 'info' },
}

const paymentForm = reactive({ name: '', type: '' as string, feeRate: 0, appId: '', appSecret: '', notifyUrl: '', sort: 0, status: 'active' })
const paymentRules: FormRules = {
  name: { required: true, message: '请输入支付名称', trigger: 'blur' },
  type: { required: true, message: '请选择支付类型', trigger: 'change' },
  feeRate: { required: true, message: '请输入手续费率', trigger: 'blur' },
  appId: { required: true, message: '请输入App ID', trigger: 'blur' },
}

function openPaymentModal(payment?: any) {
  editingPayment.value = payment || null
  if (payment) {
    Object.assign(paymentForm, { name: payment.name, type: payment.type, feeRate: payment.feeRate, appId: payment.appId, appSecret: '', notifyUrl: '', sort: payment.sort, status: payment.status })
  } else {
    Object.assign(paymentForm, { name: '', type: '', feeRate: 0, appId: '', appSecret: '', notifyUrl: '', sort: 0, status: 'active' })
  }
  paymentModalVisible.value = true
}

async function handlePaymentSubmit() {
  if (!paymentFormRef.value) return
  try { await paymentFormRef.value.validate() } catch { return }
  submitting.value = true
  try {
    ElMessage.success(editingPayment.value ? '支付方式更新成功' : '支付方式添加成功')
    paymentModalVisible.value = false
  } finally { submitting.value = false }
}

function handleTogglePayment(payment: any) {
  payment.status = payment.status === 'active' ? 'inactive' : 'active'
  ElMessage.success(`支付方式「${payment.name}」已${payment.status === 'active' ? '启用' : '禁用'}`)
}

function handleDeletePayment(id: number) {
  payments.value = payments.value.filter((p) => p.id !== id)
  ElMessage.success('支付方式已删除')
}
</script>

<style scoped>
.tab-toolbar {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 16px;
}
</style>
