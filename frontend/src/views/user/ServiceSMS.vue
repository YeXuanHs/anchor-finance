<template>
  <div class="service-sms">
    <div class="page-header">
      <h2>短信服务</h2>
      <p>管理您的短信产品</p>
    </div>
    <el-card>
      <div class="info-grid">
        <div class="info-item"><label>产品名称</label><span>{{ service.product_name || '-' }}</span></div>
        <div class="info-item"><label>剩余条数</label><span class="highlight">{{ service.remaining || 0 }}</span></div>
        <div class="info-item"><label>已发送</label><span>{{ service.sent_count || 0 }}</span></div>
        <div class="info-item"><label>状态</label><el-tag :type="service.status === 'active' ? 'success' : 'danger'" size="small">{{ service.status === 'active' ? '正常' : '已暂停' }}</el-tag></div>
      </div>
      <el-divider />
      <h4>发送记录</h4>
      <el-table :data="records" v-loading="loading">
        <el-table-column prop="phone" label="手机号" width="140" />
        <el-table-column prop="content" label="内容" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }"><el-tag :type="row.status === 'success' ? 'success' : 'danger'" size="small">{{ row.status === 'success' ? '成功' : '失败' }}</el-tag></template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="180" />
      </el-table>
    </el-card>
  </div>
</template>
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/utils/request'

const loading = ref(false)
const service = ref<any>({})
const records = ref<any[]>([])

const fetchData = async () => {
  loading.value = true
  try {
    const [svcRes, recRes] = await Promise.all([
      request.get('/api/v2/user/services/sms'),
      request.get('/api/v2/user/services/sms/records')
    ])
    service.value = svcRes.data.data || {}
    records.value = recRes.data.data || []
  } catch {} finally { loading.value = false }
}

onMounted(fetchData)
</script>
<style scoped lang="scss">
.service-sms { .page-header { margin-bottom: 24px; h2 { font-size: 20px; color: #1a365d; } p { color: #6b7280; margin-top: 4px; } } }
.info-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 20px; margin-bottom: 20px; }
.info-item { label { display: block; font-size: 13px; color: #9ca3af; margin-bottom: 4px; } span { font-size: 18px; color: #1a365d; font-weight: 600; &.highlight { color: #2563eb; } } }
h4 { font-size: 16px; color: #1a365d; margin-bottom: 12px; }
</style>
