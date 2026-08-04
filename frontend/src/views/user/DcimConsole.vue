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

        <!-- 重装系统 -->
        <div class="advanced-item">
          <div class="item-header">
            <el-icon :size="24" color="#e6a23c"><Setting /></el-icon>
            <div>
              <h4>重装系统</h4>
              <p>重新安装操作系统（数据将丢失）</p>
            </div>
          </div>
          <el-button type="warning" plain @click="openReinstallDialog" :disabled="powerStatus === 'on'">
            重装系统
          </el-button>
        </div>

        <!-- KVM/IPMI -->
        <div class="advanced-item">
          <div class="item-header">
            <el-icon :size="24" color="#67c23a"><Connection /></el-icon>
            <div>
              <h4>KVM/IPMI</h4>
              <p>通过KVM远程管理服务器硬件</p>
            </div>
          </div>
          <div class="kvm-btn-group">
            <el-button type="success" plain @click="openKvm" :loading="kvmLoading" :disabled="powerStatus === 'off'">
              KVM
            </el-button>
            <el-button plain @click="openIkv" :loading="ikvmLoading" :disabled="powerStatus === 'off'">
              iKVM
            </el-button>
          </div>
        </div>

        <!-- 任务队列 -->
        <div class="advanced-item">
          <div class="item-header">
            <el-icon :size="24" color="#409eff"><Tickets /></el-icon>
            <div>
              <h4>任务队列</h4>
              <p>查看和管理当前执行的任务</p>
              <el-tag v-if="hasActiveTask" type="warning" size="small" effect="light" style="margin-top: 4px;">
                有活跃任务
              </el-tag>
            </div>
          </div>
          <el-button type="primary" plain @click="openTaskDialog">
            查看任务
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

    <!-- 重装系统对话框 -->
    <el-dialog v-model="showReinstallDialog" title="重装系统" width="620px" :close-on-click-modal="false">
      <!-- 重装次数/权限检查提示 -->
      <el-alert
        v-if="reinstallCheckData && reinstallCheckData.canReinstall === false"
        type="warning"
        :closable="false"
        show-icon
        title="重装次数受限"
        style="margin-bottom: 16px;"
      >
        <template #default>
          <p v-if="reinstallCheckData.price">
            本周重装次数已达上限，可购买额外次数：<strong>¥{{ reinstallCheckData.price }}</strong>
          </p>
          <p v-else>本周重装次数已达最大限额，请下周重试或联系技术支持。</p>
        </template>
      </el-alert>

      <el-alert type="error" :closable="false" show-icon title="危险操作" style="margin-bottom: 16px;">
        <template #default>
          <p><strong>重装系统将清除当前系统所有数据，此操作不可恢复！</strong></p>
          <p>请确保您已备份所有重要数据。</p>
        </template>
      </el-alert>

      <el-form :model="reinstallForm" label-width="100px">
        <!-- 操作系统选择 -->
        <el-form-item label="操作系统">
          <div style="display: flex; gap: 12px; width: 100%;">
            <el-select
              v-model="reinstallForm.osGroupId"
              placeholder="系统类型"
              style="width: 40%;"
              @change="onOsGroupChange"
            >
              <el-option
                v-for="group in osGroups"
                :key="group.id"
                :label="group.name"
                :value="group.id"
              />
            </el-select>
            <el-select
              v-model="reinstallForm.osId"
              placeholder="系统版本"
              style="width: 60%;"
            >
              <el-option
                v-for="os in filteredOsList"
                :key="os.id"
                :label="os.name"
                :value="os.id"
              />
            </el-select>
          </div>
        </el-form-item>

        <!-- 密码 -->
        <el-form-item label="密码">
          <div style="display: flex; gap: 8px; width: 100%;">
            <el-input
              v-model="reinstallForm.password"
              type="password"
              placeholder="设置系统密码"
              show-password
              style="flex: 1;"
            />
            <el-button @click="generatePassword">
              <el-icon><Refresh /></el-icon>
              随机
            </el-button>
          </div>
        </el-form-item>

        <!-- SSH/RDP端口 -->
        <el-form-item label="端口">
          <el-input-number
            v-model="reinstallForm.port"
            :min="1"
            :max="65535"
            controls-position="right"
            style="width: 200px;"
          />
          <span class="form-tip" style="margin-left: 8px; color: #909399; font-size: 12px;">
            Linux默认22，Windows默认3389
          </span>
        </el-form-item>

        <!-- 分区方案 -->
        <el-form-item label="分区方案">
          <el-radio-group v-model="reinstallForm.partType">
            <el-radio :value="0">全盘格式化</el-radio>
            <el-radio :value="1">仅格式化第一分区</el-radio>
          </el-radio-group>
        </el-form-item>

        <!-- 确认备份 -->
        <el-form-item>
          <el-checkbox v-model="reinstallForm.confirmed">
            我已完成数据备份，确认重装系统
          </el-checkbox>
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="showReinstallDialog = false">取消</el-button>
        <el-button
          type="danger"
          :loading="reinstallLoading"
          :disabled="!reinstallForm.confirmed || (reinstallCheckData && reinstallCheckData.canReinstall === false)"
          @click="confirmReinstall"
        >
          确认重装
        </el-button>
      </template>
    </el-dialog>

    <!-- 任务队列对话框 -->
    <el-dialog
      v-model="showTaskDialog"
      title="任务队列"
      width="700px"
      @close="closeTaskDialog"
    >
      <div class="task-header" style="display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px;">
        <span style="color: #606266; font-size: 13px;">
          任务列表每10秒自动刷新
        </span>
        <el-button link type="primary" size="small" @click="fetchTasks" :loading="taskLoading">
          <el-icon><Refresh /></el-icon>
          刷新
        </el-button>
      </div>

      <el-table
        :data="taskList"
        stripe
        style="width: 100%"
        v-loading="taskLoading"
        empty-text="暂无任务"
      >
        <el-table-column prop="id" label="任务ID" width="80" />
        <el-table-column label="任务类型" width="120">
          <template #default="{ row }">
            <span>{{ getTaskTypeText(row.type || row.action) }}</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getTaskStatusType(row.status)" size="small" effect="light">
              {{ getTaskStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="进度" width="120">
          <template #default="{ row }">
            <el-progress
              v-if="row.progress !== undefined"
              :percentage="row.progress"
              :stroke-width="6"
              :status="row.status === 'failed' ? 'exception' : row.status === 'completed' ? 'success' : undefined"
            />
            <span v-else>-</span>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" min-width="150" />
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'running' || row.status === 'pending' || row.status === 'processing'"
              type="danger"
              link
              size="small"
              @click="cancelSpecificTask(row.id)"
            >
              取消
            </el-button>
            <span v-else style="color: #c0c4cc; font-size: 12px;">-</span>
          </template>
        </el-table-column>
      </el-table>

      <template #footer>
        <el-button @click="closeTaskDialog">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onBeforeUnmount } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Refresh, RefreshRight, VideoPlay, VideoPause, CopyDocument,
  FirstAidKit, Key, Delete, Monitor, CircleClose, TrendCharts,
  Top, Bottom, Setting, Connection, Tickets, Cpu, Warning
} from '@element-plus/icons-vue'
import request from '@/utils/request'

const { t } = useI18n()
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

// ==================== 重装系统 ====================
const showReinstallDialog = ref(false)
const reinstallLoading = ref(false)
const reinstallForm = reactive({
  osGroupId: null as number | null,
  osId: null as number | null,
  password: '',
  port: 22,
  partType: 0,
  disk: 0,
  confirmed: false
})
const osGroups = ref<any[]>([])
const osList = ref<any[]>([])
const filteredOsList = computed(() => {
  if (!reinstallForm.osGroupId) return osList.value
  return osList.value.filter((os: any) => os.group === reinstallForm.osGroupId)
})
const reinstallCheckData = ref<any>(null)

// ==================== KVM/IPMI ====================
const kvmLoading = ref(false)
const ikvmLoading = ref(false)

// ==================== 任务队列 ====================
const showTaskDialog = ref(false)
const taskLoading = ref(false)
const taskList = ref<any[]>([])
const taskTimer = ref<ReturnType<typeof setInterval> | null>(null)

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

// ==================== 重装系统功能 ====================

// 打开重装系统对话框
async function openReinstallDialog() {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) return

  reinstallForm.osGroupId = null
  reinstallForm.osId = null
  reinstallForm.password = ''
  reinstallForm.port = 22
  reinstallForm.partType = 0
  reinstallForm.disk = 0
  reinstallForm.confirmed = false
  reinstallCheckData.value = null

  showReinstallDialog.value = true
  await Promise.all([fetchCheckReinstall(), fetchOsList()])
}

// 检查重装权限和次数
async function fetchCheckReinstall() {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) return

  try {
    const { data } = await request.get(`/api/v2/hosts/${hostId}/check-reinstall`)
    reinstallCheckData.value = data?.data || null
  } catch (e: any) {
    if (e?.response?.data?.status === 400 && e.response.data.price) {
      reinstallCheckData.value = { canReinstall: false, price: e.response.data.price }
    }
  }
}

// 获取可用操作系统列表
async function fetchOsList() {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) return

  try {
    const { data } = await request.get(`/api/v2/hosts/${hostId}/os-list`)
    if (data?.data) {
      osGroups.value = data.data.groups || []
      osList.value = data.data.os || []
      if (osGroups.value.length > 0 && !reinstallForm.osGroupId) {
        reinstallForm.osGroupId = osGroups.value[0].id
      }
      if (filteredOsList.value.length > 0 && !reinstallForm.osId) {
        reinstallForm.osId = filteredOsList.value[0].id
      }
    }
  } catch (e) {
    console.error('Failed to fetch OS list:', e)
  }
}

// 生成随机密码
function generatePassword() {
  const chars = 'ABCDEFGHJKMNPQRSTUVWXYZabcdefghjkmnpqrstuvwxyz23456789'
  const specials = '!@#$%&*'
  let pass = ''
  pass += chars.charAt(Math.floor(Math.random() * 26))
  pass += specials.charAt(Math.floor(Math.random() * specials.length))
  pass += chars.charAt(26 + Math.floor(Math.random() * 26))
  for (let i = 0; i < 9; i++) {
    pass += chars.charAt(Math.floor(Math.random() * chars.length))
  }
  reinstallForm.password = pass.split('').sort(() => Math.random() - 0.5).join('')
}

// OS组变化时自动选择第一个OS
function onOsGroupChange() {
  const list = filteredOsList.value
  reinstallForm.osId = list.length > 0 ? list[0].id : null
  const group = osGroups.value.find((g: any) => g.id === reinstallForm.osGroupId)
  reinstallForm.port = group?.name?.toLowerCase() === 'windows' ? 3389 : 22
}

// 确认重装系统
async function confirmReinstall() {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) return

  if (!reinstallForm.osId) {
    ElMessage.warning('请选择操作系统')
    return
  }
  if (!reinstallForm.password) {
    ElMessage.warning('请设置密码')
    return
  }
  if (!reinstallForm.confirmed) {
    ElMessage.warning('请确认已备份数据')
    return
  }

  try {
    await ElMessageBox.confirm(
      '重装系统将清除当前系统所有数据，此操作不可恢复。确定继续吗？',
      '确认重装系统',
      {
        confirmButtonText: '确定重装',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    reinstallLoading.value = true
    await request.post(`/api/v2/hosts/${hostId}/reinstall`, {
      os: reinstallForm.osId,
      password: reinstallForm.password,
      port: reinstallForm.port,
      part_type: reinstallForm.partType,
      disk: reinstallForm.disk,
      action: 'reinstall'
    })
    showReinstallDialog.value = false
    ElMessage.success('重装系统任务已提交，请等待完成')
    fetchLogs()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e?.response?.data?.msg || e.message || '重装系统失败')
    }
  } finally {
    reinstallLoading.value = false
  }
}

// ==================== KVM/IPMI 功能 ====================

// 获取KVM连接信息并打开
async function openKvm() {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) return

  kvmLoading.value = true
  try {
    const { data } = await request.post(`/api/v2/hosts/${hostId}/kvm`)
    if (data?.data?.url) {
      window.open(data.data.url, '_blank', 'width=1280,height=800')
    } else if (data?.data?.jnl) {
      const link = document.createElement('a')
      link.href = data.data.jnl
      link.download = `kvm-${hostId}.jnlp`
      link.click()
    } else {
      ElMessage.info('KVM连接信息已获取，请查看弹出窗口')
    }
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.msg || e.message || '获取KVM连接失败')
  } finally {
    kvmLoading.value = false
  }
}

// 获取iKVM连接信息
async function openIkv() {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) return

  ikvmLoading.value = true
  try {
    const { data } = await request.post(`/api/v2/hosts/${hostId}/ikvm`)
    if (data?.data?.url) {
      window.open(data.data.url, '_blank', 'width=1280,height=800')
    } else {
      ElMessage.info('iKVM连接信息已获取')
    }
  } catch (e: any) {
    ElMessage.error(e?.response?.data?.msg || e.message || '获取iKVM连接失败')
  } finally {
    ikvmLoading.value = false
  }
}

// ==================== 任务队列功能 ====================

// 打开任务队列对话框
async function openTaskDialog() {
  showTaskDialog.value = true
  await fetchTasks()
  taskTimer.value = setInterval(fetchTasks, 10000)
}

// 关闭任务队列对话框
function closeTaskDialog() {
  showTaskDialog.value = false
  if (taskTimer.value) {
    clearInterval(taskTimer.value)
    taskTimer.value = null
  }
}

// 获取任务列表
async function fetchTasks() {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) return

  taskLoading.value = true
  try {
    const { data } = await request.get(`/api/v2/hosts/${hostId}/tasks`)
    if (data?.data) {
      taskList.value = data.data.list || data.data || []
      hasActiveTask.value = taskList.value.some((t: any) =>
        t.status === 'running' || t.status === 'pending' || t.status === 'processing'
      )
    }
  } catch (e) {
    console.error('Failed to fetch tasks:', e)
  } finally {
    taskLoading.value = false
  }
}

// 取消指定任务
async function cancelSpecificTask(taskId: string | number) {
  const hostId = route.query.host_id || route.params.id
  if (!hostId) return

  try {
    await ElMessageBox.confirm('确定取消该任务吗？', '取消任务', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.post(`/api/v2/hosts/${hostId}/tasks/${taskId}/cancel`)
    ElMessage.success('任务已取消')
    fetchTasks()
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e?.response?.data?.msg || e.message || '取消任务失败')
    }
  }
}

// 获取任务状态标签类型
function getTaskStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    'completed': 'success',
    'success': 'success',
    'running': 'warning',
    'processing': 'warning',
    'pending': 'info',
    'failed': 'danger',
    'cancelled': 'info',
    'canceled': 'info'
  }
  return map[status] || 'info'
}

// 获取任务状态中文
function getTaskStatusText(status: string) {
  const map: Record<string, string> = {
    'completed': '已完成',
    'success': '已完成',
    'running': '执行中',
    'processing': '处理中',
    'pending': '等待中',
    'failed': '失败',
    'cancelled': '已取消',
    'canceled': '已取消'
  }
  return map[status] || status
}

// 获取任务类型中文
function getTaskTypeText(type: string) {
  const map: Record<string, string> = {
    'reinstall': '重装系统',
    'rescue': '救援模式',
    'power_on': '开机',
    'power_off': '关机',
    'reboot': '重启',
    'reset_password': '密码重置',
    'format': '磁盘格式化',
    'bmc_reset': 'BMC重置'
  }
  return map[type] || type
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
  // 初始获取任务状态
  fetchTasks()
})

onBeforeUnmount(() => {
  if (taskTimer.value) {
    clearInterval(taskTimer.value)
    taskTimer.value = null
  }
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

.kvm-btn-group {
  display: flex;
  gap: 8px;
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
