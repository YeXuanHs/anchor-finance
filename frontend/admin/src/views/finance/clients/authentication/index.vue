<template>
  <div class="client-auth-page">
    <!-- 统计卡片 -->
    <el-row :gutter="16" class="stats-row">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value pending">{{ stats.pending }}</div>
          <div class="stat-label">待审核</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value approved">{{ stats.approved }}</div>
          <div class="stat-label">已通过</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value rejected">{{ stats.rejected }}</div>
          <div class="stat-label">已拒绝</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value total">{{ stats.total }}</div>
          <div class="stat-label">总计</div>
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>实名认证审核</span>
          <div class="header-actions">
            <el-button type="primary" size="small" @click="handleBatchApprove" :disabled="!selectedRows.length">批量通过</el-button>
          </div>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待审核" :value="0" />
            <el-option label="已通过" :value="1" />
            <el-option label="已拒绝" :value="2" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="姓名/身份证号" clearable />
        </el-form-item>
        <el-form-item label="认证方式">
          <el-select v-model="searchForm.certify_type" placeholder="全部" clearable>
            <el-option label="身份证" value="idcard" />
            <el-option label="人脸识别" value="face" />
            <el-option label="企业认证" value="enterprise" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="authList" v-loading="loading" stripe border @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="realname" label="认证名称" width="120" />
        <el-table-column prop="idcard" label="身份证号码" width="180" show-overflow-tooltip />
        <el-table-column prop="certify_type" label="认证方式" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ getCertifyTypeText(row.certify_type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="认证类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'company' ? 'warning' : 'info'" size="small">{{ row.type === 'company' ? '企业' : '个人' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusText(row.status) }}</el-tag>
            <div v-if="row.reason" class="reject-reason">{{ row.reason }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="提交时间" width="180" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <template v-if="row.status === 0">
              <el-button size="small" type="success" @click="handleApprove(row)">通过</el-button>
              <el-button size="small" type="danger" @click="showRejectDialog(row)">拒绝</el-button>
            </template>
            <el-button size="small" @click="handleViewDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.limit"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 16px; justify-content: flex-end"
        @size-change="fetchAuthList"
        @current-change="fetchAuthList"
      />
    </el-card>

    <!-- 拒绝对话框 -->
    <el-dialog v-model="rejectDialogVisible" title="拒绝原因" width="450px">
      <el-form label-width="80px">
        <el-form-item label="拒绝原因">
          <el-input v-model="rejectReason" type="textarea" :rows="3" placeholder="请输入拒绝原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectDialogVisible = false">取消</el-button>
        <el-button type="danger" @click="handleReject" :loading="actionLoading">确定拒绝</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="detailDialogVisible" title="认证详情" width="600px">
      <el-descriptions :column="2" border v-if="currentAuth">
        <el-descriptions-item label="用户ID">{{ currentAuth.id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ currentAuth.username }}</el-descriptions-item>
        <el-descriptions-item label="认证名称">{{ currentAuth.realname }}</el-descriptions-item>
        <el-descriptions-item label="身份证号">{{ currentAuth.idcard }}</el-descriptions-item>
        <el-descriptions-item label="认证方式">{{ getCertifyTypeText(currentAuth.certify_type) }}</el-descriptions-item>
        <el-descriptions-item label="认证类型">{{ currentAuth.type === 'company' ? '企业' : '个人' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentAuth.status)">{{ getStatusText(currentAuth.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="提交时间">{{ currentAuth.create_time }}</el-descriptions-item>
        <el-descriptions-item label="审核时间" :span="2">{{ currentAuth.check_time || '-' }}</el-descriptions-item>
        <el-descriptions-item v-if="currentAuth.reason" label="拒绝原因" :span="2">
          <span class="reject-reason">{{ currentAuth.reason }}</span>
        </el-descriptions-item>
        <el-descriptions-item v-if="currentAuth.front_image" label="正面照片" :span="2">
          <el-image :src="currentAuth.front_image" :preview-src-list="[currentAuth.front_image]" style="max-width: 300px" />
        </el-descriptions-item>
        <el-descriptions-item v-if="currentAuth.back_image" label="反面照片" :span="2">
          <el-image :src="currentAuth.back_image" :preview-src-list="[currentAuth.back_image]" style="max-width: 300px" />
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const actionLoading = ref(false)
const rejectDialogVisible = ref(false)
const detailDialogVisible = ref(false)
const rejectReason = ref('')
const currentAuth = ref<any>(null)
const selectedRows = ref<any[]>([])

const stats = reactive({ pending: 0, approved: 0, rejected: 0, total: 0 })
const searchForm = reactive({ status: undefined as number | undefined, keyword: '', certify_type: '' })
const pagination = reactive({ page: 1, limit: 10, total: 0 })
const authList = ref<any[]>([])

const getStatusType = (status: number) => {
  const map: Record<number, any> = { 0: 'warning', 1: 'success', 2: 'danger' }
  return map[status] || 'info'
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = { 0: '待审核', 1: '已通过', 2: '已拒绝' }
  return map[status] || '未知'
}

const getCertifyTypeText = (type: string) => {
  const map: Record<string, string> = { idcard: '身份证', face: '人脸识别', enterprise: '企业认证' }
  return map[type] || type || '-'
}

const fetchStats = async () => {
  try {
    const data = await request.get({ url: '/api/admin/certifications/stats' })
    if (data) Object.assign(stats, data)
  } catch { /* ignore */ }
}

const fetchAuthList = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, limit: pagination.limit }
    if (searchForm.status !== undefined) params.status = searchForm.status
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.certify_type) params.certify_type = searchForm.certify_type
    const res = await request.get({ url: '/api/admin/certifications', params })
    if (res) {
      authList.value = res.data?.list || res.list || []
      pagination.total = res.data?.total || res.total || 0
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchAuthList() }
const handleReset = () => { searchForm.status = undefined; searchForm.keyword = ''; searchForm.certify_type = ''; handleSearch() }

const handleSelectionChange = (rows: any[]) => { selectedRows.value = rows }

const handleApprove = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定通过此认证申请吗？', '提示')
    await request.post({ url: `/api/admin/certifications/${row.id}/review`, data: { status: 1 }, showSuccessMessage: true })
    fetchAuthList(); fetchStats()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('操作失败')
  }
}

const handleBatchApprove = async () => {
  try {
    await ElMessageBox.confirm(`确定批量通过 ${selectedRows.value.length} 条认证申请吗？`, '批量审核')
    const ids = selectedRows.value.filter(r => r.status === 0).map(r => r.id)
    await request.post({ url: '/api/admin/certifications/batch-review', data: { ids, status: 1 }, showSuccessMessage: true })
    fetchAuthList(); fetchStats()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('批量操作失败')
  }
}

const showRejectDialog = (row: any) => {
  currentAuth.value = row
  rejectReason.value = ''
  rejectDialogVisible.value = true
}

const handleReject = async () => {
  if (!rejectReason.value) {
    ElMessage.warning('请输入拒绝原因')
    return
  }
  actionLoading.value = true
  try {
    await request.post({
      url: `/api/admin/certifications/${currentAuth.value.id}/review`,
      data: { status: 2, reason: rejectReason.value },
      showSuccessMessage: true
    })
    rejectDialogVisible.value = false
    fetchAuthList(); fetchStats()
  } catch (error) {
    ElMessage.error('操作失败')
  } finally {
    actionLoading.value = false
  }
}

const handleViewDetail = (row: any) => {
  currentAuth.value = row
  detailDialogVisible.value = true
}

onMounted(() => { fetchAuthList(); fetchStats() })
</script>

<style scoped lang="scss">
.client-auth-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.header-actions { display: flex; gap: 8px; }
.search-form { margin-bottom: 16px; .el-form-item { margin-bottom: 0; } }
.reject-reason { font-size: 12px; color: var(--el-color-danger); margin-top: 4px; }
.stats-row { margin-bottom: 16px; }
.stat-card { text-align: center; }
.stat-value { font-size: 24px; font-weight: 700; margin-bottom: 4px; }
.stat-value.pending { color: var(--el-color-warning); }
.stat-value.approved { color: var(--el-color-success); }
.stat-value.rejected { color: var(--el-color-danger); }
.stat-value.total { color: var(--el-color-primary); }
.stat-label { color: var(--el-text-color-secondary); font-size: 13px; }
</style>
