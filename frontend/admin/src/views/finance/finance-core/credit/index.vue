<template>
  <div class="credit-management-page">
    <!-- 搜索筛选区域 -->
    <el-card shadow="never" class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item label="关键词">
          <el-input
            v-model="searchForm.keyword"
            placeholder="客户名/邮箱"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">
            <el-icon><Refresh /></el-icon>
            重置
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="username" label="客户名" width="120">
          <template #default="{ row }">
            <el-button type="primary" link @click="$router.push(`/customer-view/${row.id}`)">
              {{ row.username }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="email" label="邮箱" min-width="200" />
        <el-table-column prop="balance" label="余额" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.balance) }}</template>
        </el-table-column>
        <el-table-column prop="credit" label="信用额" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.credit) }}</template>
        </el-table-column>
        <el-table-column prop="credit_used" label="已用信用" width="120" align="right">
          <template #default="{ row }">
            <span :class="row.credit_used > 0 ? 'text-red' : ''">¥{{ formatMoney(row.credit_used) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="credit_available" label="可用信用" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.credit_available) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleAdjustCredit(row)">调整信用额</el-button>
            <el-button type="warning" link size="small" @click="handleViewLogs(row)">操作记录</el-button>
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
    <el-dialog v-model="creditDialogVisible" title="调整信用额" width="500px">
      <el-form :model="creditForm" label-width="100px">
        <el-form-item label="客户">
          <span>{{ currentClient?.username }}</span>
        </el-form-item>
        <el-form-item label="当前信用额">
          <span>¥{{ formatMoney(currentClient?.credit || 0) }}</span>
        </el-form-item>
        <el-form-item label="新信用额">
          <el-input-number v-model="creditForm.credit" :min="0" :precision="2" :step="100" />
        </el-form-item>
        <el-form-item label="原因">
          <el-input v-model="creditForm.reason" type="textarea" :rows="3" placeholder="请输入调整原因" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="creditDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSaveCredit">确定</el-button>
      </template>
    </el-dialog>

    <!-- 操作记录弹窗 -->
    <el-dialog v-model="logsDialogVisible" title="信用操作记录" width="700px">
      <el-table :data="creditLogs" border stripe>
        <el-table-column prop="created_at" label="时间" width="170" />
        <el-table-column prop="operator_name" label="操作人" width="100" />
        <el-table-column prop="type" label="类型" width="80" align="center">
          <template #default="{ row }">
            <el-tag :type="row.type === 'increase' ? 'success' : 'danger'" size="small">
              {{ row.type === 'increase' ? '增加' : '减少' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.amount) }}</template>
        </el-table-column>
        <el-table-column prop="before" label="调整前" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.before) }}</template>
        </el-table-column>
        <el-table-column prop="after" label="调整后" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.after) }}</template>
        </el-table-column>
        <el-table-column prop="reason" label="原因" min-width="150" show-overflow-tooltip />
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Refresh } from '@element-plus/icons-vue'
import request from '@/utils/http'

const loading = ref(false)
const tableData = ref([])

// 搜索表单
const searchForm = reactive({ keyword: '' })

// 分页
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

// 信用调整弹窗
const creditDialogVisible = ref(false)
const currentClient = ref<any>(null)
const creditForm = reactive({ credit: 0, reason: '' })
const saving = ref(false)

// 操作记录弹窗
const logsDialogVisible = ref(false)
const creditLogs = ref([])

// 格式化金额
const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// 获取列表数据
const fetchList = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (searchForm.keyword) params.keyword = searchForm.keyword
    const data = await request.get({ url: '/api/admin/clients/credit', params })
    tableData.value = data?.list || []
    pagination.total = data?.total || 0
  } catch (error) {
    console.error('获取信用列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 搜索
const handleSearch = () => { pagination.page = 1; fetchList() }
const handleReset = () => { searchForm.keyword = ''; pagination.page = 1; fetchList() }
const handleSizeChange = (size: number) => { pagination.page_size = size; pagination.page = 1; fetchList() }
const handlePageChange = (page: number) => { pagination.page = page; fetchList() }

// 调整信用额
const handleAdjustCredit = (row: any) => {
  currentClient.value = row
  creditForm.credit = row.credit || 0
  creditForm.reason = ''
  creditDialogVisible.value = true
}

// 保存信用额
const handleSaveCredit = async () => {
  saving.value = true
  try {
    await request.post({
      url: `/api/admin/clients/${currentClient.value.id}/credit`,
      data: creditForm
    })
    ElMessage.success('信用额调整成功')
    creditDialogVisible.value = false
    fetchList()
  } catch (error) {
    console.error('调整信用额失败:', error)
  } finally {
    saving.value = false
  }
}

// 查看操作记录
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
