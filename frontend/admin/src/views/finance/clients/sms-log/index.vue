<template>
  <div class="sms-log-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>客户短信日志</span>
          <el-button type="success" @click="handleExport">
            <el-icon><Download /></el-icon>
            导出
          </el-button>
        </div>
      </template>

      <!-- 搜索栏 -->
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="手机号">
          <el-input v-model="searchForm.phone" placeholder="请输入手机号" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="发送成功" value="success" />
            <el-option label="发送失败" value="failed" />
            <el-option label="待发送" value="pending" />
          </el-select>
        </el-form-item>
        <el-form-item label="模板">
          <el-select v-model="searchForm.template" placeholder="全部" clearable>
            <el-option label="验证码" value="verify" />
            <el-option label="通知" value="notify" />
            <el-option label="营销" value="marketing" />
            <el-option label="订单提醒" value="order" />
          </el-select>
        </el-form-item>
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="searchForm.date_range"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="phone" label="手机号" width="140" />
        <el-table-column prop="content" label="短信内容" min-width="300" show-overflow-tooltip />
        <el-table-column prop="template_name" label="模板" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTypeMap[row.status]" size="small">
              {{ statusLabelMap[row.status] }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="sent_at" label="发送时间" width="180" />
        <el-table-column prop="error_msg" label="错误信息" width="200" show-overflow-tooltip />
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleResend(row)" :disabled="row.status === 'success'">重发</el-button>
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

const statusTypeMap: Record<string, any> = {
  success: 'success',
  failed: 'danger',
  pending: 'warning'
}

const statusLabelMap: Record<string, string> = {
  success: '发送成功',
  failed: '发送失败',
  pending: '待发送'
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
    ElMessage.error('获取短信日志失败')
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
  await ElMessageBox.confirm(`确定重发短信到 ${row.phone}？`, '提示', { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/sms-logs/${row.id}/resend` })
    ElMessage.success('重发成功')
    fetchData()
  } catch (error) {
    ElMessage.error('重发失败')
  }
}

const handleExport = async () => {
  try {
    const params: any = { ...searchForm }
    const res = await request.get({ url: '/api/admin/sms-logs/export', params, responseType: 'blob' as any })
    const blob = new Blob([res], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `短信日志_${new Date().toISOString().slice(0, 10)}.xlsx`
    link.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
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
