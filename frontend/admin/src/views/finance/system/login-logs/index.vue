<template>
  <div class="login-log-page">
    <el-card shadow="never" class="search-card">
      <el-form :model="searchForm" inline>
        <el-form-item :label="$t('loginLog.username')">
          <el-input v-model="searchForm.username" :placeholder="$t('loginLog.username')" clearable style="width: 150px" @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item :label="$t('loginLog.ip')">
          <el-input v-model="searchForm.ip" :placeholder="$t('loginLog.ip')" clearable style="width: 150px" />
        </el-form-item>
        <el-form-item :label="$t('loginLog.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('common.all')" clearable style="width: 100px">
            <el-option :label="$t('loginLog.success')" value="success" />
            <el-option :label="$t('loginLog.failure')" value="failure" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('loginLog.timeRange')">
          <el-date-picker v-model="searchForm.date_range" type="daterange" :range-separator="$t('common.to')" :start-placeholder="$t('common.startDate')" :end-placeholder="$t('common.endDate')" style="width: 240px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch"><el-icon><Search /></el-icon>{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset"><el-icon><Refresh /></el-icon>{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table v-loading="loading" :data="tableData" border stripe>
        <el-table-column prop="id" label="ID" width="80" align="center" />
        <el-table-column prop="username" :label="$t('loginLog.username')" width="120" />
        <el-table-column prop="ip" :label="$t('loginLog.ip')" width="130" />
        <el-table-column prop="location" :label="$t('loginLog.location')" width="150" />
        <el-table-column prop="user_agent" :label="$t('loginLog.userAgent')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="status" :label="$t('loginLog.status')" width="80" align="center">
          <template #default="{ row }"><el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">{{ row.status === 'success' ? $t('loginLog.success') : $t('loginLog.failure') }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="failure_reason" :label="$t('loginLog.failureReason')" width="150" show-overflow-tooltip />
        <el-table-column prop="created_at" :label="$t('loginLog.loginTime')" width="170" />
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size" :page-sizes="[20, 50, 100, 200]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper" @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Search, Refresh } from '@element-plus/icons-vue'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const tableData = ref([])
const searchForm = reactive({ username: '', ip: '', status: '', date_range: null as [Date, Date] | null })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })

const fetchList = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size }
    if (searchForm.username) params.username = searchForm.username
    if (searchForm.ip) params.ip = searchForm.ip
    if (searchForm.status) params.status = searchForm.status
    if (searchForm.date_range) { params.start_date = searchForm.date_range[0].toISOString().split('T')[0]; params.end_date = searchForm.date_range[1].toISOString().split('T')[0] }
    const data = await request.get({ url: '/api/admin/login-logs', params }); tableData.value = data.list || []; pagination.total = data.total || 0
  } catch {} finally { loading.value = false }
}

const handleSearch = () => { pagination.page = 1; fetchList() }
const handleReset = () => { searchForm.username = ''; searchForm.ip = ''; searchForm.status = ''; searchForm.date_range = null; pagination.page = 1; fetchList() }
const handleSizeChange = (size: number) => { pagination.page_size = size; pagination.page = 1; fetchList() }
const handlePageChange = (page: number) => { pagination.page = page; fetchList() }

onMounted(() => { fetchList() })
</script>

<style scoped lang="scss">
.login-log-page { padding: 16px; }
.search-card { margin-bottom: 16px; :deep(.el-card__body) { padding-bottom: 0; } }
.table-card { :deep(.el-card__body) { padding: 0; } }
.pagination-wrapper { display: flex; justify-content: flex-end; padding: 16px; }
</style>
