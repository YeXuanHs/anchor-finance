<template>
  <div class="withdraw-page">
    <!-- 标签页 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="全部" name="all" />
      <el-tab-pane label="待审核" name="pending" />
      <el-tab-pane label="已批准" name="approved" />
      <el-tab-pane label="已拒绝" name="rejected" />
      <el-tab-pane label="已完成" name="completed" />
    </el-tabs>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="withdraw_no" label="提现单号" width="150" />
        <el-table-column prop="client_name" label="客户" width="120">
          <template #default="{ row }">
            <el-button type="primary" link @click="$router.push(`/customer-view/${row.client_id}`)">
              {{ row.client_name }}
            </el-button>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.amount) }}</template>
        </el-table-column>
        <el-table-column prop="fee" label="手续费" width="100" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.fee) }}</template>
        </el-table-column>
        <el-table-column prop="actual_amount" label="实际到账" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.actual_amount) }}</template>
        </el-table-column>
        <el-table-column prop="withdraw_method" label="提现方式" width="100" />
        <el-table-column prop="account_info" label="收款信息" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="170" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'pending'"
              type="success"
              link
              size="small"
              @click="handleApprove(row)"
            >
              批准
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="danger"
              link
              size="small"
              @click="handleReject(row)"
            >
              拒绝
            </el-button>
            <el-button
              v-if="row.status === 'approved'"
              type="primary"
              link
              size="small"
              @click="handleComplete(row)"
            >
              标记完成
            </el-button>
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const tableData = ref([])
const activeTab = ref('all')

// 分页
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

// 格式化金额
const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// 状态类型
const getStatusType = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { pending: 'warning', approved: 'primary', rejected: 'danger', completed: 'success' }
  return map[status] || 'info'
}

// 状态文本
const getStatusText = (status: string) => {
  const map: Record<string, string> = { pending: '待审核', approved: '已批准', rejected: '已拒绝', completed: '已完成' }
  return map[status] || '未知'
}

// 标签页切换
const handleTabChange = () => { pagination.page = 1; fetchList() }

// 获取列表数据
const fetchList = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (activeTab.value !== 'all') params.status = activeTab.value
    const data = await request.get({ url: '/api/admin/withdrawals', params })
    tableData.value = data?.list || []
    pagination.total = data?.total || 0
  } catch (error) {
    console.error('获取提现列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 分页
const handleSizeChange = (size: number) => { pagination.page_size = size; pagination.page = 1; fetchList() }
const handlePageChange = (page: number) => { pagination.page = page; fetchList() }

// 批准
const handleApprove = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要批准此提现申请吗？', '确认批准', { type: 'warning' })
    await request.post({ url: `/api/admin/withdrawals/${row.id}/approve` })
    ElMessage.success('已批准')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') console.error('批准失败:', error)
  }
}

// 拒绝
const handleReject = async (row: any) => {
  try {
    const { value: reason } = await ElMessageBox.prompt('请输入拒绝原因', '拒绝提现', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputValidator: (v) => !!v || '请输入拒绝原因'
    })
    await request.post({ url: `/api/admin/withdrawals/${row.id}/reject`, data: { reason } })
    ElMessage.success('已拒绝')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') console.error('拒绝失败:', error)
  }
}

// 标记完成
const handleComplete = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定要标记此提现为已完成吗？', '确认完成', { type: 'warning' })
    await request.post({ url: `/api/admin/withdrawals/${row.id}/complete` })
    ElMessage.success('已标记完成')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') console.error('标记失败:', error)
  }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.withdraw-page {
  padding: 16px;
}

.table-card {
  :deep(.el-card__body) { padding: 0; }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  padding: 16px;
}
</style>
