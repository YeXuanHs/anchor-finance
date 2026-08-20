<template>
  <div class="product-list-page">
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane :label="$t('productList.all')" name="all" />
      <el-tab-pane :label="$t('productList.server')" name="server" />
      <el-tab-pane :label="$t('productList.cloud')" name="cloud" />
      <el-tab-pane :label="$t('productList.hosting')" name="hosting" />
      <el-tab-pane :label="$t('productList.domain')" name="domain" />
      <el-tab-pane :label="$t('productList.ssl')" name="ssl" />
      <el-tab-pane :label="$t('productList.other')" name="other" />
    </el-tabs>

    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <div class="action-left">
          <el-button type="primary" @click="handleAdd"><el-icon><Plus /></el-icon>{{ $t('productList.addProduct') }}</el-button>
          <el-button @click="handleExport"><el-icon><Download /></el-icon>{{ $t('common.export') }}</el-button>
        </div>
        <div class="action-right">
          <el-input v-model="searchKeyword" :placeholder="$t('productList.searchPlaceholder')" clearable style="width: 200px" @keyup.enter="fetchList"><template #prefix><el-icon><Search /></el-icon></template></el-input>
          <el-button circle @click="fetchList"><el-icon><Refresh /></el-icon></el-button>
        </div>
      </div>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" :label="$t('productList.productName')" min-width="200" />
        <el-table-column prop="category" :label="$t('productList.category')" width="100" />
        <el-table-column prop="price" :label="$t('productList.price')" width="120" align="right"><template #default="{ row }">¥{{ formatMoney(row.price) }}</template></el-table-column>
        <el-table-column prop="billing_cycle" :label="$t('productList.billingCycle')" width="100" />
        <el-table-column prop="stock" :label="$t('productList.stock')" width="80" align="center" />
        <el-table-column prop="sales" :label="$t('productList.sales')" width="80" align="center" />
        <el-table-column prop="status" :label="$t('productList.status')" width="80" align="center">
          <template #default="{ row }"><el-switch v-model="row.status" :active-value="'active'" :inactive-value="'disabled'" @change="handleToggleStatus(row)" /></template>
        </el-table-column>
        <el-table-column :label="$t('productList.operations')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="warning" link size="small" @click="handlePricing(row)">{{ $t('productList.pricing') }}</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">{{ $t('common.delete') }}</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[10, 20, 50, 100]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh, Download } from '@element-plus/icons-vue'
import request from '@/utils/http'
import { $t } from '@/locales'

const router = useRouter()
const loading = ref(false)
const tableData = ref([])
const activeTab = ref('all')
const searchKeyword = ref('')
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

const formatMoney = (amount: number) => { if (!amount) return '0.00'; return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) }

const handleTabChange = () => { pagination.page = 1; fetchList() }

const fetchList = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (activeTab.value !== 'all') params.category = activeTab.value
    if (searchKeyword.value) params.keyword = searchKeyword.value
    const data = await request.get({ url: '/api/admin/products', params })
    tableData.value = data?.list || []; pagination.total = data?.total || 0
  } catch {} finally { loading.value = false }
}

const handleSizeChange = (size: number) => { pagination.page_size = size; pagination.page = 1; fetchList() }
const handlePageChange = (page: number) => { pagination.page = page; fetchList() }
const handleAdd = () => { router.push('/product-add') }
const handleEdit = (row: any) => { router.push(`/product-edit/${row.id}`) }
const handlePricing = (row: any) => { router.push(`/product-pricing/${row.id}`) }

const handleToggleStatus = async (row: any) => {
  try { await request.post({ url: `/api/admin/products/${row.id}/status`, data: { status: row.status } }); ElMessage.success($t('common.operationSuccess')) } catch { fetchList() }
}

const handleDelete = async (row: any) => {
  try { await ElMessageBox.confirm($t('productList.confirmDelete', { name: row.name }), $t('common.tips'), { type: 'warning' }); await request.del({ url: `/api/admin/products/${row.id}` }); ElMessage.success($t('common.deleteSuccess')); fetchList() } catch (error) { if (error !== 'cancel') console.error('delete failed:', error) }
}

const handleExport = async () => {
  try {
    const response = await request.get({ url: '/api/admin/products/export', responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response])); const link = document.createElement('a'); link.href = url
    link.setAttribute('download', `products_${new Date().toISOString().split('T')[0]}.xlsx`)
    document.body.appendChild(link); link.click(); document.body.removeChild(link); window.URL.revokeObjectURL(url)
  } catch {}
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.product-list-page { padding: 16px; }
.action-card { margin-bottom: 16px; }
.action-bar { display: flex; justify-content: space-between; align-items: center; }
.action-right { display: flex; gap: 8px; }
.table-card { :deep(.el-card__body) { padding: 0; } }
.pagination-wrapper { display: flex; justify-content: flex-end; padding: 16px; }
</style>
