<template>
  <div class="crons-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('cron.title') }}</span>
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            {{ $t('cron.addTask') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('cron.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('cron.taskNamePlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('cron.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('cron.all')" clearable>
            <el-option :label="$t('cron.enabled')" :value="1" />
            <el-option :label="$t('cron.disabled')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('cron.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('cron.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="name" :label="$t('cron.taskName')" min-width="150" />
        <el-table-column prop="command" :label="$t('cron.executeCommand')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="cron_expression" :label="$t('cron.cronExpression')" width="140">
          <template #default="{ row }">
            <el-tag effect="plain" size="small">{{ row.cron_expression }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('cron.status')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? $t('cron.enabled') : $t('cron.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="last_run_at" :label="$t('cron.lastRun')" width="180" />
        <el-table-column prop="last_run_status" :label="$t('cron.runResult')" width="100">
          <template #default="{ row }">
            <el-tag
              v-if="row.last_run_status !== undefined"
              :type="row.last_run_status === 1 ? 'success' : 'danger'"
              size="small"
            >
              {{ row.last_run_status === 1 ? $t('cron.success') : $t('cron.failure') }}
            </el-tag>
            <span v-else class="empty-text">-</span>
          </template>
        </el-table-column>
        <el-table-column prop="next_run_at" :label="$t('cron.nextRun')" width="180" />
        <el-table-column :label="$t('cron.operations')" width="250" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleEdit(row)">{{ $t('cron.edit') }}</el-button>
            <el-button type="success" link @click="handleRun(row)">{{ $t('cron.run') }}</el-button>
            <el-button type="warning" link @click="handleViewLogs(row)">{{ $t('cron.viewLogs') }}</el-button>
            <el-popconfirm :title="$t('cron.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link>{{ $t('cron.delete') }}</el-button>
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

    <!-- 添加/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="600px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="120px">
        <el-form-item :label="$t('cron.taskName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('cron.enterTaskName')" />
        </el-form-item>
        <el-form-item :label="$t('cron.executeCommand')" prop="command">
          <el-input v-model="formData.command" :placeholder="$t('cron.enterCommand')" />
        </el-form-item>
        <el-form-item :label="$t('cron.cronExpression')" prop="cron_expression">
          <el-input v-model="formData.cron_expression" :placeholder="$t('cron.cronPlaceholder')" />
          <div class="form-tip">{{ $t('cron.cronTip') }}</div>
        </el-form-item>
        <el-form-item :label="$t('cron.description')" prop="description">
          <el-input v-model="formData.description" type="textarea" :rows="3" :placeholder="$t('cron.enterDescription')" />
        </el-form-item>
        <el-form-item :label="$t('cron.timeout')" prop="timeout">
          <el-input-number v-model="formData.timeout" :min="0" :max="86400" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('cron.status')" prop="status">
          <el-switch v-model="formData.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('cron.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('cron.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 日志对话框 -->
    <el-dialog v-model="logDialogVisible" :title="$t('cron.taskLogs')" width="800px" top="5vh">
      <el-table :data="logData" v-loading="logLoading" style="width: 100%" max-height="400">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="started_at" :label="$t('cron.startTime')" width="180" />
        <el-table-column prop="finished_at" :label="$t('cron.endTime')" width="180" />
        <el-table-column prop="duration" :label="$t('cron.duration')" width="100">
          <template #default="{ row }">
            {{ row.duration ? `${row.duration}s` : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('cron.status')" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? $t('cron.success') : $t('cron.failure') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="output" :label="$t('cron.output')" min-width="200">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleViewOutput(row)">{{ $t('cron.viewOutput') }}</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>

    <!-- 输出详情对话框 -->
    <el-dialog v-model="outputDialogVisible" :title="$t('cron.executeOutput')" width="700px">
      <pre class="output-content">{{ currentOutput }}</pre>
      <template #footer>
        <el-button @click="outputDialogVisible = false">{{ $t('cron.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

interface Cron {
  id: number
  name: string
  command: string
  cron_expression: string
  description: string
  timeout: number
  status: number
  last_run_at: string
  last_run_status: number
  next_run_at: string
}

interface CronLog {
  id: number
  started_at: string
  finished_at: string
  duration: number
  status: number
  output: string
}

const loading = ref(false)
const submitLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref($t('cron.addTask'))
const formRef = ref<FormInstance>()

const logDialogVisible = ref(false)
const logLoading = ref(false)
const logData = ref<CronLog[]>([])

const outputDialogVisible = ref(false)
const currentOutput = ref('')

const searchForm = reactive({
  keyword: '',
  status: undefined as number | undefined
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref<Cron[]>([])

const formData = reactive({
  id: undefined as number | undefined,
  name: '',
  command: '',
  cron_expression: '',
  description: '',
  timeout: 300,
  status: 1
})

const formRules: FormRules = {
  name: [
    { required: true, message: () => $t('cron.enterTaskName'), trigger: 'blur' }
  ],
  command: [
    { required: true, message: () => $t('cron.enterCommand'), trigger: 'blur' }
  ],
  cron_expression: [
    { required: true, message: () => $t('cron.enterCronExpression'), trigger: 'blur' }
  ]
}

// 获取任务列表
const fetchCrons = async () => {
  loading.value = true
  try {
    const data = await request.get({
      url: '/api/admin/cron-tasks',
      params: {
        page: pagination.page,
        page_size: pagination.page_size,
        ...searchForm
      }
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取定时任务列表失败:', error)
    ElMessage.error($t('cron.fetchFailed'))
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchCrons()
}

// 重置
const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = undefined
  handleSearch()
}

// 添加
const handleAdd = () => {
  dialogTitle.value = $t('cron.addTask')
  formData.id = undefined
  formData.name = ''
  formData.command = ''
  formData.cron_expression = ''
  formData.description = ''
  formData.timeout = 300
  formData.status = 1
  dialogVisible.value = true
}

// 编辑
const handleEdit = (row: any) => {
  dialogTitle.value = $t('cron.editTask')
  Object.assign(formData, row)
  dialogVisible.value = true
}

// 手动执行
const handleRun = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('cron.confirmRunWithName', { name: row.name }), $t('cron.prompt'), {
      type: 'warning'
    })
    await request.post({
      url: `/api/admin/cron-tasks/${row.id}/run`,
      showSuccessMessage: true
    })
    ElMessage.success($t('cron.taskTriggered'))
    fetchCrons()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error($t('cron.runFailed'))
    }
  }
}

// 查看日志
const handleViewLogs = async (row: any) => {
  logDialogVisible.value = true
  logLoading.value = true
  try {
    const data = await request.get({
      url: `/api/admin/cron-tasks/${row.id}/logs`
    })
    logData.value = data || []
  } catch (error) {
    console.error('获取日志失败:', error)
    ElMessage.error($t('cron.fetchLogsFailed'))
  } finally {
    logLoading.value = false
  }
}

// 查看输出
const handleViewOutput = (row: any) => {
  currentOutput.value = row.output || $t('cron.noOutput')
  outputDialogVisible.value = true
}

// 删除
const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/cron-tasks/${row.id}`
    })
    ElMessage.success($t('cron.deleteSuccess'))
    fetchCrons()
  } catch (error) {
    ElMessage.error($t('cron.deleteFailed'))
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({
          url: `/api/admin/cron-tasks/${formData.id}`,
          params: { ...formData },
          showSuccessMessage: true
        })
      } else {
        await request.post({
          url: '/api/admin/cron-tasks',
          params: { ...formData },
          showSuccessMessage: true
        })
      }

      ElMessage.success(formData.id ? $t('cron.updateSuccess') : $t('cron.addSuccess'))
      dialogVisible.value = false
      fetchCrons()
    } catch (error) {
      ElMessage.error($t('cron.operationFailed'))
    } finally {
      submitLoading.value = false
    }
  })
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchCrons()
}

// 页码变化
const handlePageChange = () => {
  fetchCrons()
}

onMounted(() => {
  fetchCrons()
})
</script>

<style scoped lang="scss">
.crons-page {
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

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.form-tip {
  font-size: 12px;
  color: var(--el-text-color-placeholder);
  margin-top: 4px;
}

.empty-text {
  color: var(--el-text-color-placeholder);
  font-size: 13px;
}

.output-content {
  max-height: 400px;
  padding: 16px;
  overflow-y: auto;
  background-color: var(--el-fill-color-light);
  border-radius: 6px;
  font-family: Consolas, Monaco, 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  white-space: pre-wrap;
  word-break: break-all;
}
</style>
