<template>
  <div class="credit-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="用户">
          <el-input v-model="searchForm.username" placeholder="用户名" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="正常" value="active" />
            <el-option label="冻结" value="frozen" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div class="art-card">
      <div class="table-header">
        <h3>信用额度管理</h3>
        <el-button type="primary" @click="showDialog = true; isAdjusting = true; resetForm()">
          <el-icon><Plus /></el-icon>调整信用额度
        </el-button>
      </div>
      <el-table :data="list" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="username" label="用户名" width="140" />
        <el-table-column label="信用额度" width="130">
          <template #default="{ row }">
            <span class="amount">¥{{ row.credit_limit?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="已用额度" width="130">
          <template #default="{ row }">¥{{ row.used_credit?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="可用额度" width="130">
          <template #default="{ row }">
            <span class="amount success">¥{{ row.available_credit?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">{{ row.status === 'active' ? '正常' : '冻结' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="updated_at" label="更新时间" width="180" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="adjustCredit(row)">调整额度</el-button>
            <el-button link type="primary" @click="viewLogs(row)">信用日志</el-button>
            <el-button :type="row.status === 'active' ? 'danger' : 'success'" link @click="toggleStatus(row)">
              {{ row.status === 'active' ? '冻结' : '解冻' }}
            </el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @change="fetchData"
        style="margin-top: 16px; justify-content: flex-end"
      />
    </div>

    <el-dialog v-model="showDialog" :title="isAdjusting ? '调整信用额度' : '新增信用额度'" width="500px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="用户">
          <el-input v-model="form.username" placeholder="请输入用户名" :disabled="!!form.user_id" />
        </el-form-item>
        <el-form-item label="信用额度">
          <el-input-number v-model="form.credit_limit" :min="0" :precision="2" style="width: 200px" />
        </el-form-item>
        <el-form-item label="调整原因">
          <el-select v-model="form.adjust_reason" placeholder="请选择" clearable>
            <el-option label="新用户授信" value="new" />
            <el-option label="提升额度" value="increase" />
            <el-option label="降低额度" value="decrease" />
            <el-option label="风险控制" value="risk" />
            <el-option label="其他" value="other" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注"><el-input v-model="form.remark" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDialog = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-drawer v-model="showLogs" title="信用额度日志" size="650px">
      <el-table :data="creditLogs" v-loading="logsLoading">
        <el-table-column prop="created_at" label="时间" width="180" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ logTypeMap[row.type] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="变动金额" width="130">
          <template #default="{ row }">
            <span :style="{ color: row.amount > 0 ? '#67c23a' : '#f56c6c', fontWeight: 600 }">
              {{ row.amount > 0 ? '+' : '' }}¥{{ row.amount?.toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column label="变动后额度" width="130">
          <template #default="{ row }">¥{{ row.credit_after?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="operator" label="操作人" width="100" />
        <el-table-column prop="remark" label="备注" min-width="120" show-overflow-tooltip />
      </el-table>
      <el-pagination
        v-model:current-page="logsPage"
        v-model:page-size="logsPageSize"
        :total="logsTotal"
        layout="total, prev, pager, next"
        @change="fetchLogs"
        style="margin-top: 16px; justify-content: flex-end"
      />
    </el-drawer>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const logTypeMap: Record<string, string> = { new: '新增授信', increase: '提升额度', decrease: '降低额度', freeze: '冻结', unfreeze: '解冻', consume: '消费', repay: '还款' }

const loading = ref(false)
const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const showDialog = ref(false)
const showLogs = ref(false)
const isAdjusting = ref(false)
const logsLoading = ref(false)
const creditLogs = ref<any[]>([])
const logsTotal = ref(0)
const logsPage = ref(1)
const logsPageSize = ref(20)
const currentUserId = ref<number>(0)
const searchForm = ref({ username: '', status: '' })
const form = ref<any>({ user_id: 0, username: '', credit_limit: 0, adjust_reason: '', remark: '' })

const resetForm = () => { form.value = { user_id: 0, username: '', credit_limit: 0, adjust_reason: '', remark: '' } }
const resetSearch = () => { searchForm.value = { username: '', status: '' }; fetchData() }

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/credits', { params: { page: page.value, page_size: pageSize.value, ...searchForm.value } })
    list.value = data.data || []; total.value = data.total || 0
  } catch {} finally { loading.value = false }
}

const adjustCredit = (row: any) => {
  form.value = { user_id: row.user_id, username: row.username, credit_limit: row.credit_limit, adjust_reason: '', remark: '' }
  isAdjusting.value = true; showDialog.value = true
}

const handleSubmit = async () => {
  try {
    await request.post('/admin/api/v1/credits/adjust', form.value)
    ElMessage.success('信用额度调整成功'); showDialog.value = false; fetchData()
  } catch { ElMessage.error('操作失败') }
}

const toggleStatus = async (row: any) => {
  const action = row.status === 'active' ? '冻结' : '解冻'
  await ElMessageBox.confirm(`确定${action}该用户信用额度？`, '确认')
  try {
    await request.put(`/admin/api/v1/credits/${row.id}/status`, { status: row.status === 'active' ? 'frozen' : 'active' })
    ElMessage.success(`${action}成功`); fetchData()
  } catch { ElMessage.error('操作失败') }
}

const viewLogs = async (row: any) => {
  currentUserId.value = row.user_id || row.id
  logsPage.value = 1
  showLogs.value = true
  fetchLogs()
}

const fetchLogs = async () => {
  logsLoading.value = true
  try {
    const { data } = await request.get(`/admin/api/v1/credits/${currentUserId.value}/logs`, { params: { page: logsPage.value, page_size: logsPageSize.value } })
    creditLogs.value = data.data || []; logsTotal.value = data.total || 0
  } catch {} finally { logsLoading.value = false }
}

onMounted(fetchData)
</script>
<style scoped lang="scss">
.search-bar { margin-bottom: 16px; }
.table-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; h3 { font-size: 18px; font-weight: 600; margin: 0; } }
.amount { color: var(--danger-color); font-weight: 600; &.success { color: var(--success-color); } }
</style>
