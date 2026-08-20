<template>
  <div class="data-migration-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('dataMigration.title') }}</span>
        </div>
      </template>

      <!-- 迁移任务列表 -->
      <div class="section">
        <div class="section-header">
          <h3>{{ $t('dataMigration.migrationTasks') }}</h3>
          <el-button type="primary" size="small" @click="handleCreateTask">{{ $t('dataMigration.createMigrationTask') }}</el-button>
        </div>
        <el-table :data="tasks" v-loading="loading" style="width: 100%" border>
          <el-table-column prop="id" :label="$t('dataMigration.id')" width="80" />
          <el-table-column prop="name" :label="$t('dataMigration.taskName')" min-width="150" />
          <el-table-column prop="source_type" :label="$t('dataMigration.dataSource')" width="120" />
          <el-table-column prop="target_type" :label="$t('dataMigration.target')" width="120" />
          <el-table-column prop="status" :label="$t('dataMigration.status')" width="100">
            <template #default="{ row }">
              <el-tag :type="statusTypeMap[row.status]" size="small">
                {{ statusLabelMap[row.status as keyof typeof statusLabelMap] }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="progress" :label="$t('dataMigration.progress')" width="200">
            <template #default="{ row }">
              <el-progress :percentage="row.progress" :status="row.status === 'completed' ? 'success' : undefined" />
            </template>
          </el-table-column>
          <el-table-column prop="created_at" :label="$t('dataMigration.createTime')" width="180" />
          <el-table-column :label="$t('dataMigration.operations')" width="200" fixed="right">
            <template #default="{ row }">
              <el-button v-if="row.status === 'pending'" type="success" link @click="handleStart(row)">{{ $t('dataMigration.start') }}</el-button>
              <el-button v-if="row.status === 'running'" type="warning" link @click="handlePause(row)">{{ $t('dataMigration.pause') }}</el-button>
              <el-button type="primary" link @click="handleViewLog(row)">{{ $t('dataMigration.logs') }}</el-button>
              <el-popconfirm v-if="row.status !== 'running'" :title="$t('dataMigration.confirmDelete')" @confirm="handleDelete(row)">
                <template #reference>
                  <el-button type="danger" link>{{ $t('dataMigration.delete') }}</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <!-- 创建任务对话框 -->
    <el-dialog v-model="dialogVisible" :title="$t('dataMigration.createTaskTitle')" width="600px">
      <el-form :model="formData" label-width="100px">
        <el-form-item :label="$t('dataMigration.taskName')" required>
          <el-input v-model="formData.name" :placeholder="$t('dataMigration.inputTaskName')" />
        </el-form-item>
        <el-form-item :label="$t('dataMigration.dataSourceType')" required>
          <el-select v-model="formData.source_type" :placeholder="$t('dataMigration.pleaseSelect')">
            <el-option label="WHMCS" value="whmcs" />
            <el-option label="Blesta" value="blesta" />
            <el-option :label="$t('dataMigration.csvFile')" value="csv" />
            <el-option :label="$t('dataMigration.database')" value="database" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('dataMigration.dataSourceConfig')">
          <el-input v-model="formData.source_config" type="textarea" :rows="4" :placeholder="$t('dataMigration.jsonConfigPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('dataMigration.migrationType')">
          <el-checkbox-group v-model="formData.migration_types">
            <el-checkbox value="users">{{ $t('dataMigration.users') }}</el-checkbox>
            <el-checkbox value="products">{{ $t('dataMigration.products') }}</el-checkbox>
            <el-checkbox value="orders">{{ $t('dataMigration.orders') }}</el-checkbox>
            <el-checkbox value="tickets">{{ $t('dataMigration.tickets') }}</el-checkbox>
            <el-checkbox value="invoices">{{ $t('dataMigration.invoices') }}</el-checkbox>
          </el-checkbox-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('dataMigration.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit">{{ $t('dataMigration.create') }}</el-button>
      </template>
    </el-dialog>

    <!-- 日志查看对话框 -->
    <el-dialog v-model="logVisible" :title="$t('dataMigration.migrationLog')" width="700px">
      <pre style="max-height: 500px; overflow-y: auto; background: #f5f7fa; padding: 16px; border-radius: 8px; font-size: 13px; line-height: 1.6; white-space: pre-wrap; word-break: break-all;">{{ logData }}</pre>
      <template #footer>
        <el-button @click="logVisible = false">{{ $t('dataMigration.close') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'

const statusTypeMap: Record<string, any> = {
  pending: 'info',
  running: 'warning',
  completed: 'success',
  failed: 'danger',
  paused: 'info'
}

const statusLabelMap = computed(() => ({
  pending: $t('dataMigration.statusPending'),
  running: $t('dataMigration.statusRunning'),
  completed: $t('dataMigration.statusCompleted'),
  failed: $t('dataMigration.statusFailed'),
  paused: $t('dataMigration.statusPaused')
}))

const loading = ref(false)
const dialogVisible = ref(false)
const tasks = ref<any[]>([])

const formData = reactive({
  name: '',
  source_type: 'whmcs',
  source_config: '',
  migration_types: ['users'] as string[]
})

const fetchTasks = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/system/data-migration' })
    tasks.value = data || []
  } catch (error) {
    ElMessage.error($t('dataMigration.fetchTasksFailed'))
  } finally {
    loading.value = false
  }
}

const handleCreateTask = () => {
  formData.name = ''
  formData.source_type = 'whmcs'
  formData.source_config = ''
  formData.migration_types = ['users']
  dialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formData.name) { ElMessage.warning($t('dataMigration.inputTaskNameWarning')); return }
  try {
    await request.post({ url: '/api/admin/system/data-migration', params: { ...formData } })
    ElMessage.success($t('dataMigration.taskCreated'))
    dialogVisible.value = false
    fetchTasks()
  } catch (error) {
    ElMessage.error($t('dataMigration.createFailed'))
  }
}

const handleStart = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/system/data-migration/tasks/${row.id}/start` })
    ElMessage.success($t('dataMigration.taskStarted'))
    fetchTasks()
  } catch (error) {
    ElMessage.error($t('dataMigration.operationFailed'))
  }
}

const handlePause = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/system/data-migration/tasks/${row.id}/pause` })
    ElMessage.success($t('dataMigration.taskPaused'))
    fetchTasks()
  } catch (error) {
    ElMessage.error($t('dataMigration.operationFailed'))
  }
}

const handleViewLog = (row: any) => {
  logData.value = row.logs || $t('dataMigration.noLogData')
  logVisible.value = true
}

const logVisible = ref(false)
const logData = ref('')

const handleDelete = async (row: any) => {
  try {
    await request.del({ url: `/api/admin/system/data-migration/tasks/${row.id}` })
    ElMessage.success($t('dataMigration.deleteSuccess'))
    fetchTasks()
  } catch (error) {
    ElMessage.error($t('dataMigration.deleteFailed'))
  }
}

onMounted(() => { fetchTasks() })
</script>

<style scoped lang="scss">
.data-migration-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.section { margin-top: 24px; }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.section-header h3 { margin: 0; font-size: 16px; font-weight: 600; }
</style>
