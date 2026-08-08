<template>
  <div class="customer-list-page">
    <!-- 搜索筛选区域 -->
    <el-card shadow="never" class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="客户名/邮箱/手机号/ID"
            clearable
            style="width: 240px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable style="width: 120px">
            <el-option label="正常" value="active" />
            <el-option label="禁用" value="disabled" />
            <el-option label="待验证" value="pending" />
          </el-select>
        </el-form-item>
        <el-form-item label="客户组">
          <el-select v-model="searchForm.group_id" placeholder="全部" clearable style="width: 120px">
            <el-option v-for="group in clientGroups" :key="group.id" :label="group.name" :value="group.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="注册时间">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 240px"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 操作栏 -->
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <div class="action-left">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加客户
          </el-button>
          <el-button type="danger" :disabled="selectedIds.length === 0" @click="handleBatchDelete">
            <el-icon><Delete /></el-icon>
            批量删除
          </el-button>
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </div>
        <div class="action-right">
          <el-button circle @click="fetchList">
            <el-icon><Refresh /></el-icon>
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        @selection-change="handleSelectionChange"
        @sort-change="handleSortChange"
      >
        <el-table-column type="selection" width="50" align="center" />
        <el-table-column prop="id" label="ID" width="80" sortable="custom" align="center" />
        <el-table-column prop="username" label="客户名" min-width="120">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">
              {{ row.username }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" min-width="180" />
        <el-table-column prop="phone" label="手机号" width="130" />
        <el-table-column prop="company" label="公司" min-width="120" />
        <el-table-column prop="group_name" label="客户组" width="100" />
        <el-table-column prop="balance" label="余额" width="100" align="right">
          <template #default="{ row }">
            ¥{{ formatMoney(row.balance) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="注册时间" width="170" sortable="custom" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleView(row)">
              查看
            </el-button>
            <el-button type="primary" link size="small" @click="handleEdit(row)">
              编辑
            </el-button>
            <el-button type="primary" link size="small" @click="handleLoginAs(row)">
              登录
            </el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">
              删除
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Search, Refresh, Plus, Delete, Download } from '@element-plus/icons-vue'
import request from '@/utils/http'

const router = useRouter()
const loading = ref(false)
const tableData = ref([])
const selectedIds = ref<number[]>([])
const clientGroups = ref<{ id: number; name: string }[]>([])

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: '',
  group_id: '',
  date_range: null as [Date, Date] | null
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 排序
const sortParams = reactive({
  sort: 'created_at',
  order: 'desc'
})

// 格式化金额
const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// 状态类型
const getStatusType = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    active: 'success',
    disabled: 'danger',
    pending: 'warning'
  }
  return map[status] || 'info'
}

// 状态文本
const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: '正常',
    disabled: '禁用',
    pending: '待验证'
  }
  return map[status] || '未知'
}

// 获取列表数据
const fetchList = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      sort: sortParams.sort,
      order: sortParams.order
    }

    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.status) params.status = searchForm.status
    if (searchForm.group_id) params.group_id = searchForm.group_id
    if (searchForm.date_range) {
      params.start_date = searchForm.date_range[0].toISOString().split('T')[0]
      params.end_date = searchForm.date_range[1].toISOString().split('T')[0]
    }

    const data = await request.get({ url: '/api/admin/clients', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取客户列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 获取客户组
const fetchGroups = async () => {
  try {
    const data = await request.get({ url: '/api/admin/client-groups' })
    clientGroups.value = data || []
  } catch (error) {
    console.error('获取客户组失败:', error)
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchList()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = ''
  searchForm.group_id = ''
  searchForm.date_range = null
  pagination.page = 1
  fetchList()
}

// 选择变化
const handleSelectionChange = (selection: any[]) => {
  selectedIds.value = selection.map((item) => item.id)
}

// 排序变化
const handleSortChange = ({ prop, order }: any) => {
  sortParams.sort = prop || 'created_at'
  sortParams.order = order === 'ascending' ? 'asc' : 'desc'
  fetchList()
}

// 分页大小变化
const handleSizeChange = (size: number) => {
  pagination.page_size = size
  pagination.page = 1
  fetchList()
}

// 页码变化
const handlePageChange = (page: number) => {
  pagination.page = page
  fetchList()
}

// 添加客户
const handleAdd = () => {
  router.push('/customer-add')
}

// 查看客户
const handleView = (row: any) => {
  router.push(`/customer-view/${row.id}`)
}

// 编辑客户
const handleEdit = (row: any) => {
  router.push(`/customer-view/${row.id}`)
}

// 登录为客户
const handleLoginAs = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要以客户 "${row.username}" 的身份登录吗？`, '确认登录', {
      type: 'warning'
    })
    const data = await request.post({ url: `/api/admin/clients/${row.id}/login-as` })
    if (data.token) {
      window.open(`/?token=${data.token}`, '_blank')
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('登录失败:', error)
    }
  }
}

// 删除客户
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除客户 "${row.username}" 吗？此操作不可恢复。`, '确认删除', {
      type: 'warning'
    })
    await request.del({ url: `/api/admin/clients/${row.id}` })
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('删除失败:', error)
    }
  }
}

// 批量删除
const handleBatchDelete = async () => {
  try {
    await ElMessageBox.confirm(`确定要删除选中的 ${selectedIds.value.length} 个客户吗？此操作不可恢复。`, '确认批量删除', {
      type: 'warning'
    })
    await request.post({ url: '/api/admin/clients/batch-delete', data: { ids: selectedIds.value } })
    ElMessage.success('批量删除成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('批量删除失败:', error)
    }
  }
}

// 导出
const handleExport = async () => {
  try {
    const params: any = {}
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.status) params.status = searchForm.status
    if (searchForm.group_id) params.group_id = searchForm.group_id

    const response = await request.get({
      url: '/api/admin/clients/export',
      params,
      responseType: 'blob'
    })

    const url = window.URL.createObjectURL(new Blob([response]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `客户列表_${new Date().toISOString().split('T')[0]}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    console.error('导出失败:', error)
  }
}

onMounted(() => {
  fetchList()
  fetchGroups()
})
</script>

<style scoped lang="scss">
.customer-list-page {
  padding: 16px;
}

.search-card {
  margin-bottom: 16px;

  :deep(.el-card__body) {
    padding-bottom: 0;
  }
}

.action-card {
  margin-bottom: 16px;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.action-left {
  display: flex;
  gap: 8px;
}

.table-card {
  :deep(.el-card__body) {
    padding: 0;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px;
}
</style>
