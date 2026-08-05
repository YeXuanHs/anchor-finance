<template>
  <div class="auto-update-page">
    <!-- 版本信息卡片 -->
    <el-card shadow="never" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>系统更新管理</span>
          <el-button type="primary" @click="handleCheckUpdate" :loading="checkingLoading" :icon="Refresh">
            检查更新
          </el-button>
        </div>
      </template>

      <!-- 当前版本信息 -->
      <el-row :gutter="20" class="stat-section">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">
              <el-tag type="primary" size="large">{{ systemInfo.install_version || '-' }}</el-tag>
            </div>
            <div class="stat-label">当前版本</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">
              <el-tag :type="versionInfo.last_version_check === 'no_response' ? 'info' : 'success'" size="large">
                {{ typeof versionInfo.last_version === 'string' ? versionInfo.last_version : versionInfo.last_version?.last || '检测中...' }}
              </el-tag>
            </div>
            <div class="stat-label">最新版本</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">
              <el-tag :type="systemInfo.system_version_type === 'stable' ? 'success' : 'warning'" size="large">
                {{ systemInfo.system_version_type === 'stable' ? '稳定版' : '测试版' }}
              </el-tag>
            </div>
            <div class="stat-label">版本类型</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">
              <el-tag :type="hasUpdate ? 'warning' : 'success'" size="large">
                {{ hasUpdate ? '有可用更新' : '已是最新' }}
              </el-tag>
            </div>
            <div class="stat-label">更新状态</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 更新操作 -->
      <div class="update-actions">
        <el-button
          type="primary"
          size="large"
          :disabled="!hasUpdate || updating"
          :loading="updating"
          @click="handleStartUpdate"
        >
          {{ updating ? '正在更新...' : '立即更新' }}
        </el-button>
        <el-button size="large" @click="handleViewChangelog">查看更新日志</el-button>
      </div>

      <!-- 更新进度 -->
      <div v-if="updating" class="update-progress">
        <el-divider />
        <div class="progress-header">
          <h3>更新进度</h3>
          <el-tag :type="getProgressStatusType(updateProgress.status)">
            {{ getProgressStatusLabel(updateProgress.status) }}
          </el-tag>
        </div>
        <el-progress
          :percentage="parseProgress(updateProgress.progress)"
          :status="updateProgress.status === 400 ? 'exception' : undefined"
          :stroke-width="20"
          style="margin: 16px 0"
        />
        <div class="progress-info">
          <el-text>{{ updateProgress.msg || '正在处理...' }}</el-text>
        </div>
      </div>
    </el-card>

    <!-- 自动更新配置 -->
    <el-card shadow="never" class="section-card">
      <template #header>
        <div class="card-header">
          <span>自动更新配置</span>
        </div>
      </template>
      <el-form :model="autoUpdateConfig" label-width="140px">
        <el-form-item label="启用自动更新">
          <el-switch v-model="autoUpdateConfig.enabled" />
          <span class="form-tip">启用后系统将在指定时间自动检查并安装更新</span>
        </el-form-item>
        <el-form-item label="更新通道">
          <el-select v-model="autoUpdateConfig.channel" style="width: 300px">
            <el-option label="稳定版 (Stable)" value="stable" />
            <el-option label="测试版 (Beta)" value="beta" />
          </el-select>
        </el-form-item>
        <el-form-item label="自动更新时间">
          <el-time-picker v-model="autoUpdateConfig.time" placeholder="选择更新时间" format="HH:mm" />
          <span class="form-tip">建议设置在访问量较低的时段</span>
        </el-form-item>
        <el-form-item label="更新前自动备份">
          <el-switch v-model="autoUpdateConfig.backup_before_update" />
          <span class="form-tip">更新前自动创建系统备份，建议保持开启</span>
        </el-form-item>
        <el-form-item label="邮件通知">
          <el-switch v-model="autoUpdateConfig.email_notify" />
          <span class="form-tip">更新完成后发送通知邮件给管理员</span>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSaveAutoUpdateConfig" :loading="configLoading">保存配置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 更新日志对话框 -->
    <el-dialog v-model="changelogDialogVisible" title="更新日志" width="700px">
      <div v-loading="changelogLoading" class="changelog-content">
        <div v-if="changelogContent" v-html="changelogContent" />
        <el-empty v-else description="暂无更新日志" />
      </div>
      <template #footer>
        <el-button @click="changelogDialogVisible = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const checkingLoading = ref(false)
const updating = ref(false)
const configLoading = ref(false)
const changelogDialogVisible = ref(false)
const changelogLoading = ref(false)

let progressTimer: ReturnType<typeof setInterval> | null = null

const systemInfo = reactive({
  install_version: '',
  system_version_type: '',
  auth_status: '',
  auth_due_time: '',
  service_due_time: ''
})

const versionInfo = reactive({
  last_version: '' as any,
  last_version_check: '',
  license_type: 0
})

const updateProgress = reactive({
  progress: '0%',
  msg: '',
  status: 200
})

const autoUpdateConfig = reactive({
  enabled: false,
  channel: 'stable',
  time: new Date(2025, 0, 1, 3, 0),
  backup_before_update: true,
  email_notify: false
})

const changelogContent = ref('')

const hasUpdate = computed(() => {
  if (!systemInfo.install_version || !versionInfo.last_version) return false
  if (typeof versionInfo.last_version === 'string') return false
  return versionInfo.last_version?.last !== systemInfo.install_version
})

const parseProgress = (progress: string) => {
  return parseInt(progress) || 0
}

const getProgressStatusType = (status: number) => {
  if (status === 200) return 'primary'
  if (status === 400) return 'danger'
  return 'info'
}

const getProgressStatusLabel = (status: number) => {
  if (status === 200) return '更新中'
  if (status === 400) return '更新失败'
  return '等待中'
}

const fetchSystemInfo = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/system/info' })
    if (res) Object.assign(systemInfo, res)
  } catch {
    ElMessage.error('获取系统信息失败')
  } finally {
    loading.value = false
  }
}

const fetchLastVersion = async () => {
  try {
    const res = await request.get({ url: '/api/admin/system/updates/last-version' })
    if (res) Object.assign(versionInfo, res)
  } catch { /* ignore */ }
}

const fetchAutoUpdateConfig = async () => {
  try {
    const res = await request.get({ url: '/api/admin/system/updates/auto-config' })
    if (res) {
      autoUpdateConfig.enabled = res.enabled ?? false
      autoUpdateConfig.channel = res.channel ?? 'stable'
      autoUpdateConfig.backup_before_update = res.backup_before_update ?? true
      autoUpdateConfig.email_notify = res.email_notify ?? false
      if (res.time) {
        const [h, m] = res.time.split(':')
        autoUpdateConfig.time = new Date(2025, 0, 1, Number(h), Number(m))
      }
    }
  } catch { /* ignore */ }
}

const handleCheckUpdate = async () => {
  checkingLoading.value = true
  try {
    await fetchLastVersion()
    if (hasUpdate.value) {
      ElMessage.success('发现新版本')
    } else {
      ElMessage.info('当前已是最新版本')
    }
  } catch {
    ElMessage.error('检查更新失败')
  } finally {
    checkingLoading.value = false
  }
}

const handleStartUpdate = async () => {
  try {
    await ElMessageBox.confirm(
      '系统更新前请确保已做好备份。更新过程中请勿关闭页面。确定开始更新？',
      '系统更新确认',
      { type: 'warning' }
    )
    updating.value = true
    updateProgress.progress = '0%'
    updateProgress.msg = '开始更新...'
    updateProgress.status = 200

    await request.get({ url: '/api/admin/system/updates/auto-update' })
    startProgressPolling()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error('更新启动失败')
      updating.value = false
    }
  }
}

const startProgressPolling = () => {
  if (progressTimer) clearInterval(progressTimer)
  progressTimer = setInterval(async () => {
    try {
      const res = await request.get({ url: '/api/admin/system/updates/check-progress' })
      if (res) {
        updateProgress.progress = res.progress || '0%'
        updateProgress.msg = res.msg || ''
        updateProgress.status = res.status || 200

        if (res.progress === '100%' || res.status === 400) {
          stopProgressPolling()
          if (res.status === 200) {
            ElMessage.success('系统更新完成')
            fetchSystemInfo()
            fetchLastVersion()
          } else {
            ElMessage.error(res.msg || '更新失败')
          }
          updating.value = false
        }
      }
    } catch {
      stopProgressPolling()
      updating.value = false
    }
  }, 3000)
}

const stopProgressPolling = () => {
  if (progressTimer) {
    clearInterval(progressTimer)
    progressTimer = null
  }
}

const handleViewChangelog = async () => {
  changelogDialogVisible.value = true
  changelogLoading.value = true
  try {
    const res = await request.get({ url: '/api/admin/system/updates/content' })
    changelogContent.value = res || '暂无更新日志'
  } catch {
    changelogContent.value = '获取更新日志失败'
  } finally {
    changelogLoading.value = false
  }
}

const handleSaveAutoUpdateConfig = async () => {
  configLoading.value = true
  try {
    const timeStr = autoUpdateConfig.time
      ? `${String(autoUpdateConfig.time.getHours()).padStart(2, '0')}:${String(autoUpdateConfig.time.getMinutes()).padStart(2, '0')}`
      : '03:00'
    await request.put({
      url: '/api/admin/system/updates/auto-config',
      data: { ...autoUpdateConfig, time: timeStr },
      showSuccessMessage: true
    })
  } catch {
    ElMessage.error('保存配置失败')
  } finally {
    configLoading.value = false
  }
}

onMounted(() => {
  fetchSystemInfo()
  fetchLastVersion()
  fetchAutoUpdateConfig()
})

onUnmounted(() => {
  stopProgressPolling()
})
</script>

<style scoped lang="scss">
.auto-update-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-section {
  margin-bottom: 24px;
}

.stat-card {
  text-align: center;
  padding: 16px 0;
}

.stat-value {
  font-size: 18px;
  font-weight: 600;
  color: var(--el-color-primary);
  margin-bottom: 8px;
}

.stat-label {
  color: var(--el-text-color-secondary);
  font-size: 14px;
}

.update-actions {
  display: flex;
  gap: 16px;
  margin: 24px 0;
}

.update-progress {
  margin-top: 16px;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;

  h3 {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
  }
}

.progress-info {
  margin-top: 8px;
}

.section-card {
  margin-top: 20px;
}

.form-tip {
  margin-left: 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.changelog-content {
  min-height: 200px;
  max-height: 500px;
  overflow-y: auto;

  :deep(h1) {
    font-size: 20px;
    margin: 16px 0 8px;
  }

  :deep(h2) {
    font-size: 16px;
    margin: 12px 0 6px;
  }

  :deep(ul) {
    padding-left: 20px;
  }
}
</style>
