<template>
  <div class="combine-billing">
    <div class="page-header">
      <h2>合并账单</h2>
      <p>将多个未付账单合并为一笔支付</p>
    </div>
    <el-card>
      <el-alert title="选择要合并的未付账单，合并后将生成一笔新的订单进行统一支付" type="info" :closable="false" style="margin-bottom:20px" />
      <el-table :data="unpaidInvoices" v-loading="loading" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column prop="invoice_no" label="账单编号" width="180" />
        <el-table-column prop="description" label="描述" />
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">¥{{ row.amount?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="due_date" label="到期日" width="120" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'overdue' ? 'danger' : 'warning'" size="small">{{ row.status === 'overdue' ? '已逾期' : '待支付' }}</el-tag>
          </template>
        </el-table-column>
      </el-table>
      <div class="combine-footer" v-if="selected.length > 0">
        <div class="total">已选 {{ selected.length }} 笔，合计：<span class="amount">¥{{ totalAmount.toFixed(2) }}</span></div>
        <el-button type="primary" size="large" @click="handleCombine" :loading="submitting">合并支付</el-button>
      </div>
    </el-card>
  </div>
</template>
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const loading = ref(false)
const submitting = ref(false)
const unpaidInvoices = ref<any[]>([])
const selected = ref<any[]>([])

const totalAmount = computed(() => selected.value.reduce((sum, item) => sum + (item.amount || 0), 0))

const handleSelectionChange = (val: any[]) => { selected.value = val }

const fetchInvoices = async () => {
  loading.value = true
  try {
    const { data } = await request.get('/api/v2/user/invoices', { params: { status: 'unpaid', page_size: 100 } })
    unpaidInvoices.value = data.data || []
  } catch {} finally { loading.value = false }
}

const handleCombine = async () => {
  submitting.value = true
  try {
    const ids = selected.value.map(i => i.id)
    const { data } = await request.post('/api/v2/user/invoices/combine', { invoice_ids: ids })
    ElMessage.success('合并成功')
    router.push(`/user/orders/${data.data.order_id}`)
  } catch { ElMessage.error('合并失败') } finally { submitting.value = false }
}

onMounted(fetchInvoices)
</script>
<style scoped lang="scss">
.combine-billing { .page-header { margin-bottom: 24px; h2 { font-size: 20px; color: #1a365d; } p { color: #6b7280; margin-top: 4px; } } }
.combine-footer { display: flex; justify-content: space-between; align-items: center; margin-top: 20px; padding: 16px; background: #f8fafc; border-radius: 8px;
  .total { font-size: 16px; color: #374151; .amount { font-size: 24px; font-weight: 700; color: #2563eb; } }
}
</style>
