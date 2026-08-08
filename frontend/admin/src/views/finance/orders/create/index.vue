<template>
  <div class="create-order-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>创建订单</span>
          <el-button @click="$router.back()">返回</el-button>
        </div>
      </template>

      <el-form ref="formRef" :model="formData" :rules="rules" label-width="120px" size="default">
        <!-- 客户选择 -->
        <el-divider content-position="left">客户信息</el-divider>

        <el-form-item label="客户" prop="client_id">
          <el-select
            v-model="formData.client_id"
            filterable
            remote
            :remote-method="searchClients"
            :loading="clientSearching"
            placeholder="请输入客户名搜索"
            style="width: 400px"
          >
            <el-option
              v-for="client in clientOptions"
              :key="client.id"
              :label="`${client.username} (${client.email})`"
              :value="client.id"
            />
          </el-select>
        </el-form-item>

        <!-- 产品信息 -->
        <el-divider content-position="left">产品信息</el-divider>

        <el-form-item label="产品类型" prop="product_type">
          <el-radio-group v-model="formData.product_type">
            <el-radio value="existing">现有产品</el-radio>
            <el-radio value="custom">自定义产品</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item v-if="formData.product_type === 'existing'" label="产品" prop="product_id">
          <el-select v-model="formData.product_id" placeholder="请选择产品" style="width: 400px" @change="handleProductChange">
            <el-option
              v-for="product in products"
              :key="product.id"
              :label="product.name"
              :value="product.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item v-if="formData.product_type === 'custom'" label="产品名称" prop="custom_product_name">
          <el-input v-model="formData.custom_product_name" placeholder="请输入产品名称" style="width: 400px" />
        </el-form-item>

        <el-form-item label="计费周期" prop="billing_cycle">
          <el-select v-model="formData.billing_cycle" placeholder="请选择计费周期" style="width: 200px">
            <el-option label="月付" value="monthly" />
            <el-option label="季付" value="quarterly" />
            <el-option label="半年付" value="semi_annually" />
            <el-option label="年付" value="annually" />
            <el-option label="一次性" value="onetime" />
          </el-select>
        </el-form-item>

        <el-form-item label="数量" prop="quantity">
          <el-input-number v-model="formData.quantity" :min="1" :max="999" />
        </el-form-item>

        <!-- 价格信息 -->
        <el-divider content-position="left">价格信息</el-divider>

        <el-form-item label="单价" prop="unit_price">
          <el-input-number v-model="formData.unit_price" :min="0" :precision="2" :step="10" />
          <span class="form-tip">元</span>
        </el-form-item>

        <el-form-item label="优惠金额" prop="discount">
          <el-input-number v-model="formData.discount" :min="0" :precision="2" :step="10" />
          <span class="form-tip">元</span>
        </el-form-item>

        <el-form-item label="订单金额">
          <span class="order-amount">¥{{ formatMoney(totalAmount) }}</span>
        </el-form-item>

        <!-- 支付信息 -->
        <el-divider content-position="left">支付信息</el-divider>

        <el-form-item label="支付方式" prop="payment_method">
          <el-select v-model="formData.payment_method" placeholder="请选择支付方式" style="width: 200px">
            <el-option label="余额支付" value="balance" />
            <el-option label="线下支付" value="offline" />
            <el-option label="免费" value="free" />
          </el-select>
        </el-form-item>

        <el-form-item label="备注" prop="notes">
          <el-input v-model="formData.notes" type="textarea" :rows="3" placeholder="请输入备注" style="width: 400px" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">创建订单</el-button>
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
import request from '@/utils/http'

const router = useRouter()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const clientSearching = ref(false)
const clientOptions = ref<any[]>([])
const products = ref<any[]>([])

// 表单数据
const formData = reactive({
  client_id: null as number | null,
  product_type: 'existing',
  product_id: null as number | null,
  custom_product_name: '',
  billing_cycle: 'monthly',
  quantity: 1,
  unit_price: 0,
  discount: 0,
  payment_method: 'balance',
  notes: ''
})

// 表单验证规则
const rules: FormRules = {
  client_id: [{ required: true, message: '请选择客户', trigger: 'change' }],
  billing_cycle: [{ required: true, message: '请选择计费周期', trigger: 'change' }],
  quantity: [{ required: true, message: '请输入数量', trigger: 'blur' }],
  unit_price: [{ required: true, message: '请输入单价', trigger: 'blur' }],
  payment_method: [{ required: true, message: '请选择支付方式', trigger: 'change' }]
}

// 订单总金额
const totalAmount = computed(() => {
  return (formData.unit_price * formData.quantity) - formData.discount
})

// 格式化金额
const formatMoney = (amount: number) => {
  if (!amount) return '0.00'
  return Number(amount).toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 })
}

// 搜索客户
const searchClients = async (query: string) => {
  if (!query) return
  clientSearching.value = true
  try {
    const data = await request.get({ url: '/api/admin/clients', params: { keyword: query, page_size: 20 } })
    clientOptions.value = data?.list || []
  } catch (error) {
    console.error('搜索客户失败:', error)
  } finally {
    clientSearching.value = false
  }
}

// 获取产品列表
const fetchProducts = async () => {
  try {
    const data = await request.get({ url: '/api/admin/products' })
    products.value = data || []
  } catch (error) {
    console.error('获取产品列表失败:', error)
  }
}

// 产品变化
const handleProductChange = (productId: number) => {
  const product = products.value.find((p) => p.id === productId)
  if (product) {
    formData.unit_price = product.price || 0
  }
}

// 提交表单
const handleSubmit = async () => {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
    submitting.value = true
    await request.post({ url: '/api/admin/orders', data: { ...formData, amount: totalAmount.value } })
    ElMessage.success('订单创建成功')
    router.push('/order-list')
  } catch (error) {
    console.error('创建订单失败:', error)
  } finally {
    submitting.value = false
  }
}

// 重置表单
const handleReset = () => {
  formRef.value?.resetFields()
}

onMounted(() => {
  fetchProducts()
})
</script>

<style scoped lang="scss">
.create-order-page {
  padding: 16px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.form-tip {
  margin-left: 10px;
  font-size: 12px;
  color: #86909C;
}

.order-amount {
  font-size: 24px;
  font-weight: 600;
  color: #F59E0B;
}

:deep(.el-divider__text) {
  font-weight: 600;
  color: #1D2129;
}
</style>
