<template>
  <div class="service-nat">
    <el-card shadow="never" class="config-card" v-loading="loading">
      <template #header>
        <div class="card-header">
          <span class="card-title">NAT网关配置</span>
          <el-tag type="success" size="small" effect="light" v-if="serviceInfo.status === 'Active'">
            运行中
          </el-tag>
        </div>
      </template>

      <div class="nat-main">
        <!-- 基本配置 -->
        <div class="config-grid">
          <div class="config-section">
            <h4 class="section-title">网关信息</h4>
            <div class="config-items">
              <div class="config-item">
                <span class="label">网关名称</span>
                <span class="value">{{ serviceInfo.gateway_name || '-' }}</span>
              </div>
              <div class="config-item">
                <span class="label">公网IP</span>
                <span class="value mono">
                  {{ serviceInfo.public_ip || '-' }}
                  <el-button
                    v-if="serviceInfo.public_ip"
                    link
                    type="primary"
                    size="small"
                    @click="copyText(serviceInfo.public_ip)"
                  >
                    <el-icon><CopyDocument /></el-icon>
                  </el-button>
                </span>
              </div>
              <div class="config-item">
                <span class="label">内网IP</span>
                <span class="value mono">{{ serviceInfo.private_ip || '-' }}</span>
              </div>
              <div class="config-item">
                <span class="label">规格</span>
                <span class="value">{{ serviceInfo.spec_name || '-' }}</span>
              </div>
            </div>
          </div>

          <div class="config-section">
            <h4 class="section-title">性能指标</h4>
            <div class="config-items">
              <div class="config-item">
                <span class="label">最大连接数</span>
                <span class="value">{{ serviceInfo.max_connections || '-' }}</span>
              </div>
              <div class="config-item">
                <span class="label">新建连接数</span>
                <span class="value">{{ serviceInfo.new_connections || '-' }}/秒</span>
              </div>
              <div class="config-item">
                <span class="label">吞吐量</span>
                <span class="value">{{ serviceInfo.throughput || '-' }}</span>
              </div>
              <div class="config-item">
                <span class="label">所属VPC</span>
                <span class="value">{{ serviceInfo.vpc_name || '-' }}</span>
              </div>
            </div>
          </div>
        </div>

        <!-- 端口转发规则 -->
        <div class="rules-section">
          <div class="section-header">
            <h4 class="section-title">端口转发规则</h4>
            <el-button type="primary" size="small" @click="$emit('action', 'addRule')">
              <el-icon><Plus /></el-icon>
              添加规则
            </el-button>
          </div>

          <el-table
            :data="serviceInfo.port_rules || []"
            style="width: 100%"
            empty-text="暂无转发规则"
          >
            <el-table-column prop="protocol" label="协议" width="100">
              <template #default="{ row }">
                <el-tag size="small" :type="row.protocol === 'tcp' ? '' : 'warning'">
                  {{ row.protocol?.toUpperCase() }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="public_port" label="公网端口" width="120" />
            <el-table-column prop="private_ip" label="内网IP" min-width="150">
              <template #default="{ row }">
                <span class="mono">{{ row.private_ip }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="private_port" label="内网端口" width="120" />
            <el-table-column prop="remark" label="备注" min-width="150" />
            <el-table-column label="操作" width="120" align="center">
              <template #default="{ row, $index }">
                <el-button link type="primary" size="small" @click="$emit('action', 'editRule')">
                  编辑
                </el-button>
                <el-button link type="danger" size="small" @click="$emit('action', 'deleteRule')">
                  删除
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>

        <!-- 流量统计 -->
        <div class="flow-section">
          <h4 class="section-title">流量统计</h4>
          <div class="flow-grid">
            <div class="flow-item">
              <div class="flow-label">入站流量</div>
              <div class="flow-value">{{ serviceInfo.flow_in || '0GB' }}</div>
            </div>
            <div class="flow-item">
              <div class="flow-label">出站流量</div>
              <div class="flow-value">{{ serviceInfo.flow_out || '0GB' }}</div>
            </div>
            <div class="flow-item">
              <div class="flow-label">总流量</div>
              <div class="flow-value">{{ serviceInfo.flow_total || '0GB' }}</div>
            </div>
            <div class="flow-item">
              <div class="flow-label">流量配额</div>
              <div class="flow-value">{{ serviceInfo.flow_quota || '无限制' }}</div>
            </div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div class="action-buttons" v-if="serviceInfo.status === 'Active'">
          <el-button type="primary" @click="$emit('action', 'monitor')">
            <el-icon><DataLine /></el-icon>
            监控图表
          </el-button>
          <el-button @click="$emit('action', 'securityGroup')">
            <el-icon><Lock /></el-icon>
            安全组
          </el-button>
          <el-button @click="$emit('action', 'upgrade')">
            <el-icon><Top /></el-icon>
            升级规格
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
  CopyDocument, Plus, DataLine, Lock, Top
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

onMounted(async () => {
  const id = props.serviceId || route.params.id
  if (id) {
    loading.value = true
    try {
      const { data } = await request.get(`/api/v1/hosts/${id}`)
      localInfo.value = data.data || data || {}
    } catch (e) {
      console.error('Failed to fetch NAT data:', e)
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
.service-nat {
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

.nat-main {
  .config-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 24px;
    margin-bottom: 24px;
  }
}

.config-section, .rules-section, .flow-section {
  .section-title {
    font-size: 14px;
    font-weight: 600;
    color: #303133;
    margin: 0 0 12px 0;
    padding-bottom: 8px;
    border-bottom: 1px solid #f2f3f5;
  }
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 12px;

  .section-title {
    margin: 0;
    padding: 0;
    border: none;
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
    width: 90px;
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

.rules-section {
  margin-bottom: 24px;
}

.mono {
  font-family: 'Monaco', 'Menlo', monospace;
}

.flow-section {
  margin-bottom: 24px;
}

.flow-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
  gap: 16px;
}

.flow-item {
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
  text-align: center;

  .flow-label {
    font-size: 13px;
    color: #909399;
    margin-bottom: 8px;
  }

  .flow-value {
    font-size: 18px;
    font-weight: 600;
    color: #303133;
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

  .flow-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .action-buttons {
    .el-button {
      flex: 1;
    }
  }
}
</style>
