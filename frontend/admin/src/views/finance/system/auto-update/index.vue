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
              <el-tag :type="isStable ? 'success' : 'warning'" size="large">
                {{ isStable ? '稳定版' : '测试版' }}
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

      <!-- 有更新 -->
      <el-alert
        v-if="hasUpdate"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 20px"
      >
        <template #title>
          发现新版本 {{ remoteVersion }}，请前往 GitHub Release 下载
        </template>
        <template #default>
          <div style="margin-top: 8px">
            <el-button type="warning" size="small" @click="openGitHubRelease">
              下载最新版本
            </el-button>
            <el-button size="small" @click="openQQGroup">
              <el-icon><ChatDotRound /></el-icon> 加入交流群
            </el-button>
            <el-button size="small" @click="showChangelogDialog = true">
              <el-icon><Document /></el-icon> 查看更新日志
            </el-button>
          </div>
        </template>
      </el-alert>

      <!-- 已是最新 -->
      <el-alert
        v-else-if="checked"
        type="success"
        :closable="false"
        show-icon
        style="margin-bottom: 20px"
      >
        <template #title>
          当前已是最新版本 v{{ currentVersion }}
        </template>
      </el-alert>
    </el-card>

    <!-- 更新日志（始终显示） -->
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
          :hollow="log.version !== currentVersion"
        >
          <el-card shadow="hover">
            <h4 style="margin: 0 0 8px">
              v{{ log.version }}
              <el-tag
                :type="log.type === 'stable' ? 'success' : 'warning'"
                size="small"
                style="margin-left: 8px"
              >
                {{ log.type === 'stable' ? '稳定版' : '测试版' }}
              </el-tag>
              <el-tag
                v-if="log.version === currentVersion"
                type="primary"
                size="small"
                style="margin-left: 4px"
              >
                当前
              </el-tag>
            </h4>
            <ul style="margin: 0; padding-left: 20px">
              <li
                v-for="(change, idx) in log.changes"
                :key="idx"
                style="color: #666; line-height: 1.8"
              >
                {{ change }}
              </li>
            </ul>
          </el-card>
        </el-timeline-item>
      </el-timeline>
    </el-card>

    <!-- 更新日志弹窗 -->
    <el-dialog v-model="showChangelogDialog" title="更新日志" width="600px">
      <el-timeline v-if="changelog.length > 0">
        <el-timeline-item
          v-for="log in changelog"
          :key="log.version"
          :timestamp="log.date"
          placement="top"
        >
          <h4>
            v{{ log.version }}
            <el-tag
              v-if="log.version === currentVersion"
              type="primary"
              size="small"
              style="margin-left: 8px"
            >
              当前
            </el-tag>
          </h4>
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
import { Refresh, ChatDotRound, Document } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'

const VERSION_URL = 'https://raw.githubusercontent.com/YeXuanHs/anchor-finance/master/version.json'
const QQ_GROUP_URL = 'https://qm.qq.com/q/m3i0A7bwga'

const currentVersion = ref(__APP_VERSION__ || '1.0.0')
const remoteVersion = ref('')
const downloadUrl = ref('')
const changelog = ref<any[]>([])
const checking = ref(false)
const checked = ref(false)
const showChangelogDialog = ref(false)

const isStable = computed(() =>
  !currentVersion.value.includes('demo') && !currentVersion.value.includes('beta')
)

const hasUpdate = computed(() =>
  remoteVersion.value !== '' && remoteVersion.value !== currentVersion.value
)

const handleCheckUpdate = async () => {
  checking.value = true
  try {
    const res = await fetch(VERSION_URL + '?t=' + Date.now())
    if (!res.ok) throw new Error('Failed to fetch')
    const data = await res.json()

    remoteVersion.value = data.latest_version || ''
    downloadUrl.value = data.github_release || ''
    changelog.value = data.changelog || []
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

const openGitHubRelease = () => {
  if (downloadUrl.value) {
    window.open(downloadUrl.value, '_blank')
  } else {
    ElMessage.warning('暂无下载链接')
  }
}

onMounted(() => {
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
