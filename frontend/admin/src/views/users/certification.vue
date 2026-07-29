<template>
  <div class="certification-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="用户名/姓名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
        </el-form-item>
        <el-form-item label="证件类型">
          <el-select v-model="searchForm.id_type" placeholder="全部" clearable>
            <el-option label="身份证" value="id_card" />
            <el-option label="护照" value="passport" />
            <el-option label="营业执照" value="business_license" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>实名认证列表</h3>
        <div>
          <el-button type="success" :disabled="selectedRows.length === 0" @click="handleBatchApprove">批量通过</el-button>
          <el-button type="danger" :disabled="selectedRows.length === 0" @click="handleBatchReject">批量拒绝</el-button>
        </div>
      </div>

      <el-table :data="certifications" style="width: 100%" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user.username" label="用户名" width="120" />
        <el-table-column prop="real_name" label="真实姓名" width="120" />
        <el-table-column prop="id_type" label="证件类型" width="100">
          <template #default="{ row }">
            {{ getIdTypeText(row.id_type) }}
          </template>
        </el-table-column>
        <el-table-column prop="id_number" label="证件号码" width="180" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reviewer.username" label="审核人" width="100" />
        <el-table-column prop="reviewed_at" label="审核时间" width="180" />
        <el-table-column prop="created_at" label="申请时间" width="180" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewDetail(row)">详情</el-button>
            <template v-if="row.status === 'pending'">
              <el-button type="success" link @click="approve(row)">通过</el-button>
              <el-button type="danger" link @click="reject(row)">拒绝</el-button>
            </template>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchCertifications"
          @current-change="fetchCertifications"
        />
      </div>
    </div>

    <el-dialog v-model="showDetailDialog" title="认证详情" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="用户名">{{ currentCert.user?.username }}</el-descriptions-item>
        <el-descriptions-item label="真实姓名">{{ currentCert.real_name }}</el-descriptions-item>
        <el-descriptions-item label="证件类型">{{ getIdTypeText(currentCert.id_type) }}</el-descriptions-item>
        <el-descriptions-item label="证件号码">{{ currentCert.id_number }}</el-descriptions-item>
        <el-descriptions-item label="民族">{{ currentCert.nation || '-' }}</el-descriptions-item>
        <el-descriptions-item label="住址">{{ currentCert.address || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentCert.status)" size="small">
            {{ getStatusText(currentCert.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="审核人">{{ currentCert.reviewer?.username || '-' }}</el-descriptions-item>
        <el-descriptions-item label="申请时间">{{ currentCert.created_at }}</el-descriptions-item>
        <el-descriptions-item label="审核时间">{{ currentCert.reviewed_at || '-' }}</el-descriptions-item>
        <el-descriptions-item label="拒绝原因" :span="2" v-if="currentCert.reject_reason">
          {{ currentCert.reject_reason }}
        </el-descriptions-item>
      </el-descriptions>

      <div class="cert-images" v-if="currentCert.images?.length">
        <h4>证件照片</h4>
        <el-image
          v-for="(img, idx) in currentCert.images"
          :key="idx"
          :src="img"
          :preview-src-list="currentCert.images"
          :initial-index="idx"
          style="width: 200px; margin-right: 10px; margin-bottom: 10px;"
          fit="cover"
        />
      </div>

      <template #footer>
        <el-button @click="showDetailDialog = false">关闭</el-button>
        <template v-if="currentCert.status === 'pending'">
          <el-button type="success" @click="approve(currentCert)">通过</el-button>
          <el-button type="danger" @click="reject(currentCert)">拒绝</el-button>
        </template>
      </template>
    </el-dialog>

    <el-dialog v-model="showRejectDialog" title="拒绝原因" width="400px">
      <el-input v-model="rejectReason" type="textarea" :rows="4" placeholder="请输入拒绝原因" />
      <template #footer>
        <el-button @click="showRejectDialog = false">取消</el-button>
        <el-button type="danger" :loading="submitLoading" @click="handleReject">确定拒绝</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const certifications = ref<any[]>([])
const selectedRows = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const showDetailDialog = ref(false)
const showRejectDialog = ref(false)
const currentCert = ref<any>({})
const rejectReason = ref('')
const rejectTarget = ref<any>(null)

const searchForm = ref({ keyword: '', status: '', id_type: '' })

const getStatusType = (status: string) => {
  const map: Record<string, string> = { pending: 'warning', approved: 'success', rejected: 'danger' }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = { pending: '待审核', approved: '已通过', rejected: '已拒绝' }
  return map[status] || status
}

const getIdTypeText = (type: string) => {
  const map: Record<string, string> = { id_card: '身份证', passport: '护照', business_license: '营业执照' }
  return map[type] || type
}

const fetchCertifications = async () => {
  loading.value = true
  try {
    const params = { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    const { data } = await request.get('/admin/api/v1/users/certifications', { params })
    certifications.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch {
    ElMessage.error('获取认证列表失败')
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { currentPage.value = 1; fetchCertifications() }
const resetSearch = () => { searchForm.value = { keyword: '', status: '', id_type: '' }; handleSearch() }
const handleSelectionChange = (rows: any[]) => { selectedRows.value = rows }

const viewDetail = (cert: any) => {
  currentCert.value = cert
  showDetailDialog.value = true
}

const approve = async (cert: any) => {
  try {
    await ElMessageBox.confirm('确定要通过该认证申请吗？', '提示', { type: 'warning' })
    await request.put(`/admin/api/v1/users/certifications/${cert.id}/approve`)
    ElMessage.success('审核通过')
    showDetailDialog.value = false
    fetchCertifications()
  } catch {}
}

const reject = (cert: any) => {
  rejectTarget.value = cert
  rejectReason.value = ''
  showRejectDialog.value = true
}

const handleReject = async () => {
  if (!rejectReason.value.trim()) {
    ElMessage.warning('请输入拒绝原因')
    return
  }
  submitLoading.value = true
  try {
    await request.put(`/admin/api/v1/users/certifications/${rejectTarget.value.id}/reject`, { reason: rejectReason.value })
    ElMessage.success('已拒绝')
    showRejectDialog.value = false
    showDetailDialog.value = false
    fetchCertifications()
  } catch {
    ElMessage.error('操作失败')
  } finally {
    submitLoading.value = false
  }
}

const handleBatchApprove = async () => {
  try {
    await ElMessageBox.confirm(`确定要通过选中的 ${selectedRows.value.length} 条认证吗？`, '提示', { type: 'warning' })
    const ids = selectedRows.value.filter(r => r.status === 'pending').map(r => r.id)
    if (!ids.length) { ElMessage.warning('选中项中没有待审核的认证'); return }
    await request.post('/admin/api/v1/users/certifications/batch-approve', { ids })
    ElMessage.success('批量通过成功')
    fetchCertifications()
  } catch {}
}

const handleBatchReject = async () => {
  rejectTarget.value = null
  rejectReason.value = ''
  showRejectDialog.value = true
}

onMounted(fetchCertifications)
</script>

<style scoped lang="scss">
.certification-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
  .cert-images { margin-top: 20px; h4 { margin-bottom: 12px; } }
}
</style>
