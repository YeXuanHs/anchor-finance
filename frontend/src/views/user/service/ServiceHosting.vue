<template>
  <div class="service-hosting">
    <el-card shadow="never" class="config-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span class="card-title">虚拟主机配置</span>
          <el-tag type="success" size="small" effect="light" v-if="serviceInfo.status === 'Active'">
            运行中
          </el-tag>
        </div>
      </template>

      <div class="hosting-main">
        <!-- 基本配置 -->
        <div class="config-grid">
          <div class="config-section">
            <h4 class="section-title">主机信息</h4>
            <div class="config-items">
              <div class="config-item">
                <span class="label">主机名</span>
                <span class="value mono">
                  {{ serviceInfo.host_name || '-' }}
                  <el-button
                    v-if="serviceInfo.host_name"
                    link
                    type="primary"
                    size="small"
                    @click="copyText(serviceInfo.host_name)"
                  >
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </span>
              </div>
              <div class="config-item">
                <span class="label">IP地址</span>
                <span class="value mono">
                  {{ serviceInfo.dedicated_ip || '-' }}
                  <el-button
                    v-if="serviceInfo.dedicated_ip"
                    link
                    type="primary"
                    size="small"
                    @click="copyText(serviceInfo.dedicated_ip)"
                  >
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </span>
              </div>
              <div class="config-item">
                <span class="label">用户名</span>
                <span class="value mono">
                  {{ serviceInfo.username || '-' }}
                  <el-button
                    v-if="serviceInfo.username"
                    link
                    type="primary"
                    size="small"
                    @click="copyText(serviceInfo.username)"
                  >
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </span>
              </div>
              <div class="config-item">
                <span class="label">密码</span>
                <span class="value mono">
                  <template v-if="showPassword">{{ serviceInfo.password || '-' }}</template>
                  <template v-else>********</template>
                  <el-button
                    v-if="serviceInfo.password"
                    link
                    type="primary"
                    size="small"
                    @click="showPassword = !showPassword"
                  >
                    <el-icon><component :is="showPassword ? 'Hide' : 'View'" /></el-icon>
                  </el-button>
                  <el-button
                    v-if="serviceInfo.password"
                    link
                    type="primary"
                    size="small"
                    @click="copyText(serviceInfo.password)"
                  >
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </span>
              </div>
            </div>
          </div>

          <div class="config-section">
            <h4 class="section-title">空间配置</h4>
            <div class="config-items">
              <div class="config-item">
                <span class="label">空间大小</span>
                <span class="value">{{ serviceInfo.disk_space || '-' }}</span>
              </div>
              <div class="config-item">
                <span class="label">已用空间</span>
                <span class="value">
                  <div class="usage-bar">
                    <span>{{ serviceInfo.disk_used || '0MB' }}</span>
                    <el-progress
                      :percentage="diskPercentage"
                      :stroke-width="6"
                      :show-text="false"
                      style="width: 100px;"
                    />
                  </div>
                </span>
              </div>
              <div class="config-item">
                <span class="label">月流量</span>
                <span class="value">{{ serviceInfo.bandwidth || '-' }}</span>
              </div>
              <div class="config-item">
                <span class="label">子域名数</span>
                <span class="value">{{ serviceInfo.sub_domains || '无限制' }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 支持的功能 -->
        <div class="features-section">
          <h4 class="section-title">支持功能</h4>
          <div class="features-tags">
            <el-tag
              v-for="feature in serviceInfo.features"
              :key="feature"
              size="small"
              effect="plain"
              class="feature-tag"
            >
              {{ feature }}
            </el-tag>
            <span v-if="!serviceInfo.features?.length" class="no-features">-</span>
          </div>
        </div>

        <!-- 控制面板登录 -->
        <div class="panel-section">
          <h4 class="section-title">控制面板</h4>
          <div class="panel-info">
            <div class="panel-item">
              <span class="label">面板类型</span>
              <span class="value">{{ serviceInfo.control_panel || 'cPanel' }}</span>
            </div>
            <div class="panel-item">
              <span class="label">面板地址</span>
              <span class="value mono">
                {{ serviceInfo.panel_url || '-' }}
              </span>
            </div>
            <el-button
              v-if="serviceInfo.panel_url"
              type="primary"
              @click="openPanel"
            >
              <el-icon><Link /></el-icon>
              登录控制面板
            </el-button>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="action-buttons" v-if="serviceInfo.status === 'Active'">
          <el-button type="primary" @click="$emit('action', 'fileManager')">
            <el-icon><FolderOpened /></el-icon>
            文件管理
          </el-button>
          <el-button @click="$emit('action', 'database')">
            <el-icon><Coin /></el-icon>
            数据库管理
          </el-button>
          <el-button @click="$emit('action', 'email')">
            <el-icon><Message /></el-icon>
            邮箱管理
          </el-button>
          <el-button @click="$emit('action', 'backup')">
            <el-icon><FolderChecked /></el-icon>
            备份恢复
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
  CopyDocument, Link, FolderOpened, Coin, Message,
  FolderChecked, Hide, View
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
const showPassword = ref(false)

const serviceInfo = computed(() => ({ ...props.serviceInfo, ...localInfo.value }))

onMounted(async () => {
  const id = props.serviceId || route.params.id
  if (id) {
    loading.value = true
    try {
      const { data } = await request.get(`/api/v1/host/${id}/hosting`)
      localInfo.value = data.data || data || {}
    } catch (e) {
      console.error('Failed to fetch hosting data:', e)
    } finally {
      loading.value = false
    }
  }
})

const diskPercentage = computed(() => {
  if (!props.serviceInfo.disk_used || !props.serviceInfo.disk_space) return 0
  const used = parseFloat(props.serviceInfo.disk_used)
  const total = parseFloat(props.serviceInfo.disk_space)
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

function openPanel() {
  if (props.serviceInfo.panel_url) {
    window.open(props.serviceInfo.panel_url, '_blank')
  }
}
</script>

<style scoped lang="scss">
.service-hosting {
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

.hosting-main {
  .config-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 24px;
    margin-bottom: 24px;
  }
}

.config-section, .features-section, .panel-section {
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

.usage-bar {
  display: flex;
  align-items: center;
  gap: 12px;
}

.features-section {
  margin-bottom: 24px;
}

.features-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.feature-tag {
  background: #f0f9ff;
  border-color: #b3d8ff;
  color: #409eff;
}

.no-features {
  color: #909399;
  font-size: 14px;
}

.panel-section {
  margin-bottom: 24px;
}

.panel-info {
  display: flex;
  align-items: center;
  gap: 24px;
  flex-wrap: wrap;

  .panel-item {
    display: flex;
    align-items: center;
    gap: 8px;

    .label {
      font-size: 13px;
      color: #909399;
    }

    .value {
      font-size: 14px;
      color: #303133;

      &.mono {
        font-family: 'Monaco', 'Menlo', monospace;
      }
    }
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

  .panel-info {
    flex-direction: column;
    align-items: flex-start;
  }

  .action-buttons {
    .el-button {
      flex: 1;
    }
  }
}
</style>
