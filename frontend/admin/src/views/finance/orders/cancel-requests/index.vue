<template>
  <div class="cancel-requests-page">
    <art-card title="取消请求" shadow="never">
      <template #header>
        <div class="card-header">
          <span>产品取消请求</span>
          <el-button @click="showReasonDialog">取消原因管理</el-button>
        </div>
      </template>

      <el-table :data="cancelList" v-loading="loading" stripe border>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="username" label="用户(公司)" />
        <el-table-column prop="create_time" label="请求时间" width="180" />
        <el-table-column prop="product_name" label="产品(主机名)" />
        <el-table-column prop="ip" label="IP" width="140" />
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="row.type === 'Immediate' ? 'danger' : 'warning'">
              {{ row.type === 'Immediate' ? '立即' : '到期' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="原因" show-overflow-tooltip />
        <el-table-column prop="delete_time" label="删除时间" width="180" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'">
              {{ row.status === 1 ? '已处理' : '待处理' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="handleConfirm(row)" v-if="!row.status">确认取消</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.limit"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next"
        style="margin-top: 16px; justify-content: flex-end"
        @size-change="fetchCancelList"
        @current-change="fetchCancelList"
      />
    </art-card>

    <!-- 取消原因管理对话框 -->
    <el-dialog v-model="reasonDialogVisible" title="取消原因管理" width="600px">
      <div class="reason-toolbar">
        <el-button type="primary" size="small" @click="showAddReason">添加原因</el-button>
      </div>
      <el-table :data="reasonList" stripe border size="small">
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column prop="reason" label="原因" />
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" type="danger" @click="handleDeleteReason(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-dialog v-model="addReasonVisible" title="添加原因" width="400px" append-to-body>
        <el-input v-model="newReason" placeholder="请输入取消原因" />
        <template #footer>
          <el-button @click="addReasonVisible = false">取消</el-button>
          <el-button type="primary" @click="handleAddReason">确定</el-button>
        </template>
      </el-dialog>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/http'

const loading = ref(false)
const reasonDialogVisible = ref(false)
const addReasonVisible = ref(false)
const newReason = ref('')

const pagination = reactive({ page: 1, limit: 10, total: 0 })
const cancelList = ref<any[]>([])
const reasonList = ref<any[]>([])

const fetchCancelList = async () => {
  loading.value = true
  try {
    const res = await request.get({
      url: '/api/admin/orders/cancel-requests',
      params: { page: pagination.page, limit: pagination.limit }
    })
    if (res?.data) {
      cancelList.value = res.data.list || []
      pagination.total = res.data.total || 0
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleConfirm = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定确认此取消请求吗？产品将被取消。', '提示')
    await request.put({ url: `/api/admin/orders/cancel-requests/${row.id}/confirm`, showSuccessMessage: true })
    fetchCancelList()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('操作失败')
  }
}

const showReasonDialog = async () => {
  reasonDialogVisible.value = true
  try {
    const res = await request.get({ url: '/api/admin/orders/cancel-reasons' })
    reasonList.value = res?.data || []
  } catch (error) {
    console.error(error)
  }
}

const showAddReason = () => {
  newReason.value = ''
  addReasonVisible.value = true
}

const handleAddReason = async () => {
  if (!newReason.value) return
  try {
    await request.post({ url: '/api/admin/orders/cancel-reasons', data: { reason: newReason.value }, showSuccessMessage: true })
    addReasonVisible.value = false
    showReasonDialog()
  } catch (error) {
    ElMessage.error('添加失败')
  }
}

const handleDeleteReason = async (row: any) => {
  try {
    await ElMessageBox.confirm('确定删除此原因吗？', '提示')
    await request.delete({ url: `/api/admin/orders/cancel-reasons/${row.id}`, showSuccessMessage: true })
    showReasonDialog()
  } catch (error) {
    if (error !== 'cancel') ElMessage.error('删除失败')
  }
}

onMounted(() => fetchCancelList())
</script>

<style scoped lang="scss">
.cancel-requests-page {
  padding: 20px;
}
.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.reason-toolbar {
  margin-bottom: 12px;
}
</style>
