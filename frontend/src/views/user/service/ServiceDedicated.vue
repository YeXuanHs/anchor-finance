<template>
  <div class="service-dedicated">
    <el-card shadow="never" class="config-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span class="card-title">独立服务器配置</span>
          <div class="power-status" v-if="serviceInfo.power_status">
            <el-icon
              :size="16"
              :color="serviceInfo.power_status === 'on' ? '#67c23a' : '#909399'"
            >
              <component :is="serviceInfo.power_status === 'on' ? 'CircleCheck' : 'CircleClose'" />
            </el-icon>
            <span>{{ serviceInfo.power_status === 'on' ? '运行中' : '已关机' }}</span>
          </div>
        </div>
      </template>

      <div class="config-grid">
        <!-- 远程信息 -->
        <div class="config-section">
          <h4 class="section-title">远程信息</h4>
          <div class="config-items">
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
              <span class="label">SSH端口</span>
              <span class="value mono">{{ serviceInfo.port || '22' }}</span>
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

        <!-- 硬件配置 -->
        <div class="config-section">
          <h4 class="section-title">硬件配置</h4>
          <div class="config-items">
            <div class="config-item">
              <span class="label">CPU</span>
              <span class="value">{{ serviceInfo.cpu || '-' }}</span>
            </div>
            <div class="config-item">
              <span class="label">内存</span>
              <span class="value">{{ serviceInfo.memory || '-' }}</span>
            </div>
            <div class="config-item">
              <span class="label">系统盘</span>
              <span class="value">{{ serviceInfo.disk || '-' }}</span>
            </div>
            <div class="config-item" v-if="serviceInfo.data_disk">
              <span class="label">数据盘</span>
              <span class="value">{{ serviceInfo.data_disk }}</span>
            </div>
            <div class="config-item" v-if="serviceInfo.raid">
              <span class="label">RAID</span>
              <span class="value">{{ serviceInfo.raid }}</span>
            </div>
          </div>
        </div>

        <!-- 网络配置 -->
        <div class="config-section">
          <h4 class="section-title">网络配置</h4>
          <div class="config-items">
            <div class="config-item">
              <span class="label">机房位置</span>
              <span class="value">
                <img
                  v-if="serviceInfo.country_flag"
                  :src="`/upload/common/country/${serviceInfo.country_flag}.png`"
                  class="country-flag"
                  alt=""
                />
                {{ serviceInfo.data_center || '-' }}
              </span>
            </div>
            <div class="config-item">
              <span class="label">带宽</span>
              <span class="value">{{ serviceInfo.bandwidth ? `${serviceInfo.bandwidth}Mbps` : '-' }}</span>
            </div>
            <div class="config-item" v-if="serviceInfo.flow">
              <span class="label">流量</span>
              <span class="value">{{ serviceInfo.flow }}</span>
            </div>
            <div class="config-item" v-if="serviceInfo.peak_defence">
              <span class="label">防御峰值</span>
              <span class="value">{{ serviceInfo.peak_defence }}G</span>
            </div>
            <div class="config-item">
              <span class="label">操作系统</span>
              <span class="value">{{ serviceInfo.os || '-' }}</span>
            </div>
          </div>
        </div>

        <!-- 附加IP -->
        <div class="config-section" v-if="serviceInfo.ip_list?.length">
          <h4 class="section-title">IP列表</h4>
          <div class="ip-list">
            <div v-for="(ip, index) in serviceInfo.ip_list" :key="index" class="ip-item">
              <span class="ip-address mono">{{ ip.ip }}</span>
              <span class="ip-gateway">网关: {{ ip.gateway }}</span>
              <span class="ip-mask">掩码: {{ ip.subnet_mask }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="action-buttons" v-if="serviceInfo.status === 'Active'">
        <el-button type="primary" @click="$emit('action', 'powerOn')" :loading="actionLoading">
          <el-icon><VideoPlay /></el-icon>
          开机
        </el-button>
        <el-button type="danger" plain @click="$emit('action', 'powerOff')" :loading="actionLoading">
          <el-icon><VideoPause /></el-icon>
          关机
        </el-button>
        <el-button @click="$emit('action', 'reboot')" :loading="actionLoading">
          <el-icon><RefreshRight /></el-icon>
          重启
        </el-button>
        <el-button @click="$emit('action', 'vnc')">
          <el-icon><Monitor /></el-icon>
          VNC控制台
        </el-button>
        <el-button @click="$emit('action', 'reinstall')">
          <el-icon><FolderOpened /></el-icon>
          重装系统
        </el-button>
        <el-button @click="$emit('action', 'resetPassword')">
          <el-icon><Key /></el-icon>
          重置密码
        </el-button>
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
  CopyDocument, VideoPlay, VideoPause, RefreshRight,
  Monitor, FolderOpened, Key, Hide, View,
  CircleCheck, CircleClose
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
const actionLoading = ref(false)

const serviceInfo = computed(() => ({ ...props.serviceInfo, ...localInfo.value }))

onMounted(async () => {
  const id = props.serviceId || route.params.id
  if (id) {
    loading.value = true
    try {
      const { data } = await request.get(`/api/v1/hosts/${id}`)
      localInfo.value = data.data || data || {}
    } catch (e) {
      console.error('Failed to fetch dedicated server data:', e)
    } finally {
      loading.value = false
    }
  }
})

function copyText(text: string) {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success('已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}
</script>

<style scoped lang="scss">
.service-dedicated {
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

.power-status {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  color: #606266;
}

.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
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
    width: 70px;
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

.country-flag {
  width: 20px;
  height: 14px;
  object-fit: cover;
  border-radius: 2px;
}

.ip-list {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.ip-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 8px 12px;
  background: #f5f7fa;
  border-radius: 6px;
  font-size: 13px;

  .ip-address {
    font-weight: 500;
    color: #303133;
    min-width: 140px;
  }

  .ip-gateway, .ip-mask {
    color: #909399;
  }
}

.mono {
  font-family: 'Monaco', 'Menlo', monospace;
}

.action-buttons {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid #f2f3f5;
}

@media (max-width: 768px) {
  .config-grid {
    grid-template-columns: 1fr;
  }

  .ip-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }

  .action-buttons {
    .el-button {
      flex: 1;
    }
  }
}
</style>
