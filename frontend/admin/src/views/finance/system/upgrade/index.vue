<template>
  <div class="upgrade-page">
    <el-card shadow="never">
      <template #header><span>{{ $t('upgrade.title') }}</span></template>

      <div class="current-version">
        <el-descriptions :column="2" border :title="$t('upgrade.currentVersion')">
          <el-descriptions-item :label="$t('upgrade.currentVersion')">v{{ currentVersion }}</el-descriptions-item>
          <el-descriptions-item :label="$t('upgrade.lastCheck')">{{ lastCheck || '-' }}</el-descriptions-item>
        </el-descriptions>
      </div>

      <div class="upgrade-actions">
        <el-button type="primary" :loading="checking" @click="checkUpdate">
          <el-icon><Refresh /></el-icon>
          {{ $t('upgrade.checkUpdate') }}
        </el-button>
      </div>

      <div v-if="updateInfo.has_update" class="update-info">
        <el-alert :title="$t('upgrade.newVersionFound')" type="success" :closable="false" show-icon>
          <template #default>
            <p>{{ $t('upgrade.latestVersion') }}: <strong>v{{ updateInfo.latest_version }}</strong></p>
            <p>{{ $t('upgrade.releaseDate') }}: {{ updateInfo.release_date }}</p>
          </template>
        </el-alert>

        <el-card shadow="never" class="changelog-card">
          <template #header><span>{{ $t('upgrade.updateLog') }}</span></template>
          <div v-html="updateInfo.changelog"></div>
        </el-card>

        <div class="upgrade-btn">
          <el-button type="success" size="large" :loading="upgrading" @click="handleUpgrade">
            <el-icon><Upload /></el-icon>
            {{ $t('upgrade.upgradeNow') }}
          </el-button>
        </div>
      </div>

      <el-empty v-else-if="!checking && checked" :description="$t('upgrade.alreadyLatest')" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Upload } from '@element-plus/icons-vue'
import request from '@/utils/http'
import { $t } from '@/locales'

const currentVersion = ref('1.0.0')
const lastCheck = ref('')
const checking = ref(false)
const checked = ref(false)
const upgrading = ref(false)

const updateInfo = reactive({
  has_update: false,
  latest_version: '',
  release_date: '',
  changelog: ''
})

const fetchCurrentVersion = async () => {
  try {
    const data = await request.get({ url: '/api/admin/system/info' })
    currentVersion.value = data?.version || '1.0.0'
  } catch (error) {
    console.error('Failed to fetch version:', error)
  }
}

const checkUpdate = async () => {
  checking.value = true
  checked.value = false
  try {
    const data = await request.get({ url: '/api/admin/system/check-update' })
    Object.assign(updateInfo, data || {})
    lastCheck.value = new Date().toLocaleString()
    checked.value = true
    if (!data?.has_update) {
      ElMessage.success($t('upgrade.alreadyLatest'))
    }
  } catch (error) {
    console.error('Check update failed:', error)
  } finally {
    checking.value = false
  }
}

const handleUpgrade = async () => {
  try {
    await ElMessageBox.confirm($t('upgrade.confirmUpgrade'), $t('upgrade.confirmUpgradeTitle'), {
      type: 'warning',
      confirmButtonText: $t('upgrade.confirmUpgradeBtn'),
      cancelButtonText: $t('upgrade.cancel')
    })
    upgrading.value = true
    await request.post({ url: '/api/admin/system/upgrade' })
    ElMessage.success($t('upgrade.upgradeSuccess'))
    setTimeout(() => { window.location.reload() }, 3000)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Upgrade failed:', error)
      ElMessage.error($t('upgrade.upgradeFailed'))
    }
  } finally {
    upgrading.value = false
  }
}

onMounted(() => { fetchCurrentVersion() })
</script>

<style scoped lang="scss">
.upgrade-page { padding: 16px; }
.current-version { margin-bottom: 20px; }
.upgrade-actions { margin-bottom: 20px; }
.update-info { margin-top: 20px; }
.changelog-card { margin-top: 16px; }
.upgrade-btn { margin-top: 20px; text-align: center; }
</style>
