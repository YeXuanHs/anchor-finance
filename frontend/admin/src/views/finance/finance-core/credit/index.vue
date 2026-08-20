<template>
  <div class="credit-management-page">
    <!-- 搜索筛选区域 -->
    <el-card shadow="never" class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item :label="$t('credit.keyword')">
          <el-input
            v-model="searchForm.keyword"
            :placeholder="$t('credit.keywordPlaceholder')"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            {{ $t('common.search') }}
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            {{ $t('common.reset') }}
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="username" :label="$t('common.username')" width="120">
          <template #default="{ row }">
            <el-button type="primary" link @click="$router.push(`/customer-view/${row.id}`)">
              {{ row.username }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="email" :label="$t('common.email')" min-width="200" />
        <el-table-column prop="balance" :label="$t('common.balance')" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.balance) }}</template>
        </el-table-column>
        <el-table-column prop="credit" :label="$t('common.credit')" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.credit) }}</template>
        </el-table-column>
        <el-table-column prop="credit_used" :label="$t('credit.creditUsed')" width="120" align="right">
          <template #default="{ row }">
            <span :class="row.credit_used > 0 ? 'text-red' : ''">¥{{ formatMoney(row.credit_used) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="credit_available" :label="$t('credit.creditAvailable')" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.credit_available) }}</template>
        </el-table-column>
        <el-table-column :label="$t('common.action')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleAdjustCredit(row)">{{ $t('credit.adjustCredit') }}</el-button>
            <el-button type="warning" link size="small" @click="handleViewLogs(row)">{{ $t('credit.operationLogs') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <div class="pagination-wrapper">
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

    <!-- 调整信用额弹窗 -->
    <el-dialog v-model="creditDialogVisible" :title="$t('credit.adjustCredit')" width="500px">
      <el-form :model="creditForm" label-width="100px">
        <el-form-item :label="$t('credit.customer')">
          <span>{{ currentClient?.username }}</span>
        </el-form-item>
        <el-form-item :label="$t('credit.currentCredit')">
          <span>¥{{ formatMoney(currentClient?.credit || 0) }}</span>
        </el-form-item>
        <el-form-item :label="$t('credit.newCredit')">
          <el-input-number v-model="creditForm.credit" :min="0" :precision="2" :step="100" />
        </el-form-item>
        <el-form-item :label="$t('credit.reason')">
          <el-input v-model="creditForm.reason" type="textarea" :rows="3" :placeholder="$t('credit.reasonPlaceholder')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="creditDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" :loading="saving" @click="handleSaveCredit">{{ $t('common.confirm') }}</el-button>
      </template>
    </el-dialog>

    <!-- 操作记录弹窗 -->
    <el-dialog v-model="logsDialogVisible" :title="$t('credit.operationLogs')" width="700px">
      <el-table :data="creditLogs" border stripe>
        <el-table-column prop="created_at" :label="$t('credit.time')" width="170" />
        <el-table-column prop="operator_name" :label="$t('common.operator')" width="100" />
        <el-table-column prop="type" :label="$t('credit.type')" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.type === 'increase' ? 'success' : 'danger'" size="small">
              {{ row.type === 'increase' ? $t('credit.increase') : $t('credit.decrease') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" :label="$t('common.amount')" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.amount) }}</template>
        </el-table-column>
        <el-table-column prop="before" :label="$t('credit.beforeAdjust')" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.before) }}</template>
        </el-table-column>
        <el-table-column prop="after" :label="$t('credit.afterAdjust')" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.after) }}</template>
        </el-table-column>
        <el-table-column prop="reason" :label="$t('credit.reason')" min-width="150" show-overflow-tooltip />
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'
import { $t } from '@/locales'
import request from '@/utils/http'

const loading = ref(false)
const tableData = ref([])

const searchForm = reactive({ keyword: '' })

const pagination = reactive({ page: 1, page_size: 20, total: 0 })

const creditDialogVisible = ref(false)
const currentClient = ref<any>(null)
const creditForm = reactive({ credit: 0, reason: '' })
const saving = ref(false)

const logsDialogVisible = ref(false)
const creditLogs = ref([])

const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const fetchList = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (searchForm.keyword) params.keyword = searchForm.keyword
    const data = await request.get({ url: '/api/admin/credit/clients', params })
    tableData.value = data?.list || []
    pagination.total = data?.total || 0
  } catch (error) {
    console.error('获取信用列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchList() }
const handleReset = () => { searchForm.keyword = ''; pagination.page = 1; fetchList() }
const handleSizeChange = (size: number) => { pagination.page_size = size; pagination.page = 1; fetchList() }
const handlePageChange = (page: number) => { pagination.page = page; fetchList() }

const handleAdjustCredit = (row: any) => {
  currentClient.value = row
  creditForm.credit = row.credit || 0
  creditForm.reason = ''
  creditDialogVisible.value = true
}

const handleSaveCredit = async () => {
  saving.value = true
  try {
    await request.post({
      url: `/api/admin/clients/${currentClient.value.id}/credit`,
      data: creditForm
    })
    ElMessage.success($t('credit.adjustSuccess'))
    creditDialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('调整信用额失败:', error)
  } finally {
    saving.value = false
  }
}

const handleViewLogs = async (row: any) => {
  try {
    const data = await request.get({ url: `/api/admin/clients/${row.id}/credit-logs` })
    creditLogs.value = data || []
  } catch (error) {
    console.error('获取操作记录失败:', error)
  }
  logsDialogVisible.value = true
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.credit-management-page {
  padding: 16px;
}

.search-card {
  margin-bottom: 16px;
  :deep(.el-card__body) { padding-bottom: 0; }
}

.table-card {
  :deep(.el-card__body) { padding: 0; }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px;
}

.text-red { color: #EF4444; }
</style>
