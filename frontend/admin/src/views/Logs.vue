<template>
  <n-card :bordered="false" rounded>
    <!-- 筛选栏 -->
    <n-space justify="space-between" style="margin-bottom: 16px">
      <n-space>
        <n-select
          v-model:value="filters.type"
          :options="logTypeOptions"
          placeholder="日志类型"
          clearable
          style="width: 150px"
        />
        <n-date-picker
          v-model:value="filters.dateRange"
          type="daterange"
          clearable
          style="width: 300px"
        />
        <n-input
          v-model:value="filters.keyword"
          placeholder="搜索关键词"
          clearable
          style="width: 200px"
          @keyup.enter="handleSearch"
        />
        <n-button @click="handleSearch">搜索</n-button>
        <n-button @click="resetFilters">重置</n-button>
      </n-space>
      <n-button type="primary" @click="handleExport">导出日志</n-button>
    </n-space>

    <!-- 日志表格 -->
    <n-data-table :columns="logColumns" :data="filteredLogs" :bordered="false" :pagination="pagination" />
  </n-card>

  <!-- 日志详情抽屉 -->
  <n-modal v-model:show="showDetail" preset="card" title="日志详情" style="width: 600px">
    <n-descriptions :column="1" label-placement="left" bordered v-if="currentLog">
      <n-descriptions-item label="日志ID">{{ currentLog.id }}</n-descriptions-item>
      <n-descriptions-item label="时间">{{ currentLog.time }}</n-descriptions-item>
      <n-descriptions-item label="日志类型">
        <n-tag :type="logTypeTagMap[currentLog.type]" size="small" :bordered="false">
          {{ logTypeMap[currentLog.type] }}
        </n-tag>
      </n-descriptions-item>
      <n-descriptions-item label="操作类型">{{ currentLog.action }}</n-descriptions-item>
      <n-descriptions-item label="操作人">{{ currentLog.operator }}</n-descriptions-item>
      <n-descriptions-item label="IP 地址">{{ currentLog.ip }}</n-descriptions-item>
      <n-descriptions-item label="状态">
        <n-tag :type="currentLog.status === 'success' ? 'success' : 'error'" size="small" :bordered="false">
          {{ currentLog.status === 'success' ? '成功' : '失败' }}
        </n-tag>
      </n-descriptions-item>
      <n-descriptions-item label="描述">{{ currentLog.description }}</n-descriptions-item>
      <n-descriptions-item label="User Agent">{{ currentLog.userAgent }}</n-descriptions-item>
    </n-descriptions>
  </n-modal>
</template>

<script setup lang="ts">
import { ref, reactive, computed, h } from 'vue'
import { useMessage, NButton, NTag, NSpace } from 'naive-ui'
import type { DataTableColumns, PaginationProps } from 'naive-ui'

const message = useMessage()
const showDetail = ref(false)
const currentLog = ref<LogEntry | null>(null)

// ---- Types ----
interface LogEntry {
  id: number
  time: string
  type: 'operation' | 'login' | 'error'
  action: string
  operator: string
  ip: string
  description: string
  status: 'success' | 'failed'
  userAgent: string
}

// ---- Options ----
const logTypeOptions = [
  { label: '操作日志', value: 'operation' },
  { label: '登录日志', value: 'login' },
  { label: '错误日志', value: 'error' },
]

const logTypeMap: Record<string, string> = {
  operation: '操作日志',
  login: '登录日志',
  error: '错误日志',
}

const logTypeTagMap: Record<string, 'info' | 'warning' | 'error' | 'success'> = {
  operation: 'info',
  login: 'success',
  error: 'error',
}

// ---- Mock Data ----
const logs = ref<LogEntry[]>([
  { id: 1001, time: '2026-07-27 09:55:12', type: 'login', action: '用户登录', operator: 'admin', ip: '192.168.1.100', description: '管理员登录成功', status: 'success', userAgent: 'Chrome/126.0 Windows 10' },
  { id: 1002, time: '2026-07-27 09:50:03', type: 'operation', action: '修改产品', operator: 'admin', ip: '192.168.1.100', description: '修改了产品 "基础版" 的价格', status: 'success', userAgent: 'Chrome/126.0 Windows 10' },
  { id: 1003, time: '2026-07-27 09:45:22', type: 'operation', action: '创建订单', operator: '张三', ip: '10.0.0.55', description: '创建订单 ORD-20260727-001', status: 'success', userAgent: 'Safari/18.0 macOS' },
  { id: 1004, time: '2026-07-27 09:40:11', type: 'error', action: '支付回调', operator: '系统', ip: '127.0.0.1', description: '支付宝回调签名验证失败: invalid sign', status: 'failed', userAgent: 'Alipay-Notify/2.0' },
  { id: 1005, time: '2026-07-27 09:35:00', type: 'login', action: '用户登录', operator: '李四', ip: '10.0.0.88', description: '用户登录失败: 密码错误', status: 'failed', userAgent: 'Firefox/128.0 Windows 11' },
  { id: 1006, time: '2026-07-27 09:30:45', type: 'operation', action: '删除优惠券', operator: 'admin', ip: '192.168.1.100', description: '删除了过期优惠券 SUMMER2026', status: 'success', userAgent: 'Chrome/126.0 Windows 10' },
  { id: 1007, time: '2026-07-27 09:25:18', type: 'error', action: '邮件发送', operator: '系统', ip: '127.0.0.1', description: 'SMTP 连接超时: smtp.example.com:465', status: 'failed', userAgent: 'AnchorFinance-Worker/1.0' },
  { id: 1008, time: '2026-07-27 09:20:33', type: 'operation', action: '导出数据', operator: 'admin', ip: '192.168.1.100', description: '导出用户列表 CSV 共 1,523 条记录', status: 'success', userAgent: 'Chrome/126.0 Windows 10' },
  { id: 1009, time: '2026-07-27 09:15:09', type: 'login', action: '用户登录', operator: '王五', ip: '172.16.0.12', description: '用户通过微信 OAuth 登录成功', status: 'success', userAgent: 'WeChat/8.0 iOS 17' },
  { id: 1010, time: '2026-07-27 09:10:44', type: 'operation', action: '修改设置', operator: 'admin', ip: '192.168.1.100', description: '更新了邮件 SMTP 配置', status: 'success', userAgent: 'Chrome/126.0 Windows 10' },
  { id: 1011, time: '2026-07-27 09:05:21', type: 'error', action: 'API 请求', operator: '系统', ip: '127.0.0.1', description: '第三方 API rate limit exceeded: github.com', status: 'failed', userAgent: 'AnchorFinance-Worker/1.0' },
  { id: 1012, time: '2026-07-27 09:00:00', type: 'operation', action: '系统备份', operator: '系统', ip: '127.0.0.1', description: '每日数据库自动备份完成，大小 245MB', status: 'success', userAgent: 'AnchorFinance-Cron/1.0' },
])

// ---- Filters ----
const filters = reactive({
  type: null as string | null,
  dateRange: null as [number, number] | null,
  keyword: '',
})

const filteredLogs = computed(() => {
  let result = logs.value
  if (filters.type) {
    result = result.filter((l) => l.type === filters.type)
  }
  if (filters.keyword) {
    const kw = filters.keyword.toLowerCase()
    result = result.filter(
      (l) =>
        l.action.toLowerCase().includes(kw) ||
        l.operator.toLowerCase().includes(kw) ||
        l.description.toLowerCase().includes(kw) ||
        l.ip.includes(kw)
    )
  }
  return result
})

const pagination: PaginationProps = reactive({
  page: 1,
  pageSize: 10,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  itemCount: computed(() => filteredLogs.value.length),
  prefix: ({ itemCount }: { itemCount: number | undefined }) => `共 ${itemCount ?? 0} 条`,
})

// ---- Table Columns ----
const logColumns: DataTableColumns<LogEntry> = [
  { title: 'ID', key: 'id', width: 80 },
  { title: '时间', key: 'time', width: 180 },
  {
    title: '类型',
    key: 'type',
    width: 100,
    render: (row) =>
      h(NTag, { type: logTypeTagMap[row.type], size: 'small', bordered: false }, () => logTypeMap[row.type]),
  },
  { title: '操作', key: 'action', width: 120 },
  { title: '操作人', key: 'operator', width: 100 },
  { title: 'IP 地址', key: 'ip', width: 140 },
  { title: '描述', key: 'description', ellipsis: { tooltip: true } },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: (row) =>
      h(NTag, { type: row.status === 'success' ? 'success' : 'error', size: 'small', bordered: false }, () =>
        row.status === 'success' ? '成功' : '失败'
      ),
  },
  {
    title: '操作',
    key: 'actions',
    width: 100,
    render: (row) =>
      h(NButton, { size: 'small', quaternary: true, type: 'info', onClick: () => viewDetail(row) }, () => '详情'),
  },
]

// ---- Handlers ----
function handleSearch() {
  pagination.page = 1
}

function resetFilters() {
  filters.type = null
  filters.dateRange = null
  filters.keyword = ''
  pagination.page = 1
}

function viewDetail(row: LogEntry) {
  currentLog.value = row
  showDetail.value = true
}

function handleExport() {
  message.loading('正在生成日志导出文件...')
  setTimeout(() => {
    message.success('日志导出成功 (logs_export_20260727.csv)')
  }, 1500)
}
</script>

<style scoped>
.n-card {
  border-radius: 12px;
}
</style>
