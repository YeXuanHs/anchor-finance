<template>
  <div class="product-transfer-page">
    <div class="page-header">
      <h1 class="page-title">产品转移</h1>
      <el-button type="primary" @click="showTransferDialog = true">
        <el-icon><Promotion /></el-icon>发起转移
      </el-button>
    </div>

    <!-- 统计卡片 -->
    <div class="summary-grid">
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#409eff"><Top /></el-icon>
          <div class="summary-info">
            <span class="summary-value">{{ stats.sent }}</span>
            <span class="summary-label">已发出</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#e6a23c"><Bottom /></el-icon>
          <div class="summary-info">
            <span class="summary-value">{{ stats.received }}</span>
            <span class="summary-label">已收到</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#67c23a"><CircleCheck /></el-icon>
          <div class="summary-info">
            <span class="summary-value">{{ stats.completed }}</span>
            <span class="summary-label">已完成</span>
          </div>
        </div>
      </el-card>
      <el-card shadow="never" class="summary-card">
        <div class="summary-inner">
          <el-icon :size="28" color="#909399"><Clock /></el-icon>
          <div class="summary-info">
            <span class="summary-value">{{ stats.pending }}</span>
            <span class="summary-label">待确认</span>
          </div>
        </div>
      </el-card>
    </div>

    <!-- Tab 切换 -->
    <el-card shadow="never" class="main-card">
      <el-tabs v-model="activeTab">
        <!-- 发出的转移 -->
        <el-tab-pane label="发出的转移" name="sent">
          <el-table :data="sentList" stripe style="width: 100%" v-loading="loading" empty-text="暂无发出的转移">
            <el-table-column prop="id" label="转移号" width="140">
              <template #default="{ row }">
                <span class="mono-text">{{ row.id }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="productName" label="产品名称" min-width="160" show-overflow-tooltip />
            <el-table-column prop="targetUser" label="目标用户" width="140" />
            <el-table-column prop="transferCode" label="转移码" width="120">
              <template #default="{ row }">
                <span class="mono-text">{{ row.transferCode }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="createTime" label="发起时间" width="120" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small" effect="light" round>
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="120" fixed="right">
              <template #default="{ row }">
                <el-button v-if="row.status === 'pending'" type="danger" size="small" link @click="cancelTransfer(row)">
                  撤回
                </el-button>
                <el-button type="primary" size="small" link @click="viewDetail(row)">
                  详情
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 收到的转移 -->
        <el-tab-pane label="收到的转移" name="received">
          <el-table :data="receivedList" stripe style="width: 100%" v-loading="loading" empty-text="暂无收到的转移">
            <el-table-column prop="id" label="转移号" width="140">
              <template #default="{ row }">
                <span class="mono-text">{{ row.id }}</span>
              </template>
            </el-table-column>
            <el-table-column prop="productName" label="产品名称" min-width="160" show-overflow-tooltip />
            <el-table-column prop="fromUser" label="发起用户" width="140" />
            <el-table-column prop="createTime" label="发起时间" width="120" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small" effect="light" round>
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="160" fixed="right">
              <template #default="{ row }">
                <template v-if="row.status === 'pending'">
                  <el-button type="primary" size="small" @click="acceptTransfer(row)">接受</el-button>
                  <el-button type="danger" size="small" plain @click="rejectTransfer(row)">拒绝</el-button>
                </template>
                <el-button v-else type="primary" size="small" link @click="viewDetail(row)">详情</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <!-- 接受转移（通过转移码） -->
        <el-tab-pane label="接受转移" name="accept">
          <div class="accept-section">
            <el-card shadow="never" class="accept-card">
              <el-icon :size="48" color="#409eff"><Promotion /></el-icon>
              <h3>通过转移码接受产品</h3>
              <p>输入对方提供的转移码来接受产品转移</p>
              <el-input
                v-model="acceptCode"
                placeholder="请输入转移码"
                size="large"
                class="accept-input"
                clearable
              >
                <template #prefix>
                  <el-icon><Key /></el-icon>
                </template>
              </el-input>
              <el-button type="primary" size="large" @click="acceptByCode" :loading="acceptLoading" :disabled="!acceptCode">
                接受转移
              </el-button>
            </el-card>
          </div>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 发起转移对话框 -->
    <el-dialog v-model="showTransferDialog" title="发起产品转移" width="560px" destroy-on-close>
      <el-form ref="formRef" :model="transferForm" :rules="rules" label-width="100px">
        <el-form-item label="选择产品" prop="productId">
          <el-select v-model="transferForm.productId" placeholder="请选择要转移的产品" style="width: 100%;">
            <el-option
              v-for="product in myProducts"
              :key="product.id"
              :label="`${product.name} (${product.ip || product.domain || '-'})`"
              :value="product.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="目标用户" prop="targetUser">
          <el-input v-model="transferForm.targetUser" placeholder="请输入目标用户名或邮箱" />
        </el-form-item>
        <el-form-item label="转移备注">
          <el-input v-model="transferForm.remark" type="textarea" :rows="3" placeholder="转移备注（选填）" maxlength="200" show-word-limit />
        </el-form-item>
        <el-alert type="warning" :closable="false" show-icon title="注意事项">
          <template #default>
            <ul style="margin: 0; padding-left: 16px;">
              <li>转移成功后，产品将归属于目标用户</li>
              <li>转移期间产品不会中断服务</li>
              <li>转移后相关账单也将一并转移</li>
            </ul>
          </template>
        </el-alert>
      </el-form>
      <template #footer>
        <el-button @click="showTransferDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="submitTransfer">确认转移</el-button>
      </template>
    </el-dialog>

    <!-- 详情对话框 -->
    <el-dialog v-model="showDetailDialog" title="转移详情" width="560px">
      <el-descriptions :column="1" border v-if="currentTransfer">
        <el-descriptions-item label="转移号">{{ currentTransfer.id }}</el-descriptions-item>
        <el-descriptions-item label="产品名称">{{ currentTransfer.productName }}</el-descriptions-item>
        <el-descriptions-item label="发起用户">{{ currentTransfer.fromUser || '-' }}</el-descriptions-item>
        <el-descriptions-item label="目标用户">{{ currentTransfer.targetUser || '-' }}</el-descriptions-item>
        <el-descriptions-item label="转移码">
          <span class="mono-text">{{ currentTransfer.transferCode || '-' }}</span>
        </el-descriptions-item>
        <el-descriptions-item label="发起时间">{{ currentTransfer.createTime }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentTransfer.status)" size="small" effect="light" round>
            {{ getStatusText(currentTransfer.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="备注">{{ currentTransfer.remark || '-' }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="showDetailDialog = false">关闭</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import {
  Promotion, Top, Bottom, CircleCheck, Clock, Key
} from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import request from '@/utils/request'

interface Transfer {
  id: string
  productName: string
  productId: number
  fromUser?: string
  targetUser?: string
  transferCode?: string
  createTime: string
  status: 'pending' | 'accepted' | 'rejected' | 'cancelled'
  remark?: string
}

interface Product {
  id: number
  name: string
  ip?: string
  domain?: string
}

const activeTab = ref('sent')
const loading = ref(false)
const submitLoading = ref(false)
const acceptLoading = ref(false)
const showTransferDialog = ref(false)
const showDetailDialog = ref(false)
const acceptCode = ref('')
const currentTransfer = ref<Transfer | null>(null)
const formRef = ref<FormInstance>()

const stats = reactive({
  sent: 0,
  received: 0,
  completed: 0,
  pending: 0
})

const sentList = ref<Transfer[]>([])

const receivedList = ref<Transfer[]>([])

const myProducts = ref<Product[]>([])

const transferForm = reactive({
  productId: '',
  targetUser: '',
  remark: ''
})

const rules: FormRules = {
  productId: [{ required: true, message: '请选择要转移的产品', trigger: 'change' }],
  targetUser: [{ required: true, message: '请输入目标用户', trigger: 'blur' }]
}

function getStatusType(status: string) {
  const map: Record<string, 'success' | 'warning' | 'info' | 'danger'> = {
    pending: 'warning', accepted: 'success', rejected: 'danger', cancelled: 'info'
  }
  return map[status] || 'info'
}

function getStatusText(status: string) {
  const map: Record<string, string> = {
    pending: '待确认', accepted: '已接受', rejected: '已拒绝', cancelled: '已撤回'
  }
  return map[status] || status
}

async function fetchTransfers() {
  loading.value = true
  try {
    const { data } = await request.get('/api/v2/product-diverts')
    if (data?.data) {
      sentList.value = data.data.sent || sentList.value
      receivedList.value = data.data.received || receivedList.value
      if (data.data.stats) {
        Object.assign(stats, data.data.stats)
      }
    }
  } catch (e) {
    console.error('Failed to fetch transfers:', e)
  } finally {
    loading.value = false
  }
}

async function fetchMyProducts() {
  try {
    const { data } = await request.get('/api/v2/hosts')
    if (data?.data?.list) {
      myProducts.value = data.data.list.map((item: any) => ({
        id: item.id,
        name: item.product_name || item.name,
        ip: item.dedicated_ip,
        domain: item.domain
      }))
    }
  } catch (e) {
    console.error('Failed to fetch products:', e)
  }
}

async function submitTransfer() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        await request.post('/api/v2/product-diverts', {
          host_id: transferForm.productId,
          target_user: transferForm.targetUser,
          remark: transferForm.remark
        })
        ElMessage.success('转移请求已发送')
        showTransferDialog.value = false
        transferForm.productId = ''
        transferForm.targetUser = ''
        transferForm.remark = ''
        fetchTransfers()
      } catch (e: any) {
        ElMessage.error(e.message || '发起转移失败')
      } finally {
        submitLoading.value = false
      }
    }
  })
}

async function cancelTransfer(row: Transfer) {
  try {
    await ElMessageBox.confirm('确定要撤回该转移请求吗？', '撤回转移', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.post(`/api/v2/product-diverts/${row.id}/cancel`)
    row.status = 'cancelled'
    ElMessage.success('已撤回')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '操作失败')
    }
  }
}

async function acceptTransfer(row: Transfer) {
  try {
    await ElMessageBox.confirm('确定接受该产品转移吗？', '接受转移', {
      confirmButtonText: '确定接受',
      cancelButtonText: '取消',
      type: 'info'
    })
    await request.post(`/api/v2/product-diverts/${row.id}/accept`)
    row.status = 'accepted'
    ElMessage.success('已接受转移')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '操作失败')
    }
  }
}

async function rejectTransfer(row: Transfer) {
  try {
    await ElMessageBox.confirm('确定拒绝该产品转移吗？', '拒绝转移', {
      confirmButtonText: '确定拒绝',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await request.post(`/api/v2/product-diverts/${row.id}/reject`)
    row.status = 'rejected'
    ElMessage.success('已拒绝转移')
  } catch (e: any) {
    if (e !== 'cancel') {
      ElMessage.error(e.message || '操作失败')
    }
  }
}

async function acceptByCode() {
  if (!acceptCode.value) return
  acceptLoading.value = true
  try {
    await request.post('/api/v2/product-diverts/accept-by-code', {
      code: acceptCode.value
    })
    ElMessage.success('已接受产品转移')
    acceptCode.value = ''
    fetchTransfers()
  } catch (e: any) {
    ElMessage.error(e.message || '转移码无效或已过期')
  } finally {
    acceptLoading.value = false
  }
}

function viewDetail(row: Transfer) {
  currentTransfer.value = row
  showDetailDialog.value = true
}

onMounted(() => {
  fetchTransfers()
  fetchMyProducts()
})
</script>

<style scoped lang="scss">
.product-transfer-page {
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

.summary-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 16px;
}

.summary-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;
  background: #fff;
}

.summary-card :deep(.el-card__body) {
  padding: 20px;
}

.summary-inner {
  display: flex;
  align-items: center;
  gap: 16px;
}

.summary-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.summary-value {
  font-size: 24px;
  font-weight: 700;
  color: #303133;
}

.summary-label {
  font-size: 13px;
  color: #909399;
}

.main-card {
  border-radius: 12px;
  border: 1px solid #e8ecf1;

  :deep(.el-card__body) {
    padding: 0 20px 20px;
  }

  :deep(.el-tabs__header) {
    margin: 0;
    padding: 0 0 16px;
  }
}

.mono-text {
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 13px;
  color: #606266;
}

.accept-section {
  display: flex;
  justify-content: center;
  padding: 40px 0;
}

.accept-card {
  text-align: center;
  max-width: 480px;
  width: 100%;
  border-radius: 12px;

  h3 {
    font-size: 18px;
    font-weight: 600;
    color: #303133;
    margin: 16px 0 8px;
  }

  p {
    font-size: 14px;
    color: #909399;
    margin: 0 0 24px;
  }
}

.accept-input {
  margin-bottom: 16px;

  :deep(.el-input__inner) {
    text-align: center;
    font-family: 'Monaco', 'Menlo', monospace;
    font-size: 18px;
    letter-spacing: 2px;
  }
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: 12px;
    align-items: flex-start;
  }

  .summary-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
