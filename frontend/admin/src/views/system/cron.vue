<template>
  <div class="cron-page page-container">
    <el-tabs v-model="activeTab" type="border-card">
      <!-- 定时任务列表 -->
      <el-tab-pane label="定时任务" name="tasks">
        <div class="art-card">
          <div class="table-header">
            <h3>定时任务管理</h3>
            <el-button type="primary" @click="handleAdd">
              <el-icon><Plus /></el-icon>
              添加任务
            </el-button>
          </div>

          <el-table :data="tasks" style="width: 100%" v-loading="loading">
            <el-table-column prop="id" label="ID" width="70" />
            <el-table-column prop="name" label="任务名称" min-width="150" />
            <el-table-column prop="cron_expression" label="Cron表达式" width="160">
              <template #default="{ row }">
                <el-tag effect="plain" size="small">{{ row.cron_expression }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="last_run_at" label="上次执行" width="170" />
            <el-table-column prop="next_run_at" label="下次执行" width="170" />
            <el-table-column prop="enabled" label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-switch v-model="row.enabled" @change="handleToggle(row)" />
              </template>
            </el-table-column>
            <el-table-column label="操作" width="220" fixed="right">
              <template #default="{ row }">
                <el-button type="primary" link @click="handleRun(row)">
                  <el-icon><CaretRight /></el-icon>
                  执行
                </el-button>
                <el-button type="primary" link @click="handleEdit(row)">编辑</el-button>
                <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>

          <div class="pagination">
            <el-pagination
              v-model:current-page="currentPage"
              v-model:page-size="pageSize"
              :page-sizes="[10, 20, 50]"
              :total="total"
              layout="total, sizes, prev, pager, next"
              @size-change="fetchTasks"
              @current-change="fetchTasks"
            />
          </div>
        </div>
      </el-tab-pane>

      <!-- 执行日志 -->
      <el-tab-pane label="执行日志" name="logs">
        <div class="art-card">
          <div class="table-header">
            <h3>执行日志</h3>
            <el-button @click="fetchLogs">
              <el-icon><Refresh /></el-icon>
              刷新
            </el-button>
          </div>

          <el-table :data="logs" style="width: 100%" v-loading="logsLoading">
            <el-table-column prop="id" label="ID" width="70" />
            <el-table-column prop="task_name" label="任务名称" min-width="150" />
            <el-table-column prop="started_at" label="开始时间" width="170" />
            <el-table-column prop="duration" label="耗时" width="100">
              <template #default="{ row }">
                <span>{{ row.duration }}ms</span>
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100" align="center">
              <template #default="{ row }">
                <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small" effect="light">
                  {{ row.status === 'success' ? '成功' : '失败' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="output" label="输出" min-width="200" show-overflow-tooltip />
          </el-table>

          <div class="pagination">
            <el-pagination
              v-model:current-page="logPage"
              v-model:page-size="logPageSize"
              :page-sizes="[10, 20, 50, 100]"
              :total="logTotal"
              layout="total, sizes, prev, pager, next"
              @size-change="fetchLogs"
              @current-change="fetchLogs"
            />
          </div>
        </div>
      </el-tab-pane>
    </el-tabs>

    <!-- 编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEditing ? '编辑任务' : '添加任务'"
      width="540px"
      destroy-on-close
    >
      <el-form
        ref="dialogFormRef"
        :model="dialogForm"
        :rules="dialogRules"
        label-width="110px"
      >
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="dialogForm.name" placeholder="请输入任务名称" clearable />
        </el-form-item>
        <el-form-item label="Cron表达式" prop="cron_expression">
          <el-input v-model="dialogForm.cron_expression" placeholder="如: 0 */5 * * * *" clearable />
          <div class="cron-hint">
            格式：秒 分 时 日 月 周
          </div>
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="dialogForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入任务描述"
          />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="dialogForm.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="dialogLoading" @click="handleDialogSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Refresh, CaretRight } from '@element-plus/icons-vue'
import request from '@/utils/request'

interface CronTask {
  id: number
  name: string
  cron_expression: string
  description: string
  last_run_at: string
  next_run_at: string
  enabled: boolean
}

interface CronLog {
  id: number
  task_name: string
  started_at: string
  duration: number
  status: 'success' | 'failed'
  output: string
}

const activeTab = ref('tasks')
const loading = ref(false)
const logsLoading = ref(false)
const dialogVisible = ref(false)
const dialogLoading = ref(false)
const isEditing = ref(false)
const dialogFormRef = ref<FormInstance>()
const editingTask = ref<CronTask | null>(null)

const tasks = ref<CronTask[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const logs = ref<CronLog[]>([])
const logPage = ref(1)
const logPageSize = ref(20)
const logTotal = ref(0)

const dialogForm = reactive({
  name: '',
  cron_expression: '',
  description: '',
  enabled: true
})

const dialogRules: FormRules = {
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  cron_expression: [{ required: true, message: '请输入Cron表达式', trigger: 'blur' }]
}

const fetchTasks = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/admin/system/cron', {
      params: { page: currentPage.value, page_size: pageSize.value }
    })
    if (data?.data) {
      tasks.value = data.data.items || []
      total.value = data.data.total || 0
    }
  } catch {
    tasks.value = []
  } finally {
    loading.value = false
  }
}

const fetchLogs = async () => {
  logsLoading.value = true
  try {
    const { data } = await request.get('/api/admin/system/cron/logs', {
      params: { page: logPage.value, page_size: logPageSize.value }
    })
    if (data?.data) {
      logs.value = data.data.items || []
      logTotal.value = data.data.total || 0
    }
  } catch {
    logs.value = []
  } finally {
    logsLoading.value = false
  }
}

const handleAdd = () => {
  isEditing.value = false
  editingTask.value = null
  dialogForm.name = ''
  dialogForm.cron_expression = ''
  dialogForm.description = ''
  dialogForm.enabled = true
  dialogVisible.value = true
}

const handleEdit = (row: CronTask) => {
  isEditing.value = true
  editingTask.value = row
  dialogForm.name = row.name
  dialogForm.cron_expression = row.cron_expression
  dialogForm.description = row.description
  dialogForm.enabled = row.enabled
  dialogVisible.value = true
}

const handleRun = async (row: CronTask) => {
  try {
    await ElMessageBox.confirm(`确定立即执行任务「${row.name}」吗？`, '手动执行', {
      confirmButtonText: '执行',
      cancelButtonText: '取消',
      type: 'info'
    })
    await request.post(`/api/admin/system/cron/${row.id}/run`)
    ElMessage.success('任务已触发执行')
    fetchTasks()
    fetchLogs()
  } catch {
    // 取消操作
  }
}

const handleToggle = async (row: CronTask) => {
  try {
    await request.post(`/api/admin/system/cron/${row.id}/toggle`, { enabled: row.enabled })
    ElMessage.success(`任务已${row.enabled ? '启用' : '禁用'}`)
  } catch {
    row.enabled = !row.enabled
    ElMessage.error('操作失败')
  }
}

const handleDelete = async (row: CronTask) => {
  try {
    await ElMessageBox.confirm(`确定删除任务「${row.name}」吗？此操作不可恢复。`, '删除确认', {
      confirmButtonText: '删除',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.delete(`/api/admin/system/cron/${row.id}`)
    ElMessage.success('任务已删除')
    fetchTasks()
  } catch {
    // 取消操作
  }
}

const handleDialogSubmit = async () => {
  const valid = await dialogFormRef.value?.validate().catch(() => false)
  if (!valid) return

  dialogLoading.value = true
  try {
    if (isEditing.value && editingTask.value) {
      await request.put(`/api/admin/system/cron/${editingTask.value.id}`, { ...dialogForm })
      Object.assign(editingTask.value, dialogForm)
    } else {
      await request.post('/api/admin/system/cron', { ...dialogForm })
      await fetchTasks()
    }
    dialogVisible.value = false
    ElMessage.success(isEditing.value ? '任务已更新' : '任务已添加')
  } catch {
    ElMessage.error('操作失败，请重试')
  } finally {
    dialogLoading.value = false
  }
}

onMounted(() => {
  fetchTasks()
  fetchLogs()
})
</script>

<style scoped lang="scss">
.cron-page {
  :deep(.el-tabs__content) {
    padding: 0;
  }

  .art-card {
    padding: 20px;
  }

  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;

    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
    }
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }

  .cron-hint {
    margin-top: 4px;
    font-size: 12px;
    color: var(--el-text-color-placeholder);
  }
}
</style>
