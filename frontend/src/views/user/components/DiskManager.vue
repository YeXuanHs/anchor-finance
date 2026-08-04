<template>
  <div class="disk-manager">
    <div class="toolbar">
      <el-button type="primary" @click="showExpandDialog = true" :disabled="!selectedDisk">
        <el-icon><Top /></el-icon>
        扩容选中磁盘
      </el-button>
    </div>

    <el-table
      v-loading="loading"
      :data="diskList"
      style="width: 100%"
      empty-text="暂无磁盘"
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="50" :selectable="isSelectable" />
      <el-table-column prop="name" label="磁盘名称" min-width="150" show-overflow-tooltip />
      <el-table-column prop="type" label="类型" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="row.type === 'system' ? 'primary' : 'success'" size="small" effect="light">
            {{ row.type === 'system' ? '系统盘' : '数据盘' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="size" label="容量" width="120" align="center">
        <template #default="{ row }">
          {{ row.size ? `${row.size}GB` : '-' }}
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getDiskStatusType(row.status)" size="small" effect="light">
            {{ getDiskStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="mount_point" label="挂载点" width="150">
        <template #default="{ row }">
          <span class="mono">{{ row.mount_point || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" align="center" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="!row.mounted && row.type !== 'system'"
            link
            type="primary"
            size="small"
            :disabled="row.status !== 'available'"
            @click="handleMount(row)"
          >
            挂载
          </el-button>
          <el-button
            v-if="row.mounted && row.type !== 'system'"
            link
            type="warning"
            size="small"
            @click="handleUnmount(row)"
          >
            卸载
          </el-button>
          <el-button
            link
            type="primary"
            size="small"
            @click="handleExpandSingle(row)"
          >
            扩容
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

    <!-- 扩容弹窗 -->
    <el-dialog
      v-model="showExpandDialog"
      title="磁盘扩容"
      width="480px"
      :close-on-click-modal="false"
    >
      <el-form label-width="100px">
        <el-form-item label="当前容量">
          <span>{{ expandTarget?.size ? `${expandTarget.size}GB` : '-' }}</span>
        </el-form-item>
        <el-form-item label="扩容至(GB)" required>
          <el-input-number
            v-model="expandSize"
            :min="(expandTarget?.size || 0) + 1"
            :max="expandTarget?.max_size || 1000"
            :step="10"
            style="width: 100%"
          />
        </el-form-item>
        <el-alert
          v-if="expandSize > (expandTarget?.size || 0)"
          type="warning"
          :closable="false"
          show-icon
          :description="`扩容后容量为 ${expandSize}GB，增加 ${expandSize - (expandTarget?.size || 0)}GB`"
        />
      </el-form>
      <template #footer>
        <el-button @click="showExpandDialog = false">取消</el-button>
        <el-button type="primary" :loading="expandLoading" @click="confirmExpand">
          确认扩容
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Top } from '@element-plus/icons-vue'
import request from '@/utils/request'

const route = useRoute()

const loading = ref(false)
const diskList = ref<any[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const selectedDisks = ref<any[]>([])
const selectedDisk = ref<any>(null)

const showExpandDialog = ref(false)
const expandLoading = ref(false)
const expandTarget = ref<any>(null)
const expandSize = ref(0)

function isSelectable(row: any) {
  return row.type !== 'system'
}

function handleSelectionChange(selection: any[]) {
  selectedDisks.value = selection
  selectedDisk.value = selection.length > 0 ? selection[0] : null
}

function getDiskStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    available: 'success',
    attached: 'success',
    detaching: 'warning',
    attaching: 'warning',
    error: 'danger'
  }
  return map[status] || 'info'
}

function getDiskStatusText(status: string) {
  const map: Record<string, string> = {
    available: '可用',
    attached: '已挂载',
    detaching: '卸载中',
    attaching: '挂载中',
    error: '异常'
  }
  return map[status] || status
}

async function fetchList() {
  const id = route.params.id
  if (!id) return

  loading.value = true
  try {
    const { data } = await request.get(`/api/v2/hosts/${id}/disks`, {
      params: { page: page.value, limit: pageSize.value }
    })
    if (data?.data) {
      diskList.value = data.data.list || []
      total.value = data.data.total || 0
    }
  } catch (error) {
    console.error('获取磁盘列表失败', error)
  } finally {
    loading.value = false
  }
}

async function handleMount(row: any) {
  const id = route.params.id
  if (!id) return

  try {
    await ElMessageBox.confirm(
      `确认挂载磁盘「${row.name}」？`,
      '确认挂载',
      { type: 'info' }
    )
    await request.post(`/api/v2/hosts/${id}/disks/${row.id}/mount`)
    ElMessage.success('磁盘挂载任务已提交')
    fetchList()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '挂载磁盘失败')
    }
  }
}

async function handleUnmount(row: any) {
  const id = route.params.id
  if (!id) return

  try {
    await ElMessageBox.confirm(
      `确认卸载磁盘「${row.name}」？卸载前请确保磁盘未被使用。`,
      '确认卸载',
      { type: 'warning' }
    )
    await request.post(`/api/v2/hosts/${id}/disks/${row.id}/unmount`)
    ElMessage.success('磁盘卸载任务已提交')
    fetchList()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '卸载磁盘失败')
    }
  }
}

function handleExpandSingle(row: any) {
  expandTarget.value = row
  expandSize.value = (row.size || 0) + 10
  showExpandDialog.value = true
}

async function confirmExpand() {
  const id = route.params.id
  if (!id || !expandTarget.value) return

  if (expandSize.value <= (expandTarget.value.size || 0)) {
    ElMessage.warning('扩容后的容量必须大于当前容量')
    return
  }

  expandLoading.value = true
  try {
    await request.post(`/api/v2/hosts/${id}/disks/${expandTarget.value.id}/resize`, {
      size: expandSize.value
    })
    ElMessage.success('磁盘扩容任务已提交')
    showExpandDialog.value = false
    fetchList()
  } catch (error: any) {
    ElMessage.error(error.message || '磁盘扩容失败')
  } finally {
    expandLoading.value = false
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.disk-manager {
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

.mono {
  font-family: 'Monaco', 'Menlo', monospace;
}
</style>
