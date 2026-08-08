<template>
  <div class="customer-products-page">
    <!-- 标签页 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="全部" name="all" />
      <el-tab-pane label="使用中" name="active" />
      <el-tab-pane label="已暂停" name="suspended" />
      <el-tab-pane label="待开通" name="pending" />
      <el-tab-pane label="已到期" name="expired" />
    </el-tabs>

    <!-- 搜索 -->
    <el-card shadow="never" class="search-card">
      <el-form inline>
        <el-form-item label="关键词">
          <el-input v-model="searchKeyword" placeholder="客户名/产品名/域名" clearable style="width: 200px" @keyup.enter="fetchList" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchList"><el-icon><Search /></el-icon>搜索</el-button>
          <el-button @click="searchKeyword = ''; fetchList()"><el-icon><Refresh /></el-icon>重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="client_name" label="客户" width="120">
          <template #default="{ row }">
            <el-button type="primary" link @click="$router.push(`/customer-view/${row.client_id}`)">{{ row.client_name }}</el-button>
          </template>
        </el-table-column>
        <el-table-column prop="product_name" label="产品" min-width="150" />
        <el-table-column prop="domain" label="域名/标识" width="200" />
        <el-table-column prop="ip" label="IP地址" width="130" />
        <el-table-column prop="billing_cycle" label="计费周期" width="100" />
        <el-table-column prop="amount" label="金额" width="100" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.amount) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="expired_at" label="到期时间" width="120" />
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleManage(row)">管理</el-button>
            <el-button type="warning" link size="small" @click="handleRenew(row)">续费</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination-wrapper">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10,20,50,100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search, Refresh } from '@element-plus/icons-vue'
import request from '@/utils/http'

const loading = ref(false)
const tableData = ref([])
const activeTab = ref('all')
const searchKeyword = ref('')
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const getStatusType = (status: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' => {
  const map: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = { active: 'success', suspended: 'danger', pending: 'warning', expired: 'info' }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = { active: '使用中', suspended: '已暂停', pending: '待开通', expired: '已到期' }
  return map[status] || '未知'
}

const handleTabChange = () => { pagination.page = 1; fetchList() }

const fetchList = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (activeTab.value !== 'all') params.status = activeTab.value
    if (searchKeyword.value) params.keyword = searchKeyword.value
    const data = await request.get({ url: '/api/admin/products/services', params })
    tableData.value = data?.list || []
    pagination.total = data?.total || 0
  } catch (error) {
    console.error('获取产品列表失败:', error)
  } finally {
    loading.value = false
  }
}

const handleSizeChange = (size: number) => { pagination.page_size = size; pagination.page = 1; fetchList() }
const handlePageChange = (page: number) => { pagination.page = page; fetchList() }
const handleManage = (row: any) => { console.log('管理:', row.id) }
const handleRenew = (row: any) => { console.log('续费:', row.id) }

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.customer-products-page { padding: 16px; }
.search-card { margin-bottom: 16px; :deep(.el-card__body) { padding-bottom: 0; } }
.table-card { :deep(.el-card__body) { padding: 0; } }
.pagination-wrapper { display: flex; justify-content: flex-end; padding: 16px; }
</style>
