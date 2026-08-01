<template>
  <div class="system-info-page">
    <el-card shadow="never" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span>系统信息</span>
          <el-button type="primary" @click="fetchSystemInfo" :icon="Refresh">刷新</el-button>
        </div>
      </template>

      <!-- 服务器信息 -->
      <el-descriptions title="服务器信息" :column="2" border class="info-section">
        <el-descriptions-item label="服务器IP">{{ info.server_ip }}</el-descriptions-item>
        <el-descriptions-item label="服务器名称">{{ info.server_name }}</el-descriptions-item>
        <el-descriptions-item label="服务器端口">{{ info.server_port }}</el-descriptions-item>
        <el-descriptions-item label="操作系统">{{ info.server_system }}</el-descriptions-item>
        <el-descriptions-item label="PHP版本">
          <el-tag>{{ info.php_version }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="MySQL版本">
          <el-tag type="success">{{ info.mysql_version }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="上传限制">{{ info.upload_max_filesize }}</el-descriptions-item>
        <el-descriptions-item label="执行时间限制">{{ info.max_execution_time }}</el-descriptions-item>
        <el-descriptions-item label="内存限制">{{ info.memory_limit }}</el-descriptions-item>
        <el-descriptions-item label="当前时间">{{ info.now_time }}</el-descriptions-item>
      </el-descriptions>

      <!-- 系统版本 -->
      <el-descriptions title="系统版本" :column="2" border class="info-section">
        <el-descriptions-item label="安装版本">
          <el-tag type="primary">{{ info.install_version }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="版本类型">
          <el-tag :type="info.system_version_type === 'stable' ? 'success' : 'warning'">
            {{ info.system_version_type === 'stable' ? '稳定版' : '测试版' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="授权状态">
          <el-tag :type="info.auth_status === 'Active' ? 'success' : 'danger'">
            {{ info.auth_status === 'Active' ? '已授权' : info.auth_status || '未授权' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="授权到期时间">{{ info.auth_due_time }}</el-descriptions-item>
        <el-descriptions-item label="服务到期时间">{{ info.service_due_time }}</el-descriptions-item>
        <el-descriptions-item label="系统Token">
          <el-text truncated>{{ info.system_token || '-' }}</el-text>
        </el-descriptions-item>
      </el-descriptions>

      <!-- 操作区 -->
      <div class="action-section">
        <el-button type="primary" @click="handleCheckUpdate">检查更新</el-button>
        <el-button @click="handleDatabaseInfo">数据库信息</el-button>
        <el-button @click="handleOptimize">优化数据库</el-button>
        <el-button @click="handleBackup">备份数据库</el-button>
      </div>
    </el-card>

    <!-- 数据库信息对话框 -->
    <el-dialog v-model="dbDialogVisible" title="数据库信息" width="700px">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="数据表总数">{{ dbInfo.total_count }}</el-descriptions-item>
        <el-descriptions-item label="总行数">{{ dbInfo.total_rows }}</el-descriptions-item>
        <el-descriptions-item label="总大小">{{ dbInfo.total_size }}</el-descriptions-item>
      </el-descriptions>
      <el-table :data="dbInfo.report_array" stripe border style="margin-top: 16px" max-height="400">
        <el-table-column prop="name" label="表名" />
        <el-table-column prop="rows" label="行数" width="120" />
        <el-table-column prop="size" label="大小" width="120" />
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const dbDialogVisible = ref(false)

const info = reactive({
  server_ip: '',
  server_name: '',
  server_port: '',
  server_system: '',
  php_version: '',
  mysql_version: '',
  upload_max_filesize: '',
  max_execution_time: '',
  memory_limit: '',
  now_time: '',
  install_version: '',
  system_version_type: '',
  auth_status: '',
  auth_due_time: '',
  service_due_time: '',
  system_token: ''
})

const dbInfo = reactive({
  total_count: 0,
  total_rows: 0,
  total_size: '',
  report_array: [] as Array<{ name: string; rows: number; size: string }>
})

const fetchSystemInfo = async () => {
  loading.value = true
  try {
    const res = await request.get({ url: '/api/admin/system/info' })
    if (res) Object.assign(info, res)
  } catch (error) {
    console.error('获取系统信息失败:', error)
  } finally {
    loading.value = false
  }
}

const handleCheckUpdate = async () => {
  try {
    const res = await request.get({ url: '/api/admin/system/last-version' })
    if (res?.last_version) {
      ElMessage.info(`最新版本: ${res.last_version}`)
    }
  } catch (error) {
    ElMessage.error('检查更新失败')
  }
}

const handleDatabaseInfo = async () => {
  try {
    const res = await request.get({ url: '/api/admin/system/database-info' })
    if (res) {
      Object.assign(dbInfo, res)
      dbDialogVisible.value = true
    }
  } catch (error) {
    ElMessage.error('获取数据库信息失败')
  }
}

const handleOptimize = async () => {
  try {
    await ElMessageBox.confirm('确定要优化数据库表吗？', '提示')
    await request.post({ url: '/api/admin/system/optimize-tables' })
    ElMessage.success('数据库优化完成')
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('优化失败')
  }
}

const handleBackup = async () => {
  try {
    await ElMessageBox.confirm('确定要备份数据库吗？', '提示')
    await request.post({ url: '/api/admin/system/backup-database' })
    ElMessage.success('数据库备份完成')
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('备份失败')
  }
}

onMounted(() => fetchSystemInfo())
</script>

<style scoped lang="scss">
.system-info-page {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.info-section {
  margin-bottom: 24px;
}
.action-section {
  margin-top: 24px;
  display: flex;
  gap: 12px;
}
</style>
