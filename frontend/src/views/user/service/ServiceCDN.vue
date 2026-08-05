<template>
  <div class="service-cdn">
    <el-card shadow="never" class="config-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span class="card-title">CDN加速配置</span>
          <el-tag type="success" size="small" effect="light" v-if="serviceInfo.status === 'Active'">
            运行中
          </el-tag>
        </div>
      </template>

      <div class="cdn-main">
        <!-- 基本配置 -->
        <div class="config-grid">
          <div class="config-section">
            <h4 class="section-title">基本信息</h4>
            <div class="config-items">
              <div class="config-item">
                <span class="label">加速域名</span>
                <span class="value mono">
                  {{ serviceInfo.domain || '-' }}
                  <el-button
                    v-if="serviceInfo.domain"
                    link
                    type="primary"
                    size="small"
                    @click="copyText(serviceInfo.domain)"
                  >
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </span>
              </div>
              <div class="config-item">
                <span class="label">CNAME</span>
                <span class="value mono">
                  {{ serviceInfo.cname || '-' }}
                  <el-button
                    v-if="serviceInfo.cname"
                    link
                    type="primary"
                    size="small"
                    @click="copyText(serviceInfo.cname)"
                  >
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </span>
              </div>
              <div class="config-item">
                <span class="label">源站地址</span>
                <span class="value mono">{{ serviceInfo.origin || '-' }}</span>
              </div>
              <div class="config-item">
                <span class="label">源站端口</span>
                <span class="value mono">{{ serviceInfo.origin_port || '80' }}</span>
              </div>
            </div>
          </div>

          <div class="config-section">
            <h4 class="section-title">套餐信息</h4>
            <div class="config-items">
              <div class="config-item">
                <span class="label">套餐类型</span>
                <span class="value">{{ serviceInfo.plan_name || '-' }}</span>
              </div>
              <div class="config-item">
                <span class="label">流量配额</span>
                <span class="value">{{ serviceInfo.flow_quota || '-' }}</span>
              </div>
              <div class="config-item">
                <span class="label">已用流量</span>
                <span class="value">
                  <div class="flow-usage">
                    <span>{{ serviceInfo.flow_used || '0GB' }}</span>
                    <el-progress
                      :percentage="flowPercentage"
                      :stroke-width="6"
                      :show-text="false"
                      style="width: 100px;"
                    />
                  </div>
                </span>
              </div>
              <div class="config-item">
                <span class="label">带宽峰值</span>
                <span class="value">{{ serviceInfo.bandwidth_peak || '-' }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- CDN功能配置 -->
        <div class="features-section">
          <h4 class="section-title">功能配置</h4>
          <div class="features-grid">
            <div class="feature-item">
              <div class="feature-info">
                <span class="feature-name">HTTPS加速</span>
                <span class="feature-desc">启用HTTPS安全加速</span>
              </div>
              <el-switch v-model="features.https" @change="handleFeatureChange('https')" />
            </div>
            <div class="feature-item">
              <div class="feature-info">
                <span class="feature-name">页面缓存</span>
                <span class="feature-desc">静态资源缓存加速</span>
              </div>
              <el-switch v-model="features.cache" @change="handleFeatureChange('cache')" />
            </div>
            <div class="feature-item">
              <div class="feature-info">
                <span class="feature-name">Gzip压缩</span>
                <span class="feature-desc">自动压缩传输内容</span>
              </div>
              <el-switch v-model="features.gzip" @change="handleFeatureChange('gzip')" />
            </div>
            <div class="feature-item">
              <div class="feature-info">
                <span class="feature-name">WebSocket</span>
                <span class="feature-desc">支持WebSocket协议</span>
              </div>
              <el-switch v-model="features.websocket" @change="handleFeatureChange('websocket')" />
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="action-buttons" v-if="serviceInfo.status === 'Active'">
          <el-button type="primary" @click="$emit('action', 'purge')">
            <el-icon><Delete /></el-icon>
            刷新缓存
          </el-button>
          <el-button @click="$emit('action', 'ssl')">
            <el-icon><Lock /></el-icon>
            SSL证书管理
          </el-button>
          <el-button @click="$emit('action', 'statistics')">
            <el-icon><DataLine /></el-icon>
            访问统计
          </el-button>
          <el-button @click="$emit('action', 'log')">
            <el-icon><Document /></el-icon>
            访问日志
          </el-button>
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
import {
  CopyDocument, Delete, Lock, DataLine, Document
} from '@element-plus/icons-vue'

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

const features = ref({
  https: serviceInfo.value.features?.https || false,
  cache: serviceInfo.value.features?.cache ?? true,
  gzip: serviceInfo.value.features?.gzip ?? true,
  websocket: serviceInfo.value.features?.websocket || false
})

onMounted(async () => {
  const id = props.serviceId || route.params.id
  if (id) {
    loading.value = true
    try {
      const { data } = await request.get(`/api/v1/hosts/${id}`)
      localInfo.value = data.data || data || {}
    } catch (e) {
      console.error('Failed to fetch CDN data:', e)
    } finally {
      loading.value = false
    }
  }
})

const flowPercentage = computed(() => {
  if (!props.serviceInfo.flow_used || !props.serviceInfo.flow_quota) return 0
  const used = parseFloat(props.serviceInfo.flow_used)
  const total = parseFloat(props.serviceInfo.flow_quota)
  if (isNaN(used) || isNaN(total) || total === 0) return 0
  return Math.min(Math.round((used / total) * 100), 100)
})

function copyText(text: string) {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success('已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}

function handleFeatureChange(feature: string) {
  ElMessage.success(`${feature} 设置已更新`)
}
</script>

<style scoped lang="scss">
.service-cdn {
  .config-card {
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

.cdn-main {
  .config-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 24px;
    margin-bottom: 24px;
  }
}

.config-section {
  .section-title {
    font-size: 14px;
    font-weight: 600;
    color: #303133;
    margin: 0 0 12px 0;
    padding-bottom: 8px;
    border-bottom: 1px solid #f2f3f5;
  }
}

.config-items {
  display: flex;
  flex-direction: column;
  gap: 12px;
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
    display: flex;
    align-items: center;
    gap: 4px;
    line-height: 24px;

    &.mono {
      font-family: 'Monaco', 'Menlo', monospace;
    }
  }
}

.flow-usage {
  display: flex;
  align-items: center;
  gap: 12px;
}

.features-section {
  margin-bottom: 24px;

  .section-title {
    font-size: 14px;
    font-weight: 600;
    color: #303133;
    margin: 0 0 16px 0;
    padding-bottom: 8px;
    border-bottom: 1px solid #f2f3f5;
  }
}

.features-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 16px;
}

.feature-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 12px 16px;
  background: #f5f7fa;
  border-radius: 8px;

  .feature-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .feature-name {
    font-size: 14px;
    font-weight: 500;
    color: #303133;
  }

  .feature-desc {
    font-size: 12px;
    color: #909399;
  }
}

.action-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  padding-top: 20px;
  border-top: 1px solid #f2f3f5;
}

@media (max-width: 768px) {
  .config-grid {
    grid-template-columns: 1fr;
  }

  .features-grid {
    grid-template-columns: 1fr;
  }

  .action-buttons {
    .el-button {
      flex: 1;
    }
  }
}
</style>
