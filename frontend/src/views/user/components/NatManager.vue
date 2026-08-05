<template>
  <div class="nat-manager">
    <div class="toolbar">
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        添加NAT规则
      </el-button>
    </div>

    <el-table
      v-loading="loading"
      :data="natList"
      style="width: 100%"
      empty-text="暂无NAT规则"
    >
      <el-table-column prop="name" label="规则名称" min-width="140" show-overflow-tooltip />
      <el-table-column prop="protocol" label="协议" width="80" align="center">
        <template #default="{ row }">
          <el-tag :type="row.protocol === 'tcp' ? 'primary' : 'success'" size="small" effect="light">
            {{ row.protocol?.toUpperCase() || '-' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="外部端口" width="120" align="center">
        <template #default="{ row }">
          <span class="mono">{{ row.external_port || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="内部IP" min-width="140">
        <template #default="{ row }">
          <span class="mono">{{ row.internal_ip || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column label="内部端口" width="120" align="center">
        <template #default="{ row }">
          <span class="mono">{{ row.internal_port || '-' }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100" align="center">
        <template #default="{ row }">
          <el-tag :type="getNatStatusType(row.status)" size="small" effect="light">
            {{ getNatStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="100" align="center" fixed="right">
        <template #default="{ row }">
          <el-button
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

    <!-- 创建NAT规则弹窗 -->
    <el-dialog
      v-model="showCreateDialog"
      title="添加NAT规则"
      width="520px"
      :close-on-click-modal="false"
      @closed="resetCreateForm"
    >
      <el-form :model="createForm" label-width="100px">
        <el-form-item label="规则名称">
          <el-input
            v-model="createForm.name"
            placeholder="请输入规则名称"
            maxlength="50"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="协议" required>
          <el-radio-group v-model="createForm.protocol">
            <el-radio value="tcp">TCP</el-radio>
            <el-radio value="udp">UDP</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="外部端口" required>
          <el-input-number
            v-model="createForm.external_port"
            :min="1"
            :max="65535"
            style="width: 100%"
            placeholder="外部访问端口"
          />
        </el-form-item>
        <el-form-item label="内部IP" required>
          <el-input
            v-model="createForm.internal_ip"
            placeholder="例如: 192.168.1.100"
          />
        </el-form-item>
        <el-form-item label="内部端口" required>
          <el-input-number
            v-model="createForm.internal_port"
            :min="1"
            :max="65535"
            style="width: 100%"
            placeholder="内部服务端口"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="createLoading" @click="confirmCreate">
          确认添加
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
const natList = ref<any[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const showCreateDialog = ref(false)
const createLoading = ref(false)
const createForm = ref({
  name: '',
  protocol: 'tcp',
  external_port: undefined as number | undefined,
  internal_ip: '',
  internal_port: undefined as number | undefined
})

function getNatStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    active: 'success',
    available: 'success',
    creating: 'warning',
    deleting: 'warning',
    error: 'danger'
  }
  return map[status] || 'info'
}

function getNatStatusText(status: string) {
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
    const { data } = await request.get(`/api/v1/hosts/${id}/nat`, {
      params: { page: page.value, limit: pageSize.value }
    })
    if (data?.data) {
      natList.value = data.data.list || []
      total.value = data.data.total || 0
    }
  } catch (error) {
    console.error('获取NAT规则列表失败', error)
  } finally {
    loading.value = false
  }
}

function resetCreateForm() {
  createForm.value = {
    name: '',
    protocol: 'tcp',
    external_port: undefined,
    internal_ip: '',
    internal_port: undefined
  }
}

async function confirmCreate() {
  const id = route.params.id
  if (!id) return

  if (!createForm.value.external_port) {
    ElMessage.warning('请输入外部端口')
    return
  }
  if (!createForm.value.internal_ip.trim()) {
    ElMessage.warning('请输入内部IP')
    return
  }
  if (!createForm.value.internal_port) {
    ElMessage.warning('请输入内部端口')
    return
  }

  createLoading.value = true
  try {
    await request.post(`/api/v1/hosts/${id}/nat`, createForm.value)
    ElMessage.success('NAT规则添加成功')
    showCreateDialog.value = false
    fetchList()
  } catch (error: any) {
    ElMessage.error(error.message || '添加NAT规则失败')
  } finally {
    createLoading.value = false
  }
}

async function handleDelete(row: any) {
  const id = route.params.id
  if (!id) return

  try {
    await ElMessageBox.confirm(
      `确认删除NAT规则「${row.name || `${row.protocol}:${row.external_port}`}」？`,
      '确认删除',
      { type: 'warning', confirmButtonText: '删除', confirmButtonClass: 'el-button--danger' }
    )
    await request.delete(`/api/v1/hosts/${id}/nat/${row.id}`)
    ElMessage.success('NAT规则已删除')
    fetchList()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '删除NAT规则失败')
    }
  }
}

onMounted(() => {
  fetchList()
})
</script>

<style scoped lang="scss">
.nat-manager {
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
