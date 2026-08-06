<template>
  <div class="logs-page">
    <!-- 统计面板 -->
    <el-row :gutter="16" class="stats-row">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value total">{{ stats.total }}</div>
          <div class="stat-label">日志总量</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value info">{{ stats.info_count }}</div>
          <div class="stat-label">INFO 日志</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value warn">{{ stats.warn_count }}</div>
          <div class="stat-label">WARN 日志</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-value error">{{ stats.error_count }}</div>
          <div class="stat-label">ERROR 日志</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" class="stats-row">
      <el-col :span="16">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>模块统计</span>
            </div>
          </template>
          <div class="module-stats">
            <el-tag
              v-for="item in moduleStats"
              :key="item.module"
              class="module-tag"
              @click="filterByModule(item.module)"
            >
              {{ item.module }} ({{ item.count }})
            </el-tag>
            <el-empty v-if="!moduleStats.length" description="暂无数据" :image-size="60" />
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>快捷清理</span>
            </div>
          </template>
          <el-form label-width="80px" size="small">
            <el-form-item label="清理级别">
              <el-select v-model="cleanForm.level" placeholder="全部" clearable>
                <el-option label="INFO" value="info" />
                <el-option label="WARN" value="warn" />
                <el-option label="ERROR" value="error" />
                <el-option label="DEBUG" value="debug" />
              </el-select>
            </el-form-item>
            <el-form-item label="保留天数">
              <el-input-number v-model="cleanForm.days" :min="1" :max="365" />
            </el-form-item>
            <el-form-item>
              <el-button type="danger" @click="handleCleanByLevel">按条件清理</el-button>
            </el-form-item>
          </el-form>
        </el-card>
      </el-col>
    </el-row>
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>系统日志</span>
          <div class="header-actions">
            <el-button type="danger" size="small" @click="handleClearOld">清理旧日志</el-button>
            <el-button type="success" size="small" @click="handleExport">导出</el-button>
          </div>
        </div>
      </template>
      <el-form :inline="true" :model="searchForm" class="search-form">
        <el-form-item label="级别">
          <el-select v-model="searchForm.level" placeholder="全部" clearable>
            <el-option label="INFO" value="info" />
            <el-option label="WARN" value="warn" />
            <el-option label="ERROR" value="error" />
            <el-option label="DEBUG" value="debug" />
          </el-select>
        </el-form-item>
        <el-form-item label="模块">
          <el-input v-model="searchForm.module" placeholder="模块名" clearable />
        </el-form-item>
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="搜索关键词" clearable />
        </el-form-item>
        <el-form-item label="时间范围">
          <el-date-picker v-model="searchForm.date_range" type="daterange" range-separator="至" start-placeholder="开始" end-placeholder="结束" value-format="YYYY-MM-DD" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <el-table :data="tableData" v-loading="loading" style="width: 100%" border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="level" label="级别" width="80">
          <template #default="{ row }">
            <el-tag :type="levelTypeMap[row.level]" size="small">{{ row.level }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="module" label="模块" width="120" />
        <el-table-column prop="message" label="日志内容" min-width="300" show-overflow-tooltip />
        <el-table-column prop="ip" label="IP" width="130" />
        <el-table-column prop="created_at" label="时间" width="180" />
        <el-table-column label="操作" width="80" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleDetail(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>
      <div class="pagination-container">
        <el-pagination v-model:current-page="pagination.page" v-model:page-size="pagination.page_size"
          :page-sizes="[20, 50, 100, 200]" :total="pagination.total" layout="total, sizes, prev, pager, next, jumper"
          @size-change="handleSizeChange" @current-change="handlePageChange" />
      </div>
  </el-card>
    <el-dialog v-model="detailVisible" title="日志详情" width="700px">
      <el-descriptions :column="1" border v-if="detailData">
        <el-descriptions-item label="ID">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item label="级别"><el-tag :type="levelTypeMap[detailData.level]" size="small">{{ detailData.level }}</el-tag></el-descriptions-item>
        <el-descriptions-item label="模块">{{ detailData.module }}</el-descriptions-item>
        <el-descriptions-item label="消息">{{ detailData.message }}</el-descriptions-item>
        <el-descriptions-item label="上下文"><pre>{{ detailData.context }}</pre></el-descriptions-item>
        <el-descriptions-item label="IP">{{ detailData.ip }}</el-descriptions-item>
        <el-descriptions-item label="用户">{{ detailData.user_id || '-' }}</el-descriptions-item>
        <el-descriptions-item label="时间">{{ detailData.created_at }}</el-descriptions-item>
      </el-descriptions>
      <template #footer><el-button @click="detailVisible = false">关闭</el-button></template>
    </el-dialog>
  </div>
</template>
<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
const levelTypeMap: Record<string, string> = { info: 'primary', warn: 'warning', error: 'danger', debug: 'info' }
const loading = ref(false)
const detailVisible = ref(false)
const detailData = ref<any>(null)
const searchForm = reactive({ level: '', module: '', keyword: '', date_range: [] as string[] })
const pagination = reactive({ page: 1, page_size: 20, total: 0 })
const tableData = ref<any[]>([])
const stats = reactive({ total: 0, info_count: 0, warn_count: 0, error_count: 0 })
const moduleStats = ref<{ module: string; count: number }[]>([])
const cleanForm = reactive({ level: '', days: 30 })
const fetchStats = async () => {
  try {
    const data = await request.get({ url: '/api/admin/system-logs/stats' })
    if (data) Object.assign(stats, data)
  } catch { /* ignore */ }
}
const fetchModuleStats = async () => {
  try {
    const data = await request.get({ url: '/api/admin/system-logs/module-stats' })
    moduleStats.value = data || []
  } catch { /* ignore */ }
}
const fetchData = async () => {
  loading.value = true
  try {
    const params: any = { page: pagination.page, page_size: pagination.page_size, level: searchForm.level || undefined, module: searchForm.module || undefined, keyword: searchForm.keyword || undefined }
    if (searchForm.date_range?.length === 2) { params.start_date = searchForm.date_range[0]; params.end_date = searchForm.date_range[1] }
    const data = await request.get({ url: '/api/admin/system-logs', params })
    tableData.value = data.list || []; pagination.total = data.total || 0
  } catch { ElMessage.error('获取日志失败') } finally { loading.value = false }
}
const filterByModule = (module: string) => { searchForm.module = module; handleSearch() }
const handleSearch = () => { pagination.page = 1; fetchData() }
const handleReset = () => { Object.assign(searchForm, { level: '', module: '', keyword: '', date_range: [] }); handleSearch() }
const handleDetail = (row: any) => { detailData.value = row; detailVisible.value = true }
const handleSizeChange = () => { pagination.page = 1; fetchData() }
const handlePageChange = () => { fetchData() }

const handleClearOld = async () => {
  try {
    await ElMessageBox.confirm('确定清理30天前的日志吗？', '清理确认')
    await request.post({ url: '/api/admin/system-logs/clean', data: { days: 30 } })
    ElMessage.success('清理完成'); fetchData(); fetchStats(); fetchModuleStats()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error('清理失败') }
}
const handleCleanByLevel = async () => {
  try {
    await ElMessageBox.confirm(`确定清理${cleanForm.days}天前${cleanForm.level ? '的' + cleanForm.level.toUpperCase() + '级别' : ''}日志吗？`, '清理确认')
    await request.post({ url: '/api/admin/system-logs/clean', data: { days: cleanForm.days, level: cleanForm.level || undefined } })
    ElMessage.success('清理完成'); fetchData(); fetchStats(); fetchModuleStats()
  } catch (e: any) { if (e !== 'cancel') ElMessage.error('清理失败') }
}
const handleExport = async () => {
  try {
    const params: any = { level: searchForm.level, module: searchForm.module }
    if (searchForm.date_range?.length === 2) { params.start_date = searchForm.date_range[0]; params.end_date = searchForm.date_range[1] }
    const data = await request.get({ url: '/api/admin/system-logs/export', params })
    if (data?.url) { window.open(data.url); ElMessage.success('导出成功') }
    else ElMessage.success('导出任务已提交')
  } catch { ElMessage.error('导出失败') }
}
onMounted(() => { fetchData(); fetchStats(); fetchModuleStats() })
</script>
<style scoped lang="scss">
.logs-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.header-actions { display: flex; gap: 8px; }
.search-form { margin-bottom: 20px; }
.pagination-container { display: flex; justify-content: flex-end; margin-top: 20px; }
pre { margin: 0; white-space: pre-wrap; word-break: break-all; max-height: 300px; overflow: auto; }
.stats-row { margin-bottom: 16px; }
.stat-card { text-align: center; }
.stat-value { font-size: 24px; font-weight: 700; margin-bottom: 4px; }
.stat-value.total { color: var(--el-color-primary); }
.stat-value.info { color: var(--el-color-info); }
.stat-value.warn { color: var(--el-color-warning); }
.stat-value.error { color: var(--el-color-danger); }
.stat-label { color: var(--el-text-color-secondary); font-size: 13px; }
.module-stats { display: flex; flex-wrap: wrap; gap: 8px; }
.module-tag { cursor: pointer; transition: all .2s; &:hover { transform: scale(1.05); } }
</style>
