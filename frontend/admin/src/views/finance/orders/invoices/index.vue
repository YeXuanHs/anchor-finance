<template>
  <div class="invoice-list-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>账单管理</span>
          <el-button type="primary" @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="发票号">
          <el-input v-model="searchForm.invoice_no" placeholder="请输入发票号" clearable />
        </el-form-item>
        <el-form-item label="客户">
          <el-input v-model="searchForm.client_name" placeholder="请输入客户名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待审核" :value="0" />
            <el-option label="审核通过" :value="1" />
            <el-option label="已开票" :value="2" />
            <el-option label="已驳回" :value="3" />
            <el-option label="已作废" :value="4" />
          </el-select>
        </el-form-item>
        <el-form-item label="发票类型">
          <el-select v-model="searchForm.invoice_type" placeholder="全部" clearable>
            <el-option label="增值税普通发票" :value="1" />
            <el-option label="增值税专用发票" :value="2" />
            <el-option label="电子发票" :value="3" />
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
        <el-table-column prop="invoice_no" label="发票号" width="180">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">
              {{ row.invoice_no }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="order_no" label="关联订单" width="170" />
        <el-table-column prop="client_name" label="客户" width="120" />
        <el-table-column prop="invoice_type" label="发票类型" width="140" align="center">
          <template #default="{ row }">
            {{ getInvoiceTypeText(row.invoice_type) }}
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="110" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="tax_amount" label="税额" width="100" align="right">
          <template #default="{ row }">
            ¥{{ formatAmount(row.tax_amount) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getInvoiceStatusType(row.status)" size="small">
              {{ getInvoiceStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="170" />
        <el-table-column label="操作" width="200" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewDetail(row)">详情</el-button>
            <el-button
              v-if="row.status === 0"
              type="success"
              link
              @click="handleAudit(row, 1)"
            >
              通过
            </el-button>
            <el-button
              v-if="row.status === 0"
              type="danger"
              link
              @click="handleAudit(row, 3)"
            >
              驳回
            </el-button>
            <el-button
              v-if="row.status === 1"
              type="primary"
              link
              @click="handleIssue(row)"
            >
              开票
            </el-button>
            <el-popconfirm
              v-if="row.status === 2"
              title="确定作废该发票吗？"
              @confirm="handleVoid(row)"
            >
              <template #reference>
                <el-button type="danger" link>作废</el-button>
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

    <!-- 发票详情对话框 -->
    <el-dialog v-model="detailVisible" title="发票详情" width="800px" destroy-on-close>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="发票号">{{ detailData.invoice_no }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getInvoiceStatusType(detailData.status)" size="small">
            {{ getInvoiceStatusText(detailData.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="关联订单">{{ detailData.order_no }}</el-descriptions-item>
        <el-descriptions-item label="客户名称">{{ detailData.client_name }}</el-descriptions-item>
        <el-descriptions-item label="发票类型">{{ getInvoiceTypeText(detailData.invoice_type) }}</el-descriptions-item>
        <el-descriptions-item label="纳税人识别号">{{ detailData.tax_no || '-' }}</el-descriptions-item>
        <el-descriptions-item label="开票金额">
          <span class="amount-text">¥{{ formatAmount(detailData.amount) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="税额">¥{{ formatAmount(detailData.tax_amount) }}</el-descriptions-item>
        <el-descriptions-item label="价税合计">
          <span class="amount-text total-amount">¥{{ formatAmount((detailData.amount || 0) + (detailData.tax_amount || 0)) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="开票日期">{{ detailData.issue_date || '-' }}</el-descriptions-item>
        <el-descriptions-item label="申请时间" :span="2">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item label="公司名称" :span="2">{{ detailData.company_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="注册地址">{{ detailData.company_address || '-' }}</el-descriptions-item>
        <el-descriptions-item label="注册电话">{{ detailData.company_phone || '-' }}</el-descriptions-item>
        <el-descriptions-item label="开户银行">{{ detailData.bank_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="银行账号">{{ detailData.bank_account || '-' }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ detailData.remark || '无' }}</el-descriptions-item>
      </el-descriptions>

      <!-- 审核区域 -->
      <div v-if="detailData.status === 0" class="audit-section">
        <el-divider content-position="left">审核操作</el-divider>
        <el-input
          v-model="auditRemark"
          type="textarea"
          :rows="3"
          placeholder="请输入审核备注（可选）"
        />
        <div class="audit-actions">
          <el-button type="success" @click="handleAuditFromDetail(1)" :loading="auditLoading">
            审核通过
          </el-button>
          <el-button type="danger" @click="handleAuditFromDetail(3)" :loading="auditLoading">
            审核驳回
          </el-button>
        </div>
      </div>

      <!-- 开票确认区域 -->
      <div v-if="detailData.status === 1" class="audit-section">
        <el-divider content-position="left">开票操作</el-divider>
        <el-form label-width="100px">
          <el-form-item label="发票号码">
            <el-input v-model="issueForm.invoice_number" placeholder="请输入实际发票号码" />
          </el-form-item>
          <el-form-item label="开票日期">
            <el-date-picker
              v-model="issueForm.issue_date"
              type="date"
              placeholder="选择开票日期"
              value-format="YYYY-MM-DD"
            />
          </el-form-item>
        </el-form>
        <div class="audit-actions">
          <el-button type="primary" @click="handleIssueFromDetail" :loading="auditLoading">
            确认开票
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
import { exportToCSV } from '@/utils/export'

// 加载状态
const loading = ref(false)
const auditLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  invoice_no: '',
  client_name: '',
  status: undefined as number | undefined,
  invoice_type: undefined as number | undefined,
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

// 开票表单
const issueForm = reactive({
  invoice_number: '',
  issue_date: ''
})

// 发票类型映射
const INVOICE_TYPE_MAP: Record<number, string> = {
  1: '增值税普通发票',
  2: '增值税专用发票',
  3: '电子发票'
}

// 发票状态映射
const INVOICE_STATUS_MAP: Record<number, { text: string; type: string }> = {
  0: { text: '待审核', type: 'warning' },
  1: { text: '审核通过', type: 'primary' },
  2: { text: '已开票', type: 'success' },
  3: { text: '已驳回', type: 'danger' },
  4: { text: '已作废', type: 'info' }
}

// 获取发票类型文本
const getInvoiceTypeText = (type: number) => {
  return INVOICE_TYPE_MAP[type] || '未知'
}

// 获取发票状态文本
const getInvoiceStatusText = (status: number) => {
  return INVOICE_STATUS_MAP[status]?.text || '未知'
}

// 获取发票状态类型
const getInvoiceStatusType = (status: number) => {
  return (INVOICE_STATUS_MAP[status]?.type || 'info') as any
}

// 格式化金额
const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

// 获取发票列表
const fetchInvoices = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      invoice_no: searchForm.invoice_no || undefined,
      client_name: searchForm.client_name || undefined,
      status: searchForm.status,
      invoice_type: searchForm.invoice_type
    }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({
      url: '/api/admin/invoices',
      params
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取发票列表失败:', error)
    ElMessage.error('获取发票列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchInvoices()
}

// 重置
const handleReset = () => {
  searchForm.invoice_no = ''
  searchForm.client_name = ''
  searchForm.status = undefined
  searchForm.invoice_type = undefined
  searchForm.date_range = []
  handleSearch()
}

// 查看详情
const handleViewDetail = (row: any) => {
  detailData.value = { ...row }
  auditRemark.value = ''
  issueForm.invoice_number = ''
  issueForm.issue_date = ''
  detailVisible.value = true
}

// 快速审核
const handleAudit = async (row: any, targetStatus: number) => {
  const action = targetStatus === 1 ? '通过' : '驳回'
  try {
    await ElMessageBox.confirm(`确定要${action}该发票申请吗？`, '审核确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.put({
      url: `/api/admin/invoices/${row.id}/status`,
      params: { status: targetStatus }
    })
    ElMessage.success(`发票${action}成功`)
    fetchInvoices()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(`发票${action}失败`)
    }
  }
}

// 从详情页审核
const handleAuditFromDetail = async (targetStatus: number) => {
  const action = targetStatus === 1 ? '通过' : '驳回'
  auditLoading.value = true
  try {
    await request.put({
      url: `/api/admin/invoices/${detailData.value.id}/status`,
      params: {
        status: targetStatus,
        remark: auditRemark.value || undefined
      }
    })
    ElMessage.success(`发票${action}成功`)
    detailVisible.value = false
    fetchInvoices()
  } catch (error) {
    ElMessage.error(`发票${action}失败`)
  } finally {
    auditLoading.value = false
  }
}

// 快速开票
const handleIssue = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要对该发票进行开票操作吗？', '开票确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.put({
      url: `/api/admin/invoices/${row.id}/status`,
      params: { status: 2 }
    })
    ElMessage.success('开票成功')
    fetchInvoices()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error('开票失败')
    }
  }
}

// 从详情页开票
const handleIssueFromDetail = async () => {
  if (!issueForm.invoice_number) {
    ElMessage.warning('请输入发票号码')
    return
  }
  auditLoading.value = true
  try {
    await request.put({
      url: `/api/admin/invoices/${detailData.value.id}/status`,
      params: {
        status: 2,
        invoice_number: issueForm.invoice_number,
        issue_date: issueForm.issue_date || undefined
      }
    })
    ElMessage.success('开票成功')
    detailVisible.value = false
    fetchInvoices()
  } catch (error) {
    ElMessage.error('开票失败')
  } finally {
    auditLoading.value = false
  }
}

// 作废发票
const handleVoid = async (row: any) => {
  try {
    await request.put({
      url: `/api/admin/invoices/${row.id}/status`,
      params: { status: 4 }
    })
    ElMessage.success('发票已作废')
    fetchInvoices()
  } catch (error) {
    ElMessage.error('作废失败')
  }
}

// 导出
const handleExport = async () => {
  try {
    const data = await request.get({ url: '/api/admin/invoices', params: { page: 1, page_size: 9999 } })
    const list = data.list || data || []
    exportToCSV(list, [
      { key: 'id', title: 'ID' },
      { key: 'invoice_no', title: '发票号' },
      { key: 'user_id', title: '用户ID' },
      { key: 'total', title: '金额' },
      { key: 'status', title: '状态' },
      { key: 'created_at', title: '创建时间' }
    ], '发票列表')
    ElMessage.success('导出成功')
  } catch { ElMessage.error('导出失败') }
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchInvoices()
}

// 页码变化
const handlePageChange = () => {
  fetchInvoices()
}

onMounted(() => {
  fetchInvoices()
})
</script>

<style scoped lang="scss">
.invoice-list-page {
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

.total-amount {
  font-size: 15px;
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
