<template>
  <div class="service-domain">
    <el-card shadow="never" class="config-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span class="card-title">域名信息</span>
          <el-tag :type="getDomainStatusType(serviceInfo.domain_status)" size="small" effect="light">
            {{ getDomainStatusText(serviceInfo.domain_status) }}
          </el-tag>
        </div>
      </template>

      <div class="domain-main">
        <div class="domain-name-section">
          <h2 class="domain-name">{{ serviceInfo.domain || '-' }}</h2>
          <div class="domain-actions">
            <el-button type="primary" plain size="small" @click="$emit('action', 'dns')">
              <el-icon><Setting /></el-icon>
              DNS管理
            </el-button>
            <el-button plain size="small" @click="$emit('action', 'transfer')">
              <el-icon><Switch /></el-icon>
              域名转移
            </el-button>
          </div>
        </div>

        <el-divider />

        <div class="config-grid">
          <div class="config-items">
            <div class="config-item">
              <span class="label">域名</span>
              <span class="value mono">{{ serviceInfo.domain || '-' }}</span>
            </div>
            <div class="config-item">
              <span class="label">注册时间</span>
              <span class="value">{{ serviceInfo.register_time || '-' }}</span>
            </div>
            <div class="config-item">
              <span class="label">到期时间</span>
              <span class="value" :class="{ 'text-danger': isExpiringSoon }">
                {{ serviceInfo.due_time || '-' }}
              </span>
            </div>
            <div class="config-item">
              <span class="label">注册商</span>
              <span class="value">{{ serviceInfo.registrar || '-' }}</span>
            </div>
          </div>

          <div class="config-items">
            <div class="config-item">
              <span class="label">WHOIS保护</span>
              <span class="value">
                <el-tag :type="serviceInfo.whois_protect ? 'success' : 'info'" size="small">
                  {{ serviceInfo.whois_protect ? '已开启' : '未开启' }}
                </el-tag>
              </span>
            </div>
            <div class="config-item">
              <span class="label">自动续费</span>
              <span class="value">
                <el-switch
                  v-model="autoRenew"
                  @change="handleAutoRenewChange"
                />
              </span>
            </div>
            <div class="config-item">
              <span class="label">锁定状态</span>
              <span class="value">
                <el-tag :type="serviceInfo.domain_lock ? 'success' : 'warning'" size="small">
                  {{ serviceInfo.domain_lock ? '已锁定' : '未锁定' }}
                </el-tag>
              </span>
            </div>
            <div class="config-item">
              <span class="label">DNS服务器</span>
              <span class="value">
                <div v-if="serviceInfo.nameservers?.length" class="nameservers">
                  <span v-for="(ns, index) in serviceInfo.nameservers" :key="index" class="ns-item mono">
                    {{ ns }}
                  </span>
                </div>
                <span v-else>-</span>
              </span>
            </div>
          </div>
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
import { Setting, Switch } from '@element-plus/icons-vue'

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

const autoRenew = ref(serviceInfo.value.auto_renew || false)

onMounted(async () => {
  const id = props.serviceId || route.params.id
  if (id) {
    loading.value = true
    try {
      const { data } = await request.get(`/api/v1/host/${id}/domain`)
      localInfo.value = data.data || data || {}
    } catch (e) {
      console.error('Failed to fetch domain data:', e)
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

function getDomainStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    Active: 'success',
    Expired: 'danger',
    Pending: 'info',
    Transfer: 'warning'
  }
  return map[status] || 'info'
}

function getDomainStatusText(status: string) {
  const map: Record<string, string> = {
    Active: '正常',
    Expired: '已过期',
    Pending: '待处理',
    Transfer: '转移中'
  }
  return map[status] || status || '正常'
}

function handleAutoRenewChange(val: boolean) {
  ElMessage.success(`自动续费已${val ? '开启' : '关闭'}`)
}
</script>

<style scoped lang="scss">
.service-domain {
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

.domain-main {
  .domain-name-section {
    display: flex;
    align-items: center;
    justify-content: space-between;
    flex-wrap: wrap;
    gap: 16px;
  }

  .domain-name {
    font-size: 24px;
    font-weight: 700;
    color: #303133;
    margin: 0;
    font-family: 'Monaco', 'Menlo', monospace;
  }

  .domain-actions {
    display: flex;
    gap: 8px;
  }
}

.config-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
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

.nameservers {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.ns-item {
  font-size: 13px;
  padding: 2px 8px;
  background: #f5f7fa;
  border-radius: 4px;
  display: inline-block;
}

@media (max-width: 768px) {
  .domain-name-section {
    flex-direction: column;
    align-items: flex-start;
  }

  .config-grid {
    grid-template-columns: 1fr;
  }
}
</style>
