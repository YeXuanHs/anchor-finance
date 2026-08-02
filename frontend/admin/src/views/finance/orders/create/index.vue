<template>
  <div class="create-order-page">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>创建订单</span>
          <el-button @click="handleBack">
            <el-icon><Back /></el-icon>
            返回
          </el-button>
        </div>
      </template>

      <el-steps :active="currentStep" finish-status="success" align-center class="steps-bar">
        <el-step title="选择客户" />
        <el-step title="选择产品" />
        <el-step title="配置参数" />
        <el-step title="确认创建" />
      </el-steps>

      <!-- 步骤1: 选择客户 -->
      <div v-show="currentStep === 0" class="step-content">
        <el-form :model="orderForm" ref="step1FormRef" :rules="step1Rules" label-width="120px">
          <el-form-item label="选择客户" prop="client_id">
            <el-select
              v-model="orderForm.client_id"
              filterable
              remote
              :remote-method="searchClients"
              :loading="clientSearching"
              placeholder="请输入用户名或邮箱搜索"
              style="width: 100%"
              @change="handleClientChange"
            >
              <el-option
                v-for="item in clientOptions"
                :key="item.id"
                :label="item.username + ' (' + item.email + ')'"
                :value="item.id"
              >
                <div class="client-option">
                  <span class="client-name">{{ item.username }}</span>
                  <span class="client-email">{{ item.email }}</span>
                </div>
              </el-option>
            </el-select>
          </el-form-item>
        </el-form>

        <!-- 选中的客户信息 -->
        <el-card v-if="selectedClient" shadow="never" class="selected-info">
          <template #header>
            <span>客户信息</span>
          </template>
          <el-descriptions :column="2" size="small">
            <el-descriptions-item label="用户名">{{ selectedClient.username }}</el-descriptions-item>
            <el-descriptions-item label="邮箱">{{ selectedClient.email }}</el-descriptions-item>
            <el-descriptions-item label="手机号">{{ selectedClient.phone || '-' }}</el-descriptions-item>
            <el-descriptions-item label="余额">¥{{ formatAmount(selectedClient.balance) }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </div>

      <!-- 步骤2: 选择产品 -->
      <div v-show="currentStep === 1" class="step-content">
        <el-form :model="orderForm" ref="step2FormRef" :rules="step2Rules" label-width="120px">
          <el-form-item label="产品分组" prop="group_id">
            <el-select v-model="orderForm.group_id" placeholder="请选择分组" @change="handleGroupChange">
              <el-option v-for="group in productGroups" :key="group.id" :label="group.name" :value="group.id" />
            </el-select>
          </el-form-item>
          <el-form-item label="选择产品" prop="product_id">
            <el-select
              v-model="orderForm.product_id"
              placeholder="请选择产品"
              :disabled="!orderForm.group_id"
              @change="handleProductChange"
            >
              <el-option
                v-for="product in filteredProducts"
                :key="product.id"
                :label="product.name"
                :value="product.id"
              />
            </el-select>
          </el-form-item>
        </el-form>

        <!-- 产品信息 -->
        <el-card v-if="selectedProduct" shadow="never" class="selected-info">
          <template #header>
            <span>产品信息</span>
          </template>
          <el-descriptions :column="2" size="small">
            <el-descriptions-item label="产品名称">{{ selectedProduct.name }}</el-descriptions-item>
            <el-descriptions-item label="产品类型">{{ selectedProduct.type }}</el-descriptions-item>
            <el-descriptions-item label="计费周期">{{ selectedProduct.billing_cycle }}</el-descriptions-item>
            <el-descriptions-item label="价格">¥{{ formatAmount(selectedProduct.price) }} / {{ selectedProduct.billing_cycle }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </div>

      <!-- 步骤3: 配置参数 -->
      <div v-show="currentStep === 2" class="step-content">
        <el-form :model="orderForm" ref="step3FormRef" :rules="step3Rules" label-width="120px">
          <el-form-item label="周期" prop="billing_cycle">
            <el-radio-group v-model="orderForm.billing_cycle">
              <el-radio-button
                v-for="cycle in billingCycles"
                :key="cycle.value"
                :value="cycle.value"
              >
                {{ cycle.label }} ¥{{ formatAmount(cycle.price) }}
              </el-radio-button>
            </el-radio-group>
          </el-form-item>
          <el-form-item label="数量" prop="quantity">
            <el-input-number v-model="orderForm.quantity" :min="1" :max="999" />
          </el-form-item>
          <el-form-item label="自定义价格">
            <el-input-number v-model="orderForm.custom_price" :min="0" :precision="2" :step="10" />
            <span class="form-tip">留空则使用产品默认价格</span>
          </el-form-item>
          <el-form-item label="折扣">
            <el-input-number v-model="orderForm.discount" :min="0" :max="100" :precision="0" />
            <span class="form-tip">百分比，如 10 表示打9折</span>
          </el-form-item>

          <!-- 产品配置选项 -->
          <el-divider content-position="left">产品配置</el-divider>
          <el-form-item label="域名">
            <el-input v-model="orderForm.domain" placeholder="如: example.com" />
          </el-form-item>
          <el-form-item label="IP地址">
            <el-input v-model="orderForm.ip" placeholder="如: 1.2.3.4" />
          </el-form-item>
          <el-form-item label="配置参数">
            <el-input v-model="orderForm.config" type="textarea" :rows="3" placeholder="JSON格式配置参数" />
          </el-form-item>
          <el-form-item label="备注">
            <el-input v-model="orderForm.remark" type="textarea" :rows="2" placeholder="订单备注" />
          </el-form-item>
        </el-form>

        <!-- 价格预览 -->
        <el-card shadow="never" class="selected-info">
          <template #header>
            <span>价格预览</span>
          </template>
          <el-descriptions :column="1" size="small">
            <el-descriptions-item label="单价">¥{{ formatAmount(currentPrice) }}</el-descriptions-item>
            <el-descriptions-item label="数量">{{ orderForm.quantity }}</el-descriptions-item>
            <el-descriptions-item label="小计">¥{{ formatAmount(currentPrice * orderForm.quantity) }}</el-descriptions-item>
            <el-descriptions-item label="折扣">{{ orderForm.discount }}%</el-descriptions-item>
            <el-descriptions-item label="折后金额">
              <span class="amount-highlight">¥{{ formatAmount(finalPrice) }}</span>
            </el-descriptions-item>
          </el-descriptions>
        </el-card>
      </div>

      <!-- 步骤4: 确认创建 -->
      <div v-show="currentStep === 3" class="step-content">
        <el-card shadow="never">
          <template #header>
            <span>订单确认</span>
          </template>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="客户">
              {{ selectedClient?.username }} ({{ selectedClient?.email }})
            </el-descriptions-item>
            <el-descriptions-item label="产品">{{ selectedProduct?.name }}</el-descriptions-item>
            <el-descriptions-item label="周期">{{ orderForm.billing_cycle }}</el-descriptions-item>
            <el-descriptions-item label="数量">{{ orderForm.quantity }}</el-descriptions-item>
            <el-descriptions-item label="单价">¥{{ formatAmount(currentPrice) }}</el-descriptions-item>
            <el-descriptions-item label="折扣">{{ orderForm.discount }}%</el-descriptions-item>
            <el-descriptions-item label="订单金额">
              <span class="amount-highlight" style="font-size: 18px;">¥{{ formatAmount(finalPrice) }}</span>
            </el-descriptions-item>
            <el-descriptions-item label="域名">{{ orderForm.domain || '-' }}</el-descriptions-item>
            <el-descriptions-item label="备注" :span="2">{{ orderForm.remark || '无' }}</el-descriptions-item>
          </el-descriptions>
        </el-card>
      </div>

      <!-- 底部按钮 -->
      <div class="step-actions">
        <el-button v-if="currentStep > 0" @click="currentStep--">上一步</el-button>
        <el-button v-if="currentStep < 3" type="primary" @click="handleNext">下一步</el-button>
        <el-button v-if="currentStep === 3" type="primary" @click="handleSubmit" :loading="submitLoading">
          确认创建
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Back } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/http'

const route = useRoute()
const router = useRouter()

const currentStep = ref(0)
const submitLoading = ref(false)
const clientSearching = ref(false)

const step1FormRef = ref<FormInstance>()
const step2FormRef = ref<FormInstance>()
const step3FormRef = ref<FormInstance>()

// 客户搜索
const clientOptions = ref<any[]>([])
const selectedClient = ref<any>(null)

// 产品
const productGroups = ref<any[]>([])
const products = ref<any[]>([])
const selectedProduct = ref<any>(null)

const billingCycles = ref<any[]>([])

const orderForm = reactive({
  client_id: undefined as number | undefined,
  group_id: undefined as number | undefined,
  product_id: undefined as number | undefined,
  billing_cycle: 'monthly',
  quantity: 1,
  custom_price: undefined as number | undefined,
  discount: 0,
  domain: '',
  ip: '',
  config: '',
  remark: ''
})

const step1Rules: FormRules = {
  client_id: [{ required: true, message: '请选择客户', trigger: 'change' }]
}

const step2Rules: FormRules = {
  group_id: [{ required: true, message: '请选择产品分组', trigger: 'change' }],
  product_id: [{ required: true, message: '请选择产品', trigger: 'change' }]
}

const step3Rules: FormRules = {
  billing_cycle: [{ required: true, message: '请选择计费周期', trigger: 'change' }],
  quantity: [{ required: true, message: '请输入数量', trigger: 'blur' }]
}

const filteredProducts = computed(() => {
  if (!orderForm.group_id) return products.value
  return products.value.filter(p => p.group_id === orderForm.group_id)
})

const currentPrice = computed(() => {
  if (orderForm.custom_price !== undefined && orderForm.custom_price > 0) {
    return orderForm.custom_price
  }
  const cyclePrice = billingCycles.value.find(c => c.value === orderForm.billing_cycle)
  return cyclePrice?.price || selectedProduct.value?.price || 0
})

const finalPrice = computed(() => {
  const subtotal = currentPrice.value * orderForm.quantity
  const discountAmount = subtotal * (orderForm.discount / 100)
  return subtotal - discountAmount
})

const formatAmount = (amount: number | undefined) => {
  return amount?.toLocaleString('zh-CN', { minimumFractionDigits: 2, maximumFractionDigits: 2 }) || '0.00'
}

const searchClients = async (query: string) => {
  if (!query) {
    clientOptions.value = []
    return
  }
  clientSearching.value = true
  try {
    const data = await request.get({
      url: '/api/admin/users',
      params: { keyword: query, page_size: 20 }
    })
    clientOptions.value = data.list || []
  } catch (error) {
    console.error('搜索客户失败:', error)
  } finally {
    clientSearching.value = false
  }
}

const handleClientChange = async (clientId: number) => {
  try {
    selectedClient.value = await request.get({ url: `/api/admin/users/${clientId}` })
  } catch (error) {
    console.error('获取客户信息失败:', error)
  }
}

const handleGroupChange = () => {
  orderForm.product_id = undefined
  selectedProduct.value = null
}

const handleProductChange = async (productId: number) => {
  try {
    const data = await request.get({ url: `/api/admin/products/${productId}` })
    selectedProduct.value = data
    billingCycles.value = data.billing_cycles || []
    if (billingCycles.value.length > 0) {
      orderForm.billing_cycle = billingCycles.value[0].value
    }
  } catch (error) {
    console.error('获取产品信息失败:', error)
  }
}

const fetchProductGroups = async () => {
  try {
    productGroups.value = await request.get({ url: '/api/admin/product-groups' }) || []
  } catch (error) {
    console.error('获取产品分组失败:', error)
  }
}

const fetchProducts = async () => {
  try {
    const data = await request.get({ url: '/api/admin/products', params: { page_size: 1000 } })
    products.value = data.list || []
  } catch (error) {
    console.error('获取产品列表失败:', error)
  }
}

const handleNext = async () => {
  const formRefs: Record<number, FormInstance | undefined> = {
    0: step1FormRef.value,
    1: step2FormRef.value,
    2: step3FormRef.value
  }

  const formRef = formRefs[currentStep.value]
  if (formRef) {
    const valid = await formRef.validate().catch(() => false)
    if (!valid) return
  }

  currentStep.value++
}

const handleSubmit = async () => {
  submitLoading.value = true
  try {
    await request.post({
      url: '/api/admin/orders',
      params: orderForm
    })
    ElMessage.success('订单创建成功')
    router.push('/finance/orders/list')
  } catch (error) {
    ElMessage.error('创建失败')
  } finally {
    submitLoading.value = false
  }
}

const handleBack = () => {
  router.back()
}

onMounted(() => {
  fetchProductGroups()
  fetchProducts()

  // 如果URL带了client_id参数，自动选中
  const clientId = route.query.client_id
  if (clientId) {
    orderForm.client_id = Number(clientId)
    handleClientChange(Number(clientId))
  }
})
</script>

<style scoped lang="scss">
.create-order-page {
  padding: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.steps-bar {
  margin: 20px 0 32px;
}

.step-content {
  min-height: 300px;
  padding: 20px 0;
}

.selected-info {
  margin-top: 20px;
}

.client-option {
  display: flex;
  justify-content: space-between;

  .client-email {
    color: var(--el-text-color-secondary);
    font-size: 12px;
  }
}

.form-tip {
  margin-left: 12px;
  color: var(--el-text-color-secondary);
  font-size: 13px;
}

.amount-highlight {
  font-weight: 600;
  color: var(--el-color-primary);
}

.step-actions {
  display: flex;
  justify-content: center;
  gap: 12px;
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid var(--el-border-color-lighter);
}
</style>