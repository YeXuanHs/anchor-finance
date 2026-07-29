<template>
  <div class="service-soft">
    <div class="page-header">
      <h2>软件服务</h2>
      <p>管理您的软件授权产品</p>
    </div>
    <el-card>
      <div class="info-grid">
        <div class="info-item"><label>软件名称</label><span>{{ service.product_name || '-' }}</span></div>
        <div class="info-item"><label>授权码</label><span class="mono">{{ service.license_key || '-' }}</span></div>
        <div class="info-item"><label>到期时间</label><span>{{ service.expired_at || '-' }}</span></div>
        <div class="info-item"><label>状态</label><el-tag :type="service.status === 'active' ? 'success' : 'danger'" size="small">{{ service.status === 'active' ? '正常' : '已过期' }}</el-tag></div>
      </div>
      <el-divider />
      <h4>操作</h4>
      <el-space>
        <el-button type="primary" @click="handleRenew">续费</el-button>
        <el-button @click="handleDownload">下载</el-button>
        <el-button @click="handleResetKey">重置授权码</el-button>
      </el-space>
    </el-card>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const service = ref<any>({})

const fetchData = async () => {
  try {
    const { data } = await request.get('/api/user/services/software')
    service.value = data.data || {}
  } catch {}
}

const handleRenew = () => { ElMessage.info('跳转到续费页面') }
const handleDownload = () => { ElMessage.info('开始下载') }
const handleResetKey = async () => {
  await ElMessageBox.confirm('重置后旧的授权码将立即失效，确定继续？', '确认')
  try { await request.post('/api/user/services/software/reset-key'); ElMessage.success('已重置'); fetchData() } catch {}
}

onMounted(fetchData)
</script>
<style scoped lang="scss">
.service-soft { .page-header { margin-bottom: 24px; h2 { font-size: 20px; color: #1a365d; } p { color: #6b7280; margin-top: 4px; } } }
.info-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; margin-bottom: 20px; }
.info-item { label { display: block; font-size: 13px; color: #9ca3af; margin-bottom: 4px; } span { font-size: 16px; color: #1a365d; font-weight: 600; &.mono { font-family: monospace; font-size: 14px; } } }
h4 { font-size: 16px; color: #1a365d; margin-bottom: 12px; }
</style>
