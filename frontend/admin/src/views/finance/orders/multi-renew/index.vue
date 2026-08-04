<template>
  <div class="page-container">
    <art-card title="批量续费管理" shadow="never">
      <!-- 操作栏 -->
      <div class="mb-4 flex justify-between">
        <el-button type="primary" @click="handleCreate">创建批量续费任务</el-button>
        <el-button @click="fetchData" :icon="Refresh">刷新</el-button>
      </div>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="任务ID" width="100" />
        <el-table-column prop="name" label="任务名称" min-width="150" />
        <el-table-column prop="service_count" label="服务数量" width="100" />
        <el-table-column prop="duration" label="续费时长" width="100">
          <template #default="{ row }">{{ row.duration }}个月</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="progress" label="进度" width="150">
          <template #default="{ row }">
            <el-progress :percentage="row.progress || 0" :status="row.status === 3 ? 'success' : undefined" />
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column prop="operator" label="操作人" width="120" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="success" @click="handleExecute(row)" :disabled="row.status !== 0">执行</el-button>
            <el-button size="small" type="danger" @click="handleCancel(row)" :disabled="row.status !== 0 && row.status !== 1">取消</el-button>
            <el-button size="small" @click="handleViewLog(row)">日志</el-button>
            <el-button size="small" type="info" @click="handleDelete(row)" :disabled="row.status === 1">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <!-- 分页 -->
      <el-pagination
        class="mt-4 justify-end"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50, 100]"
        :total="pagination.total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchData"
        @current-change="fetchData"
      />
    </art-card>

    <!-- 创建任务对话框 -->
    <el-dialog v-model="createDialogVisible" title="创建批量续费任务" width="600px" destroy-on-close>
      <el-form :model="createForm" :rules="createFormRules" ref="createFormRef" label-width="120px">
        <el-form-item label="任务名称" prop="name">
          <el-input v-model="createForm.name" placeholder="请输入任务名称" />
        </el-form-item>
        <el-form-item label="续费时长" prop="duration">
          <el-select v-model="createForm.duration" placeholder="请选择续费时长" style="width: 100%">
            <el-option label="1个月" :value="1" />
            <el-option label="3个月" :value="3" />
            <el-option label="6个月" :value="6" />
            <el-option label="12个月" :value="12" />
          </el-select>
        </el-form-item>
        <el-form-item label="筛选条件" prop="filter_type">
          <el-radio-group v-model="createForm.filter_type">
            <el-radio value="ids">按服务ID</el-radio>
            <el-radio value="product">按产品</el-radio>
            <el-radio value="due_date">按到期时间</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item v-if="createForm.filter_type === 'ids'" label="服务ID列表" prop="service_ids">
          <el-input v-model="createForm.service_ids" type="textarea" :rows="4" placeholder="请输入服务ID，多个用逗号分隔" />
        </el-form-item>
        <el-form-item v-if="createForm.filter_type === 'product'" label="选择产品" prop="product_id">
          <el-select v-model="createForm.product_id" placeholder="请选择产品" style="width: 100%">
            <el-option v-for="p in productList" :key="p.id" :label="p.name" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item v-if="createForm.filter_type === 'due_date'" label="到期时间范围" prop="due_date_range">
          <el-date-picker v-model="createForm.due_date_range" type="daterange" range-separator="至" start-placeholder="开始日期" end-placeholder="结束日期" style="width: 100%" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="createForm.remark" type="textarea" :rows="2" placeholder="备注信息（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreateSubmit" :loading="createSubmitLoading">创建任务</el-button>
      </template>
    </el-dialog>

    <!-- 日志抽屉 -->
    <el-drawer v-model="logDrawerVisible" title="任务日志" size="600px" destroy-on-close>
      <el-timeline>
        <el-timeline-item
          v-for="log in taskLogs"
          :key="log.id"
          :timestamp="log.created_at"
          :type="log.type === 'error' ? 'danger' : log.type === 'success' ? 'success' : 'primary'"
          placement="top"
        >
          <p>{{ log.content }}</p>
        </el-timeline-item>
      </el-timeline>
      <el-empty v-if="taskLogs.length === 0" description="暂无日志" />
    </el-drawer>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import request from '@/utils/http'

defineOptions({ name: 'MultiRenew' })

const loading = ref(false)
const tableData = ref([])
const productList = ref<any[]>([])

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 创建对话框
const createDialogVisible = ref(false)
const createSubmitLoading = ref(false)
const createFormRef = ref<FormInstance>()
const createForm = reactive({
  name: '',
  duration: 1 as number,
  filter_type: 'ids' as 'ids' | 'product' | 'due_date',
  service_ids: '',
  product_id: undefined as number | undefined,
  due_date_range: null as any,
  remark: ''
})
const createFormRules: FormRules = {
  name: [{ required: true, message: '请输入任务名称', trigger: 'blur' }],
  duration: [{ required: true, message: '请选择续费时长', trigger: 'change' }],
  filter_type: [{ required: true, message: '请选择筛选条件', trigger: 'change' }]
}

// 日志抽屉
const logDrawerVisible = ref(false)
const taskLogs = ref<any[]>([])

const getStatusType = (status: number) => {
  const map: Record<number, string> = { 0: 'info', 1: 'warning', 2: 'danger', 3: 'success', 4: 'info' }
  return map[status] || 'info'
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = { 0: '待执行', 1: '执行中', 2: '已取消', 3: '已完成', 4: '失败' }
  return map[status] || '未知'
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get({
      url: '/api/admin/multi-renew',
      params: { page: pagination.page, page_size: pagination.pageSize }
    })
    tableData.value = res?.list || []
    pagination.total = res?.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchProducts = async () => {
  try {
    const res = await request.get({ url: '/api/admin/products' })
    productList.value = res || []
  } catch (error) {
    console.error(error)
  }
}

const handleCreate = () => {
  createForm.name = ''
  createForm.duration = 1
  createForm.filter_type = 'ids'
  createForm.service_ids = ''
  createForm.product_id = undefined
  createForm.due_date_range = null
  createForm.remark = ''
  createDialogVisible.value = true
  fetchProducts()
}

const handleCreateSubmit = async () => {
  if (!createFormRef.value) return
  await createFormRef.value.validate(async (valid) => {
    if (!valid) return
    createSubmitLoading.value = true
    try {
      await request.post({ url: '/api/admin/multi-renew', params: createForm })
      ElMessage.success('任务创建成功')
      createDialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error('创建失败')
    } finally {
      createSubmitLoading.value = false
    }
  })
}

const handleExecute = async (row: any) => {
  await ElMessageBox.confirm('确定执行该批量续费任务？', '提示', { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/multi-renew/${row.id}/execute` })
    ElMessage.success('任务已开始执行')
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleCancel = async (row: any) => {
  await ElMessageBox.confirm('确定取消该任务？', '提示', { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/multi-renew/${row.id}/cancel` })
    ElMessage.success('任务已取消')
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleViewLog = async (row: any) => {
  logDrawerVisible.value = true
  try {
    const res = await request.get({ url: `/api/admin/multi-renew/${row.id}/logs` })
    taskLogs.value = res || []
  } catch (error) {
    taskLogs.value = []
  }
}

const handleDelete = async (row: any) => {
  await ElMessageBox.confirm('确定删除该任务？', '警告', { type: 'error' })
  try {
    await request.del({ url: `/api/admin/multi-renew/${row.id}` })
    ElMessage.success('删除成功')
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

onMounted(() => fetchData())
</script>

<style scoped lang="scss">
.page-container {
  padding: 20px;
}
</style>
