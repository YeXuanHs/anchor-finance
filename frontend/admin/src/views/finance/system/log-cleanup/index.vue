<template>
  <div class="log-cleanup-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('logCleanup.title') }}</span>
        </div>
      </template>

      <!-- 清理规则列表 -->
      <div class="section">
        <div class="section-header">
          <h3>{{ $t('logCleanup.autoCleanRules') }}</h3>
          <el-button type="primary" size="small" @click="handleAddRule">{{ $t('logCleanup.addRule') }}</el-button>
        </div>
        <el-table :data="rules" v-loading="loading" style="width: 100%" border>
          <el-table-column prop="id" :label="$t('logCleanup.id')" width="70" />
          <el-table-column prop="name" :label="$t('logCleanup.ruleName')" min-width="150" />
          <el-table-column prop="log_type" :label="$t('logCleanup.logType')" width="120">
            <template #default="{ row }">
              <el-tag size="small">{{ logTypeMap[row.log_type as keyof typeof logTypeMap] || row.log_type }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="retention_days" :label="$t('logCleanup.retentionDays')" width="100" align="center" />
          <el-table-column prop="max_count" :label="$t('logCleanup.maxCount')" width="100" align="center">
            <template #default="{ row }">{{ row.max_count || $t('logCleanup.unlimited') }}</template>
          </el-table-column>
          <el-table-column prop="schedule" :label="$t('logCleanup.schedule')" width="100" />
          <el-table-column prop="last_run_at" :label="$t('logCleanup.lastRun')" width="180" />
          <el-table-column prop="deleted_count" :label="$t('logCleanup.cleanedCount')" width="100" align="center" />
          <el-table-column prop="is_enabled" :label="$t('logCleanup.status')" width="80">
            <template #default="{ row }">
              <el-tag :type="row.is_enabled ? 'success' : 'info'" size="small">{{ row.is_enabled ? $t('logCleanup.enabled') : $t('logCleanup.disabled') }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column :label="$t('logCleanup.operations')" width="220" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="handleEditRule(row)">{{ $t('logCleanup.edit') }}</el-button>
              <el-button type="success" link @click="handleRunNow(row)">{{ $t('logCleanup.runNow') }}</el-button>
              <el-popconfirm :title="$t('logCleanup.confirmDelete')" @confirm="handleDeleteRule(row)">
                <template #reference><el-button type="danger" link>{{ $t('logCleanup.delete') }}</el-button></template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 手动清理 -->
      <div class="section">
        <div class="section-header">
          <h3>{{ $t('logCleanup.manualClean') }}</h3>
        </div>
        <el-card shadow="hover" class="manual-clean-card">
          <el-form :model="manualForm" label-width="100px" inline>
            <el-form-item :label="$t('logCleanup.logType')">
              <el-select v-model="manualForm.log_type" :placeholder="$t('logCleanup.all')" clearable style="width: 150px">
                <el-option v-for="(label, key) in logTypeMap" :key="key" :label="label" :value="key" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('logCleanup.logStatus')">
              <el-select v-model="manualForm.status" :placeholder="$t('logCleanup.all')" clearable style="width: 120px">
                <el-option :label="$t('logCleanup.success')" value="success" />
                <el-option :label="$t('logCleanup.failure')" value="error" />
              </el-select>
            </el-form-item>
            <el-form-item :label="$t('logCleanup.cleanMethod')">
              <el-radio-group v-model="manualForm.clean_by">
                <el-radio value="days">{{ $t('logCleanup.byDays') }}</el-radio>
                <el-radio value="count">{{ $t('logCleanup.byCount') }}</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item :label="manualForm.clean_by === 'days' ? $t('logCleanup.cleanRetentionDays') : $t('logCleanup.cleanRetentionCount')">
              <el-input-number v-if="manualForm.clean_by === 'days'" v-model="manualForm.days" :min="1" :max="365" />
              <el-input-number v-else v-model="manualForm.count" :min="100" :step="1000" />
            </el-form-item>
            <el-form-item>
              <el-button type="danger" @click="handleManualClean" :loading="manualLoading">{{ $t('logCleanup.executeClean') }}</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </div>

      <!-- 清理统计 -->
      <div class="section">
        <h3>{{ $t('logCleanup.cleanStats') }}</h3>
        <el-row :gutter="20">
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-value">{{ stats.total_cleaned }}</div>
              <div class="stat-label">{{ $t('logCleanup.totalCleaned') }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-value">{{ stats.freed_space }}</div>
              <div class="stat-label">{{ $t('logCleanup.freedSpace') }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-value">{{ stats.last_cleanup_at || '-' }}</div>
              <div class="stat-label">{{ $t('logCleanup.lastCleanTime') }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-value">{{ stats.next_cleanup_at || '-' }}</div>
              <div class="stat-label">{{ $t('logCleanup.nextCleanTime') }}</div>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-card>

    <!-- 添加/编辑规则对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item :label="$t('logCleanup.ruleName')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('logCleanup.inputRuleName')" />
        </el-form-item>
        <el-form-item :label="$t('logCleanup.logType')" prop="log_type">
          <el-select v-model="formData.log_type" :placeholder="$t('logCleanup.selectLogType')" style="width: 100%">
            <el-option v-for="(label, key) in logTypeMap" :key="key" :label="label" :value="key" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('logCleanup.retentionDays')" prop="retention_days">
          <el-input-number v-model="formData.retention_days" :min="1" :max="365" />
          <span class="form-tip">{{ $t('logCleanup.retentionDaysTip') }}</span>
        </el-form-item>
        <el-form-item :label="$t('logCleanup.maxCount')">
          <el-input-number v-model="formData.max_count" :min="0" :step="1000" />
          <span class="form-tip">{{ $t('logCleanup.maxCountTip') }}</span>
        </el-form-item>
        <el-form-item :label="$t('logCleanup.schedule')" prop="schedule">
          <el-select v-model="formData.schedule" :placeholder="$t('logCleanup.selectSchedule')" style="width: 100%">
            <el-option :label="$t('logCleanup.daily')" value="daily" />
            <el-option :label="$t('logCleanup.weekly')" value="weekly" />
            <el-option :label="$t('logCleanup.monthly')" value="monthly" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('logCleanup.status')">
          <el-switch v-model="formData.is_enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('logCleanup.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">{{ $t('logCleanup.confirm') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const logTypeMap = computed(() => ({
  operation: $t('logCleanup.logTypeOperation'),
  login: $t('logCleanup.logTypeLogin'),
  system: $t('logCleanup.logTypeSystem'),
  api: $t('logCleanup.logTypeApi'),
  error: $t('logCleanup.logTypeError'),
  ticket: $t('logCleanup.logTypeTicket'),
  order: $t('logCleanup.logTypeOrder'),
  all: $t('logCleanup.all')
}))

const loading = ref(false)
const submitLoading = ref(false)
const manualLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref($t('logCleanup.addRuleTitle'))
const formRef = ref<FormInstance>()

const rules = ref<any[]>([])
const stats = reactive({ total_cleaned: 0, freed_space: '0 MB', last_cleanup_at: '', next_cleanup_at: '' })

const formData = reactive({
  id: undefined as number | undefined,
  name: '', log_type: 'operation', retention_days: 30, max_count: 0, schedule: 'daily', is_enabled: true
})

const manualForm = reactive({ log_type: '', status: '', clean_by: 'days', days: 30, count: 1000 })

const formRules = computed<FormRules>(() => ({
  name: [{ required: true, message: $t('logCleanup.inputRuleName'), trigger: 'blur' }],
  log_type: [{ required: true, message: $t('logCleanup.selectLogType'), trigger: 'change' }],
  retention_days: [{ required: true, message: $t('logCleanup.retentionDaysTip'), trigger: 'blur' }],
  schedule: [{ required: true, message: $t('logCleanup.selectSchedule'), trigger: 'change' }]
}))

const fetchRules = async () => {
  loading.value = true
  try { rules.value = await request.get({ url: '/api/admin/log-cleaner/rules' }) || [] }
  catch { ElMessage.error($t('logCleanup.fetchRulesFailed')) } finally { loading.value = false }
}

const fetchStats = async () => {
  try { const data = await request.get({ url: '/api/admin/log-cleaner/stats' }); if (data) Object.assign(stats, data) }
  catch { /* ignore */ }
}

const handleAddRule = () => {
  dialogTitle.value = $t('logCleanup.addRuleTitle')
  Object.assign(formData, { id: undefined, name: '', log_type: 'operation', retention_days: 30, max_count: 0, schedule: 'daily', is_enabled: true })
  dialogVisible.value = true
}

const handleEditRule = (row: any) => {
  dialogTitle.value = $t('logCleanup.editRuleTitle')
  Object.assign(formData, { id: row.id, name: row.name, log_type: row.log_type, retention_days: row.retention_days, max_count: row.max_count || 0, schedule: row.schedule, is_enabled: row.is_enabled })
  dialogVisible.value = true
}

const handleDeleteRule = async (row: any) => {
  try { await request.del({ url: `/api/admin/log-cleaner/rules/${row.id}` }); ElMessage.success($t('logCleanup.deleteSuccess')); fetchRules() }
  catch { ElMessage.error($t('logCleanup.deleteFailed')) }
}

const handleRunNow = async (row: any) => {
  try {
    await ElMessageBox.confirm($t('logCleanup.confirmRunNow', { name: row.name }), $t('logCleanup.confirmRunTitle'))
    await request.post({ url: `/api/admin/log-cleaner/rules/${row.id}/run` })
    ElMessage.success($t('logCleanup.cleanTaskSubmitted')); fetchRules(); fetchStats()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error($t('logCleanup.runFailed')) }
}

const handleManualClean = async () => {
  try {
    const logTypeLabel = manualForm.log_type ? logTypeMap.value[manualForm.log_type as keyof typeof logTypeMap.value] : $t('logCleanup.all')
    const msg = manualForm.clean_by === 'days'
      ? `${$t('logCleanup.confirmCleanDays', { days: manualForm.days, logType: logTypeLabel })}`
      : `${$t('logCleanup.confirmCleanCount', { count: manualForm.count, logType: logTypeLabel })}`
    await ElMessageBox.confirm(msg, $t('logCleanup.manualCleanConfirm'))
    manualLoading.value = true
    const params: any = { clean_by: manualForm.clean_by }
    if (manualForm.clean_by === 'days') params.days = manualForm.days
    else params.count = manualForm.count
    if (manualForm.log_type) params.log_type = manualForm.log_type
    if (manualForm.status) params.status = manualForm.status
    await request.post({ url: '/api/admin/log-cleaner/manual', data: params })
    ElMessage.success($t('logCleanup.cleanCompleted')); fetchStats()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error($t('logCleanup.cleanFailed')) } finally { manualLoading.value = false }
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      if (formData.id) {
        await request.put({ url: `/api/admin/log-cleaner/rules/${formData.id}`, data: { ...formData }, showSuccessMessage: true })
      } else {
        await request.post({ url: '/api/admin/log-cleaner/rules', data: { ...formData }, showSuccessMessage: true })
      }
      dialogVisible.value = false; fetchRules()
    } catch { ElMessage.error($t('logCleanup.operationFailed')) } finally { submitLoading.value = false }
  })
}

onMounted(() => { fetchRules(); fetchStats() })
</script>

<style scoped lang="scss">
.log-cleanup-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.section { margin-top: 24px; &:first-child { margin-top: 0; } }
.section-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; h3 { margin: 0; font-size: 16px; font-weight: 600; } }
h3 { margin: 0 0 16px; font-size: 16px; font-weight: 600; }
.stat-card { text-align: center; padding: 8px 0; }
.stat-value { font-size: 20px; font-weight: 600; color: var(--el-color-primary); margin-bottom: 4px; }
.stat-label { color: var(--el-text-color-secondary); font-size: 14px; }
.form-tip { margin-left: 12px; color: var(--el-text-color-secondary); font-size: 13px; }
.manual-clean-card { margin-top: 8px; }
</style>
