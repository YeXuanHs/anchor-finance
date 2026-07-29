<template>
  <div class="payment-settings page-container">
    <div class="art-card">
      <div class="table-header">
        <h3>支付方式配置</h3>
        <el-button type="primary" @click="showAddDialog = true">
          <el-icon><Plus /></el-icon>
          添加支付方式
        </el-button>
      </div>
      
      <el-table :data="paymentMethods" style="width: 100%">
        <el-table-column prop="name" label="支付方式" />
        <el-table-column prop="gateway" label="网关" />
        <el-table-column prop="status" label="状态">
          <template #default="{ row }">
            <el-switch v-model="row.enabled" @change="togglePayment(row)" />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="150">
          <template #default="{ row }">
            <el-button type="primary" link @click="editPayment(row)">配置</el-button>
            <el-button type="danger" link @click="deletePayment(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
    
    <!-- 配置对话框 -->
    <el-dialog v-model="showConfigDialog" :title="`${currentPayment?.name} 配置`" width="600px">
      <el-form :model="configForm" label-width="120px">
        <el-form-item label="App ID">
          <el-input v-model="configForm.app_id" />
        </el-form-item>
        <el-form-item label="App Secret">
          <el-input v-model="configForm.app_secret" type="password" show-password />
        </el-form-item>
        <el-form-item label="商户号">
          <el-input v-model="configForm.merchant_id" />
        </el-form-item>
        <el-form-item label="回调URL">
          <el-input v-model="configForm.callback_url" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showConfigDialog = false">取消</el-button>
        <el-button type="primary" @click="saveConfig">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const showAddDialog = ref(false)
const showConfigDialog = ref(false)
const currentPayment = ref<any>(null)

const paymentMethods = ref<any[]>([])

const fetchPaymentMethods = async () => {
  try {
    const { data } = await request.get('/admin/payment-gateways')
    if (data?.data) {
      paymentMethods.value = data.data
    }
  } catch {
    ElMessage.error('获取支付方式列表失败')
  }
}

const configForm = ref({
  app_id: '',
  app_secret: '',
  merchant_id: '',
  callback_url: ''
})

const togglePayment = async (payment: any) => {
  try {
    await request.put(`/admin/api/v1/system/payment/${payment.id}/toggle`, { enabled: payment.enabled })
    ElMessage.success(`${payment.name} 已${payment.enabled ? '启用' : '禁用'}`)
  } catch {
    payment.enabled = !payment.enabled
    ElMessage.error('操作失败')
  }
}

const editPayment = async (payment: any) => {
  currentPayment.value = payment
  showConfigDialog.value = true
  try {
    const { data } = await request.get(`/admin/api/v1/system/payment/${payment.id}/config`)
    if (data.data) {
      configForm.value = { ...configForm.value, ...data.data }
    }
  } catch {}
}

const deletePayment = async (payment: any) => {
  try {
    await ElMessageBox.confirm(`确认删除 ${payment.name}？`, '提示', { type: 'warning' })
    await request.delete(`/admin/api/v1/system/payment/${payment.id}`)
    paymentMethods.value = paymentMethods.value.filter((p: any) => p.id !== payment.id)
    ElMessage.success('删除成功')
  } catch {}
}

const saveConfig = async () => {
  try {
    await request.put(`/admin/api/v1/system/payment/${currentPayment.value.id}/config`, configForm.value)
    showConfigDialog.value = false
    ElMessage.success('配置已保存')
  } catch {
    ElMessage.error('保存失败')
  }
}

onMounted(() => {
  fetchPaymentMethods()
})
</script>

<style scoped lang="scss">
.payment-settings {
  .table-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    
    h3 {
      margin: 0;
      font-size: 18px;
      font-weight: 600;
    }
  }
}
</style>
