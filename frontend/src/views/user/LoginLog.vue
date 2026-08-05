<template>
  <div class="login-log-page">
    <div class="page-header">
      <h1 class="page-title">登录日志</h1>
    </div>

    <el-card shadow="never" class="filter-card">
      <el-form :model="filters" inline>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="filters.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            style="width: 300px"
          />
        </el-form-item>
        <el-form-item label="登录状态">
          <el-select v-model="filters.status" placeholder="全部" clearable style="width: 130px">
            <el-option label="成功" value="success" />
            <el-option label="失败" value="failed" />
          </el-select>
        </el-form-item>
        <el-form-item label="IP 地址">
          <el-input v-model="filters.ip" placeholder="输入 IP 地址" clearable style="width: 180px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never" class="table-card">
      <el-table :data="logs" style="width: 100%" v-loading="loading" stripe>
        <el-table-column prop="time" label="登录时间" width="170" sortable />
        <el-table-column prop="ip" label="IP 地址" width="150" />
        <el-table-column prop="location" label="登录地点" width="140" />
        <el-table-column prop="device" label="设备信息" min-width="200" show-overflow-tooltip>
          <template #default="{ row }">
            <div class="device-info">
              <el-icon :size="16" color="#909399">
                <component :is="deviceIcon(row.deviceType)" />
              </el-icon>
              <span>{{ row.device }}</span>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="browser" label="浏览器" width="130" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small" effect="light" round>
              {{ row.status === 'success' ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="failReason" label="失败原因" width="150">
          <template #default="{ row }">
            <span v-if="row.status === 'failed'" class="fail-reason">{{ row.failReason }}</span>
            <span v-else class="text-muted">-</span>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :page-sizes="[20, 50, 100]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        style="margin-top: 20px; justify-content: flex-end"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import request from '@/utils/request'
import { Search, Monitor, Iphone, Platform } from '@element-plus/icons-vue'

interface LoginLog {
  id: string
  time: string
  ip: string
  location: string
  device: string
  deviceType: 'desktop' | 'mobile' | 'tablet'
  browser: string
  status: 'success' | 'failed'
  failReason: string
}

const loading = ref(false)
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)

const filters = reactive({
  dateRange: null as any,
  status: '',
  ip: ''
})

const logs = ref<LoginLog[]>([])

function deviceIcon(type: string) {
  const map: Record<string, any> = {
    desktop: Monitor,
    mobile: Iphone,
    tablet: Platform
  }
  return map[type] || Monitor
}

async function loadData() {
  loading.value = true
  try {
    const res = await request.get('/api/v1/login-logs', {
      params: { page: currentPage.value, page_size: pageSize.value, status: filters.status, ip: filters.ip }
    })
    logs.value = res.data?.data || []
    total.value = res.data?.total || 0
  } catch { /* ignore */ }
  loading.value = false
}

function handleSearch() {
  currentPage.value = 1
  loadData()
}

function handleReset() {
  filters.dateRange = null
  filters.status = ''
  filters.ip = ''
}

onMounted(() => { loadData() })
</script>

<style scoped lang="scss">
.login-log-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.filter-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;

  :deep(.el-form-item) {
    margin-bottom: 0;
  }
}

.table-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.device-info {
  display: flex;
  align-items: center;
  gap: 8px;
}

.fail-reason {
  color: #f56c6c;
  font-size: 13px;
}

.text-muted {
  color: #c0c4cc;
}
</style>
