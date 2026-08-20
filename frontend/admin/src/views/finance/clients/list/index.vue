<template>
  <div class="customer-list-page">
    <!-- 页面描述 -->
    <el-card shadow="never" class="desc-card">
      <div class="desc-content">
        <span>{{ $t('customerList.description') }}</span>
        <el-link type="primary" href="https://bbs.idcsmart.com/forum.php?mod=viewthread&tid=136" target="_blank">{{ $t('customerList.helpDoc') }}</el-link>
      </div>
      <div class="action-buttons">
        <el-button type="primary" @click="handleAdd">
          <el-icon><Plus /></el-icon>
          {{ $t('customerList.addCustomer') }}
        </el-button>
        <el-button @click="showAdvancedSearch = !showAdvancedSearch">
          {{ $t('customerList.advancedSearch') }}
        </el-button>
      </div>
    </el-card>

    <!-- 高级搜索区域 -->
    <el-card v-if="showAdvancedSearch" shadow="never" class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item :label="$t('customerList.keyword')">
          <el-input
            v-model="searchForm.keyword"
            :placeholder="$t('customerList.keywordPlaceholder')"
            clearable
            style="width: 200px"
            @keyup.enter="handleSearch"
          />
        </el-form-item>
        <el-form-item :label="$t('customerList.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('common.all')" clearable style="width: 100px">
            <el-option :label="$t('customerList.active')" value="active" />
            <el-option :label="$t('customerList.disabled')" value="disabled" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('customerList.clientGroup')">
          <el-select v-model="searchForm.group_id" :placeholder="$t('common.all')" clearable style="width: 120px">
            <el-option v-for="group in clientGroups" :key="group.id" :label="group.name" :value="group.id" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 数据表格 -->
    <el-card shadow="never" class="table-card">
      <el-table
        v-loading="loading"
        :data="tableData"
        border
        stripe
        @sort-change="handleSortChange"
      >
        <el-table-column prop="id" label="ID" width="70" sortable="custom" align="center">
          <template #default="{ row }">
            <el-link type="primary" @click="handleView(row)">{{ row.id }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="username" :label="$t('customerList.name')" min-width="120">
          <template #default="{ row }">
            <el-link type="primary" @click="handleView(row)">{{ row.username }}</el-link>
          </template>
        </el-table-column>
        <el-table-column :label="$t('customerList.phoneEmail')" min-width="180">
          <template #default="{ row }">
            <div>{{ row.phone || '-' }}</div>
            <div class="text-secondary">{{ row.email || '-' }}</div>
          </template>
        </el-table-column>
        <el-table-column :label="$t('customerList.service')" width="80" align="center">
          <template #default="{ row }">
            <el-link type="primary" @click="handleViewProducts(row)">
              {{ row.products_count || 0 }}({{ row.active_products || 0 }})
            </el-link>
          </template>
        </el-table-column>
        <el-table-column :label="$t('customerList.incomeExpense')" width="120" align="right">
          <template #default="{ row }">
            <div>¥{{ formatMoney(row.total_income) }}</div>
            <div class="text-secondary">¥{{ formatMoney(row.total_expense) }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="balance" :label="$t('customerList.balance')" width="90" align="right" sortable="custom">
          <template #default="{ row }">
            ¥{{ formatMoney(row.balance) }}
          </template>
        </el-table-column>
        <el-table-column prop="group_name" :label="$t('customerList.clientGroup')" width="100">
          <template #default="{ row }">
            {{ row.group_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('customerList.status')" width="80" align="center" sortable="custom">
          <template #default="{ row }">
            {{ getStatusText(row.status) }}
          </template>
        </el-table-column>
        <el-table-column prop="sales_name" :label="$t('customerList.sales')" width="80">
          <template #default="{ row }">
            {{ row.sales_name || '-' }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('customerList.creditLimit')" width="140" align="right">
          <template #default="{ row }">
            <div>¥{{ formatMoney(row.credit_used) }}</div>
            <div class="text-secondary">¥{{ formatMoney(row.credit_total) }}</div>
          </template>
        </el-table-column>
        <el-table-column prop="api_created_at" :label="$t('customerList.apiCreatedAt')" width="120" sortable="custom">
          <template #default="{ row }">
            {{ row.api_created_at || '-' }}
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
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/http'
import { $t } from '@/locales'

const router = useRouter()
const loading = ref(false)
const tableData = ref([])
const clientGroups = ref<{ id: number; name: string }[]>([])
const showAdvancedSearch = ref(false)

const searchForm = reactive({
  keyword: '',
  status: '',
  group_id: ''
})

const pagination = reactive({
  page: 1,
  page_size: 100,
  total: 0
})

const sortParams = reactive({
  sort: 'id',
  order: 'DESC'
})

const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

const getStatusText = (status: string) => {
  const map: Record<string, () => string> = {
    active: () => $t('customerList.active'),
    disabled: () => $t('customerList.disabled'),
    pending: () => $t('customerList.pending')
  }
  return map[status]?.() || $t('common.unknown')
}

const fetchList = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      sort: sortParams.sort,
      order: sortParams.order
    }
    if (searchForm.keyword) params.keyword = searchForm.keyword
    if (searchForm.status) params.status = searchForm.status
    if (searchForm.group_id) params.group_id = searchForm.group_id
    const res = await request.get({ url: '/api/admin/clients', params })
    tableData.value = res?.data || res?.list || []
    pagination.total = res?.total || 0
  } catch (error) {
    console.error('fetch client list failed:', error)
  } finally {
    loading.value = false
  }
}

const fetchGroups = async () => {
  try {
    const data = await request.get({ url: '/api/admin/client-groups' })
    clientGroups.value = data || []
  } catch (error) {
    console.error('fetch client groups failed:', error)
  }
}

const handleSearch = () => { pagination.page = 1; fetchList() }
const handleReset = () => { searchForm.keyword = ''; searchForm.status = ''; searchForm.group_id = ''; handleSearch() }
const handleSizeChange = () => { pagination.page = 1; fetchList() }
const handlePageChange = () => { fetchList() }
const handleSortChange = ({ prop, order }: any) => { sortParams.sort = prop || 'id'; sortParams.order = order === 'ascending' ? 'ASC' : 'DESC'; fetchList() }
const handleAdd = () => { router.push('/customer-add') }
const handleView = (row: any) => { router.push(`/customer-view/${row.id}`) }
const handleViewProducts = (row: any) => { router.push(`/customer-view/${row.id}?tab=products`) }

onMounted(() => { fetchList(); fetchGroups() })
</script>

<style scoped lang="scss">
.customer-list-page {
  padding: 16px;
}
.desc-card {
  margin-bottom: 16px;
  .desc-content {
    display: flex;
    align-items: center;
    gap: 8px;
    margin-bottom: 12px;
  }
  .action-buttons {
    display: flex;
    gap: 8px;
  }
}
.search-card {
  margin-bottom: 16px;
}
.table-card {
  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    margin-top: 16px;
  }
}
.text-secondary {
  color: var(--el-text-color-secondary);
  font-size: 12px;
}
</style>
