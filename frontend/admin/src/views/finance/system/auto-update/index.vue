<template>
  <div class="auto-update-page">
    <!-- 版本信息卡片 -->
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>系统更新管理</span>
          <el-button type="primary" @click="handleCheckUpdate" :loading="checking" :icon="Refresh">
            检查更新
          </el-button>
        </div>
      </template>

      <!-- 当前版本信息 -->
      <el-row :gutter="20" class="stat-section">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">
              <el-tag type="primary" size="large">v{{ currentVersion }}</el-tag>
            </div>
            <div class="stat-label">当前版本</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">
              <el-tag :type="remoteVersion ? 'success' : 'info'" size="large">
                {{ remoteVersion || '未检测' }}
              </el-tag>
            </div>
            <div class="stat-label">最新版本</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">
              <el-tag :type="versionType === 'stable' ? 'success' : 'warning'" size="large">
                {{ currentVersion.includes('demo') || currentVersion.includes('beta') ? '测试版' : '稳定版' }}
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

      <!-- 更新提示 -->
      <el-alert
        v-if="hasUpdate"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 20px"
      >
        <template #title>
          发现新版本 {{ remoteVersion }}，请加入交流群下载最新安装包
        </template>
        <template #default>
          <div style="margin-top: 8px">
            <el-button type="primary" size="small" @click="openQQGroup">
              加入交流群下载
            </el-button>
            <el-button size="small" @click="showChangelog = true">
              查看更新日志
            </el-button>
          </div>
        </template>
      </el-alert>

      <el-alert
        v-else-if="checked && !hasUpdate"
        type="success"
        :closable="false"
        show-icon
      >
        <template #title>
          当前已是最新版本 v{{ currentVersion }}
        </template>
      </el-alert>
    </el-card>

    <!-- 更新日志 -->
    <el-card shadow="never" class="section-card" v-if="changelog.length > 0">
      <template #header>
        <div class="card-header">
          <span>更新日志</span>
        </div>
      </template>
      <el-timeline>
        <el-timeline-item
          v-for="log in changelog"
          :key="log.version"
          :timestamp="log.date"
          placement="top"
          :type="log.type === 'stable' ? 'success' : 'warning'"
        >
          <el-card shadow="hover">
            <h4 style="margin: 0 0 8px">
              v{{ log.version }}
              <el-tag :type="log.type === 'stable' ? 'success' : 'warning'" size="small" style="margin-left: 8px">
                {{ log.type === 'stable' ? '稳定版' : '测试版' }}
              </el-tag>
            </h4>
            <ul style="margin: 0; padding-left: 20px">
              <li v-for="(change, idx) in log.changes" :key="idx" style="color: #666; line-height: 1.8">
                {{ change }}
              </li>
            </ul>
          </el-card>
        </el-timeline-item>
      </el-timeline>
    </el-card>

    <!-- 更新日志弹窗 -->
    <el-dialog v-model="showChangelog" title="更新日志" width="600px">
      <el-timeline v-if="changelog.length > 0">
        <el-timeline-item
          v-for="log in changelog"
          :key="log.version"
          :timestamp="log.date"
          placement="top"
        >
          <h4>v{{ log.version }}</h4>
          <ul>
            <li v-for="(change, idx) in log.changes" :key="idx">{{ change }}</li>
          </ul>
        </el-timeline-item>
      </el-timeline>
      <el-empty v-else description="暂无更新日志" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

// GitHub raw URL for version.json
const VERSION_URL = 'https://raw.githubusercontent.com/YeXuanHs/anchor-finance/master/version.json'
const QQ_GROUP_URL = 'https://qm.qq.com/q/m3i0A7bwga'

const currentVersion = ref(__APP_VERSION__ || '1.0.0')
const remoteVersion = ref('')
const versionType = ref('')
const changelog = ref<any[]>([])
const downloadTip = ref('')
const checking = ref(false)
const checked = ref(false)
const showChangelog = ref(false)

const hasUpdate = computed(() => {
  if (!remoteVersion.value) return false
  return remoteVersion.value !== currentVersion.value
})

const handleCheckUpdate = async () => {
  checking.value = true
  try {
    const res = await fetch(VERSION_URL + '?t=' + Date.now())
    if (!res.ok) throw new Error('Failed to fetch')
    const data = await res.json()

    remoteVersion.value = data.latest_version || ''
    changelog.value = data.changelog || []
    downloadTip.value = data.download_tip || '请加入交流群下载最新安装包'
    checked.value = true

    if (hasUpdate.value) {
      ElMessage.success(`发现新版本 ${remoteVersion.value}`)
    } else {
      ElMessage.info('当前已是最新版本')
    }
  } catch {
    ElMessage.error('检查更新失败，请检查网络连接')
    checked.value = true
  } finally {
    checking.value = false
  }
}

const openQQGroup = () => {
  window.open(QQ_GROUP_URL, '_blank')
}

onMounted(() => {
  // 自动检查一次
  handleCheckUpdate()
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

.section-card {
  margin-top: 20px;
}
</style>
