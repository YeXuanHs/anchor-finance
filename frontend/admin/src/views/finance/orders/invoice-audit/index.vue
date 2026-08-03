<template>
  <div class="invoice-audit-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>发票审核</span>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="发票号">
          <el-input v-model="searchForm.invoice_no" placeholder="请输入发票号" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待审核" :value="0" />
            <el-option label="审核通过" :value="1" />
            <el-option label="审核驳回" :value="2" />
            <el-option label="已开票" :value="3" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="invoice_no" label="发票号" width="180" />
        <el-table-column prop="client_name" label="客户" width="120" />
        <el-table-column prop="title" label="发票抬头" min-width="200" />
        <el-table-column prop="tax_no" label="税号" width="160" />
        <el-table-column prop="amount" label="金额" width="120" align="right">
          <template #default="{ row }">
            <span class="amount-text">¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="statusTypeMap[row.status]" size="small">
              {{ statusLabelMap[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">详情</el-button>
            <el-button v-if="row.status === 0" type="success" link @click="handleAudit(row, 1)">通过</el-button>
            <el-button v-if="row.status === 0" type="danger" link @click="handleAudit(row, 2)">驳回</el-button>
            <el-button v-if="row.status === 1" type="warning" link @click="handleMarkIssued(row)">标记已开票</el-button>
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

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="发票详情" width="650px">
      <el-descriptions :column="2" border v-if="detailData">
        <el-descriptions-item label="发票号">{{ detailData.invoice_no }}</el-descriptions-item>
        <el-descriptions-item label="客户">{{ detailData.client_name }}</el-descriptions-item>
        <el-descriptions-item label="发票抬头" :span="2">{{ detailData.title }}</el-descriptions-item>
        <el-descriptions-item label="税号">{{ detailData.tax_no || '-' }}</el-descriptions-item>
        <el-descriptions-item label="金额">
          <span class="amount-text">¥{{ formatAmount(detailData.amount) }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="开户银行">{{ detailData.bank_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="银行账号">{{ detailData.bank_account || '-' }}</el-descriptions-item>
        <el-descriptions-item label="地址" :span="2">{{ detailData.address || '-' }}</el-descriptions-item>
        <el-descriptions-item label="电话">{{ detailData.phone || '-' }}</el-descriptions-item>
        <el-descriptions-item label="邮箱">{{ detailData.email || '-' }}</el-descriptions-item>
        <el-descriptions-item label="申请时间" :span="2">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ detailData.remark || '无' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<!-- TODO: 完善发票审核页面功能
  1. 添加审核备注输入
  2. 添加批量审核功能
  3. 添加发票信息导出
  4. 添加关联订单查看
-->
<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const statusTypeMap: Record<number, string> = {
  0: 'warning',
  1: 'success',
  2: 'danger',
  3: 'info'
}

const statusLabelMap: Record<number, string> = {
  0: '待审核',
  1: '审核通过',
  2: '审核驳回',
  3: '已开票'
}

const loading = ref(false)
const detailVisible = ref(false)
const detailData = ref<any>(null)

const searchForm = reactive({
  invoice_no: '',
  status: undefined as number | undefined
})

const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

const fetchData = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/invoices/audit',
      params: { page: pagination.page, page_size: pagination.page_size, ...searchForm }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取发票列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { searchForm.invoice_no = ''; searchForm.status = undefined; handleSearch() }
const handleView = (row: any) => { detailData.value = row; detailVisible.value = true }

const handleAudit = async (row: any, status: number) => {
  const action = status === 1 ? '通过' : '驳回'
  try {
    await ElMessageBox.confirm(`确定${action}该发票申请吗？`, '审核确认')
    await request.post({ url: `/api/admin/invoices/${row.id}/audit`, params: { status } })
    ElMessage.success(`审核${action}成功`)
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel') ElMessage.error('操作失败')
  }
}

const handleMarkIssued = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/invoices/${row.id}/issued` })
    ElMessage.success('已标记为已开票')
    fetchData()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.invoice-audit-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.amount-text { font-weight: 600; color: var(--el-color-primary); }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
