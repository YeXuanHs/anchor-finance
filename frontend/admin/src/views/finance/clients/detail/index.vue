<template>
  <div class="client-detail-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>客户详情</span>
          <div class="header-actions">
            <el-button @click="handleBack">
              <el-icon><Back /></el-icon>
              返回
            </el-button>
            <el-button type="primary" @click="handleEdit">
              <el-icon><Edit /></el-icon>
              编辑
            </el-button>
          </div>
        </div>
      </template>

      <div v-loading="loading" class="loading-container">
        <div v-if="client" class="client-content">
          <!-- 客户基本信息摘要 -->
          <div class="client-summary">
            <div class="summary-left">
              <div class="avatar">
                <el-avatar :size="80">{{ client.username?.charAt(0) }}</el-avatar>
              </div>
              <div class="info">
                <h2>{{ client.username }}</h2>
                <p class="email">{{ client.email }}</p>
                <p class="phone" v-if="client.phone">{{ client.phone }}</p>
                <p class="group">
                  <el-tag size="small">{{ client.group_name || '默认分组' }}</el-tag>
                  <el-tag :type="client.status === 1 ? 'success' : 'danger'" size="small">
                    {{ client.status === 1 ? '正常' : '禁用' }}
                  </el-tag>
                </p>
              </div>
            </div>
            <div class="summary-right">
              <div class="stat-item">
                <div class="stat-label">余额</div>
                <div class="stat-value">¥{{ formatAmount(client.balance) }}</div>
              </div>
              <div class="stat-item">
                <div class="stat-label">信用额度</div>
                <div class="stat-value">¥{{ formatAmount(client.credit) }}</div>
              </div>
              <div class="stat-item">
                <div class="stat-label">注册时间</div>
                <div class="stat-value">{{ client.created_at }}</div>
              </div>
              <div class="stat-item">
                <div class="stat-label">最后登录</div>
                <div class="stat-value">{{ client.last_login_at || '未登录' }}</div>
              </div>
            </div>
          </div>

          <!-- 标签页 -->
          <el-tabs v-model="activeTab" @tab-change="handleTabChange">
            <!-- 资料标签页 -->
            <el-tab-pane label="资料" name="profile">
              <div class="tab-content">
                <el-form :model="client" label-width="120px" class="profile-form">
                  <el-row :gutter="20">
                    <el-col :span="12">
                      <el-form-item label="用户名">{{ client.username }}</el-form-item>
                      <el-form-item label="邮箱">{{ client.email }}</el-form-item>
                      <el-form-item label="手机号">{{ client.phone || '-' }}</el-form-item>
                      <el-form-item label="公司">{{ client.company || '-' }}</el-form-item>
                      <el-form-item label="职位">{{ client.position || '-' }}</el-form-item>
                    </el-col>
                    <el-col :span="12">
                      <el-form-item label="姓名">{{ client.fullname || '-' }}</el-form-item>
                      <el-form-item label="地址">{{ client.address || '-' }}</el-form-item>
                      <el-form-item label="城市">{{ client.city || '-' }}</el-form-item>
                      <el-form-item label="省份">{{ client.state || '-' }}</el-form-item>
                      <el-form-item label="邮编">{{ client.postcode || '-' }}</el-form-item>
                    </el-col>
                  </el-row>
                  <el-form-item label="国家">{{ client.country || '-' }}</el-form-item>
                  <el-form-item label="时区">{{ client.timezone || '默认' }}</el-form-item>
                  <el-form-item label="语言">{{ client.language || 'zh-CN' }}</el-form-item>
                </el-form>
              </div>
            </el-tab-pane>

            <!-- 服务标签页 -->
            <el-tab-pane label="服务" name="services">
              <div class="tab-content">
                <el-table :data="services" v-loading="servicesLoading" style="width: 100%">
                  <el-table-column prop="id" label="ID" width="80" />
                  <el-table-column prop="product_name" label="产品" min-width="150" />
                  <el-table-column prop="domain" label="域名" min-width="200" />
                  <el-table-column prop="billing_cycle" label="周期" width="100" />
                  <el-table-column prop="amount" label="金额" width="120">
                    <template #default="{ row }">
                      ¥{{ formatAmount(row.amount) }}
                    </template>
                  </el-table-column>
                  <el-table-column prop="status" label="状态" width="100">
                    <template #default="{ row }">
                      <el-tag :type="getServiceStatusType(row.status)" size="small">
                        {{ getServiceStatusText(row.status) }}
                      </el-tag>
                    </template>
                  </el-table-column>
                  <el-table-column prop="created_at" label="创建时间" width="180" />
                  <el-table-column prop="due_date" label="到期时间" width="180" />
                </el-table>
                <el-empty v-if="!services.length && !servicesLoading" description="暂无服务" />
              </div>
            </el-tab-pane>

            <!-- 账单标签页 -->
            <el-tab-pane label="账单" name="bills">
              <div class="tab-content">
                <el-table :data="bills" v-loading="billsLoading" style="width: 100%">
                  <el-table-column prop="id" label="ID" width="80" />
                  <el-table-column prop="bill_no" label="账单号" width="180">
                    <template #default="{ row }">
                      <el-button type="primary" link @click="handleViewBill(row)">
                        {{ row.bill_no }}
                      </el-button>
                    </template>
                  </el-table-column>
                  <el-table-column prop="amount" label="金额" width="120">
                    <template #default="{ row }">
                      ¥{{ formatAmount(row.amount) }}
                    </template>
                  </el-table-column>
                  <el-table-column prop="paid_amount" label="已付金额" width="120">
                    <template #default="{ row }">
                      ¥{{ formatAmount(row.paid_amount) }}
                    </template>
                  </el-table-column>
                  <el-table-column prop="status" label="状态" width="100">
                    <template #default="{ row }">
                      <el-tag :type="getBillStatusType(row.status)" size="small">
                        {{ getBillStatusText(row.status) }}
                      </el-tag>
                    </template>
                  </el-table-column>
                  <el-table-column prop="created_at" label="创建时间" width="180" />
                  <el-table-column prop="due_date" label="到期时间" width="180" />
                  <el-table-column label="操作" width="150">
                    <template #default="{ row }">
                      <el-button type="primary" link @click="handleViewBill(row)">查看</el-button>
                      <el-button
                        v-if="row.status === 0 || row.status === 1"
                        type="success"
                        link
                        @click="handleSendBill(row)"
                      >
                        发送
                      </el-button>
                    </template>
                  </el-table-column>
                </el-table>
                <el-empty v-if="!bills.length && !billsLoading" description="暂无账单" />
              </div>
            </el-tab-pane>

            <!-- 交易标签页 -->
            <el-tab-pane label="交易" name="transactions">
              <div class="tab-content">
                <el-table :data="transactions" v-loading="transactionsLoading" style="width: 100%">
                  <el-table-column prop="id" label="ID" width="80" />
                  <el-table-column prop="transaction_no" label="交易号" width="200" />
                  <el-table-column prop="type" label="类型" width="100">
                    <template #default="{ row }">
                      <el-tag :type="getTransactionType(row.type)" size="small">
                        {{ getTransactionTypeText(row.type) }}
                      </el-tag>
                    </template>
                  </el-table-column>
                  <el-table-column prop="amount" label="金额" width="120">
                    <template #default="{ row }">
                      <span :class="row.amount >= 0 ? 'text-success' : 'text-danger'">
                        {{ row.amount >= 0 ? '+' : '' }}¥{{ formatAmount(row.amount) }}
                      </span>
                    </template>
                  </el-table-column>
                  <el-table-column prop="balance" label="余额" width="120">
                    <template #default="{ row }">
                      ¥{{ formatAmount(row.balance) }}
                    </template>
                  </el-table-column>
                  <el-table-column prop="gateway" label="支付网关" width="120" />
                  <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
                  <el-table-column prop="created_at" label="时间" width="180" />
                </el-table>
                <el-empty v-if="!transactions.length && !transactionsLoading" description="暂无交易记录" />
              </div>
            </el-tab-pane>

            <!-- 信用标签页 -->
            <el-tab-pane label="信用" name="credit">
              <div class="tab-content">
                <el-card shadow="never">
                  <template #header>
                    <span>信用信息</span>
                  </template>
                  <el-form :model="client" label-width="120px">
                    <el-row :gutter="20">
                      <el-col :span="12">
                        <el-form-item label="当前余额">
                          <span class="amount-highlight">¥{{ formatAmount(client.balance) }}</span>
                        </el-form-item>
                        <el-form-item label="信用额度">
                          <span class="amount-highlight">¥{{ formatAmount(client.credit) }}</span>
                        </el-form-item>
                        <el-form-item label="可用额度">
                          <span class="amount-highlight">¥{{ formatAmount(client.credit - client.balance) }}</span>
                        </el-form-item>
                      </el-col>
                      <el-col :span="12">
                        <el-form-item label="总充值">
                          <span class="amount-highlight">¥{{ formatAmount(client.total_deposit) }}</span>
                        </el-form-item>
                        <el-form-item label="总消费">
                          <span class="amount-highlight">¥{{ formatAmount(client.total_spend) }}</span>
                        </el-form-item>
                        <el-form-item label="总退款">
                          <span class="amount-highlight">¥{{ formatAmount(client.total_refund) }}</span>
                        </el-form-item>
                      </el-col>
                    </el-row>
                  </el-form>
                </el-card>

                <el-card shadow="never" style="margin-top: 16px;">
                  <template #header>
                    <span>信用调整</span>
                  </template>
                  <el-form :model="creditForm" label-width="100px">
                    <el-form-item label="调整类型">
                      <el-radio-group v-model="creditForm.type">
                        <el-radio value="deposit">充值</el-radio>
                        <el-radio value="refund">退款</el-radio>
                        <el-radio value="deduct">扣款</el-radio>
                      </el-radio-group>
                    </el-form-item>
                    <el-form-item label="金额">
                      <el-input-number v-model="creditForm.amount" :min="0.01" :precision="2" />
                    </el-form-item>
                    <el-form-item label="备注">
                      <el-input v-model="creditForm.remark" type="textarea" :rows="2" placeholder="请输入备注" />
                    </el-form-item>
                    <el-form-item>
                      <el-button type="primary" @click="handleCreditAdjust" :loading="creditLoading">
                        提交调整
                      </el-button>
                    </el-form-item>
                  </el-form>
                </el-card>
              </div>
            </el-tab-pane>

            <!-- 工单标签页 -->
            <el-tab-pane label="工单" name="tickets">
              <div class="tab-content">
                <el-table :data="tickets" v-loading="ticketsLoading" style="width: 100%">
                  <el-table-column prop="id" label="ID" width="80" />
                  <el-table-column prop="subject" label="标题" min-width="200" show-overflow-tooltip />
                  <el-table-column prop="department_name" label="部门" width="120" />
                  <el-table-column prop="priority" label="优先级" width="80">
                    <template #default="{ row }">
                      <el-tag :type="getPriorityType(row.priority)" size="small">
                        {{ getPriorityText(row.priority) }}
                      </el-tag>
                    </template>
                  </el-table-column>
                  <el-table-column prop="status" label="状态" width="100">
                    <template #default="{ row }">
                      <el-tag :type="getTicketStatusType(row.status)" size="small">
                        {{ getTicketStatusText(row.status) }}
                      </el-tag>
                    </template>
                  </el-table-column>
                  <el-table-column prop="created_at" label="创建时间" width="180" />
                  <el-table-column label="操作" width="100">
                    <template #default="{ row }">
                      <el-button type="primary" link @click="handleViewTicket(row)">查看</el-button>
                    </template>
                  </el-table-column>
                </el-table>
                <el-empty v-if="!tickets.length && !ticketsLoading" description="暂无工单" />
              </div>
            </el-tab-pane>

            <!-- 日志标签页 -->
            <el-tab-pane label="日志" name="logs">
              <div class="tab-content">
                <el-table :data="logs" v-loading="logsLoading" style="width: 100%">
                  <el-table-column prop="id" label="ID" width="80" />
                  <el-table-column prop="action" label="操作" width="120" />
                  <el-table-column prop="description" label="描述" min-width="300" show-overflow-tooltip />
                  <el-table-column prop="ip" label="IP地址" width="150" />
                  <el-table-column prop="created_at" label="时间" width="180" />
                </el-table>
                <el-empty v-if="!logs.length && !logsLoading" description="暂无日志" />
              </div>
            </el-tab-pane>

            <!-- 推介标签页 -->
            <el-tab-pane label="推介" name="referrals">
              <div class="tab-content">
                <el-card shadow="never">
                  <template #header>
                    <span>推介信息</span>
                  </template>
                  <el-form :model="client" label-width="120px">
                    <el-form-item label="推介码">
                      <el-input v-model="client.referral_code" readonly>
                        <template #append>
                          <el-button @click="handleCopyCode">复制</el-button>
                        </template>
                      </el-input>
                    </el-form-item>
                    <el-form-item label="推介链接">
                      <el-input :model-value="referralLink" readonly>
                        <template #append>
                          <el-button @click="handleCopyLink">复制</el-button>
                        </template>
                      </el-input>
                    </el-form-item>
                    <el-form-item label="推介人数">
                      <span>{{ client.referral_count || 0 }} 人</span>
                    </el-form-item>
                    <el-form-item label="推介佣金">
                      <span class="amount-highlight">¥{{ formatAmount(client.referral_commission) }}</span>
                    </el-form-item>
                  </el-form>
                </el-card>

                <el-card shadow="never" style="margin-top: 16px;">
                  <template #header>
                    <span>推介记录</span>
                  </template>
                  <el-table :data="referrals" v-loading="referralsLoading" style="width: 100%">
                    <el-table-column prop="id" label="ID" width="80" />
                    <el-table-column prop="referred_username" label="被推介用户" width="150" />
                    <el-table-column prop="commission" label="佣金" width="120">
                      <template #default="{ row }">
                        ¥{{ formatAmount(row.commission) }}
                      </template>
                    </el-table-column>
                    <el-table-column prop="status" label="状态" width="100">
                      <template #default="{ row }">
                        <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                          {{ row.status === 1 ? '已结算' : '待结算' }}
                        </el-tag>
                      </template>
                    </el-table-column>
                    <el-table-column prop="created_at" label="时间" width="180" />
                  </el-table>
                  <el-empty v-if="!referrals.length && !referralsLoading" description="暂无推介记录" />
                </el-card>
              </div>
            </el-tab-pane>
          </el-tabs>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Back, Edit } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const activeTab = ref('profile')
const client = ref<any>({})

const services = ref<any[]>([])
const bills = ref<any[]>([])
const transactions = ref<any[]>([])
const tickets = ref<any[]>([])
const logs = ref<any[]>([])
const referrals = ref<any[]>([])

const servicesLoading = ref(false)
const billsLoading = ref(false)
const transactionsLoading = ref(false)
const ticketsLoading = ref(false)
const logsLoading = ref(false)
const referralsLoading = ref(false)
const creditLoading = ref(false)

const creditForm = reactive({
  type: 'deposit',
  amount: 0,
  remark: ''
})

const referralLink = computed(() => {
  return `${window.location.origin}/register?ref=${client.value?.referral_code || ''}`
})

const fetchClient = async () => {
  const id = route.params.id
  if (!id) return
  loading.value = true
  try {
    const data = await request.get({ url: `/api/admin/users/${id}` })
    client.value = data
  } catch (error) {
    console.error('获取客户详情失败:', error)
    ElMessage.error('获取客户详情失败')
  } finally {
    loading.value = false
  }
}

const fetchTabData = async (tab: string) => {
  const id = route.params.id
  if (!id) return

  const fetchMap: Record<string, () => void> = {
    services: async () => {
      servicesLoading.value = true
      try {
        services.value = await request.get({ url: `/api/admin/client-services?user_id=${id}` }) || []
      } finally { servicesLoading.value = false }
    },
    bills: async () => {
      billsLoading.value = true
      try {
        bills.value = await request.get({ url: `/api/admin/invoices?user_id=${id}` }) || []
      } finally { billsLoading.value = false }
    },
    transactions: async () => {
      transactionsLoading.value = true
      try {
        transactions.value = await request.get({ url: `/api/admin/accounts?user_id=${id}` }) || []
      } finally { transactionsLoading.value = false }
    },
    tickets: async () => {
      ticketsLoading.value = true
      try {
        tickets.value = await request.get({ url: `/api/admin/tickets?user_id=${id}` }) || []
      } finally { ticketsLoading.value = false }
    },
    logs: async () => {
      logsLoading.value = true
      try {
        logs.value = await request.get({ url: `/api/admin/log-records?user_id=${id}` }) || []
      } finally { logsLoading.value = false }
    },
    referrals: async () => {
      referralsLoading.value = true
      try {
        referrals.value = await request.get({ url: `/api/admin/affiliate?user_id=${id}` }) || []
      } finally { referralsLoading.value = false }
    }
  }

  if (fetchMap[tab]) fetchMap[tab]()
}

const handleTabChange = (tab: string) => {
  fetchTabData(tab)
}

const handleBack = () => {
  router.back()
}

const handleEdit = () => {
  router.push(`/finance/clients/list?edit=${route.params.id}`)
}

const handleViewBill = (row: any) => {
  router.push(`/finance/orders/invoice-detail/${row.id}`)
}

const handleSendBill = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/invoices/${row.id}/email` })
    ElMessage.success('账单发送成功')
    fetchTabData('bills')
  } catch (error) {
    ElMessage.error('发送失败')
  }
}

const handleViewTicket = (row: any) => {
  router.push(`/finance/tickets/detail/${row.id}`)
}

const handleCreditAdjust = async () => {
  const id = route.params.id
  if (!id) return
  if (creditForm.amount <= 0) {
    ElMessage.error('请输入正确的金额')
    return
  }
  creditLoading.value = true
  try {
    await request.post({ url: `/api/admin/credit/users/${id}/adjust`, params: creditForm })
    ElMessage.success('信用调整成功')
    fetchClient()
    creditForm.amount = 0
    creditForm.remark = ''
  } catch (error) {
    ElMessage.error('调整失败')
  } finally {
    creditLoading.value = false
  }
}

const handleCopyCode = () => {
  if (client.value?.referral_code) {
    navigator.clipboard.writeText(client.value.referral_code)
    ElMessage.success('推介码已复制')
  }
}

const handleCopyLink = () => {
  navigator.clipboard.writeText(referralLink.value)
  ElMessage.success('推介链接已复制')
}

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

const getServiceStatusType = (status: string) => {
  const map: Record<string, string> = { active: 'success', suspended: 'warning', terminated: 'danger', pending: 'info' }
  return (map[status] || 'info') as any
}

const getServiceStatusText = (status: string) => {
  const map: Record<string, string> = { active: '活跃', suspended: '暂停', terminated: '已终止', pending: '待处理' }
  return map[status] || status
}

const getBillStatusType = (status: number) => {
  const map: Record<number, string> = { 0: 'info', 1: 'warning', 2: 'success', 3: 'success', 4: 'danger' }
  return (map[status] || 'info') as any
}

const getBillStatusText = (status: number) => {
  const map: Record<number, string> = { 0: '待支付', 1: '已发送', 2: '已支付', 3: '部分支付', 4: '已取消' }
  return map[status] || '未知'
}

const getTransactionType = (type: string) => {
  const map: Record<string, string> = { deposit: 'success', payment: 'warning', refund: 'info', deduction: 'danger' }
  return (map[type] || 'info') as any
}

const getTransactionTypeText = (type: string) => {
  const map: Record<string, string> = { deposit: '充值', payment: '支付', refund: '退款', deduction: '扣款' }
  return map[type] || type
}

const getPriorityType = (priority: string) => {
  const map: Record<string, string> = { low: 'info', medium: '', high: 'warning', urgent: 'danger' }
  return (map[priority] || 'info') as any
}

const getPriorityText = (priority: string) => {
  const map: Record<string, string> = { low: '低', medium: '中', high: '高', urgent: '紧急' }
  return map[priority] || priority
}

const getTicketStatusType = (status: string) => {
  const map: Record<string, string> = { open: 'warning', in_progress: '', replied: 'success', closed: 'info' }
  return (map[status] || 'info') as any
}

const getTicketStatusText = (status: string) => {
  const map: Record<string, string> = { open: '待处理', in_progress: '处理中', replied: '已回复', closed: '已关闭' }
  return map[status] || status
}

onMounted(() => {
  fetchClient()
})
</script>

<style scoped lang="scss">
.client-detail-page {
  padding: 20px;
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

.client-summary {
  display: flex;
  justify-content: space-between;
  padding: 20px;
  background: var(--el-fill-color-lighter);
  border-radius: 8px;
  margin-bottom: 20px;

  .summary-left {
    display: flex;
    gap: 20px;

    .info {
      h2 {
        margin: 0 0 8px;
        font-size: 24px;
      }

      .email, .phone {
        margin: 0 0 4px;
        color: var(--el-text-color-secondary);
      }

      .group {
        margin: 8px 0 0;
        display: flex;
        gap: 8px;
      }
    }
  }

  .summary-right {
    display: flex;
    gap: 40px;

    .stat-item {
      text-align: center;

      .stat-label {
        color: var(--el-text-color-secondary);
        font-size: 14px;
        margin-bottom: 4px;
      }

      .stat-value {
        font-size: 18px;
        font-weight: 600;
        color: var(--el-text-color-primary);
      }
    }
  }
}

.tab-content {
  padding: 20px 0;
}

.profile-form {
  max-width: 800px;
}

.amount-highlight {
  font-size: 18px;
  font-weight: 600;
  color: var(--el-color-primary);
}

.text-success {
  color: var(--el-color-success);
}

.text-danger {
  color: var(--el-color-danger);
}
</style>