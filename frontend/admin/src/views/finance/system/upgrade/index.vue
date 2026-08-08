<template>
  <div class="upgrade-page">
    <el-card shadow="never">
      <template #header><span>系统升级</span></template>

      <div class="current-version">
        <el-descriptions :column="2" border title="当前版本">
          <el-descriptions-item label="当前版本">v{{ currentVersion }}</el-descriptions-item>
          <el-descriptions-item label="最后检查">{{ lastCheck || '-' }}</el-descriptions-item>
        </el-descriptions>
      </div>

      <div class="upgrade-actions">
        <el-button type="primary" :loading="checking" @click="checkUpdate">
          <el-icon><Refresh /></el-icon>
          检查更新
        </el-button>
      </div>

      <div v-if="updateInfo.has_update" class="update-info">
        <el-alert title="发现新版本" type="success" :closable="false" show-icon>
          <template #default>
            <p>最新版本: <strong>v{{ updateInfo.latest_version }}</strong></p>
            <p>发布日期: {{ updateInfo.release_date }}</p>
          </template>
        </el-alert>

        <el-card shadow="never" class="changelog-card">
          <template #header><span>更新日志</span></template>
          <div v-html="updateInfo.changelog"></div>
        </el-card>

        <div class="upgrade-btn">
          <el-button type="success" size="large" :loading="upgrading" @click="handleUpgrade">
            <el-icon><Upload /></el-icon>
            立即升级
          </el-button>
        </div>
      </div>

      <el-empty v-else-if="!checking && checked" description="当前已是最新版本" />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Upload } from '@element-plus/icons-vue'
import request from '@/utils/http'

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
    console.error('获取版本失败:', error)
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
      ElMessage.success('当前已是最新版本')
    }
  } catch (error) {
    console.error('检查更新失败:', error)
  } finally {
    checking.value = false
  }
}

const handleUpgrade = async () => {
  try {
    await ElMessageBox.confirm('升级前建议备份数据库。确定要开始升级吗？', '确认升级', {
      type: 'warning',
      confirmButtonText: '确定升级',
      cancelButtonText: '取消'
    })
    upgrading.value = true
    await request.post({ url: '/api/admin/system/upgrade' })
    ElMessage.success('升级成功，系统将自动重启')
    setTimeout(() => { window.location.reload() }, 3000)
  } catch (error) {
    if (error !== 'cancel') {
      console.error('升级失败:', error)
      ElMessage.error('升级失败，请查看日志')
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
