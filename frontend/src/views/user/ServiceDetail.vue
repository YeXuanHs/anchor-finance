<template>
  <div class="service-detail-page">
    <!-- 页面头部 -->
    <div class="page-header">
      <el-page-header @back="$router.push('/user/products')">
        <template #content>
          <div class="header-content">
            <span class="page-title">{{ serviceInfo.name || '服务详情' }}</span>
            <el-tag
              v-if="serviceInfo.status"
              :type="getStatusType(serviceInfo.status)"
              size="small"
              effect="light"
              round
              class="status-tag"
            >
              {{ getStatusText(serviceInfo.status) }}
            </el-tag>
          </div>
        </template>
        <template #extra>
          <div class="header-actions">
            <el-button
              v-if="serviceInfo.status === 'Active'"
              type="primary"
              plain
              size="small"
              @click="handleRenew"
            >
              <el-icon><RefreshRight /></el-icon>
              续费
            </el-button>
            <el-button
              v-if="serviceInfo.status === 'Active'"
              plain
              size="small"
              @click="handleUpgrade"
            >
              <el-icon><Top /></el-icon>
              升级
            </el-button>
            <el-dropdown
              v-if="serviceInfo.status === 'Active'"
              trigger="click"
              @command="handleCommand"
            >
              <el-button plain size="small">
                更多操作
                <el-icon class="el-icon--right"><ArrowDown /></el-icon>
              </el-button>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="reinstall">重装系统</el-dropdown-item>
                  <el-dropdown-item command="resetPassword">重置密码</el-dropdown-item>
                  <el-dropdown-item command="rescue">救援模式</el-dropdown-item>
                  <el-dropdown-item command="console">控制台</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </div>
        </template>
      </el-page-header>
    </div>

    <!-- 骨架屏加载 -->
    <template v-if="loading">
      <div class="skeleton-wrapper">
        <el-skeleton :rows="6" animated />
      </div>
    </template>

    <!-- 主要内容 -->
    <template v-else>
      <!-- 基本信息卡片 -->
      <div class="info-section">
        <el-card class="info-card" shadow="never">
          <template #header>
            <div class="card-header">
              <span class="card-title">基本信息</span>
              <el-button
                v-if="serviceInfo.remark !== undefined"
                link
                type="primary"
                size="small"
                @click="showRemarkDialog = true"
              >
                <el-icon><Edit /></el-icon>
                {{ serviceInfo.remark ? '编辑备注' : '添加备注' }}
              </el-button>
            </div>
          </template>

          <div class="info-grid">
            <div class="info-item">
              <span class="info-label">产品名称</span>
              <span class="info-value">{{ serviceInfo.product_name || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">产品类型</span>
              <span class="info-value">{{ serviceInfo.type_name || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">开通时间</span>
              <span class="info-value">{{ serviceInfo.active_time || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">到期时间</span>
              <span
                class="info-value"
                :class="{ 'text-danger': isExpired }"
              >
                {{ serviceInfo.due_time || '-' }}
              </span>
            </div>
            <div class="info-item">
              <span class="info-label">计费周期</span>
              <span class="info-value">{{ serviceInfo.billing_cycle_name || '-' }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">续费价格</span>
              <span class="info-value price">
                {{ serviceInfo.renew_price ? `¥${serviceInfo.renew_price}` : '-' }}
              </span>
            </div>
            <div v-if="serviceInfo.remark" class="info-item full-width">
              <span class="info-label">备注</span>
              <span class="info-value">{{ serviceInfo.remark }}</span>
            </div>
          </div>
        </el-card>
      </div>

      <!-- 动态组件渲染不同类型的服务详情 -->
      <component
        :is="currentComponent"
        v-if="currentComponent"
        :service-info="serviceInfo"
        @action="handleServiceAction"
      />

      <!-- Tabs 区域 -->
      <div class="tabs-section">
        <el-card shadow="never" class="tabs-card">
          <el-tabs v-model="activeTab" @tab-change="handleTabChange">
            <!-- 账单 -->
            <el-tab-pane label="账单" name="billing">
              <div class="tab-content">
                <el-table
                  v-loading="billingLoading"
                  :data="billingList"
                  style="width: 100%"
                  empty-text="暂无账单记录"
                >
                  <el-table-column prop="pay_time" label="支付时间" min-width="160" />
                  <el-table-column prop="type" label="类型" width="120" />
                  <el-table-column prop="amount" label="金额" width="120">
                    <template #default="{ row }">
                      <span class="price">¥{{ row.amount }}</span>
                    </template>
                  </el-table-column>
                  <el-table-column prop="trans_id" label="交易号" min-width="180" />
                  <el-table-column prop="gateway" label="支付方式" width="120" />
                </el-table>
                <div v-if="billingTotal > 0" class="pagination-wrapper">
                  <el-pagination
                    v-model:current-page="billingPage"
                    v-model:page-size="billingPageSize"
                    :total="billingTotal"
                    :page-sizes="[10, 20, 50]"
                    layout="total, sizes, prev, pager, next"
                    @size-change="fetchBilling"
                    @current-change="fetchBilling"
                  />
                </div>
              </div>
            </el-tab-pane>

            <!-- 日志 -->
            <el-tab-pane label="日志" name="log">
              <div class="tab-content">
                <el-table
                  v-loading="logLoading"
                  :data="logList"
                  style="width: 100%"
                  empty-text="暂无操作日志"
                >
                  <el-table-column prop="create_time" label="操作时间" min-width="160" />
                  <el-table-column prop="description" label="操作详情" min-width="200" show-overflow-tooltip />
                  <el-table-column prop="user" label="操作人" width="120" />
                  <el-table-column prop="ipaddr" label="IP地址" width="150" />
                </el-table>
                <div v-if="logTotal > 0" class="pagination-wrapper">
                  <el-pagination
                    v-model:current-page="logPage"
                    v-model:page-size="logPageSize"
                    :total="logTotal"
                    :page-sizes="[10, 20, 50]"
                    layout="total, sizes, prev, pager, next"
                    @size-change="fetchLog"
                    @current-change="fetchLog"
                  />
                </div>
              </div>
            </el-tab-pane>

            <!-- 下载 -->
            <el-tab-pane label="下载" name="download">
              <div class="tab-content">
                <el-table
                  v-loading="downloadLoading"
                  :data="downloadList"
                  style="width: 100%"
                  empty-text="暂无可下载文件"
                >
                  <el-table-column prop="title" label="文件名" min-width="200">
                    <template #default="{ row }">
                      <div class="file-name">
                        <el-icon :size="16" :color="getFileIconColor(row.type)">
                          <Document />
                        </el-icon>
                        <span>{{ row.title }}</span>
                      </div>
                    </template>
                  </el-table-column>
                  <el-table-column prop="create_time" label="上传时间" width="180" />
                  <el-table-column prop="downloads" label="下载次数" width="100" align="center" />
                  <el-table-column label="操作" width="100" align="center">
                    <template #default="{ row }">
                      <el-button link type="primary" @click="handleDownload(row)">
                        <el-icon><Download /></el-icon>
                        下载
                      </el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
            </el-tab-pane>

            <!-- 升级 -->
            <el-tab-pane label="升级" name="upgrade">
              <div class="tab-content">
                <div v-loading="upgradeLoading" class="upgrade-section">
                  <div v-if="upgradePlans.length > 0" class="upgrade-plans">
                    <div
                      v-for="plan in upgradePlans"
                      :key="plan.id"
                      class="plan-card"
                      :class="{ 'is-selected': selectedUpgradePlan?.id === plan.id }"
                      @click="selectedUpgradePlan = plan"
                    >
                      <div class="plan-name">{{ plan.name }}</div>
                      <div class="plan-price">
                        <span class="price-amount">¥{{ plan.price }}</span>
                        <span class="price-period">/{{ plan.billing_cycle_name }}</span>
                      </div>
                      <div v-if="plan.diff_price" class="price-diff">
                        差价: <span class="diff-amount">¥{{ plan.diff_price }}</span>
                      </div>
                      <div v-if="plan.description" class="plan-desc">{{ plan.description }}</div>
                    </div>
                  </div>
                  <el-empty v-else-if="!upgradeLoading" description="暂无可升级套餐" />
                  <div v-if="selectedUpgradePlan" class="upgrade-actions">
                    <el-button type="primary" size="large" @click="confirmUpgrade">
                      确认升级到 {{ selectedUpgradePlan.name }}
                    </el-button>
                  </div>
                </div>
              </div>
            </el-tab-pane>

            <!-- 云服务器高级管理 Tab（仅云服务器类型显示） -->
            <template v-if="isCloudType">
              <el-tab-pane label="快照" name="snapshot">
                <div class="tab-content">
                  <SnapshotManager />
                </div>
              </el-tab-pane>

              <el-tab-pane label="备份" name="backup">
                <div class="tab-content">
                  <BackupManager />
                </div>
              </el-tab-pane>

              <el-tab-pane label="磁盘" name="disk">
                <div class="tab-content">
                  <DiskManager />
                </div>
              </el-tab-pane>

              <el-tab-pane label="VPC网络" name="vpc">
                <div class="tab-content">
                  <VpcManager />
                </div>
              </el-tab-pane>

              <el-tab-pane label="NAT" name="nat">
                <div class="tab-content">
                  <NatManager />
                </div>
              </el-tab-pane>

              <el-tab-pane label="SSH密钥" name="sshkey">
                <div class="tab-content">
                  <SshKeyManager />
                </div>
              </el-tab-pane>
            </template>
          </el-tabs>
        </el-card>
      </div>
    </template>

    <!-- 修改备注弹窗 -->
    <el-dialog
      v-model="showRemarkDialog"
      title="修改备注"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-input
        v-model="remarkValue"
        placeholder="请输入备注信息"
        maxlength="100"
        show-word-limit
      />
      <template #footer>
        <el-button @click="showRemarkDialog = false">取消</el-button>
        <el-button type="primary" :loading="remarkLoading" @click="submitRemark">
          确定
        </el-button>
      </template>
    </el-dialog>

    <!-- 重装系统弹窗 -->
    <el-dialog
      v-model="showReinstallDialog"
      title="重装系统"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-form :model="reinstallForm" label-width="100px">
        <el-form-item label="操作系统" required>
          <el-select v-model="reinstallForm.os" placeholder="请选择操作系统" style="width: 100%">
            <el-option
              v-for="item in osOptions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="系统版本" required>
          <el-select v-model="reinstallForm.version" placeholder="请选择系统版本" style="width: 100%">
            <el-option
              v-for="item in filteredVersions"
              :key="item.value"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="登录密码" required>
          <el-input
            v-model="reinstallForm.password"
            type="password"
            placeholder="请输入新的登录密码"
            show-password
          />
        </el-form-item>
        <el-alert
          type="warning"
          :closable="false"
          show-icon
          title="警告"
          description="重装系统将清除所有数据，请提前备份重要文件！"
        />
      </el-form>
      <template #footer>
        <el-button @click="showReinstallDialog = false">取消</el-button>
        <el-button type="danger" :loading="reinstallLoading" @click="confirmReinstall">
          确认重装
        </el-button>
      </template>
    </el-dialog>

    <!-- 重置密码弹窗 -->
    <el-dialog
      v-model="showResetPasswordDialog"
      title="重置密码"
      width="450px"
      :close-on-click-modal="false"
    >
      <template v-if="!newPassword">
        <p style="margin-bottom: 16px; color: #606266;">
          确认要重置此服务的登录密码吗？重置后将生成新的随机密码。
        </p>
      </template>
      <template v-else>
        <el-result icon="success" title="密码重置成功">
          <template #sub-title>
            <div style="margin-top: 12px;">
              <p style="color: #606266; margin-bottom: 12px;">新密码如下，请妥善保管：</p>
              <el-input v-model="newPassword" readonly>
                <template #append>
                  <el-button @click="copyPassword">复制</el-button>
                </template>
              </el-input>
            </div>
          </template>
        </el-result>
      </template>
      <template #footer>
        <el-button @click="showResetPasswordDialog = false; newPassword = ''">
          {{ newPassword ? '关闭' : '取消' }}
        </el-button>
        <el-button
          v-if="!newPassword"
          type="warning"
          :loading="resetPasswordLoading"
          @click="confirmResetPassword"
        >
          确认重置
        </el-button>
      </template>
    </el-dialog>

    <!-- 救援模式弹窗 -->
    <el-dialog
      v-model="showRescueDialog"
      title="救援模式"
      width="500px"
      :close-on-click-modal="false"
    >
      <el-alert type="info" :closable="false" show-icon title="关于救援模式">
        <template #default>
          <p>救援模式会临时启动一个独立的操作系统环境，用于修复原系统无法正常启动的问题。</p>
          <p style="margin-top: 8px;">进入救援模式后：</p>
          <ul style="margin: 8px 0 0 20px;">
            <li>原系统磁盘将被挂载到 /mnt 目录</li>
            <li>您可以通过 SSH 连接到救援系统</li>
            <li>修复完成后，重启服务器即可恢复正常</li>
          </ul>
        </template>
      </el-alert>
      <el-form :model="rescueForm" label-width="80px" style="margin-top: 16px;">
        <el-form-item label="SSH密码">
          <el-input
            v-model="rescueForm.password"
            type="password"
            placeholder="设置救援系统的SSH密码"
            show-password
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRescueDialog = false">取消</el-button>
        <el-button type="warning" :loading="rescueLoading" @click="confirmRescue">
          进入救援模式
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  RefreshRight, Top, ArrowDown, Edit, Document, Download
} from '@element-plus/icons-vue'
import request from '@/utils/request'

// 服务详情子组件
import ServiceCloud from './service/ServiceCloud.vue'
import ServiceDedicated from './service/ServiceDedicated.vue'
import ServiceDomain from './service/ServiceDomain.vue'
import ServiceSSL from './service/ServiceSSL.vue'
import ServiceCDN from './service/ServiceCDN.vue'
import ServiceHosting from './service/ServiceHosting.vue'
import ServiceNAT from './service/ServiceNAT.vue'

// 云服务器高级管理组件
import SnapshotManager from './components/SnapshotManager.vue'
import BackupManager from './components/BackupManager.vue'
import DiskManager from './components/DiskManager.vue'
import VpcManager from './components/VpcManager.vue'
import NatManager from './components/NatManager.vue'
import SshKeyManager from './components/SshKeyManager.vue'

const route = useRoute()
const router = useRouter()

// 服务类型与组件映射
const componentMap: Record<string, any> = {
  cloud: ServiceCloud,
  server: ServiceDedicated,
  dedicated: ServiceDedicated,
  domain: ServiceDomain,
  ssl: ServiceSSL,
  cdn: ServiceCDN,
  hosting: ServiceHosting,
  hostingaccount: ServiceHosting,
  nat: ServiceNAT
}

// 状态
const loading = ref(true)
const serviceInfo = ref<any>({})
const activeTab = ref('billing')

// 备注相关
const showRemarkDialog = ref(false)
const remarkValue = ref('')
const remarkLoading = ref(false)

// 账单相关
const billingLoading = ref(false)
const billingList = ref<any[]>([])
const billingPage = ref(1)
const billingPageSize = ref(10)
const billingTotal = ref(0)

// 日志相关
const logLoading = ref(false)
const logList = ref<any[]>([])
const logPage = ref(1)
const logPageSize = ref(10)
const logTotal = ref(0)

// 下载相关
const downloadLoading = ref(false)
const downloadList = ref<any[]>([])

// 升级相关
const upgradeLoading = ref(false)
const upgradePlans = ref<any[]>([])
const selectedUpgradePlan = ref<any>(null)

// 重装系统相关
const showReinstallDialog = ref(false)
const reinstallLoading = ref(false)
const reinstallForm = reactive({
  os: '',
  version: '',
  password: ''
})
const osOptions = ref([
  { label: 'CentOS', value: 'centos' },
  { label: 'Ubuntu', value: 'ubuntu' },
  { label: 'Debian', value: 'debian' },
  { label: 'Windows Server', value: 'windows' }
])
const osVersions = ref<Record<string, { label: string; value: string }[]>>({
  centos: [
    { label: 'CentOS 7.9', value: 'centos7.9' },
    { label: 'CentOS 8.5', value: 'centos8.5' },
    { label: 'CentOS Stream 9', value: 'centos-stream9' }
  ],
  ubuntu: [
    { label: 'Ubuntu 20.04 LTS', value: 'ubuntu20.04' },
    { label: 'Ubuntu 22.04 LTS', value: 'ubuntu22.04' },
    { label: 'Ubuntu 24.04 LTS', value: 'ubuntu24.04' }
  ],
  debian: [
    { label: 'Debian 11', value: 'debian11' },
    { label: 'Debian 12', value: 'debian12' }
  ],
  windows: [
    { label: 'Windows Server 2019', value: 'windows2019' },
    { label: 'Windows Server 2022', value: 'windows2022' }
  ]
})

// 重置密码相关
const showResetPasswordDialog = ref(false)
const resetPasswordLoading = ref(false)
const newPassword = ref('')

// 救援模式相关
const showRescueDialog = ref(false)
const rescueLoading = ref(false)
const rescueForm = reactive({ password: '' })

// 计算属性
const currentComponent = computed(() => {
  const type = serviceInfo.value.type?.toLowerCase() || ''
  return componentMap[type] || null
})

const isExpired = computed(() => {
  if (!serviceInfo.value.due_time) return false
  return new Date(serviceInfo.value.due_time) < new Date()
})

const filteredVersions = computed(() => {
  return osVersions.value[reinstallForm.os] || []
})

// 是否为云服务器类型（显示高级管理Tab）
const isCloudType = computed(() => {
  const type = serviceInfo.value.type?.toLowerCase() || ''
  return ['cloud', 'server', 'dedicated'].includes(type)
})

// 方法
function getStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    Active: 'success',
    Suspended: 'warning',
    Expired: 'danger',
    Pending: 'info'
  }
  return map[status] || 'info'
}

function getStatusText(status: string) {
  const map: Record<string, string> = {
    Active: '运行中',
    Suspended: '已暂停',
    Expired: '已过期',
    Pending: '待开通'
  }
  return map[status] || status
}

function getFileIconColor(type: string) {
  const map: Record<string, string> = {
    '1': '#e6a23c',
    '2': '#67c23a',
    '3': '#909399'
  }
  return map[type] || '#409eff'
}

// 获取服务详情
async function fetchServiceDetail() {
  const id = route.params.id
  if (!id) return

  loading.value = true
  try {
    const { data } = await request.get(`/api/v1/hosts/${id}`)
    if (data?.data) {
      serviceInfo.value = data.data
      remarkValue.value = data.data.remark || ''
    }
  } catch (error) {
    ElMessage.error('获取服务详情失败')
  } finally {
    loading.value = false
  }
}

// 获取账单列表
async function fetchBilling() {
  const id = route.params.id
  if (!id) return

  billingLoading.value = true
  try {
    const { data } = await request.get(`/api/v1/hosts/${id}/billing`, {
      params: {
        page: billingPage.value,
        limit: billingPageSize.value
      }
    })
    if (data?.data) {
      billingList.value = data.data.list || []
      billingTotal.value = data.data.total || 0
    }
  } catch (error) {
    console.error('获取账单失败', error)
  } finally {
    billingLoading.value = false
  }
}

// 获取日志列表
async function fetchLog() {
  const id = route.params.id
  if (!id) return

  logLoading.value = true
  try {
    const { data } = await request.get(`/api/v1/hosts/${id}/log`, {
      params: {
        page: logPage.value,
        limit: logPageSize.value
      }
    })
    if (data?.data) {
      logList.value = data.data.list || []
      logTotal.value = data.data.total || 0
    }
  } catch (error) {
    console.error('获取日志失败', error)
  } finally {
    logLoading.value = false
  }
}

// 获取下载列表
async function fetchDownload() {
  const id = route.params.id
  if (!id) return

  downloadLoading.value = true
  try {
    const { data } = await request.get(`/api/v1/hosts/${id}/download`)
    if (data?.data) {
      downloadList.value = data.data || []
    }
  } catch (error) {
    console.error('获取下载列表失败', error)
  } finally {
    downloadLoading.value = false
  }
}

// Tab 切换
function handleTabChange(tab: string) {
  switch (tab) {
    case 'billing':
      if (billingList.value.length === 0) fetchBilling()
      break
    case 'log':
      if (logList.value.length === 0) fetchLog()
      break
    case 'download':
      if (downloadList.value.length === 0) fetchDownload()
      break
    case 'upgrade':
      if (upgradePlans.value.length === 0) fetchUpgradePlans()
      break
  }
}

// 提交备注
async function submitRemark() {
  const id = route.params.id
  if (!id) return

  remarkLoading.value = true
  try {
    await request.post(`/api/v1/hosts/${id}/remark`, {
      remark: remarkValue.value
    })
    serviceInfo.value.remark = remarkValue.value
    showRemarkDialog.value = false
    ElMessage.success('备注修改成功')
  } catch (error) {
    ElMessage.error('备注修改失败')
  } finally {
    remarkLoading.value = false
  }
}

// 操作处理
function handleRenew() {
  const id = route.params.id
  router.push(`/user/upgrade?host=${id}&action=renew`)
}

function handleUpgrade() {
  const id = route.params.id
  router.push(`/user/upgrade?host=${id}&action=upgrade`)
}

function handleCommand(command: string) {
  switch (command) {
    case 'reinstall':
      reinstallForm.os = ''
      reinstallForm.version = ''
      reinstallForm.password = ''
      showReinstallDialog.value = true
      break
    case 'resetPassword':
      newPassword.value = ''
      showResetPasswordDialog.value = true
      break
    case 'rescue':
      rescueForm.password = ''
      showRescueDialog.value = true
      break
    case 'console':
      openConsole()
      break
  }
}

function handleServiceAction(action: string) {
  ElMessage.info(`操作: ${action}`)
}

function handleDownload(file: any) {
  if (file.down_link) {
    window.open(file.down_link, '_blank')
  }
}

// 获取可升级套餐列表
async function fetchUpgradePlans() {
  const id = route.params.id
  if (!id) return

  upgradeLoading.value = true
  try {
    const { data } = await request.get(`/api/v1/upgrades/available/${id}`)
    if (data?.data) {
      upgradePlans.value = data.data.list || []
    }
  } catch (error) {
    console.error('获取升级套餐失败', error)
  } finally {
    upgradeLoading.value = false
  }
}

// 确认升级
async function confirmUpgrade() {
  const id = route.params.id
  if (!id || !selectedUpgradePlan.value) return

  try {
    await ElMessageBox.confirm(
      `确认升级到 ${selectedUpgradePlan.value.name}？差价: ¥${selectedUpgradePlan.value.diff_price}`,
      '确认升级',
      { type: 'warning' }
    )
    await request.post('/api/v1/upgrades', {
      host_id: id,
      plan_id: selectedUpgradePlan.value.id
    })
    ElMessage.success('升级成功')
    fetchServiceDetail()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '升级失败')
    }
  }
}

// 确认重装系统
async function confirmReinstall() {
  const id = route.params.id
  if (!id) return

  if (!reinstallForm.os || !reinstallForm.version || !reinstallForm.password) {
    ElMessage.warning('请填写完整信息')
    return
  }

  try {
    await ElMessageBox.confirm('重装系统将清除所有数据，确认继续？', '确认重装', {
      type: 'error',
      confirmButtonText: '确认重装',
      confirmButtonClass: 'el-button--danger'
    })
    reinstallLoading.value = true
    await request.post(`/api/v1/hosts/${id}/reinstall`, {
      os: reinstallForm.version,
      password: reinstallForm.password
    })
    showReinstallDialog.value = false
    ElMessage.success('重装系统任务已提交，请稍后查看')
    fetchServiceDetail()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '重装系统失败')
    }
  } finally {
    reinstallLoading.value = false
  }
}

// 确认重置密码
async function confirmResetPassword() {
  const id = route.params.id
  if (!id) return

  resetPasswordLoading.value = true
  try {
    const { data } = await request.post(`/api/v1/hosts/${id}/reset-password`)
    if (data?.data?.password) {
      newPassword.value = data.data.password
    }
    ElMessage.success('密码重置成功')
  } catch (error: any) {
    ElMessage.error(error.message || '密码重置失败')
  } finally {
    resetPasswordLoading.value = false
  }
}

// 复制密码
function copyPassword() {
  navigator.clipboard.writeText(newPassword.value)
  ElMessage.success('密码已复制到剪贴板')
}

// 确认进入救援模式
async function confirmRescue() {
  const id = route.params.id
  if (!id) return

  if (!rescueForm.password) {
    ElMessage.warning('请设置SSH密码')
    return
  }

  rescueLoading.value = true
  try {
    await request.post(`/api/v1/hosts/${id}/rescue`, {
      password: rescueForm.password
    })
    showRescueDialog.value = false
    ElMessage.success('救援模式已启动，请通过SSH连接')
    fetchServiceDetail()
  } catch (error: any) {
    ElMessage.error(error.message || '启动救援模式失败')
  } finally {
    rescueLoading.value = false
  }
}

// 打开控制台
async function openConsole() {
  const id = route.params.id
  if (!id) return

  try {
    const { data } = await request.get(`/api/v1/hosts/${id}/operations`)
    if (data?.data?.url) {
      window.open(data.data.url, '_blank')
    } else {
      ElMessage.warning('获取控制台地址失败')
    }
  } catch (error: any) {
    ElMessage.error(error.message || '获取控制台地址失败')
  }
}

// 初始化
onMounted(() => {
  fetchServiceDetail()
  fetchBilling()
})
</script>

<style scoped lang="scss">
.service-detail-page {
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

.status-tag {
  vertical-align: middle;
}

.header-actions {
  display: flex;
  gap: 8px;
}

.skeleton-wrapper {
  background: #fff;
  border-radius: 12px;
  padding: 24px;
}

.info-section {
  .info-card {
    border-radius: 12px;
    border: 1px solid #e8ecf1;

    :deep(.el-card__header) {
      padding: 16px 20px;
      border-bottom: 1px solid #f2f3f5;
    }

    :deep(.el-card__body) {
      padding: 20px;
    }
  }
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
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 16px;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 4px;

  &.full-width {
    grid-column: 1 / -1;
  }
}

.info-label {
  font-size: 13px;
  color: #909399;
}

.info-value {
  font-size: 14px;
  color: #303133;
  font-weight: 500;

  &.price {
    color: #fa8c16;
    font-weight: 600;
  }

  &.text-danger {
    color: #f56c6c;
  }
}

.tabs-section {
  .tabs-card {
    border-radius: 12px;
    border: 1px solid #e8ecf1;

    :deep(.el-card__body) {
      padding: 0 20px 20px;
    }

    :deep(.el-tabs__header) {
      margin: 0;
      padding: 0 0 16px;
    }

    :deep(.el-tabs__nav-wrap::after) {
      display: none;
    }
  }
}

.tab-content {
  min-height: 200px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.file-name {
  display: flex;
  align-items: center;
  gap: 8px;
}

.price {
  color: #f56c6c;
  font-weight: 600;
}

.upgrade-section {
  padding: 16px 0;
}

.upgrade-plans {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 16px;
}

.plan-card {
  border: 2px solid #e4e7ed;
  border-radius: 12px;
  padding: 20px;
  cursor: pointer;
  transition: all 0.3s;
  background: #fff;

  &:hover {
    border-color: #409eff;
    box-shadow: 0 4px 12px rgba(64, 158, 255, 0.15);
  }

  &.is-selected {
    border-color: #409eff;
    background: #ecf5ff;
  }
}

.plan-name {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 12px;
}

.plan-price {
  margin-bottom: 8px;
}

.price-amount {
  font-size: 24px;
  font-weight: 700;
  color: #f56c6c;
}

.price-period {
  font-size: 14px;
  color: #909399;
}

.price-diff {
  font-size: 14px;
  color: #606266;
  margin-bottom: 8px;
}

.diff-amount {
  color: #f56c6c;
  font-weight: 600;
}

.plan-desc {
  font-size: 13px;
  color: #909399;
  line-height: 1.5;
}

.upgrade-actions {
  margin-top: 24px;
  text-align: center;
}

@media (max-width: 768px) {
  .page-header {
    :deep(.el-page-header__main) {
      flex-direction: column;
      align-items: flex-start;
    }
  }

  .header-actions {
    margin-top: 12px;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }
}
</style>
