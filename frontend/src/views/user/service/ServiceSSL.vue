<template>
  <div class="service-ssl">
    <el-card shadow="never" class="config-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span class="card-title">SSL证书信息</span>
          <el-tag :type="getSSLStatusType(serviceInfo.ssl_status)" size="small" effect="light">
            {{ getSSLStatusText(serviceInfo.ssl_status) }}
          </el-tag>
        </div>
      </template>

      <div class="ssl-main">
        <div class="ssl-brand-section">
          <div class="ssl-icon">
            <el-icon :size="32" color="#67c23a"><Lock /></el-icon>
          </div>
          <div class="ssl-basic">
            <h3 class="ssl-product-name">{{ serviceInfo.product_name || 'SSL证书' }}</h3>
            <p class="ssl-domain">{{ serviceInfo.domain || '-' }}</p>
          </div>
        </div>

        <el-divider />

        <div class="config-grid">
          <div class="config-items">
            <div class="config-item">
              <span class="label">证书类型</span>
              <span class="value">{{ serviceInfo.cert_type || '-' }}</span>
            </div>
            <div class="config-item">
              <span class="label">验证级别</span>
              <span class="value">
                <el-tag size="small">{{ serviceInfo.verify_level || '-' }}</el-tag>
              </span>
            </div>
            <div class="config-item">
              <span class="label">域名数量</span>
              <span class="value">{{ serviceInfo.domain_count || 1 }}</span>
            </div>
            <div class="config-item">
              <span class="label">签发时间</span>
              <span class="value">{{ serviceInfo.issue_time || '-' }}</span>
            </div>
          </div>

          <div class="config-items">
            <div class="config-item">
              <span class="label">到期时间</span>
              <span class="value" :class="{ 'text-danger': isExpiringSoon }">
                {{ serviceInfo.due_time || '-' }}
              </span>
            </div>
            <div class="config-item">
              <span class="label">签发机构</span>
              <span class="value">{{ serviceInfo.issuer || '-' }}</span>
            </div>
            <div class="config-item">
              <span class="label">加密算法</span>
              <span class="value mono">{{ serviceInfo.algorithm || 'RSA 2048' }}</span>
            </div>
            <div class="config-item">
              <span class="label">证书状态</span>
              <span class="value">
                <el-tag :type="getSSLStatusType(serviceInfo.ssl_status)" size="small">
                  {{ getSSLStatusText(serviceInfo.ssl_status) }}
                </el-tag>
              </span>
            </div>
          </div>
        </div>

        <!-- 证书操作 -->
        <div class="cert-actions" v-if="serviceInfo.status === 'Active'">
          <el-button type="primary" @click="$emit('action', 'download')" v-if="serviceInfo.ssl_status === 'Issued'">
            <el-icon><Download /></el-icon>
            下载证书
          </el-button>
          <el-button @click="$emit('action', 'reissue')">
            <el-icon><RefreshRight /></el-icon>
            重新签发
          </el-button>
          <el-button @click="$emit('action', 'renew')">
            <el-icon><Refresh /></el-icon>
            续费证书
          </el-button>
        </div>
      </div>
    </el-card>

    <!-- 证书详情 -->
    <el-card v-if="serviceInfo.cert_detail" shadow="never" class="detail-card">
      <template #header>
        <span class="card-title">证书详情</span>
      </template>
      <div class="cert-detail">
        <div class="detail-item">
          <span class="label">证书内容 (Certificate)</span>
          <el-input
            v-model="serviceInfo.cert_detail.certificate"
            type="textarea"
            :rows="4"
            readonly
            class="cert-textarea"
          >
            <template #append>
              <el-button @click="copyText(serviceInfo.cert_detail.certificate)">复制</el-button>
            </template>
          </el-input>
        </div>
        <div class="detail-item">
          <span class="label">私钥 (Private Key)</span>
          <el-input
            v-model="serviceInfo.cert_detail.private_key"
            type="textarea"
            :rows="4"
            readonly
            class="cert-textarea"
          >
            <template #append>
              <el-button @click="copyText(serviceInfo.cert_detail.private_key)">复制</el-button>
            </template>
          </el-input>
        </div>
        <div class="detail-item">
          <span class="label">证书链 (CA Bundle)</span>
          <el-input
            v-model="serviceInfo.cert_detail.ca_bundle"
            type="textarea"
            :rows="4"
            readonly
            class="cert-textarea"
          >
            <template #append>
              <el-button @click="copyText(serviceInfo.cert_detail.ca_bundle)">复制</el-button>
            </template>
          </el-input>
        </div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'
import { Lock, Download, RefreshRight, Refresh } from '@element-plus/icons-vue'

const route = useRoute()

const props = withDefaults(defineProps<{
  serviceId?: number | string
  serviceInfo?: any
}>(), {
  serviceInfo: () => ({})
})

const emit = defineEmits<{
  (e: 'action', action: string): void
}>()

const loading = ref(false)
const localInfo = ref<any>({})

const serviceInfo = computed(() => ({ ...props.serviceInfo, ...localInfo.value }))

onMounted(async () => {
  const id = props.serviceId || route.params.id
  if (id) {
    loading.value = true
    try {
      const { data } = await request.get(`/api/v2/hosts/${id}`)
      localInfo.value = data.data || data || {}
    } catch (e) {
      console.error('Failed to fetch SSL certificate data:', e)
    } finally {
      loading.value = false
    }
  }
})

const isExpiringSoon = computed(() => {
  if (!props.serviceInfo.due_time) return false
  const dueDate = new Date(props.serviceInfo.due_time)
  const now = new Date()
  const diffDays = (dueDate.getTime() - now.getTime()) / (1000 * 60 * 60 * 24)
  return diffDays <= 30
})

function getSSLStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    Issued: 'success',
    Pending: 'warning',
    Expired: 'danger',
    Revoked: 'info'
  }
  return map[status] || 'info'
}

function getSSLStatusText(status: string) {
  const map: Record<string, string> = {
    Issued: '已签发',
    Pending: '待验证',
    Expired: '已过期',
    Revoked: '已吊销'
  }
  return map[status] || status || '待签发'
}

function copyText(text: string) {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success('已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}
</script>

<style scoped lang="scss">
.service-ssl {
  display: flex;
  flex-direction: column;
  gap: 20px;

  .config-card, .detail-card {
    border-radius: 12px;
    border: 1px solid #e8ecf1;

    :deep(.el-card__header) {
      padding: 16px 20px;
      border-bottom: 1px solid #f2f3f5;
    }

    :deep(.el-card__body) {
      padding: 20px;
    }
  }
}

.card-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.card-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
}

.ssl-main {
  .ssl-brand-section {
    display: flex;
    align-items: center;
    gap: 16px;
  }

  .ssl-icon {
    width: 56px;
    height: 56px;
    background: linear-gradient(135deg, #67c23a 0%, #4caf50 100%);
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .ssl-basic {
    .ssl-product-name {
      font-size: 18px;
      font-weight: 600;
      color: #303133;
      margin: 0 0 4px 0;
    }

    .ssl-domain {
      font-size: 14px;
      color: #909399;
      margin: 0;
      font-family: 'Monaco', 'Menlo', monospace;
    }
  }
}

.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(280px, 1fr));
  gap: 24px;
}

.config-items {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.config-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;

  .label {
    width: 80px;
    flex-shrink: 0;
    font-size: 13px;
    color: #909399;
    line-height: 24px;
  }

  .value {
    flex: 1;
    font-size: 14px;
    color: #303133;
    line-height: 24px;

    &.mono {
      font-family: 'Monaco', 'Menlo', monospace;
    }

    &.text-danger {
      color: #f56c6c;
    }
  }
}

.cert-actions {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #f2f3f5;
}

.cert-detail {
  display: flex;
  flex-direction: column;
  gap: 16px;

  .detail-item {
    .label {
      display: block;
      font-size: 13px;
      color: #909399;
      margin-bottom: 8px;
    }

    .cert-textarea {
      :deep(.el-textarea__inner) {
        font-family: 'Monaco', 'Menlo', monospace;
        font-size: 12px;
      }
    }
  }
}

@media (max-width: 768px) {
  .config-grid {
    grid-template-columns: 1fr;
  }

  .cert-actions {
    .el-button {
      flex: 1;
    }
  }
}
</style>
