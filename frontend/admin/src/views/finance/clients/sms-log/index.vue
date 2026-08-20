<template>
  <div class="sms-log-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('clientsSmsLog.title') }}</span>
          <el-button type="success" @click="handleExport">
            <el-icon><Download /></el-icon>
            {{ $t('common.export') }}
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item :label="$t('common.phone')">
          <el-input v-model="searchForm.phone" :placeholder="$t('clientsSmsLog.enterPhone')" clearable />
        </el-form-item>
        <el-form-item :label="$t('common.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('clientsSmsLog.sendSuccess')" value="success" />
            <el-option :label="$t('clientsSmsLog.sendFailed')" value="failed" />
            <el-option :label="$t('clientsSmsLog.pending')" value="pending" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('clientsSmsLog.template')">
          <el-select v-model="searchForm.template" :placeholder="$t('common.all')" clearable>
            <el-option :label="$t('clientsSmsLog.verifyCode')" value="verify" />
            <el-option :label="$t('clientsSmsLog.notification')" value="notify" />
            <el-option :label="$t('clientsSmsLog.marketing')" value="marketing" />
            <el-option :label="$t('clientsSmsLog.orderReminder')" value="order" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('common.dateRange')">
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
          <el-button type="primary" @click="handleSearch">{{ $t('common.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('common.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="phone" :label="$t('common.phone')" width="140" />
        <el-table-column prop="content" :label="$t('clientsSmsLog.smsContent')" min-width="300" show-overflow-tooltip />
        <el-table-column prop="template_name" :label="$t('clientsSmsLog.template')" width="120" />
        <el-table-column prop="status" :label="$t('common.status')" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTypeMap[row.status]" size="small">
              {{ statusLabelMap[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sent_at" :label="$t('clientsSmsLog.sentAt')" width="180" />
        <el-table-column prop="error_msg" :label="$t('clientsSmsLog.errorMsg')" width="200" show-overflow-tooltip />
        <el-table-column :label="$t('common.action')" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleResend(row)" :disabled="row.status === 'success'">{{ $t('clientsSmsLog.resend') }}</el-button>
          </template>
        </el-table-column>
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download } from '@element-plus/icons-vue'
import request from '@/utils/http'
import { $t } from '@/locales'

const statusTypeMap: Record<string, any> = {
  success: 'success',
  failed: 'danger',
  pending: 'warning'
}

const statusLabelMap: Record<string, string> = {
  success: $t('clientsSmsLog.sendSuccess'),
  failed: $t('clientsSmsLog.sendFailed'),
  pending: $t('clientsSmsLog.pending')
}

const loading = ref(false)
const searchForm = reactive({
  phone: '',
  status: '',
  template: '',
  date_range: [] as string[]
})
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])

const fetchData = async () => {
  loading.value = true
  try {
    const params: any = {
      page: pagination.page,
      page_size: pagination.page_size,
      phone: searchForm.phone || undefined,
      status: searchForm.status || undefined,
      template: searchForm.template || undefined
    }
    if (searchForm.date_range?.length === 2) {
      params.start_date = searchForm.date_range[0]
      params.end_date = searchForm.date_range[1]
    }
    const data = await request.get({ url: '/api/admin/sms/logs', params })
    tableData.value = data.list || []
    pagination.total = data.total || 0
  } catch (error) {
    ElMessage.error($t('clientsSmsLog.fetchFailed'))
  } finally {
    loading.value = false
  }
}

const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => {
  searchForm.phone = ''
  searchForm.status = ''
  searchForm.template = ''
  searchForm.date_range = []
  handleSearch()
}

const handleResend = async (row: any) => {
  await ElMessageBox.confirm($t('clientsSmsLog.confirmResend', { phone: row.phone }), $t('common.tips'), { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/sms/send`, data: { phone: row.phone, content: row.content, template_id: row.template_id } })
    ElMessage.success($t('clientsSmsLog.resendSuccess'))
    fetchData()
  } catch (error) {
    ElMessage.error($t('clientsSmsLog.resendFailed'))
  }
}

const handleExport = async () => {
  try {
    const params: any = { ...searchForm }
    const res = await request.get({ url: '/api/admin/sms/logs', params, responseType: 'blob' as any })
    const blob = new Blob([res], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `sms_log_${new Date().toISOString().slice(0, 10)}.xlsx`
    link.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success($t('common.exportSuccess'))
  } catch (error) {
    ElMessage.error($t('common.exportFailed'))
  }
}
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

onMounted(() => { fetchData() })
</script>

<style scoped lang="scss">
.sms-log-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
</style>
