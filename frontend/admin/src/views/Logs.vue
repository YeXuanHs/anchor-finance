<template>
  <div class="admin-page">
    <el-card shadow="hover">
      <template #header>
        <div class="card-header">
          <span class="card-title">系统日志</span>
          <div class="card-actions">
            <el-select v-model="filterType" placeholder="日志类型" clearable style="width: 130px">
              <el-option v-for="o in logTypeOptions" :key="o.value" :label="o.label" :value="o.value" />
            </el-select>
            <el-date-picker v-model="dateRange" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" style="width: 260px" clearable />
            <el-input v-model="searchKeyword" placeholder="搜索关键词" clearable style="width: 180px">
              <template #prefix><el-icon><Search /></el-icon></template>
            </el-input>
            <el-button @click="handleExport"><el-icon><Download /></el-icon>导出日志</el-button>
            <el-button type="danger" @click="handleClearLogs">清空日志</el-button>
          </div>
        </div>
      </template>

      <el-table :data="filteredLogs" stripe size="small">
        <el-table-column prop="time" label="时间" width="180" sortable />
        <el-table-column prop="type" label="类型" width="110">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">{{ row.typeLabel }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="operator" label="操作人" width="120" />
        <el-table-column prop="ip" label="IP地址" width="140" />
        <el-table-column prop="description" label="描述" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="130" fixed="right">
          <template #default="{ row }">
            <el-button text type="primary" @click="viewDetail(row)">详情</el-button>
            <el-popconfirm title="确定删除此日志？" @confirm="deleteLog(row.id)">
              <template #reference>
                <el-button text type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="filteredLogs.length"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
        />
      </div>
    </el-card>

    <!-- Log Detail Dialog -->
    <el-dialog v-model="showDetail" title="日志详情" width="640px">
      <el-descriptions :column="1" border size="small">
        <el-descriptions-item label="日志ID">{{ detailLog?.id }}</el-descriptions-item>
        <el-descriptions-item label="时间">{{ detailLog?.time }}</el-descriptions-item>
        <el-descriptions-item label="操作类型">
          <el-tag :type="getTypeTag(detailLog?.type || '')" size="small">{{ detailLog?.typeLabel }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="操作人">{{ detailLog?.operator }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ detailLog?.ip }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="detailLog?.status === 'success' ? 'success' : 'danger'" size="small">
            {{ detailLog?.status === 'success' ? '成功' : '失败' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="User Agent">{{ detailLog?.userAgent }}</el-descriptions-item>
        <el-descriptions-item label="描述">{{ detailLog?.description }}</el-descriptions-item>
        <el-descriptions-item label="详细信息">
          <pre style="background: #f5f7fa; padding: 12px; border-radius: 6px; font-size: 12px; overflow-x: auto">{{ detailLog?.detail }}</pre>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Search, Download } from '@element-plus/icons-vue'

const filterType = ref<string | null>(null)
const dateRange = ref<any>(null)
const searchKeyword = ref('')
const showDetail = ref(false)
const detailLog = ref<any>(null)
const currentPage = ref(1)
const pageSize = ref(10)

const logTypeOptions = [
  { label: '操作日志', value: 'operation' },
  { label: '登录日志', value: 'login' },
  { label: '错误日志', value: 'error' },
]

function getTypeTag(type: string) {
  const map: Record<string, string> = { operation: 'info', login: 'success', error: 'danger' }
  return (map[type] || 'info') as any
}

const mockLogs = ref([
  { id: '10001', time: '2026-07-27 10:15:32', type: 'login', typeLabel: '登录日志', operator: 'admin', ip: '192.168.1.100', description: '管理员登录系统', status: 'success', userAgent: 'Mozilla/5.0 Chrome/126.0', detail: 'Login successful for user admin from IP 192.168.1.100' },
  { id: '10002', time: '2026-07-27 09:45:18', type: 'operation', typeLabel: '操作日志', operator: 'admin', ip: '192.168.1.100', description: '修改系统基本设置', status: 'success', userAgent: 'Mozilla/5.0 Chrome/126.0', detail: 'Updated settings: siteName, siteDescription, timezone' },
  { id: '10003', time: '2026-07-27 09:30:05', type: 'error', typeLabel: '错误日志', operator: 'system', ip: '127.0.0.1', description: '邮件发送失败：SMTP连接超时', status: 'error', userAgent: 'AnchorFinance/1.0', detail: 'Error: SMTP connection timeout after 30s\nHost: smtp.example.com:465' },
  { id: '10004', time: '2026-07-27 09:12:44', type: 'operation', typeLabel: '操作日志', operator: 'admin', ip: '192.168.1.100', description: '新增OAuth提供商：GitHub', status: 'success', userAgent: 'Mozilla/5.0 Chrome/126.0', detail: 'Added OAuth provider: github\nAppID: Iv1.abc123def456' },
  { id: '10005', time: '2026-07-27 08:55:21', type: 'login', typeLabel: '登录日志', operator: 'user_zhang', ip: '10.0.0.55', description: '用户登录失败：密码错误', status: 'error', userAgent: 'Safari/605.1', detail: 'Login failed: incorrect password\nAttempt: 2/5' },
  { id: '10006', time: '2026-07-27 08:30:00', type: 'operation', typeLabel: '操作日志', operator: 'admin', ip: '192.168.1.100', description: '导出系统日志报表', status: 'success', userAgent: 'Mozilla/5.0 Chrome/126.0', detail: 'Exported logs from 2026-07-01 to 2026-07-27\nFormat: CSV' },
  { id: '10007', time: '2026-07-26 23:59:01', type: 'error', typeLabel: '错误日志', operator: 'system', ip: '127.0.0.1', description: '数据库备份失败：磁盘空间不足', status: 'error', userAgent: 'AnchorFinance/1.0 CronJob', detail: 'Backup failed: insufficient disk space\nRequired: 2.5GB\nAvailable: 1.2GB' },
  { id: '10008', time: '2026-07-26 18:20:33', type: 'operation', typeLabel: '操作日志', operator: 'editor_li', ip: '10.0.0.88', description: '发布财务报告：2026年Q2', status: 'success', userAgent: 'Edge/126.0', detail: 'Published report: Q2-2026-Financial-Report.pdf' },
])

const filteredLogs = computed(() => {
  let logs = [...mockLogs.value]
  if (filterType.value) logs = logs.filter((l) => l.type === filterType.value)
  if (searchKeyword.value) {
    const kw = searchKeyword.value.toLowerCase()
    logs = logs.filter((l) => l.description.toLowerCase().includes(kw) || l.operator.toLowerCase().includes(kw) || l.ip.includes(kw))
  }
  return logs
})

function viewDetail(row: any) { detailLog.value = row; showDetail.value = true }
function deleteLog(id: string) { mockLogs.value = mockLogs.value.filter((l) => l.id !== id); ElMessage.success('日志已删除') }
function handleExport() { ElMessage.success('日志导出中，请稍候...') }
function handleClearLogs() { mockLogs.value = []; ElMessage.success('日志已清空') }
</script>

<style scoped>
.card-header { display: flex; align-items: center; justify-content: space-between; flex-wrap: wrap; gap: 12px; }
.card-title { font-size: 16px; font-weight: 600; }
.card-actions { display: flex; align-items: center; gap: 12px; flex-wrap: wrap; }
</style>
