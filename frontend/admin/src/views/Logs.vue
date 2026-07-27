<template>
  <n-card :bordered="false" rounded>
    <!-- 筛选栏 -->
    <n-space justify="space-between" align="center" style="margin-bottom: 16px; flex-wrap: wrap; gap: 12px">
      <n-space style="flex-wrap: wrap; gap: 12px">
        <n-select
          v-model:value="filterType"
          :options="logTypeOptions"
          placeholder="日志类型"
          clearable
          style="width: 150px"
        />
        <n-date-picker
          v-model:value="dateRange"
          type="daterange"
          clearable
          style="width: 280px"
        />
        <n-input
          v-model:value="searchKeyword"
          placeholder="搜索关键词"
          clearable
          style="width: 200px"
        >
          <template #prefix>
            <n-icon><SearchOutline /></n-icon>
          </template>
        </n-input>
      </n-space>
      <n-space>
        <n-button @click="handleExport">
          <template #icon><n-icon><DownloadOutline /></n-icon></template>
          导出日志
        </n-button>
        <n-button type="error" @click="handleClearLogs">
          清空日志
        </n-button>
      </n-space>
    </n-space>

    <!-- 日志表格 -->
    <n-data-table
      :columns="columns"
      :data="filteredLogs"
      :bordered="false"
      :single-line="false"
      :pagination="pagination"
      striped
    />
  </n-card>

  <!-- 日志详情对话框 -->
  <n-modal
    v-model:show="showDetail"
    title="日志详情"
    preset="card"
    style="max-width: 640px"
    :bordered="false"
  >
    <n-descriptions :column="1" label-placement="left" bordered size="small">
      <n-descriptions-item label="日志ID">{{ detailLog?.id }}</n-descriptions-item>
      <n-descriptions-item label="时间">{{ detailLog?.time }}</n-descriptions-item>
      <n-descriptions-item label="操作类型">
        <n-tag :type="getTypeTag(detailLog?.type || '')" size="small" :bordered="false">
          {{ detailLog?.typeLabel }}
        </n-tag>
      </n-descriptions-item>
      <n-descriptions-item label="操作人">{{ detailLog?.operator }}</n-descriptions-item>
      <n-descriptions-item label="IP地址">{{ detailLog?.ip }}</n-descriptions-item>
      <n-descriptions-item label="状态">
        <n-tag :type="detailLog?.status === 'success' ? 'success' : 'error'" size="small" :bordered="false">
          {{ detailLog?.status === 'success' ? '成功' : '失败' }}
        </n-tag>
      </n-descriptions-item>
      <n-descriptions-item label="User Agent">{{ detailLog?.userAgent }}</n-descriptions-item>
      <n-descriptions-item label="描述">{{ detailLog?.description }}</n-descriptions-item>
      <n-descriptions-item label="详细信息">
        <n-log :log="detailLog?.detail || ''" language="log" :rows="6" />
      </n-descriptions-item>
    </n-descriptions>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, computed, h } from 'vue'
import { useMessage, NTag, NButton, NPopconfirm, NSpace } from 'naive-ui'
import { SearchOutline, DownloadOutline } from '@vicons/ionicons5'
import type { DataTableColumns, PaginationProps } from 'naive-ui'

const message = useMessage()
const filterType = ref<string | null>(null)
const dateRange = ref<[number, number] | null>(null)
const searchKeyword = ref('')
const showDetail = ref(false)
const detailLog = ref<LogEntry | null>(null)

const logTypeOptions = [
  { label: '操作日志', value: 'operation' },
  { label: '登录日志', value: 'login' },
  { label: '错误日志', value: 'error' },
]

const typeLabelMap: Record<string, string> = {
  operation: '操作日志',
  login: '登录日志',
  error: '错误日志',
}

interface LogEntry {
  id: string
  time: string
  type: string
  typeLabel: string
  operator: string
  ip: string
  description: string
  status: 'success' | 'error'
  userAgent: string
  detail: string
}

const mockLogs = ref<LogEntry[]>([
  {
    id: '10001',
    time: '2026-07-27 10:15:32',
    type: 'login',
    typeLabel: '登录日志',
    operator: 'admin',
    ip: '192.168.1.100',
    description: '管理员登录系统',
    status: 'success',
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/126.0',
    detail: 'Login successful for user admin from IP 192.168.1.100',
  },
  {
    id: '10002',
    time: '2026-07-27 09:45:18',
    type: 'operation',
    typeLabel: '操作日志',
    operator: 'admin',
    ip: '192.168.1.100',
    description: '修改系统基本设置',
    status: 'success',
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/126.0',
    detail: 'Updated settings: siteName, siteDescription, timezone',
  },
  {
    id: '10003',
    time: '2026-07-27 09:30:05',
    type: 'error',
    typeLabel: '错误日志',
    operator: 'system',
    ip: '127.0.0.1',
    description: '邮件发送失败：SMTP连接超时',
    status: 'error',
    userAgent: 'AnchorFinance/1.0',
    detail: 'Error: SMTP connection timeout after 30s\nHost: smtp.example.com:465\nRetry count: 3',
  },
  {
    id: '10004',
    time: '2026-07-27 09:12:44',
    type: 'operation',
    typeLabel: '操作日志',
    operator: 'admin',
    ip: '192.168.1.100',
    description: '新增OAuth提供商：GitHub',
    status: 'success',
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/126.0',
    detail: 'Added OAuth provider: github\nAppID: Iv1.abc123def456\nCallback: /oauth/github/callback',
  },
  {
    id: '10005',
    time: '2026-07-27 08:55:21',
    type: 'login',
    typeLabel: '登录日志',
    operator: 'user_zhang',
    ip: '10.0.0.55',
    description: '用户登录失败：密码错误',
    status: 'error',
    userAgent: 'Mozilla/5.0 (Macintosh) Safari/605.1',
    detail: 'Login failed for user_zhang: incorrect password\nAttempt: 2/5',
  },
  {
    id: '10006',
    time: '2026-07-27 08:30:00',
    type: 'operation',
    typeLabel: '操作日志',
    operator: 'admin',
    ip: '192.168.1.100',
    description: '导出系统日志报表',
    status: 'success',
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/126.0',
    detail: 'Exported logs from 2026-07-01 to 2026-07-27\nFormat: CSV\nTotal records: 1523',
  },
  {
    id: '10007',
    time: '2026-07-26 23:59:01',
    type: 'error',
    typeLabel: '错误日志',
    operator: 'system',
    ip: '127.0.0.1',
    description: '数据库备份失败：磁盘空间不足',
    status: 'error',
    userAgent: 'AnchorFinance/1.0 CronJob',
    detail: 'Backup failed: insufficient disk space\nRequired: 2.5GB\nAvailable: 1.2GB\nPath: /backups/db/',
  },
  {
    id: '10008',
    time: '2026-07-26 18:20:33',
    type: 'operation',
    typeLabel: '操作日志',
    operator: 'editor_li',
    ip: '10.0.0.88',
    description: '发布财务报告：2026年Q2',
    status: 'success',
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Edge/126.0',
    detail: 'Published report: Q2-2026-Financial-Report.pdf\nCategory: Financial Reports\nVisibility: public',
  },
  {
    id: '10009',
    time: '2026-07-26 15:10:12',
    type: 'login',
    typeLabel: '登录日志',
    operator: 'admin',
    ip: '192.168.1.100',
    description: '管理员登录系统',
    status: 'success',
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/126.0',
    detail: 'Login successful for user admin from IP 192.168.1.100\n2FA verified: true',
  },
  {
    id: '10010',
    time: '2026-07-26 14:05:48',
    type: 'operation',
    typeLabel: '操作日志',
    operator: 'admin',
    ip: '192.168.1.100',
    description: '修改邮件SMTP配置',
    status: 'success',
    userAgent: 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) Chrome/126.0',
    detail: 'Updated SMTP settings\nHost: smtp.example.com\nPort: 465\nEncryption: SSL',
  },
])

const filteredLogs = computed(() => {
  let logs = [...mockLogs.value]
  if (filterType.value) {
    logs = logs.filter((l) => l.type === filterType.value)
  }
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    logs = logs.filter(
      (l) =>
        l.description.toLowerCase().includes(kw) ||
        l.operator.toLowerCase().includes(kw) ||
        l.ip.includes(kw)
    )
  }
  return logs
})

const pagination: PaginationProps = {
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  itemCount: computed(() => filteredLogs.value.length) as unknown as number,
}

function getTypeTag(type: string): 'success' | 'warning' | 'error' | 'info' | 'default' {
  const map: Record<string, 'success' | 'warning' | 'error' | 'info' | 'default'> = {
    operation: 'info',
    login: 'success',
    error: 'error',
  }
  return map[type] || 'default'
}

const columns: DataTableColumns<LogEntry> = [
  {
    title: '时间',
    key: 'time',
    width: 180,
    sorter: (a, b) => a.time.localeCompare(b.time),
  },
  {
    title: '类型',
    key: 'type',
    width: 110,
    render(row) {
      return h(NTag, {
        type: getTypeTag(row.type),
        size: 'small',
        bordered: false,
      }, { default: () => row.typeLabel })
    },
  },
  {
    title: '操作人',
    key: 'operator',
    width: 120,
  },
  {
    title: 'IP地址',
    key: 'ip',
    width: 140,
  },
  {
    title: '描述',
    key: 'description',
    ellipsis: { tooltip: true },
  },
  {
    title: '状态',
    key: 'status',
    width: 90,
    render(row) {
      return h(NTag, {
        type: row.status === 'success' ? 'success' : 'error',
        size: 'small',
        bordered: false,
      }, { default: () => row.status === 'success' ? '成功' : '失败' })
    },
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    render(row) {
      return h(NSpace, { size: 'small' }, {
        default: () => [
          h(NButton, {
            size: 'small',
            type: 'info',
            onClick: () => viewDetail(row),
          }, { default: () => '详情' }),
          h(NPopconfirm, {
            onPositiveClick: () => deleteLog(row.id),
          }, {
            trigger: () => h(NButton, { size: 'small', type: 'error' }, { default: () => '删除' }),
            default: () => '确定删除此日志？',
          }),
        ],
      })
    },
  },
]

function viewDetail(row: LogEntry) {
  detailLog.value = row
  showDetail.value = true
}

function deleteLog(id: string) {
  mockLogs.value = mockLogs.value.filter((l) => l.id !== id)
  message.success('日志已删除')
}

function handleExport() {
  // TODO: real export
  message.success('日志导出中，请稍候...')
}

function handleClearLogs() {
  // TODO: real clear with confirm
  mockLogs.value = []
  message.success('日志已清空')
}
</script>

<style scoped>
.n-card {
  border-radius: 12px;
}
</style>
