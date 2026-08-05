<template>
  <div class="authorize-page">
    <!-- 授权状态卡片 -->
    <el-card shadow="never" v-loading="loading" class="status-card">
      <template #header>
        <div class="card-header">
          <span>系统授权管理</span>
          <div class="header-actions">
            <el-button type="primary" @click="handleRefresh" :icon="Refresh">刷新状态</el-button>
            <el-button type="warning" @click="handleVerify">验证授权</el-button>
          </div>
        </div>
      </template>

      <!-- 授权状态概览 -->
      <el-row :gutter="20" class="stat-section">
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">
              <el-tag :type="authInfo.auth_status === 'Active' ? 'success' : 'danger'" size="large">
                {{ authInfo.auth_status === 'Active' ? '已授权' : authInfo.auth_status || '未授权' }}
              </el-tag>
            </div>
            <div class="stat-label">授权状态</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ authInfo.auth_due_time || '-' }}</div>
            <div class="stat-label">授权到期时间</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">{{ authInfo.service_due_time || '-' }}</div>
            <div class="stat-label">服务到期时间</div>
          </el-card>
        </el-col>
        <el-col :span="6">
          <el-card shadow="hover" class="stat-card">
            <div class="stat-value">
              <el-tag :type="authInfo.license_type === 1 ? 'success' : 'info'">
                {{ authInfo.license_type === 1 ? '专业版' : '免费版' }}
              </el-tag>
            </div>
            <div class="stat-label">许可证类型</div>
          </el-card>
        </el-col>
      </el-row>

      <!-- 授权详情 -->
      <el-descriptions title="授权详情" :column="2" border class="info-section">
        <el-descriptions-item label="授权码">
          <div class="license-display">
            <el-text truncated>{{ authInfo.system_license || '未设置' }}</el-text>
            <el-button type="primary" link @click="showLicenseDialog = true">更换</el-button>
          </div>
        </el-descriptions-item>
        <el-descriptions-item label="服务器IP">{{ authInfo.server_ip || '-' }}</el-descriptions-item>
        <el-descriptions-item label="服务器名称">{{ authInfo.server_name || '-' }}</el-descriptions-item>
        <el-descriptions-item label="系统版本">
          <el-tag>{{ authInfo.install_version || '-' }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="版本类型">
          <el-tag :type="authInfo.system_version_type === 'stable' ? 'success' : 'warning'">
            {{ authInfo.system_version_type === 'stable' ? '稳定版' : '测试版' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="暂停原因" v-if="authInfo.auth_suspend_reason">
          <el-text type="danger">{{ authInfo.auth_suspend_reason }}</el-text>
        </el-descriptions-item>
      </el-descriptions>

      <!-- 已授权应用 -->
      <div class="section" v-if="authInfo.auth_app?.length">
        <div class="section-header">
          <h3>已授权应用</h3>
        </div>
        <el-table :data="authInfo.auth_app" border style="width: 100%">
          <el-table-column prop="name" label="应用名称" />
          <el-table-column prop="version" label="版本" width="120" />
          <el-table-column prop="status" label="状态" width="120">
            <template #default="{ row }">
              <el-tag :type="row.status === 'active' ? 'success' : 'danger'" size="small">
                {{ row.status === 'active' ? '已激活' : '未激活' }}
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
          <span>授权规则管理</span>
          <el-button type="primary" @click="fetchAuthRules" :icon="Refresh">刷新</el-button>
        </div>
      </template>
      <el-table :data="authRules" v-loading="rulesLoading" border style="width: 100%">
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="规则名称" min-width="150" />
        <el-table-column prop="description" label="规则描述" min-width="200" show-overflow-tooltip />
        <el-table-column prop="type" label="规则类型" width="120">
          <template #default="{ row }">
            <el-tag size="small">{{ ruleTypeMap[row.type] || row.type }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100" align="center">
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
          <span>版本切换</span>
        </div>
      </template>
      <el-form :model="versionForm" label-width="120px">
        <el-form-item label="当前版本类型">
          <el-tag :type="authInfo.system_version_type === 'stable' ? 'success' : 'warning'" size="large">
            {{ authInfo.system_version_type === 'stable' ? '稳定版 (Stable)' : '测试版 (Beta)' }}
          </el-tag>
        </el-form-item>
        <el-form-item label="切换到">
          <el-radio-group v-model="versionForm.type">
            <el-radio value="stable">稳定版</el-radio>
            <el-radio value="beta">测试版</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleToggleVersion" :loading="versionLoading">切换版本</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 更换授权码对话框 -->
    <el-dialog v-model="showLicenseDialog" title="更换授权码" width="500px">
      <el-form :model="licenseForm" label-width="100px">
        <el-form-item label="授权码" required>
          <el-input v-model="licenseForm.license" placeholder="请输入新的授权码" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <el-alert type="warning" :closable="false" style="margin-top: 16px">
        更换授权码后将自动验证，请确保授权码正确。错误的授权码将导致系统无法正常使用。
      </el-alert>
      <template #footer>
        <el-button @click="showLicenseDialog = false">取消</el-button>
        <el-button type="primary" @click="handleUpdateLicense" :loading="licenseLoading">确认更换</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const rulesLoading = ref(false)
const licenseLoading = ref(false)
const versionLoading = ref(false)
const showLicenseDialog = ref(false)

const ruleTypeMap: Record<string, string> = {
  access: '访问规则',
  feature: '功能限制',
  api: 'API限制',
  storage: '存储限制'
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
    const res = await request.get({ url: '/api/admin/system/authorize/common-info' })
    if (res) authInfo.license_type = res.license_type ?? 0
  } catch { /* ignore */ }
}

const fetchAuthRules = async () => {
  rulesLoading.value = true
  try {
    const res = await request.get({ url: '/api/admin/system/authorize/rules' })
    authRules.value = res?.list || res || []
  } catch {
    ElMessage.error('获取授权规则失败')
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
    ElMessage.warning('请输入授权码')
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
      `确定要切换到${versionForm.type === 'stable' ? '稳定版' : '测试版'}吗？切换后系统将更新到对应版本。`,
      '版本切换确认'
    )
    versionLoading.value = true
    await request.post({
      url: '/api/admin/system/authorize/toggle-version',
      data: { type: versionForm.type },
      showSuccessMessage: true
    })
    fetchSystemInfo()
  } catch (e: any) {
    if (e !== 'cancel') ElMessage.error('版本切换失败')
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
