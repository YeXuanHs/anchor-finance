<template>
  <div class="customer-detail-page">
    <!-- 顶部客户信息卡片 -->
    <el-card shadow="never" class="customer-info-card">
      <el-row :gutter="20">
        <el-col :span="6">
          <div class="customer-basic">
            <el-avatar :size="64" :src="customer.avatar">
              {{ customer.username?.charAt(0) }}
            </el-avatar>
            <div class="customer-name">
              <h3>{{ customer.username }}</h3>
              <el-tag :type="getStatusType(customer.status)" size="small">
                {{ getStatusText(customer.status) }}
              </el-tag>
            </div>
          </div>
        </el-col>
        <el-col :span="18">
          <el-row :gutter="20">
            <el-col :span="6">
              <div class="stat-item">
                <div class="stat-label">余额</div>
                <div class="stat-value">¥{{ formatMoney(customer.balance) }}</div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="stat-item">
                <div class="stat-label">信用额</div>
                <div class="stat-value">¥{{ formatMoney(customer.credit) }}</div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="stat-item">
                <div class="stat-label">订单数</div>
                <div class="stat-value">{{ customer.orders_count || 0 }}</div>
              </div>
            </el-col>
            <el-col :span="6">
              <div class="stat-item">
                <div class="stat-label">工单数</div>
                <div class="stat-value">{{ customer.tickets_count || 0 }}</div>
              </div>
            </el-col>
          </el-row>
        </el-col>
      </el-row>
    </el-card>

    <!-- 标签页 -->
    <el-card shadow="never" class="detail-card">
      <el-tabs v-model="activeTab">
        <!-- 概览 -->
        <el-tab-pane label="概览" name="overview">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="客户ID">{{ customer.id }}</el-descriptions-item>
            <el-descriptions-item label="用户名">{{ customer.username }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ customer.email }}</el-descriptions-item>
            <el-descriptions-item label="手机号">{{ customer.phone || '-' }}</el-descriptions-item>
            <el-descriptions-item label="公司">{{ customer.company || '-' }}</el-descriptions-item>
            <el-descriptions-item label="客户组">{{ customer.group_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="注册时间">{{ customer.created_at }}</el-descriptions-item>
            <el-descriptions-item label="最后登录">{{ customer.last_login_at || '-' }}</el-descriptions-item>
            <el-descriptions-item label="备注" :span="2">{{ customer.notes || '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>

        <!-- 产品/服务 -->
        <el-tab-pane label="产品/服务" name="products">
          <el-table :data="products" border stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="product_name" label="产品名称" min-width="150" />
            <el-table-column prop="domain" label="域名/标识" width="200" />
            <el-table-column prop="status" label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getProductStatusType(row.status)" size="small">
                  {{ getProductStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="billing_cycle" label="计费周期" width="100" />
            <el-table-column prop="amount" label="金额" width="100" align="right">
              <template #default="{ row }">¥{{ formatMoney(row.amount) }}</template>
            </el-table-column>
            <el-table-column prop="next_due_date" label="下次到期" width="120" />
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <el-button type="primary" link size="small">查看</el-button>
                <el-button type="warning" link size="small">管理</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 订单 -->
        <el-tab-pane label="订单" name="orders">
          <el-table :data="orders" border stripe>
            <el-table-column prop="order_no" label="订单号" width="150">
              <template #default="{ row }">
                <el-button type="primary" link @click="$router.push(`/order-detail/${row.id}`)">
                  {{ row.order_no }}
                </el-button>
              </template>
            </el-table-column>
            <el-table-column prop="product_name" label="产品" min-width="150" />
            <el-table-column prop="type" label="类型" width="80" align="center">
              <template #default="{ row }">
                <el-tag size="small">{{ getTypeText(row.type) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="amount" label="金额" width="100" align="right">
              <template #default="{ row }">¥{{ formatMoney(row.amount) }}</template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getOrderStatusType(row.status)" size="small">
                  {{ getOrderStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="下单时间" width="170" />
          </el-table>
        </el-tab-pane>

        <!-- 工单 -->
        <el-tab-pane label="工单" name="tickets">
          <el-table :data="tickets" border stripe>
            <el-table-column prop="ticket_no" label="工单号" width="130">
              <template #default="{ row }">
                <el-button type="primary" link @click="$router.push(`/ticket-detail/${row.id}`)">
                  {{ row.ticket_no }}
                </el-button>
              </template>
            </el-table-column>
            <el-table-column prop="subject" label="主题" min-width="200" />
            <el-table-column prop="priority" label="优先级" width="80" align="center">
              <template #default="{ row }">
                <el-tag :type="getPriorityType(row.priority)" size="small">
                  {{ getPriorityText(row.priority) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getTicketStatusType(row.status)" size="small">
                  {{ getTicketStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="170" />
          </el-table>
        </el-tab-pane>

        <!-- 交易记录 -->
        <el-tab-pane label="交易记录" name="transactions">
          <el-table :data="transactions" border stripe>
            <el-table-column prop="transaction_no" label="流水号" width="180" />
            <el-table-column prop="type" label="类型" width="80" align="center">
              <template #default="{ row }">
                <el-tag :type="getTransactionType(row.type)" size="small">
                  {{ getTransactionText(row.type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="amount" label="金额" width="120" align="right">
              <template #default="{ row }">
                <span :class="row.type === 'expense' ? 'text-red' : 'text-green'">
                  {{ row.type === 'expense' ? '-' : '+' }}¥{{ formatMoney(row.amount) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="balance_after" label="余额" width="120" align="right">
              <template #default="{ row }">¥{{ formatMoney(row.balance_after) }}</template>
            </el-table-column>
            <el-table-column prop="description" label="描述" min-width="200" />
            <el-table-column prop="created_at" label="时间" width="170" />
          </el-table>
        </el-tab-pane>

        <!-- 账单 -->
        <el-tab-pane label="账单" name="invoices">
          <el-table :data="invoices" border stripe>
            <el-table-column prop="invoice_no" label="账单号" width="150" />
            <el-table-column prop="description" label="描述" min-width="200" />
            <el-table-column prop="amount" label="金额" width="120" align="right">
              <template #default="{ row }">¥{{ formatMoney(row.amount) }}</template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="getInvoiceStatusType(row.status)" size="small">
                  {{ getInvoiceStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="due_date" label="到期日" width="120" />
            <el-table-column prop="paid_at" label="付款时间" width="170" />
          </el-table>
        </el-tab-pane>

        <!-- 操作日志 -->
        <el-tab-pane label="操作日志" name="logs">
          <el-timeline>
            <el-timeline-item
              v-for="log in logs"
              :key="log.id"
              :timestamp="log.created_at"
              placement="top"
            >
              <el-card shadow="never">
                <div class="log-content">
                  <span class="log-action">{{ log.action }}</span>
                  <span class="log-detail">{{ log.detail }}</span>
                </div>
                <div class="log-operator">操作人: {{ log.operator_name }}</div>
              </el-card>
            </el-timeline-item>
          </el-timeline>
        </el-tab-pane>

        <!-- 实名认证 -->
        <el-tab-pane label="实名认证" name="verification">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="认证状态">
              <el-tag :type="customer.verified ? 'success' : 'warning'" size="small">
                {{ customer.verified ? '已认证' : '未认证' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="真实姓名">{{ customer.real_name || '-' }}</el-descriptions-item>
            <el-descriptions-item label="证件类型">{{ customer.id_type || '-' }}</el-descriptions-item>
            <el-descriptions-item label="证件号码">{{ customer.id_number || '-' }}</el-descriptions-item>
            <el-descriptions-item label="认证时间">{{ customer.verified_at || '-' }}</el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>

        <!-- 自定义字段 -->
        <el-tab-pane label="自定义字段" name="custom_fields">
          <el-descriptions :column="2" border>
            <el-descriptions-item
              v-for="field in customFields"
              :key="field.id"
              :label="field.name"
            >
              {{ field.value || '-' }}
            </el-descriptions-item>
          </el-descriptions>
        </el-tab-pane>
      </el-tabs>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import request from '@/utils/http'

const route = useRoute()
const customerId = route.params.id
const activeTab = ref('overview')

// 客户信息
const customer = ref<any>({})
const products = ref([])
const orders = ref([])
const tickets = ref([])
const transactions = ref([])
const invoices = ref([])
const logs = ref([])
const customFields = ref([])

// 格式化金额
const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// 状态类型
const getStatusType = (status: string) => {
  const map: Record<string, string> = { active: 'success', disabled: 'danger', pending: 'warning' }
  return map[status] || 'info'
}

// 状态文本
const getStatusText = (status: string) => {
  const map: Record<string, string> = { active: '正常', disabled: '禁用', pending: '待验证' }
  return map[status] || '未知'
}

// 产品状态
const getProductStatusType = (status: string) => {
  const map: Record<string, string> = { active: 'success', suspended: 'danger', pending: 'warning', terminated: 'info' }
  return map[status] || 'info'
}

const getProductStatusText = (status: string) => {
  const map: Record<string, string> = { active: '正常', suspended: '暂停', pending: '待开通', terminated: '已终止' }
  return map[status] || '未知'
}

// 订单类型
const getTypeText = (type: string) => {
  const map: Record<string, string> = { new: '新购', renewal: '续费', upgrade: '升级', refund: '退款' }
  return map[type] || '未知'
}

// 订单状态
const getOrderStatusType = (status: string) => {
  const map: Record<string, string> = { pending_payment: 'warning', pending_activation: 'primary', active: 'success', completed: 'success', cancelled: 'info', refunded: 'danger' }
  return map[status] || 'info'
}

const getOrderStatusText = (status: string) => {
  const map: Record<string, string> = { pending_payment: '待付款', pending_activation: '待开通', active: '进行中', completed: '已完成', cancelled: '已取消', refunded: '已退款' }
  return map[status] || '未知'
}

// 工单优先级
const getPriorityType = (priority: number) => {
  const map: Record<number, string> = { 1: 'info', 2: 'primary', 3: 'warning', 4: 'danger' }
  return map[priority] || 'info'
}

const getPriorityText = (priority: number) => {
  const map: Record<number, string> = { 1: '低', 2: '普通', 3: '高', 4: '紧急' }
  return map[priority] || '未知'
}

// 工单状态
const getTicketStatusType = (status: string) => {
  const map: Record<string, string> = { open: 'warning', in_progress: 'primary', replied: 'success', closed: 'info' }
  return map[status] || 'info'
}

const getTicketStatusText = (status: string) => {
  const map: Record<string, string> = { open: '待处理', in_progress: '处理中', replied: '已回复', closed: '已关闭' }
  return map[status] || '未知'
}

// 交易类型
const getTransactionType = (type: string) => {
  const map: Record<string, string> = { income: 'success', expense: 'danger', refund: 'warning', recharge: 'primary' }
  return map[type] || 'info'
}

const getTransactionText = (type: string) => {
  const map: Record<string, string> = { income: '收入', expense: '支出', refund: '退款', recharge: '充值' }
  return map[type] || '未知'
}

// 账单状态
const getInvoiceStatusType = (status: string) => {
  const map: Record<string, string> = { unpaid: 'warning', paid: 'success', cancelled: 'info', refunded: 'danger' }
  return map[status] || 'info'
}

const getInvoiceStatusText = (status: string) => {
  const map: Record<string, string> = { unpaid: '待付款', paid: '已付款', cancelled: '已取消', refunded: '已退款' }
  return map[status] || '未知'
}

// 获取客户详情
const fetchCustomer = async () => {
  try {
    const data = await request.get({ url: `/api/admin/clients/${customerId}` })
    customer.value = data
  } catch (error) {
    console.error('获取客户详情失败:', error)
  }
}

// 获取产品列表
const fetchProducts = async () => {
  try {
    const data = await request.get({ url: `/api/admin/clients/${customerId}/products` })
    products.value = data || []
  } catch (error) {
    console.error('获取产品列表失败:', error)
  }
}

// 获取订单列表
const fetchOrders = async () => {
  try {
    const data = await request.get({ url: `/api/admin/clients/${customerId}/orders` })
    orders.value = data || []
  } catch (error) {
    console.error('获取订单列表失败:', error)
  }
}

// 获取工单列表
const fetchTickets = async () => {
  try {
    const data = await request.get({ url: `/api/admin/clients/${customerId}/tickets` })
    tickets.value = data || []
  } catch (error) {
    console.error('获取工单列表失败:', error)
  }
}

// 获取交易记录
const fetchTransactions = async () => {
  try {
    const data = await request.get({ url: `/api/admin/clients/${customerId}/transactions` })
    transactions.value = data || []
  } catch (error) {
    console.error('获取交易记录失败:', error)
  }
}

// 获取账单
const fetchInvoices = async () => {
  try {
    const data = await request.get({ url: `/api/admin/clients/${customerId}/invoices` })
    invoices.value = data || []
  } catch (error) {
    console.error('获取账单失败:', error)
  }
}

// 获取操作日志
const fetchLogs = async () => {
  try {
    const data = await request.get({ url: `/api/admin/clients/${customerId}/logs` })
    logs.value = data || []
  } catch (error) {
    console.error('获取操作日志失败:', error)
  }
}

// 获取自定义字段
const fetchCustomFields = async () => {
  try {
    const data = await request.get({ url: `/api/admin/clients/${customerId}/custom-fields` })
    customFields.value = data || []
  } catch (error) {
    console.error('获取自定义字段失败:', error)
  }
}

onMounted(() => {
  fetchCustomer()
  fetchProducts()
  fetchOrders()
  fetchTickets()
  fetchTransactions()
  fetchInvoices()
  fetchLogs()
  fetchCustomFields()
})
</script>

<style scoped lang="scss">
.customer-detail-page {
  padding: 16px;
}

.customer-info-card {
  margin-bottom: 16px;
}

.customer-basic {
  display: flex;
  align-items: center;
  gap: 16px;
}

.customer-name {
  h3 {
    margin: 0 0 8px 0;
    font-size: 18px;
  }
}

.stat-item {
  text-align: center;
  padding: 10px 0;
}

.stat-label {
  font-size: 14px;
  color: #86909C;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #1D2129;
}

.detail-card {
  :deep(.el-tabs__content) {
    padding: 20px 0;
  }
}

.log-content {
  display: flex;
  gap: 8px;
}

.log-action {
  font-weight: 500;
}

.log-detail {
  color: #86909C;
}

.log-operator {
  margin-top: 8px;
  font-size: 12px;
  color: #86909C;
}

.text-green { color: #36D391; }
.text-red { color: #EF4444; }
</style>
