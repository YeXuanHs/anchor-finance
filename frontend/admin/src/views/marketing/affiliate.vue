<template>
  <div class="affiliate-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="用户ID/推介码" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="正常" value="active" />
            <el-option label="禁用" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>推广联盟管理</h3>
        <div>
          <el-button @click="showConfigDialog = true">
            <el-icon><Setting /></el-icon>
            佣金配置
          </el-button>
          <el-button @click="showWithdrawDialog = true">
            <el-icon><Wallet /></el-icon>
            提现审核
          </el-button>
        </div>
      </div>

      <el-table :data="tableData" style="width: 100%" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="affiliate_code" label="推介码" width="140">
          <template #default="{ row }">
            <el-tag type="info" size="small">{{ row.affiliate_code }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="referral_count" label="下线数" width="100" />
        <el-table-column prop="commission_rate" label="佣金比例" width="120">
          <template #default="{ row }">{{ (row.commission_rate * 100).toFixed(1) }}%</template>
        </el-table-column>
        <el-table-column prop="total_commission" label="累计佣金" width="140">
          <template #default="{ row }">
            <span class="amount">¥{{ row.total_commission?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="pending_commission" label="待结算" width="120">
          <template #default="{ row }">
            <span style="color: var(--el-color-warning);">¥{{ row.pending_commission?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'active' ? 'success' : 'info'" size="small">
              {{ row.status === 'active' ? '正常' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="copyLink(row)">推广链接</el-button>
            <el-button type="primary" link @click="openEditDialog(row)">编辑</el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchData"
          @size-change="fetchData"
        />
      </div>
    </div>

    <el-dialog v-model="editDialogVisible" title="编辑推广员" width="500px" destroy-on-close>
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <el-form-item label="佣金比例" prop="commission_rate">
          <el-input-number v-model="formData.commission_rate" :min="0" :max="1" :step="0.01" :precision="2" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">0~1，如 0.1 表示 10%</span>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-switch v-model="formData.status" active-value="active" inactive-value="disabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="editDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showConfigDialog" title="佣金配置" width="550px">
      <el-form :model="configForm" label-width="120px">
        <el-form-item label="默认佣金比例">
          <el-input-number v-model="configForm.default_rate" :min="0" :max="1" :step="0.01" :precision="2" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">0~1</span>
        </el-form-item>
        <el-form-item label="最低提现金额">
          <el-input-number v-model="configForm.min_payout" :min="0" :precision="2" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">元</span>
        </el-form-item>
        <el-form-item label="结算周期">
          <el-select v-model="configForm.settle_cycle">
            <el-option label="实时结算" value="realtime" />
            <el-option label="按月结算" value="monthly" />
            <el-option label="按季度结算" value="quarterly" />
          </el-select>
        </el-form-item>
        <el-form-item label="Cookie有效期">
          <el-input-number v-model="configForm.cookie_days" :min="1" />
          <span style="margin-left: 8px; color: var(--el-text-color-secondary);">天</span>
        </el-form-item>
        <el-form-item label="启用推广计划">
          <el-switch v-model="configForm.enabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showConfigDialog = false">取消</el-button>
        <el-button type="primary" @click="saveConfig">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showWithdrawDialog" title="提现审核" width="700px">
      <el-table :data="withdrawList" style="width: 100%" v-loading="withdrawLoading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户" width="120" />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">¥{{ row.amount?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="method" label="提现方式" width="120" />
        <el-table-column prop="account" label="提现账号" min-width="160" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'pending' ? 'warning' : row.status === 'approved' ? 'success' : 'danger'" size="small">
              {{ { pending: '待审核', approved: '已通过', rejected: '已拒绝' }[row.status as string] || row.status }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <template v-if="row.status === 'pending'">
              <el-button type="success" link @click="reviewWithdraw(row, 'approved')">通过</el-button>
              <el-button type="danger" link @click="reviewWithdraw(row, 'rejected')">拒绝</el-button>
            </template>
            <span v-else style="color: var(--el-text-color-secondary);">已处理</span>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Plus, Setting, Wallet } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const submitLoading = ref(false)
const tableData = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const editDialogVisible = ref(false)
const editId = ref<number | null>(null)
const formRef = ref<FormInstance>()
const showConfigDialog = ref(false)
const showWithdrawDialog = ref(false)
const withdrawLoading = ref(false)
const withdrawList = ref<any[]>([])

const searchForm = ref({ keyword: '', status: '' })

const formData = reactive({
  commission_rate: 0.1,
  status: 'active' as string
})

const formRules: FormRules = {
  commission_rate: [{ required: true, message: '请输入佣金比例', trigger: 'blur' }]
}

const configForm = ref({
  default_rate: 0.1,
  min_payout: 100,
  settle_cycle: 'monthly',
  cookie_days: 30,
  enabled: true
})

const handleSearch = () => { page.value = 1; fetchData() }

const resetSearch = () => {
  searchForm.value = { keyword: '', status: '' }
  handleSearch()
}

const fetchData = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/affiliates', {
      params: { page: page.value, page_size: pageSize.value, ...searchForm.value }
    })
    tableData.value = data.data || []
    total.value = data.total || 0
  } catch {} finally {
    loading.value = false
  }
}

const fetchWithdrawList = async () => {
  withdrawLoading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/affiliates/withdrawals', {
      params: { status: 'pending' }
    })
    withdrawList.value = data.data || []
  } catch {} finally {
    withdrawLoading.value = false
  }
}

const openEditDialog = (row: any) => {
  editId.value = row.id
  formData.commission_rate = row.commission_rate
  formData.status = row.status
  editDialogVisible.value = true
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate()
  submitLoading.value = true
  try {
    await request.put(`/admin/api/v1/affiliates/${editId.value}`, formData)
    ElMessage.success('更新成功')
    editDialogVisible.value = false
    fetchData()
  } catch {} finally {
    submitLoading.value = false
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm(`确定删除推广员「${row.username}」吗？`, '提示', { type: 'warning' })
  try {
    await request.delete(`/admin/api/v1/affiliates/${row.id}`)
    ElMessage.success('删除成功')
    fetchData()
  } catch {}
}

const copyLink = (row: any) => {
  const link = `${window.location.origin}/ref/${row.affiliate_code}`
  navigator.clipboard.writeText(link).then(() => {
    ElMessage.success('推广链接已复制到剪贴板')
  })
}

const saveConfig = async () => {
  try {
    await request.put('/admin/api/v1/affiliates/config', configForm.value)
    ElMessage.success('配置已保存')
    showConfigDialog.value = false
  } catch {}
}

const reviewWithdraw = async (row: any, status: string) => {
  const action = status === 'approved' ? '通过' : '拒绝'
  await ElMessageBox.confirm(`确定${action}该笔提现申请？`, '提示', { type: 'warning' })
  try {
    await request.put(`/admin/api/v1/affiliates/withdrawals/${row.id}`, { status })
    ElMessage.success(`已${action}`)
    fetchWithdrawList()
  } catch {}
}

onMounted(() => {
  fetchData()
  fetchWithdrawList()
})
</script>

<style scoped lang="scss">
.affiliate-page {
  .search-bar { margin-bottom: 16px; }
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 18px; font-weight: 600; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
  .amount { color: var(--el-color-danger); font-weight: 600; }
}
</style>
