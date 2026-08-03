<template>
  <div class="marketplace-logs-page">
    <div class="page-header">
      <h1 class="page-title">操作日志</h1>
      <div class="header-right">
        <el-select v-model="actionFilter" placeholder="操作类型" clearable style="width: 130px">
          <el-option label="全部" value="" />
          <el-option label="上架" value="list" />
          <el-option label="下架" value="delist" />
          <el-option label="修改" value="update" />
          <el-option label="售出" value="sold" />
          <el-option label="购买" value="purchase" />
        </el-select>
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          style="width: 280px"
        />
        <el-button type="primary" @click="fetchList">
          <el-icon><Search /></el-icon> 搜索
        </el-button>
        <el-button @click="handleReset">重置</el-button>
      </div>
    </div>

    <!-- 日志时间线 -->
    <el-card shadow="never" class="logs-card">
      <el-table :data="tableData" stripe v-loading="loading" empty-text="暂无操作日志">
        <el-table-column prop="id" label="日志ID" width="100">
          <template #default="{ row }">
            <span class="mono">{{ row.id }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="操作时间" width="170" />
        <el-table-column prop="action" label="操作类型" width="100">
          <template #default="{ row }">
            <el-tag :type="actionTagType(row.action)" size="small">
              {{ actionText(row.action) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="product_name" label="商品名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="listing_id" label="商品ID" width="120">
          <template #default="{ row }">
            <span class="mono">{{ row.listing_id || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="detail" label="操作详情" min-width="250" show-overflow-tooltip />
        <el-table-column prop="operator" label="操作人" width="120" />
        <el-table-column prop="ip" label="IP地址" width="140">
          <template #default="{ row }">
            <span class="mono">{{ row.ip || '-' }}</span>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Search } from '@element-plus/icons-vue'
import request from '@/utils/request'

const loading = ref(false)
const actionFilter = ref('')
const dateRange = ref<string[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const tableData = ref<any[]>([])

const actionTagType = (action: string) => {
  const map: Record<string, string> = {
    list: 'success',
    delist: 'warning',
    update: 'primary',
    sold: 'danger',
    purchase: 'primary'
  }
  return map[action] || 'info'
}

const actionText = (action: string) => {
  const map: Record<string, string> = {
    list: '上架',
    delist: '下架',
    update: '修改',
    sold: '售出',
    purchase: '购买'
  }
  return map[action] || action
}

function handleReset() {
  actionFilter.value = ''
  dateRange.value = []
  currentPage.value = 1
  fetchList()
}

async function fetchList() {
  loading.value = true
  try {
    const params: any = { page: currentPage.value, page_size: pageSize.value }
    if (actionFilter.value) params.action = actionFilter.value
    if (dateRange.value?.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    const res = await request.get('/api/v1/marketplace/logs', { params })
    tableData.value = res.data?.data?.list || res.data?.list || []
    total.value = res.data?.data?.total || 0
  } catch {
    tableData.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.marketplace-logs-page {
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    flex-wrap: wrap;
    gap: 12px;

    .page-title {
      font-size: 20px;
      font-weight: 700;
      color: #303133;
      margin: 0;
    }

    .header-right {
      display: flex;
      gap: 12px;
      align-items: center;
    }
  }

  .logs-card {
    border-radius: 12px;
    border: 1px solid #e8ecf1;
  }

  .logs-card :deep(.el-card__body) { padding: 0; }
  .logs-card :deep(.el-table th.el-table__cell) { background: #fafbfc; font-weight: 600; }

  .mono {
    font-family: 'Monaco', 'Menlo', monospace;
    font-size: 13px;
    color: #606266;
  }

  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    padding: 16px 20px;
    border-top: 1px solid #e8ecf1;
  }

  @media (max-width: 768px) {
    .page-header { flex-direction: column; align-items: flex-start; }
    .header-right { width: 100%; flex-wrap: wrap; }
  }
}
</style>
