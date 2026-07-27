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
import { ref } from 'vue'
import { ElMessage } from 'element-plus'

const showAddDialog = ref(false)
const showConfigDialog = ref(false)
const currentPayment = ref<any>(null)

const paymentMethods = ref([
  { id: 1, name: '支付宝', gateway: 'alipay', enabled: true },
  { id: 2, name: '微信支付', gateway: 'wechat', enabled: true },
  { id: 3, name: '虎皮椒', gateway: 'xunhupay', enabled: false },
  { id: 4, name: '易支付', gateway: 'epay', enabled: false },
  { id: 5, name: 'USDT', gateway: 'usdt', enabled: false },
  { id: 6, name: 'PayPal', gateway: 'paypal', enabled: false },
  { id: 7, name: 'Stripe', gateway: 'stripe', enabled: false },
  { id: 8, name: '余额支付', gateway: 'balance', enabled: true }
])

const configForm = ref({
  app_id: '',
  app_secret: '',
  merchant_id: '',
  callback_url: ''
})

const togglePayment = (payment: any) => {
  // TODO: 切换支付方式状态
  ElMessage.success(`${payment.name} 已${payment.enabled ? '启用' : '禁用'}`)
}

const editPayment = (payment: any) => {
  currentPayment.value = payment
  showConfigDialog.value = true
  // TODO: 加载配置
}

const deletePayment = async (payment: any) => {
  // TODO: 删除支付方式
}

const saveConfig = async () => {
  // TODO: 保存配置
  showConfigDialog.value = false
  ElMessage.success('配置已保存')
}
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
