<template>
  <div class="payment-gateways-page">
    <!-- 顶部操作栏 -->
    <div class="page-header">
      <div class="header-left">
        <h2>支付方式管理</h2>
        <span class="subtitle">管理用户可用的支付方式，每种方式对应一个支付接口+支付类型的组合</span>
      </div>
      <el-button type="primary" @click="handleAdd">
        <el-icon><Plus /></el-icon>
        新增支付方式
      </el-button>
    </div>

    <!-- 支付方式列表 -->
    <el-card class="table-card">
      <el-table :data="tableData" v-loading="loading" stripe>
        <el-table-column prop="id" label="ID" width="60" />
        <el-table-column label="图标" width="60">
          <template #default="{ row }">
            <img :src="getIconUrl(row.code)" class="payment-icon" />
          </template>
        </el-table-column>
        <el-table-column prop="title" label="显示名称" min-width="150">
          <template #default="{ row }">
            <span class="payment-title">{{ row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="标识" width="150" />
        <el-table-column label="接口/类型" width="180">
          <template #default="{ row }">
            <el-tag size="small" type="info">{{ getGatewayLabel(row.gateway) }}</el-tag>
            <el-tag size="small" class="ml-1">{{ getCodeLabel(row.code) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="手续费" width="100">
          <template #default="{ row }">
            {{ row.fee_rate > 0 ? (row.fee_rate * 100).toFixed(2) + '%' : '无' }}
          </template>
        </el-table-column>
        <el-table-column label="金额限制" width="150">
          <template #default="{ row }">
            <template v-if="row.min_amount > 0 || row.max_amount > 0">
              {{ row.min_amount || 0 }} - {{ row.max_amount || '不限' }}
            </template>
            <span v-else class="text-muted">不限</span>
          </template>
        </el-table-column>
        <el-table-column prop="sort_order" label="排序" width="80" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-switch
              v-model="row.is_enabled"
              @change="handleToggleStatus(row)"
              :active-text="''"
              :inactive-text="''"
            />
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link size="small" @click="handleEdit(row)">编辑</el-button>
            <el-button type="info" link size="small" @click="handleTest(row)">测试</el-button>
            <el-popconfirm title="确定删除该支付方式？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button type="danger" link size="small">删除</el-button>
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
      :title="isEdit ? '编辑支付方式' : '新增支付方式'"
      width="680px"
      :close-on-click-modal="false"
    >
      <el-form ref="formRef" :model="formData" :rules="formRules" label-width="100px">
        <!-- 接口选择 -->
        <el-form-item label="支付接口" prop="gateway">
          <el-select v-model="formData.gateway" placeholder="选择支付接口" @change="handleGatewayChange" :disabled="isEdit">
            <el-option
              v-for="(label, value) in supportedInfo.gateways"
              :key="value"
              :label="label"
              :value="value"
            />
          </el-select>
          <div class="form-hint">处理支付的后端接口，创建后不可修改</div>
        </el-form-item>

        <!-- 支付类型选择 -->
        <el-form-item label="支付类型" prop="code">
          <el-select v-model="formData.code" placeholder="选择支付类型" @change="handleCodeChange" :disabled="isEdit">
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
          <div class="form-hint">用户看到的支付方式类型，创建后不可修改</div>
        </el-form-item>

        <!-- 标识名 -->
        <el-form-item label="唯一标识" prop="name">
          <el-input v-model="formData.name" placeholder="如 EpayAlipay" :disabled="isEdit" />
          <div class="form-hint">系统内部标识，创建后不可修改</div>
        </el-form-item>

        <!-- 显示名称 -->
        <el-form-item label="显示名称" prop="title">
          <el-input v-model="formData.title" placeholder="如 支付宝-易支付" />
          <div class="form-hint">用户看到的名称，格式建议：支付类型-接口名称</div>
        </el-form-item>

        <!-- 接口配置 -->
        <el-divider content-position="left">接口配置</el-divider>

        <!-- 易支付配置 -->
        <template v-if="formData.gateway === 'epay'">
          <el-form-item label="商户ID" prop="config.pid">
            <el-input v-model="configForm.pid" placeholder="易支付商户ID" />
          </el-form-item>
          <el-form-item label="商户密钥" prop="config.key">
            <el-input v-model="configForm.key" placeholder="易支付商户密钥" show-password />
          </el-form-item>
          <el-form-item label="接口地址" prop="config.api_url">
            <el-input v-model="configForm.api_url" placeholder="https://pay.example.com" />
          </el-form-item>
        </template>

        <!-- 迅虎支付配置 -->
        <template v-if="formData.gateway === 'xunhupay'">
          <el-form-item label="应用ID" prop="config.api_key">
            <el-input v-model="configForm.api_key" placeholder="迅虎应用ID (api_key)" />
          </el-form-item>
          <el-form-item label="应用密钥" prop="config.secret_key">
            <el-input v-model="configForm.secret_key" placeholder="迅虎应用密钥 (secret_key)" show-password />
          </el-form-item>
          <el-form-item label="接口地址" prop="config.gateway">
            <el-input v-model="configForm.gateway" placeholder="https://api.xunhupay.com" />
          </el-form-item>
        </template>

        <!-- 支付宝官方配置 -->
        <template v-if="formData.gateway === 'alipay'">
          <el-form-item label="应用ID" prop="config.app_id">
            <el-input v-model="configForm.app_id" placeholder="支付宝开放平台应用ID" />
          </el-form-item>
          <el-form-item label="应用私钥" prop="config.private_key">
            <el-input v-model="configForm.private_key" type="textarea" :rows="4" placeholder="应用私钥（RSA2）" />
          </el-form-item>
          <el-form-item label="支付宝公钥" prop="config.alipay_public_key">
            <el-input v-model="configForm.alipay_public_key" type="textarea" :rows="4" placeholder="支付宝公钥" />
          </el-form-item>
          <el-form-item label="支付产品">
            <el-checkbox v-model="configForm.product_pc">PC网页支付</el-checkbox>
            <el-checkbox v-model="configForm.product_wap">手机网页支付</el-checkbox>
            <el-checkbox v-model="configForm.product_qr">扫码支付</el-checkbox>
          </el-form-item>
        </template>

        <!-- 微信支付官方配置 -->
        <template v-if="formData.gateway === 'wxpay'">
          <el-form-item label="公众号ID" prop="config.app_id">
            <el-input v-model="configForm.app_id" placeholder="微信公众号/小程序AppID" />
          </el-form-item>
          <el-form-item label="商户号" prop="config.mch_id">
            <el-input v-model="configForm.mch_id" placeholder="微信支付商户号" />
          </el-form-item>
          <el-form-item label="API密钥" prop="config.api_key">
            <el-input v-model="configForm.api_key" placeholder="微信支付API密钥" show-password />
          </el-form-item>
          <el-form-item label="应用密钥" prop="config.app_secret">
            <el-input v-model="configForm.app_secret" placeholder="公众号/小程序AppSecret" show-password />
          </el-form-item>
          <el-form-item label="支付产品">
            <el-checkbox v-model="configForm.product_native">扫码支付</el-checkbox>
            <el-checkbox v-model="configForm.product_jsapi">JSAPI支付</el-checkbox>
            <el-checkbox v-model="configForm.product_wap">H5支付</el-checkbox>
          </el-form-item>
        </template>

        <!-- 余额支付无配置 -->
        <template v-if="formData.gateway === 'balance'">
          <el-alert title="余额支付无需额外配置" type="info" :closable="false" show-icon />
        </template>

        <!-- 其他设置 -->
        <el-divider content-position="left">其他设置</el-divider>

        <el-form-item label="手续费率">
          <el-input-number v-model="formData.fee_rate" :min="0" :max="1" :step="0.01" :precision="4" />
          <span class="ml-2">%</span>
          <div class="form-hint">如 0.02 表示 2% 手续费</div>
        </el-form-item>

        <el-form-item label="金额限制">
          <div class="amount-range">
            <el-input-number v-model="formData.min_amount" :min="0" :precision="2" placeholder="最低" />
            <span class="range-sep">至</span>
            <el-input-number v-model="formData.max_amount" :min="0" :precision="2" placeholder="最高" />
          </div>
          <div class="form-hint">为 0 表示不限制</div>
        </el-form-item>

        <el-form-item label="排序">
          <el-input-number v-model="formData.sort_order" :min="0" />
          <div class="form-hint">数值越小越靠前</div>
        </el-form-item>

        <el-form-item label="启用状态">
          <el-switch v-model="formData.is_enabled" active-text="启用" inactive-text="禁用" />
        </el-form-item>
      </el-form>

      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleSubmit" :loading="submitting">
          {{ isEdit ? '保存' : '创建' }}
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
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
  name: [{ required: true, message: '请输入唯一标识', trigger: 'blur' }],
  title: [{ required: true, message: '请输入显示名称', trigger: 'blur' }],
  gateway: [{ required: true, message: '请选择支付接口', trigger: 'change' }],
  code: [{ required: true, message: '请选择支付类型', trigger: 'change' }]
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
const handleEdit = (row: PaymentGateway) => {
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
      ElMessage.success('更新成功')
    } else {
      await request.post({
        url: '/api/admin/payment-gateways',
        data
      })
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

// 切换状态
const handleToggleStatus = async (row: PaymentGateway) => {
  try {
    await request.post({
      url: `/api/admin/payment-gateways/${row.id}/status`
    })
    ElMessage.success('状态已更新')
  } catch (error) {
    row.is_enabled = !row.is_enabled
    ElMessage.error('更新状态失败')
  }
}

// 测试接口
const handleTest = async (row: PaymentGateway) => {
  try {
    await request.post({
      url: `/api/admin/payment-gateways/${row.id}/test`
    })
    ElMessage.success('接口测试通过')
  } catch (error: any) {
    ElMessage.error(error.message || '接口测试失败')
  }
}

// 删除
const handleDelete = async (row: PaymentGateway) => {
  try {
    await request.delete({
      url: `/api/admin/payment-gateways/${row.id}`
    })
    ElMessage.success('删除成功')
    fetchData()
  } catch (error: any) {
    ElMessage.error(error.message || '删除失败')
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
