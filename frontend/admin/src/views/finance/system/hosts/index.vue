<template>
  <div class="hosts-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>主机管理</span>
          <el-button type="primary" @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="域名/IP/用户名" clearable />
        </el-form-item>
        <el-form-item label="用户ID">
          <el-input v-model="searchForm.user_id" placeholder="用户ID" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="活跃" value="active" />
            <el-option label="暂停" value="suspended" />
            <el-option label="待开通" value="pending" />
            <el-option label="已终止" value="terminated" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="user_id" label="用户ID" width="90" align="center" />
        <el-table-column prop="username" label="用户名" width="120" show-overflow-tooltip />
        <el-table-column prop="product_name" label="产品名称" min-width="150" show-overflow-tooltip />
        <el-table-column prop="domain" label="域名/IP" min-width="180" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="开通时间" width="170" />
        <el-table-column prop="due_date" label="到期时间" width="170">
          <template #default="{ row }">
            <span :class="{ 'danger-text': isExpiringSoon(row.due_date) }">
              {{ row.due_date || '-' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="220" fixed="right" align="center">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleView(row)">详情</el-button>
            <el-popconfirm
              v-if="row.status === 'active'"
              title="确定暂停该主机吗？"
              @confirm="handleSuspend(row)"
            >
              <template #reference>
                <el-button type="warning" link>暂停</el-button>
              </template>
            </el-popconfirm>
            <el-button
              v-if="row.status === 'suspended'"
              type="success"
              link
              @click="handleActivate(row)"
            >
              激活
            </el-button>
            <el-popconfirm
              v-if="row.status !== 'terminated'"
              title="确定终止该主机吗？此操作不可恢复！"
              @confirm="handleTerminate(row)"
            >
              <template #reference>
                <el-button type="danger" link>终止</el-button>
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

    <!-- 详情对话框 -->
    <el-dialog v-model="detailVisible" title="主机详情" width="700px" destroy-on-close>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="主机ID">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item label="用户ID">{{ detailData.user_id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ detailData.username || '-' }}</el-descriptions-item>
        <el-descriptions-item label="产品名称">{{ detailData.product_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="域名/IP" :span="2">{{ detailData.domain || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(detailData.status)" size="small">
            {{ getStatusText(detailData.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ detailData.ip || '-' }}</el-descriptions-item>
        <el-descriptions-item label="开通时间">{{ detailData.created_at || '-' }}</el-descriptions-item>
        <el-descriptions-item label="到期时间">{{ detailData.due_date || '-' }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ detailData.remark || '无' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const detailVisible = ref(false)
const detailData = ref<any>({})

const searchForm = reactive({
  keyword: '',
  user_id: '',
  status: '' as string
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<any[]>([])

const STATUS_MAP: Record<string, { text: string; type: string }> = {
  active: { text: '活跃', type: 'success' },
  suspended: { text: '暂停', type: 'warning' },
  pending: { text: '待开通', type: 'info' },
  terminated: { text: '已终止', type: 'danger' }
}

const getStatusText = (status: string) => STATUS_MAP[status]?.text || status

const getStatusType = (status: string) => (STATUS_MAP[status]?.type || 'info') as any

const isExpiringSoon = (dueDate: string) => {
  if (!dueDate) return false
  const diff = new Date(dueDate).getTime() - Date.now()
  return diff > 0 && diff < 7 * 24 * 60 * 60 * 1000
}

const fetchList = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/hosts',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        keyword: searchForm.keyword || undefined,
        user_id: searchForm.user_id || undefined,
        status: searchForm.status || undefined
      }
    })
    tableData.value = data.list || data || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取主机列表失败:', error)
    ElMessage.error('获取主机列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchList()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.user_id = ''
  searchForm.status = ''
  handleSearch()
}

const handleView = (row: any) => {
  detailData.value = { ...row }
  detailVisible.value = true
}

const handleSuspend = async (row: any) => {
  try {
    await request.put({ url: `/api/admin/hosts/${row.id}/status`, params: { status: 'suspended' } })
    ElMessage.success('主机已暂停')
    fetchList()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleActivate = async (row: any) => {
  try {
    await request.put({ url: `/api/admin/hosts/${row.id}/status`, params: { status: 'active' } })
    ElMessage.success('主机已激活')
    fetchList()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleTerminate = async (row: any) => {
  try {
    await request.put({ url: `/api/admin/hosts/${row.id}/status`, params: { status: 'terminated' } })
    ElMessage.success('主机已终止')
    fetchList()
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const handleExport = () => {
  ElMessage.info('导出功能开发中...')
}

const handleSizeChange = () => {
  pagination.page = 1
  fetchList()
}

const handlePageChange = () => {
  fetchList()
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.hosts-page {
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

.danger-text {
  color: #f56c6c;
  font-weight: 600;
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}
</style>
