<template>
  <div class="cancel-requests-page">
    <art-card :title="$t('cancelRequests.title')" shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('cancelRequests.headerTitle') }}</span>
          <el-button @click="showReasonDialog">{{ $t('cancelRequests.reasonManagement') }}</el-button>
        </div>
      </template>

      <el-table :data="cancelList" v-loading="loading" stripe border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" :label="$t('cancelRequests.columnUser')" />
        <el-table-column prop="create_time" :label="$t('cancelRequests.columnRequestTime')" width="180" />
        <el-table-column prop="product_name" :label="$t('cancelRequests.columnProduct')" />
        <el-table-column prop="ip" label="IP" width="140" />
        <el-table-column prop="type" :label="$t('cancelRequests.columnType')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'Immediate' ? 'danger' : 'warning'">
              {{ row.type === 'Immediate' ? $t('cancelRequests.typeImmediate') : $t('cancelRequests.typeDue') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reason" :label="$t('cancelRequests.columnReason')" show-overflow-tooltip />
        <el-table-column prop="delete_time" :label="$t('cancelRequests.columnDeleteTime')" width="180" />
        <el-table-column :label="$t('cancelRequests.columnStatus')" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? $t('cancelRequests.statusProcessed') : $t('cancelRequests.statusPending') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('cancelRequests.columnActions')" width="120" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="handleConfirm(row)" v-if="!row.status">{{ $t('cancelRequests.btnConfirmCancel') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.limit"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next"
        style="margin-top: 16px; justify-content: flex-end"
        @size-change="fetchCancelList"
        @current-change="fetchCancelList"
      />
    </art-card>

    <el-dialog v-model="reasonDialogVisible" :title="$t('cancelRequests.reasonManagement')" width="600px">
      <div class="reason-toolbar">
        <el-button type="primary" size="small" @click="showAddReason">{{ $t('cancelRequests.btnAddReason') }}</el-button>
      </div>
      <el-table :data="reasonList" stripe border size="small">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="reason" :label="$t('cancelRequests.columnReason')" />
        <el-table-column :label="$t('cancelRequests.columnActions')" width="100">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="handleDeleteReason(row)">{{ $t('cancelRequests.btnDelete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-dialog v-model="addReasonVisible" :title="$t('cancelRequests.btnAddReason')" width="400px" append-to-body>
        <el-input v-model="newReason" :placeholder="$t('cancelRequests.inputReasonPlaceholder')" />
        <template #footer>
          <el-button @click="addReasonVisible = false">{{ $t('cancelRequests.btnCancel') }}</el-button>
          <el-button type="primary" @click="handleAddReason">{{ $t('cancelRequests.btnConfirm') }}</el-button>
        </template>
      </el-dialog>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const reasonDialogVisible = ref(false)
const addReasonVisible = ref(false)
const newReason = ref('')

const pagination = reactive({ page: 1, limit: 10, total: 0 })
const cancelList = ref<any[]>([])
const reasonList = ref<any[]>([])

const fetchCancelList = async () => {
  loading.value = true
  try {
    const res = await request.get({
      url: '/api/admin/user-manage/cancel-requests',
      params: { page: pagination.page, limit: pagination.limit }
    })
    if (res?.data) {
      cancelList.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleConfirm = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('cancelRequests.confirmCancelMessage'), $t('cancelRequests.tip'))
    await request.post({ url: `/api/admin/user-manage/cancel-requests/${row.id}`, showSuccessMessage: true })
    fetchCancelList()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error($t('cancelRequests.operationFailed'))
  }
}

const showReasonDialog = async () => {
  reasonDialogVisible.value = true
  try {
    const res = await request.get({ url: '/api/admin/cancel-reasons' })
    reasonList.value = res?.data || []
  } catch (error) {
    console.error(error)
  }
}

const showAddReason = () => {
  newReason.value = ''
  addReasonVisible.value = true
}

const handleAddReason = async () => {
  if (!newReason.value) return
  try {
    await request.post({ url: '/api/admin/cancel-reasons', params: { reason: newReason.value }, showSuccessMessage: true })
    addReasonVisible.value = false
    showReasonDialog()
  } catch (error) {
    ElMessage.error($t('cancelRequests.addFailed'))
  }
}

const handleDeleteReason = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('cancelRequests.confirmDeleteMessage'), $t('cancelRequests.tip'))
    await request.del({ url: `/api/admin/cancel-reasons/${row.id}`, showSuccessMessage: true })
    showReasonDialog()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error($t('cancelRequests.deleteFailed'))
  }
}

onMounted(() => fetchCancelList())
</script>

<style scoped lang="scss">
.cancel-requests-page {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.reason-toolbar {
  margin-bottom: 12px;
}
</style>
