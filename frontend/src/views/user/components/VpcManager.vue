<template>
  <div class="vpc-manager">
    <div class="toolbar">
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        创建VPC
      </el-button>
    </div>

    <el-table
      v-loading="loading"
      :data="vpcList"
      style="width: 100%"
      empty-text="暂无VPC网络"
    >
      <el-table-column prop="name" label="VPC名称" min-width="150" show-overflow-tooltip />
      <el-table-column prop="cidr" label="CIDR" width="180">
        <template #default="{ row }">
          <span class="mono">{{ row.cidr || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="vlan_id" label="VLAN ID" width="100" align="center" />
      <el-table-column prop="status" label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getVpcStatusType(row.status)" size="small" effect="light">
            {{ getVpcStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="is_default" label="默认" width="80" align="center">
        <template #default="{ row }">
          <el-tag v-if="row.is_default" type="success" size="small" effect="plain">是</el-tag>
          <span v-else>-</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="220" align="center" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="!row.is_default"
            link
            type="primary"
            size="small"
            @click="handleSwitch(row)"
          >
            切换为默认
          </el-button>
          <el-button
            v-if="!row.is_default"
            link
            type="danger"
            size="small"
            @click="handleDelete(row)"
          >
            删除
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div v-if="total > 0" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next"
        @size-change="fetchList"
        @current-change="fetchList"
      />
    </div>

    <!-- 创建VPC弹窗 -->
    <el-dialog
      v-model="showCreateDialog"
      title="创建VPC"
      width="520px"
      :close-on-click-modal="false"
      @closed="resetCreateForm"
    >
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="VPC名称" required>
          <el-input
            v-model="createForm.name"
            placeholder="请输入VPC名称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="CIDR" required>
          <el-input
            v-model="createForm.cidr"
            placeholder="例如: 10.0.0.0/24"
          />
        </el-form-item>
        <el-form-item label="描述">
          <el-input
            v-model="createForm.description"
            type="textarea"
            :rows="3"
            placeholder="请输入VPC描述（可选）"
            maxlength="200"
            show-word-limit
          />
        </el-form-item>
        <el-alert type="info" :closable="false" show-icon>
          <template #default>
            <p>CIDR 常用范围：</p>
            <p style="margin-top: 4px; font-size: 12px; color: #909399;">
              10.0.0.0/8、172.16.0.0/12、192.168.0.0/16
            </p>
          </template>
        </el-alert>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="createLoading" @click="confirmCreate">
          确认创建
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import request from '@/utils/request'

const route = useRoute()

const loading = ref(false)
const vpcList = ref<any[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const showCreateDialog = ref(false)
const createLoading = ref(false)
const createForm = ref({
  name: '',
  cidr: '',
  description: ''
})

function getVpcStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    active: 'success',
    available: 'success',
    creating: 'warning',
    deleting: 'warning',
    error: 'danger'
  }
  return map[status] || 'info'
}

function getVpcStatusText(status: string) {
  const map: Record<string, string> = {
    active: '活跃',
    available: '可用',
    creating: '创建中',
    deleting: '删除中',
    error: '异常'
  }
  return map[status] || status
}

async function fetchList() {
  const id = route.params.id
  if (!id) return

  loading.value = true
  try {
    const { data } = await request.get(`/api/v2/hosts/${id}/vpcs`, {
      params: { page: page.value, limit: pageSize.value }
    })
    if (data?.data) {
      vpcList.value = data.data.list || []
      total.value = data.data.total || 0
    }
  } catch (error) {
    console.error('获取VPC列表失败', error)
  } finally {
    loading.value = false
  }
}

function resetCreateForm() {
  createForm.value = { name: '', cidr: '', description: '' }
}

async function confirmCreate() {
  const id = route.params.id
  if (!id) return

  if (!createForm.value.name.trim()) {
    ElMessage.warning('请输入VPC名称')
    return
  }
  if (!createForm.value.cidr.trim()) {
    ElMessage.warning('请输入CIDR')
    return
  }

  createLoading.value = true
  try {
    await request.post(`/api/v2/hosts/${id}/vpcs`, createForm.value)
    ElMessage.success('VPC创建成功')
    showCreateDialog.value = false
    fetchList()
  } catch (error: any) {
    ElMessage.error(error.message || '创建VPC失败')
  } finally {
    createLoading.value = false
  }
}

async function handleSwitch(row: any) {
  const id = route.params.id
  if (!id) return

  try {
    await ElMessageBox.confirm(
      `确认将「${row.name}」切换为默认VPC？切换后服务器将使用新的网络配置。`,
      '确认切换',
      { type: 'warning' }
    )
    await request.post(`/api/v2/hosts/${id}/vpcs/${row.id}/switch`)
    ElMessage.success('VPC切换成功')
    fetchList()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '切换VPC失败')
    }
  }
}

async function handleDelete(row: any) {
  const id = route.params.id
  if (!id) return

  try {
    await ElMessageBox.confirm(
      `确认删除VPC「${row.name}」？删除后该VPC下的网络配置将被清除。`,
      '确认删除',
      { type: 'warning', confirmButtonText: '删除', confirmButtonClass: 'el-button--danger' }
    )
    await request.delete(`/api/v2/hosts/${id}/vpcs/${row.id}`)
    ElMessage.success('VPC已删除')
    fetchList()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除VPC失败')
    }
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.vpc-manager {
  .toolbar {
    display: flex;
    justify-content: flex-end;
    margin-bottom: 16px;
  }
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.mono {
  font-family: 'Monaco', 'Menlo', monospace;
}
</style>
