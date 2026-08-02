<template>
  <div class="invoice-detail-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>账单详情</span>
          <div class="header-actions">
            <el-button @click="handleBack">
              <el-icon><Back /></el-icon>
              返回
            </el-button>
          </div>
        </div>
      </template>

      <div v-loading="loading" class="loading-container">
        <div v-if="bill" class="bill-content">
          <!-- 账单基本信息 -->
          <el-descriptions :column="2" border class="bill-info">
            <el-descriptions-item label="账单编号">{{ bill.bill_no }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="statusTypeMap[bill.status]" size="default">
                {{ statusLabelMap[bill.status] }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="关联客户">
              <el-button type="primary" link @click="handleViewClient">
                {{ bill.client_username }}
              </el-button>
            </el-descriptions-item>
            <el-descriptions-item label="关联产品">{{ bill.product_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="账单金额">
              <span class="amount-text">¥{{ formatAmount(bill.amount) }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="已付金额">
              <span class="amount-text">¥{{ formatAmount(bill.paid_amount) }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="待付金额">
              <span class="amount-text" :class="{ 'text-danger': bill.amount - bill.paid_amount > 0 }">
                ¥{{ formatAmount(bill.amount - bill.paid_amount) }}
              </span>
            </el-descriptions-item>
            <el-descriptions-item label="优惠金额">¥{{ formatAmount(bill.discount_amount || 0) }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ bill.created_at }}</el-descriptions-item>
            <el-descriptions-item label="到期时间">{{ bill.due_date }}</el-descriptions-item>
            <el-descriptions-item label="支付时间">{{ bill.paid_at || '-' }}</el-descriptions-item>
            <el-descriptions-item label="备注" :span="2">{{ bill.remark || '无' }}</el-descriptions-item>
          </el-descriptions>

          <!-- 账单明细列表 -->
          <div class="section">
            <h3>账单明细</h3>
            <el-table :data="bill.items || []" style="width: 100%" border>
              <el-table-column prop="id" label="ID" width="80" align="center" />
              <el-table-column prop="description" label="描述" min-width="200" />
              <el-table-column prop="quantity" label="数量" width="100" align="center" />
              <el-table-column prop="unit_price" label="单价" width="120" align="right">
                <template #default="{ row }">
                  ¥{{ formatAmount(row.unit_price) }}
                </template>
              </el-table-column>
              <el-table-column prop="discount" label="折扣" width="100" align="right">
                <template #default="{ row }">
                  {{ row.discount || '0%' }}
                </template>
              </el-table-column>
              <el-table-column prop="amount" label="金额" width="120" align="right">
                <template #default="{ row }">
                  <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
                </template>
              </el-table-column>
            </el-table>
            <div class="bill-total">
              <span>合计: </span>
              <span class="amount-text">¥{{ formatAmount(bill.amount) }}</span>
            </div>
          </div>

          <!-- 支付记录 -->
          <div class="section">
            <h3>支付记录</h3>
            <el-table :data="payments" v-loading="paymentsLoading" style="width: 100%" border>
              <el-table-column prop="id" label="ID" width="80" align="center" />
              <el-table-column prop="payment_no" label="支付流水号" width="200" />
              <el-table-column prop="gateway" label="支付网关" width="120" />
              <el-table-column prop="amount" label="支付金额" width="120" align="right">
                <template #default="{ row }">
                  <span class="text-success">¥{{ formatAmount(row.amount) }}</span>
                </template>
              </el-table-column>
              <el-table-column prop="status" label="状态" width="100">
                <template #default="{ row }">
                  <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
                    {{ row.status === 1 ? '成功' : '失败' }}
                  </el-tag>
                </template>
              </el-table-column>
              <el-table-column prop="created_at" label="支付时间" width="180" />
              <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
            </el-table>
            <el-empty v-if="!payments.length && !paymentsLoading" description="暂无支付记录" />
          </div>
        </div>
      </div>
    </el-card>

    <!-- 底部操作栏 -->
    <div class="action-bar" v-if="bill">
      <el-button
        v-if="bill.status === 0 || bill.status === 1"
        type="primary"
        @click="handleSend"
        :loading="actionLoading"
      >
        <el-icon><Promotion /></el-icon>
        发送账单
      </el-button>
      <el-button
        v-if="bill.status === 2 || bill.status === 3"
        type="success"
        @click="handleRefund"
        :loading="actionLoading"
      >
        <el-icon><RefreshRight /></el-icon>
        退款
      </el-button>
      <el-popconfirm
        v-if="bill.status !== 4 && bill.status !== 5"
        title="确定取消该账单吗？"
        @confirm="handleCancel"
      >
        <template #reference>
          <el-button type="danger" :loading="actionLoading">
            <el-icon><CircleClose /></el-icon>
            取消账单
          </el-button>
        </template>
      </el-popconfirm>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Back, Promotion, RefreshRight, CircleClose } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const paymentsLoading = ref(false)
const actionLoading = ref(false)
const bill = ref<any>({})
const payments = ref<any[]>([])

const statusTypeMap: Record<number, string> = {
  0: 'info',
  1: 'warning',
  2: 'success',
  3: 'success',
  4: 'danger'
}

const statusLabelMap: Record<number, string> = {
  0: '待支付',
  1: '已发送',
  2: '已支付',
  3: '部分支付',
  4: '已取消'
}

const fetchBill = async () => {
  const id = route.params.id
  if (!id) return

  loading.value = true
  try {
    const data = await request.get({ url: `/api/admin/bills/${id}` })
    bill.value = data.bill || data
  } catch (error) {
    console.error('获取账单详情失败:', error)
    ElMessage.error('获取账单详情失败')
  } finally {
    loading.value = false
  }
}

const fetchPayments = async () => {
  const id = route.params.id
  if (!id) return

  paymentsLoading.value = true
  try {
    const data = await request.get({ url: `/api/admin/bills/${id}/payments` })
    payments.value = data || []
  } catch (error) {
    console.error('获取支付记录失败:', error)
  } finally {
    paymentsLoading.value = false
  }
}

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

const handleBack = () => {
  router.back()
}

const handleViewClient = () => {
  if (bill.value?.client_id) {
    router.push(`/finance/clients/detail/${bill.value.client_id}`)
  }
}

const handleSend = async () => {
  const id = route.params.id
  if (!id) return

  actionLoading.value = true
  try {
    await request.post({ url: `/api/admin/bills/${id}/send` })
    ElMessage.success('账单发送成功')
    fetchBill()
  } catch (error) {
    ElMessage.error('发送失败')
  } finally {
    actionLoading.value = false
  }
}

const handleRefund = async () => {
  const id = route.params.id
  if (!id) return

  try {
    await ElMessageBox.confirm('确定要对该账单进行退款吗？', '退款确认', {
      type: 'warning'
    })

    actionLoading.value = true
    await request.post({ url: `/api/admin/bills/${id}/refund` })
    ElMessage.success('退款成功')
    fetchBill()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('退款失败')
    }
  } finally {
    actionLoading.value = false
  }
}

const handleCancel = async () => {
  const id = route.params.id
  if (!id) return

  actionLoading.value = true
  try {
    await request.post({ url: `/api/admin/bills/${id}/cancel` })
    ElMessage.success('账单已取消')
    fetchBill()
  } catch (error) {
    ElMessage.error('取消失败')
  } finally {
    actionLoading.value = false
  }
}

onMounted(() => {
  fetchBill()
  fetchPayments()
})
</script>

<style scoped lang="scss">
.invoice-detail-page {
  padding: 20px;
  padding-bottom: 80px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.loading-container {
  min-height: 400px;
}

.bill-info {
  margin-bottom: 24px;
}

.amount-text {
  font-weight: 600;
  color: var(--el-color-primary);
}

.text-success {
  color: var(--el-color-success);
}

.text-danger {
  color: var(--el-color-danger);
}

.section {
  margin-top: 24px;

  h3 {
    margin: 0 0 16px;
    font-size: 16px;
    font-weight: 600;
  }

  .bill-total {
    text-align: right;
    margin-top: 16px;
    padding: 12px 16px;
    background: var(--el-fill-color-lighter);
    border-radius: 4px;
    font-size: 16px;
  }
}

.action-bar {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 16px 20px;
  background: var(--el-bg-color);
  border-top: 1px solid var(--el-border-color-lighter);
  display: flex;
  gap: 12px;
  z-index: 100;
}
</style>