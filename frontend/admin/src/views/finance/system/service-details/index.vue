<template>
  <div class="page-container">
    <art-card title="服务详情管理" shadow="never">
      <!-- 搜索筛选 -->
      <el-form :model="searchForm" inline class="mb-4">
        <el-form-item label="关键词">
          <el-input v-model="searchForm.keyword" placeholder="用户名/产品名/服务ID" clearable @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="全部" clearable>
            <el-option label="使用中" :value="1" />
            <el-option label="已暂停" :value="2" />
            <el-option label="已终止" :value="3" />
            <el-option label="已删除" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">搜索</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="product_name" label="产品" min-width="150" />
        <el-table-column prop="domain" label="域名/IP" width="150" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column prop="due_date" label="到期时间" width="180" />
        <el-table-column label="操作" width="320" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleViewDetail(row)">详情</el-button>
            <el-button size="small" type="success" @click="handleRenew(row)">续费</el-button>
            <el-button size="small" type="warning" @click="handleSuspend(row)" :disabled="row.status !== 1">暂停</el-button>
            <el-button size="small" type="info" @click="handleUnsuspend(row)" :disabled="row.status !== 2">解暂停</el-button>
            <el-button size="small" type="danger" @click="handleTerminate(row)" :disabled="row.status === 0 || row.status === 3">终止</el-button>
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

    <!-- 服务详情抽屉 -->
    <el-drawer v-model="detailDrawerVisible" title="服务详情" size="600px" destroy-on-close>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="服务ID">{{ currentRow.id }}</el-descriptions-item>
        <el-descriptions-item label="用户ID">{{ currentRow.user_id }}</el-descriptions-item>
        <el-descriptions-item label="用户名">{{ currentRow.username }}</el-descriptions-item>
        <el-descriptions-item label="产品名称">{{ currentRow.product_name }}</el-descriptions-item>
        <el-descriptions-item label="域名/IP">{{ currentRow.domain }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentRow.status)">{{ getStatusText(currentRow.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ currentRow.created_at }}</el-descriptions-item>
        <el-descriptions-item label="到期时间">{{ currentRow.due_date }}</el-descriptions-item>
        <el-descriptions-item label="备注">{{ currentRow.remark || '-' }}</el-descriptions-item>
      </el-descriptions>

      <!-- 操作日志 -->
      <div class="mt-4">
        <h4 class="mb-2">操作日志</h4>
        <el-timeline>
          <el-timeline-item
            v-for="log in serviceLogs"
            :key="log.id"
            :timestamp="log.created_at"
            placement="top"
          >
            <el-card shadow="never">
              <p>{{ log.content }}</p>
              <p class="text-gray-400 text-sm">操作人：{{ log.operator }}</p>
            </el-card>
          </el-timeline-item>
        </el-timeline>
        <el-empty v-if="serviceLogs.length === 0" description="暂无日志" />
      </div>
    </el-drawer>

    <!-- 续费对话框 -->
    <el-dialog v-model="renewDialogVisible" title="服务续费" width="500px" destroy-on-close>
      <el-form :model="renewForm" :rules="renewFormRules" ref="renewFormRef" label-width="100px">
        <el-form-item label="服务ID">
          <el-input :model-value="renewForm.service_id" disabled />
        </el-form-item>
        <el-form-item label="用户名">
          <el-input :model-value="renewForm.username" disabled />
        </el-form-item>
        <el-form-item label="产品">
          <el-input :model-value="renewForm.product_name" disabled />
        </el-form-item>
        <el-form-item label="续费时长" prop="duration">
          <el-select v-model="renewForm.duration" placeholder="请选择续费时长" style="width: 100%">
            <el-option label="1个月" :value="1" />
            <el-option label="3个月" :value="3" />
            <el-option label="6个月" :value="6" />
            <el-option label="12个月" :value="12" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input v-model="renewForm.remark" type="textarea" :rows="3" placeholder="续费备注（可选）" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="renewDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleRenewSubmit" :loading="renewSubmitLoading">确认续费</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

defineOptions({ name: 'ServiceDetails' })

const loading = ref(false)
const tableData = ref([])

// 搜索表单
const searchForm = reactive({
  keyword: '',
  status: undefined as number | undefined
})

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 详情抽屉
const detailDrawerVisible = ref(false)
const currentRow = ref<any>({})
const serviceLogs = ref<any[]>([])

// 续费对话框
const renewDialogVisible = ref(false)
const renewSubmitLoading = ref(false)
const renewFormRef = ref<FormInstance>()
const renewForm = reactive({
  service_id: 0,
  username: '',
  product_name: '',
  duration: 1 as number,
  remark: ''
})
const renewFormRules: FormRules = {
  duration: [{ required: true, message: '请选择续费时长', trigger: 'change' }]
}

const getStatusType = (status: number) => {
  const map: Record<number, string> = { 0: 'info', 1: 'success', 2: 'warning', 3: 'danger' }
  return map[status] || 'info'
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = { 0: '已删除', 1: '使用中', 2: '已暂停', 3: '已终止' }
  return map[status] || '未知'
}

const fetchData = async () => {
  loading.value = true
  try {
    const res = await request.get({
      url: '/api/admin/service-details',
      params: {
        page: pagination.page,
        page_size: pagination.pageSize,
        ...searchForm
      }
    })
    tableData.value = res?.list || []
    pagination.total = res?.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchData()
}

const handleReset = () => {
  searchForm.keyword = ''
  searchForm.status = undefined
  handleSearch()
}

const handleViewDetail = async (row: any) => {
  currentRow.value = row
  detailDrawerVisible.value = true
  try {
    const res = await request.get({ url: `/api/admin/service-details/${row.id}/logs` })
    serviceLogs.value = res || []
  } catch (error) {
    serviceLogs.value = []
  }
}

const handleRenew = (row: any) => {
  renewForm.service_id = row.id
  renewForm.username = row.username
  renewForm.product_name = row.product_name
  renewForm.duration = 1
  renewForm.remark = ''
  renewDialogVisible.value = true
}

const handleRenewSubmit = async () => {
  if (!renewFormRef.value) return
  await renewFormRef.value.validate(async (valid) => {
    if (!valid) return
    renewSubmitLoading.value = true
    try {
      await request.post({
        url: `/api/admin/service-details/${renewForm.service_id}/renew`,
        params: { duration: renewForm.duration, remark: renewForm.remark }
      })
      ElMessage.success('续费成功')
      renewDialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error('续费失败')
    } finally {
      renewSubmitLoading.value = false
    }
  })
}

const handleSuspend = async (row: any) => {
  await ElMessageBox.confirm('确定暂停该服务？', '提示', { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/service-details/${row.id}/suspend` })
    ElMessage.success('暂停成功')
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleUnsuspend = async (row: any) => {
  await ElMessageBox.confirm('确定解除暂停该服务？', '提示', { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/service-details/${row.id}/unsuspend` })
    ElMessage.success('解暂停成功')
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleTerminate = async (row: any) => {
  await ElMessageBox.confirm('确定终止该服务？此操作不可逆！', '警告', { type: 'error' })
  try {
    await request.post({ url: `/api/admin/service-details/${row.id}/terminate` })
    ElMessage.success('终止成功')
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
