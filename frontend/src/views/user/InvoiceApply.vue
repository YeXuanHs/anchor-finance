<template>
  <div class="invoice-apply-page">
    <div class="page-header">
      <h1 class="page-title">申请发票</h1>
      <el-button @click="$router.back()">返回列表</el-button>
    </div>

    <el-card shadow="never" class="form-card">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        label-position="right"
        class="apply-form"
      >
        <!-- 发票类型 -->
        <el-divider content-position="left">发票类型</el-divider>
        <el-form-item label="发票类型" prop="invoiceType">
          <el-radio-group v-model="form.invoiceType">
            <el-radio value="normal">增值税普通发票</el-radio>
            <el-radio value="special">增值税专用发票</el-radio>
          </el-radio-group>
          <div class="form-tip">增值税专用发票可用于抵扣进项税额</div>
        </el-form-item>

        <!-- 关联订单 -->
        <el-divider content-position="left">关联订单</el-divider>
        <el-form-item label="关联订单" prop="orderIds">
          <el-select
            v-model="form.orderIds"
            multiple
            filterable
            placeholder="请选择需要开票的订单"
            style="width: 100%;"
          >
            <el-option
              v-for="order in availableOrders"
              :key="order.id"
              :label="`${order.orderNo} - ${order.description} (¥${order.amount})`"
              :value="order.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="开票金额">
          <span class="calculated-amount">¥{{ calculatedAmount }}</span>
          <span class="form-tip" style="margin-left: 12px;">根据所选订单自动计算</span>
        </el-form-item>

        <!-- 发票信息 -->
        <el-divider content-position="left">发票信息</el-divider>
        <el-form-item label="发票抬头" prop="title">
          <el-input v-model="form.title" placeholder="请输入发票抬头" maxlength="100" show-word-limit />
        </el-form-item>
        <el-form-item label="纳税人识别号" prop="taxNo">
          <el-input v-model="form.taxNo" placeholder="请输入纳税人识别号" maxlength="20" />
        </el-form-item>
        <template v-if="form.invoiceType === 'special'">
          <el-form-item label="开户银行" prop="bankName">
            <el-input v-model="form.bankName" placeholder="请输入开户银行" />
          </el-form-item>
          <el-form-item label="银行账号" prop="bankAccount">
            <el-input v-model="form.bankAccount" placeholder="请输入银行账号" />
          </el-form-item>
          <el-form-item label="注册地址" prop="registerAddress">
            <el-input v-model="form.registerAddress" placeholder="请输入注册地址" />
          </el-form-item>
          <el-form-item label="注册电话" prop="registerPhone">
            <el-input v-model="form.registerPhone" placeholder="请输入注册电话" />
          </el-form-item>
        </template>

        <!-- 邮寄信息 -->
        <el-divider content-position="left">邮寄信息</el-divider>
        <el-form-item label="收件人" prop="receiverName">
          <el-input v-model="form.receiverName" placeholder="请输入收件人姓名" />
        </el-form-item>
        <el-form-item label="联系电话" prop="receiverPhone">
          <el-input v-model="form.receiverPhone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="邮寄地址" prop="receiverAddress">
          <el-input v-model="form.receiverAddress" type="textarea" :rows="2" placeholder="请输入详细邮寄地址" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3" placeholder="请输入备注信息（选填）" maxlength="200" show-word-limit />
        </el-form-item>

        <!-- 提交 -->
        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">提交申请</el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const loadingOrders = ref(false)

interface OrderOption {
  id: number
  orderNo: string
  description: string
  amount: string
}

const availableOrders = ref<OrderOption[]>([])

onMounted(async () => {
  loadingOrders.value = true
  try {
    const { data } = await request.get('/api/v1/orders')
    availableOrders.value = data.data?.list || data.list || data.data || []
  } catch (e) { console.error(e) } finally { loadingOrders.value = false }
})

const form = reactive({
  invoiceType: 'normal' as 'normal' | 'special',
  orderIds: [] as number[],
  title: '',
  taxNo: '',
  bankName: '',
  bankAccount: '',
  registerAddress: '',
  registerPhone: '',
  receiverName: '',
  receiverPhone: '',
  receiverAddress: '',
  remark: ''
})

const rules: FormRules = {
  invoiceType: [{ required: true, message: '请选择发票类型', trigger: 'change' }],
  orderIds: [{ required: true, message: '请选择关联订单', trigger: 'change' }],
  title: [{ required: true, message: '请输入发票抬头', trigger: 'blur' }],
  taxNo: [{ required: true, message: '请输入纳税人识别号', trigger: 'blur' }],
  bankName: [{ required: true, message: '请输入开户银行', trigger: 'blur' }],
  bankAccount: [{ required: true, message: '请输入银行账号', trigger: 'blur' }],
  registerAddress: [{ required: true, message: '请输入注册地址', trigger: 'blur' }],
  registerPhone: [{ required: true, message: '请输入注册电话', trigger: 'blur' }],
  receiverName: [{ required: true, message: '请输入收件人姓名', trigger: 'blur' }],
  receiverPhone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }],
  receiverAddress: [{ required: true, message: '请输入邮寄地址', trigger: 'blur' }]
}

const calculatedAmount = computed(() => {
  const total = form.orderIds.reduce((sum, id) => {
    const order = availableOrders.value.find(o => o.id === id)
    return sum + (order ? parseFloat(order.amount) : 0)
  }, 0)
  return total.toFixed(2)
})

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        await request.post('/api/v1/invoices/apply', {
          invoiceType: form.invoiceType,
          orderIds: form.orderIds,
          title: form.title,
          taxNo: form.taxNo,
          bankName: form.bankName,
          bankAccount: form.bankAccount,
          registerAddress: form.registerAddress,
          registerPhone: form.registerPhone,
          receiverName: form.receiverName,
          receiverPhone: form.receiverPhone,
          receiverAddress: form.receiverAddress,
          remark: form.remark
        })
        ElMessage.success('发票申请已提交，请等待审核')
        router.push('/user/invoice/list')
      } catch (e: any) { ElMessage.error(e?.message || '提交失败，请重试') } finally { submitting.value = false }
    }
  })
}

function handleReset() {
  if (formRef.value) {
    formRef.value.resetFields()
  }
  form.remark = ''
}
</script>

<style scoped>
.invoice-apply-page {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-size: 20px;
  font-weight: 700;
  color: #303133;
  margin: 0;
}

.form-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
}

.form-card :deep(.el-card__body) {
  padding: 20px 40px;
}

.apply-form {
  max-width: 700px;
}

.apply-form :deep(.el-divider__text) {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
}

.apply-form :deep(.el-divider) {
  margin: 24px 0 20px;
}

.form-tip {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}

.calculated-amount {
  font-size: 24px;
  font-weight: 700;
  color: #0056FF;
}

@media (max-width: 768px) {
  .page-header { flex-direction: column; gap: 12px; align-items: flex-start; }
  .form-card :deep(.el-card__body) { padding: 16px; }
  .apply-form :deep(.el-form-item__label) { width: 100px !important; }
  .apply-form :deep(.el-form-item__content) { margin-left: 100px !important; }
}
</style>
