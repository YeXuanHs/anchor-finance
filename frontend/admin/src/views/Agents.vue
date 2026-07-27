<template>
  <div class="agents">
    <!-- Stats Overview -->
    <n-grid :cols="4" :x-gap="16" :y-gap="16" responsive="screen" :item-responsive="true">
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card" :bordered="false" rounded>
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">总代理数</span>
              <span class="stat-value">{{ agentStats.total }}</span>
            </div>
            <div class="stat-icon blue">
              <n-icon size="28"><PeopleIcon /></n-icon>
            </div>
          </div>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card" :bordered="false" rounded>
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">活跃代理</span>
              <span class="stat-value">{{ agentStats.active }}</span>
            </div>
            <div class="stat-icon green">
              <n-icon size="28"><FlashIcon /></n-icon>
            </div>
          </div>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card" :bordered="false" rounded>
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">累计佣金</span>
              <span class="stat-value">¥{{ agentStats.totalCommission.toLocaleString() }}</span>
            </div>
            <div class="stat-icon cyan">
              <n-icon size="28"><CashIcon /></n-icon>
            </div>
          </div>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card" :bordered="false" rounded>
          <div class="stat-content">
            <div class="stat-info">
              <span class="stat-label">待审核提现</span>
              <span class="stat-value text-orange">{{ agentStats.pendingWithdraw }}</span>
            </div>
            <div class="stat-icon orange">
              <n-icon size="28"><TimeIcon /></n-icon>
            </div>
          </div>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- Main Content -->
    <n-grid :cols="1" :x-gap="16" style="margin-top: 16px">
      <n-gi>
        <n-card :bordered="false" rounded>
          <n-tabs v-model:value="activeTab" type="line" animated>
            <!-- Agent List Tab -->
            <n-tab-pane name="list" tab="代理列表">
              <template #header-extra>
                <n-space>
                  <n-input
                    v-model:value="searchKeyword"
                    placeholder="搜索代理名称"
                    size="small"
                    clearable
                    style="width: 200px"
                  >
                    <template #prefix>
                      <n-icon size="16"><SearchIcon /></n-icon>
                    </template>
                  </n-input>
                  <n-button type="primary" size="small" @click="showAddAgent = true">
                    <template #icon><n-icon><AddIcon /></n-icon></template>
                    添加代理
                  </n-button>
                </n-space>
              </template>
              <n-data-table
                :columns="agentColumns"
                :data="filteredAgents"
                :bordered="false"
                :pagination="agentPagination"
                size="small"
              />
            </n-tab-pane>

            <!-- Commission Settings Tab -->
            <n-tab-pane name="commission" tab="佣金设置">
              <div class="commission-settings">
                <n-card title="默认佣金配置" :bordered="false" size="small">
                  <n-form label-placement="left" label-width="120px">
                    <n-form-item label="默认佣金比例">
                      <n-input-number
                        v-model:value="commissionConfig.defaultRate"
                        :min="0"
                        :max="100"
                        :step="0.5"
                        style="width: 200px"
                      >
                        <template #suffix>%</template>
                      </n-input-number>
                    </n-form-item>
                    <n-form-item label="最低提现金额">
                      <n-input-number
                        v-model:value="commissionConfig.minWithdraw"
                        :min="0"
                        :step="10"
                        style="width: 200px"
                      >
                        <template #prefix>¥</template>
                      </n-input-number>
                    </n-form-item>
                    <n-form-item label="结算周期">
                      <n-input-number
                        v-model:value="commissionConfig.settlementDays"
                        :min="1"
                        :max="90"
                        style="width: 200px"
                      >
                        <template #suffix>天</template>
                      </n-input-number>
                    </n-form-item>
                    <n-form-item label="自动结算">
                      <n-switch v-model:value="commissionConfig.autoSettle" />
                    </n-form-item>
                    <n-form-item>
                      <n-button type="primary" @click="saveCommissionConfig">保存配置</n-button>
                    </n-form-item>
                  </n-form>
                </n-card>

                <n-card title="等级佣金配置" :bordered="false" size="small" style="margin-top: 16px">
                  <n-data-table
                    :columns="levelColumns"
                    :data="levelData"
                    :bordered="false"
                    :pagination="false"
                    size="small"
                  />
                </n-card>
              </div>
            </n-tab-pane>

            <!-- Withdraw Review Tab -->
            <n-tab-pane name="withdraw" tab="提现审核">
              <n-data-table
                :columns="withdrawColumns"
                :data="withdrawList"
                :bordered="false"
                :pagination="withdrawPagination"
                size="small"
              />
            </n-tab-pane>
          </n-tabs>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- Agent Detail Modal -->
    <n-modal v-model:show="showDetail" preset="card" title="代理详情" style="width: 640px" :bordered="false">
      <div v-if="currentAgent" class="agent-detail">
        <n-grid :cols="2" :x-gap="24" :y-gap="16">
          <n-gi>
            <div class="detail-item">
              <span class="detail-label">代理名称</span>
              <span class="detail-value">{{ currentAgent.name }}</span>
            </div>
          </n-gi>
          <n-gi>
            <div class="detail-item">
              <span class="detail-label">代理等级</span>
              <n-tag :type="levelTagType(currentAgent.level)" size="small" round>
                {{ currentAgent.levelName }}
              </n-tag>
            </div>
          </n-gi>
          <n-gi>
            <div class="detail-item">
              <span class="detail-label">佣金比例</span>
              <span class="detail-value">{{ currentAgent.commissionRate }}%</span>
            </div>
          </n-gi>
          <n-gi>
            <div class="detail-item">
              <span class="detail-label">状态</span>
              <n-tag :type="currentAgent.status === 'active' ? 'success' : 'default'" size="small" round>
                {{ currentAgent.status === 'active' ? '正常' : '已禁用' }}
              </n-tag>
            </div>
          </n-gi>
          <n-gi>
            <div class="detail-item">
              <span class="detail-label">累计佣金</span>
              <span class="detail-value highlight">¥{{ currentAgent.totalCommission.toLocaleString() }}</span>
            </div>
          </n-gi>
          <n-gi>
            <div class="detail-item">
              <span class="detail-label">已提现</span>
              <span class="detail-value">¥{{ currentAgent.withdrawn.toLocaleString() }}</span>
            </div>
          </n-gi>
          <n-gi>
            <div class="detail-item">
              <span class="detail-label">可提现余额</span>
              <span class="detail-value highlight">¥{{ currentAgent.balance.toLocaleString() }}</span>
            </div>
          </n-gi>
          <n-gi>
            <div class="detail-item">
              <span class="detail-label">下线人数</span>
              <span class="detail-value">{{ currentAgent.referrals }}人</span>
            </div>
          </n-gi>
          <n-gi>
            <div class="detail-item">
              <span class="detail-label">注册时间</span>
              <span class="detail-value">{{ currentAgent.createdAt }}</span>
            </div>
          </n-gi>
          <n-gi>
            <div class="detail-item">
              <span class="detail-label">联系方式</span>
              <span class="detail-value">{{ currentAgent.contact || '-' }}</span>
            </div>
          </n-gi>
        </n-grid>
      </div>
    </n-modal>

    <!-- Add Agent Modal -->
    <n-modal v-model:show="showAddAgent" preset="card" title="添加代理" style="width: 520px" :bordered="false">
      <n-form ref="addFormRef" :model="addForm" label-placement="left" label-width="100px">
        <n-form-item label="代理名称" path="name" :rule="{ required: true, message: '请输入代理名称', trigger: 'blur' }">
          <n-input v-model:value="addForm.name" placeholder="请输入代理名称" />
        </n-form-item>
        <n-form-item label="联系邮箱" path="email">
          <n-input v-model:value="addForm.email" placeholder="请输入联系邮箱" />
        </n-form-item>
        <n-form-item label="联系电话" path="phone">
          <n-input v-model:value="addForm.phone" placeholder="请输入联系电话" />
        </n-form-item>
        <n-form-item label="代理等级" path="level">
          <n-input-number v-model:value="addForm.level" :min="1" :max="5" style="width: 100%" />
        </n-form-item>
        <n-form-item label="佣金比例" path="commissionRate" :rule="{ required: true, message: '请设置佣金比例', trigger: 'blur' }">
          <n-input-number v-model:value="addForm.commissionRate" :min="0" :max="100" :step="0.5" style="width: 100%">
            <template #suffix>%</template>
          </n-input-number>
        </n-form-item>
      </n-form>
      <template #action>
        <n-space justify="end">
          <n-button @click="showAddAgent = false">取消</n-button>
          <n-button type="primary" @click="handleAddAgent">确认添加</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { h, reactive, ref, computed } from 'vue'
import {
  PeopleOutline as PeopleIcon,
  FlashOutline as FlashIcon,
  CashOutline as CashIcon,
  TimeOutline as TimeIcon,
  SearchOutline as SearchIcon,
  AddOutline as AddIcon,
  CreateOutline as EditIcon,
  EyeOutline as ViewIcon,
} from '@vicons/ionicons5'
import { NTag, NButton, NSpace, NPopconfirm, NIcon, useMessage, type DataTableColumns } from 'naive-ui'

const message = useMessage()

// ---- Tab ----
const activeTab = ref('list')

// ---- Agent Stats ----
const agentStats = reactive({
  total: 128,
  active: 96,
  totalCommission: 586200,
  pendingWithdraw: 12,
})

// ---- Search ----
const searchKeyword = ref('')

// ---- Agent List ----
interface Agent {
  id: string
  name: string
  level: number
  levelName: string
  commissionRate: number
  totalCommission: number
  withdrawn: number
  balance: number
  referrals: number
  status: 'active' | 'disabled'
  createdAt: string
  contact: string
}

const agentList = ref<Agent[]>([
  { id: 'AGT001', name: '科技云代理', level: 1, levelName: '金牌代理', commissionRate: 15, totalCommission: 128600, withdrawn: 98000, balance: 30600, referrals: 256, status: 'active', createdAt: '2025-06-15', contact: '13800138001' },
  { id: 'AGT002', name: '星辰网络', level: 1, levelName: '金牌代理', commissionRate: 15, totalCommission: 96800, withdrawn: 82000, balance: 14800, referrals: 198, status: 'active', createdAt: '2025-08-20', contact: '13800138002' },
  { id: 'AGT003', name: '云端数据', level: 2, levelName: '银牌代理', commissionRate: 12, totalCommission: 68500, withdrawn: 55000, balance: 13500, referrals: 142, status: 'active', createdAt: '2025-10-08', contact: '13800138003' },
  { id: 'AGT004', name: '极速互联', level: 2, levelName: '银牌代理', commissionRate: 12, totalCommission: 52300, withdrawn: 42000, balance: 10300, referrals: 108, status: 'active', createdAt: '2025-11-22', contact: '13800138004' },
  { id: 'AGT005', name: '蓝海信息', level: 3, levelName: '铜牌代理', commissionRate: 8, totalCommission: 35200, withdrawn: 28000, balance: 7200, referrals: 76, status: 'active', createdAt: '2026-01-10', contact: '13800138005' },
  { id: 'AGT006', name: '智慧云服', level: 3, levelName: '铜牌代理', commissionRate: 8, totalCommission: 28600, withdrawn: 22000, balance: 6600, referrals: 58, status: 'active', createdAt: '2026-02-18', contact: '13800138006' },
  { id: 'AGT007', name: '数据先锋', level: 3, levelName: '铜牌代理', commissionRate: 8, totalCommission: 18900, withdrawn: 15000, balance: 3900, referrals: 42, status: 'disabled', createdAt: '2026-03-05', contact: '13800138007' },
  { id: 'AGT008', name: '网络驿站', level: 2, levelName: '银牌代理', commissionRate: 12, totalCommission: 45800, withdrawn: 38000, balance: 7800, referrals: 95, status: 'active', createdAt: '2026-04-12', contact: '13800138008' },
  { id: 'AGT009', name: '创想科技', level: 1, levelName: '金牌代理', commissionRate: 15, totalCommission: 112800, withdrawn: 92000, balance: 20800, referrals: 230, status: 'active', createdAt: '2025-07-28', contact: '13800138009' },
])

const filteredAgents = computed(() => {
  if (!searchKeyword.value) return agentList.value
  return agentList.value.filter((a) =>
    a.name.toLowerCase().includes(searchKeyword.value.toLowerCase())
  )
})

const agentPagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => { agentPagination.page = page },
  onUpdatePageSize: (pageSize: number) => {
    agentPagination.pageSize = pageSize
    agentPagination.page = 1
  },
})

function levelTagType(level: number): 'warning' | 'info' | 'success' | 'default' {
  const map: Record<number, 'warning' | 'info' | 'success' | 'default'> = {
    1: 'warning',
    2: 'info',
    3: 'success',
    4: 'default',
    5: 'default',
  }
  return map[level] || 'default'
}

// ---- Agent Detail ----
const showDetail = ref(false)
const currentAgent = ref<Agent | null>(null)

function viewAgent(agent: Agent) {
  currentAgent.value = agent
  showDetail.value = true
}

// ---- Agent Columns ----
const agentColumns: DataTableColumns<Agent> = [
  { title: 'ID', key: 'id', width: 100, ellipsis: { tooltip: true } },
  { title: '代理名称', key: 'name', width: 140 },
  {
    title: '等级',
    key: 'level',
    width: 100,
    render: (row) =>
      h(NTag, { type: levelTagType(row.level), size: 'small', round: true, bordered: false }, { default: () => row.levelName }),
  },
  {
    title: '佣金比例',
    key: 'commissionRate',
    width: 100,
    render: (row) => h('span', { style: 'font-weight:600;color:#1890ff' }, `${row.commissionRate}%`),
  },
  {
    title: '累计佣金',
    key: 'totalCommission',
    width: 120,
    render: (row) => h('span', { style: 'font-weight:600' }, `¥${row.totalCommission.toLocaleString()}`),
  },
  {
    title: '已提现',
    key: 'withdrawn',
    width: 110,
    render: (row) => h('span', null, `¥${row.withdrawn.toLocaleString()}`),
  },
  { title: '下线人数', key: 'referrals', width: 90 },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) =>
      h(NTag, { type: row.status === 'active' ? 'success' : 'default', size: 'small', round: true, bordered: false }, { default: () => (row.status === 'active' ? '正常' : '已禁用') }),
  },
  {
    title: '操作',
    key: 'action',
    width: 150,
    render: (row) =>
      h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, { text: true, type: 'primary', size: 'small', onClick: () => viewAgent(row) }, {
            icon: () => h(NIcon, { size: 14 }, { default: () => h(ViewIcon) }),
            default: () => '详情',
          }),
          h(NPopconfirm, { onPositiveClick: () => toggleAgentStatus(row) }, {
            trigger: () =>
              h(NButton, { text: true, type: row.status === 'active' ? 'warning' : 'success', size: 'small' }, {
                default: () => (row.status === 'active' ? '禁用' : '启用'),
              }),
            default: () => `确认${row.status === 'active' ? '禁用' : '启用'}该代理？`,
          }),
        ],
      }),
  },
]

function toggleAgentStatus(agent: Agent) {
  agent.status = agent.status === 'active' ? 'disabled' : 'active'
  message.success(`已${agent.status === 'active' ? '启用' : '禁用'}代理：${agent.name}`)
}

// ---- Add Agent ----
const showAddAgent = ref(false)
const addFormRef = ref()
const addForm = reactive({
  name: '',
  email: '',
  phone: '',
  level: 3,
  commissionRate: 8,
})

function handleAddAgent() {
  if (!addForm.name) {
    message.warning('请输入代理名称')
    return
  }
  const levelNames: Record<number, string> = { 1: '金牌代理', 2: '银牌代理', 3: '铜牌代理', 4: '普通代理', 5: '初级代理' }
  agentList.value.unshift({
    id: `AGT${String(agentList.value.length + 1).padStart(3, '0')}`,
    name: addForm.name,
    level: addForm.level,
    levelName: levelNames[addForm.level] || '普通代理',
    commissionRate: addForm.commissionRate,
    totalCommission: 0,
    withdrawn: 0,
    balance: 0,
    referrals: 0,
    status: 'active',
    createdAt: new Date().toISOString().split('T')[0],
    contact: addForm.phone || addForm.email,
  })
  message.success(`代理 "${addForm.name}" 添加成功`)
  showAddAgent.value = false
  addForm.name = ''
  addForm.email = ''
  addForm.phone = ''
  addForm.level = 3
  addForm.commissionRate = 8
}

// ---- Commission Config ----
const commissionConfig = reactive({
  defaultRate: 8,
  minWithdraw: 100,
  settlementDays: 30,
  autoSettle: true,
})

function saveCommissionConfig() {
  message.success('佣金配置已保存')
}

interface LevelConfig {
  level: number
  name: string
  commissionRate: string
  minReferrals: number
}

const levelData = ref<LevelConfig[]>([
  { level: 1, name: '金牌代理', commissionRate: '15%', minReferrals: 200 },
  { level: 2, name: '银牌代理', commissionRate: '12%', minReferrals: 100 },
  { level: 3, name: '铜牌代理', commissionRate: '8%', minReferrals: 0 },
])

const levelColumns: DataTableColumns<LevelConfig> = [
  { title: '等级', key: 'level', width: 80, render: (row) => h('span', null, `Lv.${row.level}`) },
  { title: '名称', key: 'name', width: 120 },
  { title: '佣金比例', key: 'commissionRate', width: 120 },
  { title: '最低下线数', key: 'minReferrals', width: 120 },
  {
    title: '操作',
    key: 'action',
    width: 80,
    render: () =>
      h(NButton, { text: true, type: 'primary', size: 'small' }, {
        icon: () => h(NIcon, { size: 14 }, { default: () => h(EditIcon) }),
        default: () => '编辑',
      }),
  },
]

// ---- Withdraw Review ----
interface WithdrawRequest {
  id: string
  agentName: string
  amount: number
  method: string
  account: string
  status: 'pending' | 'approved' | 'rejected'
  createdAt: string
}

const withdrawList = ref<WithdrawRequest[]>([
  { id: 'WD001', agentName: '科技云代理', amount: 5000, method: '支付宝', account: '138****8001', status: 'pending', createdAt: '2026-07-27 08:30' },
  { id: 'WD002', agentName: '星辰网络', amount: 3000, method: '银行卡', account: '6222****5678', status: 'pending', createdAt: '2026-07-27 07:15' },
  { id: 'WD003', agentName: '创想科技', amount: 8000, method: '微信', account: 'cx****@wx', status: 'pending', createdAt: '2026-07-26 18:45' },
  { id: 'WD004', agentName: '云端数据', amount: 2000, method: '支付宝', account: '138****8003', status: 'approved', createdAt: '2026-07-26 14:20' },
  { id: 'WD005', agentName: '极速互联', amount: 1500, method: '银行卡', account: '6222****9012', status: 'rejected', createdAt: '2026-07-25 16:30' },
  { id: 'WD006', agentName: '蓝海信息', amount: 4000, method: '支付宝', account: '138****8005', status: 'approved', createdAt: '2026-07-25 10:15' },
])

const withdrawPagination = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  onChange: (page: number) => { withdrawPagination.page = page },
  onUpdatePageSize: (pageSize: number) => {
    withdrawPagination.pageSize = pageSize
    withdrawPagination.page = 1
  },
})

const statusTagMap: Record<string, { label: string; type: 'warning' | 'success' | 'error' }> = {
  pending: { label: '待审核', type: 'warning' },
  approved: { label: '已通过', type: 'success' },
  rejected: { label: '已拒绝', type: 'error' },
}

function handleWithdraw(id: string, action: 'approved' | 'rejected') {
  const item = withdrawList.value.find((w) => w.id === id)
  if (item) {
    item.status = action
    message.success(`提现已${action === 'approved' ? '通过' : '拒绝'}`)
  }
}

const withdrawColumns: DataTableColumns<WithdrawRequest> = [
  { title: '申请ID', key: 'id', width: 100 },
  { title: '代理名称', key: 'agentName', width: 120 },
  {
    title: '提现金额',
    key: 'amount',
    width: 120,
    render: (row) => h('span', { style: 'font-weight:600;color:#1890ff' }, `¥${row.amount.toLocaleString()}`),
  },
  { title: '提现方式', key: 'method', width: 100 },
  { title: '收款账户', key: 'account', width: 140 },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render: (row) => {
      const s = statusTagMap[row.status]
      return h(NTag, { type: s.type, size: 'small', round: true, bordered: false }, { default: () => s.label })
    },
  },
  { title: '申请时间', key: 'createdAt', width: 155 },
  {
    title: '操作',
    key: 'action',
    width: 140,
    render: (row) => {
      if (row.status !== 'pending') {
        return h('span', { style: 'color: #bfbfbf; font-size: 12px' }, '已处理')
      }
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(NPopconfirm, { onPositiveClick: () => handleWithdraw(row.id, 'approved') }, {
            trigger: () => h(NButton, { text: true, type: 'success', size: 'small' }, { default: () => '通过' }),
            default: () => '确认通过该提现申请？',
          }),
          h(NPopconfirm, { onPositiveClick: () => handleWithdraw(row.id, 'rejected') }, {
            trigger: () => h(NButton, { text: true, type: 'error', size: 'small' }, { default: () => '拒绝' }),
            default: () => '确认拒绝该提现申请？',
          }),
        ],
      })
    },
  },
]
</script>

<style scoped>
.agents {
  padding: 0;
}

.stat-card {
  border-radius: 12px;
  transition: box-shadow 0.3s, transform 0.2s;
}

.stat-card:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.08);
  transform: translateY(-2px);
}

.stat-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.stat-info {
  display: flex;
  flex-direction: column;
}

.stat-label {
  font-size: 13px;
  color: #8c8c8c;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #1a1a2e;
  line-height: 1.2;
}

.stat-value.text-orange {
  color: #fa8c16;
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon.blue {
  background: rgba(24, 144, 255, 0.1);
  color: #1890ff;
}

.stat-icon.cyan {
  background: rgba(19, 194, 194, 0.1);
  color: #13c2c2;
}

.stat-icon.green {
  background: rgba(82, 196, 26, 0.1);
  color: #52c41a;
}

.stat-icon.orange {
  background: rgba(250, 140, 22, 0.1);
  color: #fa8c16;
}

.commission-settings {
  max-width: 800px;
}

.agent-detail {
  padding: 8px 0;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-label {
  font-size: 12px;
  color: #8c8c8c;
}

.detail-value {
  font-size: 14px;
  color: #1a1a2e;
  font-weight: 500;
}

.detail-value.highlight {
  color: #1890ff;
  font-weight: 700;
  font-size: 16px;
}
</style>
