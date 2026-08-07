<template>
  <div class="page-container">
    <art-card title="产品转移/分流" shadow="never">
      <el-tabs v-model="activeTab">
        <!-- 转移请求列表 -->
        <el-tab-pane label="转移请求" name="requests">
          <el-table :data="requestList" v-loading="loading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="user_id" label="用户ID" width="100" />
            <el-table-column prop="username" label="用户名" width="120" />
            <el-table-column prop="service_id" label="服务ID" width="100" />
            <el-table-column prop="source_product" label="原产品" min-width="150" />
            <el-table-column prop="target_product" label="目标产品" min-width="150" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getRequestStatusType(row.status)">{{ getRequestStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="申请时间" width="180" />
            <el-table-column label="操作" width="200" fixed="right">
              <template #default="{ row }">
                <template v-if="row.status === 0">
                  <el-button size="small" type="success" @click="handleApprove(row)">接受</el-button>
                  <el-button size="small" type="danger" @click="handleReject(row)">拒绝</el-button>
                </template>
                <span v-else class="text-gray-400">已处理</span>
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
        <el-tab-pane label="转移码管理" name="codes">
          <div class="mb-4">
            <el-button type="primary" @click="handleGenerateCode">生成转移码</el-button>
          </div>
          <el-table :data="codeList" v-loading="codesLoading" stripe>
            <el-table-column prop="id" label="ID" width="80" />
            <el-table-column prop="code" label="转移码" width="200">
              <template #default="{ row }">
                <el-text class="font-mono">{{ row.code }}</el-text>
                <el-button size="small" link @click="handleCopyCode(row.code)">复制</el-button>
              </template>
            </el-table-column>
            <el-table-column prop="service_id" label="服务ID" width="100" />
            <el-table-column prop="username" label="用户名" width="120" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 0 ? 'success' : 'info'">{{ row.status === 0 ? '未使用' : '已使用' }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="expires_at" label="过期时间" width="180" />
            <el-table-column prop="created_at" label="创建时间" width="180" />
            <el-table-column label="操作" width="150" fixed="right">
              <template #default="{ row }">
                <el-button size="small" type="warning" @click="handleRegenerateCode(row)" :disabled="row.status === 1">重新生成</el-button>
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
    <el-dialog v-model="generateDialogVisible" title="生成转移码" width="500px" destroy-on-close>
      <el-form :model="generateForm" :rules="generateFormRules" ref="generateFormRef" label-width="100px">
        <el-form-item label="服务ID" prop="service_id">
          <el-input v-model.number="generateForm.service_id" placeholder="请输入服务ID" />
        </el-form-item>
        <el-form-item label="有效期" prop="expires_days">
          <el-select v-model="generateForm.expires_days" placeholder="请选择有效期" style="width: 100%">
            <el-option label="1天" :value="1" />
            <el-option label="3天" :value="3" />
            <el-option label="7天" :value="7" />
            <el-option label="30天" :value="30" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="generateDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleGenerateSubmit" :loading="generateSubmitLoading">生成</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

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
  service_id: [{ required: true, message: '请输入服务ID', trigger: 'blur' }],
  expires_days: [{ required: true, message: '请选择有效期', trigger: 'change' }]
}

const getRequestStatusType = (status: number) => {
  const map: Record<number, any> = { 0: 'warning', 1: 'success', 2: 'danger' }
  return map[status] || 'info'
}

const getRequestStatusText = (status: number) => {
  const map: Record<number, string> = { 0: '待审核', 1: '已通过', 2: '已拒绝' }
  return map[status] || '未知'
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
  await ElMessageBox.confirm('确定接受该转移请求？', '提示', { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/product-diverts/requests/${row.id}/approve` })
    ElMessage.success('已接受')
    fetchRequests()
  } catch (error) {
    console.error(error)
  }
}

const handleReject = async (row: any) => {
  await ElMessageBox.confirm('确定拒绝该转移请求？', '提示', { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/product-diverts/requests/${row.id}/reject` })
    ElMessage.success('已拒绝')
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
      ElMessage.success('转移码生成成功')
      generateDialogVisible.value = false
      fetchCodes()
    } catch (error) {
      ElMessage.error('生成失败')
    } finally {
      generateSubmitLoading.value = false
    }
  })
}

const handleRegenerateCode = async (row: any) => {
  await ElMessageBox.confirm('确定重新生成该转移码？原转移码将失效', '提示', { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/product-diverts/codes/${row.id}/regenerate` })
    ElMessage.success('重新生成成功')
    fetchCodes()
  } catch (error) {
    console.error(error)
  }
}

const handleCopyCode = (code: string) => {
  navigator.clipboard.writeText(code)
  ElMessage.success('已复制到剪贴板')
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
