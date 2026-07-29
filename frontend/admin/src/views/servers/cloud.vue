<template>
  <div class="cloud-servers-page page-container">
    <div class="search-bar">
      <el-form :inline="true" :model="searchForm">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="名称/IP" clearable />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="运行中" value="active" />
            <el-option label="已暂停" value="suspended" />
            <el-option label="已停止" value="stopped" />
            <el-option label="创建中" value="creating" />
          </el-select>
        </el-form-item>
        <el-form-item label="用户ID">
          <el-input v-model="searchForm.user_id" placeholder="用户ID" clearable />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="resetSearch">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="art-card">
      <div class="table-header">
        <h3>云服务器管理</h3>
        <div class="header-actions">
          <el-button type="primary" @click="handleBatchAction('start')" :disabled="!selectedIds.length">
            <el-icon><VideoPlay /></el-icon>批量启动
          </el-button>
          <el-button type="warning" @click="handleBatchAction('stop')" :disabled="!selectedIds.length">
            <el-icon><VideoPause /></el-icon>批量停止
          </el-button>
          <el-button type="danger" @click="handleBatchDelete" :disabled="!selectedIds.length">
            <el-icon><Delete /></el-icon>批量删除
          </el-button>
        </div>
      </div>

      <el-table :data="list" style="width: 100%" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="50" />
        <el-table-column prop="id" label="ID" width="70" />
        <el-table-column prop="name" label="名称" min-width="120" show-overflow-tooltip />
        <el-table-column prop="ip" label="IP地址" width="140" />
        <el-table-column prop="cpu" label="CPU" width="70">
          <template #default="{ row }">{{ row.cpu }}核</template>
        </el-table-column>
        <el-table-column prop="memory" label="内存" width="80">
          <template #default="{ row }">{{ row.memory }}MB</template>
        </el-table-column>
        <el-table-column prop="disk" label="硬盘" width="80">
          <template #default="{ row }">{{ row.disk }}GB</template>
        </el-table-column>
        <el-table-column prop="user_id" label="所属用户" width="100" />
        <el-table-column prop="datacenter_name" label="机房" width="120" show-overflow-tooltip />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="cloudStatusType(row.status)" size="small">
              {{ cloudStatusLabel(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="expire_at" label="到期时间" width="170" />
        <el-table-column label="操作" width="220" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewDetail(row)">详情</el-button>
            <el-button
              :type="row.status === 'active' ? 'warning' : 'success'"
              link
              @click="toggleStatus(row)"
            >
              {{ row.status === 'active' ? '暂停' : '启动' }}
            </el-button>
            <el-button type="danger" link @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </div>

    <el-dialog v-model="detailVisible" title="云服务器详情" width="700px" destroy-on-close>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="ID">{{ detailData.id }}</el-descriptions-item>
        <el-descriptions-item label="名称">{{ detailData.name }}</el-descriptions-item>
        <el-descriptions-item label="IP地址">{{ detailData.ip }}</el-descriptions-item>
        <el-descriptions-item label="所属用户">{{ detailData.user_id }}</el-descriptions-item>
        <el-descriptions-item label="CPU">{{ detailData.cpu }}核</el-descriptions-item>
        <el-descriptions-item label="内存">{{ detailData.memory }}MB</el-descriptions-item>
        <el-descriptions-item label="硬盘">{{ detailData.disk }}GB</el-descriptions-item>
        <el-descriptions-item label="带宽">{{ detailData.bandwidth }}Mbps</el-descriptions-item>
        <el-descriptions-item label="机房">{{ detailData.datacenter_name }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="cloudStatusType(detailData.status)" size="small">
            {{ cloudStatusLabel(detailData.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ detailData.created_at }}</el-descriptions-item>
        <el-descriptions-item label="到期时间">{{ detailData.expire_at }}</el-descriptions-item>
      </el-descriptions>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { VideoPlay, VideoPause, Delete } from '@element-plus/icons-vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const list = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(20)
const total = ref(0)
const selectedIds = ref<number[]>([])

const searchForm = ref({ keyword: '', status: '', user_id: '' })

const detailVisible = ref(false)
const detailData = reactive<any>({})

const cloudStatusLabel = (status: string) => {
  const map: Record<string, string> = { active: '运行中', suspended: '已暂停', stopped: '已停止', creating: '创建中' }
  return map[status] || status
}

const cloudStatusType = (status: string) => {
  const map: Record<string, string> = { active: 'success', suspended: 'danger', stopped: 'info', creating: 'warning' }
  return (map[status] as any) || 'info'
}

const fetchList = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/admin/api/v1/servers/cloud', {
      params: { page: currentPage.value, page_size: pageSize.value, ...searchForm.value }
    })
    list.value = data.data?.list || []
    total.value = data.data?.total || 0
  } catch { /* handled */ } finally {
    loading.value = false
  }
}

const handleSearch = () => { currentPage.value = 1; fetchList() }

const resetSearch = () => {
  searchForm.value = { keyword: '', status: '', user_id: '' }
  handleSearch()
}

const handleSelectionChange = (rows: any[]) => {
  selectedIds.value = rows.map((r) => r.id)
}

const viewDetail = (row: any) => {
  Object.assign(detailData, row)
  detailVisible.value = true
}

const toggleStatus = async (row: any) => {
  const newStatus = row.status === 'active' ? 'suspended' : 'active'
  const action = newStatus === 'active' ? '启动' : '暂停'
  await ElMessageBox.confirm(`确定要${action}云服务器「${row.name}」吗？`, '提示', { type: 'warning' })
  try {
    await request.put(`/admin/api/v1/servers/cloud/${row.id}/status`, { status: newStatus })
    ElMessage.success(`${action}成功`)
    fetchList()
  } catch { /* handled */ }
}

const handleBatchAction = async (action: string) => {
  const label = action === 'start' ? '启动' : '停止'
  await ElMessageBox.confirm(`确定要批量${label} ${selectedIds.value.length} 台服务器吗？`, '提示', { type: 'warning' })
  try {
    await request.post('/admin/api/v1/servers/cloud/batch-status', {
      ids: selectedIds.value,
      status: action === 'start' ? 'active' : 'stopped'
    })
    ElMessage.success(`批量${label}成功`)
    fetchList()
  } catch { /* handled */ }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm(`确定要删除云服务器「${row.name}」吗？`, '警告', { type: 'error' })
  try {
    await request.delete(`/admin/api/v1/servers/cloud/${row.id}`)
    ElMessage.success('删除成功')
    fetchList()
  } catch { /* handled */ }
}

const handleBatchDelete = async () => {
  await ElMessageBox.confirm(`确定要删除选中的 ${selectedIds.value.length} 台服务器吗？`, '警告', { type: 'error' })
  try {
    await request.post('/admin/api/v1/servers/cloud/batch-delete', { ids: selectedIds.value })
    ElMessage.success('批量删除成功')
    fetchList()
  } catch { /* handled */ }
}

onMounted(fetchList)
</script>

<style scoped lang="scss">
.cloud-servers-page {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    h3 { margin: 0; font-size: 16px; font-weight: 600; }
    .header-actions { display: flex; gap: 8px; }
  }
  .pagination { margin-top: 20px; display: flex; justify-content: flex-end; }
}
</style>
