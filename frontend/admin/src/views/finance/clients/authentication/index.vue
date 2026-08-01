<template>
  <div class="client-auth-page">
    <art-card title="实名认证审核" shadow="never">
      <template #header>
        <div class="card-header">
          <span>实名认证审核</span>
        </div>
      </template>

      <el-table :data="authList" v-loading="loading" stripe border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="姓名" width="120" />
        <el-table-column prop="realname" label="认证名称" width="120" />
        <el-table-column prop="idcard" label="身份证号码" width="180" />
        <el-table-column prop="certify_type" label="认证方式" width="120">
          <template #default="{ row }">
            <el-tag>{{ row.certify_type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="认证类型" width="100" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
            <div v-if="row.reason" class="reject-reason">{{ row.reason }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="create_time" label="提交时间" width="180" />
        <el-table-column label="操作" width="200" fixed="right">
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
        :total="pagination.total"
        layout="total, sizes, prev, pager, next"
        style="margin-top: 16px; justify-content: flex-end"
        @size-change="fetchAuthList"
        @current-change="fetchAuthList"
      />
    </art-card>

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
    <el-dialog v-model="detailDialogVisible" title="认证详情" width="500px">
      <el-descriptions :column="1" border v-if="currentAuth">
        <el-descriptions-item label="用户ID">{{ currentAuth.id }}</el-descriptions-item>
        <el-descriptions-item label="姓名">{{ currentAuth.username }}</el-descriptions-item>
        <el-descriptions-item label="认证名称">{{ currentAuth.realname }}</el-descriptions-item>
        <el-descriptions-item label="身份证号">{{ currentAuth.idcard }}</el-descriptions-item>
        <el-descriptions-item label="认证方式">{{ currentAuth.certify_type }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentAuth.status)">{{ getStatusText(currentAuth.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="提交时间">{{ currentAuth.create_time }}</el-descriptions-item>
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

const pagination = reactive({ page: 1, limit: 10, total: 0 })
const authList = ref<any[]>([])

const getStatusType = (status: number) => {
  const map: Record<number, string> = { 0: 'warning', 1: 'success', 2: 'danger' }
  return map[status] || 'info'
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = { 0: '待审核', 1: '已通过', 2: '已拒绝' }
  return map[status] || '未知'
}

const fetchAuthList = async () => {
  loading.value = true
  try {
    const res = await request.get({
      url: '/api/admin/clients/authentication',
      params: { page: pagination.page, limit: pagination.limit }
    })
    if (res) {
      authList.value = res.data?.list || []
      pagination.total = res.data?.total || 0
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleApprove = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定通过此认证申请吗？', '提示')
    await request.put({ url: `/api/admin/clients/authentication/${row.id}/approve`, showSuccessMessage: true })
    fetchAuthList()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('操作失败')
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
    await request.put({
      url: `/api/admin/clients/authentication/${currentAuth.value.id}/reject`,
      data: { reason: rejectReason.value },
      showSuccessMessage: true
    })
    rejectDialogVisible.value = false
    fetchAuthList()
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

onMounted(() => fetchAuthList())
</script>

<style scoped lang="scss">
.client-auth-page {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.reject-reason {
  font-size: 12px;
  color: var(--el-color-danger);
  margin-top: 4px;
}
</style>
