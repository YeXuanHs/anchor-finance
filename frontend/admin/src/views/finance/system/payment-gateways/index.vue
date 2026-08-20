<template>
  <div class="payment-gateways-page">
    <!-- 顶部操作栏 -->
    <div class="page-header">
      <div class="header-left">
        <h2>{{ $t('paymentGateways.title') }}</h2>
        <span class="subtitle">{{ $t('paymentGateways.subtitle') }}</span>
      </div>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        {{ $t('paymentGateways.addPayment') }}
      </el-button>
    </div>

    <!-- 支付方式列表 -->
    <el-card class="table-card">
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column :label="$t('paymentGateways.icon')" width="60">
          <template #default="{ row }">
            <img :src="getIconUrl(row.code)" class="payment-icon" />
          </template>
        </el-table-column>
        <el-table-column prop="title" :label="$t('paymentGateways.displayName')" min-width="150">
          <template #default="{ row }">
            <span class="payment-title">{{ row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" :label="$t('paymentGateways.identifier')" width="150" />
        <el-table-column :label="$t('paymentGateways.gatewayType')" width="180">
          <template #default="{ row }">
            <el-tag size="small" type="info">{{ getGatewayLabel(row.gateway) }}</el-tag>
            <el-tag size="small" class="ml-1">{{ getCodeLabel(row.code) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('paymentGateways.fee')" width="100">
          <template #default="{ row }">
            {{ row.fee_rate > 0 ? (row.fee_rate * 100).toFixed(2) + '%' : $t('common.none') }}
          </template>
        </el-table-column>
        <el-table-column :label="$t('paymentGateways.amountLimit')" width="150">
          <template #default="{ row }">
            <template v-if="row.min_amount > 0 || row.max_amount > 0">
              {{ row.min_amount || 0 }} - {{ row.max_amount || $t('common.none') }}
            </template>
            <span v-else class="text-muted">{{ $t('common.none') }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" :label="$t('common.sortOrder')" width="80" />
        <el-table-column :label="$t('common.status')" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_enabled"
              @change="handleToggleStatus(row)"
              :active-text="''"
              :inactive-text="''"
            />
          </template>
        </el-table-column>
        <el-table-column :label="$t('common.action')" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">{{ $t('common.edit') }}</el-button>
            <el-button type="info" link size="small" @click="handleTest(row)">{{ $t('paymentGateways.test') }}</el-button>
            <el-popconfirm :title="$t('paymentGateways.confirmDelete')" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link size="small">{{ $t('common.delete') }}</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[20, 50, 100]"
          layout="total, sizes, prev, pager, next"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? $t('paymentGateways.editPayment') : $t('paymentGateways.addPayment')"
      width="680px"
      :close-on-click-modal="false"
    >
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <!-- 接口选择 -->
        <el-form-item :label="$t('paymentGateways.gateway')" prop="gateway">
          <el-select v-model="formData.gateway" :placeholder="$t('paymentGateways.selectGateway')" @change="handleGatewayChange" :disabled="isEdit">
            <el-option
              v-for="(label, value) in supportedInfo.gateways"
              :key="value"
              :label="label"
              :value="value"
            />
          </el-select>
          <div class="form-hint">{{ $t('paymentGateways.gatewayHint') }}</div>
        </el-form-item>

        <!-- 支付类型选择 -->
        <el-form-item :label="$t('paymentGateways.paymentType')" prop="code">
          <el-select v-model="formData.code" :placeholder="$t('paymentGateways.selectType')" @change="handleCodeChange" :disabled="isEdit">
            <el-option
              v-for="method in availableCodes"
              :key="method.value"
              :label="method.label"
              :value="method.value"
            >
              <div style="display: flex; align-items: center; gap: 8px;">
                <img :src="getIconUrl(method.value)" style="width: 20px; height: 20px;" />
                <span>{{ method.label }}</span>
              </div>
            </el-option>
          </el-select>
          <div class="form-hint">{{ $t('paymentGateways.typeHint') }}</div>
        </el-form-item>

        <!-- 标识名 -->
        <el-form-item :label="$t('paymentGateways.uniqueId')" prop="name">
          <el-input v-model="formData.name" :placeholder="$t('paymentGateways.idPlaceholder')" :disabled="isEdit" />
          <div class="form-hint">{{ $t('paymentGateways.idHint') }}</div>
        </el-form-item>

        <!-- 显示名称 -->
        <el-form-item :label="$t('paymentGateways.displayName')" prop="title">
          <el-input v-model="formData.title" :placeholder="$t('paymentGateways.namePlaceholder')" />
          <div class="form-hint">{{ $t('paymentGateways.nameHint') }}</div>
        </el-form-item>

        <!-- 接口配置 -->
        <el-divider content-position="left">{{ $t('paymentGateways.gatewayConfig') }}</el-divider>

        <!-- 易支付配置 -->
        <template v-if="formData.gateway === 'epay'">
          <el-form-item :label="$t('paymentGateways.merchantId')" prop="config.pid">
            <el-input v-model="configForm.pid" :placeholder="$t('paymentGateways.merchantIdPlaceholder')" />
          </el-form-item>
          <el-form-item :label="$t('paymentGateways.merchantKey')" prop="config.key">
            <el-input v-model="configForm.key" :placeholder="$t('paymentGateways.merchantKeyPlaceholder')" show-password />
          </el-form-item>
          <el-form-item :label="$t('paymentGateways.apiUrl')" prop="config.api_url">
            <el-input v-model="configForm.api_url" placeholder="https://pay.example.com" />
          </el-form-item>
        </template>

        <!-- 迅虎支付配置 -->
        <template v-if="formData.gateway === 'xunhupay'">
          <el-form-item :label="$t('paymentGateways.appId')" prop="config.api_key">
            <el-input v-model="configForm.api_key" :placeholder="$t('paymentGateways.xunhupayAppIdPlaceholder')" />
          </el-form-item>
          <el-form-item :label="$t('paymentGateways.appSecret')" prop="config.secret_key">
            <el-input v-model="configForm.secret_key" :placeholder="$t('paymentGateways.xunhupaySecretPlaceholder')" show-password />
          </el-form-item>
          <el-form-item :label="$t('paymentGateways.apiUrl')" prop="config.gateway">
            <el-input v-model="configForm.gateway" placeholder="https://api.xunhupay.com" />
          </el-form-item>
        </template>

        <!-- 支付宝官方配置 -->
        <template v-if="formData.gateway === 'alipay'">
          <el-form-item :label="$t('paymentGateways.appId')" prop="config.app_id">
            <el-input v-model="configForm.app_id" :placeholder="$t('paymentGateways.alipayAppIdPlaceholder')" />
          </el-form-item>
          <el-form-item :label="$t('paymentGateways.privateKey')" prop="config.private_key">
            <el-input v-model="configForm.private_key" type="textarea" :rows="4" :placeholder="$t('paymentGateways.privateKeyPlaceholder')" />
          </el-form-item>
          <el-form-item :label="$t('paymentGateways.alipayPublicKey')" prop="config.alipay_public_key">
            <el-input v-model="configForm.alipay_public_key" type="textarea" :rows="4" :placeholder="$t('paymentGateways.alipayPublicKeyPlaceholder')" />
          </el-form-item>
          <el-form-item :label="$t('paymentGateways.paymentProduct')">
            <el-checkbox v-model="configForm.product_pc">{{ $t('paymentGateways.pcPayment') }}</el-checkbox>
            <el-checkbox v-model="configForm.product_wap">{{ $t('paymentGateways.wapPayment') }}</el-checkbox>
            <el-checkbox v-model="configForm.product_qr">{{ $t('paymentGateways.qrPayment') }}</el-checkbox>
          </el-form-item>
        </template>

        <!-- 微信支付官方配置 -->
        <template v-if="formData.gateway === 'wxpay'">
          <el-form-item :label="$t('paymentGateways.officialAccountId')" prop="config.app_id">
            <el-input v-model="configForm.app_id" :placeholder="$t('paymentGateways.wxAppIdPlaceholder')" />
          </el-form-item>
          <el-form-item :label="$t('paymentGateways.merchantNo')" prop="config.mch_id">
            <el-input v-model="configForm.mch_id" :placeholder="$t('paymentGateways.wxMerchantPlaceholder')" />
          </el-form-item>
          <el-form-item :label="$t('paymentGateways.apiKey')" prop="config.api_key">
            <el-input v-model="configForm.api_key" :placeholder="$t('paymentGateways.wxApiKeyPlaceholder')" show-password />
          </el-form-item>
          <el-form-item :label="$t('paymentGateways.appSecret')" prop="config.app_secret">
            <el-input v-model="configForm.app_secret" :placeholder="$t('paymentGateways.wxAppSecretPlaceholder')" show-password />
          </el-form-item>
          <el-form-item :label="$t('paymentGateways.paymentProduct')">
            <el-checkbox v-model="configForm.product_native">{{ $t('paymentGateways.qrPayment') }}</el-checkbox>
            <el-checkbox v-model="configForm.product_jsapi">{{ $t('paymentGateways.jsapiPayment') }}</el-checkbox>
            <el-checkbox v-model="configForm.product_wap">{{ $t('paymentGateways.h5Payment') }}</el-checkbox>
          </el-form-item>
        </template>

        <!-- 余额支付无配置 -->
        <template v-if="formData.gateway === 'balance'">
          <el-alert :title="$t('paymentGateways.balanceNoConfig')" type="info" :closable="false" show-icon />
        </template>

        <!-- 其他设置 -->
        <el-divider content-position="left">{{ $t('paymentGateways.otherSettings') }}</el-divider>

        <el-form-item :label="$t('paymentGateways.feeRate')">
          <el-input-number v-model="formData.fee_rate" :min="0" :max="1" :step="0.01" :precision="4" />
          <span class="ml-2">%</span>
          <div class="form-hint">{{ $t('paymentGateways.feeRateHint') }}</div>
        </el-form-item>

        <el-form-item :label="$t('paymentGateways.amountLimit')">
          <div class="amount-range">
            <el-input-number v-model="formData.min_amount" :min="0" :precision="2" :placeholder="$t('paymentGateways.minAmount')" />
            <span class="range-sep">{{ $t('common.to') }}</span>
            <el-input-number v-model="formData.max_amount" :min="0" :precision="2" :placeholder="$t('paymentGateways.maxAmount')" />
          </div>
          <div class="form-hint">{{ $t('paymentGateways.amountHint') }}</div>
        </el-form-item>

        <el-form-item :label="$t('common.sortOrder')">
          <el-input-number v-model="formData.sort_order" :min="0" />
          <div class="form-hint">{{ $t('paymentGateways.sortHint') }}</div>
        </el-form-item>

        <el-form-item :label="$t('paymentGateways.enabledStatus')">
          <el-switch v-model="formData.is_enabled" :active-text="$t('common.enable')" :inactive-text="$t('common.disable')" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">{{ $t('common.cancel') }}</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEdit ? $t('common.save') : $t('common.add') }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { $t } from '@/locales'
import request from '@/utils/http'

// 支付方式数据
interface PaymentGateway {
  id: number
  name: string
  title: string
  gateway: string
  code: string
  config: string
  fee_rate: number
  min_amount: number
  max_amount: number
  sort_order: number
  is_enabled: boolean
}

// 支持的接口和类型
interface SupportedInfo {
  gateways: Record<string, string>
  codes: Record<string, string>
}

const loading = ref(false)
const submitting = ref(false)
const dialogVisible = ref(false)
const isEdit = ref(false)
const editId = ref(0)
const tableData = ref<PaymentGateway[]>([])
const supportedInfo = ref<SupportedInfo>({ gateways: {}, codes: {} })

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

// 表单数据
const formData = reactive({
  name: '',
  title: '',
  gateway: '',
  code: '',
  config: '',
  fee_rate: 0,
  min_amount: 0,
  max_amount: 0,
  sort_order: 0,
  is_enabled: true
})

// 接口配置表单
const configForm = reactive<Record<string, any>>({})

// 表单验证规则
const formRules = {
  name: [{ required: true, message: $t('paymentGateways.inputId'), trigger: 'blur' }],
  title: [{ required: true, message: $t('paymentGateways.inputName'), trigger: 'blur' }],
  gateway: [{ required: true, message: $t('paymentGateways.selectGateway'), trigger: 'change' }],
  code: [{ required: true, message: $t('paymentGateways.selectType'), trigger: 'change' }]
}

// 获取支付类型图标
const getIconUrl = (code: string) => {
  const icons: Record<string, string> = {
    alipay: '/assets/payment/alipay.png',
    wechat: '/assets/payment/wechat.png',
    qqpay: '/assets/payment/qqpay.png',
    usdt: '/assets/payment/usdt.png',
    bank: '/assets/payment/bank.png',
    balance: '/assets/payment/balance.png'
  }
  return icons[code] || '/assets/payment/default.png'
}

// 获取接口显示名
const getGatewayLabel = (gateway: string) => {
  return supportedInfo.value.gateways[gateway] || gateway
}

// 获取支付类型显示名
const getCodeLabel = (code: string) => {
  return supportedInfo.value.codes[code] || code
}

// 当前接口支持的支付类型
const availableCodes = computed(() => {
  if (!formData.gateway) return []
  // 根据接口类型返回支持的支付类型（参考zjmf）
  const gatewayCodes: Record<string, string[]> = {
    epay: ['alipay', 'wechat', 'qqpay', 'usdt', 'bank'],  // 易支付支持所有类型
    xunhupay: ['alipay', 'wechat'],  // 迅虎支付只支持支付宝和微信
    alipay: ['alipay'],
    wxpay: ['wechat'],
    balance: ['balance']
  }
  const codes = gatewayCodes[formData.gateway] || []
  return codes.map(code => ({
    value: code,
    label: supportedInfo.value.codes[code] || code
  }))
})

// 接口配置默认值
const defaultConfigs: Record<string, Record<string, any>> = {
  epay: { pid: '', key: '', api_url: '' },
  xunhupay: { api_key: '', secret_key: '', gateway: 'https://api.xunhupay.com' },
  alipay: { app_id: '', private_key: '', alipay_public_key: '', product_pc: true, product_wap: true, product_qr: false },
  wxpay: { app_id: '', mch_id: '', api_key: '', app_secret: '', product_native: true, product_jsapi: false, product_wap: false },
  balance: {}
}

// 接口切换
const handleGatewayChange = (gateway: string) => {
  formData.code = ''
  Object.keys(configForm).forEach(key => delete configForm[key])
  Object.assign(configForm, defaultConfigs[gateway] || {})
  updateName()
}

// 支付类型切换
const handleCodeChange = (code: string) => {
  updateName()
}

// 自动生成名称
const updateName = () => {
  if (formData.gateway && formData.code) {
    const gatewayName = formData.gateway.charAt(0).toUpperCase() + formData.gateway.slice(1)
    const codeName = formData.code.charAt(0).toUpperCase() + formData.code.slice(1)
    formData.name = `${gatewayName}${codeName}`
    formData.title = `${getCodeLabel(formData.code)}-${getGatewayLabel(formData.gateway)}`
  }
}

// 获取数据
const fetchData = async () => {
  loading.value = true
  try {
    const [listRes, infoRes] = await Promise.all([
      request.get({
        url: '/api/admin/payment-gateways',
        params: {
          page: pagination.page,
          page_size: pagination.pageSize
        }
      }),
      request.get({ url: '/api/admin/payment-gateways/supported' })
    ])

    if (listRes.data) {
      tableData.value = listRes.data.list || listRes.data.items || []
      pagination.total = listRes.data.total || 0
    }
    if (infoRes.data) {
      supportedInfo.value = infoRes.data
    }
  } catch (error) {
    console.error('获取支付方式列表失败:', error)
  } finally {
    loading.value = false
  }
}

// 打开新增对话框
const handleAdd = () => {
  isEdit.value = false
  editId.value = 0
  Object.assign(formData, {
    name: '',
    title: '',
    gateway: '',
    code: '',
    config: '',
    fee_rate: 0,
    min_amount: 0,
    max_amount: 0,
    sort_order: 0,
    is_enabled: true
  })
  Object.keys(configForm).forEach(key => delete configForm[key])
  dialogVisible.value = true
}

// 打开编辑对话框
const handleEdit = (row: any) => {
  isEdit.value = true
  editId.value = row.id
  Object.assign(formData, {
    name: row.name,
    title: row.title,
    gateway: row.gateway,
    code: row.code,
    config: row.config,
    fee_rate: row.fee_rate,
    min_amount: row.min_amount,
    max_amount: row.max_amount,
    sort_order: row.sort_order,
    is_enabled: row.is_enabled
  })
  // 解析配置
  Object.keys(configForm).forEach(key => delete configForm[key])
  if (row.config) {
    try {
      const config = JSON.parse(row.config)
      Object.assign(configForm, config)
    } catch (e) {
      console.error('解析配置失败:', e)
    }
  }
  dialogVisible.value = true
}

// 提交表单
const handleSubmit = async () => {
  submitting.value = true
  try {
    const data = {
      ...formData,
      config: JSON.stringify(configForm)
    }

    if (isEdit.value) {
      await request.put({
        url: `/api/admin/payment-gateways/${editId.value}`,
        data
      })
      ElMessage.success($t('common.updateSuccess'))
    } else {
      await request.post({
        url: '/api/admin/payment-gateways',
        data
      })
      ElMessage.success($t('common.addSuccess'))
    }
    dialogVisible.value = false
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || $t('common.operateFailed'))
  } finally {
    submitting.value = false
  }
}

// 切换状态
const handleToggleStatus = async (row: any) => {
  try {
    await request.post({
      url: `/api/admin/payment-gateways/${row.id}/status`
    })
    ElMessage.success($t('common.updateSuccess'))
  } catch (error) {
    row.is_enabled = !row.is_enabled
    ElMessage.error($t('common.updateFailed'))
  }
}

// 测试接口
const handleTest = async (row: any) => {
  try {
    await request.post({
      url: `/api/admin/payment-gateways/${row.id}/test`
    })
    ElMessage.success($t('paymentGateways.testSuccess'))
  } catch (error: any) {
    ElMessage.error(error.message || $t('paymentGateways.testFailed'))
  }
}

// 删除
const handleDelete = async (row: any) => {
  try {
    await request.del({
      url: `/api/admin/payment-gateways/${row.id}`
    })
    ElMessage.success($t('common.deleteSuccess'))
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || $t('common.deleteFailed'))
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped>
.payment-gateways-page {
  padding: 20px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
}

.header-left h2 {
  margin: 0 0 8px 0;
  font-size: 20px;
  font-weight: 600;
}

.subtitle {
  color: #909399;
  font-size: 14px;
}

.table-card {
  background: #fff;
}

.payment-icon {
  width: 32px;
  height: 32px;
  object-fit: contain;
}

.payment-title {
  font-weight: 500;
}

.ml-1 {
  margin-left: 4px;
}

.ml-2 {
  margin-left: 8px;
}

.text-muted {
  color: #909399;
}

.pagination-wrap {
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
}

.form-hint {
  color: #909399;
  font-size: 12px;
  margin-top: 4px;
}

.amount-range {
  display: flex;
  align-items: center;
  gap: 8px;
}

.range-sep {
  color: #909399;
}
</style>
