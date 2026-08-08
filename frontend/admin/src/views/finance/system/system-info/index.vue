<template>
  <div class="system-info-page">
    <el-card shadow="never">
      <template #header>
        <span>关于系统</span>
      </template>

      <div class="system-info">
        <div class="info-header">
          <div class="logo">
            <el-icon :size="64" color="var(--el-color-primary)"><Monitor /></el-icon>
          </div>
          <div class="title">
            <h1>锚点财务管理系统</h1>
            <p class="version">版本: v{{ version }}</p>
            <p class="desc">兼容智简魔方(zjmf)财务系统的全能IDC管理平台</p>
          </div>
        </div>

        <el-divider />

        <el-descriptions :column="2" border title="系统信息">
          <el-descriptions-item label="系统名称">锚点财务</el-descriptions-item>
          <el-descriptions-item label="当前版本">v{{ version }}</el-descriptions-item>
          <el-descriptions-item label="Go版本">{{ systemInfo.go_version || '-' }}</el-descriptions-item>
          <el-descriptions-item label="数据库版本">{{ systemInfo.db_version || '-' }}</el-descriptions-item>
          <el-descriptions-item label="操作系统">{{ systemInfo.os || '-' }}</el-descriptions-item>
          <el-descriptions-item label="服务器时间">{{ systemInfo.server_time || '-' }}</el-descriptions-item>
          <el-descriptions-item label="启动时间">{{ systemInfo.started_at || '-' }}</el-descriptions-item>
          <el-descriptions-item label="运行时长">{{ systemInfo.uptime || '-' }}</el-descriptions-item>
        </el-descriptions>

        <el-divider />

        <el-descriptions :column="2" border title="许可信息">
          <el-descriptions-item label="授权类型">
            <el-tag type="success">开源版</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="授权状态">永久有效</el-descriptions-item>
        </el-descriptions>

        <el-divider />

        <div class="links">
          <h3>相关链接</h3>
          <el-space wrap>
            <el-button type="primary" plain @click="openLink('https://github.com/anchorfinance')">
              GitHub仓库
            </el-button>
            <el-button type="primary" plain @click="openLink('https://docs.anchorfinance.dev')">
              开发文档
            </el-button>
            <el-button type="primary" plain @click="checkUpdate">
              检查更新
            </el-button>
          </el-space>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Monitor } from '@element-plus/icons-vue'
import request from '@/utils/http'

const version = ref('1.0.0')
const systemInfo = ref<any>({})

// 获取系统信息
const fetchSystemInfo = async () => {
  try {
    const data = await request.get({ url: '/api/admin/system/info' })
    systemInfo.value = data || {}
    if (data?.version) {
      version.value = data.version
    }
  } catch (error) {
    console.error('获取系统信息失败:', error)
  }
}

// 打开链接
const openLink = (url: string) => {
  window.open(url, '_blank')
}

// 检查更新
const checkUpdate = async () => {
  try {
    const data = await request.get({ url: '/api/admin/system/check-update' })
    if (data?.has_update) {
      ElMessage.info(`发现新版本 v${data.latest_version}，请前往更新`)
    } else {
      ElMessage.success('当前已是最新版本')
    }
  } catch (error) {
    console.error('检查更新失败:', error)
  }
}

onMounted(() => {
  fetchSystemInfo()
})
</script>

<style scoped lang="scss">
.system-info-page {
  padding: 16px;
}

.system-info {
  max-width: 800px;
  margin: 0 auto;
}

.info-header {
  display: flex;
  align-items: center;
  gap: 24px;
  margin-bottom: 20px;
}

.logo {
  width: 80px;
  height: 80px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--el-color-primary-light-9);
  border-radius: 16px;
}

.title {
  h1 {
    margin: 0 0 8px 0;
    font-size: 24px;
    font-weight: 600;
  }

  .version {
    margin: 0 0 4px 0;
    font-size: 16px;
    color: var(--el-color-primary);
  }

  .desc {
    margin: 0;
    font-size: 14px;
    color: #86909C;
  }
}

.links {
  h3 {
    margin: 0 0 12px 0;
    font-size: 16px;
  }
}
</style>
