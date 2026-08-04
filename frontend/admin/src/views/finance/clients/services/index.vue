<template>
  <div class="page-container">
    <art-card title="客户服务" shadow="never">
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="user_id" label="用户ID" width="100" />
        <el-table-column prop="username" label="用户名" width="120" />
        <el-table-column prop="product_name" label="产品" min-width="150" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="250" fixed="right">
          <template #default="{ row }">
            <el-button size="small" @click="handleSuspend(row)" :disabled="row.status !== 1">暂停</el-button>
            <el-button size="small" @click="handleTerminate(row)" :disabled="row.status === 0">终止</el-button>
            <el-button size="small" @click="handleRenew(row)">续费</el-button>
          </template>
        </el-table-column>
      </el-table>
    </art-card>

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

const loading = ref(false)
const tableData = ref([])

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
    const res = await request.get({ url: '/api/admin/client-services' })
    tableData.value = res || []
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleSuspend = async (row: any) => {
  await ElMessageBox.confirm('确定暂停该服务？', '提示', { type: 'warning' })
  try {
    await request.post({ url: `/api/admin/client-services/${row.id}/suspend` })
    ElMessage.success('暂停成功')
    fetchData()
  } catch (error) {
    console.error(error)
  }
}

const handleTerminate = async (row: any) => {
  await ElMessageBox.confirm('确定终止该服务？此操作不可逆！', '警告', { type: 'error' })
  try {
    await request.post({ url: `/api/admin/client-services/${row.id}/terminate` })
    ElMessage.success('终止成功')
    fetchData()
  } catch (error) {
    console.error(error)
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
        url: `/api/admin/client-services/${renewForm.service_id}/renew`,
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

onMounted(() => fetchData())
</script>

<style scoped lang="scss">
.page-container {
  padding: 20px;
}
</style>
