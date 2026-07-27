<template>
  <div class="agents">
    <!-- Stats Summary -->
    <n-grid :cols="4" :x-gap="16" :y-gap="16" responsive="screen" :item-responsive="true" style="margin-bottom: 24px">
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card" :bordered="false" rounded>
          <n-statistic label="总代理数" :value="agentStats.totalAgents" />
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card active" :bordered="false" rounded>
          <n-statistic label="活跃代理" :value="agentStats.activeAgents" />
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card commission" :bordered="false" rounded>
          <n-statistic label="本月佣金" :value="agentStats.monthCommission">
            <template #prefix>¥</template>
          </n-statistic>
        </n-card>
      </n-gi>
      <n-gi span="4 m:2 l:1">
        <n-card class="stat-card pending" :bordered="false" rounded>
          <n-statistic label="待审核提现" :value="agentStats.pendingWithdrawals">
            <template #suffix>
              <span class="stat-badge">需处理</span>
            </template>
          </n-statistic>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- Tabs -->
    <n-tabs v-model:value="activeTab" type="line" animated>
      <!-- Agent List -->
      <n-tab-pane name="agents" tab="代理列表">
        <n-card :bordered="false" rounded class="table-card">
          <template #header-extra>
            <n-space>
              <n-input v-model:value="searchKeyword" placeholder="搜索代理名称" clearable style="width: 200px" />
              <n-button type="primary" @click="showAddAgent = true">
                添加代理
              </n-button>
            </n-space>
          </template>
          <n-data-table
            :columns="agentColumns"
            :data="filteredAgents"
            :bordered="false"
            :pagination="{ pageSize: 10 }"
            size="small"
          />
        </n-card>
      </n-tab-pane>

      <!-- Withdrawal Review -->
      <n-tab-pane name="withdrawals" tab="提现审核">
        <n-card :bordered="false" rounded class="table-card">
          <n-data-table
            :columns="withdrawalColumns"
            :data="withdrawalRequests"
            :bordered="false"
            :pagination="{ pageSize: 10 }"
            size="small"
          />
        </n-card>
      </n-tab-pane>

      <!-- Commission Settings -->
      <n-tab-pane name="settings" tab="佣金设置">
        <n-card :bordered="false" rounded class="settings-card">
          <n-form ref="commissionFormRef" :model="commissionSettings" label-placement="left" label-width="120">
            <n-form-item label="默认佣金比例">
              <n-input-number v-model:value="commissionSettings.defaultRate" :min="0" :max="100" :step="0.5">
                <template #suffix>%</template>
              </n-input-number>
            </n-form-item>
            <n-form-item label="一级代理佣金">
              <n-input-number v-model:value="commissionSettings.level1Rate" :min="0" :max="100" :step="0.5">
                <template #suffix>%</template>
              </n-input-number>
            </n-form-item>
            <n-form-item label="二级代理佣金">
              <n-input-number v-model:value="commissionSettings.level2Rate" :min="0" :max="100" :step="0.5">
                <template #suffix>%</template>
              </n-input-number>
            </n-form-item>
            <n-form-item label="最低提现金额">
              <n-input-number v-model:value="commissionSettings.minWithdrawal" :min="0" :step="100">
                <template #prefix>¥</template>
              </n-input-number>
            </n-form-item>
            <n-form-item label="提现手续费">
              <n-input-number v-model:value="commissionSettings.withdrawalFee" :min="0" :max="100" :step="0.5">
                <template #suffix>%</template>
              </n-input-number>
            </n-form-item>
            <n-form-item label="结算周期">
              <n-select v-model:value="commissionSettings.settlementCycle" :options="settlementOptions" />
            </n-form-item>
            <n-form-item>
              <n-button type="primary" @click="saveCommissionSettings">保存设置</n-button>
            </n-form-item>
          </n-form>
        </n-card>
      </n-tab-pane>
    </n-tabs>

    <!-- Agent Detail Modal -->
    <n-modal v-model:show="showDetail" preset="card" title="代理详情" style="width: 600px" :bordered="false">
      <n-descriptions :column="2" bordered label-placement="left">
        <n-descriptions-item label="代理名称">{{ currentAgent?.name }}</n-descriptions-item>
        <n-descriptions-item label="代理等级">
          <n-tag :type="getLevelType(currentAgent?.level)" size="small">
            {{ currentAgent?.level }}级代理
          </n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="联系人">{{ currentAgent?.contact }}</n-descriptions-item>
        <n-descriptions-item label="联系电话">{{ currentAgent?.phone }}</n-descriptions-item>
        <n-descriptions-item label="佣金比例">{{ currentAgent?.commissionRate }}%</n-descriptions-item>
        <n-descriptions-item label="累计佣金">¥{{ currentAgent?.totalCommission?.toLocaleString() }}</n-descriptions-item>
        <n-descriptions-item label="已提现">¥{{ currentAgent?.withdrawn?.toLocaleString() }}</n-descriptions-item>
        <n-descriptions-item label="可提现">¥{{ currentAgent?.available?.toLocaleString() }}</n-descriptions-item>
        <n-descriptions-item label="下线人数">{{ currentAgent?.subordinates }}人</n-descriptions-item>
        <n-descriptions-item label="状态">
          <n-tag :type="currentAgent?.status === 'active' ? 'success' : 'default'" size="small">
            {{ currentAgent?.status === 'active' ? '活跃' : '禁用' }}
          </n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="注册时间">{{ currentAgent?.createTime }}</n-descriptions-item>
        <n-descriptions-item label="最后活跃">{{ currentAgent?.lastActive }}</n-descriptions-item>
      </n-descriptions>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showDetail = false">关闭</n-button>
          <n-button type="primary" @click="editAgent(currentAgent)">编辑</n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- Add/Edit Agent Modal -->
    <n-modal v-model:show="showAddAgent" preset="card" :title="editingAgent ? '编辑代理' : '添加代理'" style="width: 500px" :bordered="false">
      <n-form ref="agentFormRef" :model="agentForm" :rules="agentRules" label-placement="left" label-width="80">
        <n-form-item label="代理名称" path="name">
          <n-input v-model:value="agentForm.name" placeholder="请输入代理名称" />
        </n-form-item>
        <n-form-item label="联系人" path="contact">
          <n-input v-model:value="agentForm.contact" placeholder="请输入联系人" />
        </n-form-item>
        <n-form-item label="联系电话" path="phone">
          <n-input v-model:value="agentForm.phone" placeholder="请输入联系电话" />
        </n-form-item>
        <n-form-item label="代理等级" path="level">
          <n-select v-model:value="agentForm.level" :options="levelOptions" />
        </n-form-item>
        <n-form-item label="佣金比例" path="commissionRate">
          <n-input-number v-model:value="agentForm.commissionRate" :min="0" :max="100" :step="0.5">
            <template #suffix>%</template>
          </n-input-number>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showAddAgent = false">取消</n-button>
          <n-button type="primary" @click="saveAgent">保存</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { h, ref, reactive, computed } from 'vue'
import {
  NTag,
  NButton,
  NSpace,
  NPopconfirm,
  type DataTableColumns,
  type FormRules,
} from 'naive-ui'

// Active tab
const activeTab = ref('agents')
const searchKeyword = ref('')
const showDetail = ref(false)
const showAddAgent = ref(false)
const currentAgent = ref<any>(null)
const editingAgent = ref<any>(null)

// Stats
const agentStats = reactive({
  totalAgents: 156,
  activeAgents: 89,
  monthCommission: 45680,
  pendingWithdrawals: 12,
})

// Agent list
const agents = ref([
  { id: 1, name: '华东代理商', level: 1, commissionRate: 15, totalCommission: 128560, withdrawn: 98000, available: 30560, subordinates: 25, status: 'active', contact: '张经理', phone: '13800138001', createTime: '2025-06-15', lastActive: '2026-07-27' },
  { id: 2, name: '华南代理商', level: 1, commissionRate: 15, totalCommission: 96420, withdrawn: 75000, available: 21420, subordinates: 18, status: 'active', contact: '李经理', phone: '13800138002', createTime: '2025-08-20', lastActive: '2026-07-26' },
  { id: 3, name: '北京分销商', level: 2, commissionRate: 10, totalCommission: 45680, withdrawn: 35000, available: 10680, subordinates: 12, status: 'active', contact: '王总', phone: '13800138003', createTime: '2025-10-10', lastActive: '2026-07-25' },
  { id: 4, name: '上海代理商', level: 1, commissionRate: 15, totalCommission: 78950, withdrawn: 60000, available: 18950, subordinates: 15, status: 'active', contact: '赵经理', phone: '13800138004', createTime: '2025-12-01', lastActive: '2026-07-24' },
  { id: 5, name: '西南分销商', level: 2, commissionRate: 10, totalCommission: 23450, withdrawn: 18000, available: 5450, subordinates: 8, status: 'inactive', contact: '孙总', phone: '13800138005', createTime: '2026-01-15', lastActive: '2026-06-20' },
  { id: 6, name: '深圳代理商', level: 1, commissionRate: 12, totalCommission: 56780, withdrawn: 45000, available: 11780, subordinates: 10, status: 'active', contact: '周经理', phone: '13800138006', createTime: '2026-02-28', lastActive: '2026-07-23' },
  { id: 7, name: '杭州分销商', level: 2, commissionRate: 10, totalCommission: 34560, withdrawn: 28000, available: 6560, subordinates: 6, status: 'active', contact: '吴总', phone: '13800138007', createTime: '2026-03-10', lastActive: '2026-07-22' },
  { id: 8, name: '成都代理商', level: 2, commissionRate: 10, totalCommission: 18920, withdrawn: 15000, available: 3920, subordinates: 5, status: 'inactive', contact: '郑经理', phone: '13800138008', createTime: '2026-04-05', lastActive: '2026-05-15' },
])

// Filtered agents
const filteredAgents = computed(() => {
  if (!searchKeyword.value) return agents.value
  return agents.value.filter((a) => a.name.includes(searchKeyword.value))
})

// Agent table columns
const agentColumns: DataTableColumns<any> = [
  { title: '代理名称', key: 'name', width: 150 },
  {
    title: '等级',
    key: 'level',
    width: 100,
    render: (row) =>
      h(NTag, { type: getLevelType(row.level), size: 'small', round: true, bordered: false }, { default: () => `${row.level}级代理` }),
  },
  {
    title: '佣金比例',
    key: 'commissionRate',
    width: 100,
    render: (row) => h('span', { style: 'font-weight:600;color:#18a058' }, `${row.commissionRate}%`),
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
    width: 120,
    render: (row) => h('span', {}, `¥${row.withdrawn.toLocaleString()}`),
  },
  { title: '下线人数', key: 'subordinates', width: 90 },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) =>
      h(
        NTag,
        { type: row.status === 'active' ? 'success' : 'default', size: 'small', round: true, bordered: false },
        { default: () => (row.status === 'active' ? '活跃' : '禁用') }
      ),
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    fixed: 'right',
    render: (row) =>
      h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, { size: 'small', type: 'info', quaternary: true, onClick: () => viewAgent(row) }, { default: () => '详情' }),
          h(NButton, { size: 'small', type: 'primary', quaternary: true, onClick: () => editAgent(row) }, { default: () => '编辑' }),
          h(
            NPopconfirm,
            { onPositiveClick: () => toggleAgentStatus(row) },
            {
              trigger: () =>
                h(NButton, { size: 'small', type: row.status === 'active' ? 'warning' : 'success', quaternary: true }, {
                  default: () => (row.status === 'active' ? '禁用' : '启用'),
                }),
              default: () => `确定${row.status === 'active' ? '禁用' : '启用'}该代理?`,
            }
          ),
        ],
      }),
  },
]

// Level type helper
function getLevelType(level: number | undefined): 'success' | 'warning' | 'info' | 'default' {
  if (level === 1) return 'success'
  if (level === 2) return 'warning'
  return 'default'
}

// View agent detail
function viewAgent(agent: any) {
  currentAgent.value = agent
  showDetail.value = true
}

// Edit agent
function editAgent(agent: any) {
  editingAgent.value = agent
  agentForm.name = agent.name
  agentForm.contact = agent.contact
  agentForm.phone = agent.phone
  agentForm.level = agent.level
  agentForm.commissionRate = agent.commissionRate
  showAddAgent.value = true
}

// Toggle agent status
function toggleAgentStatus(agent: any) {
  agent.status = agent.status === 'active' ? 'inactive' : 'active'
}

// Agent form
const agentForm = reactive({
  name: '',
  contact: '',
  phone: '',
  level: 2,
  commissionRate: 10,
})

const agentRules: FormRules = {
  name: { required: true, message: '请输入代理名称', trigger: 'blur' },
  contact: { required: true, message: '请输入联系人', trigger: 'blur' },
  phone: { required: true, message: '请输入联系电话', trigger: 'blur' },
}

const levelOptions = [
  { label: '一级代理', value: 1 },
  { label: '二级代理', value: 2 },
  { label: '三级代理', value: 3 },
]

// Save agent
function saveAgent() {
  if (editingAgent.value) {
    // Update existing agent
    Object.assign(editingAgent.value, agentForm)
  } else {
    // Add new agent
    agents.value.push({
      id: Date.now(),
      ...agentForm,
      totalCommission: 0,
      withdrawn: 0,
      available: 0,
      subordinates: 0,
      status: 'active',
      createTime: new Date().toISOString().split('T')[0],
      lastActive: new Date().toISOString().split('T')[0],
    })
  }
  showAddAgent.value = false
  editingAgent.value = null
  resetForm()
}

// Reset form
function resetForm() {
  agentForm.name = ''
  agentForm.contact = ''
  agentForm.phone = ''
  agentForm.level = 2
  agentForm.commissionRate = 10
}

// Withdrawal requests
const withdrawalRequests = ref([
  { id: 'WD20260727001', agentName: '华东代理商', amount: 5000, bankInfo: '工商银行 ****8888', status: 'pending', requestTime: '2026-07-27 09:30' },
  { id: 'WD20260727002', agentName: '华南代理商', amount: 3000, bankInfo: '建设银行 ****6666', status: 'pending', requestTime: '2026-07-27 08:45' },
  { id: 'WD20260726003', agentName: '北京分销商', amount: 2000, bankInfo: '农业银行 ****5555', status: 'approved', requestTime: '2026-07-26 16:20' },
  { id: 'WD20260726004', agentName: '上海代理商', amount: 8000, bankInfo: '中国银行 ****4444', status: 'approved', requestTime: '2026-07-26 14:15' },
  { id: 'WD20260725005', agentName: '深圳代理商', amount: 1500, bankInfo: '招商银行 ****3333', status: 'rejected', requestTime: '2026-07-25 11:30' },
])

// Withdrawal status
const withdrawalStatusMap: Record<string, { label: string; type: string }> = {
  pending: { label: '待审核', type: 'warning' },
  approved: { label: '已通过', type: 'success' },
  rejected: { label: '已拒绝', type: 'error' },
}

// Withdrawal table columns
const withdrawalColumns: DataTableColumns<any> = [
  { title: '提现单号', key: 'id', width: 150, ellipsis: { tooltip: true } },
  { title: '代理名称', key: 'agentName', width: 120 },
  {
    title: '提现金额',
    key: 'amount',
    width: 120,
    render: (row) => h('span', { style: 'font-weight:600;color:#18a058' }, `¥${row.amount.toLocaleString()}`),
  },
  { title: '银行信息', key: 'bankInfo', width: 150 },
  {
    title: '状态',
    key: 'status',
    width: 100,
    render: (row) => {
      const s = withdrawalStatusMap[row.status]
      return h(NTag, { type: s.type as any, size: 'small', round: true, bordered: false }, { default: () => s.label })
    },
  },
  { title: '申请时间', key: 'requestTime', width: 150 },
  {
    title: '操作',
    key: 'actions',
    width: 180,
    fixed: 'right',
    render: (row) => {
      if (row.status !== 'pending') return h('span', { style: 'color:#8c8c8c' }, '已处理')
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(
            NPopconfirm,
            { onPositiveClick: () => approveWithdrawal(row) },
            {
              trigger: () => h(NButton, { size: 'small', type: 'success', quaternary: true }, { default: () => '通过' }),
              default: () => `确定通过该提现申请?`,
            }
          ),
          h(
            NPopconfirm,
            { onPositiveClick: () => rejectWithdrawal(row) },
            {
              trigger: () => h(NButton, { size: 'small', type: 'error', quaternary: true }, { default: () => '拒绝' }),
              default: () => `确定拒绝该提现申请?`,
            }
          ),
        ],
      })
    },
  },
]

// Approve withdrawal
function approveWithdrawal(row: any) {
  row.status = 'approved'
}

// Reject withdrawal
function rejectWithdrawal(row: any) {
  row.status = 'rejected'
}

// Commission settings
const commissionSettings = reactive({
  defaultRate: 10,
  level1Rate: 15,
  level2Rate: 10,
  minWithdrawal: 1000,
  withdrawalFee: 1,
  settlementCycle: 'monthly',
})

const settlementOptions = [
  { label: '每周结算', value: 'weekly' },
  { label: '每两周结算', value: 'biweekly' },
  { label: '每月结算', value: 'monthly' },
]

// Save commission settings
function saveCommissionSettings() {
  // TODO: Save to backend
  window.$message?.success('佣金设置已保存')
}
</script>

<style scoped>
.agents {
  padding: 0;
}

.stat-card {
  border-radius: 12px;
  transition: all 0.3s;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.2);
}

.stat-card.active {
  background: linear-gradient(135deg, rgba(24, 160, 88, 0.15) 0%, rgba(24, 160, 88, 0.05) 100%);
  border: 1px solid rgba(24, 160, 88, 0.2);
}

.stat-card.commission {
  background: linear-gradient(135deg, rgba(32, 128, 240, 0.15) 0%, rgba(32, 128, 240, 0.05) 100%);
  border: 1px solid rgba(32, 128, 240, 0.2);
}

.stat-card.pending {
  background: linear-gradient(135deg, rgba(240, 160, 32, 0.15) 0%, rgba(240, 160, 32, 0.05) 100%);
  border: 1px solid rgba(240, 160, 32, 0.2);
}

.stat-badge {
  font-size: 12px;
  padding: 2px 8px;
  border-radius: 4px;
  background: rgba(240, 160, 32, 0.2);
  color: #f0a020;
  margin-left: 8px;
}

.table-card {
  border-radius: 12px;
}

.settings-card {
  border-radius: 12px;
  max-width: 600px;
}
</style>
