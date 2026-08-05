<template>
  <div class="cache-clear-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>缓存清理管理</span>
        </div>
      </template>

      <!-- 缓存状态概览 -->
      <div class="section">
        <h3>缓存状态概览</h3>
        <el-row :gutter="16">
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-label">总缓存数</div>
              <div class="stat-value">{{ stats.total_items || 0 }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-label">缓存大小</div>
              <div class="stat-value primary">{{ stats.total_size || '0 MB' }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-label">命中率</div>
              <div class="stat-value success">{{ stats.hit_rate || '0%' }}</div>
            </el-card>
          </el-col>
          <el-col :span="6">
            <el-card shadow="hover" class="stat-card">
              <div class="stat-label">上次清理</div>
              <div class="stat-value warning" style="font-size: 16px">{{ stats.last_clear_at || '-' }}</div>
            </el-card>
          </el-col>
        </el-row>
      </div>

      <!-- 缓存列表 -->
      <div class="section">
        <div class="section-header">
          <h3>缓存类型列表</h3>
          <el-button type="danger" @click="handleClearAll" :loading="clearAllLoading">
            <el-icon><Delete /></el-icon>
            清理全部缓存
          </el-button>
        </div>
        <el-table :data="cacheList" v-loading="loading" style="width: 100%" border>
          <el-table-column prop="name" label="缓存名称" min-width="150" />
          <el-table-column prop="key" label="缓存标识" width="180">
            <template #default="{ row }">
              <el-tag size="small">{{ row.key }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="description" label="描述" min-width="200" show-overflow-tooltip />
          <el-table-column prop="items" label="缓存条数" width="100" align="center" />
          <el-table-column prop="size" label="缓存大小" width="100" align="center" />
          <el-table-column prop="last_cleared_at" label="上次清理" width="180" />
          <el-table-column label="操作" width="120" fixed="right">
            <template #default="{ row }">
              <el-popconfirm :title="`确定清理 ${row.name} 吗？`" @confirm="handleClearSingle(row)">
                <template #reference>
                  <el-button type="danger" link :loading="row.clearing">清理</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>

      <!-- 清理结果 -->
      <div class="section" v-if="clearResult">
        <h3>清理结果</h3>
        <el-card shadow="hover" class="result-card">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="清理时间">{{ clearResult.cleared_at }}</el-descriptions-item>
            <el-descriptions-item label="清理类型">
              <el-tag :type="clearResult.type === 'all' ? 'danger' : 'warning'">
                {{ clearResult.type === 'all' ? '全部清理' : '单类清理' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="清理缓存">{{ clearResult.cache_name || '全部' }}</el-descriptions-item>
            <el-descriptions-item label="清理条数">{{ clearResult.cleared_items || 0 }}</el-descriptions-item>
            <el-descriptions-item label="释放空间">{{ clearResult.freed_size || '0 MB' }}</el-descriptions-item>
            <el-descriptions-item label="耗时">{{ clearResult.duration || '0ms' }}</el-descriptions-item>
          </el-descriptions>
          <div class="result-status" :class="clearResult.success ? 'success' : 'error'">
            <el-icon v-if="clearResult.success"><CircleCheckFilled /></el-icon>
            <el-icon v-else><CircleCloseFilled /></el-icon>
            {{ clearResult.success ? '清理成功' : '清理失败' }}
            <span v-if="clearResult.error_msg" class="error-msg">：{{ clearResult.error_msg }}</span>
          </div>
        </el-card>
      </div>

      <!-- 操作日志 -->
      <div class="section">
        <div class="section-header">
          <h3>清理操作日志</h3>
          <el-button size="small" @click="fetchLogs">刷新</el-button>
        </div>
        <el-table :data="logs" v-loading="logsLoading" style="width: 100%" border size="small">
          <el-table-column prop="id" label="ID" width="70" />
          <el-table-column prop="operator" label="操作人" width="120" />
          <el-table-column prop="cache_name" label="清理目标" width="150" />
          <el-table-column prop="cleared_items" label="清理条数" width="100" align="center" />
          <el-table-column prop="freed_size" label="释放空间" width="100" align="center" />
          <el-table-column prop="status" label="状态" width="80" align="center">
            <template #default="{ row }">
              <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">
                {{ row.status === 'success' ? '成功' : '失败' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="created_at" label="操作时间" width="180" />
          <el-table-column prop="remark" label="备注" min-width="150" show-overflow-tooltip />
        </el-table>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Delete, CircleCheckFilled, CircleCloseFilled } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'CacheClearManage' })

const loading = ref(false)
const clearAllLoading = ref(false)
const logsLoading = ref(false)

const stats = reactive({ total_items: 0, total_size: '0 MB', hit_rate: '0%', last_clear_at: '' })
const cacheList = ref<any[]>([])
const logs = ref<any[]>([])
const clearResult = ref<any>(null)

const fetchStats = async () => {
  try {
    const data = await request.get({ url: '/api/admin/system/clear-cache/stats' })
    Object.assign(stats, data)
  } catch (error) {
    console.error('获取缓存统计失败:', error)
  }
}

const fetchCacheList = async () => {
  loading.value = true
  try {
    const data = await request.get({ url: '/api/admin/system/clear-cache/list' })
    cacheList.value = (data || []).map((item: any) => ({ ...item, clearing: false }))
  } catch (error) {
    ElMessage.error('获取缓存列表失败')
  } finally {
    loading.value = false
  }
}

const fetchLogs = async () => {
  logsLoading.value = true
  try {
    const data = await request.get({ url: '/api/admin/system/clear-cache/logs' })
    logs.value = data || []
  } catch (error) {
    console.error('获取操作日志失败:', error)
  } finally {
    logsLoading.value = false
  }
}

const handleClearSingle = async (row: any) => {
  row.clearing = true
  try {
    const data = await request.post({ url: '/api/admin/system/clear-cache', params: { cache_key: row.key } })
    ElMessage.success(`${row.name} 清理成功`)
    clearResult.value = {
      ...data,
      cache_name: row.name,
      type: 'single',
      success: true,
      cleared_at: new Date().toLocaleString()
    }
    fetchCacheList()
    fetchStats()
    fetchLogs()
  } catch (error: any) {
    clearResult.value = {
      cache_name: row.name,
      type: 'single',
      success: false,
      error_msg: error.message || '清理失败',
      cleared_at: new Date().toLocaleString()
    }
    ElMessage.error('清理失败')
  } finally {
    row.clearing = false
  }
}

const handleClearAll = async () => {
  try {
    await ElMessageBox.confirm('确定清理全部缓存吗？此操作不可撤销。', '清理全部缓存', { type: 'warning' })
    clearAllLoading.value = true
    const data = await request.post({ url: '/api/admin/system/clear-cache', params: { cache_key: 'all' } })
    ElMessage.success('全部缓存清理成功')
    clearResult.value = {
      ...data,
      type: 'all',
      success: true,
      cleared_at: new Date().toLocaleString()
    }
    fetchCacheList()
    fetchStats()
    fetchLogs()
  } catch (error: any) {
    if (error !== 'cancel') {
      clearResult.value = {
        type: 'all',
        success: false,
        error_msg: error.message || '清理失败',
        cleared_at: new Date().toLocaleString()
      }
      ElMessage.error('清理失败')
    }
  } finally {
    clearAllLoading.value = false
  }
}

onMounted(() => { fetchStats(); fetchCacheList(); fetchLogs() })
</script>

<style scoped lang="scss">
.cache-clear-page { padding: 20px; }
.card-header { display: flex; justify-content: space-between; align-items: center; }
.section { margin-top: 24px; &:first-child { margin-top: 0; } }
.section-header {
  display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;
  h3 { margin: 0; font-size: 16px; font-weight: 600; }
}
h3 { margin: 0 0 16px; font-size: 16px; font-weight: 600; }
.stat-card { text-align: center; padding: 8px 0; }
.stat-value {
  font-size: 24px; font-weight: 600; color: var(--el-text-color-primary); margin-bottom: 4px;
  &.success { color: var(--el-color-success); }
  &.primary { color: var(--el-color-primary); }
  &.warning { color: var(--el-color-warning); }
}
.stat-label { color: var(--el-text-color-secondary); font-size: 14px; }
.result-card { margin-top: 8px; }
.result-status {
  margin-top: 16px; padding: 12px; border-radius: 4px; display: flex; align-items: center; gap: 8px;
  &.success { background: var(--el-color-success-light-9); color: var(--el-color-success); }
  &.error { background: var(--el-color-danger-light-9); color: var(--el-color-danger); }
}
.error-msg { font-size: 13px; }
</style>
