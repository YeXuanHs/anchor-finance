<template>
  <div class="snapshot-manager">
    <div class="toolbar">
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        创建快照
      </el-button>
    </div>

    <el-table
      v-loading="loading"
      :data="snapshotList"
      style="width: 100%"
      empty-text="暂无快照"
    >
      <el-table-column prop="name" label="快照名称" min-width="160" show-overflow-tooltip />
      <el-table-column prop="description" label="描述" min-width="180" show-overflow-tooltip />
      <el-table-column prop="size" label="大小" width="100" align="center">
        <template #default="{ row }">
          {{ row.size ? `${row.size}GB` : '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getSnapshotStatusType(row.status)" size="small" effect="light">
            {{ getSnapshotStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="220" align="center" fixed="right">
        <template #default="{ row }">
          <el-button
            link
            type="primary"
            size="small"
            :disabled="row.status !== 'available'"
            @click="handleRestore(row)"
          >
            恢复
          </el-button>
          <el-button
            link
            type="danger"
            size="small"
            @click="handleDelete(row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div v-if="total > 0" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="fetchList"
        @current-change="fetchList"
      />
    </div>

    <!-- 创建快照弹窗 -->
    <el-dialog
      v-model="showCreateDialog"
      title="创建快照"
      width="480px"
      :close-on-click-modal="false"
      @closed="resetCreateForm"
    >
      <el-form :model="createForm" label-width="90px">
        <el-form-item label="快照名称">
          <el-input
            v-model="createForm.name"
            placeholder="请输入快照名称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="createForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入快照描述（可选）"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="createLoading" @click="confirmCreate">
          确认创建
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/request'

const route = useRoute()

const loading = ref(false)
const snapshotList = ref<any[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const showCreateDialog = ref(false)
const createLoading = ref(false)
const createForm = ref({
  name: '',
  description: ''
})

function getSnapshotStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    available: 'success',
    creating: 'warning',
    restoring: 'warning',
    error: 'danger'
  }
  return map[status] || 'info'
}

function getSnapshotStatusText(status: string) {
  const map: Record<string, string> = {
    available: '可用',
    creating: '创建中',
    restoring: '恢复中',
    error: '异常'
  }
  return map[status] || status
}

async function fetchList() {
  const id = route.params.id
  if (!id) return

  loading.value = true
  try {
    const { data } = await request.get(`/api/v2/hosts/${id}/snapshots`, {
      params: { page: page.value, limit: pageSize.value }
    })
    if (data?.data) {
      snapshotList.value = data.data.list || []
      total.value = data.data.total || 0
    }
  } catch (error) {
    console.error('获取快照列表失败', error)
  } finally {
    loading.value = false
  }
}

function resetCreateForm() {
  createForm.value = { name: '', description: '' }
}

async function confirmCreate() {
  const id = route.params.id
  if (!id) return

  if (!createForm.value.name.trim()) {
    ElMessage.warning('请输入快照名称')
    return
  }

  createLoading.value = true
  try {
    await request.post(`/api/v2/hosts/${id}/snapshots`, createForm.value)
    ElMessage.success('快照创建任务已提交')
    showCreateDialog.value = false
    fetchList()
  } catch (error: any) {
    ElMessage.error(error.message || '创建快照失败')
  } finally {
    createLoading.value = false
  }
}

async function handleRestore(row: any) {
  const id = route.params.id
  if (!id) return

  try {
    await ElMessageBox.confirm(
      `确认恢复快照「${row.name}」？恢复后服务器将重启。`,
      '确认恢复',
      { type: 'warning' }
    )
    await request.post(`/api/v2/hosts/${id}/snapshots/${row.id}/restore`)
    ElMessage.success('快照恢复任务已提交')
    fetchList()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '恢复快照失败')
    }
  }
}

async function handleDelete(row: any) {
  const id = route.params.id
  if (!id) return

  try {
    await ElMessageBox.confirm(
      `确认删除快照「${row.name}」？此操作不可恢复。`,
      '确认删除',
      { type: 'warning', confirmButtonText: '删除', confirmButtonClass: 'el-button--danger' }
    )
    await request.delete(`/api/v2/hosts/${id}/snapshots/${row.id}`)
    ElMessage.success('快照已删除')
    fetchList()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除快照失败')
    }
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.snapshot-manager {
  .toolbar {
    display: flex;
    justify-content: flex-end;
    margin-bottom: 16px;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}
</style>
