<template>
  <div class="api-manage-page">
    <div class="page-header">
      <h1 class="page-title">API 管理</h1>
    </div>

    <!-- API功能开关 -->
    <el-card shadow="never" class="api-card">
      <template #header>
        <div class="card-header">
          <span>API 设置</span>
        </div>
      </template>

      <div class="api-toggle-section">
        <div class="toggle-row">
          <div class="toggle-info">
            <span class="toggle-label">API 功能</span>
            <span class="toggle-desc">开启后可使用 API 接口进行对接</span>
          </div>
          <el-switch
            v-model="apiOpen"
            :active-value="1"
            :inactive-value="0"
            active-text="开启"
            inactive-text="关闭"
            @change="handleToggleAPI"
            :loading="toggling"
          />
        </div>
      </div>
    </el-card>

    <!-- API密钥信息（仅开启时显示） -->
    <el-card v-if="summary" shadow="never" class="api-card">
      <template #header>
        <div class="card-header">
          <span>API 密钥</span>
          <el-button type="warning" size="small" @click="handleResetKey" :loading="resetting">
            <el-icon><RefreshRight /></el-icon>
            重置密钥
          </el-button>
        </div>
      </template>

      <div class="api-key-section">
        <el-alert
          title="请妥善保管您的 API 密钥，不要泄露给他人。如发现密钥泄露，请立即重置。"
          type="warning"
          :closable="false"
          show-icon
          class="security-alert"
        />

        <div class="key-display-row">
          <label class="key-label">API 密钥</label>
          <div class="key-copy-row">
            <code class="key-text">{{ summary.api_password || '未生成' }}</code>
            <el-button type="primary" size="small" @click="copyKey(summary.api_password)" :disabled="!summary.api_password">
              复制
            </el-button>
          </div>
        </div>

        <div v-if="summary.api_create_time" class="key-meta">
          <span>开启时间：{{ formatTime(summary.api_create_time) }}</span>
        </div>
      </div>
    </el-card>

    <!-- API使用统计（仅开启时显示） -->
    <el-card v-if="summary" shadow="never" class="api-card">
      <template #header>
        <div class="card-header">
          <span>使用统计</span>
        </div>
      </template>

      <div class="stats-grid">
        <div class="stat-item">
          <div class="stat-value">{{ summary.host_count || 0 }}</div>
          <div class="stat-label">主机总数</div>
        </div>
        <div class="stat-item">
          <div class="stat-value">{{ summary.active_count || 0 }}</div>
          <div class="stat-label">活跃主机</div>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { RefreshRight } from '@element-plus/icons-vue'
import request from '@/utils/request'

interface APISummary {
  api_password: string
  api_create_time: string | null
  host_count: number
  active_count: number
}

const apiOpen = ref<number>(0)
const summary = ref<APISummary | null>(null)
const toggling = ref(false)
const resetting = ref(false)

// 获取用户信息（包含api_open状态）
const fetchUserInfo = async () => {
  try {
    const { data } = await request.get('/api/v1/user/profile')
    if (data?.data) {
      apiOpen.value = data.data.api_open || 0
      if (apiOpen.value === 1) {
        fetchSummary()
      }
    }
  } catch (e) {
    console.error('Failed to fetch user info:', e)
  }
}

// 获取API摘要
const fetchSummary = async () => {
  try {
    const { data } = await request.get('/api/v1/user/api/summary')
    if (data?.data) {
      summary.value = data.data
    }
  } catch (e: any) {
    // 如果返回错误说明API未开启
    summary.value = null
  }
}

onMounted(() => {
  fetchUserInfo()
})

// 开关API
const handleToggleAPI = async (val: number) => {
  toggling.value = true
  try {
    const { data } = await request.post('/api/v1/user/api/open', { open: val })
    if (data?.code === 0 || data?.status === 200) {
      ElMessage.success(val === 1 ? 'API 已开启' : 'API 已关闭')
      if (val === 1) {
        fetchSummary()
      } else {
        summary.value = null
      }
    } else {
      // 回滚状态
      apiOpen.value = val === 1 ? 0 : 1
      ElMessage.error(data?.msg || data?.message || '操作失败')
    }
  } catch (e: any) {
    apiOpen.value = val === 1 ? 0 : 1
    ElMessage.error(e?.response?.data?.msg || e?.message || '操作失败')
  } finally {
    toggling.value = false
  }
}

// 重置密钥
const handleResetKey = async () => {
  try {
    await ElMessageBox.confirm(
      '重置密钥后，旧密钥将立即失效，确定要重置吗？',
      '确认重置',
      { confirmButtonText: '确定重置', cancelButtonText: '取消', type: 'warning' }
    )
  } catch {
    return // 用户取消
  }

  resetting.value = true
  try {
    const { data } = await request.post('/api/v1/user/api/reset')
    if (data?.code === 0 || data?.status === 200) {
      ElMessage.success('密钥已重置')
      fetchSummary()
    } else {
      ElMessage.error(data?.msg || data?.message || '重置失败')
    }
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.msg || e?.message || '重置失败')
  } finally {
    resetting.value = false
  }
}

// 复制密钥
const copyKey = (key: string) => {
  if (!key) return
  navigator.clipboard.writeText(key).then(() => {
    ElMessage.success('已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败，请手动复制')
  })
}

// 格式化时间
const formatTime = (t: string) => {
  if (!t) return ''
  return new Date(t).toLocaleString('zh-CN')
}
</script>

<style scoped lang="scss">
.api-manage-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
  max-width: 800px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.api-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-weight: 600;
}

.api-toggle-section {
  padding: 8px 0;
}

.toggle-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.toggle-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.toggle-label {
  font-size: 15px;
  font-weight: 500;
  color: #303133;
}

.toggle-desc {
  font-size: 13px;
  color: #909399;
}

.api-key-section {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.security-alert {
  border-radius: 8px;
}

.key-display-row {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.key-label {
  font-size: 13px;
  color: #909399;
  font-weight: 500;
}

.key-copy-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.key-text {
  flex: 1;
  font-family: 'Courier New', monospace;
  font-size: 14px;
  background: #f5f7fa;
  padding: 10px 14px;
  border-radius: 8px;
  border: 1px solid #dcdfe6;
  color: #303133;
  word-break: break-all;
  user-select: all;
}

.key-meta {
  font-size: 13px;
  color: #909399;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
}

.stat-item {
  text-align: center;
  padding: 20px;
  background: #f5f7fa;
  border-radius: 10px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #409eff;
  margin-bottom: 6px;
}

.stat-label {
  font-size: 13px;
  color: #909399;
}
</style>
