<template>
  <div class="page-container">
    <art-card :title="$t('serviceDetails.title')" shadow="never">
      <!-- 搜索筛选 -->
      <el-form :model="searchForm" inline class="mb-4">
        <el-form-item :label="$t('serviceDetails.keyword')">
          <el-input v-model="searchForm.keyword" :placeholder="$t('serviceDetails.keywordPlaceholder')" clearable @keyup.enter="handleSearch" />
        </el-form-item>
        <el-form-item :label="$t('serviceDetails.status')">
          <el-select v-model="searchForm.status" :placeholder="$t('serviceDetails.all')" clearable>
            <el-option :label="$t('serviceDetails.active')" :value="1" />
            <el-option :label="$t('serviceDetails.suspended')" :value="2" />
            <el-option :label="$t('serviceDetails.terminated')" :value="3" />
            <el-option :label="$t('serviceDetails.deleted')" :value="0" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">{{ $t('serviceDetails.search') }}</el-button>
          <el-button @click="handleReset">{{ $t('serviceDetails.reset') }}</el-button>
        </el-form-item>
      </el-form>

      <!-- 数据表格 -->
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" :label="$t('serviceDetails.id')" width="80" />
        <el-table-column prop="user_id" :label="$t('serviceDetails.userId')" width="100" />
        <el-table-column prop="username" :label="$t('serviceDetails.username')" width="120" />
        <el-table-column prop="product_name" :label="$t('serviceDetails.product')" min-width="150" />
        <el-table-column prop="domain" :label="$t('serviceDetails.domainIp')" width="150" />
        <el-table-column prop="status" :label="$t('serviceDetails.statusColumn')" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" :label="$t('serviceDetails.createTime')" width="180" />
        <el-table-column prop="due_date" :label="$t('serviceDetails.dueDate')" width="180" />
        <el-table-column :label="$t('serviceDetails.operations')" width="320" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleViewDetail(row)">{{ $t('serviceDetails.detail') }}</el-button>
            <el-button size="small" type="success" @click="handleRenew(row)">{{ $t('serviceDetails.renew') }}</el-button>
            <el-button size="small" type="warning" @click="handleSuspend(row)" :disabled="row.status !== 1">{{ $t('serviceDetails.suspend') }}</el-button>
            <el-button size="small" type="info" @click="handleUnsuspend(row)" :disabled="row.status !== 2">{{ $t('serviceDetails.unsuspend') }}</el-button>
            <el-button size="small" type="danger" @click="handleTerminate(row)" :disabled="row.status === 0 || row.status === 3">{{ $t('serviceDetails.terminate') }}</el-button>
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
    <el-drawer v-model="detailDrawerVisible" :title="$t('serviceDetails.serviceDetail')" size="600px" destroy-on-close>
      <el-descriptions :column="1" border>
        <el-descriptions-item :label="$t('serviceDetails.serviceId')">{{ currentRow.id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('serviceDetails.userId')">{{ currentRow.user_id }}</el-descriptions-item>
        <el-descriptions-item :label="$t('serviceDetails.username')">{{ currentRow.username }}</el-descriptions-item>
        <el-descriptions-item :label="$t('serviceDetails.productName')">{{ currentRow.product_name }}</el-descriptions-item>
        <el-descriptions-item :label="$t('serviceDetails.domainIp')">{{ currentRow.domain }}</el-descriptions-item>
        <el-descriptions-item :label="$t('serviceDetails.statusColumn')">
          <el-tag :type="getStatusType(currentRow.status)">{{ getStatusText(currentRow.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item :label="$t('serviceDetails.createTime')">{{ currentRow.created_at }}</el-descriptions-item>
        <el-descriptions-item :label="$t('serviceDetails.dueDate')">{{ currentRow.due_date }}</el-descriptions-item>
        <el-descriptions-item :label="$t('serviceDetails.remark')">{{ currentRow.remark || '-' }}</el-descriptions-item>
      </el-descriptions>

      <!-- 操作日志 -->
      <div class="mt-4">
        <h4 class="mb-2">{{ $t('serviceDetails.operationLog') }}</h4>
        <el-timeline>
          <el-timeline-item
            v-for="log in serviceLogs"
            :key="log.id"
            :timestamp="log.created_at"
            placement="top"
          >
            <el-card shadow="never">
              <p>{{ log.content }}</p>
              <p class="text-gray-400 text-sm">{{ $t('serviceDetails.username') }}：{{ log.operator }}</p>
            </el-card>
          </el-timeline-item>
        </el-timeline>
        <el-empty v-if="serviceLogs.length === 0" :description="$t('serviceDetails.noLogs')" />
      </div>
    </el-drawer>

    <!-- 续费对话框 -->
    <el-dialog v-model="renewDialogVisible" :title="$t('serviceDetails.serviceRenew')" width="500px" destroy-on-close>
      <el-form :model="renewForm" :rules="renewFormRules" ref="renewFormRef" label-width="100px">
        <el-form-item :label="$t('serviceDetails.serviceId')">
          <el-input :model-value="renewForm.service_id" disabled />
        </el-form-item>
        <el-form-item :label="$t('serviceDetails.username')">
          <el-input :model-value="renewForm.username" disabled />
        </el-form-item>
        <el-form-item :label="$t('serviceDetails.product')">
          <el-input :model-value="renewForm.product_name" disabled />
        </el-form-item>
        <el-form-item :label="$t('serviceDetails.selectDuration')" prop="duration">
          <el-select v-model="renewForm.duration" :placeholder="$t('serviceDetails.selectDuration')" style="width: 100%">
            <el-option :label="$t('serviceDetails.duration1Month')" :value="1" />
            <el-option :label="$t('serviceDetails.duration3Months')" :value="3" />
            <el-option :label="$t('serviceDetails.duration6Months')" :value="6" />
            <el-option :label="$t('serviceDetails.duration12Months')" :value="12" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('serviceDetails.remark')" prop="remark">
          <el-input v-model="renewForm.remark" type="textarea" :rows="3" :placeholder="$t('serviceDetails.renewRemark')" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="renewDialogVisible = false">{{ $t('serviceDetails.cancel') }}</el-button>
        <el-button type="primary" @click="handleRenewSubmit" :loading="renewSubmitLoading">{{ $t('serviceDetails.confirmRenew') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

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
const renewFormRules = computed<FormRules>(() => ({
  duration: [{ required: true, message: $t('serviceDetails.selectDuration'), trigger: 'change' }]
}))

const getStatusType = (status: number) => {
  const map: Record<number, any> = { 0: 'info', 1: 'success', 2: 'warning', 3: 'danger' }
  return map[status] || 'info'
}

const getStatusText = (status: number) => {
  const map: Record<number, string> = {
    0: $t('serviceDetails.deleted'),
    1: $t('serviceDetails.active'),
    2: $t('serviceDetails.suspended'),
    3: $t('serviceDetails.terminated')
  }
  return map[status] || $t('serviceDetails.unknown')
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
      ElMessage.success($t('serviceDetails.renewSuccess'))
      renewDialogVisible.value = false
      fetchData()
    } catch (error) {
      ElMessage.error($t('serviceDetails.renewFailed'))
    } finally {
      renewSubmitLoading.value = false
    }
  })
}

const handleSuspend = async (row: any) => {
  await ElMessageBox.confirm($t('serviceDetails.confirmSuspend'), $t('serviceDetails.tips'), { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/service-details/${row.id}/suspend` })
    ElMessage.success($t('serviceDetails.suspendSuccess'))
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleUnsuspend = async (row: any) => {
  await ElMessageBox.confirm($t('serviceDetails.confirmUnsuspend'), $t('serviceDetails.tips'), { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/service-details/${row.id}/unsuspend` })
    ElMessage.success($t('serviceDetails.unsuspendSuccess'))
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleTerminate = async (row: any) => {
  await ElMessageBox.confirm($t('serviceDetails.confirmTerminate'), $t('serviceDetails.warning'), { type: 'error' })
  try {
    await request.post({ url: `/api/admin/service-details/${row.id}/terminate` })
    ElMessage.success($t('serviceDetails.terminateSuccess'))
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
