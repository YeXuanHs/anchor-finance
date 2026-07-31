<template>
  <div class="balance-page art-full-height">
    <!-- 搜索栏 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="用户ID">
          <el-input v-model="searchForm.user_id" placeholder="请输入用户ID" clearable />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input v-model="searchForm.username" placeholder="请输入用户名" clearable />
        </el-form-item>
        <el-form-item label="变动类型">
          <el-select v-model="searchForm.type" placeholder="全部" clearable>
            <el-option label="充值" value="recharge" />
            <el-option label="扣款" value="deduct" />
            <el-option label="消费" value="payment" />
            <el-option label="退款" value="refund" />
          </el-select>
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
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
          <span>余额记录列表</span>
          <el-space>
            <el-button type="primary" @click="showRechargeDialog">
              <el-icon><Plus /></el-icon>
              充值
            </el-button>
            <el-button type="warning" @click="showDeductDialog">
              <el-icon><Minus /></el-icon>
              扣款
            </el-button>
          </el-space>
        </div>
      </template>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="type" label="变动类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">
              {{ getTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            <span :class="row.type === 'recharge' || row.type === 'refund' ? 'text-green' : 'text-red'">
              {{ row.type === 'recharge' || row.type === 'refund' ? '+' : '-' }}¥{{ row.amount?.toFixed(2) || '0.00' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="balance_before" label="变动前余额" width="120">
          <template #default="{ row }">
            ¥{{ row.balance_before?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="balance_after" label="变动后余额" width="120">
          <template #default="{ row }">
            ¥{{ row.balance_after?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column prop="operator" label="操作人" width="100" />
        <el-table-column prop="created_at" label="创建时间" width="180" />
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

    <!-- 充值对话框 -->
    <el-dialog v-model="rechargeDialogVisible" title="用户充值" width="500px">
      <el-form :model="rechargeForm" :rules="rechargeRules" ref="rechargeFormRef" label-width="100px">
        <el-form-item label="用户ID" prop="user_id">
          <el-input v-model="rechargeForm.user_id" placeholder="请输入用户ID" />
        </el-form-item>
        <el-form-item label="充值金额" prop="amount">
          <el-input-number v-model="rechargeForm.amount" :min="0.01" :precision="2" :step="100" style="width: 100%" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="rechargeForm.remark" type="textarea" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rechargeDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleRecharge" :loading="submitLoading">确定</el-button>
      </template>
    </el-dialog>

    <!-- 扣款对话框 -->
    <el-dialog v-model="deductDialogVisible" title="用户扣款" width="500px">
      <el-form :model="deductForm" :rules="deductRules" ref="deductFormRef" label-width="100px">
        <el-form-item label="用户ID" prop="user_id">
          <el-input v-model="deductForm.user_id" placeholder="请输入用户ID" />
        </el-form-item>
        <el-form-item label="扣款金额" prop="amount">
          <el-input-number v-model="deductForm.amount" :min="0.01" :precision="2" :step="100" style="width: 100%" />
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="deductForm.remark" type="textarea" placeholder="请输入备注" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="deductDialogVisible = false">取消</el-button>
        <el-button type="warning" @click="handleDeduct" :loading="submitLoading">确定扣款</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search, Plus, Minus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'BalanceManage' })

// 加载状态
const loading = ref(false)
const submitLoading = ref(false)

// 搜索表单
const searchForm = reactive({
  user_id: '',
  username: '',
  type: undefined as string | undefined,
  date_range: [] as string[]
})

// 分页
const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

// 表格数据
const tableData = ref([])

// 充值对话框
const rechargeDialogVisible = ref(false)
const rechargeFormRef = ref<FormInstance>()
const rechargeForm = reactive({
  user_id: '',
  amount: 0,
  remark: ''
})

// 扣款对话框
const deductDialogVisible = ref(false)
const deductFormRef = ref<FormInstance>()
const deductForm = reactive({
  user_id: '',
  amount: 0,
  remark: ''
})

// 充值表单验证规则
const rechargeRules: FormRules = {
  user_id: [
    { required: true, message: '请输入用户ID', trigger: 'blur' }
  ],
  amount: [
    { required: true, message: '请输入充值金额', trigger: 'blur' }
  ]
}

// 扣款表单验证规则
const deductRules: FormRules = {
  user_id: [
    { required: true, message: '请输入用户ID', trigger: 'blur' }
  ],
  amount: [
    { required: true, message: '请输入扣款金额', trigger: 'blur' }
  ]
}

// 获取类型标签
const getTypeTag = (type: string) => {
  const map: Record<string, string> = {
    recharge: 'success',
    deduct: 'danger',
    payment: 'warning',
    refund: 'info'
  }
  return map[type] || 'info'
}

// 获取类型文本
const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    recharge: '充值',
    deduct: '扣款',
    payment: '消费',
    refund: '退款'
  }
  return map[type] || '未知'
}

// 获取余额记录列表
const fetchBalances = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size
    }
    if (searchForm.user_id) params.user_id = searchForm.user_id
    if (searchForm.username) params.username = searchForm.username
    if (searchForm.type) params.type = searchForm.type
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }

    const data = await request.get({
      url: '/api/admin/balances',
      params
    })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    console.error('获取余额记录失败:', error)
    ElMessage.error('获取余额记录失败')
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => {
  pagination.page = 1
  fetchBalances()
}

// 重置
const handleReset = () => {
  searchForm.user_id = ''
  searchForm.username = ''
  searchForm.type = undefined
  searchForm.date_range = []
  handleSearch()
}

// 显示充值对话框
const showRechargeDialog = () => {
  rechargeForm.user_id = ''
  rechargeForm.amount = 0
  rechargeForm.remark = ''
  rechargeDialogVisible.value = true
}

// 显示扣款对话框
const showDeductDialog = () => {
  deductForm.user_id = ''
  deductForm.amount = 0
  deductForm.remark = ''
  deductDialogVisible.value = true
}

// 处理充值
const handleRecharge = async () => {
  if (!rechargeFormRef.value) return

  await rechargeFormRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      await request.post({
        url: '/api/admin/balances/recharge',
        params: rechargeForm,
        showSuccessMessage: true
      })
      rechargeDialogVisible.value = false
      fetchBalances()
    } catch (error) {
      console.error('充值失败:', error)
    } finally {
      submitLoading.value = false
    }
  })
}

// 处理扣款
const handleDeduct = async () => {
  if (!deductFormRef.value) return

  await deductFormRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      await request.post({
        url: '/api/admin/balances/deduct',
        params: deductForm,
        showSuccessMessage: true
      })
      deductDialogVisible.value = false
      fetchBalances()
    } catch (error) {
      console.error('扣款失败:', error)
    } finally {
      submitLoading.value = false
    }
  })
}

// 分页大小变化
const handleSizeChange = () => {
  pagination.page = 1
  fetchBalances()
}

// 页码变化
const handlePageChange = () => {
  fetchBalances()
}

onMounted(() => {
  fetchBalances()
})
</script>

<style scoped lang="scss">
.balance-page {
  padding: 20px;
}

.search-card {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.search-form {
  .el-form-item {
    margin-bottom: 0;
  }
}

.pagination-container {
  display: flex;
  justify-content: flex-end;
  margin-top: 20px;
}

.text-green {
  color: #67c23a;
  font-weight: 600;
}

.text-red {
  color: #f56c6c;
  font-weight: 600;
}
</style>
