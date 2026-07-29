<template>
  <div class="contract-host-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>合同关联主机</span>
        </div>
      </template>

      <div class="contract-list">
        <div
          v-for="contract in contracts"
          :key="contract.id"
          class="contract-item"
        >
          <div class="contract-header">
            <div class="contract-info">
              <h3 class="contract-title">{{ contract.contract_no }}</h3>
              <el-tag :type="getStatusType(contract.status)" size="small">
                {{ getStatusText(contract.status) }}
              </el-tag>
            </div>
            <div class="contract-meta">
              <span>产品：{{ contract.product_name }}</span>
              <span>金额：¥{{ contract.amount?.toFixed(2) }}</span>
              <span>{{ contract.start_date }} ~ {{ contract.end_date }}</span>
            </div>
          </div>

          <el-divider />

          <div class="host-section">
            <div class="section-header">
              <span>关联主机</span>
              <el-button
                v-if="contract.status === 'active'"
                type="primary"
                size="small"
                plain
                @click="showBindHostDialog(contract)"
              >
                <el-icon><Plus /></el-icon> 关联主机
              </el-button>
            </div>

            <div v-if="contract.hosts && contract.hosts.length > 0" class="host-list">
              <div
                v-for="host in contract.hosts"
                :key="host.id"
                class="host-item"
              >
                <div class="host-info">
                  <el-icon :size="20" color="#409eff"><Monitor /></el-icon>
                  <div class="host-text">
                    <div class="host-name">{{ host.name }}</div>
                    <div class="host-domain">{{ host.domain }}</div>
                  </div>
                </div>
                <div class="host-specs">
                  <span>{{ host.cpu }}核CPU</span>
                  <span>{{ host.memory }}GB内存</span>
                  <span>{{ host.disk }}GB{{ host.disk_type }}</span>
                </div>
                <el-tag :type="getHostStatusType(host.status)" size="small">
                  {{ getHostStatusText(host.status) }}
                </el-tag>
                <el-button
                  v-if="contract.status === 'active'"
                  type="danger"
                  size="small"
                  plain
                  @click="unbindHost(contract, host)"
                >
                  解绑
                </el-button>
              </div>
            </div>
            <el-empty v-else description="暂未关联主机" :image-size="60" />
          </div>
        </div>
      </div>

      <el-empty v-if="contracts.length === 0" description="暂无合同" />
    </el-card>

    <!-- 关联主机对话框 -->
    <el-dialog v-model="bindDialogVisible" title="关联主机" width="600px">
      <div class="bind-host-content">
        <div class="available-hosts">
          <h4>可关联主机</h4>
          <div class="host-select-list">
            <div
              v-for="host in availableHosts"
              :key="host.id"
              class="host-select-item"
              :class="{ selected: selectedHosts.includes(host.id) }"
              @click="toggleHostSelection(host.id)"
            >
              <el-checkbox :model-value="selectedHosts.includes(host.id)" @click.stop />
              <div class="host-select-info">
                <div class="host-select-name">{{ host.name }}</div>
                <div class="host-select-domain">{{ host.domain }}</div>
              </div>
              <div class="host-select-specs">
                {{ host.cpu }}核 / {{ host.memory }}GB / {{ host.disk }}GB
              </div>
            </div>
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="bindDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="confirmBindHost">
          确认关联 ({{ selectedHosts.length }})
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Monitor } from '@element-plus/icons-vue'
import request from '@/utils/request'

interface Host {
  id: number
  name: string
  domain: string
  cpu: number
  memory: number
  disk: number
  disk_type: string
  status: string
}

interface Contract {
  id: number
  contract_no: string
  product_name: string
  amount: number
  status: string
  start_date: string
  end_date: string
  hosts: Host[]
}

const contracts = ref<Contract[]>([])
const availableHosts = ref<Host[]>([])
const selectedHosts = ref<number[]>([])
const currentContract = ref<Contract | null>(null)
const bindDialogVisible = ref(false)
const submitting = ref(false)

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    active: 'success',
    expired: 'info',
    pending: 'warning'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    active: '生效中',
    expired: '已过期',
    pending: '待签署'
  }
  return map[status] || status
}

const getHostStatusType = (status: string) => {
  const map: Record<string, string> = {
    running: 'success',
    stopped: 'danger',
    suspended: 'warning'
  }
  return map[status] || 'info'
}

const getHostStatusText = (status: string) => {
  const map: Record<string, string> = {
    running: '运行中',
    stopped: '已停止',
    suspended: '已暂停'
  }
  return map[status] || status
}

const showBindHostDialog = (contract: Contract) => {
  currentContract.value = contract
  selectedHosts.value = []
  loadAvailableHosts()
  bindDialogVisible.value = true
}

const toggleHostSelection = (hostId: number) => {
  const index = selectedHosts.value.indexOf(hostId)
  if (index > -1) {
    selectedHosts.value.splice(index, 1)
  } else {
    selectedHosts.value.push(hostId)
  }
}

const confirmBindHost = async () => {
  if (selectedHosts.value.length === 0) {
    ElMessage.warning('请选择要关联的主机')
    return
  }
  submitting.value = true
  try {
    await request.post(`/api/v1/contracts/${currentContract.value?.id}/bind-hosts`, {
      host_ids: selectedHosts.value
    })
    ElMessage.success('主机关联成功')
    bindDialogVisible.value = false
    loadContracts()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.message || '关联失败')
  } finally {
    submitting.value = false
  }
}

const unbindHost = async (contract: Contract, host: Host) => {
  try {
    await ElMessageBox.confirm(
      `确定要解绑主机"${host.name}"吗？`,
      '确认解绑',
      { type: 'warning' }
    )
    await request.post(`/api/v1/contracts/${contract.id}/unbind-host`, {
      host_id: host.id
    })
    contract.hosts = contract.hosts.filter(h => h.id !== host.id)
    ElMessage.success('解绑成功')
  } catch (e: any) {
    if (e !== 'cancel' && e?.message !== 'cancel') {
      ElMessage.error(e.response?.data?.message || '解绑失败')
    }
  }
}

const loadContracts = async () => {
  try {
    const { data } = await request.get('/api/v1/contracts')
    contracts.value = data?.data?.list || data?.data?.items || data?.data || []
  } catch {
    contracts.value = []
  }
}

const loadAvailableHosts = async () => {
  try {
    const { data } = await request.get('/api/v1/user/products')
    availableHosts.value = data?.data?.list || data?.data?.items || data?.data || []
  } catch {
    availableHosts.value = []
  }
}

onMounted(() => {
  loadContracts()
})
</script>

<style scoped lang="scss">
.contract-host-page {
  .contract-list {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .contract-item {
    border: 1px solid #ebeef5;
    border-radius: 8px;
    padding: 20px;

    &:hover {
      border-color: #c0c4cc;
    }

    .contract-header {
      .contract-info {
        display: flex;
        align-items: center;
        gap: 12px;
        margin-bottom: 8px;

        .contract-title {
          font-size: 16px;
          font-weight: 600;
          margin: 0;
        }
      }

      .contract-meta {
        display: flex;
        gap: 24px;
        font-size: 14px;
        color: #606266;
      }
    }

    .host-section {
      .section-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 16px;
        font-weight: 600;
      }

      .host-list {
        display: flex;
        flex-direction: column;
        gap: 12px;
      }

      .host-item {
        display: flex;
        align-items: center;
        gap: 16px;
        padding: 12px 16px;
        background: #f5f7fa;
        border-radius: 6px;

        .host-info {
          display: flex;
          align-items: center;
          gap: 12px;
          flex: 1;

          .host-text {
            .host-name {
              font-weight: 500;
            }

            .host-domain {
              font-size: 13px;
              color: #909399;
            }
          }
        }

        .host-specs {
          display: flex;
          gap: 12px;
          font-size: 13px;
          color: #606266;
        }
      }
    }
  }
}

.bind-host-content {
  h4 {
    margin: 0 0 16px 0;
  }

  .host-select-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
    max-height: 400px;
    overflow-y: auto;
  }

  .host-select-item {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 16px;
    border: 1px solid #ebeef5;
    border-radius: 8px;
    cursor: pointer;
    transition: all 0.3s;

    &:hover {
      border-color: #c0c4cc;
    }

    &.selected {
      border-color: #409eff;
      background: #ecf5ff;
    }

    .host-select-info {
      flex: 1;

      .host-select-name {
        font-weight: 500;
        margin-bottom: 4px;
      }

      .host-select-domain {
        font-size: 13px;
        color: #909399;
      }
    }

    .host-select-specs {
      font-size: 13px;
      color: #606266;
    }
  }
}
</style>
