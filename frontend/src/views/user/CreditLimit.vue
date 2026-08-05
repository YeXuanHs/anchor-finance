<template>
  <div class="credit-limit-page">
    <el-card class="page-card">
      <template #header>
        <div class="card-header">
          <span>信用额度</span>
          <el-button type="primary" @click="showApplyDialog">申请提额</el-button>
        </div>
      </template>

      <div class="credit-overview">
        <el-row :gutter="20">
          <el-col :span="8">
            <div class="credit-card">
              <div class="credit-label">总额度</div>
              <div class="credit-value">¥{{ creditInfo.total?.toFixed(2) || '0.00' }}</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="credit-card">
              <div class="credit-label">已用额度</div>
              <div class="credit-value used">¥{{ creditInfo.used?.toFixed(2) || '0.00' }}</div>
            </div>
          </el-col>
          <el-col :span="8">
            <div class="credit-card">
              <div class="credit-label">可用额度</div>
              <div class="credit-value available">¥{{ creditInfo.available?.toFixed(2) || '0.00' }}</div>
            </div>
          </el-col>
        </el-row>

        <div class="credit-progress">
          <div class="progress-label">
            <span>额度使用率</span>
            <span>{{ usagePercent }}%</span>
          </div>
          <el-progress
            :percentage="usagePercent"
            :color="usagePercent > 80 ? '#f56c6c' : usagePercent > 60 ? '#e6a23c' : '#67c23a'"
            :stroke-width="10"
          />
        </div>
      </div>

      <el-divider />

      <div class="credit-actions">
        <h3>额度管理</h3>
        <div class="action-cards">
          <div class="action-card" @click="showApplyDialog">
            <el-icon :size="32" color="#409eff"><Document /></el-icon>
            <div class="action-text">
              <div class="action-title">申请提额</div>
              <div class="action-desc">提交资料提升信用额度</div>
            </div>
          </div>
          <div class="action-card" @click="showUsageDialog">
            <el-icon :size="32" color="#67c23a"><DataLine /></el-icon>
            <div class="action-text">
              <div class="action-title">额度使用明细</div>
              <div class="action-desc">查看额度使用详情</div>
            </div>
          </div>
          <div class="action-card" @click="showRepayDialog">
            <el-icon :size="32" color="#e6a23c"><Wallet /></el-icon>
            <div class="action-text">
              <div class="action-title">还款</div>
              <div class="action-desc">归还已使用的信用额度</div>
            </div>
          </div>
        </div>
      </div>

      <el-divider />

      <h3>额度使用记录</h3>
      <el-table :data="creditLogs" style="width: 100%">
        <el-table-column prop="created_at" label="时间" width="180" />
        <el-table-column prop="type" label="类型">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)">{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="amount" label="金额">
          <template #default="{ row }">
            <span :class="row.type === 'increase' || row.type === 'repay' ? 'text-success' : 'text-danger'">
              {{ row.type === 'increase' || row.type === 'repay' ? '+' : '-' }}¥{{ row.amount?.toFixed(2) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="balance_after" label="操作后余额">
          <template #default="{ row }">
            <span>¥{{ row.balance_after?.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="reason" label="原因" />
        <el-table-column prop="order_no" label="关联订单" />
      </el-table>

      <div class="pagination-wrap">
        <el-pagination
          v-model:current-page="currentPage"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next"
          @size-change="loadCreditLogs"
          @current-change="loadCreditLogs"
        />
      </div>
    </el-card>

    <!-- 申请提额对话框 -->
    <el-dialog v-model="applyDialogVisible" title="申请提额" width="520px">
      <el-form :model="applyForm" label-width="100px">
        <el-form-item label="当前额度">
          <span class="current-amount">¥{{ creditInfo.total?.toFixed(2) }}</span>
        </el-form-item>
        <el-form-item label="申请额度" required>
          <el-input-number
            v-model="applyForm.amount"
            :min="creditInfo.total || 0"
            :max="100000"
            :step="1000"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="申请原因" required>
          <el-input
            v-model="applyForm.reason"
            type="textarea"
            :rows="4"
            placeholder="请说明申请提额的原因"
          />
        </el-form-item>
        <el-form-item label="补充材料">
          <el-upload
            action="#"
            :auto-upload="false"
            :on-change="handleFileChange"
            :file-list="applyForm.files"
          >
            <el-button type="primary" plain>选择文件</el-button>
            <template #tip>
              <div class="el-upload__tip">支持上传营业执照、资产证明等材料</div>
            </template>
          </el-upload>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="applyDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitApply">提交申请</el-button>
      </template>
    </el-dialog>

    <!-- 还款对话框 -->
    <el-dialog v-model="repayDialogVisible" title="信用额度还款" width="420px">
      <el-form :model="repayForm" label-width="80px">
        <el-form-item label="已用额度">
          <span class="current-amount text-danger">¥{{ creditInfo.used?.toFixed(2) }}</span>
        </el-form-item>
        <el-form-item label="还款金额" required>
          <el-input-number
            v-model="repayForm.amount"
            :min="1"
            :max="creditInfo.used || 0"
            :precision="2"
            style="width: 100%"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="repayDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitRepay">确认还款</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Document, DataLine, Wallet } from '@element-plus/icons-vue'
import request from '@/utils/request'

const creditInfo = ref({
  total: 10000,
  used: 3500,
  available: 6500
})

const creditLogs = ref<any[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const submitting = ref(false)

// 申请提额
const applyDialogVisible = ref(false)
const applyForm = ref({
  amount: 15000,
  reason: '',
  files: [] as any[]
})

// 还款
const repayDialogVisible = ref(false)
const repayForm = ref({
  amount: 0
})

const usagePercent = computed(() => {
  if (!creditInfo.value.total) return 0
  return Math.round((creditInfo.value.used / creditInfo.value.total) * 100)
})

const getTypeTag = (type: string) => {
  const map: Record<string, string> = {
    increase: 'success',
    decrease: 'danger',
    use: 'warning',
    repay: 'success'
  }
  return map[type] || 'info'
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    increase: '额度增加',
    decrease: '额度减少',
    use: '额度使用',
    repay: '额度归还'
  }
  return map[type] || type
}

const showApplyDialog = () => {
  applyForm.value = {
    amount: (creditInfo.value.total || 0) + 5000,
    reason: '',
    files: []
  }
  applyDialogVisible.value = true
}

const showUsageDialog = () => {
  const el = document.querySelector('.credit-limit-page h3:last-of-type')
  if (el) {
    el.scrollIntoView({ behavior: 'smooth', block: 'start' })
  }
}

const showRepayDialog = () => {
  repayForm.value = { amount: creditInfo.value.used || 0 }
  repayDialogVisible.value = true
}

const handleFileChange = (file: any) => {
  applyForm.value.files.push(file)
}

const submitApply = async () => {
  if (!applyForm.value.reason) {
    ElMessage.warning('请填写申请原因')
    return
  }
  submitting.value = true
  try {
    await request.post('/api/v1/credit/apply', {
      amount: applyForm.value.amount,
      reason: applyForm.value.reason
    })
    ElMessage.success('提额申请已提交')
    applyDialogVisible.value = false
  } catch (e: any) {
    ElMessage.error(e.response?.data?.message || '申请提交失败')
  } finally {
    submitting.value = false
  }
}

const submitRepay = async () => {
  if (!repayForm.value.amount || repayForm.value.amount <= 0) {
    ElMessage.warning('请输入有效的还款金额')
    return
  }
  submitting.value = true
  try {
    await request.post('/api/v1/credit/repay', { amount: repayForm.value.amount })
    creditInfo.value.used -= repayForm.value.amount
    creditInfo.value.available += repayForm.value.amount
    ElMessage.success('还款成功')
    repayDialogVisible.value = false
    loadCreditLogs()
  } catch (e: any) {
    ElMessage.error(e.response?.data?.message || '还款失败')
  } finally {
    submitting.value = false
  }
}

const loadCreditLogs = async () => {
  try {
    const { data } = await request.get('/api/v1/credit/logs', {
      params: { page: currentPage.value, page_size: pageSize.value }
    })
    creditLogs.value = data?.data?.list || data?.data?.items || []
    total.value = data?.data?.total || 0
  } catch {
    creditLogs.value = []
    total.value = 0
  }
}

onMounted(async () => {
  try {
    const { data } = await request.get('/api/v1/credit')
    if (data?.data) {
      creditInfo.value = {
        total: data.data.limit || data.data.total || 0,
        used: data.data.used || 0,
        available: (data.data.limit || data.data.total || 0) - (data.data.used || 0)
      }
    }
  } catch {}
  loadCreditLogs()
})
</script>

<style scoped lang="scss">
.credit-limit-page {
  .credit-overview {
    margin-bottom: 20px;
  }

  .credit-card {
    background: #f5f7fa;
    padding: 20px;
    border-radius: 8px;
    text-align: center;

    .credit-label {
      color: #909399;
      font-size: 14px;
      margin-bottom: 8px;
    }

    .credit-value {
      font-size: 24px;
      font-weight: bold;
      color: #303133;

      &.used {
        color: #e6a23c;
      }

      &.available {
        color: #67c23a;
      }
    }
  }

  .credit-progress {
    margin-top: 20px;
    padding: 16px;
    background: #f5f7fa;
    border-radius: 8px;

    .progress-label {
      display: flex;
      justify-content: space-between;
      margin-bottom: 8px;
      font-size: 14px;
      color: #606266;
    }
  }

  .credit-actions {
    h3 {
      font-size: 16px;
      font-weight: 600;
      margin: 0 0 16px 0;
    }

    .action-cards {
      display: grid;
      grid-template-columns: repeat(3, 1fr);
      gap: 16px;
    }

    .action-card {
      display: flex;
      align-items: center;
      gap: 16px;
      padding: 20px;
      border: 1px solid #ebeef5;
      border-radius: 8px;
      cursor: pointer;
      transition: all 0.3s;

      &:hover {
        border-color: #409eff;
        box-shadow: 0 2px 12px rgba(64, 158, 255, 0.1);
      }

      .action-text {
        .action-title {
          font-weight: 600;
          margin-bottom: 4px;
        }

        .action-desc {
          font-size: 13px;
          color: #909399;
        }
      }
    }
  }

  .text-success {
    color: #67c23a;
    font-weight: bold;
  }

  .text-danger {
    color: #f56c6c;
    font-weight: bold;
  }

  .current-amount {
    font-size: 18px;
    font-weight: bold;
  }

  .pagination-wrap {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
