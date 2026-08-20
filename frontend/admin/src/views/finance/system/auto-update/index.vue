<template>
  <div class="auto-update-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>{{ $t('autoUpdate.title') }}</span>
          <el-button type="primary" @click="handleCheckUpdate" :loading="checking" :icon="Refresh">
            {{ $t('autoUpdate.checkUpdate') }}
          </el-button>
        </div>
      </template>

      <el-row :gutter="20" class="stat-section">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">
              <el-tag type="primary" size="large">v{{ currentVersion }}</el-tag>
            </div>
            <div class="stat-label">{{ $t('autoUpdate.currentVersion') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">
              <el-tag :type="remoteVersion ? 'success' : 'info'" size="large">
                {{ remoteVersion || $t('autoUpdate.notChecked') }}
              </el-tag>
            </div>
            <div class="stat-label">{{ $t('autoUpdate.latestVersion') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">
              <el-tag :type="isStable ? 'success' : 'warning'" size="large">
                {{ isStable ? $t('autoUpdate.stable') : $t('autoUpdate.beta') }}
              </el-tag>
            </div>
            <div class="stat-label">{{ $t('autoUpdate.versionType') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">
              <el-tag :type="hasUpdate ? 'warning' : 'success'" size="large">
                {{ hasUpdate ? $t('autoUpdate.hasUpdate') : $t('autoUpdate.isLatest') }}
              </el-tag>
            </div>
            <div class="stat-label">{{ $t('autoUpdate.updateStatus') }}</div>
          </el-card>
        </el-col>
      </el-row>

      <el-alert
        v-if="hasUpdate"
        type="warning"
        :closable="false"
        show-icon
        style="margin-bottom: 20px"
      >
        <template #title>
          {{ $t('autoUpdate.newVersionFound', { version: remoteVersion }) }}
        </template>
        <template #default>
          <div style="margin-top: 8px">
            <el-button type="warning" size="small" @click="openGitHubRelease">
              {{ $t('autoUpdate.downloadLatest') }}
            </el-button>
            <el-button size="small" @click="openQQGroup">
              <el-icon><ChatDotRound /></el-icon> {{ $t('autoUpdate.joinGroup') }}
            </el-button>
            <el-button size="small" @click="showChangelogDialog = true">
              <el-icon><Document /></el-icon> {{ $t('autoUpdate.viewChangelog') }}
            </el-button>
          </div>
        </template>
      </el-alert>

      <el-alert
        v-else-if="checked"
        type="success"
        :closable="false"
        show-icon
        style="margin-bottom: 20px"
      >
        <template #title>
          {{ $t('autoUpdate.currentIsLatest', { version: currentVersion }) }}
        </template>
      </el-alert>
    </el-card>

    <el-card shadow="never" class="section-card" v-if="changelog.length > 0">
      <template #header>
        <div class="card-header">
          <span>{{ $t('autoUpdate.updateLog') }}</span>
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
                {{ log.type === 'stable' ? $t('autoUpdate.stable') : $t('autoUpdate.beta') }}
              </el-tag>
              <el-tag
                v-if="log.version === currentVersion"
                type="primary"
                size="small"
                style="margin-left: 4px"
              >
                {{ $t('autoUpdate.current') }}
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

    <el-dialog v-model="showChangelogDialog" :title="$t('autoUpdate.updateLog')" width="600px">
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
              {{ $t('autoUpdate.current') }}
            </el-tag>
          </h4>
          <ul>
            <li v-for="(change, idx) in log.changes" :key="idx">{{ change }}</li>
          </ul>
        </el-timeline-item>
      </el-timeline>
      <el-empty v-else :description="$t('autoUpdate.noChangelog')" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Refresh, ChatDotRound, Document } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { $t } from '@/locales'

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
      ElMessage.success($t('autoUpdate.updateFoundMsg', { version: remoteVersion.value }))
    } else {
      ElMessage.info($t('autoUpdate.isLatestVersion'))
    }
  } catch {
    ElMessage.error($t('autoUpdate.checkUpdateFailed'))
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
    ElMessage.warning($t('autoUpdate.noDownloadLink'))
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
