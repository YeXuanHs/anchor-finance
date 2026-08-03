<template>
  <div class="order-list-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>订单列表</span>
          <el-button type="primary" @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="订单号">
          <el-input v-model="searchForm.order_no" placeholder="请输入订单号" clearable />
        </el-form-item>
        <el-form-item label="客户">
          <el-input v-model="searchForm.client_name" placeholder="请输入客户名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待付款" :value="0" />
            <el-option label="待审核" :value="1" />
            <el-option label="审核通过" :value="2" />
            <el-option label="已开通" :value="3" />
            <el-option label="已完成" :value="4" />
            <el-option label="已取消" :value="5" />
            <el-option label="已退款" :value="6" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="70" align="center" />
        <el-table-column prop="order_no" label="订单号" width="170">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">
              {{ row.order_no }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="client_name" label="客户" width="120" />
        <el-table-column prop="product_name" label="产品/服务" min-width="160" show-overflow-tooltip />
        <el-table-column prop="amount" label="金额" width="110" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="pay_method" label="支付方式" width="100" align="center">
          <template #default="{ row }">
            {{ row.pay_method || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getOrderStatusType(row.status)" size="small">
              {{ getOrderStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="下单时间" width="170" />
        <el-table-column label="操作" width="180" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">详情</el-button>
            <el-button
              v-if="row.status === 1"
              type="success"
              link
              @click="handleAudit(row, 2)"
            >
              通过
            </el-button>
            <el-button
              v-if="row.status === 1"
              type="danger"
              link
              @click="handleAudit(row, 5)"
            >
              驳回
            </el-button>
            <el-popconfirm
              v-if="row.status === 0"
              title="确定取消该订单吗？"
              @confirm="handleAudit(row, 5)"
            >
              <template #reference>
                <el-button type="danger" link>取消</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.page_size"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange"
          @current-change="handlePageChange"
        />
      </div>
    </el-card>

    <!-- 订单详情对话框 -->
    <el-dialog v-model="detailVisible" title="订单详情" width="750px" destroy-on-close>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="订单号">{{ detailData.order_no }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getOrderStatusType(detailData.status)" size="small">
            {{ getOrderStatusText(detailData.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="客户名称">{{ detailData.client_name }}</el-descriptions-item>
        <el-descriptions-item label="客户邮箱">{{ detailData.client_email || '-' }}</el-descriptions-item>
        <el-descriptions-item label="产品/服务" :span="2">{{ detailData.product_name }}</el-descriptions-item>
        <el-descriptions-item label="金额">
          <span class="amount-text">¥{{ formatAmount(detailData.amount) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="支付方式">{{ detailData.pay_method || '-' }}</el-descriptions-item>
        <el-descriptions-item label="下单时间">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item label="更新时间">{{ detailData.updated_at || '-' }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ detailData.remark || '无' }}</el-descriptions-item>
      </el-descriptions>

      <!-- 审核区域 -->
      <div v-if="detailData.status === 1" class="audit-section">
        <el-divider content-position="left">审核操作</el-divider>
        <el-input
          v-model="auditRemark"
          type="textarea"
          :rows="3"
          placeholder="请输入审核备注（可选）"
        />
        <div class="audit-actions">
          <el-button type="success" @click="handleAuditFromDetail(2)" :loading="auditLoading">
            审核通过
          </el-button>
          <el-button type="danger" @click="handleAuditFromDetail(5)" :loading="auditLoading">
            审核驳回
          </el-button>
        </div>
      </div>

      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Download } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

// 加载状态
const loading = ref(false)
const auditLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  order_no: '',
  client_name: '',
  status: undefined as number | undefined,
  date_range: [] as string[]
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 表格数据
const tableData = ref<any[]>([])

// 详情对话框
const detailVisible = ref(false)
const detailData = ref<any>({})
const auditRemark = ref('')

// 订单状态映射
const ORDER_STATUS_MAP: Record<number, { text: string; type: string }> = {
  0: { text: '待付款', type: 'warning' },
  1: { text: '待审核', type: 'primary' },
  2: { text: '审核通过', type: 'success' },
  3: { text: '已开通', type: 'success' },
  4: { text: '已完成', type: 'info' },
  5: { text: '已取消', type: 'info' },
  6: { text: '已退款', type: 'danger' }
}

// 获取订单状态文本
const getOrderStatusText = (status: number) => {
  return ORDER_STATUS_MAP[status]?.text || '未知'
}

// 获取订单状态类型
const getOrderStatusType = (status: number) => {
  return (ORDER_STATUS_MAP[status]?.type || 'info') as any
}

// 格式化金额
const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

// 获取订单列表
const fetchOrders = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      order_no: searchForm.order_no || undefined,
      client_name: searchForm.client_name || undefined,
      status: searchForm.status
    }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({
      url: '/api/admin/orders',
      params
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取订单列表失败:', error)
    ElMessage.error('获取订单列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchOrders()
}

// 重置
const handleReset = () => {
  searchForm.order_no = ''
  searchForm.client_name = ''
  searchForm.status = undefined
  searchForm.date_range = []
  handleSearch()
}

// 查看详情
const handleViewDetail = (row: any) => {
  detailData.value = { ...row }
  auditRemark.value = ''
  detailVisible.value = true
}

// 快速审核
const handleAudit = async (row: any, targetStatus: number) => {
  const action = targetStatus === 2 ? '通过' : '驳回'
  try {
    await ElMessageBox.confirm(`确定要${action}该订单吗？`, '审核确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.post({
      url: `/api/admin/orders/${row.id}/status`,
      params: { status: targetStatus }
    })
    ElMessage.success(`订单${action}成功`)
    fetchOrders()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(`订单${action}失败`)
    }
  }
}

// 从详情页审核
const handleAuditFromDetail = async (targetStatus: number) => {
  const action = targetStatus === 2 ? '通过' : '驳回'
  auditLoading.value = true
  try {
    await request.post({
      url: `/api/admin/orders/${detailData.value.id}/status`,
      params: {
        status: targetStatus,
        remark: auditRemark.value || undefined
      }
    })
    ElMessage.success(`订单${action}成功`)
    detailVisible.value = false
    fetchOrders()
  } catch (error) {
    ElMessage.error(`订单${action}失败`)
  } finally {
    auditLoading.value = false
  }
}

// 导出
const handleExport = () => {
  ElMessage.info('导出功能开发中...')
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchOrders()
}

// 页码变化
const handlePageChange = () => {
  fetchOrders()
}

onMounted(() => {
  fetchOrders()
})
</script>

<style scoped lang="scss">
.order-list-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  margin-bottom: 20px;
}

.amount-text {
  font-weight: 600;
  color: #f56c6c;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.audit-section {
  margin-top: 16px;

  .audit-actions {
    margin-top: 12px;
    display: flex;
    justify-content: flex-end;
    gap: 12px;
  }
}
</style>
