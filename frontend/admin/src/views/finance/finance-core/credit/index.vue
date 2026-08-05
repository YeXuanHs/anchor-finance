<template>
  <div class="credit-page art-full-height">
    <!-- 搜索栏 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="用户ID">
          <el-input v-model="searchForm.user_id" placeholder="请输入用户ID" clearable />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="searchForm.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="正常" :value="1" />
            <el-option label="禁用" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never" class="art-table-card">
      <!-- 表格头部 -->
      <template #header>
        <div class="card-header">
          <span>信用额度列表</span>
          <div>
            <el-button type="primary" size="small" @click="handleConfig">
              <el-icon><Setting /></el-icon>
              额度配置
            </el-button>
          </div>
        </div>
      </template>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="credit_limit" label="信用额度" width="150">
          <template #default="{ row }">
            <span class="text-primary">¥{{ row.credit_limit?.toFixed(2) || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="used_amount" label="已用额度" width="150">
          <template #default="{ row }">
            <span class="text-warning">¥{{ row.used_amount?.toFixed(2) || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="available_amount" label="可用额度" width="150">
          <template #default="{ row }">
            <span class="text-success">¥{{ row.available_amount?.toFixed(2) || '0.00' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="使用率" width="150">
          <template #default="{ row }">
            <el-progress :percentage="getUsagePercent(row)" :color="getUsageColor(row)" :stroke-width="10" />
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'danger'" size="small">
              {{ row.status === 1 ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="expire_at" label="到期时间" width="180" />
        <el-table-column prop="updated_at" label="更新时间" width="180" />
        <el-table-column label="操作" width="280" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleSetCredit(row)">设置额度</el-button>
            <el-button type="warning" link @click="handleAdjustCredit(row)">调整额度</el-button>
            <el-button type="info" link @click="handleViewLog(row)">额度日志</el-button>
            <el-popconfirm title="确定禁用该用户信用额度吗？" @confirm="handleToggleStatus(row)">
              <template #reference>
                <el-button :type="row.status === 1 ? 'danger' : 'success'" link>
                  {{ row.status === 1 ? '禁用' : '启用' }}
                </el-button>
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

    <!-- 设置额度对话框 -->
    <el-dialog v-model="setDialogVisible" title="设置信用额度" width="500px">
      <el-form :model="setForm" :rules="setRules" ref="setFormRef" label-width="100px">
        <el-form-item label="用户">
          <el-input :value="setForm.username" disabled />
        </el-form-item>
        <el-form-item label="当前额度">
          <el-input :value="'¥' + (setForm.current_limit?.toFixed(2) || '0.00')" disabled />
        </el-form-item>
        <el-form-item label="新额度" prop="credit_limit">
          <el-input-number v-model="setForm.credit_limit" :min="0" :precision="2" :step="1000" style="width: 100%" />
        </el-form-item>
        <el-form-item label="到期时间" prop="expire_at">
          <el-date-picker v-model="setForm.expire_at" type="datetime" placeholder="选择到期时间" value-format="YYYY-MM-DD HH:mm:ss" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="setDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitSet" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 调整额度对话框 -->
    <el-dialog v-model="adjustDialogVisible" title="调整信用额度" width="500px">
      <el-form :model="adjustForm" :rules="adjustRules" ref="adjustFormRef" label-width="100px">
        <el-form-item label="用户">
          <el-input :value="adjustForm.username" disabled />
        </el-form-item>
        <el-form-item label="当前额度">
          <el-input :value="'¥' + (adjustForm.current_limit?.toFixed(2) || '0.00')" disabled />
        </el-form-item>
        <el-form-item label="调整方式">
          <el-radio-group v-model="adjustForm.adjust_type">
            <el-radio value="increase">增加</el-radio>
            <el-radio value="decrease">减少</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="调整金额" prop="amount">
          <el-input-number v-model="adjustForm.amount" :min="0.01" :precision="2" :step="100" style="width: 100%" />
        </el-form-item>
        <el-form-item label="调整原因" prop="reason">
          <el-input v-model="adjustForm.reason" type="textarea" placeholder="请输入调整原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="adjustDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmitAdjust" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 额度日志对话框 -->
    <el-dialog v-model="logDialogVisible" title="信用额度日志" width="800px">
      <div class="log-header" v-if="logUserInfo">
        <span>用户：<strong>{{ logUserInfo.username }}</strong></span>
        <span style="margin-left: 20px">当前余额：<strong class="text-primary">¥{{ logUserInfo.credit?.toFixed(2) || '0.00' }}</strong></span>
      </div>
      <el-table :data="logData" v-loading="logLoading" style="width: 100%" border size="small">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="description" label="描述" min-width="250" show-overflow-tooltip />
        <el-table-column prop="amount" label="金额" width="120" align="right">
          <template #default="{ row }">
            <span :class="row.amount >= 0 ? 'text-success' : 'text-danger'">
              {{ row.amount >= 0 ? '+' : '' }}{{ row.amount?.toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="balance" label="余额" width="120" align="right">
          <template #default="{ row }">
            ¥{{ row.balance?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="180" />
      </el-table>
      <div class="pagination-container">
        <el-pagination
          v-model:current-page="logPagination.page"
          v-model:page-size="logPagination.page_size"
          :total="logPagination.total"
          layout="total, prev, pager, next"
          small
          @current-change="fetchCreditLog"
        />
      </div>
    </el-dialog>

    <!-- 额度配置对话框 -->
    <el-dialog v-model="configDialogVisible" title="信用额度配置" width="500px">
      <el-form :model="configForm" label-width="120px">
        <el-form-item label="默认额度">
          <el-input-number v-model="configForm.default_credit" :min="0" :precision="2" :step="100" style="width: 100%" />
          <div class="form-tip">新用户默认信用额度，0表示不启用</div>
        </el-form-item>
        <el-form-item label="最大额度">
          <el-input-number v-model="configForm.max_credit" :min="0" :precision="2" :step="10000" style="width: 100%" />
          <div class="form-tip">单用户最大信用额度上限</div>
        </el-form-item>
        <el-form-item label="自动续期">
          <el-switch v-model="configForm.auto_renew" />
          <div class="form-tip">到期后自动续期</div>
        </el-form-item>
        <el-form-item label="续期天数" v-if="configForm.auto_renew">
          <el-input-number v-model="configForm.renew_days" :min="1" :max="365" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="configDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSaveConfig" :loading="configLoading">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search, Setting } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'CreditManage' })

const loading = ref(false)
const submitLoading = ref(false)
const configLoading = ref(false)
const logLoading = ref(false)

const searchForm = reactive({ user_id: '', username: '', status: undefined as number | undefined })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

// 设置额度
const setDialogVisible = ref(false)
const setFormRef = ref<FormInstance>()
const setForm = reactive({ id: 0, username: '', current_limit: 0, credit_limit: 0, expire_at: '' })

// 调整额度
const adjustDialogVisible = ref(false)
const adjustFormRef = ref<FormInstance>()
const adjustForm = reactive({ id: 0, username: '', current_limit: 0, adjust_type: 'increase', amount: 0, reason: '' })

// 额度日志
const logDialogVisible = ref(false)
const logData = ref<any[]>([])
const logUserInfo = ref<any>(null)
const logPagination = reactive({ page: 1, page_size: 20, total: 0 })
const logUserId = ref(0)

// 配置
const configDialogVisible = ref(false)
const configForm = reactive({ default_credit: 0, max_credit: 0, auto_renew: false, renew_days: 30 })

const setRules: FormRules = { credit_limit: [{ required: true, message: '请输入信用额度', trigger: 'blur' }] }
const adjustRules: FormRules = { amount: [{ required: true, message: '请输入调整金额', trigger: 'blur' }] }

const getUsagePercent = (row: any) => {
  if (!row.credit_limit || row.credit_limit === 0) return 0
  return Math.min(100, Math.round((row.used_amount / row.credit_limit) * 100))
}

const getUsageColor = (row: any) => {
  const percent = getUsagePercent(row)
  if (percent >= 90) return '#f56c6c'
  if (percent >= 70) return '#e6a23c'
  return '#67c23a'
}

const fetchCredits = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (searchForm.user_id) params.user_id = searchForm.user_id
    if (searchForm.username) params.username = searchForm.username
    if (searchForm.status !== undefined) params.status = searchForm.status
    const data = await request.get({ url: '/api/admin/credit/index', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error('获取信用额度列表失败')
  } finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchCredits() }
const handleReset = () => { searchForm.user_id = ''; searchForm.username = ''; searchForm.status = undefined; handleSearch() }

const handleSetCredit = (row: any) => {
  setForm.id = row.id; setForm.username = row.username
  setForm.current_limit = row.credit_limit; setForm.credit_limit = row.credit_limit
  setForm.expire_at = row.expire_at || ''
  setDialogVisible.value = true
}

const handleAdjustCredit = (row: any) => {
  adjustForm.id = row.id; adjustForm.username = row.username
  adjustForm.current_limit = row.credit_limit; adjustForm.adjust_type = 'increase'
  adjustForm.amount = 0; adjustForm.reason = ''
  adjustDialogVisible.value = true
}

const handleViewLog = async (row: any) => {
  logUserId.value = row.user_id
  logUserInfo.value = row
  logPagination.page = 1
  logDialogVisible.value = true
  fetchCreditLog()
}

const fetchCreditLog = async () => {
  logLoading.value = true
  try {
    const data = await request.get({ url: `/api/admin/credit/${logUserId.value}/log`, params: { page: logPagination.page, page_size: logPagination.page_size } })
    logData.value = data.list || data.data || []
    logPagination.total = data.total || data.count || 0
    if (data.user) logUserInfo.value = { ...logUserInfo.value, ...data.user }
  } catch { ElMessage.error('获取额度日志失败') } finally { logLoading.value = false }
}

const handleSubmitSet = async () => {
  if (!setFormRef.value) return
  await setFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      await request.put({ url: `/api/admin/credit/users/${setForm.id}/settings`, params: { credit_limit: setForm.credit_limit, expire_at: setForm.expire_at }, showSuccessMessage: true })
      setDialogVisible.value = false; fetchCredits()
    } finally { submitLoading.value = false }
  })
}

const handleSubmitAdjust = async () => {
  if (!adjustFormRef.value) return
  await adjustFormRef.value.validate(async (valid) => {
    if (!valid) return
    submitLoading.value = true
    try {
      await request.post({ url: `/api/admin/credit/users/${adjustForm.id}/adjust`, params: { amount: adjustForm.amount, adjust_type: adjustForm.adjust_type, reason: adjustForm.reason }, showSuccessMessage: true })
      adjustDialogVisible.value = false; fetchCredits()
    } finally { submitLoading.value = false }
  })
}

const handleToggleStatus = async (row: any) => {
  try {
    await request.put({ url: `/api/admin/credit/users/${row.id}/settings`, params: { status: row.status === 1 ? 0 : 1 }, showSuccessMessage: true })
    fetchCredits()
  } catch { ElMessage.error('操作失败') }
}

const handleConfig = async () => {
  try {
    const data = await request.get({ url: '/api/admin/credit/config' })
    if (data) Object.assign(configForm, data)
  } catch { /* use defaults */ }
  configDialogVisible.value = true
}

const handleSaveConfig = async () => {
  configLoading.value = true
  try {
    await request.post({ url: '/api/admin/credit/config', data: configForm, showSuccessMessage: true })
    configDialogVisible.value = false
  } catch { ElMessage.error('保存失败') } finally { configLoading.value = false }
}

const handleSizeChange = () => { pagination.page = 1; fetchCredits() }
const handlePageChange = () => { fetchCredits() }

onMounted(() => { fetchCredits() })
</script>

<style scoped lang="scss">
.credit-page { padding: 20px; }
.search-card { margin-bottom: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { .el-form-item { margin-bottom: 0; } }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
.text-primary { color: #409eff; font-weight: 600; }
.text-warning { color: #e6a23c; font-weight: 600; }
.text-success { color: #67c23a; font-weight: 600; }
.text-danger { color: #f56c6c; font-weight: 600; }
.log-header { margin-bottom: 16px; padding: 12px; background: #f5f7fa; border-radius: 4px; }
.form-tip { font-size: 12px; color: #909399; margin-top: 4px; }
</style>
