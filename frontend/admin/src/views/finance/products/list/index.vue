<template>
  <div class="product-list-page">
    <!-- 标签页 -->
    <el-tabs v-model="activeTab" @tab-change="handleTabChange">
      <el-tab-pane label="全部产品" name="all" />
      <el-tab-pane label="独立服务器" name="server" />
      <el-tab-pane label="云服务器" name="cloud" />
      <el-tab-pane label="虚拟主机" name="hosting" />
      <el-tab-pane label="域名" name="domain" />
      <el-tab-pane label="SSL证书" name="ssl" />
      <el-tab-pane label="其他" name="other" />
    </el-tabs>

    <!-- 操作栏 -->
    <el-card shadow="never" class="action-card">
      <div class="action-bar">
        <div class="action-left">
          <el-button type="primary" @click="handleAdd">
            <el-icon><Plus /></el-icon>
            添加产品
          </el-button>
          <el-button @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </div>
        <div class="action-right">
          <el-input
            v-model="searchKeyword"
            placeholder="搜索产品名称"
            clearable
            style="width: 200px"
            @keyup.enter="fetchList"
          >
            <template #prefix>
              <el-icon><Search /></el-icon>
            </template>
          </el-input>
          <el-button circle @click="fetchList">
            <el-icon><Refresh /></el-icon>
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="name" label="产品名称" min-width="200" />
        <el-table-column prop="category" label="分类" width="100" />
        <el-table-column prop="price" label="价格" width="120" align="right">
          <template #default="{ row }">¥{{ formatMoney(row.price) }}</template>
        </el-table-column>
        <el-table-column prop="billing_cycle" label="计费周期" width="100" />
        <el-table-column prop="stock" label="库存" width="80" align="center" />
        <el-table-column prop="sales" label="销量" width="80" align="center" />
        <el-table-column prop="status" label="状态" width="80" align="center">
          <template #default="{ row }">
            <el-switch
              v-model="row.status"
              :active-value="'active'"
              :inactive-value="'disabled'"
              @change="handleToggleStatus(row)"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="warning" link size="small" @click="handlePricing(row)">定价</el-button>
            <el-button type="danger" link size="small" @click="handleDelete(row)">删除</el-button>
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
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search, Refresh, Download } from '@element-plus/icons-vue'
import request from '@/utils/http'

const router = useRouter()
const loading = ref(false)
const tableData = ref([])
const activeTab = ref('all')
const searchKeyword = ref('')

// 分页
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

// 格式化金额
const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// 标签页切换
const handleTabChange = () => { pagination.page = 1; fetchList() }

// 获取列表数据
const fetchList = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (activeTab.value !== 'all') params.category = activeTab.value
    if (searchKeyword.value) params.keyword = searchKeyword.value
    const data = await request.get({ url: '/api/admin/products', params })
    tableData.value = data?.list || []
    pagination.total = data?.total || 0
  } catch (error) {
    console.error('获取产品列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 分页
const handleSizeChange = (size: number) => { pagination.page_size = size; pagination.page = 1; fetchList() }
const handlePageChange = (page: number) => { pagination.page = page; fetchList() }

// 添加产品
const handleAdd = () => { router.push('/product-add') }

// 编辑产品
const handleEdit = (row: any) => { router.push(`/product-edit/${row.id}`) }

// 定价
const handlePricing = (row: any) => { router.push(`/product-pricing/${row.id}`) }

// 切换状态
const handleToggleStatus = async (row: any) => {
  try {
    await request.post({ url: `/api/admin/products/${row.id}/status`, data: { status: row.status } })
    ElMessage.success('状态更新成功')
  } catch (error) {
    console.error('更新状态失败:', error)
    fetchList()
  }
}

// 删除产品
const handleDelete = async (row: any) => {
  try {
    await ElMessageBox.confirm(`确定要删除产品 "${row.name}" 吗？`, '确认删除', { type: 'warning' })
    await request.del({ url: `/api/admin/products/${row.id}` })
    ElMessage.success('删除成功')
    fetchList()
  } catch (error) {
    if (error !== 'cancel') console.error('删除失败:', error)
  }
}

// 导出
const handleExport = async () => {
  try {
    const response = await request.get({ url: '/api/admin/products/export', responseType: 'blob' })
    const url = window.URL.createObjectURL(new Blob([response]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `产品列表_${new Date().toISOString().split('T')[0]}.xlsx`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
  } catch (error) {
    console.error('导出失败:', error)
  }
}

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.product-list-page {
  padding: 16px;
}

.action-card {
  margin-bottom: 16px;
}

.action-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.action-right {
  display: flex;
  gap: 8px;
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
