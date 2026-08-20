<template>
  <div class="page-container">
    <art-card :title="$t('productDiverts.title')" shadow="never">
      <el-tabs v-model="activeTab">
        <!-- 转移请求列表 -->
        <el-tab-pane :label="$t('productDiverts.requestTab')" name="requests">
          <el-table :data="requestList" v-loading="loading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="user_id" :label="$t('productDiverts.userId')" width="100" />
            <el-table-column prop="username" :label="$t('common.username')" width="120" />
            <el-table-column prop="service_id" :label="$t('productDiverts.serviceId')" width="100" />
            <el-table-column prop="source_product" :label="$t('productDiverts.sourceProduct')" min-width="150" />
            <el-table-column prop="target_product" :label="$t('productDiverts.targetProduct')" min-width="150" />
            <el-table-column prop="status" :label="$t('common.status')" width="100">
              <template #default="{ row }">
                <el-tag :type="getRequestStatusType(row.status)">{{ getRequestStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" :label="$t('productDiverts.applyTime')" width="180" />
            <el-table-column :label="$t('common.action')" width="200" fixed="right">
              <template #default="{ row }">
                <template v-if="row.status === 0">
                  <el-button size="small" type="success" @click="handleApprove(row)">{{ $t('productDiverts.accept') }}</el-button>
                  <el-button size="small" type="danger" @click="handleReject(row)">{{ $t('productDiverts.reject') }}</el-button>
                </template>
                <span v-else class="text-gray-400">{{ $t('productDiverts.processed') }}</span>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            class="mt-4 justify-end"
            v-model:current-page="requestPagination.page"
            v-model:page-size="requestPagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="requestPagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="fetchRequests"
            @current-change="fetchRequests"
          />
        </el-tab-pane>

        <!-- 转移码管理 -->
        <el-tab-pane :label="$t('productDiverts.codeTab')" name="codes">
          <div class="mb-4">
            <el-button type="primary" @click="handleGenerateCode">{{ $t('productDiverts.generateCode') }}</el-button>
          </div>
          <el-table :data="codeList" v-loading="codesLoading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="code" :label="$t('productDiverts.code')" width="200">
              <template #default="{ row }">
                <el-text class="font-mono">{{ row.code }}</el-text>
                <el-button size="small" link @click="handleCopyCode(row.code)">{{ $t('productDiverts.copy') }}</el-button>
              </template>
            </el-table-column>
            <el-table-column prop="service_id" :label="$t('productDiverts.serviceId')" width="100" />
            <el-table-column prop="username" :label="$t('common.username')" width="120" />
            <el-table-column prop="status" :label="$t('common.status')" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 0 ? 'success' : 'info'">{{ row.status === 0 ? $t('productDiverts.unused') : $t('productDiverts.used') }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="expires_at" :label="$t('productDiverts.expireTime')" width="180" />
            <el-table-column prop="created_at" :label="$t('common.createdAt')" width="180" />
            <el-table-column :label="$t('common.action')" width="150" fixed="right">
              <template #default="{ row }">
                <el-button size="small" type="warning" @click="handleRegenerateCode(row)" :disabled="row.status === 1">{{ $t('productDiverts.regenerate') }}</el-button>
              </template>
            </el-table-column>
          </el-table>
          <el-pagination
            class="mt-4 justify-end"
            v-model:current-page="codePagination.page"
            v-model:page-size="codePagination.pageSize"
            :page-sizes="[10, 20, 50, 100]"
            :total="codePagination.total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="fetchCodes"
            @current-change="fetchCodes"
          />
        </el-tab-pane>
      </el-tabs>
    </art-card>

    <!-- 生成转移码对话框 -->
    <el-dialog v-model="generateDialogVisible" :title="$t('productDiverts.generateCode')" width="500px" destroy-on-close>
      <el-form :model="generateForm" :rules="generateFormRules" ref="generateFormRef" label-width="100px">
        <el-form-item :label="$t('productDiverts.serviceId')" prop="service_id">
          <el-input v-model.number="generateForm.service_id" :placeholder="$t('productDiverts.enterServiceId')" />
        </el-form-item>
        <el-form-item :label="$t('productDiverts.validity')" prop="expires_days">
          <el-select v-model="generateForm.expires_days" :placeholder="$t('productDiverts.selectValidity')" style="width: 100%">
            <el-option :label="$t('productDiverts.day1')" :value="1" />
            <el-option :label="$t('productDiverts.day3')" :value="3" />
            <el-option :label="$t('productDiverts.day7')" :value="7" />
            <el-option :label="$t('productDiverts.day30')" :value="30" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="generateDialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleGenerateSubmit" :loading="generateSubmitLoading">{{ $t('productDiverts.generate') }}</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'
import { $t } from '@/locales'

defineOptions({ name: 'ProductDiverts' })

const activeTab = ref('requests')

// 转移请求
const loading = ref(false)
const requestList = ref([])
const requestPagination = reactive({ page: 1, pageSize: 20, total: 0 })

// 转移码
const codesLoading = ref(false)
const codeList = ref([])
const codePagination = reactive({ page: 1, pageSize: 20, total: 0 })

// 生成对话框
const generateDialogVisible = ref(false)
const generateSubmitLoading = ref(false)
const generateFormRef = ref<FormInstance>()
const generateForm = reactive({
  service_id: undefined as number | undefined,
  expires_days: 7 as number
})
const generateFormRules: FormRules = {
  service_id: [{ required: true, message: $t('productDiverts.enterServiceId'), trigger: 'blur' }],
  expires_days: [{ required: true, message: $t('productDiverts.selectValidity'), trigger: 'change' }]
}

const getRequestStatusType = (status: number) => {
  const map: Record<number, any> = { 0: 'warning', 1: 'success', 2: 'danger' }
  return map[status] || 'info'
}

const getRequestStatusText = (status: number) => {
  const map: Record<number, () => string> = { 0: () => $t('productDiverts.pending'), 1: () => $t('productDiverts.approved'), 2: () => $t('productDiverts.rejected') }
  return map[status]?.() || $t('common.unknown')
}

const fetchRequests = async () => {
  loading.value = true
  try {
    const res = await request.get({
      url: '/api/admin/product-diverts/requests',
      params: { page: requestPagination.page, page_size: requestPagination.pageSize }
    })
    requestList.value = res?.list || []
    requestPagination.total = res?.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const fetchCodes = async () => {
  codesLoading.value = true
  try {
    const res = await request.get({
      url: '/api/admin/product-diverts/codes',
      params: { page: codePagination.page, page_size: codePagination.pageSize }
    })
    codeList.value = res?.list || []
    codePagination.total = res?.total || 0
  } catch (error) {
    console.error(error)
  } finally {
    codesLoading.value = false
  }
}

const handleApprove = async (row: any) => {
  await ElMessageBox.confirm($t('productDiverts.confirmAccept'), $t('common.tips'), { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/product-diverts/requests/${row.id}/approve` })
    ElMessage.success($t('productDiverts.accepted'))
    fetchRequests()
  } catch (error) {
    console.error(error)
  }
}

const handleReject = async (row: any) => {
  await ElMessageBox.confirm($t('productDiverts.confirmReject'), $t('common.tips'), { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/product-diverts/requests/${row.id}/reject` })
    ElMessage.success($t('productDiverts.rejected'))
    fetchRequests()
  } catch (error) {
    console.error(error)
  }
}

const handleGenerateCode = () => {
  generateForm.service_id = undefined
  generateForm.expires_days = 7
  generateDialogVisible.value = true
}

const handleGenerateSubmit = async () => {
  if (!generateFormRef.value) return
  await generateFormRef.value.validate(async (valid) => {
    if (!valid) return
    generateSubmitLoading.value = true
    try {
      await request.post({ url: '/api/admin/product-diverts/codes', params: generateForm })
      ElMessage.success($t('productDiverts.generateSuccess'))
      generateDialogVisible.value = false
      fetchCodes()
    } catch (error) {
      ElMessage.error($t('productDiverts.generateFailed'))
    } finally {
      generateSubmitLoading.value = false
    }
  })
}

const handleRegenerateCode = async (row: any) => {
  await ElMessageBox.confirm($t('productDiverts.confirmRegenerate'), $t('common.tips'), { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/product-diverts/codes/${row.id}/regenerate` })
    ElMessage.success($t('productDiverts.regenerateSuccess'))
    fetchCodes()
  } catch (error) {
    console.error(error)
  }
}

const handleCopyCode = (code: string) => {
  navigator.clipboard.writeText(code)
  ElMessage.success($t('productDiverts.copied'))
}

onMounted(() => {
  fetchRequests()
  fetchCodes()
})
</script>

<style scoped lang="scss">
.page-container {
  padding: 20px;
}
</style>
