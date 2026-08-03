<template>
  <div class="dcim-console-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <el-page-header @back="$router.back()">
        <template #content>
          <div class="header-content">
            <span class="page-title">DCIM 控制台</span>
            <el-tag
              :type="powerStatus === 'on' ? 'success' : 'info'"
              size="small"
              effect="light"
              round
            >
              {{ powerStatus === 'on' ? '运行中' : '已关机' }}
            </el-tag>
          </div>
        </template>
        <template #extra>
          <el-button @click="refreshStatus" :loading="refreshing">
            <el-icon><Refresh /></el-icon>刷新状态
          </el-button>
        </template>
      </el-page-header>
    </div>

    <!-- 服务器信息卡片 -->
    <el-card shadow="never" class="info-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span class="card-title">服务器信息</span>
          <el-tag type="info" size="small">{{ serverInfo.rack || '-' }}</el-tag>
        </div>
      </template>
      <div class="info-grid">
        <div class="info-section">
          <h4 class="section-title">基本信息</h4>
          <div class="info-items">
            <div class="info-item">
              <span class="label">服务器名称</span>
              <span class="value">{{ serverInfo.name || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">序列号</span>
              <span class="value mono">{{ serverInfo.serialNumber || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">机房位置</span>
              <span class="value">{{ serverInfo.dataCenter || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">机架位置</span>
              <span class="value mono">{{ serverInfo.rack || '-' }}</span>
            </div>
          </div>
        </div>

        <div class="info-section">
          <h4 class="section-title">硬件配置</h4>
          <div class="info-items">
            <div class="info-item">
              <span class="label">CPU</span>
              <span class="value">{{ serverInfo.cpu || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">内存</span>
              <span class="value">{{ serverInfo.memory || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">系统盘</span>
              <span class="value">{{ serverInfo.systemDisk || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">数据盘</span>
              <span class="value">{{ serverInfo.dataDisk || '-' }}</span>
            </div>
          </div>
        </div>

        <div class="info-section">
          <h4 class="section-title">网络信息</h4>
          <div class="info-items">
            <div class="info-item">
              <span class="label">主IP</span>
              <span class="value mono">
                {{ serverInfo.mainIp || '-' }}
                <el-button v-if="serverInfo.mainIp" link type="primary" size="small" @click="copyText(serverInfo.mainIp)">
                  <el-icon><CopyDocument /></el-icon>
                </el-button>
              </span>
            </div>
            <div class="info-item">
              <span class="label">子网掩码</span>
              <span class="value mono">{{ serverInfo.subnetMask || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">网关</span>
              <span class="value mono">{{ serverInfo.gateway || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="label">带宽</span>
              <span class="value">{{ serverInfo.bandwidth || '-' }}</span>
            </div>
          </div>
        </div>
      </div>
    </el-card>

    <!-- 电源操作 -->
    <el-card shadow="never" class="action-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">电源管理</span>
          <div class="power-indicator">
            <span class="power-dot" :class="powerStatus"></span>
            <span>{{ powerStatus === 'on' ? '运行中' : '已关机' }}</span>
          </div>
        </div>
      </template>
      <div class="action-grid">
        <el-button
          type="success"
          size="large"
          class="action-btn"
          @click="handlePowerAction('on')"
          :loading="actionLoading === 'on'"
          :disabled="powerStatus === 'on'"
        >
          <el-icon><VideoPlay /></el-icon>
          <span>开机</span>
        </el-button>
        <el-button
          type="danger"
          size="large"
          class="action-btn"
          @click="handlePowerAction('off')"
          :loading="actionLoading === 'off'"
          :disabled="powerStatus === 'off'"
        >
          <el-icon><VideoPause /></el-icon>
          <span>硬关机</span>
        </el-button>
        <el-button
          type="warning"
          size="large"
          class="action-btn"
          @click="handlePowerAction('reboot')"
          :loading="actionLoading === 'reboot'"
          :disabled="powerStatus === 'off'"
        >
          <el-icon><RefreshRight /></el-icon>
          <span>硬重启</span>
        </el-button>
        <el-button
          size="large"
          class="action-btn"
          @click="handlePowerAction('reset')"
          :loading="actionLoading === 'reset'"
          :disabled="powerStatus === 'off'"
        >
          <el-icon><Refresh /></el-icon>
          <span>BMC重置</span>
        </el-button>
      </div>
    </el-card>

    <!-- 高级操作 -->
    <el-card shadow="never" class="action-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">高级操作</span>
          <el-tag type="warning" size="small" effect="light">请谨慎操作</el-tag>
        </div>
      </template>
      <div class="advanced-grid">
        <!-- 救援系统 -->
        <div class="advanced-item">
          <div class="item-header">
            <el-icon :size="24" color="#e6a23c"><FirstAidKit /></el-icon>
            <div>
              <h4>救援系统</h4>
              <p>启动临时救援系统进行故障修复</p>
            </div>
          </div>
          <el-button type="warning" plain @click="showRescueDialog = true" :disabled="powerStatus === 'off'">
            进入救援
          </el-button>
        </div>

        <!-- 密码重置 -->
        <div class="advanced-item">
          <div class="item-header">
            <el-icon :size="24" color="#409eff"><Key /></el-icon>
            <div>
              <h4>密码重置</h4>
              <p>重置系统root/Administrator密码</p>
            </div>
          </div>
          <el-button type="primary" plain @click="showPasswordDialog = true">
            重置密码
          </el-button>
        </div>

        <!-- 格式化 -->
        <div class="advanced-item">
          <div class="item-header">
            <el-icon :size="24" color="#f56c6c"><Delete /></el-icon>
            <div>
              <h4>磁盘格式化</h4>
              <p>格式化系统盘或数据盘（数据将丢失）</p>
            </div>
          </div>
          <el-button type="danger" plain @click="showFormatDialog = true" :disabled="powerStatus === 'on'">
            格式化磁盘
          </el-button>
        </div>

        <!-- VNC 控制台 -->
        <div class="advanced-item">
          <div class="item-header">
            <el-icon :size="24" color="#67c23a"><Monitor /></el-icon>
            <div>
              <h4>VNC 控制台</h4>
              <p>通过浏览器远程访问服务器桌面</p>
            </div>
          </div>
          <el-button type="success" plain @click="openVnc" :disabled="powerStatus === 'off'">
            打开VNC
          </el-button>
        </div>

        <!-- 取消任务 -->
        <div class="advanced-item">
          <div class="item-header">
            <el-icon :size="24" color="#909399"><CircleClose /></el-icon>
            <div>
              <h4>取消任务</h4>
              <p>取消当前正在执行的任务</p>
            </div>
          </div>
          <el-button plain @click="cancelTask" :disabled="!hasActiveTask">
            取消任务
          </el-button>
        </div>

        <!-- 流量监控 -->
        <div class="advanced-item">
          <div class="item-header">
            <el-icon :size="24" color="#0056FF"><TrendCharts /></el-icon>
            <div>
              <h4>流量监控</h4>
              <p>查看流入/流出流量统计</p>
            </div>
          </div>
          <el-button plain @click="showTrafficDialog = true">
            查看流量
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 操作日志 -->
    <el-card shadow="never" class="log-card">
      <template #header>
        <div class="card-header">
          <span class="card-title">操作日志</span>
          <el-button link type="primary" size="small" @click="fetchLogs">刷新</el-button>
        </div>
      </template>
      <el-table :data="logs" stripe style="width: 100%" v-loading="logLoading" empty-text="暂无操作日志">
        <el-table-column prop="time" label="操作时间" width="180" />
        <el-table-column prop="action" label="操作类型" width="120">
          <template #default="{ row }">
            <el-tag :type="getLogTagType(row.action)" size="small" effect="light">
              {{ row.action }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="detail" label="操作详情" min-width="200" show-overflow-tooltip />
        <el-table-column prop="operator" label="操作人" width="120" />
        <el-table-column prop="ip" label="IP地址" width="140">
          <template #default="{ row }">
            <span class="mono-text">{{ row.ip }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="result" label="结果" width="100">
          <template #default="{ row }">
            <el-tag :type="row.result === '成功' ? 'success' : 'danger'" size="small" effect="light">
              {{ row.result }}
            </el-tag>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- 救援系统对话框 -->
    <el-dialog v-model="showRescueDialog" title="救援系统" width="500px" :close-on-click-modal="false">
      <el-alert type="info" :closable="false" show-icon title="关于救援模式">
        <template #default>
          <p>救援模式会临时启动一个独立的Linux系统环境，用于修复原系统无法正常启动的问题。</p>
          <ul style="margin: 8px 0 0 20px;">
            <li>原系统磁盘将被挂载到 /mnt 目录</li>
            <li>您可以通过 SSH 连接到救援系统</li>
            <li>修复完成后，重启服务器即可恢复正常</li>
          </ul>
        </template>
      </el-alert>
      <el-form :model="rescueForm" label-width="100px" style="margin-top: 16px;">
        <el-form-item label="救援系统类型">
          <el-select v-model="rescueForm.type" style="width: 100%;">
            <el-option label="Linux 救援系统" value="linux" />
            <el-option label="Windows PE 救援" value="winpe" />
          </el-select>
        </el-form-item>
        <el-form-item label="SSH密码">
          <el-input v-model="rescueForm.password" type="password" placeholder="设置救援系统的SSH密码" show-password />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRescueDialog = false">取消</el-button>
        <el-button type="warning" :loading="rescueLoading" @click="confirmRescue">进入救援模式</el-button>
      </template>
    </el-dialog>

    <!-- 密码重置对话框 -->
    <el-dialog v-model="showPasswordDialog" title="重置密码" width="480px" :close-on-click-modal="false">
      <el-form :model="passwordForm" label-width="100px">
        <el-form-item label="用户类型">
          <el-select v-model="passwordForm.user" style="width: 100%;">
            <el-option label="root / Administrator" value="root" />
            <el-option label="自定义用户" value="custom" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="passwordForm.user === 'custom'" label="用户名">
          <el-input v-model="passwordForm.username" placeholder="请输入用户名" />
        </el-form-item>
        <el-form-item label="新密码">
          <el-input v-model="passwordForm.password" type="password" placeholder="请输入新密码" show-password />
        </el-form-item>
      </el-form>
      <el-alert type="warning" :closable="false" show-icon title="注意" description="密码重置后需要重启服务器才能生效" style="margin-top: 8px;" />
      <template #footer>
        <el-button @click="showPasswordDialog = false">取消</el-button>
        <el-button type="primary" :loading="passwordLoading" @click="confirmResetPassword">确认重置</el-button>
      </template>
    </el-dialog>

    <!-- 格式化对话框 -->
    <el-dialog v-model="showFormatDialog" title="磁盘格式化" width="500px" :close-on-click-modal="false">
      <el-alert type="error" :closable="false" show-icon title="危险操作">
        <template #default>
          <p><strong>格式化将清除磁盘上的所有数据，此操作不可恢复！</strong></p>
          <p>请确保您已备份所有重要数据。</p>
        </template>
      </el-alert>
      <el-form :model="formatForm" label-width="100px" style="margin-top: 16px;">
        <el-form-item label="格式化类型">
          <el-radio-group v-model="formatForm.type">
            <el-radio value="full">全盘格式化</el-radio>
            <el-radio value="system">仅系统分区</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="确认操作">
          <el-checkbox v-model="formatForm.confirmed">我已完成数据备份，确认格式化</el-checkbox>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showFormatDialog = false">取消</el-button>
        <el-button type="danger" :loading="formatLoading" :disabled="!formatForm.confirmed" @click="confirmFormat">
          确认格式化
        </el-button>
      </template>
    </el-dialog>

    <!-- 流量监控对话框 -->
    <el-dialog v-model="showTrafficDialog" title="流量监控" width="600px">
      <div class="traffic-info">
        <div class="traffic-item">
          <div class="traffic-label">
            <el-icon color="#67c23a"><Bottom /></el-icon>
            <span>流入流量</span>
          </div>
          <div class="traffic-value">{{ trafficData.inbound }}</div>
          <el-progress :percentage="trafficData.inboundPercent" :stroke-width="8" color="#67c23a" />
        </div>
        <div class="traffic-item">
          <div class="traffic-label">
            <el-icon color="#409eff"><Top /></el-icon>
            <span>流出流量</span>
          </div>
          <div class="traffic-value">{{ trafficData.outbound }}</div>
          <el-progress :percentage="trafficData.outboundPercent" :stroke-width="8" />
        </div>
        <div class="traffic-total">
          <span>本月流量配额：{{ trafficData.total }}</span>
          <span>已使用：{{ trafficData.used }}</span>
        </div>
      </div>
      <template #footer>
        <el-button @click="showTrafficDialog = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Refresh, RefreshRight, VideoPlay, VideoPause, CopyDocument,
  FirstAidKit, Key, Delete, Monitor, CircleClose, TrendCharts,
  Top, Bottom
} from '@element-plus/icons-vue'
import request from '@/utils/request'

const route = useRoute()
const router = useRouter()

const loading = ref(false)
const refreshing = ref(false)
const logLoading = ref(false)
const powerStatus = ref<'on' | 'off'>('on')
const actionLoading = ref<string | null>(null)
const hasActiveTask = ref(false)

// 对话框状态
const showRescueDialog = ref(false)
const showPasswordDialog = ref(false)
const showFormatDialog = ref(false)
const showTrafficDialog = ref(false)

// 加载状态
const rescueLoading = ref(false)
const passwordLoading = ref(false)
const formatLoading = ref(false)

// 服务器信息
const serverInfo = reactive({
  name: '',
  serialNumber: '',
  dataCenter: '',
  rack: '',
  cpu: '',
  memory: '',
  systemDisk: '',
  dataDisk: '',
  mainIp: '',
  subnetMask: '',
  gateway: '',
  bandwidth: ''
})

// 表单数据
const rescueForm = reactive({ type: 'linux', password: '' })
const passwordForm = reactive({ user: 'root', username: '', password: '' })
const formatForm = reactive({ type: 'full', confirmed: false })

// 流量数据
const trafficData = reactive({
  inbound: '128.5 GB',
  outbound: '45.2 GB',
  inboundPercent: 42,
  outboundPercent: 15,
  total: '500 GB',
  used: '173.7 GB'
})

// 操作日志
const logs = ref<any[]>([])

// 获取服务器信息
async function fetchServerInfo() {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) return

  loading.value = true
  try {
    const { data } = await request.get(`/api/v2/hosts/${hostId}`)
    if (data?.data) {
      const info = data.data
      serverInfo.name = info.product_name || info.name || ''
      serverInfo.serialNumber = info.serial_number || ''
      serverInfo.dataCenter = info.data_center || ''
      serverInfo.rack = info.rack || ''
      serverInfo.cpu = info.cpu || ''
      serverInfo.memory = info.memory || ''
      serverInfo.systemDisk = info.disk || ''
      serverInfo.dataDisk = info.data_disk || ''
      serverInfo.mainIp = info.dedicated_ip || ''
      serverInfo.subnetMask = info.subnet_mask || ''
      serverInfo.gateway = info.gateway || ''
      serverInfo.bandwidth = info.bandwidth ? `${info.bandwidth}Mbps` : ''
      powerStatus.value = info.power_status === 'on' ? 'on' : 'off'
    }
  } catch (e) {
    console.error('Failed to fetch server info:', e)
  } finally {
    loading.value = false
  }
}

// 刷新状态
async function refreshStatus() {
  refreshing.value = true
  await fetchServerInfo()
  refreshing.value = false
  ElMessage.success('状态已刷新')
}

// 电源操作
async function handlePowerAction(action: string) {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) return

  const actionMap: Record<string, { title: string; msg: string; api: string }> = {
    on: { title: '确认开机', msg: '确定要开启服务器吗？', api: 'power_on' },
    off: { title: '确认硬关机', msg: '硬关机相当于直接断电，可能导致数据丢失。确定继续吗？', api: 'power_off' },
    reboot: { title: '确认硬重启', msg: '硬重启相当于强制断电重启，可能导致数据丢失。确定继续吗？', api: 'reboot' },
    reset: { title: '确认BMC重置', msg: '确定要重置BMC管理口吗？', api: 'reset_bmc' }
  }

  const config = actionMap[action]
  if (!config) return

  try {
    await ElMessageBox.confirm(config.msg, config.title, {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: action === 'off' ? 'error' : 'warning'
    })

    actionLoading.value = action
    await request.post(`/api/v2/hosts/${hostId}/${config.api}`)
    ElMessage.success('操作已提交')

    // 更新状态
    if (action === 'on') powerStatus.value = 'on'
    else if (action === 'off') powerStatus.value = 'off'
    else if (action === 'reboot') {
      powerStatus.value = 'off'
      setTimeout(() => { powerStatus.value = 'on' }, 5000)
    }
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '操作失败')
    }
  } finally {
    actionLoading.value = null
  }
}

// 救援系统
async function confirmRescue() {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) return

  if (!rescueForm.password) {
    ElMessage.warning('请设置SSH密码')
    return
  }

  rescueLoading.value = true
  try {
    await request.post(`/api/v2/hosts/${hostId}/rescue`, {
      type: rescueForm.type,
      password: rescueForm.password
    })
    showRescueDialog.value = false
    ElMessage.success('救援系统已启动')
  } catch (e: any) {
    ElMessage.error(e.message || '启动救援系统失败')
  } finally {
    rescueLoading.value = false
  }
}

// 密码重置
async function confirmResetPassword() {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) return

  if (!passwordForm.password) {
    ElMessage.warning('请输入新密码')
    return
  }

  passwordLoading.value = true
  try {
    await request.post(`/api/v2/hosts/${hostId}/reset-password`, {
      user: passwordForm.user === 'custom' ? passwordForm.username : passwordForm.user,
      password: passwordForm.password
    })
    showPasswordDialog.value = false
    ElMessage.success('密码已重置，重启后生效')
  } catch (e: any) {
    ElMessage.error(e.message || '密码重置失败')
  } finally {
    passwordLoading.value = false
  }
}

// 格式化磁盘
async function confirmFormat() {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) return

  try {
    await ElMessageBox.confirm('此操作将清除磁盘所有数据且不可恢复，确定继续吗？', '最终确认', {
      confirmButtonText: '确定格式化',
      cancelButtonText: '取消',
      type: 'error',
      confirmButtonClass: 'el-button--danger'
    })

    formatLoading.value = true
    await request.post(`/api/v2/hosts/${hostId}/format`, {
      type: formatForm.type
    })
    showFormatDialog.value = false
    ElMessage.success('格式化任务已提交')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '格式化失败')
    }
  } finally {
    formatLoading.value = false
  }
}

// 打开VNC
function openVnc() {
  const hostId = route.query.host_id || route.params.id
  router.push({ path: '/user/vnc-console', query: { host_id: hostId } })
}

// 取消任务
async function cancelTask() {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) return

  try {
    await ElMessageBox.confirm('确定取消当前正在执行的任务吗？', '取消任务', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.post(`/api/v2/hosts/${hostId}/cancel-task`)
    ElMessage.success('任务已取消')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '操作失败')
    }
  }
}

// 获取日志
async function fetchLogs() {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) return

  logLoading.value = true
  try {
    const { data } = await request.get(`/api/v2/hosts/${hostId}/log`)
    if (data?.data?.list) {
      logs.value = data.data.list
    }
  } catch (e) {
    console.error('Failed to fetch logs:', e)
  } finally {
    logLoading.value = false
  }
}

// 复制文本
function copyText(text: string) {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success('已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

// 日志标签类型
function getLogTagType(action: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    '开机': 'success', '关机': 'danger', '重启': 'warning',
    '密码重置': '', '救援模式': 'warning', '格式化': 'danger'
  }
  return map[action] || 'info'
}

onMounted(() => {
  fetchServerInfo()
  fetchLogs()
})
</script>

<style scoped lang="scss">
.dcim-console-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  :deep(.el-page-header) {
    width: 100%;
  }

  :deep(.el-page-header__content) {
    flex: 1;
  }

  :deep(.el-page-header__extra) {
    flex: none;
  }
}

.header-content {
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #303133;
}

.info-card,
.action-card,
.log-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
}

.info-section {
  .section-title {
    font-size: 14px;
    font-weight: 600;
    color: #303133;
    margin: 0 0 12px 0;
    padding-bottom: 8px;
    border-bottom: 1px solid #f2f3f5;
  }
}

.info-items {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.info-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;

  .label {
    width: 70px;
    flex-shrink: 0;
    font-size: 13px;
    color: #909399;
    line-height: 24px;
  }

  .value {
    flex: 1;
    font-size: 14px;
    color: #303133;
    display: flex;
    align-items: center;
    gap: 4px;
    line-height: 24px;

    &.mono {
      font-family: 'Monaco', 'Menlo', monospace;
    }
  }
}

.power-indicator {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 14px;
  color: #606266;
}

.power-dot {
  width: 10px;
  height: 10px;
  border-radius: 50%;

  &.on {
    background: #67c23a;
    box-shadow: 0 0 8px rgba(103, 194, 58, 0.5);
  }

  &.off {
    background: #909399;
  }
}

.action-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.action-btn {
  height: 80px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  border-radius: 12px;
  font-size: 14px;

  .el-icon {
    font-size: 24px;
  }
}

.advanced-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 16px;
}

.advanced-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
  gap: 16px;

  .item-header {
    display: flex;
    align-items: center;
    gap: 12px;
    flex: 1;

    h4 {
      font-size: 14px;
      font-weight: 600;
      color: #303133;
      margin: 0 0 4px 0;
    }

    p {
      font-size: 12px;
      color: #909399;
      margin: 0;
    }
  }
}

.mono-text {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 13px;
  color: #606266;
}

.traffic-info {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.traffic-item {
  .traffic-label {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 14px;
    color: #303133;
    margin-bottom: 8px;
  }

  .traffic-value {
    font-size: 24px;
    font-weight: 700;
    color: #303133;
    margin-bottom: 8px;
  }
}

.traffic-total {
  display: flex;
  justify-content: space-between;
  padding: 12px 16px;
  background: #f5f7fa;
  border-radius: 8px;
  font-size: 13px;
  color: #606266;
}

@media (max-width: 768px) {
  .info-grid {
    grid-template-columns: 1fr;
  }

  .action-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .advanced-grid {
    grid-template-columns: 1fr;
  }
}
</style>
