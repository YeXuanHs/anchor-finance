<template>
  <div class="log-cleanup-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>日志清理管理</span>
        </div>
      </template>

      <!-- 清理规则列表 -->
      <div class="section">
        <div class="section-header">
          <h3>自动清理规则</h3>
          <el-button type="primary" size="small" @click="handleAddRule">添加规则</el-button>
        </div>
        <el-table :data="rules" v-loading="loading" style="width: 100%" border>
          <el-table-column prop="id" label="ID" width="70" />
          <el-table-column prop="name" label="规则名称" min-width="150" />
          <el-table-column prop="log_type" label="日志类型" width="120">
            <template #default="{ row }">
              <el-tag size="small">{{ logTypeMap[row.log_type] || row.log_type }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="retention_days" label="保留天数" width="100" align="center" />
          <el-table-column prop="max_count" label="最大条数" width="100" align="center">
            <template #default="{ row }">{{ row.max_count || '不限' }}</template>
          </el-table-column>
          <el-table-column prop="schedule" label="执行周期" width="100" />
          <el-table-column prop="last_run_at" label="上次执行" width="180" />
          <el-table-column prop="deleted_count" label="已清理条数" width="100" align="center" />
          <el-table-column prop="is_enabled" label="状态" width="80">
            <template #default="{ row }">
              <el-tag :type="row.is_enabled ? 'success' : 'info'" size="small">{{ row.is_enabled ? '启用' : '禁用' }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="220" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="handleEditRule(row)">编辑</el-button>
              <el-button type="success" link @click="handleRunNow(row)">立即执行</el-button>
              <el-popconfirm title="确定删除该规则吗？" @confirm="handleDeleteRule(row)">
                <template #reference><el-button type="danger" link>删除</el-button></template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 手动清理 -->
      <div class="section">
        <div class="section-header">
          <h3>手动清理</h3>
        </div>
        <el-card shadow="hover" class="manual-clean-card">
          <el-form :model="manualForm" label-width="100px" inline>
            <el-form-item label="日志类型">
              <el-select v-model="manualForm.log_type" placeholder="全部" clearable style="width: 150px">
                <el-option v-for="(label, key) in logTypeMap" :key="key" :label="label" :value="key" />
              </el-select>
            </el-form-item>
            <el-form-item label="日志状态">
              <el-select v-model="manualForm.status" placeholder="全部" clearable style="width: 120px">
                <el-option label="成功" value="success" />
                <el-option label="失败" value="error" />
              </el-select>
            </el-form-item>
            <el-form-item label="清理方式">
              <el-radio-group v-model="manualForm.clean_by">
                <el-radio value="days">按天数</el-radio>
                <el-radio value="count">按数量</el-radio>
              </el-radio-group>
            </el-form-item>
            <el-form-item :label="manualForm.clean_by === 'days' ? '保留天数' : '保留条数'">
              <el-input-number v-if="manualForm.clean_by === 'days'" v-model="manualForm.days" :min="1" :max="365" />
              <el-input-number v-else v-model="manualForm.count" :min="100" :step="1000" />
            </el-form-item>
            <el-form-item>
              <el-button type="danger" @click="handleManualClean" :loading="manualLoading">执行清理</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </div>

      <!-- 清理统计 -->
      <div class="section">
        <h3>清理统计</h3>
        <el-row :gutter="20">
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-value">{{ stats.total_cleaned }}</div>
              <div class="stat-label">总清理条数</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-value">{{ stats.freed_space }}</div>
              <div class="stat-label">释放空间</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-value">{{ stats.last_cleanup_at || '-' }}</div>
              <div class="stat-label">上次清理时间</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-value">{{ stats.next_cleanup_at || '-' }}</div>
              <div class="stat-label">下次清理时间</div>
            </el-card>
          </el-col>
        </el-row>
      </div>
    </el-card>

    <!-- 添加/编辑规则对话框 -->
    <el-dialog v-model="dialogVisible" :title="dialogTitle" width="500px">
      <el-form :model="formData" :rules="formRules" ref="formRef" label-width="100px">
        <el-form-item label="规则名称" prop="name">
          <el-input v-model="formData.name" placeholder="请输入规则名称" />
        </el-form-item>
        <el-form-item label="日志类型" prop="log_type">
          <el-select v-model="formData.log_type" placeholder="请选择日志类型" style="width: 100%">
            <el-option v-for="(label, key) in logTypeMap" :key="key" :label="label" :value="key" />
          </el-select>
        </el-form-item>
        <el-form-item label="保留天数" prop="retention_days">
          <el-input-number v-model="formData.retention_days" :min="1" :max="365" />
          <span class="form-tip">超过此天数的日志将被清理</span>
        </el-form-item>
        <el-form-item label="最大条数">
          <el-input-number v-model="formData.max_count" :min="0" :step="1000" />
          <span class="form-tip">0表示不限制，超过此数量将被清理</span>
        </el-form-item>
        <el-form-item label="执行周期" prop="schedule">
          <el-select v-model="formData.schedule" placeholder="请选择执行周期" style="width: 100%">
            <el-option label="每天" value="daily" />
            <el-option label="每周" value="weekly" />
            <el-option label="每月" value="monthly" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="formData.is_enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const logTypeMap: Record<string, string> = {
  operation: '操作日志', login: '登录日志', system: '系统日志', api: 'API日志',
  error: '错误日志', ticket: '工单日志', order: '订单日志', all: '全部'
}

const loading = ref(false)
const submitLoading = ref(false)
const manualLoading = ref(false)
const dialogVisible = ref(false)
const dialogTitle = ref('添加规则')
const formRef = ref<FormInstance>()

const rules = ref<any[]>([])
const stats = reactive({ total_cleaned: 0, freed_space: '0 MB', last_cleanup_at: '', next_cleanup_at: '' })

const formData = reactive({
  id: undefined as number | undefined,
  name: '', log_type: 'operation', retention_days: 30, max_count: 0, schedule: 'daily', is_enabled: true
})

const manualForm = reactive({ log_type: '', status: '', clean_by: 'days', days: 30, count: 1000 })

const formRules: FormRules = {
  name: [{ required: true, message: '请输入规则名称', trigger: 'blur' }],
  log_type: [{ required: true, message: '请选择日志类型', trigger: 'change' }],
  retention_days: [{ required: true, message: '请输入保留天数', trigger: 'blur' }],
  schedule: [{ required: true, message: '请选择执行周期', trigger: 'change' }]
}

const fetchRules = async () => {
  loading.value = true
  try { rules.value = await request.get({ url: '/api/admin/log-cleaner/rules' }) || [] }
  catch { ElMessage.error('获取清理规则失败') } finally { loading.value = false }
}

const fetchStats = async () => {
  try { const data = await request.get({ url: '/api/admin/log-cleaner/stats' }); if (data) Object.assign(stats, data) }
  catch { /* ignore */ }
}

const handleAddRule = () => {
  dialogTitle.value = '添加规则'
  Object.assign(formData, { id: undefined, name: '', log_type: 'operation', retention_days: 30, max_count: 0, schedule: 'daily', is_enabled: true })
  dialogVisible.value = true
}

const handleEditRule = (row: any) => {
  dialogTitle.value = '编辑规则'
  Object.assign(formData, { id: row.id, name: row.name, log_type: row.log_type, retention_days: row.retention_days, max_count: row.max_count || 0, schedule: row.schedule, is_enabled: row.is_enabled })
  dialogVisible.value = true
}

const handleDeleteRule = async (row: any) => {
  try { await request.del({ url: `/api/admin/log-cleaner/rules/${row.id}` }); ElMessage.success('删除成功'); fetchRules() }
  catch { ElMessage.error('删除失败') }
}

const handleRunNow = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定立即执行清理规则 "${row.name}" 吗？`, '确认执行')
    await request.post({ url: `/api/admin/log-cleaner/rules/${row.id}/run` })
    ElMessage.success('清理任务已提交'); fetchRules(); fetchStats()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error('执行失败') }
}

const handleManualClean = async () => {
  try {
    const msg = manualForm.clean_by === 'days'
      ? `确定清理${manualForm.days}天前${manualForm.log_type ? '的' + logTypeMap[manualForm.log_type] : '所有'}日志吗？`
      : `确定只保留最近${manualForm.count}条${manualForm.log_type ? logTypeMap[manualForm.log_type] : '所有'}日志吗？`
    await ElMessageBox.confirm(msg, '手动清理确认')
    manualLoading.value = true
    const params: any = { clean_by: manualForm.clean_by }
    if (manualForm.clean_by === 'days') params.days = manualForm.days
    else params.count = manualForm.count
    if (manualForm.log_type) params.log_type = manualForm.log_type
    if (manualForm.status) params.status = manualForm.status
    await request.post({ url: '/api/admin/log-cleaner/manual', data: params })
    ElMessage.success('清理完成'); fetchStats()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error('清理失败') } finally { manualLoading.value = false }
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
    } catch { ElMessage.error('操作失败') } finally { submitLoading.value = false }
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
