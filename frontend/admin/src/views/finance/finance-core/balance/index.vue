<template>
  <div class="balance-page art-full-height">
    <!-- 搜索栏 -->
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('balance.userId')">
          <el-input v-model="searchForm.user_id" :placeholder="$t('balance.userIdPlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('common.username')">
          <el-input v-model="searchForm.username" :placeholder="$t('balance.usernamePlaceholder')" clearable />
        </el-form-item>
        <el-form-item :label="$t('balance.changeType')">
          <el-select v-model="searchForm.type" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('balance.recharge')" value="recharge" />
            <el-option :label="$t('balance.deduct')" value="deduct" />
            <el-option :label="$t('balance.payment')" value="payment" />
            <el-option :label="$t('balance.refund')" value="refund" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('balance.dateRange')">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            :range-separator="$t('common.to')"
            :start-placeholder="$t('common.startDate')"
            :end-placeholder="$t('common.endDate')"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            {{ $t('common.search') }}
          </el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never" class="art-table-card">
      <!-- 表格头部 -->
      <template #header>
        <div class="card-header">
          <span>{{ $t('balance.recordList') }}</span>
          <el-space>
            <el-button type="primary" @click="showRechargeDialog">
              <el-icon><Plus /></el-icon>
              {{ $t('balance.recharge') }}
            </el-button>
            <el-button type="warning" @click="showDeductDialog">
              <el-icon><Minus /></el-icon>
              {{ $t('balance.deduct') }}
            </el-button>
          </el-space>
        </div>
      </template>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" :label="$t('balance.userId')" width="100" />
        <el-table-column prop="username" :label="$t('common.username')" width="120" />
        <el-table-column prop="type" :label="$t('balance.changeType')" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">
              {{ getTypeText(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" :label="$t('common.amount')" width="120">
          <template #default="{ row }">
            <span :class="row.type === 'recharge' || row.type === 'refund' ? 'text-green' : 'text-red'">
              {{ row.type === 'recharge' || row.type === 'refund' ? '+' : '-' }}¥{{ row.amount?.toFixed(2) || '0.00' }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="balance_before" :label="$t('balance.balanceBefore')" width="120">
          <template #default="{ row }">
            ¥{{ row.balance_before?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="balance_after" :label="$t('balance.balanceAfter')" width="120">
          <template #default="{ row }">
            ¥{{ row.balance_after?.toFixed(2) || '0.00' }}
          </template>
        </el-table-column>
        <el-table-column prop="remark" :label="$t('common.remark')" min-width="150" show-overflow-tooltip />
        <el-table-column prop="operator" :label="$t('common.operator')" width="100" />
        <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180" />
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
    <el-dialog v-model="rechargeDialogVisible" :title="$t('balance.userRecharge')" width="500px">
      <el-form :model="rechargeForm" :rules="rechargeRules" ref="rechargeFormRef" label-width="100px">
        <el-form-item :label="$t('balance.userId')" prop="user_id">
          <el-input v-model="rechargeForm.user_id" :placeholder="$t('balance.userIdPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('balance.rechargeAmount')" prop="amount">
          <el-input-number v-model="rechargeForm.amount" :min="0.01" :precision="2" :step="100" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('common.remark')" prop="remark">
          <el-input v-model="rechargeForm.remark" type="textarea" :placeholder="$t('balance.remarkPlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rechargeDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleRecharge" :loading="submitLoading">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 扣款对话框 -->
    <el-dialog v-model="deductDialogVisible" :title="$t('balance.userDeduct')" width="500px">
      <el-form :model="deductForm" :rules="deductRules" ref="deductFormRef" label-width="100px">
        <el-form-item :label="$t('balance.userId')" prop="user_id">
          <el-input v-model="deductForm.user_id" :placeholder="$t('balance.userIdPlaceholder')" />
        </el-form-item>
        <el-form-item :label="$t('balance.deductAmount')" prop="amount">
          <el-input-number v-model="deductForm.amount" :min="0.01" :precision="2" :step="100" style="width: 100%" />
        </el-form-item>
        <el-form-item :label="$t('common.remark')" prop="remark">
          <el-input v-model="deductForm.remark" type="textarea" :placeholder="$t('balance.remarkPlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="deductDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="warning" @click="handleDeduct" :loading="submitLoading">{{ $t('balance.confirmDeduct') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search, Plus, Minus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { $t } from '@/locales'
import request from '@/utils/http'

defineOptions({ name: 'BalanceManage' })

const loading = ref(false)
const submitLoading = ref(false)

const searchForm = reactive({
  user_id: '',
  username: '',
  type: undefined as string | undefined,
  date_range: [] as string[]
})

const pagination = reactive({
  page: 1,
  page_size: 20,
  total: 0
})

const tableData = ref([])

const rechargeDialogVisible = ref(false)
const rechargeFormRef = ref<FormInstance>()
const rechargeForm = reactive({
  user_id: '',
  amount: 0,
  remark: ''
})

const deductDialogVisible = ref(false)
const deductFormRef = ref<FormInstance>()
const deductForm = reactive({
  user_id: '',
  amount: 0,
  remark: ''
})

const rechargeRules: FormRules = {
  user_id: [
    { required: true, message: $t('balance.userIdPlaceholder'), trigger: 'blur' }
  ],
  amount: [
    { required: true, message: $t('balance.inputRechargeAmount'), trigger: 'blur' }
  ]
}

const deductRules: FormRules = {
  user_id: [
    { required: true, message: $t('balance.userIdPlaceholder'), trigger: 'blur' }
  ],
  amount: [
    { required: true, message: $t('balance.inputDeductAmount'), trigger: 'blur' }
  ]
}

const getTypeTag = (type: string) => {
  const map: Record<string, any> = {
    recharge: 'success',
    deduct: 'danger',
    payment: 'warning',
    refund: 'info'
  }
  return map[type] || 'info'
}

const getTypeText = (type: string) => {
  const map: Record<string, () => string> = {
    recharge: () => $t('balance.recharge'),
    deduct: () => $t('balance.deduct'),
    payment: () => $t('balance.payment'),
    refund: () => $t('balance.refund')
  }
  return map[type]?.() || $t('common.unknown')
}

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
    ElMessage.error($t('balance.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchBalances()
}

const handleReset = () => {
  searchForm.user_id = ''
  searchForm.username = ''
  searchForm.type = undefined
  searchForm.date_range = []
  handleSearch()
}

const showRechargeDialog = () => {
  rechargeForm.user_id = ''
  rechargeForm.amount = 0
  rechargeForm.remark = ''
  rechargeDialogVisible.value = true
}

const showDeductDialog = () => {
  deductForm.user_id = ''
  deductForm.amount = 0
  deductForm.remark = ''
  deductDialogVisible.value = true
}

const handleRecharge = async () => {
  if (!rechargeFormRef.value) return

  await rechargeFormRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      await request.post({
        url: `/api/admin/user-manage/${rechargeForm.user_id}/balance`,
        params: {
          amount: rechargeForm.amount,
          description: rechargeForm.remark || $t('balance.adminRecharge')
        },
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

const handleDeduct = async () => {
  if (!deductFormRef.value) return

  await deductFormRef.value.validate(async (valid) => {
    if (!valid) return

    submitLoading.value = true
    try {
      await request.post({
        url: `/api/admin/user-manage/${deductForm.user_id}/balance`,
        params: {
          amount: -deductForm.amount,
          description: deductForm.remark || $t('balance.adminDeduct')
        },
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

const handleSizeChange = () => {
  pagination.page = 1
  fetchBalances()
}

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
