<template>
  <div class="authorize-page">
    <!-- 授权状态卡片 -->
    <el-card shadow="never" v-loading="loading" class="status-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('page.authorize.title') }}</span>
          <div class="header-actions">
            <el-button type="primary" @click="handleRefresh" :icon="Refresh">{{ $t('page.authorize.refreshStatus') }}</el-button>
            <el-button type="warning" @click="handleVerify">{{ $t('page.authorize.verifyAuth') }}</el-button>
          </div>
        </div>
      </template>

      <!-- 授权状态概览 -->
      <el-row :gutter="20" class="stat-section">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">
              <el-tag :type="authInfo.auth_status === 'Active' ? 'success' : 'danger'" size="large">
                {{ authInfo.auth_status === 'Active' ? $t('page.authorize.authorized') : authInfo.auth_status || $t('page.authorize.unauthorized') }}
              </el-tag>
            </div>
            <div class="stat-label">{{ $t('page.authorize.authStatus') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ authInfo.auth_due_time || '-' }}</div>
            <div class="stat-label">{{ $t('page.authorize.authDueTime') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ authInfo.service_due_time || '-' }}</div>
            <div class="stat-label">{{ $t('page.authorize.serviceDueTime') }}</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">
              <el-tag :type="authInfo.license_type === 1 ? 'success' : 'info'">
                {{ authInfo.license_type === 1 ? $t('page.authorize.proVersion') : $t('page.authorize.freeVersion') }}
              </el-tag>
            </div>
            <div class="stat-label">{{ $t('page.authorize.licenseType') }}</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 授权详情 -->
      <el-descriptions :title="$t('page.authorize.authDetails')" :column="2" border class="info-section">
        <el-descriptions-item :label="$t('page.authorize.licenseCode')">
          <div class="license-display">
            <el-text truncated>{{ authInfo.system_license || $t('page.authorize.notSet') }}</el-text>
            <el-button type="primary" link @click="showLicenseDialog = true">{{ $t('page.authorize.change') }}</el-button>
          </div>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('page.authorize.serverIP')">{{ authInfo.server_ip || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.authorize.serverName')">{{ authInfo.server_name || '-' }}</el-descriptions-item>
        <el-descriptions-item :label="$t('page.authorize.systemVersion')">
          <el-tag>{{ authInfo.install_version || '-' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('page.authorize.versionType')">
          <el-tag :type="authInfo.system_version_type === 'stable' ? 'success' : 'warning'">
            {{ authInfo.system_version_type === 'stable' ? $t('page.authorize.stable') : $t('page.authorize.beta') }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('page.authorize.suspendReason')" v-if="authInfo.auth_suspend_reason">
          <el-text type="danger">{{ authInfo.auth_suspend_reason }}</el-text>
        </el-descriptions-item>
      </el-descriptions>

      <!-- 已授权应用 -->
      <div class="section" v-if="authInfo.auth_app?.length">
        <div class="section-header">
          <h3>{{ $t('page.authorize.authorizedApps') }}</h3>
        </div>
        <el-table :data="authInfo.auth_app" border style="width: 100%">
          <el-table-column prop="name" :label="$t('page.authorize.appName')" />
          <el-table-column prop="version" :label="$t('page.authorize.version')" width="120" />
          <el-table-column prop="status" :label="$t('page.authorize.status')" width="120">
            <template #default="{ row }">
              <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
                {{ row.status === 'active' ? $t('page.authorize.activated') : $t('page.authorize.notActivated') }}
              </el-tag>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-card>

    <!-- 授权规则管理 -->
    <el-card shadow="never" class="section-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('page.authorize.authRules') }}</span>
          <el-button type="primary" @click="fetchAuthRules" :icon="Refresh">{{ $t('page.authorize.refresh') }}</el-button>
        </div>
      </template>
      <el-table :data="authRules" v-loading="rulesLoading" border style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" :label="$t('page.authorize.ruleName')" min-width="150" />
        <el-table-column prop="description" :label="$t('page.authorize.ruleDescription')" min-width="200" show-overflow-tooltip />
        <el-table-column prop="type" :label="$t('page.authorize.ruleType')" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ ruleTypeMap[row.type] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" :label="$t('page.authorize.status')" width="100" align="center">
          <template #default="{ row }">
            <el-switch v-model="row.status" :active-value="1" :inactive-value="0" @change="handleToggleRule(row)" />
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 版本切换 -->
    <el-card shadow="never" class="section-card">
      <template #header>
        <div class="card-header">
          <span>{{ $t('page.authorize.versionSwitch') }}</span>
        </div>
      </template>
      <el-form :model="versionForm" label-width="120px">
        <el-form-item :label="$t('page.authorize.currentVersionType')">
          <el-tag :type="authInfo.system_version_type === 'stable' ? 'success' : 'warning'" size="large">
            {{ authInfo.system_version_type === 'stable' ? $t('page.authorize.stableWithEn') : $t('page.authorize.betaWithEn') }}
          </el-tag>
        </el-form-item>
        <el-form-item :label="$t('page.authorize.switchTo')">
          <el-radio-group v-model="versionForm.type">
            <el-radio value="stable">{{ $t('page.authorize.stable') }}</el-radio>
            <el-radio value="beta">{{ $t('page.authorize.beta') }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleToggleVersion" :loading="versionLoading">{{ $t('page.authorize.switchVersion') }}</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 更换授权码对话框 -->
    <el-dialog v-model="showLicenseDialog" :title="$t('page.authorize.changeLicense')" width="500px">
      <el-form :model="licenseForm" label-width="100px">
        <el-form-item :label="$t('page.authorize.licenseCode')" required>
          <el-input v-model="licenseForm.license" :placeholder="$t('page.authorize.enterLicense')" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <el-alert type="warning" :closable="false" style="margin-top: 16px">
        {{ $t('page.authorize.licenseWarning') }}
      </el-alert>
      <template #footer>
        <el-button @click="showLicenseDialog = false">{{ $t('page.authorize.cancel') }}</el-button>
        <el-button type="primary" @click="handleUpdateLicense" :loading="licenseLoading">{{ $t('page.authorize.confirmChange') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

const loading = ref(false)
const rulesLoading = ref(false)
const licenseLoading = ref(false)
const versionLoading = ref(false)
const showLicenseDialog = ref(false)

const ruleTypeMap: Record<string, string> = {
  access: $t('page.authorize.ruleTypeAccess'),
  feature: $t('page.authorize.ruleTypeFeature'),
  api: $t('page.authorize.ruleTypeApi'),
  storage: $t('page.authorize.ruleTypeStorage')
}

const authInfo = reactive({
  server_ip: '',
  server_name: '',
  server_port: '',
  server_system: '',
  install_version: '',
  system_version_type: '',
  system_license: '',
  auth_status: '',
  auth_suspend_reason: '',
  auth_app: [] as Array<{ name: string; version: string; status: string }>,
  auth_due_time: '',
  service_due_time: '',
  license_type: 0
})

const authRules = ref<any[]>([])
const versionForm = reactive({ type: 'stable' })
const licenseForm = reactive({ license: '' })

const fetchSystemInfo = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/system/info' })
    if (res) {
      Object.assign(authInfo, res)
      versionForm.type = res.system_version_type || 'stable'
    }
  } catch (error) {
    console.error('获取系统信息失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchCommonInfo = async () => {
  try {
    const res = await request.get({ url: '/api/admin/system/authorize' })
    if (res) authInfo.license_type = res.license_type ?? 0
  } catch { /* ignore */ }
}

const fetchAuthRules = async () => {
  rulesLoading.value = true
  try {
    const res = await request.get({ url: '/api/admin/system/authorize/rules' })
    authRules.value = res?.list || res || []
  } catch {
    ElMessage.error($t('page.authorize.fetchRulesFailed'))
  } finally {
    rulesLoading.value = false
  }
}

const handleRefresh = () => {
  fetchSystemInfo()
  fetchCommonInfo()
}

const handleVerify = async () => {
  try {
    await request.get({ url: '/api/admin/system/authorize/verify', showSuccessMessage: true })
    fetchSystemInfo()
  } catch { /* error handled by request */ }
}

const handleToggleRule = async (row: any) => {
  try {
    await request.put({
      url: `/api/admin/system/authorize/rules/${row.id}`,
      data: { status: row.status },
      showSuccessMessage: true
    })
  } catch {
    row.status = row.status === 1 ? 0 : 1
  }
}

const handleUpdateLicense = async () => {
  if (!licenseForm.license.trim()) {
    ElMessage.warning($t('page.authorize.enterLicenseCode'))
    return
  }
  licenseLoading.value = true
  try {
    await request.put({
      url: '/api/admin/system/authorize/license',
      data: { license: licenseForm.license },
      showSuccessMessage: true
    })
    showLicenseDialog.value = false
    licenseForm.license = ''
    fetchSystemInfo()
    fetchCommonInfo()
  } catch { /* error handled by request */ } finally {
    licenseLoading.value = false
  }
}

const handleToggleVersion = async () => {
  try {
    await ElMessageBox.confirm(
      $t('page.authorize.switchVersionConfirm', { type: versionForm.type === 'stable' ? $t('page.authorize.stable') : $t('page.authorize.beta') }),
      $t('page.authorize.switchVersionTitle')
    )
    versionLoading.value = true
    await request.post({
      url: '/api/admin/system/authorize/toggle-version',
      data: { type: versionForm.type },
      showSuccessMessage: true
    })
    fetchSystemInfo()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error($t('page.authorize.switchVersionFailed'))
  } finally {
    versionLoading.value = false
  }
}

onMounted(() => {
  fetchSystemInfo()
  fetchCommonInfo()
  fetchAuthRules()
})
</script>

<style scoped lang="scss">
.authorize-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;

  .header-actions {
    display: flex;
    gap: 8px;
  }
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

.info-section {
  margin-bottom: 24px;
}

.section {
  margin-top: 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;

  h3 {
    margin: 0;
    font-size: 16px;
    font-weight: 600;
  }
}

.section-card {
  margin-top: 20px;
}

.license-display {
  display: flex;
  align-items: center;
  gap: 12px;
}
</style>
